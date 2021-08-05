---
sidebar_position: 1
---

This part of setting up dq-vault can be done using two methods. You may follow any one of your choices.

- ## Method 1:-
  Using `Docker` to get your vault server up and running. You can find it in this [link](https://github.com/deqode/dq-vault/tree/main/setup). We have provided the required docker files in the setup folder.
- ## Method 2:-
  Setting up Vault manually. The steps are given below in this document, starting from vault installation to creating your own vault server by using the CLI.

:::info
If you are already done with setting up the vault server using method 1, you may go directly to ** [part 2](https://deqode.github.io/dq-vault/docs/guides/usage)** which elaborates the usage of the vault as an application server
:::
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

Assuming that you have golang installed and your `GOPATH` configured, get the plugin repository and run the build command in that folder:

```sh
  $ go build
```

This will you give you a binary executable file with the name `dq-vault`.

Now move this binary file to a directory which the vault will use as its plugin directory. The plugin directory is where the vault looks up for available plugins.

```sh
  $ mv dq-vault /etc/vault/plugins/vault_plugin
```

:::info
The above path is just an example, you can change the etc path to your own desired path.
:::