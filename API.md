## Vault Ethereum API

Vault provides a CLI that wraps the Vault REST interface. Any HTTP client (including the Vault CLI) can be used for accessing the API. Since the REST API produces JSON, I use the wonderful [jq](https://stedolan.github.io/jq/) for the examples.

* [List Accounts](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#list-accounts)
* [Read Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#read-account)
* [Create Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#create-account)
* [Update Account/Re-Encrypt Keystore](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#update-accountre-encrypt-keystore)
* [Delete Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#delete-account)
* [Import Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#import-account)
* [Sign Ethereum Contract](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#sign-ethereum-contract)
* [Sign Data](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#sign-data)
* [Send Ethereum/Debit Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#send-ethereumdebit-account)

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


### SIGN DATA

This endpoint will sign the provided data.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:name/sign`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `data` (`string: <required>`) - Some data.

#### Sample Payload

```sh

{
  "data": "this is very important"
}
```

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test9/sign | jq .
```

#### Sample Response

The example below shows output for the successful signing of some data by the private key associated with  `/ethereum/accounts/test2`.

```
{
  "request_id": "5491a21c-7541-f48c-d573-0d241f12bfd3",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "signature": "0xb5d19f676208d20861336cfc38da1012716d15ca8e23a18fd46f65a18e6fef8f313d2ba6aa424f9c096076ceb8d6cd4bd48fac520e9df592e51869fd5cebad0801"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### SEND ETHEREUM/DEBIT ACCOUNT

This endpoint will debit an Ethereum account.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:name/debit`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `to` (`string: <required>`) - A Hex string specifying the Ethereum address to send the Ether `to`.
* `value` (`string: <required>`) - The amount of ether - in wei.
* `gas_price` (`string: <optional> - defaults to 20000000000`) - The price in gas for the transaction.
* `gas_limit` (`string: <optional> - defaults to 50000`) - The gas limit for the transaction.

#### Sample Payload

The following sends 10 ETH to `0xa152E7a09267bcFf6C33388cAab403b76B889939`.

```sh

{
  "value":"10000000000000000000",
  "to": "0x0374E76DA2f0bE85a9FdC6763864c1087e6Ed28b"
}
```

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test2/debit | jq .
```

#### Sample Response

The example below shows the output for the successfully sending ETH from `/ethereum/accounts/test2`. The Transaction hash is returned.

```
{
  "request_id": "3660838f-2ddc-f92b-8796-9351f7d123dd",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "tx_hash": "0xdd675f368e5002212f8bdb50f17a0cd8e4433dd0fda9d7dd181a4c28e4dccb83"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}

```
