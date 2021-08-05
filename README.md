# dq-vault - Hashicorp vault BTC/ETH plugin

<p align="center"><img src="https://deqode.github.io/dq-vault/assets/images/vault-dq-192x192-202df720d6d8d239d0fbf4cdc208c1c8.png"></p>

![GitHub](https://img.shields.io/github/license/deqode/dq-vault)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/deqode/dq-vault)](https://pkg.go.dev/github.com/github.com/deqode/dq-vault)
[![Go Report Card](https://goreportcard.com/badge/github.com/deqode/dq-vault)](https://goreportcard.com/report/github.com/deqode/dq-vault)
![GitHub last commit](https://img.shields.io/github/last-commit/deqode/codeanalyser)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/deqode/dq-vault)

This vault plugin stores a user's mnemonic inside vault in an encrypted manner. The plugin uses this stored mnemonic to derive a private key based on an HD wallet path provided by the user and signs a raw transaction given as input using that private key. All this process happens inside the vault and the user never knows the mnemonic (unless he has provided it manually) or the private key derived. All he needs to do is give a raw transaction as input and the vault returns a signed transaction. A particular user is identified in the vault using a UUID generated when the user is initially registered in the vault.

There will be two roles communicating with vault:

1. Admin: The one who sets up the vault.
2. Application Server: The one who uses vault to read and update data.

**The application server can communicate with a vault server using API requests/calls. Both CLI commands and API call methods have been included in this guide.**


Visit this [link](https://deqode.github.io/dq-vault/) for full documentation of ```dq-vault```

# Installation

This part of setting up vault can be done using two methods. You may follow any one of your choices.

- ## Method 1:-
  Using `Docker` to get your vault server up and running. You can find it in this [link](https://gitlab.com/arout/Vault/blob/master/setup/README.md). We have provided the required docker files in the setup folder.
- ## Method 2:-
  Setting up Vault manually. The steps are given below in this README, starting from vault installation to creating your own vault server by using the CLI.

**If you are already done with setting up the vault server using method 1, you may go directly to [part 2](http://localhost:3000/dq-vault/docs/guides/plugin-usage) which elaborates the usage of the vault as an application server.**

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

## License
```
Copyright 2021, DeqodeLabs (https://deqode.com/)

Licensed under the MIT License(the "License");

```
<p align="center"><img src="https://deqode.com/wp-content/uploads/presskit-logo.png" width="400"></p>