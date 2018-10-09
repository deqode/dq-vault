# PART 1 :- Setting up Vault using Docker

We are assuming that you have docker and docker-compose installed in your system. If not, set up your docker environment first or refer to the original [README.md](https://gitlab.com/arout/Vault/blob/master/README.md) for manual setting up of the vault.

Also, you need to have Postgres set up in your system.

All the required docker files are provided in this repo.

- `config.hcl` is where we define our Vault configurations.
- `docker-compose.yml` sets up our postgres and vault containers.
- `init.sql` inside postgres_docker folder is used to create a table in our postgres database where all our vault encrypted data will be stored.

## Setting Vault Configurations

- Edit the `config.hcl` file. Given below is a sample config.hcl file that we have provided:

  config.hcl

  ```
  storage "postgresql" {
    connection_url = "postgres://vault:vault@database:5432/databasename?sslmode=disable"
    table="vault_kv_store"
  }

  listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = 1
  }

  "api_addr" = "http://127.0.0.1:8200"
  "plugin_directory" ="/vault/plugins"
  "disable_mlock"=true
  ```

  - `api_addr` defines the access port of vault. All the API requests to vault will be done via this port.
  - `Storage` defines the backend-storage type that vault will use to store all the encrypted data. Since this backend-storage is not a part of the vault, we define the access port of PostgreSQL server and the table name. Change the role, password, databaseName, and tableName according to your postgres parameters. Note that we have disabled SSL for database requests.
  - In the `listener` part we have disabled TLS. If activated, TLS certificates and keys have to be provided here also. The example above listens on localhost port 8200 without TLS.
  - Lastly, we have defined the `plugin directory` where the vault looks for plugins.

  For more on vault configurations, you may refer to this [link](https://www.vaultproject.io/docs/configuration/index.html)

- Edit the `docker-compose.yml` file. You just need to change the Postgres environment variables according to your choice.

**Keep the Postgres environment variables same in both config.hcl and docker-compose.yml.**

- Edit the `init.sql` file if you want to change the table name.

**Keep the table-name same in both init.sql and config.hcl.**

## Steps

Run the following command to set-up the vault and Postgres containers:

```sh
  $ docker-compose up --build
```

The command does the following things:-

- Creates Postgres container with the provided database and the `vault-kv-store` table.
- Builds vault container.
- Links the Postgres container with vault container.
- The vault server is created with proper configurations according to config.hcl and listens to requests on the port as specified.

You should see the following response on the terminal:

```
vault_docker | ==> Vault server configuration:
vault_docker |
vault_docker |              Api Address: http://127.0.0.1:8200
vault_docker |                      Cgo: disabled
vault_docker |          Cluster Address: https://127.0.0.1:8201
vault_docker |               Listener 1: tcp (addr: "0.0.0.0:8200", cluster address: "0.0.0.0:8201", max_request_duration: "1m30s", max_request_size: "33554432", tls: "disabled")
vault_docker |                Log Level: (not set)
vault_docker |                    Mlock: supported: true, enabled: false
vault_docker |                  Storage: postgresql
vault_docker |                  Version: Vault v0.11.1
vault_docker |              Version Sha: 8575f8fedcf8f5a6eb2b4701cb527b99574b5286
vault_docker |
vault_docker | ==> Vault server started! Log data will stream in below:
vault_docker |
```

Note: If the above response is not displayed, then run the docker-compose again.

## Vault Setup

Now that the vault server is up and running on port 8200, it is actually in a sealed state, that is vault functionalities can't be accessed yet. To access vault we need to unseal it. First, open another terminal window and initialize vault by running the following commands:

```sh
  $ curl \
    --request PUT \
    --data '{ "secret_shares": 10, "secret_threshold": 5}' \
    http://127.0.0.1:8200/v1/sys/init
```

What the above command does is gives a set of 10 Shamir keys which have the capability to unseal vault and an initial root token. Here vault is initialized in such a way that any 5 keys out of 10 are enough to unseal vault. The root token is used to login into the vault. Only after logging in, you can start using vault. Store these in a safe place for later use.

The initial root token you receive in the response gives you admin access to the vault.

Note: You can set your own secret_shares and secret_threshold. Secret_shares is the total number of keys that will be generated and secret_threshold is the minimum number of keys required to unseal vault.

**Note: Use base64 encoded keys to unseal vault**

Start the unsealing process by running the command:

```sh
  $ curl \
    --request PUT \
    --data '{"key": "..."}' \
    http://127.0.0.1:8200/v1/sys/unseal
```

According to the above example, run this command 5 times using 5 different keys.

To check if the vault is unsealed or not, run the following command:

```sh
  curl http://127.0.0.1:8200/v1/sys/seal-status
```

Response:

```
  {
    "type":"shamir",
    "sealed":false, <------this
    "t":10,
    "n":5,
    "progress":5,
    "nonce":"",
    "version":"0.11.1"
  }
```

If you find the sealed attribute to be false, then your vault is unsealed.

## Enable Plugin

- If you previously have enabled this plugin, you need to disable it.

  ```sh
    $ curl -X DELETE \
    http://127.0.0.1:8200/v1/sys/mounts/api \
    -H 'content-type: application/json' \
    -H 'x-vault-token: ...'
  ```

- Register the plugin in vault's plugin catalog:

  ```sh
    curl -X PUT \
    http://127.0.0.1:8200/v1/sys/plugins/catalog/secrets-api \
    -H 'content-type: application/json' \
    -H 'x-vault-token: ...' \
    -d '{"sha_256": "8dc7e0f1df9e2e183a7579c7eb102ce40e8a2de44c5ab9378bb348b8dd332358","command": "vault_plugin"}'
  ```

  Provide your initial root token as the x-vault-token.

- Mount the secrets engine

  ```sh
    curl -X POST \
  http://127.0.0.1:8200/v1/sys/mounts/api \
  -H 'content-type: application/json' \
  -H 'x-vault-token: bd1c4051-cfc2-b6d4-6547-c0541b74d0bf' \
  -d '{"plugin_name":"secrets-api", "type":"plugin"}'
  ```

  Our plugin must be enabled now.

## Creating policies for application server

We need to define policies for the application server that will be using our vault. We don't want our application server to have complete root access of vault, rather, it should just have the capability to update our Vault API plugin that we just enabled. For that, we have provided an application.json file in the config folder.

To register this policy in vault, open terminal in the directory containing application.json and run the following command:

```sh
  $ curl \
  --header "X-Vault-Token: ..." \
  --request PUT \
  --data @application.json \
  http://127.0.0.1:8200/v1/sys/policy/application
```

Now we can use application policy to define access capabilities of anyone using vault. For more details refer to this [link](https://www.vaultproject.io/docs/concepts/policies.html)

## Enable userpass authentication method

We want our application server to login into vault using a particular `username` and `password` and should have access capabilities defined by the `application` policy we created earlier. In order to do this, we will be enabling userpass authentication method.

```sh
  $ cat payload.json
  {
    "type": "userpass"
  }

  $ curl \
  --header "X-Vault-Token: ..." \
  --request POST \
  --data @payload.json \
  http://127.0.0.1:8200/v1/sys/auth/userpass
```

We then create a username and password using which our application server will log in. We also attach the application policy in this command. The following command creates a user with username-"appserver" and password-"secret" with application policy attached:

```sh
  $ cat payload.json
  {
    "password": "secret",
    "policies": "application, default"
  }

  $ curl \
  --header "X-Vault-Token: ..." \
  --request POST \
  --data @payload.json \
  http://127.0.0.1:8200/v1/auth/userpass/users/appserver
```

We then give these credentials to the application server, who will use this username and password to login into the vault. Note that anyone logged in by this method will have capabilities defined by the application policy.

We can easily create multiple user login credentials for different application servers.

For information on plugin usage, follow this [link](https://gitlab.com/arout/Vault/blob/master/README.md#part-2-using-vault).
