![Immutability](/docs/tagline.png?raw=true "Changes Everything")

## Vault Ethereum API

Vault provides a CLI that wraps the Vault REST interface. Any HTTP client (including the Vault CLI) can be used for accessing the API. Since the REST API produces JSON, I use the wonderful [jq](https://stedolan.github.io/jq/) for the examples.

&nbsp;&nbsp;&nbsp;&nbsp;`└── ethereum `&nbsp;&nbsp;([install](./README.md#install-plugin))  
&nbsp;&nbsp;&nbsp;&nbsp;`    ├── accounts `&nbsp;&nbsp;([list](./API.md#list-accounts))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   ├── <NAME> `&nbsp;&nbsp;([create](./API.md#create-account), [update](./API.md#update-account), [read](./API.md#read-account), [delete](./API.md#delete-account))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   │   ├── debit `&nbsp;&nbsp;([create](./API.md#debit-account))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   │   ├── contracts `&nbsp;&nbsp;([list](./API.md#list-contracts))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   │   │   └── <NAME> `&nbsp;&nbsp;([create](./API.md#deploy-contract), [read](./API.md#read-contract), [delete](./API.md#delete-contract))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   │   ├── sign `&nbsp;&nbsp;([create](./API.md#sign))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   │   ├── transfer `&nbsp;&nbsp;([create](./API.md#transfer))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   │   └── verify `&nbsp;&nbsp;([create](./API.md#verify))  
&nbsp;&nbsp;&nbsp;&nbsp;`    ├── addresses `&nbsp;&nbsp;([list](./API.md#list-addresses))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   ├── <ADDRESS> `&nbsp;&nbsp;([read](./API.md#read-address))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   │   └── verify `&nbsp;&nbsp;([create](./API.md#verify-by-address))  
&nbsp;&nbsp;&nbsp;&nbsp;`    ├── block `  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NUMBER> `&nbsp;&nbsp;([read](./API.md#read-block))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── transactions `&nbsp;&nbsp;([read](./API.md#read-block-transactions))  
&nbsp;&nbsp;&nbsp;&nbsp;`    ├── config `&nbsp;&nbsp;([create](./API.md#create-config), [update](./API.md#update-config), [read](./API.md#read-config))  
&nbsp;&nbsp;&nbsp;&nbsp;`    ├── convert `&nbsp;&nbsp;([update](./API.md#convert))  
&nbsp;&nbsp;&nbsp;&nbsp;`    ├── export `  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME> `&nbsp;&nbsp;([create](./API.md#export))  
&nbsp;&nbsp;&nbsp;&nbsp;`    ├── import `  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME>  `&nbsp;&nbsp;([create](./API.md#import))  
&nbsp;&nbsp;&nbsp;&nbsp;`    ├── names `&nbsp;&nbsp;([list](./API.md#list-names))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │   └──  <NAME> `&nbsp;&nbsp;([read](./API.md#read-name))  
&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── verify `&nbsp;&nbsp;([create](./API.md#verify-by-name))  
&nbsp;&nbsp;&nbsp;&nbsp;`    └── transaction `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`        └── <TRANSACTION_HASH> `&nbsp;&nbsp;([read](./API.md#read-transaction))  

### LIST ACCOUNTS

This endpoint will list all accounts stores at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `LIST`  | `:mount-path/accounts`  | `200 application/json` |

#### Parameters

* `mount-path` (`string: <required>`) - Specifies the path of the accounts to list. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request LIST \
    https://localhost:8200/v1/ethereum/accounts | jq .
```

#### Sample Response

The example below shows output for a query path of `/ethereum/accounts/` when there are 4 accounts.

```
{
  "request_id": "6c33d1c9-d599-179c-f20c-307a47e129f4",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "keys": [
      "lesswow",
      "test",
      "muchwow",
      "morewow"
    ]
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
* `spending_limit_tx` (`string: <optional> `) - The total amount of Wei allowed to be spent in a single transaction.
* `spending_limit_total` (`string: <optional>`) - The total amount of Wei allowed to be spent for this account.
* `whitelist` (`string array: <optional>`) - Comma delimited list of allowed accounts.
* `blacklist` (`string array: <optional>`) - Comma delimited list of disallowed accounts. Note: `blacklist` overrides `whitelist`.

#### Sample Payload

```
{
  "whitelist": ["0xD9E025bFb6ef48919D9C1a49834b7BA859714cD8","0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572"]
}
```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test3 | jq .
```

#### Sample Response

The example below shows output for the successful creation of `/ethereum/accounts/test3`.

```
{
  "request_id": "af43bafd-5afa-9577-ffec-5fb7d1316bfa",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0xd850d887f803b64a18c98a3acfbe0496f02fe9f5",
    "blacklist": null,
    "spending_limit_total": "0",
    "spending_limit_tx": "0",
    "total_spend": "0",
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

This endpoint will update an accounts constraints.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `PUT`  | `:mount-path/accounts/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to create. This is specified as part of the URL.
* `spending_limit_tx` (`string: <optional> `) - The total amount of Wei allowed to be spent in a single transaction.
* `spending_limit_total` (`string: <optional>`) - The total amount of Wei allowed to be spent for this account.
* `whitelist` (`string array: <optional>`) - Comma delimited list of allowed accounts.
* `blacklist` (`string array: <optional>`) - Comma delimited list of disallowed accounts. Note: `blacklist` overrides `whitelist`.

#### Sample Payload

```
{
  "whitelist": ["0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572"]
}
```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @update.json \
    https://localhost:8200/v1/ethereum/accounts/test3 | jq .
```

#### Sample Response

The example below shows output for the successful creation of `/ethereum/accounts/test3`.

```
{
  "request_id": "28286c4a-f40e-4842-c52c-7970571c9ddd",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0xd850d887f803b64a18c98a3acfbe0496f02fe9f5",
    "blacklist": null,
    "spending_limit_total": "0",
    "spending_limit_tx": "0",
    "total_spend": "0",
    "whitelist": [
      "0x58e9043a873EdBa4c5C865Bf1c65dcB3473f7572"
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
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/accounts/muchwow | jq .
```

#### Sample Response

The example below shows output for a read of `/ethereum/accounts/test`.

```
{
  "request_id": "58c331ea-a63d-b9ea-75d7-bc9e1933708c",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x7b715f8748ef586b98d3e7c88f326b5a8f409cd8",
    "balance": 799958000000000000,
    "blacklist": null,
    "spending_limit_total": "",
    "spending_limit_tx": "",
    "total_spend": "",
    "whitelist": null
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
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request DELETE \
    https://localhost:8200/v1/ethereum/accounts/test3
```

#### Sample Response

There is no response payload.

### DEBIT ACCOUNT


This endpoint will debit an Ethereum account.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:name/debit`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `address_to` (`string: <required>`) - A Hex string specifying the Ethereum address to send the ETH to.
* `amount` (`string: <required>`) - The amount of ether - in wei.
* `gas_price` (`string: <optional>`) - The price in gas for the transaction. If omitted, we will use the suggested gas price.
* `gas_limit` (`string: <optional>`) - The gas limit for the transaction. If omitted, we will estimate the gas limit.

#### Sample Payload

The following sends 0.2 ETH to `0x36D1F896E55a6577C62FDD6b84fbF74582266700`.

```sh

{
  "amount":"200000000000000000",
  "to": "0x36D1F896E55a6577C62FDD6b84fbF74582266700"
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
    "amount": "200000000000000000"
    "from_address": "0x7b715f8748ef586b98d3e7c88f326b5a8f409cd8"
    "gas_limit": "21000"
    "gas_price": "2000000000"
    "balance": "1000000000000000000"
    "to_address": "0x36D1F896E55a6577C62FDD6b84fbF74582266700"
    "total_spend": "200000000000000000"
    "transaction_hash": "0x0b4938a1a44f545deeea500d50761c22bfe2bc006b26be8adf4dcd4fc0597769"
  },
  "warnings": null
}
```

### LIST CONTRACTS

### DEPLOY CONTRACT

This endpoint will sign a provided Ethereum contract.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:account_name/contracts/:contract_name`  | `200 application/json` |

#### Parameters

* `account_name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `contract_name` (`string: <required>`) - Specifies the name of the contract. This is specified as part of the URL.
* `transaction_data` (`string: <required>`) - The compiled Ethereum contract.
* `amount` (`string: <required>`) - The amount of ether in wei to fund the contract with.
* `nonce` (`string: <optional> - defaults to "1"`) - The nonce for the transaction
* `gas_price` (`string: <required>`) - The price in gas for the transaction in wei.
* `gas_limit` (`string: <required>`) - The gas limit for the transaction.

#### Sample Payload

```
{
  "transaction_data": "6060604052341561000f57600080fd5b60d38061001d6000396000f3006060604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c14606e575b600080fd5b3415605857600080fd5b606c60048080359060200190919050506094565b005b3415607857600080fd5b607e609e565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a72305820d4b4961183894cf1196bcafbbe4d2573a925296dff82a9dcbc0e8bd8027b153f0029",
  "amount":"10000000000",
  "gas_limit":"1500000",
  "nonce":"1"

}
```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test6/contracts/helloworld | jq .
```

#### Sample Response

The example below shows output for the successful deployment of a contract by the account at `/ethereum/accounts/test6/contracts/helloworld`.

```
{
  "request_id": "af4a743e-73ea-ddbd-dac1-351303ac8430",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "transaction_hash": "0x5edffe3d8e1c43dff0d17f720219721582e16bd82ddfe4d3c9b7e70cefb968d3"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}

```


### READ CONTRACT
### DELETE CONTRACT
### SIGN
### TRANSFER
### VERIFY
### LIST ADDRESSES
### READ ADDRESS
### VERIFY BY ADDRESS
### READ BLOCK
### READ BLOCK TRANSACTIONS
### CREATE CONFIG
### UPDATE CONFIG
### READ CONFIG
### CONVERT
### EXPORT

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
  "path":"/Users/cypherhat/.ethereum/keystore"
}
```
#### Sample Request

```sh
$$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/export/test | jq .
```

#### Sample Response

The example below shows output for the successful export of the keystore for `/ethereum/accounts/test`.

```
{
  "request_id": "47e3ef56-e2ba-3895-d61b-d9615df1560c",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "passphrase": "routing-explode-slander-satiable-stardom-cope-cranium-upriver-unfold",
    "path": "/Users/cypherhat/.ethereum/keystore/UTC--2018-08-19T20-05-13.985145605Z--36d1f896e55a6577c62fdd6b84fbf74582266700"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### IMPORT

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
{"path":"/Users/cypherhat/.ethereum/keystore/UTC--2017-12-01T23-13-37.315592353Z--a152e7a09267bcff6c33388caab403b76b889939", "passphrase":"$PASSPHRASE"}
EOF
unset PASSPHRASE
```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data $PAYLOAD_WITH_PASSPHRASE \
    https://localhost:8200/v1/ethereum/import/test3 | jq .
    unset PAYLOAD_WITH_PASSPHRASE
```
#### Sample Response

The example below shows output for the successful creation of `/ethereum/accounts/test3`.

```
{
  "request_id": "139fc04e-a4be-7775-7df1-ec4c789ccc53",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x36d1f896e55a6577c62fdd6b84fbf74582266700",
    "blacklist": null,
    "spending_limit_total": "",
    "spending_limit_tx": "",
    "total_spend": "",
    "whitelist": null
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### LIST NAMES
### READ NAME
### VERIFY BY NAME
### READ TRANSACTION  
