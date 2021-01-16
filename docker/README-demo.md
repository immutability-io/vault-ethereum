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
  "request_id": "bce6f02f-b468-296e-89d8-de1ddc691789",
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
  "request_id": "4fd1733b-3adb-f66f-5b31-bd99adb2ede2",
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
  "request_id": "227d3bfa-56c1-d0f3-d3ba-e8e60beddf2e",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "address": "0x8b470Eb0Cd23121814805886864b80085C58889d",
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
curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" -d '{"amount":"0.5","unit_from":"ETH","unit_to":"WEI"}' https://localhost:9200/v1/vault-ethereum/convert
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
0x8b470Eb0Cd23121814805886864b80085C58889d
```

### STEP 3: SEND BOB'S ETH TO ALICE:
### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/bob/transfer to='0x8b470Eb0Cd23121814805886864b80085C58889d' amount='500000000000000000'

```
### Using Curl

```
curl -X PUT -H "X-Vault-Request: true" -H "X-Vault-Token: $(vault print token)" -d '{"amount":"500000000000000000","to":"0x8b470Eb0Cd23121814805886864b80085C58889d"}' https://localhost:9200/v1/vault-ethereum/accounts/bob/transfer
```

### Sample Response


```json
{
  "request_id": "9f658fdc-5d10-a7ce-556f-d7f80ad105c4",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "amount": "500000000000000000",
    "from": "0x90259301a101A380F7138B50b6828cfFfd4Cbf60",
    "gas_limit": "21000",
    "gas_price": "20000000000",
    "nonce": "0",
    "signed_transaction": "0xf86c808504a817c800825208948b470eb0cd23121814805886864b80085c58889d8806f05b59d3b20000801ba03c1f34e95053ebf858042bbb014134b42023f53fab8fe0fab5fe463826641fa0a075792ae1ce94d8101d108359defb3c593e18edc66493208b9b1c29a4597cf4a8",
    "to": "0x8b470Eb0Cd23121814805886864b80085C58889d",
    "transaction_hash": "0x45e71141170d61444e58e8622e91e4d35ce535d1e0bebd798bfc978b5dd886c6"
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
curl -H "X-Vault-Token: $(vault print token)" -H "X-Vault-Request: true" https://localhost:9200/v1/vault-ethereum/accounts/bob/balance
```

### Sample Response


```json
{
  "request_id": "47aee751-f12f-0aed-4b0c-87974fdc7551",
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
  "request_id": "7a31fb0e-93af-0bfb-6932-c835f9c6f29a",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "address": "0x8b470Eb0Cd23121814805886864b80085C58889d",
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
  "request_id": "4079bbb0-3235-b25c-796b-3cc27d5125b1",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "amount": "20000000000000000",
    "from": "0x8b470Eb0Cd23121814805886864b80085C58889d",
    "gas_limit": "21000",
    "gas_price": "20000000000",
    "nonce": "0",
    "signed_transaction": "0xf9014a808504a817c8008252089490259301a101a380f7138b50b6828cfffd4cbf6087470de4df820000b8de6638366438323032623238343737333539343030383235323038393434353932643866386437623030316537326362323661373365346661313830366135316163373964383830646530623662336137363430303030383032636130353932346264653765663130616138386462396336366464346635666231366234366466663233313962393936386265393833313138623537626235303536326130303162323462333130313030303466313364396132366233323038343532353761366366633262663831396133643535653366633836323633633566303737321ba01a47a3f579adc7d51c5bd183dceb1dcd3cce54a9d2f72c5045f2e384ef6fbc41a03909e60d1390eb4b45069a53f277cd5c1bedc2ebf15a66c749cb4e8e3c24de09",
    "to": "0x90259301a101A380F7138B50b6828cfFfd4Cbf60",
    "transaction_hash": "0x20e7eebf1ea4b0967647493c364a497aff515e94d8b630e9b59a0d6d8ef32380"
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
  "request_id": "9d43285d-85f8-44b5-1866-a694be9cb68f",
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
### Sample Response


```json
{
  "request_id": "203fec8b-7eb6-e2f5-6b54-83a5785cb067",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "contract": "0x2bEeBc3DedbD94A0521eEe7C8D2ae4214EB15b44",
    "from": "0x90259301a101A380F7138B50b6828cfFfd4Cbf60",
    "gas_limit": "1380521",
    "gas_price": "20000000000",
    "nonce": "1",
    "signed_transaction": "0xf918ab018504a817c800831510a98080b9185860806040523480156200001157600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040518060400160405280600581526020017f4649584544000000000000000000000000000000000000000000000000000000815250600290805190602001906200009f92919062000221565b506040518060400160405280601a81526020017f4578616d706c6520466978656420537570706c7920546f6b656e00000000000081525060039080519060200190620000ed92919062000221565b506012600460006101000a81548160ff021916908360ff160217905550600460009054906101000a900460ff1660ff16600a0a620f424002600581905550600554600660008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6005546040518082815260200191505060405180910390a3620002d0565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200026457805160ff191683800117855562000295565b8280016001018555821562000295579182015b828111156200029457825182559160200191906001019062000277565b5b509050620002a49190620002a8565b5090565b620002cd91905b80821115620002c9576000816000905550600101620002af565b5090565b90565b61157880620002e06000396000f3fe6080604052600436106100e85760003560e01c80638da5cb5b1161008a578063d4ee1d9011610059578063d4ee1d90146105bf578063dc39d06d14610616578063dd62ed3e14610689578063f2fde38b1461070e576100e8565b80638da5cb5b1461035b57806395d89b41146103b2578063a9059cbb14610442578063cae9ca51146104b5576100e8565b806323b872dd116100c657806323b872dd1461021b578063313ce567146102ae57806370a08231146102df57806379ba509714610344576100e8565b806306fdde03146100ed578063095ea7b31461017d57806318160ddd146101f0575b600080fd5b3480156100f957600080fd5b5061010261075f565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610142578082015181840152602081019050610127565b50505050905090810190601f16801561016f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561018957600080fd5b506101d6600480360360408110156101a057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506107fd565b604051808215151515815260200191505060405180910390f35b3480156101fc57600080fd5b506102056108ef565b6040518082815260200191505060405180910390f35b34801561022757600080fd5b506102946004803603606081101561023e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061094a565b604051808215151515815260200191505060405180910390f35b3480156102ba57600080fd5b506102c3610bf5565b604051808260ff1660ff16815260200191505060405180910390f35b3480156102eb57600080fd5b5061032e6004803603602081101561030257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610c08565b6040518082815260200191505060405180910390f35b34801561035057600080fd5b50610359610c51565b005b34801561036757600080fd5b50610370610dee565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103be57600080fd5b506103c7610e13565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156104075780820151818401526020810190506103ec565b50505050905090810190601f1680156104345780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561044e57600080fd5b5061049b6004803603604081101561046557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610eb1565b604051808215151515815260200191505060405180910390f35b3480156104c157600080fd5b506105a5600480360360608110156104d857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019064010000000081111561051f57600080fd5b82018360208201111561053157600080fd5b8035906020019184600183028401116401000000008311171561055357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929050505061104c565b604051808215151515815260200191505060405180910390f35b3480156105cb57600080fd5b506105d461127f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561062257600080fd5b5061066f6004803603604081101561063957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506112a5565b604051808215151515815260200191505060405180910390f35b34801561069557600080fd5b506106f8600480360360408110156106ac57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506113eb565b6040518082815260200191505060405180910390f35b34801561071a57600080fd5b5061075d6004803603602081101561073157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611472565b005b60038054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107f55780601f106107ca576101008083540402835291602001916107f5565b820191906000526020600020905b8154815290600101906020018083116107d857829003601f168201915b505050505081565b600081600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a36001905092915050565b6000610945600660008073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205460055461150f90919063ffffffff16565b905090565b600061099e82600660008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461150f90919063ffffffff16565b600660008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610a7082600760008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461150f90919063ffffffff16565b600760008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610b4282600660008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461152990919063ffffffff16565b600660008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a3600190509392505050565b600460009054906101000a900460ff1681565b6000600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610cab57600080fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff166000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60028054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610ea95780601f10610e7e57610100808354040283529160200191610ea9565b820191906000526020600020905b815481529060010190602001808311610e8c57829003601f168201915b505050505081565b6000610f0582600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461150f90919063ffffffff16565b600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610f9a82600660008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461152990919063ffffffff16565b600660008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a36001905092915050565b600082600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925856040518082815260200191505060405180910390a38373ffffffffffffffffffffffffffffffffffffffff16638f4ffcb1338530866040518563ffffffff1660e01b8152600401808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018481526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561120d5780820151818401526020810190506111f2565b50505050905090810190601f16801561123a5780820380516001836020036101000a031916815260200191505b5095505050505050600060405180830381600087803b15801561125c57600080fd5b505af1158015611270573d6000803e3d6000fd5b50505050600190509392505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461130057600080fd5b8273ffffffffffffffffffffffffffffffffffffffff1663a9059cbb6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff16846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156113a857600080fd5b505af11580156113bc573d6000803e3d6000fd5b505050506040513d60208110156113d257600080fd5b8101908080519060200190929190505050905092915050565b6000600760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146114cb57600080fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008282111561151e57600080fd5b818303905092915050565b600081830190508281101561153d57600080fd5b9291505056fea265627a7a723058200926f39699dae5a1f57cb3fbcf7a56a3b06e907ae6e0e6b2d897f543a7dc32c864736f6c634300050a00321ba08498aa9879821e1bb08f2df57596c0ad5b3ca5d8f90ea53110db4a95592b3dcba02c55de6e48381ff2b9e8aeee11975aa08378b724c81ca38658c86f8d81a635ba",
    "transaction_hash": "0x32c2a45a2b412abb5089b6a1239b4e79056d48fb910aeecd2f1c1445a3983776"
  },
  "warnings": null
}
```

## DEPLOY CONTRACT Owned


### Using the Vault CLI

```
vault write -format=json vault-ethereum/accounts/bob/deploy abi=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/Owned.abi bin=@/Users/immutability/projects/go/src/github.com/immutability-io/vault-ethereum/docker/erc20/build/Owned.bin

```
### Sample Response


```json
{
  "request_id": "fe39f5f5-5ebd-e1a3-7332-8463c360b97c",
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
  "request_id": "96a51aaa-8d29-d722-bbcd-8319760c608c",
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
  "request_id": "34d063c6-e435-926d-06ef-2420e80679c8",
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

