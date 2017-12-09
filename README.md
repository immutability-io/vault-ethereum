# Ethereum plugin for Vault

This repository contains code for a HashiCorp Vault Plugin. I am building this Plugin for a few reasons:

1. To gain familiarity with using `geth` as a library.
2. To get experience with Vault Plugins.
3. As a *first step* towards developing a DevOps workflow for building Ethereum smart contracts.

## Credits

None of this would have been possible without the fantastic [tutorial](https://www.hashicorp.com/blog/building-a-vault-secure-plugin) on Vault Plugins by Seth Vargo. Seth is one of those rare individuals who can communicate the simple essence of a complex technology in practical terms.

I have been developing DevOps workflow solutions using the **everything-as-code** mantra for a while. I build GitHub oriented workflows that leverage extensive automation to build VMs and containers, provision infrastructure and manage policies and credentials using immutable architecture principles. When I moved into the Ethereum ecosystem - where immutability is baked in from the ground up - I was a bit chagrined to discover that the developer (much less DevOps) experience seemed a bit immature in comparison: nearly every tutorial on smart contract deployment talked about "pasting code" into a Wallet GUI.

I had the great fortune to attend DevCon3 in November and hear Andy Milenius speak with clarity and vision about how the Ethereum developer ecosystem should embrace the Unix philosophy - the same philosophy that makes **everything-as-code** possibly: simple tools, with clear focus and purpose, driven by repeatable and interoperable mechanics. So, when I returned from DevCon3 (and dug out from my work backlog - a week away is hard) I installed `seth` and `dapp` and found inspiration.

The community chat that the [dapphub](https://dapphub.com/) guys run (esp. Andy and Mikael and Daniel Brockman) is a super warm and welcoming place that pointed me towards code that greatly helped this experiment.

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

## ToDo

More to come soon...

## License

This code is licensed under the MPLv2 license.
