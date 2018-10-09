# Hashicorp Vault Secrets Plugin

## Introduction

This vault plugin stores a user's mnemonic inside vault in an encrypted manner. The plugin uses this stored mnemonic to derive a private key based on an HD wallet path provided by the user and signs a raw transaction given as input using that private key. All this process happens inside the vault and the user never knows the mnemonic (unless he has provided it manually) or the private key derived. All he needs to do is give a raw transaction as input and the vault returns a signed transaction. A particular user is identified in the vault using a UUID generated when the user is initially registered in the vault.

There will be two roles communicating with vault:

1. Admin: The one who sets up the vault.
2. Application Server: The one who uses vault to read and update data.

This readme guides you through two processes:-

1. How to set up the vault server (For Admin)
2. How to use vault to register a user and create a signature on demand (for application server)

The application server can communicate with a vault server using API requests/calls. Both CLI commands and API call methods have been included in this guide.

# PART 1:- SETTING UP VAULT

This part of setting up vault can be done using two methods. You may follow any one of your choices.

- ## Method 1:-
  Using `Docker` to get your vault server up and running. You can find it in this [link](https://gitlab.com/arout/Vault/blob/master/setup/README.md). We have provided the required docker files in the setup folder.
- ## Method 2:-
  Setting up Vault manually. The steps are given below in this README, starting from vault installation to creating your own vault server by using the CLI.

**If you are already done with setting up the vault server using method 1, you may go directly to [part 2](https://gitlab.com/arout/Vault/blob/master/README.md#part-2-using-vault) which elaborates the usage of the vault as an application server.**

## Vault installation

The first thing you need to do is to install vault to set-up a vault server.

- To install Vault, find the [appropriate package](https://www.vaultproject.io/downloads.html) for your system and download it. Vault is packaged as a zip archive.
- After downloading Vault, unzip the package. Vault runs as a single binary named vault.
- Copy the vault binary to your `PATH`. In Ubuntu, PATH should be the `usr/bin` directory.
- To verify the installation, type vault in your terminal. You should see help output similar to the following:

  ```
    $ vault
    Usage: vault <command> [args]

    Common commands:
        read        Read data and retrieves secrets
        write       Write data, configuration, and secrets
        delete      Delete secrets and configuration
        list        List data or secrets
        login       Authenticate locally
        server      Start a Vault server
        status      Print seal and HA status
        unwrap      Unwrap a wrapped secret

    Other commands:
        audit          Interact with audit devices
        auth           Interact with auth methods
        lease          Interact with leases
        operator       Perform operator-specific tasks
        path-help      Retrieve API help for paths
        policy         Interact with policies
        secrets        Interact with secrets engines
        ssh            Initiate an SSH session
        token          Interact with tokens
  ```

- You can find the official installation guide [here](https://www.vaultproject.io/intro/getting-started/install.html)

## Get go files and Build plugin

Assuming that you have golang installed and your GOPATH configured, get the plugin repository and run the build command in that folder:

```sh
  $ go build
```

This will you give you a binary executable file with the name `Vault`.

Now move this binary file to a directory which the vault will use as its plugin directory. The plugin directory is where the vault looks up for available plugins.

```sh
  $ mv Vault /etc/vault/plugins/vault_plugin
```

**The above path is just an example, you can change the etc path to your own desired path.**

## Set up postgres

Assuming that you have PostgreSQL installed in your system, you need to create a table which will be used by Vault to store it's encrypted data.

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

This is the step where we define the vault configurations. Vault supports .hcl files to write your configurations and set-up vault accordingly.

Given below is an example of a config.hcl file:

config.hcl-

```
 "api_addr" = "http://127.0.0.1:8200"

 storage "postgresql" {
   connection_url = "postgres://role:password@localhost:5432/databaseName?sslmode=disable"
   table = "tableName"
 }

 listener "tcp" {
 address     = "127.0.0.1:8200"
 tls_disable = 1
 }

 "plugin_directory" = "/etc/vault/plugins"
```

- `api_addr` defines the access port of vault. All the requests to vault will be done via this port.

- `Storage` defines the backend-storage type that vault will use to store all the encrypted data. Since this backend-storage is not a part of the vault, we define the access port of PostgreSQL server and the table name which is already created. Change the role, password, databaseName, and tableName according to your postgres parameters. **Note that we have disabled SSL for database requests**.

- In the `listener` part we have disabled TLS. If activated, TLS certificates and keys have to be provided here also. The example above listens on localhost port 8200 without TLS.

- Lastly, we have defined the `plugin directory` where the vault looks for plugins. Remember to change this according to your desired path where you stored the Vault bin file earlier.

## Vault Setup

Once you have created the config.hcl configuration file, we can now start our vault server. Open the terminal in the folder containing the config file and start the server by running the following command:

```sh
    $ sudo vault server -config=config.hcl
```

Now that the vault server is up and running, it is actually in a sealed state, that is vault functionalities can't be accessed yet. To access vault we need to unseal it. First, open another terminal window and initialize vault by running the following commands:

```sh
    $ export VAULT_ADDR='http://127.0.0.1:8200'

    $ vault operator init
```

The first command is required for non-TLS mode. The output is a set of 5 shamir keys which have the capability to unseal vault and an initial root token. Here vault is initialized in such a way that any 3 keys out of 5 are enough to unseal vault. The root token is used to login into the vault. Only after logging in, you can start using vault. Store these in a safe place for later use.

Start the unsealing process by running the command:

```sh
    $ vault operator unseal
```

The vault will ask you for an unseal key. Provide any one of the above 5. Run this command two more times and provide two other keys. Vault should be unsealed now. To verify run the following command:

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

Provide your `root-token` in the above command and you should be logged in to vault as admin. Now we can send requests to vault and set-up our plugin.

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

We need to define policies for the application server that will be using our vault. We don't want our application server to have complete root access of vault, rather, it should just have the capability to update our Vault API plugin that we just enabled. For that, we need to create another .hcl file (application.hcl as an example) to define the policies.

application.hcl:-

```
  path "api/*"
  {
    capabilities = ["read", "update"]
  }
```

To register this policy in vault, open terminal in the directory containing application.hcl and run the following command:

```sh
  $ vault policy write application application.hcl
```

Now we can use application policy to define access capabilities of anyone using vault. More on that later.

## Enable userpass authentication method

We want our application server to login into vault using a particular `username` and `password` and should have access capabilities defined by the `application` policy we created earlier. In order to do this we will be enabling userpass authentication method.

```sh
  $ vault auth enable userpass
```

We then create a username and password using which our application server will log in. We also attach the application policy in this command. The following command creates a user with username-"appserver" and password-"secret" with application policy attached:

```sh
  $ vault write auth/userpass/users/appserver password=secret policies=application
```

We then give these credentials to the application server, who will use this username and password to login into the vault. Note that anyone logged in by this method will have capabilities defined by the application policy.

We can easily create multiple user login credentials for different application servers.

# PART 2:- USING VAULT

The following part contains both cURL requests and CLI commands to communicate with vault server. If you have used docker for vault setup, then use cURL requests.

## Login as the application server

Log-in into the vault as application server using the following command:

```sh
  $ vault login -method=userpass username=appserver password=secret
```

API call

```sh
  $ curl \
    --request POST \
    --data '{"password": "secret"}' \
    http://127.0.0.1:8200/v1/auth/userpass/login/appserver
```

The command will return a token which will be used to keep the application server authenticated.

## Plugin Usage

Once we are logged in as an application server, we can use our API plugin to store mnemonic of HD wallet keys and also to sign raw transactions using those keys.

### Register user

Registers a user and stores the corresponding user's mnemonic in the vault. The request returns a unique id(UUID) of the user which will be later used to access the user's keys stored in the vault.

| Method | Path            | Produces                 |
| :----- | :-------------- | :----------------------- |
| `POST` | `/api/register` | `200 (application/json)` |

#### Parameters

- `username` `(string)` `optional` - Specifies the user-name of the user being registered.

- `mnemonic` `(string)` `optional` - Specifies the mnemonic to be stored. If not provided, a random mnemonic will be generated and stored.

- `passphrase` `(string)` `optional` - Specifies the passphrase.

#### CLI

```sh
  $ vault write api/register username=user

  Key     Value
  ---     -----
  uuid    c3f394de-919d-4a66-a1b3-7686642be430
```

#### API call

```
  $ cat payload.json
  {
    "username": "user",
    "mnemonic": "",
    "passphrase": ""
  }

  $ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/api/register
```

**The X-Vault-Token being passed in the header is the token received when the application server logged in. For API calls, the token has to be always passed in the header for authentication.**

```
  Response:

  {
    "request_id":"03cf1ec0-dbd3-9ce8-1663-067e91d680ab",
    "lease_id":"",
    "renewable":false,"lease_duration":0,
    "data":{
      "uuid":"9af93fcc-c41f-4c30-828f-c4b774573205"
      },
    "wrap_info":null,
    "warnings":null,
    "auth":null
  }
```

### Obtain user address and public key

Once a user is registered, we can obtain the address and public key of that corresponding user.

| Method | Path           | Produces                 |
| :----- | :------------- | :----------------------- |
| `POST` | `/api/address` | `200 (application/json)` |

#### Parameters

- `uuid` `(string)` `required` - Specifies the uuid of the user who will sign a transaction.

- `path` `(string)` `required` - Specifies the HD-wallet path.

- `coinType` `(uint16)` `required` - Specifies the coin-type Value of the coin to be used.

#### coinType

```
  - Bitcoin:0
  - Bitcoin Testnet:1
  - Ethereum:60
```

#### CLI

```sh
  $ vault write api/address uuid="214b8190-f9f5-4cf8-a4f1-3874ed5b05d0" path="m/44'/0'/0'/0/0" coinType=0

  Key          Value
  ---          -----
  address      1Ax1vX6dASU35gqyPDzX16np4rifvm8oF2
  publicKey    030df204bde019cbbe79bec99aecfacd3c5f321f404d2df47b42573c987393df80
  uuid         214b8190-f9f5-4cf8-a4f1-3874ed5b05d0
```

#### API call

```sh
  $ cat payload.json
  {
    "uuid": "9af93fcc-c41f-4c30-828f-c4b774573205",
    "path": "m/44'/0'/0'/0/0",
    "coinType": 0
  }

  $ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/api/address
```

```
  Response:

  {
    "request_id": "a98b3299-9b03-e730-eae4-6a0f727816aa",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "address": "1Ax1vX6dASU35gqyPDzX16np4rifvm8oF2",
        "publicKey": "030df204bde019cbbe79bec99aecfacd3c5f321f404d2df47b42573c987393df80",
        "uuid": "214b8190-f9f5-4cf8-a4f1-3874ed5b05d0"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
  }
```

### Create signature

Once a user is registered, we can now sign raw transactions just by using the user's UUID(which accesses the stored keys). As of now Bitcoin, Bitcoin Testnet and Ethereum transactions are supported.

| Method | Path             | Produces                 |
| :----- | :--------------- | :----------------------- |
| `POST` | `/api/signature` | `200 (application/json)` |

#### Parameters

- `uuid` `(string)` `required` - Specifies the uuid of the user who will sign a transaction.

- `path` `(string)` `required` - Specifies the HD-wallet path.

- `coinType` `(uint16)` `required` - Specifies the coin-type Value of the coin to be used.

- `payload` `(string)` `required` - Contains the raw transaction to be signed in JSON format.

#### coinType

```
  - Bitcoin:0
  - Bitcoin Testnet:1
  - Ethereum:60
```

#### Payload

Since payload contains the raw transaction, it's structure differs for Bitcoin and ethereum.

##### Bitcoin

```
  {
    inputs: [] of {txhash: string, vout: uint32}
    outputs: [] of {address: string, amount: int 64}
  }
```

- txhash refers to the txid containing the UTXO and vout points to the index of that UTXO.
- address refers to the payee address and amount refers to the amount of BTC you want to send.

Example payload:

```
  {
    "inputs":[{
        "txhash":"81b4c832d70cb56ff957589752eb412a4cab78a25a8fc52d6a09e5bd4404d48a",
        "vout":0
      },{
          "txhash":"9dd5264b09bd4aebc1d74b776e6669ba3f0e381ef2992c9434e4d0bee3068edb",
          "vout":0
        }],

    "outputs":[{
        "address":"1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa",
        "amount":91234
      },{
          "address":"1HPvK7CAYeHzCdBMBkuXeEsXdvX64yMkoE",
          "amount":91234
        }]
  }
```

##### Ethereum

```
  {
    Nonce : uint64
    Value: uint64
    GasLimit: uint64
    GasPrice: uint64
    To: string
    Data: string
    ChainID: int64
  }
```

Example payload:

```
  {
    "nonce":0,
    "value":1000000000,
    "gasLimit":21000,
    "gasPrice":30000000000,
    "to":"0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d",
    "data":"",
    "chainId":1
  }
```

The request finally returns a signature of a raw transaction which was signed inside vault using the following things:

- uuid of the user proposing the transaction.
- Stored mnemonic corresponding to the provided uuid.
- HD wallet path.

#### CLI

BTC:

```sh
  $ vault write api/signature uuid="c3f394de-919d-4a66-a1b3-7686642be430" \
    path="m/0/0" \
    payload="{\"inputs\":[{\"txhash\":\"b31695ff693b196d41600266d82bdf1092a4a55be608f41e1bde985408b16774\",\"vout\":0}],\"outputs\":[{\"address\":\"3BGgKxAsqoFyouTgUJGW3TAJdvYrk43Jr5\",\"amount\":91234}]}" \
    coinType=0

  Key          Value
  ---          -----
  signature    01000000017467b1085498de1b1ef408e65ba5a49210df2bd8660260416d193b69ff9516b3000000006a47304402200cd2c06db98cb1a71cbb7558506815d20933e4451ffda2760971b5e477c7766902206dc6aa33f3c05305a992fcf3f19d58953b55398f8052a0ae1f061ad8b38b3135012103e1a150d41f5d6871da8310e5ea8226f105716639483e3e2c79981d65392ce499ffffffff01626401000000000017a9146916ea9f8135de454ecb1c22ade111ff48fb7c9f8700000000
```

BTC testnet:

```sh
  $ vault write api/signature uuid="c3f394de-919d-4a66-a1b3-7686642be430" \
    path="m/0/0" \
    payload="{\"inputs\":[{\"txhash\":\"b31695ff693b196d41600266d82bdf1092a4a55be608f41e1bde985408b16774\",\"vout\":0}],\"outputs\":[{\"address\":\"3BGgKxAsqoFyouTgUJGW3TAJdvYrk43Jr5\",\"amount\":91234}]}" \
    coinType=1

  Key          Value
  ---          -----
  signature    01000000017467b1085498de1b1ef408e65ba5a49210df2bd8660260416d193b69ff9516b3000000006b483045022100d3323c41f117c4c1ef3e52fde37bc01b24fc6090de8dbeb6918a494bfea21ef602206ee496d9933eb5a9246808b96cbad4c8b84b9b5ad7a66afe045acc72f033e2d6012103c023b44933371f7d208bc0ff8a65505d67bf8750de913d21af8d194585ac7af0ffffffff01626401000000000017a9146916ea9f8135de454ecb1c22ade111ff48fb7c9f8700000000
```

ETH:

```sh
  $ vault write api/signature uuid="c3f394de-919d-4a66-a1b3-7686642be430" \
  path="m/0/0" \
  payload="{\"nonce\":0,\"value\":1000000000,\"gasLimit\":21000,\"gasPrice\":30000000000,\"to\":\"0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d\",\"data\":\"\",\"chainId\":1}" \
  coinType=60

  Key          Value
  ---          -----
  signature    0xf868808506fc23ac00825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d843b9aca008026a08a465e9d1c707d02f72360ab21d1a1be5faf84671413e7df0402e954a666cd79a04ab6481295d13f31fc4265888e8bd9962e200062889f162b320caf4c697f96c4
```

#### API call

BTC:

```sh
  $ cat payload.json
  {
    "uuid": "9af93fcc-c41f-4c30-828f-c4b774573205",
    "path": "m/0/0",
    "payload": "{\"inputs\":[{\"txhash\":\"b31695ff693b196d41600266d82bdf1092a4a55be608f41e1bde985408b16774\",\"vout\":0}],\"outputs\":[{\"address\":\"3BGgKxAsqoFyouTgUJGW3TAJdvYrk43Jr5\",\"amount\":91234}]}",
    "coinType": 0
  }

  $ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/api/signature
```

```
  Response:

  {
    "request_id":"c9f6916d-5985-6320-770b-cc3fb22b0b37",
    "lease_id":"",
    "renewable":false,
    "lease_duration":0,
    "data":{
      "signature":"01000000017467b1085498de1b1ef408e65ba5a49210df2bd8660260416d193b69ff9516b3000000006a47304402204355e8a9cd9f2e4cac867c8ab55a63f020404249051385f4170788b35d246d9602201b2ff4bbc7a9dd9eb5ebbfcffeed3de0a9138ca1606d302b407d3e99c092ac1e0121027276b9edee40a02957f237d79536205524c3864d0d351909cdf519adc60de6d4ffffffff01626401000000000017a9146916ea9f8135de454ecb1c22ade111ff48fb7c9f8700000000"
      },
    "wrap_info":null,
    "warnings":null,
    "auth":null
  }
```

BTC Testnet:

```sh
  $ cat payload.json
  {
    "uuid": "9af93fcc-c41f-4c30-828f-c4b774573205",
    "path": "m/0/0",
    "payload": "{\"inputs\":[{\"txhash\":\"b31695ff693b196d41600266d82bdf1092a4a55be608f41e1bde985408b16774\",\"vout\":0}],\"outputs\":[{\"address\":\"3BGgKxAsqoFyouTgUJGW3TAJdvYrk43Jr5\",\"amount\":91234}]}",
    "coinType": 1
  }

  $ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/api/signature
```

```
  Response:

  {
    "request_id":"c9f6916d-5985-6320-770b-cc3fb22b0b37",
    "lease_id":"",
    "renewable":false,
    "lease_duration":0,
    "data":{
      "signature":"01000000017467b1085498de1b1ef408e65ba5a49210df2bd8660260416d193b69ff9516b3000000006b483045022100d3323c41f117c4c1ef3e52fde37bc01b24fc6090de8dbeb6918a494bfea21ef602206ee496d9933eb5a9246808b96cbad4c8b84b9b5ad7a66afe045acc72f033e2d6012103c023b44933371f7d208bc0ff8a65505d67bf8750de913d21af8d194585ac7af0ffffffff01626401000000000017a9146916ea9f8135de454ecb1c22ade111ff48fb7c9f8700000000"
      },
    "wrap_info":null,
    "warnings":null,
    "auth":null
  }
```

ETH:

```sh
  $ cat payload.json
  {
    "uuid": "9af93fcc-c41f-4c30-828f-c4b774573205",
    "path": "m/0/0",
    "payload": "{\"nonce\":0,\"value\":1000000000,\"gasLimit\":21000,\"gasPrice\":30000000000,\"to\":\"0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d\",\"data\":\"\",\"chainId\":1}",
    "coinType": 60
  }

  $ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/api/signature
```

```
  Response:

  {
    "request_id":"c9f6916d-5985-6320-770b-cc3fb22b0b37",
    "lease_id":"",
    "renewable":false,
    "lease_duration":0,
    "data":{
      "signature":"0xf868808506fc23ac00825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d843b9aca008026a08a465e9d1c707d02f72360ab21d1a1be5faf84671413e7df0402e954a666cd79a04ab6481295d13f31fc4265888e8bd9962e200062889f162b320caf4c697f96c4"
      },
    "wrap_info":null,
    "warnings":null,
    "auth":null
  }
```

cURL requests for all of the above commands can be imported to postman via vault_postman_collection.json.

Note: If using postman, replace the existing token with yours in the header section of each request.
