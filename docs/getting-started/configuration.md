---
sidebar_position: 2
---

## Set up postgres

Assuming that you have PostgreSQL installed in your system, you need to create a table which will be used by Vault to
store it's encrypted data.

Once you are into PostgreSQL shell prompt, run the following commands to create a table:

```
  postgres=# create database vault;

  postgres=# \c vault

  vault=# CREATE TABLE vault_kv_store (
  parent_path TEXT COLLATE "C" NOT NULL,
  path        TEXT COLLATE "C",
  key         TEXT COLLATE "C",
  value       BYTEA,
  CONSTRAINT pkey PRIMARY KEY (path,key)
  );

  vault=# CREATE INDEX parent_path_idx ON vault_kv_store (parent_path);
```

## Populate config file

This is the step where we define the vault configurations. Vault supports .hcl files to write your configurations and
set-up vault accordingly.

Given below is an example of a config.hcl file:

**config.hcl**-

```
 disable_mlock = true
 ui            = true

 storage "postgresql" {
   connection_url = "postgres://role:password@localhost:5432/databaseName?sslmode=disable"
   table = "tableName"
 }

 
 listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = "true"
 }

 plugin_directory = "/$HOME/etc/vault/plugins/vault_plugin"
 api_addr = "http://127.0.0.1:8200"
```

**config.hcl** without postgres storage-

```
 disable_mlock = true
 ui            = true

 listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = "true"
 }

 storage "file" {
  path = "/$HOME/etc/tmp/vault-data"
 }

 plugin_directory = "/$HOME/etc/vault/plugins/vault_plugin"
 api_addr = "http://127.0.0.1:8200"
```

- `api_addr` defines the access port of vault. All the requests to vault will be done via this port.

- `Storage` defines the backend-storage type that vault will use to store all the encrypted data. Since this
  backend-storage is not a part of the vault, we define the access port of PostgreSQL server and the table name which is
  already created. Change the role, password, databaseName, and tableName according to your postgres parameters. **Note
  that we have disabled SSL for database requests**.

- In the `listener` part we have disabled TLS. If activated, TLS certificates and keys have to be provided here also.
  The example above listens on localhost port 8200 without TLS.

- Lastly, we have defined the `plugin directory` where the vault looks for plugins. Remember to change this according to
  your desired path where you stored the Vault bin file earlier.

## Vault Setup

Once you have created the config.hcl configuration file, we can now start our vault server. Open the terminal in the
folder containing the config file and start the server by running the following command:

```sh
    $ sudo vault server -config=config.hcl
```

Now that the vault server is up and running, it is actually in a sealed state, that is vault functionalities can't be
accessed yet. To access vault we need to unseal it. First, open another terminal window and initialize vault by running
the following commands:

```sh
    $ export VAULT_ADDR='http://127.0.0.1:8200'

    $ vault operator init
```

The first command is required for non-TLS mode. The output is a set of 5 shamir keys which have the capability to unseal
vault and an initial root token. Here vault is initialized in such a way that any 3 keys out of 5 are enough to unseal
vault. The root token is used to login into the vault. Only after logging in, you can start using vault. Store these in
a safe place for later use.

Start the unsealing process by running the command:

```sh
    $ vault operator unseal
```

The vault will ask you for an unseal key. Provide any one of the above 5. Run this command two more times and provide
two other keys. Vault should be unsealed now. To verify run the following command:

```sh
    $ vault status
```

The output should be something like this

```
  Key             Value
  ---             -----
  Seal Type       shamir
  Sealed          false <----this
  Total Shares    5
  Threshold       3
  Version         0.11.0
  Cluster Name    vault-cluster-8dea58da
  Cluster ID      8ac011b1-a830-663f-715a-cf5b3f87ae54
  HA Enabled      false
```

If you see the sealed key to have value false, vault is unsealed.

Now it's time to log in as the root admin. Run the command:

```sh
    $ vault login 85de6efd-d036-9f0d-1c64-5e18e63adee9
```

Provide your `root-token` in the above command and you should be logged in to vault as admin. Now we can send requests
to vault and set-up our plugin.

## Enable plugin

- If you previously have enabled this plugin, you need to disable it.

  ```sh
    $ vault secrets disable /api
  ```

- Calculate the SHA256 of the plugin and register it in Vault's plugin catalog.

  ```sh
    $ export SHA256=$(shasum -a 256 "/etc/vault/plugins/vault_plugin" | cut -d' ' -f1)

    $ vault write sys/plugins/catalog/secrets-api \
        sha_256="${SHA256}" \
        command="vault_plugin"
  ```

- Mount the secrets engine

  ```sh
    $ vault secrets enable \
      -path="api" \
      -plugin-name="secrets-api" \
      plugin
  ```

## Creating policies for application server

We need to define policies for the application server that will be using our vault. We don't want our application server
to have complete root access of vault, rather, it should just have the capability to update our Vault API plugin that we
just enabled. For that, we need to create another .hcl file (application.hcl as an example) to define the policies.

application.hcl:-

```
  path "api/*"
  {
    capabilities = ["read", "update"]
  }
```

To register this policy in vault, open terminal in the directory containing application.hcl and run the following
command:

```sh
  $ vault policy write application application.hcl
```

Now we can use application policy to define access capabilities of anyone using vault. More on that later.

## Enable userpass authentication method

We want our application server to login into vault using a particular `username` and `password` and should have access
capabilities defined by the `application` policy we created earlier. In order to do this we will be enabling userpass
authentication method.

```sh
  $ vault auth enable userpass
```

We then create a username and password using which our application server will log in. We also attach the application
policy in this command. The following command creates a user with username-"appserver" and password-"secret" with
application policy attached:

```sh
  $ vault write auth/userpass/users/appserver password=secret policies=application
```

We then give these credentials to the application server, who will use this username and password to login into the
vault. Note that anyone logged in by this method will have capabilities defined by the application policy.

We can easily create multiple user login credentials for different application servers.
