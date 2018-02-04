#!/usr/bin/env bats

@test "test read $FUNDED_ACCOUNT balance" {
  pending_balance="$(vault read -format=json ethereum/accounts/$FUNDED_ACCOUNT/balance | jq .data.pending_balance | tr -d '"')"
    [ "$pending_balance" != "0" ]
}

@test "test send ETH from $FUNDED_ACCOUNT" {
  [ "$FUNDED_ACCOUNT" != "" ]
  balance_result="$(vault read -format=json ethereum/accounts/$FUNDED_ACCOUNT/balance | jq .data)"
  pending_balance="$(echo $balance_result | jq .pending_balance | tr -d '"')"
  sender_address="$(echo $balance_result | jq .address | tr -d '"')"
  send_amount="${pending_balance:0:${#pending_balance} - 3}"
  recipient_address="$(vault write -format=json ethereum/accounts/recipient chain_id=1977 | jq .data.address | tr -d '"')"
  result="$(vault write -format=json ethereum/accounts/$FUNDED_ACCOUNT/debit to=$recipient_address amount=$send_amount | jq .data)"
  from_address="$(echo $result | jq .from_address | tr -d '"')"
  to_address="$(echo $result | jq .to_address | tr -d '"')"
    [ "$from_address" = "$sender_address" ]
    [ "$to_address" = "$recipient_address" ]
}

@test "test deploy contract from $FUNDED_ACCOUNT" {
  [ "$FUNDED_ACCOUNT" != "" ]
  result="$(vault write ethereum/accounts/$FUNDED_ACCOUNT/contracts/helloworld @send_contract.json | jq .data)"
  tx_hash="$(echo $result | jq .tx_hash | tr -d '"')"
  read_result="$(vault read ethereum/accounts/$FUNDED_ACCOUNT/contracts/helloworld | jq .data)"
  tx_hash_read="$(echo $read_result | jq .tx_hash | tr -d '"')"
    [ "$tx_hash" = "$tx_hash_read" ]
}
