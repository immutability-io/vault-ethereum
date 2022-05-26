default_lease_ttl = "24h"
disable_mlock     = "true"
max_lease_ttl     = "43800h"

#backend "raft" {
#  path    = "/var/raft/data"
#  node_id = "node01"
#}
// from https://www.vaultproject.io/docs/configuration/storage/postgresql
storage "postgresql" {
  connection_url = "myurl_connection"
  table = "vault_kv_store"
  ha_table = "vault_ha_locks"
}
api_addr     = "https://localhost:9200"
cluster_name = "vault"
cluster_addr = "https://127.0.0.1:9201"
ui           = "true"

plugin_directory = "/home/vault/plugins"
listener "tcp" {
  address            = "0.0.0.0:9200"
  tls_cert_file      = "/home/vault/config/vault.crt"
  tls_client_ca_file = "/home/vault/config/root.crt"
  tls_key_file       = "/home/vault/config/vault.key"
}