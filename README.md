# Ethereum plugin for Vault

The Ethereum secret backend is intended to provide many of the capabilities of an Ethereum wallet. It is designed to support smart contract continuous development practices including contract deployment and testing. It has only been exercised on private Ethereum chains and the Rinkeby testnet. Some of the functionality (creating accounts and signing contract creation transactions) can happen without any local `geth` node to be running. Other functionality (deploying contracts and sending transactions - still in development) will require the geth RPC interface.

## Features

This plugin provides services to:

* Create new externally controlled accounts (using a provided passphrase or a generated one.)
* Import JSON keystores (with provided passphrase.)
* Export JSON keystores.
* Sign transactions
* Send transactions (in development.)

All secrets in Vault are encrypted. However, for ease of integration with `geth`, the plugin stores the Ethereum private key in encrypted (JSON keystore) format. It is not necessary for this plugin to use a passphrase to protect private keys, however, at present that is the design choice.

![Vault and Geth Topology](/doc/vault-geth.png?raw=true "Vault Ethereum Plugin")

## Quick start

Building the plugin and installing it will be covered later, but let's assume that has been done. It is important to note that the Vault Ethereum plugin can be mounted at any path. A common model is to use a well defined namespace for mounting Vault backends - for example, using the GitHub org/repo namespace: `ethereum/immutability-io/world-domination-token`. For this discussion, we assume that the Vault Ethereum plugin has been mounted at `ethereum`.

Let's create an Ethereum account:

```sh
$ vault write ethereum/accounts/test generate_passphrase=true
```

That's all that is needed. The passphrase will be generated and stored in Vault. **NOTE:** we also assume that the vault client has been authenticated and has permission `write` to `ethereum/accounts/test`. Discussion of the [Vault authentication model can be found here](https://www.vaultproject.io/docs/concepts/auth.html).

The Ethereum plugin will return the following information from the above command:

```
Key     	Value
---     	-----
account 	0xD010BB32d6243d70Eb863610674a50EaEdfF8474
chain_id	4
keystore	{"address":"d010bb32d6243d70eb863610674a50eaedff8474","crypto":{"cipher":"aes-128-ctr","ciphertext":"fddf50de1041e87d45049fb7c1a2826487d08fc4f0664ab1decbf271e141d706","cipherparams":{"iv":"50d40092713acc1fb915d95b7e896f8b"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c274e45dfdfc8c912f96130936b48102dfe5c216f3fab417d3109446de448f0a"},"mac":"a51ce95867fefadc2846b30ddb8ba3d911faf5649718c9344ecb61c337ae806c"},"id":"a04d13d2-2319-481f-82e1-e86b3fc6a86a","version":3}
rpc_url 	localhost:8545
```

The parameters `chain_id` and `rpc_url` are defaults and can be customized when writing an account. The `keystore` can be copied to the `geth` keystore directory, if desired. Also note that the passphrase that is used to encrypt the `keystore` is **NOT** returned.

We can read the account stored at `ethereum/accounts/test` as follows:

```sh
$ vault read ethereum/accounts/test
```

```
Key     	Value
---     	-----
address 	0xD010BB32d6243d70Eb863610674a50EaEdfF8474
chain_id	4
keystore	{"address":"d010bb32d6243d70eb863610674a50eaedff8474","crypto":{"cipher":"aes-128-ctr","ciphertext":"fddf50de1041e87d45049fb7c1a2826487d08fc4f0664ab1decbf271e141d706","cipherparams":{"iv":"50d40092713acc1fb915d95b7e896f8b"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c274e45dfdfc8c912f96130936b48102dfe5c216f3fab417d3109446de448f0a"},"mac":"a51ce95867fefadc2846b30ddb8ba3d911faf5649718c9344ecb61c337ae806c"},"id":"a04d13d2-2319-481f-82e1-e86b3fc6a86a","version":3}
rpc_url 	localhost:8545
```

If we need to access the passphrase, we can do the following:

```sh
$ vault read ethereum/accounts/test/passphrase
```

```
Key       	Value
---       	-----
passphrase	durable-wrongdoer-keenness-clergyman-dorsal-registrar
```

The passphrase is accessible at a different path than the account. We do this because Vault ACLs are path based and this allows Vault administrators to parcel out different policies to different actors based on their roles.

Now suppose we have an Ethereum contract we need to sign - the compiled binary in the file `./out/Helloworld.bin`. Signing is simple:

```sh
$ vault write ethereum/accounts/test/sign-contract transaction_data=@./out/Helloworld.bin value=3 gas_limit=1000000 gas_price=500000 nonce=1
```

```
Key      	Value
---      	-----
signed_tx	0xf90231018307a120830f42408003b901e03630363036303430353233343135363130303066353736303030383066643562363064333830363130303164363030303339363030306633303036303630363034303532363030343336313036303439353736303030333537633031303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303039303034363366666666666666663136383036333630666534376231313436303465353738303633366434636536336331343630366535373562363030303830666435623334313536303538353736303030383066643562363036633630303438303830333539303630323030313930393139303530353036303934353635623030356233343135363037383537363030303830666435623630376536303965353635623630343035313830383238313532363032303031393135303530363034303531383039313033393066333562383036303030383139303535353035303536356236303030383035343930353039303536303061313635363237613761373233303538323064346234393631313833383934636631313936626361666262653464323537336139323532393664666638326139646362633065386264383032376231353366303032392ba0d26ac3ecfa8e7f23dea90e87d56b9985717c39ef66754cd103549ff0c211861da079fa8aff47bd2adc8d7549b043354203eff44a035b8c8d0216b9eb7bbbe35731
```

Lastly, suppose we already have Ethereum keystores and we are convinced that storing them (and their passphrases) in Vault is something we want to do. The plugin supports importing JSON keystores. **NOTE:** you have to provide the path to a single keystore - this plugin doesn't support importing an entire directory yet.

```sh
$ ls -la ~/.ethereum/keystore
```

```
total 24
drwxr-xr-x  5 immutability  admin  170 Dec  2 11:57 .
drwxr-xr-x  3 immutability  admin  102 Dec  2 11:55 ..
-rw-r--r--  1 immutability  admin  492 Dec  2 11:56 UTC--2017-12-01T23-13-37.315592353Z--a152e7a09267bcff6c33388caab403b76b889939
-rw-r--r--  1 immutability  admin  492 Dec  2 11:56 UTC--2017-12-01T23-13-56.838050955Z--0374e76da2f0be85a9fdc6763864c1087e6ed28b
-rw-r--r--  1 immutability  admin  492 Dec  2 11:57 UTC--2017-12-01T23-14-16.032409548Z--f19a9a9b2ad60c66429451075046869a9b7014f7
```

As will be discussed in the next section, handling passphrases is always problematic. Care should be taken when importing a keystore not to leak the passphrase to the shell's history file or to the environment:

```sh
$ read PASSPHRASE; vault write ethereum/import/test2 path=/Users/immutability/.ethereum/keystore/UTC--2017-12-01T23-13-37.315592353Z--a152e7a09267bcff6c33388caab403b76b889939 passphrase=$PASSPHRASE; unset PASSPHRASE
```

```
Key    	Value
---    	-----
address	0xa152E7a09267bcFf6C33388cAab403b76B889939
```

Now, we can use the imported account as we did with our generated account:

```sh
$ vault write ethereum/accounts/test2/sign-contract transaction_data=@./out/Helloworld.bin value=3 gas_limit=1000000 gas_price=500000 nonce=1
```

## Storing passphrases

Keeping passphrases on the same storage medium as the encrypted private key is probably the most controversial part of this design. The justification is based on the context in which this plugin is expected to be used.

In a DevOps environment, we leverage automation across the pipeline. We often have non-human actors engaged in the process of deployment and testing. A typical practice in the Ethereum community is to `unlock` an account for a period of time. Since there is no authentication needed to use this `unlocked` account, this creates a window of opportunity for bad actors to send transactions. Audit controls in this scenario are limited as well.

Another alternative to `unlocking` an account is to sign a transaction in the context of a human user providing a passphrase. This makes automation ineffective.

Also, having users handling passphrases with any frequency - the kind of frequency that we have in a typical development environment - makes exposure of passphrases likely. A tired developer will forget that they exported a variable or put a passphrase in a file.

### Vault can help

Every interaction with the Vault Ethereum backed needs to be [authenticated](https://www.vaultproject.io/docs/concepts/auth.html). Because Vault decouples authentication from storage, you can tailor the authentication mechanism to fit your needs:

* A variety of authentication providers (AppRole, AWS, Google Cloud, Kubernetes, GitHub, LDAP, MFA, Okta, RADIUS, TLS Certificates, Tokens, and Username & Password) each intended to support the unique context of the workflow.
* A sophisticated single-use token mechanism (https://www.vaultproject.io/docs/concepts/response-wrapping.html).

Every path in Vault can be protected with ACLs: You can allow some authenticated identities to import keystores, others to export them, and segregate access by account. Every access to Vault is audited as well, so it is pretty easy to diagnose access issues.

Vault encrypts all data and provides an excellent **cold storage** solution - when the Vault is sealed it requires a quorum of Shamir secret shards to bring it back online. This is functionally equivalent to what the Ethereum community would call a *multi-sig wallet*.

Furthermore, if you are an enterprise and capable of paying for [Enterprise Vault](https://www.hashicorp.com/products/vault) you can leverage HSMs as a persistence mechanism for Vault keys. This makes Vault equivalent to what the Ethereum folks call a hardware wallet. (It is very comparable to what [Gemalto and Ledger](https://www.gemalto.com/press/Pages/Gemalto-and-Ledger-Join-Forces-to-Provide--Security-Infrastructure-for-Cryptocurrency-Based-Activities-.aspx) have developed.)

## Vault Ethereum API

Vault provides a CLI that wraps the Vault REST interface. Any HTTP client (including the Vault CLI) can be used for accessing the API.

### LIST ACCOUNTS

This endpoint will list all accounts stores at a path.



## Plugin Setup

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
