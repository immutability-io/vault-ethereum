## Vault Ethereum API

Vault provides a CLI that wraps the Vault REST interface. Any HTTP client (including the Vault CLI) can be used for accessing the API. Since the REST API produces JSON, I use the wonderful [jq](https://stedolan.github.io/jq/) for the examples.

* [List Accounts](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#list-accounts)
* [Read Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#read-account)
* [Read Account Balance](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#read-account-balance)
* [Create Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#create-account)
* [Update Account/Re-Encrypt Keystore](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#update-accountre-encrypt-keystore)
* [Delete Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#delete-account)
* [Import Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#import-account)
* [Export Account](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#export-account)
* [Deploy Ethereum Contract](https://github.com/immutability-io/vault-ethereum/blob/master/API.md#deploy-ethereum-contract)
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

This endpoint will list details about the Ethereum account at a path.

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

The example below shows output for a read of `/ethereum/accounts/test`.

```
{
  "request_id": "fe52ec63-80a4-08f5-3780-ac8bd68a8450",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x3943FF61FF803316cF02938b5b0b3Ba3bbE183e4",
    "blacklist": null,
    "chain_id": "4",
    "rpc_url": "http://localhost:8545",
    "whitelist": [
      "0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8",
      "0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572"
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### READ ACCOUNT BALANCE

This endpoint will list the current balance of the Ethereum account at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/accounts/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to read. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/accounts/test/balance | jq .
```

#### Sample Response

The example below shows output for a read of `/ethereum/accounts/test`.

```
{
  "request_id": "018a03db-5560-acc5-f1d8-d568521dcff0",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x3943FF61FF803316cF02938b5b0b3Ba3bbE183e4",
    "pending_balance": "0",
    "pending_nonce": "0",
    "pending_tx_count": "0"
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
* `whitelist` (`string array: <optional>`) - Comma delimited list of allowed accounts.
* `blacklist` (`string array: <optional>`) - Comma delimited list of disallowed accounts. Note: `blacklist` overrides `whitelist`.

#### Sample Payload

```
{
  "rpc_url": "localhost:8545",
  "chain_id": "1977",
  "whitelist": ["0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8","0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572"]
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

The example below shows output for the successful creation of `/ethereum/accounts/test3`.

```
{
  "request_id": "8bfbe4f9-5f8b-1599-27da-172b04c5b8df",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0xb7633a740Df793CbF7530b251c89aecA4F4df748",
    "blacklist": null,
    "chain_id": "1977",
    "rpc_url": "localhost:8545",
    "whitelist": [
      "0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8",
      "0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572"
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### UPDATE ACCOUNT

This endpoint will re-encrypt the keystore for an Ethereum account with a new passphrase.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `PUT`  | `:mount-path/accounts/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to update. This is specified as part of the URL.
* `rpc_url` (`string: <optional> default:"localhost:8545"`) - Specifies the URL of the 'geth' node.
* `chain_id` (`string: <optional> default:"4"`) - Specifies the Ethereum network. Defaults to Rinkeby.
* `whitelist` (`string array: <optional>`) - Comma delimited list of allowed accounts.
* `blacklist` (`string array: <optional>`) - Comma delimited list of disallowed accounts. Note: `blacklist` overrides `whitelist`.

#### Sample Payload

```
{
  "rpc_url": "localhost:8545",
  "chain_id": "1977",
  "whitelist": ["0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8","0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572"]
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

The example below shows output for the successful re-encryption of the keystore for `/ethereum/accounts/test3`.

```
{
  "request_id": "c4d7bae9-269a-d8b9-0171-c6284524c2b5",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0xb7633a740Df793CbF7530b251c89aecA4F4df748",
    "blacklist": null,
    "chain_id": "1977",
    "rpc_url": "localhost:8545",
    "whitelist": [
      "0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8",
      "0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572"
    ]
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
* `passphrase` (`string: <required>`) - The `passphrase` that was used to encrypt the keystore.

#### Sample Payload

Be careful with those passphrases!

```sh
read -s PASSPHRASE; read  PAYLOAD_WITH_PASSPHRASE <<EOF
{"path":"/Users/immutability/.ethereum/keystore/UTC--2017-12-01T23-13-37.315592353Z--a152e7a09267bcff6c33388caab403b76b889939", "passphrase":"$PASSPHRASE"}
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

The example below shows output for the successful creation of `/ethereum/accounts/test3`.

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

### EXPORT ACCOUNT

This endpoint will export a JSON Keystore for use in another wallet.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:name/export`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to export. This is specified as part of the URL.
* `path` (`string: <required>`) - The directory where the JSON keystore file will be exported to.

#### Sample Payload

```sh
{
  "path":"/Users/immutability/.ethereum/keystore"
}
```
#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test/export | jq .
```

#### Sample Response

The example below shows output for the successful export of the keystore for `/ethereum/accounts/test`.

```
{
  "request_id": "9443b8cf-9bde-0790-5b5f-1a01e14629bc",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "passphrase": "synthesis-augmented-playhouse-squeeze-reapply-curry-sprite-surround-coleslaw",
    "path": "/Users/immutability/.ethereum/keystore/UTC--2018-02-02T00-19-34.618912520Z--060b8e95956b8e0423b011ea496e69eec0db136f"
  },
  "warnings": null
}
```


### DEPLOY ETHEREUM CONTRACT

This endpoint will sign a provided Ethereum contract.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:account_name/contracts/:contract_name`  | `200 application/json` |

#### Parameters

* `account_name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `contract_name` (`string: <required>`) - Specifies the name of the contract. This is specified as part of the URL.
* `transaction_data` (`string: <required>`) - The compiled Ethereum contract.
* `value` (`string: <required>`) - The amount of ether in wei.
* `nonce` (`string: <optional> - defaults to "1"`) - The nonce for the transaction
* `gas_price` (`string: <required>`) - The price in gas for the transaction in wei.
* `gas_limit` (`string: <required>`) - The gas limit for the transaction.

#### Sample Payload

```
{
  "transaction_data": "6060604052341561000f57600080fd5b60d38061001d6000396000f3006060604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c14606e575b600080fd5b3415605857600080fd5b606c60048080359060200190919050506094565b005b3415607857600080fd5b607e609e565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a72305820d4b4961183894cf1196bcafbbe4d2573a925296dff82a9dcbc0e8bd8027b153f0029",
  "value":"10000000000",
  "gas_limit":"1500000",
  "gas_price":"21000000000",
  "nonce":"1"

}
```

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test6/contracts/helloworld | jq .
```

#### Sample Response

The example below shows output for the successful deployment of a contract by the account at `/ethereum/accounts/test6`.

```
{
  "request_id": "af4a743e-73ea-ddbd-dac1-351303ac8430",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "tx_hash": "0x5edffe3d8e1c43dff0d17f720219721582e16bd82ddfe4d3c9b7e70cefb968d3"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}

```

### READ ETHEREUM CONTRACT

This endpoint will sign a provided Ethereum contract.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/accounts/:account_name/contracts/:contract_name`  | `200 application/json` |

#### Parameters

* `account_name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `contract_name` (`string: <required>`) - Specifies the name of the contract. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/accounts/test6/contracts/helloworld | jq .
```

**NOTE**: If the transaction hasn't been included in a block yet, the contract address will show as: `Receipt not available`

#### Sample Response

```
{
  "request_id": "da4be9f4-b9fd-90c2-b981-80553cc2359a",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0xCA3986C32beaD6c434773CD41107537f7dDe0c98",
    "tx_hash": "0x62fd378e374ea1166ccb2087ffca49cf1ffcb5162ff3a9651c5b77a781fdfeab"
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
* `raw` (`boolean: <optional>- defaults to false`) - if true, data is expected to be raw hashed transaction data in hex encoding; otherwise data is treated as regular text

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
    https://localhost:8200/v1/ethereum/accounts/test2/sign | jq .
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


### VERIFY SIGNATURE

This endpoint will verify that this account signed some data.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:name/verify`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `data` (`string: <required>`) - Some data.
* `raw` (`boolean: <optional>- defaults to false`) - if true, data is expected to be raw hashed transaction data in hex encoding; otherwise data is treated as regular text
* `signature` (`string: <required>`) - The signature to verify.

#### Sample Payload

```sh

{
  "data": "this is very important"
  "signature": "0xdb6f22f068ae23473beb9b71bb1a2df64a71cb2e51fc43d67558ba8934da572d49b3faa8da387703870474c92beb8c53e89bbd02ba2356b5fc8fa5b342d8fb7b00"
}
```

#### Sample Request

```sh
$ curl -s --cacert /etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test/verify | jq .
```

#### Sample Response

The example below shows output for the successful verification of a signature created by `/ethereum/accounts/test`.

```
{
  "request_id": "0862028d-0810-5355-a583-de002477b26a",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "verified": false
  },
  "warnings": null
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
* `gas_price` (`string: <optional>`) - The price in gas for the transaction. If omitted, we will use the suggested gas price.
* `gas_limit` (`string: <optional>`) - The gas limit for the transaction. If omitted, we will estimate the gas limit.

#### Sample Payload

The following sends 10 ETH to `0xa152E7a09267bcFf6C33388cAab403b76B889939`.

```sh

{
  "amount":"100000000000",
  "to": "0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572"
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
  "request_id": "ac79079d-9e8c-e340-b718-fe19a27ff914",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "from_address": "0x060B8e95956b8E0423b011ea496e69EeC0db136F",
    "to_address": "0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572",
    "tx_hash": "0x536e6ed12214886fa546baa8d72c67cdf45b8bc07d42676b794b474d021a43ff"
  },
  "warnings": null
}
```
