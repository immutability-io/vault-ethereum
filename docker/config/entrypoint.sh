#!/bin/bash

CONFIG_DIR="/home/vault/config"
INIT_SCRIPT="/home/vault/config/init.sh"
CA_CERT="$CONFIG_DIR/root.crt"
CA_KEY="$CONFIG_DIR/root.key"
TLS_KEY="$CONFIG_DIR/vault.key"
TLS_CERT="$CONFIG_DIR/vault.crt"
OPENSSL_CONFIG="$CONFIG_DIR/vault.cnf"
CSR="$CONFIG_DIR/vault.csr"

export VAULT_ADDR="https://127.0.0.1:9200"
export VAULT_CACERT="$CA_CERT"

function create_config {

	cat > "$OPENSSL_CONFIG" << EOF

[req]
default_bits = 2048
encrypt_key  = no
default_md   = sha256
prompt       = no
utf8         = yes

# Speify the DN here so we aren't prompted (along with prompt = no above).
distinguished_name = req_distinguished_name

# Extensions for SAN IP and SAN DNS
req_extensions = v3_req

# Be sure to update the subject to match your organization.
[req_distinguished_name]
C  = US
ST = Maryland
L  = Immutability
O  = Immutability LLC
CN = localhost

# Allow client and server auth. You may want to only allow server auth.
# Link to SAN names.
[v3_req]
basicConstraints     = CA:FALSE
subjectKeyIdentifier = hash
keyUsage             = digitalSignature, keyEncipherment
extendedKeyUsage     = clientAuth, serverAuth
subjectAltName       = @alt_names

# Alternative names are specified as IP.# and DNS.# for IPs and
# DNS accordingly.
[alt_names]
IP.1  = 127.0.0.1
DNS.1 = localhost
EOF
}

function gencerts {

    create_config
	openssl req \
	-new \
	-sha256 \
	-newkey rsa:2048 \
	-days 120 \
	-nodes \
	-x509 \
	-subj "/C=US/ST=Maryland/L=Immutability/O=Immutability LLC" \
	-keyout "$CA_KEY" \
	-out "$CA_CERT"

	openssl genrsa -out "$TLS_KEY" 2048

	openssl req \
	-new -key "$TLS_KEY" \
	-out "$CSR" \
	-config "$OPENSSL_CONFIG"

	openssl x509 \
	-req \
	-days 120 \
	-in "$CSR" \
	-CA "$CA_CERT" \
	-CAkey "$CA_KEY" \
	-CAcreateserial \
	-sha256 \
	-extensions v3_req \
	-extfile "$OPENSSL_CONFIG" \
	-out "$TLS_CERT"

	openssl x509 -in "$TLS_CERT" -noout -text
    chown -R nobody:nobody $CONFIG_DIR && chmod -R 777 $CONFIG_DIR
}

mkdir -p $CONFIG_DIR
gencerts

nohup vault server -log-level=debug -config /home/vault/config/vault.hcl &
VAULT_PID=$!

which bash

if [ -f "$INIT_SCRIPT" ]; then
    /bin/bash $INIT_SCRIPT
fi

wait $VAULT_PID
