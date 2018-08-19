![Immutability](/docs/tagline.png?raw=true "Changes Everything")

## Vault Ethereum API

Vault provides a CLI that wraps the Vault REST interface. Any HTTP client (including the Vault CLI) can be used for accessing the API. Since the REST API produces JSON, I use the wonderful [jq](https://stedolan.github.io/jq/) for the examples.

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
### LIST CONTRACTS
### DEPLOY CONTRACT
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
### IMPORT
### LIST NAMES
### READ NAME
### VERIFY BY NAME
### READ TRANSACTION  
