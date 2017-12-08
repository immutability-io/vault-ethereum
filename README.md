# Ethereum plugin for Vault


This repository contains code for a HashiCorp Vault Plugin. It is
intended to support Ethereum smart contract deployment - especially in a DevOps setting. Ultimately,
it is intended to provide an enterprise *wallet* capability.

## Setup

I assume some familiarity with Vault and Vault's plugin
ecosystem. If you are not familiar, please [refer to this](https://www.vaultproject.io/guides/plugin-backends.html)

You must have a Vault server already running, unsealed, and
authenticated. I will provide a Vagrant box and scripts that configure a Vault server that supports this plugin.

1. Download and decompress the latest plugin binary from the Releases tab on
GitHub. Alternatively you can compile the plugin from source.

1. Move the compiled plugin into Vault's configured `plugin_directory`:

  ```sh
  $ mv vault-ethereum /etc/vault.d/vault_plugins/vault-ethereum
  ```

1. Calculate the SHA256 of the plugin and register it in Vault's plugin catalog.
If you are downloading the pre-compiled binary, it is highly recommended that
you use the published checksums to verify integrity.

  ```sh
  $ export SHA256=$(shasum -a 256 "/etc/vault.d/vault_plugins/vault-ethereum" | cut -d' ' -f1)

  $ vault write sys/plugins/catalog/ethereum-plugin \
      sha_256="${SHA256}" \
      command="vault-ethereum --ca-cert=/etc/vault.d/root.crt --client-cert=/etc/vault.d/vault.crt --client-key=/etc/vault.d/vault.key"
  ```

1. Mount the auth method:

  ```sh
  $ vault mount -path="ethereum" -plugin-name="ethereum-plugin" plugin
  ```

## Authenticating with the Shared Secret

To authenticate, the user supplies the shared secret:

```sh
$ vault write ethereum/accounts/foo value=bar

```

The response will be a standard auth response with some token metadata:

```text
Key  	Value
---  	-----
value	bar
```
## License

This code is licensed under the MPLv2 license.
