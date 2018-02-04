#!/usr/bin/env bats

@test "disable ethereum secrets plugin" {
  run vault secrets disable ethereum
    [ "$status" -eq 0 ]
}

@test "delete ethereum secrets plugin from catalog" {
  run vault delete sys/plugins/catalog/ethereum-plugin
    [ "$status" -eq 0 ]
}

@test "write ethereum secrets plugin to catalog" {
  sha256=$(shasum -a 256 "$HOME/etc/vault.d/vault_plugins/vault-ethereum" | cut -d' ' -f1)
  vault_command="vault-ethereum --ca-cert=$HOME/etc/vault.d/root.crt --client-cert=$HOME/etc/vault.d/vault.crt --client-key=$HOME/etc/vault.d/vault.key"
  run vault write sys/plugins/catalog/ethereum-plugin sha_256="$sha256" command="$vault_command"
    [ "$status" -eq 0 ]
  results=$(vault read -format=json sys/plugins/catalog/ethereum-plugin | jq .data)
  plugin_name="$(echo $results | jq .name | tr -d '"')"
  command_name="$(echo $results | jq .command | tr -d '"')"
  sha256_from_catalog="$(echo $results | jq .sha256 | tr -d '"')"
    [ "$plugin_name" = "ethereum-plugin" ]
    [ "$command_name" = "vault-ethereum" ]
    [ "$sha256_from_catalog" = "$sha256" ]
}

@test "enable ethereum secrets plugin" {
  run vault secrets enable -path=ethereum -plugin-name=ethereum-plugin plugin
    [ "$status" -eq 0 ]
}
