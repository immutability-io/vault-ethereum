# Ethereum plugin for Vault

The Ethereum secret backend is intended to provide many of the capabilities of an Ethereum wallet. It is designed to support smart contract continuous development practices including contract deployment and testing. It has only been exercised on private Ethereum chains and the Rinkeby testnet. Some of the functionality (creating accounts and signing contract creation transactions) can happen without any local `geth` node to be running. Other functionality (deploying contracts and sending transactions - still in development) will require the geth RPC interface.

## Features

This plugin provides services to:

* Create new externally controlled accounts (using a provided passphrase or a generated one.)
* Import JSON keystores (with provided passphrase.)
* Export JSON keystores (in development.)
* Sign transactions for contract deployment
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

Vault provides a CLI that wraps the Vault REST interface. Any HTTP client (including the Vault CLI) can be used for accessing the API. Since the REST API produces JSON, I use the wonderful [jq](https://stedolan.github.io/jq/) for the examples.

### LIST ACCOUNTS

This endpoint will list all accounts stores at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `LIST`  | `:mount-path/accounts`  | `200 application/json` |

#### Parameters

* `path` (`string: <required>`) - Specifies the path of the accounts to list. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request LIST \
    https://localhost:8200/v1/ethereum/accounts | jq .
```

#### Sample Response

The example below shows output for a query path of `/ethereum/accounts/` when there are 2 accounts at `/ethereum/accounts/test` and `/ethereum/accounts/test`.

```
{
  "request_id": "f5689b77-ff54-8aed-27e0-1be52ab4fd61",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "keys": [
      "test",
      "test2"
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}

```

### READ ACCOUNT

This endpoint will list details about the Ethereum account at a path. The passphrase will **NOT** be revealed.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/accounts/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to read. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/accounts/test | jq .
```

#### Sample Response

The example below shows output for a read of `/ethereum/accounts/test`. Note the encoding of the keystore.

```
{
  "request_id": "f6f15161-12f6-e0bf-32de-700d5a40bab7",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x87f12ea8D186B9aDd3209C0Ee8B8C4672d8b1A43",
    "chain_id": "4",
    "keystore": "{\"address\":\"87f12ea8d186b9add3209c0ee8b8c4672d8b1a43\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"d440dcacd5d74bf2aa7d716ee2381e165f3434d3c5f42948e7aef315daea430d\",\"cipherparams\":{\"iv\":\"1d2dba8aae7f213634d175b29f2598ce\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"eb625d810cf5813e15de74f23f72802bcb5aadef4557a24097e6d9ff1c482fd0\"},\"mac\":\"2d9d3b5242971336e5966e0e94622c889b53e6fffbea43deb7ba0738a31dd63a\"},\"id\":\"006f9432-b125-4c2c-9ad3-edbac905b671\",\"version\":3}",
    "rpc_url": "localhost:8545"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### CREATE ACCOUNT

This endpoint will create an Ethereum account at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to create. This is specified as part of the URL.
* `rpc_url` (`string: <optional> default:"localhost:8545"`) - Specifies the URL of the 'geth' node.
* `chain_id` (`string: <optional> default:"4"`) - Specifies the Ethereum network. Defaults to Rinkeby.
* `generate_passphrase` (`boolean: <optional> default:false`) - Determines whether the passphrase will be generated.
* `passphrase` (`string: <optional>`) - If `generate_passphrase` is false, a `passphrase` must be provided.
* `words` (`integer: <optional> default:"6"`) - Specifies the number of words to use in the generated passphrase.
* `separator` (`string: <optional> default:"-"`) - Specifies the delimiter used to separate the words in the generated passphrase.

#### Sample Payload

```
{
  "rpc_url": "localhost:8545",
  "chain_id": "1977",
  "generate_passphrase": true
}
```

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test3 | jq .
```

#### Sample Response

The example below shows output for the successful creation of `/ethereum/accounts/test3`. Note the encoding of the keystore.

```
{
  "request_id": "914c5797-815e-3d4e-d9de-b4978ac1e267",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "account": "0x55BcB4ba4BdE352B828deFaA45ae1880DbDb9A22",
    "chain_id": "1977",
    "keystore": "{\"address\":\"55bcb4ba4bde352b828defaa45ae1880dbdb9a22\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"3b81d2e79fdd76400c5fa2e2afe3c425f2c063027d51f0e006fb9575da54c70a\",\"cipherparams\":{\"iv\":\"d1646e44ece77140a9bdf86f01444329\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"87af64f3696ebb0d4595fefca4070f960361701a7600e5fa90b75e63d1094e90\"},\"mac\":\"63ca10b230fa2d438ebb873a133368eb7bb972e528fe2fd6985b1ac4bfca7dc8\"},\"id\":\"49e2446b-48fd-4159-bc6f-4476662dbc83\",\"version\":3}",
    "rpc_url": "localhost:8545"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### UPDATE ACCOUNT/RE-ENCRYPT KEYSTORE

This endpoint will re-encrypt the keystore for an Ethereum account with a new passphrase.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `PUT`  | `:mount-path/accounts/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to update. This is specified as part of the URL.
* `generate_passphrase` (`boolean: <optional> default:false`) - Determines whether the passphrase will be generated.
* `passphrase` (`string: <optional>`) - If `generate_passphrase` is false, a `passphrase` must be provided.
* `words` (`integer: <optional> default:"6"`) - Specifies the number of words to use in the generated passphrase.
* `separator` (`string: <optional> default:"-"`) - Specifies the delimiter used to separate the words in the generated passphrase.

#### Sample Payload

```
{
  "generate_passphrase": true
}
```

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request PUT \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test3 | jq .
```

#### Sample Response

The example below shows output for the successful re-encryption of the keystore for `/ethereum/accounts/test3`. Note the encoding of the keystore.

```
{
  "request_id": "4dd998b7-40e0-fa86-23f2-b39da925cbfb",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x55BcB4ba4BdE352B828deFaA45ae1880DbDb9A22",
    "chain_id": "1977",
    "keystore": "{\"address\":\"55bcb4ba4bde352b828defaa45ae1880dbdb9a22\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"d52bf4c5fe0bed07e489c2463646b0cef28e0f825d15d828bf10cff7191075e6\",\"cipherparams\":{\"iv\":\"508706acf516376cac47a94d4134888b\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"94b89865d53bec30d05c534d53a16553bea14b9c3797571ad67f3735eca2350c\"},\"mac\":\"657340088ef77777b7dcf789e2f669acc723f58b9dd74f89bc4d97b2867330b8\"},\"id\":\"49e2446b-48fd-4159-bc6f-4476662dbc83\",\"version\":3}",
    "rpc_url": "localhost:8545"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}

```

### DELETE ACCOUNT

This endpoint will delete the account - and its passphrase - from Vault.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `DELETE`  | `:mount-path/accounts/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to update. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request DELETE \
    https://localhost:8200/v1/ethereum/accounts/test3
```

#### Sample Response

There is no response payload.

### IMPORT ACCOUNT

This endpoint will import a JSON Keystore and passphrase into Vault at a path. It will create an account and map it to the `:mount-path/accounts/:name`. If an account already exists for this name, the operation fails.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/import/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to create. This is specified as part of the URL.
* `path` (`string: <required>`) - The path of the JSON keystore file.
* `rpc_url` (`string: <optional> default:"localhost:8545"`) - Specifies the URL of the 'geth' node.
* `chain_id` (`string: <optional> default:"4"`) - Specifies the Ethereum network. Defaults to Rinkeby.
* `passphrase` (`string: <required>`) - The `passphrase` that was used to encrypt the keystore.

#### Sample Payload

Be careful with those passphrases!

```sh
read PASSPHRASE; read  PAYLOAD_WITH_PASSPHRASE <<EOF
{"path":"/Users/tssbi08/.ethereum/keystore/UTC--2017-12-01T23-13-37.315592353Z--a152e7a09267bcff6c33388caab403b76b889939", "passphrase":"$PASSPHRASE"}
EOF
unset PASSPHRASE
```

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data $PAYLOAD_WITH_PASSPHRASE \
    https://localhost:8200/v1/ethereum/import/test3 | jq .
    unset PAYLOAD_WITH_PASSPHRASE
```

#### Sample Response

The example below shows output for the successful creation of `/ethereum/accounts/test3`. Note the encoding of the keystore.

```
{
  "request_id": "c8b79326-74eb-c75e-a602-bd0609ba9a10",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0xa152E7a09267bcFf6C33388cAab403b76B889939"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### SIGN ETHEREUM CONTRACT

This endpoint will sign a provided Ethereum contract.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:name/sign-contract`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `transaction_data` (`string: <required>`) - The compiled Ethereum contract.
* `value` (`string: <required>`) - The amount of ether.
* `nonce` (`string: <optional> - defaults to "1"`) - The nonce for the transaction
* `gas_price` (`string: <required>`) - The price in gas for the transaction.
* `gas_limit` (`string: <required>`) - The gas limit for the transaction.

#### Sample Payload

```sh

{
  "transaction_data": "6060604052341561000f57600080fd5b60d38061001d6000396000f3006060604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c14606e575b600080fd5b3415605857600080fd5b606c60048080359060200190919050506094565b005b3415607857600080fd5b607e609e565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a72305820d4b4961183894cf1196bcafbbe4d2573a925296dff82a9dcbc0e8bd8027b153f0029",
  "value":"3",
  "gas_limit":"1000000",
  "gas_price":"500000",
  "nonce":"1"
}
```

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test2/sign-contract | jq .
```

#### Sample Response

The example below shows output for the successful signing of a contract by the private key associated with  `/ethereum/accounts/test2`.

```
{
  "request_id": "494f7e52-1e1b-e4b1-677d-acfd43e9c317",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "signed_tx": "0xf90231018307a120830f42408003b901e03630363036303430353233343135363130303066353736303030383066643562363064333830363130303164363030303339363030306633303036303630363034303532363030343336313036303439353736303030333537633031303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303039303034363366666666666666663136383036333630666534376231313436303465353738303633366434636536336331343630366535373562363030303830666435623334313536303538353736303030383066643562363036633630303438303830333539303630323030313930393139303530353036303934353635623030356233343135363037383537363030303830666435623630376536303965353635623630343035313830383238313532363032303031393135303530363034303531383039313033393066333562383036303030383139303535353035303536356236303030383035343930353039303536303061313635363237613761373233303538323064346234393631313833383934636631313936626361666262653464323537336139323532393664666638326139646362633065386264383032376231353366303032392ca0c63156377cc040bbf2be7d3a045bf4b8fa88f4969159f0d4377dfd0ac6fd76e2a02fa4f5dd0058d4343a4402918bfcb858a5da3fcb4023ebeb4de1bb469cb1122a"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

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
