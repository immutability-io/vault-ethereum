## HOW TO CONFIGURE THE PLUGIN AFTER IT HAS BEEN ENABLED


### Using Curl

```
curl -X PUT -H "X-Vault-Token: $(vault print token)" -H "X-Vault-Request: true" -d '{"chain_id":"5777","rpc_url":"http://ganache:8545"}' https://localhost:9200/v1/vault-ethereum/config
```

### Using the Vault CLI

```
vault write -format=json vault-ethereum/config  rpc_url='http://ganache:8545' chain_id='5777'

```
### Sample Response


```json
{
  "request_id": "ccc15935-a663-cdf8-cf13-6515c0099da7",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "bound_cidr_list": null,
    "chain_id": "5777",
    "exclusions": null,
    "inclusions": null,
    "rpc_url": "http://ganache:8545"
  },
  "warnings": null
}
```

## HOW TO CREATE AN ACCOUNT NAMED BOB USING A MNEMONIC


### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/bob mnemonic='volcano story trust file before member board recycle always draw fiction when'

```
### Using Curl

```
curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" -d '{"mnemonic":"volcano story trust file before member board recycle always draw fiction when"}' https://localhost:9200/v1/vault-ethereum/accounts/bob
```

### Sample Response


```json
{
  "request_id": "c730cdf2-3b32-b32d-a411-ccdcc351a4b8",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "address": "0x90259301a101A380F7138B50b6828cfFfd4Cbf60",
    "exclusions": null,
    "inclusions": null
  },
  "warnings": null
}
```

### STEP 1: GET BOB'S ADDRESS:
### Using the Vault CLI

```
vault read -field=address vault-ethereum/accounts/bob

```

### Using Curl

```
curl -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" https://localhost:9200/v1/vault-ethereum/accounts/bob
```

### Sample Response


```
0x90259301a101A380F7138B50b6828cfFfd4Cbf60
```

## HOW TO CREATE AN ACCOUNT NAMED ALICE WITH NO MNEMONIC


### Using the Vault CLI

```
vault write -f -format=json vault-ethereum/accounts/alice

```
### Using Curl

```
curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" -d 'null' https://localhost:9200/v1/vault-ethereum/accounts/alice
```

### Sample Response


```json
{
  "request_id": "56dff6db-26e4-a104-4c22-51433eb3ecf8",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "address": "0x310BC85C2593D567C31b919af6F4fb2257bD4aE8",
    "exclusions": null,
    "inclusions": null
  },
  "warnings": null
}
```

## HOW TO TRANSFER 0.5 ETH FROM BOB TO ALICE


### STEP 1: CONVERT 0.5 ETH TO WEI:
### Using the Vault CLI

```
vault write -field=amount_to vault-ethereum/convert unit_from='ETH' unit_to='WEI' amount=0.5

```

### Using Curl

```
curl -X PUT -H "X-Vault-Token: $(vault print token)" -H "X-Vault-Request: true" -d '{"amount":"0.5","unit_from":"ETH","unit_to":"WEI"}' https://localhost:9200/v1/vault-ethereum/convert
```

### Sample Response


```
500000000000000000
```

### STEP 2: GET ALICE'S ADDRESS:
### Using the Vault CLI

```
vault read -field=address vault-ethereum/accounts/alice

```

### Using Curl

```
curl -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" https://localhost:9200/v1/vault-ethereum/accounts/alice
```

### Sample Response


```
0x310BC85C2593D567C31b919af6F4fb2257bD4aE8
```

### STEP 3: SEND BOB'S ETH TO ALICE:
### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/bob/transfer to='0x310BC85C2593D567C31b919af6F4fb2257bD4aE8' amount='500000000000000000'

```
### Using Curl

```
curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" -d '{"amount":"500000000000000000","to":"0x310BC85C2593D567C31b919af6F4fb2257bD4aE8"}' https://localhost:9200/v1/vault-ethereum/accounts/bob/transfer
```

### Sample Response


```json
{
  "request_id": "e8b375a8-ec8c-ab41-0539-9c1df78d6cc3",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "amount": "500000000000000000",
    "from": "0x90259301a101A380F7138B50b6828cfFfd4Cbf60",
    "gas_limit": "21000",
    "gas_price": "20000000000",
    "nonce": "0",
    "signed_transaction": "0xf86c808504a817c80082520894310bc85c2593d567c31b919af6f4fb2257bd4ae88806f05b59d3b20000801ca0a1b2942d128b0a6deaac2b0270e7a3ab61fe2937f2c4719e83b6228949279518a0420ec6e1456cbacfc3562c4629d192963c10f64fa57e5ccef01165e5d284ce6c",
    "to": "0x310BC85C2593D567C31b919af6F4fb2257bD4aE8",
    "transaction_hash": "0x8327d986bc56ab274fef332d7cf05e4da3ab1771e3f2f3891449e9787dae9d53"
  },
  "warnings": null
}
```

## HOW TO CHECK BALANCE OF THE ACCOUNT NAMED BOB


### Using the Vault CLI

```
vault read -format=json vault-ethereum/accounts/bob/balance

```
### Using Curl

```
curl -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" https://localhost:9200/v1/vault-ethereum/accounts/bob/balance
```

### Sample Response


```json
{
  "request_id": "98e57f32-1d3e-f1f9-8e19-7f958acf2d9c",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "address": "0x90259301a101A380F7138B50b6828cfFfd4Cbf60",
    "balance": "99499580000000000000"
  },
  "warnings": null
}
```

## HOW TO CHECK BALANCE OF THE ACCOUNT NAMED ALICE


### Using the Vault CLI

```
vault read -format=json vault-ethereum/accounts/alice/balance

```
### Using Curl

```
curl -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" https://localhost:9200/v1/vault-ethereum/accounts/alice/balance
```

### Sample Response


```json
{
  "request_id": "80bc289b-5a10-d348-a545-565d56ad1657",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "address": "0x310BC85C2593D567C31b919af6F4fb2257bD4aE8",
    "balance": "500000000000000000"
  },
  "warnings": null
}
```

## HOW TO SIGN A TRANSACTION WITH DATA


### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/alice/sign-tx to='0x90259301a101A380F7138B50b6828cfFfd4Cbf60' data='f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772' amount='20000000000000000'

```
### Using Curl

```
curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" -d '{"amount":"20000000000000000","data":"f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772","to":"0x90259301a101A380F7138B50b6828cfFfd4Cbf60"}' https://localhost:9200/v1/vault-ethereum/accounts/alice/sign-tx
```

### Sample Response


```json
{
  "request_id": "086240db-a4a3-87a2-09bb-a6d9edbe6dd9",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "amount": "20000000000000000",
    "from": "0x310BC85C2593D567C31b919af6F4fb2257bD4aE8",
    "gas_limit": "21000",
    "gas_price": "20000000000",
    "nonce": "0",
    "signed_transaction": "0xf9014a808504a817c8008252089490259301a101a380f7138b50b6828cfffd4cbf6087470de4df820000b8de6638366438323032623238343737333539343030383235323038393434353932643866386437623030316537326362323661373365346661313830366135316163373964383830646530623662336137363430303030383032636130353932346264653765663130616138386462396336366464346635666231366234366466663233313962393936386265393833313138623537626235303536326130303162323462333130313030303466313364396132366233323038343532353761366366633262663831396133643535653366633836323633633566303737321ca0dbdc5a672821a88e14d39e8e23d77488073237d8fb3f776df5a08f0636e6813ea0653158e861e0e140ee8bf056f4956a17d599e9636fc5e5838789a309ad283cf3",
    "to": "0x90259301a101A380F7138B50b6828cfFfd4Cbf60",
    "transaction_hash": "0x7d5b42500829391e5a80b97d3ba5a158174e2e9bcf438cc4f12ff97fffc07d2b"
  },
  "warnings": null
}
```

## HOW TO SIGN MESSAGE. Message signature can be verify on: https://etherscan.io/verifySig/2156


### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/bob/sign message='HOLA VAULT'

```
### Using Curl

```
curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" -d '{"message":"HOLA VAULT"}' https://localhost:9200/v1/vault-ethereum/accounts/bob/sign
```

### Sample Response


```json
{
  "request_id": "19709fa8-3c48-ff1c-6a2b-50de07d6952c",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "address": "0x90259301a101a380f7138b50b6828cfffd4cbf60",
    "hashedMessage": "0x916195f5e434a5ac99bda30311516c98e3896c96e8e63e0fb7771584b7bd6b07",
    "signature": "0x79747d2859bb3c7305b33b4dae171c1f7d0f01821cbf0577577956e135c046bf0e53618cb373afd0ab92a7a13e566ebdb94f39d3ec738885d6d7f26d3c12393400"
  },
  "warnings": null
}
```

## DEPLOY CONTRACT FixedSupplyToken


### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/bob/deploy abi=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/FixedSupplyToken.abi bin=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/FixedSupplyToken.bin

```
## DEPLOY CONTRACT Owned


### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/bob/deploy abi=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/Owned.abi bin=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/Owned.bin

```
### Sample Response


```json
{
  "request_id": "d252fe9d-5568-91a1-234c-722be020869a",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "contract": "0xC5bb9B170339E5220b5CA7f4a12fDd08Ae4a172f",
    "from": "0x90259301a101A380F7138B50b6828cfFfd4Cbf60",
    "gas_limit": "292163",
    "gas_price": "20000000000",
    "nonce": "2",
    "signed_transaction": "0xf904a0028504a817c800830475438080b9044d608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506103ed806100606000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806379ba5097146100515780638da5cb5b1461005b578063d4ee1d90146100a5578063f2fde38b146100ef575b600080fd5b610059610133565b005b6100636102d0565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100ad6102f5565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101316004803603602081101561010557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061031b565b005b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461018d57600080fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff166000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461037457600080fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505056fea265627a7a72305820d48bb346f756f1b4c514d0cb6670a26b563bf5d9eb222548a5eb56d86218f01b64736f6c634300050a00321ca093f9b58cb112acec726c1558b1a5eadc13dcdfd9f2767d76d0510dd75050b75ca067ed985ef898fa6885f379cb5ee67d03dfe7a859356e3bdd51131e60433aeddc",
    "transaction_hash": "0xd0729631e73c3861388901c666af8a2b16d09aae303e5093a591edba6d741ff8"
  },
  "warnings": null
}
```

## DEPLOY CONTRACT SafeMath


### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/bob/deploy abi=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/SafeMath.abi bin=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/SafeMath.bin

```
### Sample Response


```json
{
  "request_id": "d2ee52ea-74f5-952d-83c9-255f75ab9879",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "contract": "0x43ce456fA428B4FB5AEA1762221348BBA8FA7097",
    "from": "0x90259301a101A380F7138B50b6828cfFfd4Cbf60",
    "gas_limit": "71714",
    "gas_price": "20000000000",
    "nonce": "3",
    "signed_transaction": "0xf8ca038504a817c800830118228080b87860556023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a72305820a8567376629c47a3133e19c9c636472d6accc6b72662776af08e7395c3adea9164736f6c634300050a00321ca09ae510e20c0d104b416706b5db0c2dd0673fd54f32b904f89049211f5fba67e0a00bed48ff364282f770b2f2eccb241118fb15dd65a7b26e07a5b75e84d9c6a4ce",
    "transaction_hash": "0x81cf7772b118c89bcab6fdec1179b231b2bdb0a0e8b59e821c99cb513ce973dd"
  },
  "warnings": null
}
```

## READ TOKEN "0x2bEeBc3DedbD94A0521eEe7C8D2ae4214EB15b44" BALANCE


### Using the Vault CLI

```
vault read -format=json vault-ethereum/accounts/bob/erc20/totalSupply contract="0x2bEeBc3DedbD94A0521eEe7C8D2ae4214EB15b44"

```
### Sample Response


```json
{
  "request_id": "c9acf4bb-64c3-2372-6ab0-8d34bc2dc05b",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "contract": "0x2bEeBc3DedbD94A0521eEe7C8D2ae4214EB15b44",
    "name": "Example Fixed Supply Token",
    "symbol": "FIXED",
    "total_supply": "1000000"
  },
  "warnings": null
}
```

## DEPLOY CONTRACT FixedSupplyToken


### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/bob/deploy abi=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/FixedSupplyToken.abi bin=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/FixedSupplyToken.bin

```
