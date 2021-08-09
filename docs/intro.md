---
sidebar_position: 1
---

# Introduction 👇️

This vault plugin stores a user's mnemonic inside vault in an encrypted manner.
The plugin uses this stored mnemonic to derive a private key based on an HD wallet path provided by the user and signs a raw transaction given as input using that private key. All this process happens inside the vault and the user never knows the mnemonic (unless he has provided it manually) or the private key derived. All he needs to do is give a raw transaction as input and the vault returns a signed transaction. A particular user is identified in the vault using a UUID generated when the user is initially registered in the vault.

:::tip
This plugin inherits all security and encryption provided by best in class battle-tested [`HashiCorp Vault`](https://www.vaultproject.io/) for blockchain
:::
There will be two roles communicating with vault:

- **Admin**: The one who sets up the vault. 
- **Application Server**: The one who uses vault to read and update data.

This documentation guides you through two processes:-

1. How to set up the vault server (**For Admin**)
2. How to use vault to register a user and create a signature on demand (**for application server**)

The application server can communicate with a vault server using API requests/calls. Both CLI commands and API call methods have been included in this guide.

![Docusaurus](/img/vault-dq-192x192.png)

# Features 👌️
- Currently, in support for BTC/ETH
- Manage User's Secrets and Protect Sensitive Data of blockchain
- Supports both, Admin and Application users
- API-driven design
- Open Source
- Extend and integrate
- Easy setup using Docker