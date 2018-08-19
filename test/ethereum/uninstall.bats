#!/usr/bin/env bats

@test "disable ethereum secrets plugin" {
  run vault secrets disable ethereum
    [ "$status" -eq 0 ]
}

@test "delete ethereum secrets plugin from catalog" {
  run vault delete sys/plugins/catalog/ethereum-plugin
    [ "$status" -eq 0 ]
}
