![Immutability](/docs/tagline.png?raw=true "Changes Everything")

## Vault Ethereum API

Vault provides a CLI that wraps the Vault REST interface. Any HTTP client (including the Vault CLI) can be used for accessing the API. Since the REST API produces JSON, I use the wonderful [jq](https://stedolan.github.io/jq/) for the examples.


&nbsp;&nbsp;&nbsp;&nbsp;`└── <MOUNT> `&nbsp;&nbsp;([install](./README.md#install-plugin))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── config `&nbsp;&nbsp;([create](./API.md#create-config), [update](./API.md#update-config), [read](./API.md#read-config))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── accounts `&nbsp;&nbsp;([list](./API.md#list-accounts))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME> `&nbsp;&nbsp;([create](./API.md#create-account), [update](./API.md#update-account), [read](./API.md#read-account), [delete](./API.md#delete-account))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── debit `&nbsp;&nbsp;([update](./API.md#debit-account))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── sign `&nbsp;&nbsp;([update](./API.md#sign))
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── sign-tx `&nbsp;&nbsp;([update](./API.md#sign-tx))
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── transfer `&nbsp;&nbsp;([update](./API.md#transfer))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── verify `&nbsp;&nbsp;([update](./API.md#verify))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── addresses `&nbsp;&nbsp;([list](./API.md#list-addresses))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <ADDRESS> `&nbsp;&nbsp;([read](./API.md#read-address))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── balance `&nbsp;&nbsp;([update](./API.md#balance-by-address))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── verify `&nbsp;&nbsp;([update](./API.md#verify-by-address))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── block `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NUMBER> `&nbsp;&nbsp;([read](./API.md#read-block))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── transactions `&nbsp;&nbsp;([read](./API.md#read-block-transactions))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── convert `&nbsp;&nbsp;([update](./API.md#convert))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── export `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME> `&nbsp;&nbsp;([create](./API.md#export))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── deploy `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME> `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── contracts `&nbsp;&nbsp;([list](./API.md#list-contracts))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │           └── <NAME> `&nbsp;&nbsp;([create](./API.md#deploy-contract), [read](./API.md#read-contract), [delete](./API.md#delete-contract))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── import `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └── <NAME>  `&nbsp;&nbsp;([create](./API.md#import))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    ├── names `&nbsp;&nbsp;([list](./API.md#list-names))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │   └──  <NAME> `&nbsp;&nbsp;([read](./API.md#read-name))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       ├── balance `&nbsp;&nbsp;([update](./API.md#balance-by-name))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    │       └── verify `&nbsp;&nbsp;([update](./API.md#verify-by-name))  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`    └── transaction `  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`        └── <TRANSACTION_HASH> `&nbsp;&nbsp;([read](./API.md#read-transaction))  

### Swagger Doc:

http://petstore.swagger.io/?url=https://raw.githubusercontent.com/zambien/vault-ethereum/got_that_swagger/swagger.json


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

This endpoint will update an account's constraints.

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
    "balance_in_usd": 239.053494011038634,
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
* `send` (`bool: <optional - defaults to true>`) - Indicates whether the transaction should be sent to the network. 

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
  "request_id": "b921207e-c0d9-a3c1-442b-ef8b1884238d",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "amount": "100000000000000000",
    "amount_in_usd": "0",
    "address_from": "0x4169c9508728285e8a9f7945d08645bb6b3576e5",
    "address_to": "0x8AC5e6617F65c071f6dD5d7bD400bf4a46434D41",
    "gas_limit": "21000",
    "gas_price": "1000000000",
    "signed_transaction": "0xf86b06843b9aca00825208948ac5e6617f65c071f6dd5d7bd400bf4a46434d4188016345785d8a0000802ca0ff3fccbde1964047db6be33410436a9220c91ea4080b0e14489dc35fbdabd008a0448fe3ec216a639e1b0eb87b0e4b20aab2e5ec46dad4c38cfc81a1c54e309d21",
    "starting_balance": 8460893507395267000,
    "starting_balance_in_usd": "0",
    "total_spend": "100000000000000000",
    "transaction_hash": "0x3a103587ea6bdeee944e5f68f90ed7b1f4c7699236167d1b1d29495b0319fb26"
  },
  "warnings": null
}
```

### LIST CONTRACTS

This endpoint will list all account contracts.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `LIST`  | `:mount-path/deploy/:account_name/contracts`  | `200 application/json` |

#### Parameters

* `path` (`string: <required>`) - Specifies the mount point. This is specified as part of the URL.
* `account_name` (`string: <required>`) - Specifies the name of the account from which to deploy. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request LIST \
    https://localhost:8200/v1/ethereum/deploy/test/contracts | jq .
```

#### Sample Response

The example below shows output for a query path of `ethereum/deploy/test/contracts` when there are 2 contracts at `/ethereum/deploy/test`.

```
{
  "request_id": "3e53a65b-9910-f5b7-1d42-b6bb83fcfc2c",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "keys": [
      "ponzi",
      "cryptomonkeys",
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```


### DEPLOY CONTRACT

This endpoint will sign a provided Ethereum contract.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/deploy/:account_name/contracts/:contract_name`  | `200 application/json` |

#### Parameters

* `account_name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `contract_name` (`string: <required>`) - Specifies the name of the contract. This is specified as part of the URL.
* `transaction_data` (`string: <required>`) - The compiled Ethereum contract.
* `amount` (`string: <required>`) - The amount of ether in wei to fund the contract with.
* `nonce` (`string: <optional> - defaults to "1"`) - The nonce for the transaction
* `gas_price` (`string: <required>`) - The price in gas for the transaction in wei.
* `gas_limit` (`string: <required>`) - The gas limit for the transaction.
* `send` (`bool: <optional - defaults to true>`) - Indicates whether the transaction should be sent to the network. 

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
    https://localhost:8200/v1/ethereum/deploy/test6/contracts/helloworld | jq .
```

#### Sample Response

The example below shows output for the successful deployment of a contract by the account at `/ethereum/deploy/test6/contracts/helloworld`.

```
{
  "request_id": "af4a743e-73ea-ddbd-dac1-351303ac8430",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "amount": "10000000000",
    "amount_in_usd": "0",
    "total_spend": "10000000000",
    "transaction_hash": "0xee051ae8e9a5afefe94853254de2ea512d88dc4a455334a8e286464c0fa9e767"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}

```

### READ CONTRACT

This endpoint will return the address of an Ethereum contract (if available.)

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/deploy/:account_name/contracts/:contract_name`  | `200 application/json` |

#### Parameters

* `account_name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `contract_name` (`string: <required>`) - Specifies the name of the contract. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/deploy/test6/contracts/helloworld | jq .
```

**NOTE**: If the transaction hasn't been included in a block yet, the contract address will show as: `receipt not available`

#### Sample Response

```
{
  "request_id": "da4be9f4-b9fd-90c2-b981-80553cc2359a",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0xCA3986C32beaD6c434773CD41107537f7dDe0c98",
    "transaction_hash": "0xee051ae8e9a5afefe94853254de2ea512d88dc4a455334a8e286464c0fa9e767"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### DELETE CONTRACT
### SIGN

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
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test2/sign | jq .
```

#### Sample Response

The example below shows output for the successful signing of some data by the private key associated with  `/ethereum/accounts/test2`.

```
{
  "request_id": "d99b7948-453a-67c2-9111-178e6a731812",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x36d1f896e55a6577c62fdd6b84fbf74582266700",
    "signature": "0x90a8712c948b5dfe412ca7e2082be9ef6ddf318a9aaf9183b702c0d1ee180d9d1f97683cb52026dc0de0b6033237cf421a27e88e7d0e608ac4778a9dcfd8818000"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### SIGN-TX

This endpoint will sign transaction.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:name/sign-tx`  | `200 application/json` |

#### Parameters
* `name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `nonce` (`string: <required>`) - The nonce for the transaction.
* `address_to` (`string: <required>`) - The address of the account to send ETH to.
* `value` (`string: <required>`) - Value of ETH (in wei).
* `tx_data` (`string: <optional>`) - Transaction data in HEX string.
* `gas_limit` (`string: <optional>`) - The gas limit for the transaction - defaults to 21000.
* `gas_price` (`string: <optional>`) - The gas price for the transaction in wei. Default - 0.
* `chain_id` (`string: <required>`) - Specifies the Ethereum network.

#### Sample Payload

```json

{
  "nonce":"1",
  "gas_price":"1000000000000000000",
  "gas_limit":"21000",
  "address_to":"0xF4E6e6fa97E10ddc057c94F501B94C1d24EF85Aa",
  "value":"9000000000000000000",
  "tx_data":"0x",
  "chain_id":3
}

```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/main/sign-tx | jq .
```

#### Sample Response

The example below shows output for the successful signing of transaction by the private key associated with  `/ethereum/accounts/main`.
Signed raw transaction in hex format and ready for broadcasting to network.

```json

{
  "request_id": "61b9b5c8-866f-99d6-97f6-5ee0b6489a35",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x84b1858e0bdd6346db6d0f14b574b9be243cc4da",
    "signed_tx": "0xf86f01880de0b6b3a764000082520894f4e6e6fa97e10ddc057c94f501b94c1d24ef85aa887ce66c50e28400008029a043e16edbcaf7c066372fc7eb665d1fe04cf254303df142ccf64c332c730bfac5a0794263d33fd9ed8770aeb534f6a9de513366cce278589c6d2623845c6f901b84"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}

```

### TRANSFER
### VERIFY

This endpoint will verify that this account signed some data.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:name/verify`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `data` (`string: <required>`) - Some data.
* `signature` (`string: <required>`) - The signature to verify.

#### Sample Payload

```sh

{
  "data": "this is very important"
  "signature": "0x90a8712c948b5dfe412ca7e2082be9ef6ddf318a9aaf9183b702c0d1ee180d9d1f97683cb52026dc0de0b6033237cf421a27e88e7d0e608ac4778a9dcfd8818000"
}
```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/accounts/test/verify | jq .
```

#### Sample Response

The example below shows output for the successful verification of a signature created by `/ethereum/accounts/test`.

```
{
  "request_id": "f806212d-087c-7378-dd85-676701aeabb7",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x36d1f896e55a6577c62fdd6b84fbf74582266700",
    "signature": "0x90a8712c948b5dfe412ca7e2082be9ef6ddf318a9aaf9183b702c0d1ee180d9d1f97683cb52026dc0de0b6033237cf421a27e88e7d0e608ac4778a9dcfd8818000",
    "verified": true
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### LIST ADDRESSES

This endpoint will list all account addresses.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `LIST`  | `:mount-path/addresses`  | `200 application/json` |

#### Parameters

* `path` (`string: <required>`) - Specifies the path of the accounts to list. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request LIST \
    https://localhost:8200/v1/ethereum/addresses | jq .
```

#### Sample Response

The example below shows output for a query path of `/ethereum/addresses/` when there are 3 addresses.

```
{
  "request_id": "3e53a65b-9910-f5b7-1d42-b6bb83fcfc2c",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "keys": [
      "0x36d1f896e55a6577c62fdd6b84fbf74582266700",
      "0x7b715f8748ef586b98d3e7c88f326b5a8f409cd8",
      "0x4169c9508728285e8a9f7945d08645bb6b3576e5"
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### READ ADDRESS

This endpoint will list the names associated with an address.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/names/:address`  | `200 application/json` |

#### Parameters

* `address` (`string: <required>`) - Specifies the address of the account to read. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/addresses/0xb56b2dd44073d87cbac5d4a3655354b3762178ee | jq .
```

#### Sample Response

The example below shows output for a read of `/ethereum/addresses/0xb56b2dd44073d87cbac5d4a3655354b3762178ee`.

```
{
  "request_id": "087d361d-127b-a277-d023-283208f62743",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "names": [
      "muchwow"
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### BALANCE BY ADDRESS

This endpoint will return the balance for an address.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/addresses/:address/balance`  | `200 application/json` |

#### Parameters

* `address` (`string: <required>`) - Specifies the address of the account. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/mainnet/addresses/0x4169c9508728285e8a9f7945d08645bb6b3576e5/balance | jq .
```

#### Sample Response

The example below shows output for the successful verification of a signature created by `/ethereum/accounts/test`.

```
{
  "request_id": "d772820e-bbc1-acd1-2aed-1107e72a857e",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x4169c9508728285e8a9f7945d08645bb6b3576e5",
    "balance": "10000000000000000",
    "balance_in_usd": "2.99183586901"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### VERIFY BY ADDRESS

This endpoint will verify that this account signed some data.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/addresses/:address/verify`  | `200 application/json` |

#### Parameters

* `address` (`string: <required>`) - Specifies the address of the account to use for signing. This is specified as part of the URL.
* `data` (`string: <required>`) - Some data.
* `signature` (`string: <required>`) - The signature to verify.

#### Sample Payload

```sh

{
  "data": "this is very important"
  "signature": "0x90a8712c948b5dfe412ca7e2082be9ef6ddf318a9aaf9183b702c0d1ee180d9d1f97683cb52026dc0de0b6033237cf421a27e88e7d0e608ac4778a9dcfd8818000"
}
```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/addresses/0xb56b2dd44073d87cbac5d4a3655354b3762178ee/verify | jq .
```

#### Sample Response

The example below shows output for the successful verification of a signature created by `/ethereum/accounts/test`.

```
{
  "request_id": "f806212d-087c-7378-dd85-676701aeabb7",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0xb56b2dd44073d87cbac5d4a3655354b3762178ee",
    "signature": "0x90a8712c948b5dfe412ca7e2082be9ef6ddf318a9aaf9183b702c0d1ee180d9d1f97683cb52026dc0de0b6033237cf421a27e88e7d0e608ac4778a9dcfd8818000",
    "verified": true
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### READ BLOCK

This endpoint will read details associated with a block.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/block/:number`  | `200 application/json` |

#### Parameters

* `number` (`string: <required>`) - Specifies the number of the block to read. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/block/2800568 | jq .
```

#### Sample Response

The example below shows output for a read of `/ethereum/block/2800568`.

```
{
  "request_id": "5ac451bd-d66a-5f03-9541-1ccecde26223",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "block": 2800568,
    "block_hash": "0x00890448fbdc2000e3e70a66d0b9ac8eaa8d18606512af268e2542ca9d550e3d",
    "difficulty": 2,
    "time": 1534078628,
    "transaction_count": 13
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### READ BLOCK TRANSACTIONS

This endpoint will read transactions associated with a block.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/block/:number/transactions`  | `200 application/json` |

#### Parameters

* `number` (`string: <required>`) - Specifies the number of the block to read. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/block/2800568/transactions | jq .

```

#### Sample Response

The example below shows output for a read of `/ethereum/block/2800568/transactions`.

```
{
  "request_id": "67f5e340-da54-aa64-bbab-598483c9f213",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "0x0c10dc4158ddb4b3464b3519695661a3b73556bcd1a6d5833aa46d0275749e39": {
      "address_to": "0xc0e15E11306334258d61fEe52a22D15e6c9C59e0",
      "gas": 3000000,
      "gas_price": 1000000000,
      "nonce": 73902,
      "value": "0"
    },
    "0x3c45da76247826473e2227b4a8eb3a9fb45280c05b4ecef1d6993087e91cfb58": {
      "address_to": "0xc51b536AD6169bb7F8c893FF4b744d03433455cB",
      "gas": 277468,
      "gas_price": 1000000000,
      "nonce": 134957,
      "value": "0"
    },
    "0x3d26e1d8470d85556fd6a04ec832b5e8cfcbfb306b8ebb94bf030996ca3c7cd4": {
      "address_to": "0xc51b536AD6169bb7F8c893FF4b744d03433455cB",
      "gas": 373558,
      "gas_price": 1000000000,
      "nonce": 134958,
      "value": "0"
    },
    "0x4773f64900d9ae460c72ca478fe2d122027b9dbc0433a40deb62d48ca068b108": {
      "address_to": "0x97e3bA6cC43b2aF2241d4CAD4520DA8266170988",
      "gas": 1200000,
      "gas_price": 1000000000,
      "nonce": 27861,
      "value": "0"
    },
    "0x4fcb5de0ad524d63df206593fd0c0b6e74d9d62919685b2ad6e3bbdeb9753c2e": {
      "address_to": "0x58dcf18084A320670F9abb059312AE60610bda58",
      "gas": 4500000,
      "gas_price": 1000000000,
      "nonce": 13537,
      "value": "0"
    },
    "0x5aa13eba92ebed85199c7b563cca567dc2ab927cff01f6785274f4174acc93a8": {
      "address_to": "0x17dA6A8B86578CEC4525945A355E8384025fa5Af",
      "gas": 100000,
      "gas_price": 1000000000,
      "nonce": 277,
      "value": "1000000000000000000"
    },
    "0x718fa59d7ef2bc9f2ee9642d9357449509ce4e6f14a8e2d0c820d7cfc7dc7535": {
      "address_to": "0xc0e15E11306334258d61fEe52a22D15e6c9C59e0",
      "gas": 3000000,
      "gas_price": 1000000000,
      "nonce": 73963,
      "value": "0"
    },
    "0xa173a5237e034fc54638a1dcf10343f29c0dc2dd8b8b1a25a7d370539e183bd9": {
      "address_to": "0xc51b536AD6169bb7F8c893FF4b744d03433455cB",
      "gas": 377718,
      "gas_price": 1000000000,
      "nonce": 52914,
      "value": "0"
    },
    "0xa737cb951ab32c1b71a521809bab2f4b34234d27abd0ac0c6fac29d480ed6748": {
      "address_to": "0x08Fe64C8B476c9cC776d8A0Ad5bB50e29D83A970",
      "gas": 3000000,
      "gas_price": 1000000000,
      "nonce": 340,
      "value": "0"
    },
    "0xaa322a7ebc3225029ef24a9fde80b2869deacc0a64b14533d07e52fff4542141": {
      "address_to": "0xcB912023AaEB5057BeDB13c937E0519cED0D627A",
      "gas": 400000,
      "gas_price": 1000000000,
      "nonce": 32859,
      "value": "0"
    },
    "0xb11f44837f56205f6695c258c7e25ce6dc81f3c746d55ece8a22375147cfa34f": {
      "address_to": "0x39394Ad63206DE2cd86881e1D20ff1621566e482",
      "gas": 500000,
      "gas_price": 1000000000,
      "nonce": 22685,
      "value": "0"
    },
    "0xc4c44fc474f2aa27e9ac651fdf3a109984c464a5859db47d5cb2cfdc93f52788": {
      "address_to": "0xc51b536AD6169bb7F8c893FF4b744d03433455cB",
      "gas": 275388,
      "gas_price": 1000000000,
      "nonce": 70743,
      "value": "0"
    },
    "0xd64b823fcb714b6811c9790e7219ef1836e5568e5bdf2e7385a1b910178c7de8": {
      "address_to": "0xCFe52FEDF5fc3b92ABA3D43b96B4Ae1d0b39062c",
      "gas": 3000000,
      "gas_price": 1000000000,
      "nonce": 14,
      "value": "0"
    }
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### CREATE CONFIG

This endpoint configure the plugin at a mount.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/config`  | `200 application/json` |

#### Parameters

* `rpc_url` (`string: <optional> default:"https://rinkeby.infura.io"`) - Specifies the RPC URL of the Ethereum node.
* `chain_id` (`string: <optional> default:"4"`) - Specifies the Ethereum network. Defaults to Rinkeby.
* `bound_cidr_list` (`string array: <optional>`) - Comma delimited list of allowed CIDR blocks.
* `api_key` (`string: <optional>`) - The Infura API key.

#### Sample Payload

```

{
  "chain_id": "4",
  "rpc_url": "https://rinkeby.infura.io"
}

```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @config.json \
    https://localhost:8200/v1/ethereum/config | jq .
```

#### Sample Response

The example below shows output for the successful creation of a mount's configuration.

```
{
  "request_id": "086c19b4-e07a-55ba-5ba3-3819b9c0b1da",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "api_key": "",
    "bound_cidr_list": null,
    "chain_id": "4",
    "rpc_url": "https://rinkeby.infura.io"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### UPDATE CONFIG

This endpoint reconfigures the plugin at a mount.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `PUT`  | `:mount-path/config`  | `200 application/json` |

#### Parameters

* `rpc_url` (`string: <optional> default:"https://rinkeby.infura.io"`) - Specifies the RPC URL of the Ethereum node.
* `chain_id` (`string: <optional> default:"4"`) - Specifies the Ethereum network. Defaults to Rinkeby.
* `bound_cidr_list` (`string array: <optional>`) - Comma delimited list of allowed CIDR blocks.
* `api_key` (`string: <optional>`) - The Infura API key.

#### Sample Payload

```

{
  "chain_id": "4",
  "rpc_url": "https://rinkeby.infura.io"
}

```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request PUT \
    --data @config.json \
    https://localhost:8200/v1/ethereum/config | jq .
```

#### Sample Response

The example below shows output for the successful update of a mount's configuration.

```
{
  "request_id": "086c19b4-e07a-55ba-5ba3-3819b9c0b1da",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "api_key": "",
    "bound_cidr_list": null,
    "chain_id": "4",
    "rpc_url": "https://rinkeby.infura.io"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### READ CONFIG

This endpoint returns the configuration of a plugin at a mount.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `PUT`  | `:mount-path/config`  | `200 application/json` |

#### Parameters

None


#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/config | jq .
```

#### Sample Response

The example below shows output for the mount's configuration.

```
{
  "request_id": "086c19b4-e07a-55ba-5ba3-3819b9c0b1da",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "api_key": "",
    "bound_cidr_list": null,
    "chain_id": "4",
    "rpc_url": "https://rinkeby.infura.io"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### CONVERT

This endpoint will convert one Ethereum unit to another.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/convert`  | `200 application/json` |

#### Parameters

* `amount_from` (`string: <required>`) - Specifies amount to convert.
* `unit_from` (`string: <required>`) - Specifies unit to convert from.
* `unit_to` (`string: <required>`) - Specifies unit to convert to.

#### Sample Payload

```sh
{
  "unit_from": "wei",
  "unit_to": "eth",
  "amount": "200000000000000000"
}
```
#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request PUT \
    --data @convert.json \
    https://localhost:8200/v1/ethereum/convert | jq .
```

#### Sample Response

The example below shows output for the successful export of the keystore for `/ethereum/accounts/test`.

```
{
  "request_id": "125314b7-7d7d-5c30-daa2-05d3680a68ea",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "amount_from": "200000000000000000",
    "amount_to": "0.2",
    "unit_from": "wei",
    "unit_to": "ether"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

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
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
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

This endpoint will list all account names.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `LIST`  | `:mount-path/names`  | `200 application/json` |

#### Parameters

* `path` (`string: <required>`) - Specifies the path of the accounts to list. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request LIST \
    https://localhost:8200/v1/ethereum/names | jq .
```

#### Sample Response

The example below shows output for a query path of `/ethereum/names/` when there are 3 addresses.

```
{
  "request_id": "3e53a65b-9910-f5b7-1d42-b6bb83fcfc2c",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "keys": [
      "test",
      "test3",
      "lesswow"
    ]
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### READ NAME

This endpoint will list the addreses associated with a named account.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/names/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to read. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/names/muchwow | jq .
```

#### Sample Response

The example below shows output for a read of `/ethereum/names/muchwow`.

```
{
  "request_id": "57556ed7-da99-ee4d-fbf0-2feaed17e5b9",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0xb56b2dd44073d87cbac5d4a3655354b3762178ee"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```


### BALANCE BY NAME

This endpoint will return the balance for an address.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/names/:name/balance`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/mainnet/addresses/immutability/balance | jq .
```

#### Sample Response

The example below shows output for the successful verification of a signature created by `/ethereum/accounts/test`.

```
{
  "request_id": "d772820e-bbc1-acd1-2aed-1107e72a857e",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x4169c9508728285e8a9f7945d08645bb6b3576e5",
    "balance": "10000000000000000",
    "balance_in_usd": "2.99183586901"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### VERIFY BY NAME

This endpoint will verify that this account signed some data.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/names/:name/verify`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the account to use for signing. This is specified as part of the URL.
* `data` (`string: <required>`) - Some data.
* `signature` (`string: <required>`) - The signature to verify.

#### Sample Payload

```sh

{
  "data": "this is very important"
  "signature": "0x90a8712c948b5dfe412ca7e2082be9ef6ddf318a9aaf9183b702c0d1ee180d9d1f97683cb52026dc0de0b6033237cf421a27e88e7d0e608ac4778a9dcfd8818000"
}
```

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request POST \
    --data @payload.json \
    https://localhost:8200/v1/ethereum/name/test/verify | jq .
```

#### Sample Response

The example below shows output for the successful verification of a signature created by `/ethereum/accounts/test`.

```
{
  "request_id": "f806212d-087c-7378-dd85-676701aeabb7",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address": "0x36d1f896e55a6577c62fdd6b84fbf74582266700",
    "signature": "0x90a8712c948b5dfe412ca7e2082be9ef6ddf318a9aaf9183b702c0d1ee180d9d1f97683cb52026dc0de0b6033237cf421a27e88e7d0e608ac4778a9dcfd8818000",
    "verified": true
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

### READ TRANSACTION  

This endpoint will read details associated with a transaction hash.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/transaction/:transaction_hase`  | `200 application/json` |

#### Parameters

* `transaction_hase` (`string: <required>`) - Specifies the transaction hash to read. This is specified as part of the URL.

#### Sample Request

```sh
$ curl -s --cacert ~/etc/vault.d/root.crt --header "X-Vault-Token: $VAULT_TOKEN" \
    --request GET \
    https://localhost:8200/v1/ethereum/transaction/0x4773f64900d9ae460c72ca478fe2d122027b9dbc0433a40deb62d48ca068b108 | jq .
```

#### Sample Response

The example below shows output for a read of `/ethereum/transaction/0x4773f64900d9ae460c72ca478fe2d122027b9dbc0433a40deb62d48ca068b108`.

```
{
  "request_id": "fc524980-a91b-a1f8-2815-ba1aa67aa70f",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "address_from": "0x674647242239941b2D35368e66A4EdC39b161Da9",
    "address_to": "0x97e3bA6cC43b2aF2241d4CAD4520DA8266170988",
    "gas": 1200000,
    "gas_price": 1000000000,
    "nonce": 27861,
    "pending": false,
    "receipt_status": 1,
    "transaction_hash": "0x4773f64900d9ae460c72ca478fe2d122027b9dbc0433a40deb62d48ca068b108",
    "value": "0"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```


### Swagger local setup

If you are using a local vault you must trust the cert for the local vault.  For example on a Mac assuming you followed the Vault setup instructions in Immutability:

`sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/etc/vault.d/root.crt`

You must also allow CORS in Vault.  For example:
`vault write sys/config/cors enabled=true allowed_origins="*"`


```
# swagger built binary is no good for golang 1.11 so we will build and install it
go get -u github.com/go-swagger/go-swagger/cmd/swagger
cd $GOPATH/src/github.com/go-swagger/go-swagger/cmd/swagger \
    && git checkout 0.16.0 \ # because master is broke
    && go build \
    && sudo cp swagger /usr/local/bin/swagger \
    && cd -
    
# generate and serve swagger json    
swagger generate spec -o ./swagger.json --scan-models \
  && swagger serve -F=swagger swagger.json
  
# or if you prefer ReDoc
swagger generate spec -o ./swagger.json --scan-models \
  && swagger serve swagger.json      
```    