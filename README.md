# Ethereum plugin for Vault

The Ethereum secret backend is intended to provide many of the capabilities of an Ethereum wallet. It is designed to support smart contract continuous development practices including contract deployment and testing. It has only been exercised on private Ethereum chains and the Rinkeby testnet. Some of the functionality (creating accounts and signing contract creation transactions) can happen without any local `geth` node to be running. Other functionality (deploying contracts and sending transactions - still in development) will require the geth RPC interface.

## Features

This plugin provides services to:

* Create new externally controlled accounts (using a provided passphrase or a generated one.)
* Import JSON keystores (with provided passphrase.)
* Export JSON keystores
* Sign transactions for contract deployment
* Sign arbitrary data
* Send Ethereum
* Deploy contracts

All secrets in Vault are encrypted. However, for ease of integration with `geth`, the plugin stores the Ethereum private key in encrypted (JSON keystore) format. 

![Vault and Geth Topology](/doc/vault-geth.png?raw=true "Vault Ethereum Plugin")

## Quick start

Building the plugin and installing it will be covered later, but let's assume that has been done. It is important to note that the Vault Ethereum plugin can be mounted at any path. A common model is to use a well defined namespace for mounting Vault backends - for example, using the GitHub org/repo namespace: `ethereum/immutability-io/world-domination-token`. For this discussion, we assume that the Vault Ethereum plugin has been mounted at `ethereum`.

### Create new externally controlled accounts

Let's create an Ethereum account:

```sh
$ vault write ethereum/accounts/test generate_passphrase=true
```

That's all that is needed. The passphrase will be generated and stored in Vault. **NOTE:** we also assume that the vault client has been authenticated and has permission `write` to `ethereum/accounts/test`. Discussion of the [Vault authentication model can be found here](https://www.vaultproject.io/docs/concepts/auth.html).

The Ethereum plugin will return the following information from the above command:

```
Key     	Value
---     	-----
account 	0x2D9A87873C3735207bBcBBcb6f8Bc320CfcA8A5e
chain_id	4
rpc_url 	http://localhost:8545
```

The parameters `chain_id` and `rpc_url` are defaults and can be customized when writing an account. **NOTE**: the passphrase that is used to encrypt the keystore is **NOT** returned.

### Read externally controlled accounts

We can read the account stored at `ethereum/accounts/test` as follows:

```sh
$ vault read ethereum/accounts/test
```

```
Key             	Value
---             	-----
address         	0x2D9A87873C3735207bBcBBcb6f8Bc320CfcA8A5e
chain_id        	4
pending_balance 	0
pending_nonce   	0
pending_tx_count	0
rpc_url         	http://localhost:8545
```

### Read passphrase

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

### Create contracts

Suppose you have written a smart contract. Likely, it is only one or 2 deployment cycles away from yielding ICO riches. So, you better deploy it. The Vault plugin allows you to deploy a compiled smart contract.

Sending any transaction on the Ethereum network requires the payment of fees. So, you send the transaction that deploys a contract **from** an Ethereum account with a positive balance.

Assume that the compiled binary in the file `./out/Helloworld.bin`. Deployment is simple:

```sh
$ vault write ethereum/accounts/test/contracts/helloworld transaction_data=@Helloworld.bin value=10000000000000000000 gas_price=21000000000 gas_limit=1500000
```

The above command says: *Deploy a contract, named `helloworld`, from the account named `test`*

```
Key             	Value
---             	-----
account_address 	0x206d4B8aB00F1D3FdD3683A318776942f82A7F28
pending_balance 	200779500000000000000
pending_nonce   	7
pending_tx_count	0
tx_hash         	0x206ba52b1edd32510e6ab607bbfbba70369595210d22885b3067868a376e9677
```

When you deploy a contract, the contract address isn't immediately available. What is returned from the Vault-Ethereum plugin after a contract deployment is:

* `account_address`: The account that was used to sign and deploy the contract.
* `pending_balance`: The *pending* balance on the account that was used to sign and deploy the contract.
* `pending_nonce`: The *pending* nonce  on the account that was used to sign and deploy the contract.
* `pending_tx_count`: The *pending* transaction count  on the account that was used to sign and deploy the contract.
* `tx_hash`: The hash of the contract deployment transaction.

### Read contract address

Since the contract address isn't known at the point when the transaction is sent, so you have to **revisit** the contract (with a read operation) to determine the address:

```sh
$ vault read ethereum/accounts/test/contracts/helloworld
```

```
Key             	Value
---             	-----
contract_address	0x9dC730499BbAe80F4241a2523C516919C69339Af
tx_hash         	0x206ba52b1edd32510e6ab607bbfbba70369595210d22885b3067868a376e9677
```

### Import keystores from other wallets

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

Now, we can use the imported account as we did with our generated account.

### Export keystores

If we wish to export Vault managed keystores into a external wallet, we can:

```sh
$ vault write ethereum/accounts/test/export directory=/Users/immutability/.ethereum/keystore
```

```
Key 	Value
--- 	-----
path	/Users/immutability/.ethereum/keystore/UTC--2017-12-01T23-13-37.315592353Z--a152e7a09267bcff6c33388caab403b76b889939
```


### Signing arbitrary data

We can also sign arbitrary data using the `sign` endpoint:

```sh
$ vault write ethereum/accounts/test2/sign  data=@../data/test.txt
```

```
Key      	Value
---      	-----
signature	0xe81d649f2a295aa58ad0d67b2adf0f5f336e11a46bd69347f197f073244863406027daed083675b5af5c99b3f1608b53620cd02ca51a65b67773b1580552deb501
```

### Sending Ethereum

Now that we have accounts in Vault, we can drain them! We can send ETH to other accounts on the network. (In my case, it must be emphasized: this is a private test network or Rinkeby.) Assuming there are funds in the account, we can send ETH to another address. In this case, we write a `debit` to the `test3` account:

```sh
$ vault write ethereum/accounts/test3/debit to=0x0374E76DA2f0bE85a9FdC6763864c1087e6Ed28b value=10000000000000000000
```

```
Key    	Value
---    	-----
tx_hash	0xe99f3de1dfbae82121a009b9d3a2a60174f2904721ec114a8fc5454a96e62ba8

```

This defaults `gas_limit` to 50000 with a default `gas_price` of 20 gwei.


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

## [Plugin API](https://github.com/immutability-io/vault-ethereum/blob/master/API.md)

The complete API to the plugin is documented [here](https://github.com/immutability-io/vault-ethereum/blob/master/API.md). Each API is exemplified using curl as a sample REST client.

## Plugin Setup

I assume some familiarity with Vault and Vault's plugin ecosystem. If you are not familiar, please [refer to this](https://www.vaultproject.io/guides/plugin-backends.html). I realize that it is a lot to ask for someone to be so familiar with something so new. I will be writing a series of tutorials around this space in the near future. I will link them here when done. I will (eventually) provide a Vagrant box and scripts that configure a Vault server that supports this plugin.

For this to work, you must have a Vault server already running, unsealed, and authenticated.

### Build the plugin

You can use the `Makefile` or simply us `go build` from this project's root directory.

## Install the plugin

It is assumed that your Vault configuration specifies a `plugin_directory`. Mine is:

```
$ cat vault-config.hcl
...
plugin_directory="/etc/vault.d/vault_plugins"
...
```

Move the compiled plugin into Vault's configured `plugin_directory`:

  ```sh
  $ mv vault-ethereum /etc/vault.d/vault_plugins/vault-ethereum
  ```

Calculate the SHA256 of the plugin and register it in Vault's plugin catalog.

  ```sh
  $ export SHA256=$(shasum -a 256 "/etc/vault.d/vault_plugins/vault-ethereum" | cut -d' ' -f1)

  $ vault write sys/plugins/catalog/ethereum-plugin \
      sha_256="${SHA256}" \
      command="vault-ethereum --ca-cert=/etc/vault.d/root.crt --client-cert=/etc/vault.d/vault.crt --client-key=/etc/vault.d/vault.key"
  ```

If you are using Vault in `dev` mode, you don't need to supply the certificate parameters. For any real Vault installation, however, you will be using TLS.

## Mount the Ethereum secret backend

  ```sh
  $ vault mount -path="ethereum" -plugin-name="ethereum-plugin" plugin
  ```

## ToDo

More (much) to come soon...

## Credits

None of this would have been possible without the fantastic [tutorial](https://www.hashicorp.com/blog/building-a-vault-secure-plugin) on Vault Plugins by Seth Vargo. Seth is one of those rare individuals who can communicate the simple essence of a complex technology in practical terms.

I had the great fortune to attend DevCon3 in November and hear Andy Milenius speak with clarity and vision about how the Ethereum developer ecosystem should embrace the Unix philosophy - the same philosophy that makes **everything-as-code** possibly: simple tools, with clear focus and purpose, driven by repeatable and interoperable mechanics. So, when I returned from DevCon3 (and dug out from my work backlog - a week away is hard) I installed `seth` and `dapp` and found inspiration.

The community chat that the [dapphub](https://dapphub.com/) guys run (esp. Andy and Mikael and Daniel Brockman) is a super warm and welcoming place that pointed me towards code that greatly helped this experiment.

## License

This code is licensed under the MPLv2 license. Please feel free to use it. Please feel free to contribute.
