#!/usr/bin/env bats

@test "test write convert" {
  conversion="$(vault write -format=json ethereum/convert @convert.json| jq .data)"
  amount_to="$(echo $conversion | jq -r .amount_to)"
  valid="0.0000000000000000000000000000044"
    [ "$amount_to" = "$valid" ]
}
