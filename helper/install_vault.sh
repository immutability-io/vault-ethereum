#!/bin/bash

PLUGIN_VERSION="0.1.0"
VAULT_VERSION="0.11.0"

function print_help {
    echo "Usage: bash install_vault.sh OPTIONS"
    echo -e "\nOPTIONS:"
    echo -e "  --linux\tInstall Linux version"
    echo -e "  --darwin\tInstall Darwin (MacOS) version"
}

function gencerts {
  TMPDIR="$(pwd)/tls"

  # Optional: Ensure the target directory exists and is empty.
  rm -rf "${TMPDIR}"
  mkdir -p "${TMPDIR}"

  cat > "${TMPDIR}/openssl.cnf" << EOF
[req]
default_bits = 2048
encrypt_key  = no # Change to encrypt the private key using des3 or similar
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
L  = Vault
O  = Immutability
CN = localhost

# Allow client and server auth. You may want to only allow server auth.
# Link to SAN names.
[v3_req]
basicConstraints     = CA:FALSE
subjectKeyIdentifier = hash
keyUsage             = digitalSignature, keyEncipherment
extendedKeyUsage     = clientAuth, serverAuth
subjectAltName       = @alt_names

# Alternative names are specified as IP.# and DNS.# for IP addresses and
# DNS accordingly. 
[alt_names]
IP.1  = 127.0.0.7
DNS.1 = localhost
EOF

  openssl req \
    -new \
    -newkey rsa:2048 \
    -days 120 \
    -nodes \
    -x509 \
    -subj "/C=US/ST=Maryland/L=Vault/O=My Company CA" \
    -keyout "${TMPDIR}/ca.key" \
    -out "${TMPDIR}/ca.crt"

  # Generate the private key for the service. Again, you may want to increase
  # the bits to 4096.
  openssl genrsa -out "${TMPDIR}/my-service.key" 2048

  # Generate a CSR using the configuration and the key just generated. We will
  # give this CSR to our CA to sign.
  openssl req \
    -new -key "${TMPDIR}/my-service.key" \
    -out "${TMPDIR}/my-service.csr" \
    -config "${TMPDIR}/openssl.cnf"
    
  # Sign the CSR with our CA. This will generate a new certificate that is signed
  # by our CA.
  openssl x509 \
    -req \
    -days 120 \
    -in "${TMPDIR}/my-service.csr" \
    -CA "${TMPDIR}/ca.crt" \
    -CAkey "${TMPDIR}/ca.key" \
    -CAcreateserial \
    -extensions v3_req \
    -extfile "${TMPDIR}/openssl.cnf" \
    -out "${TMPDIR}/my-service.crt"

  openssl x509 -in "${TMPDIR}/my-service.crt" -noout -text

  mv $TMPDIR/ca.crt $HOME/etc/vault.d/root.crt
  mv $TMPDIR/ca.key $HOME/etc/vault.d/root.key
  mv $TMPDIR/my-service.crt $HOME/etc/vault.d/vault.crt
  mv $TMPDIR/my-service.key $HOME/etc/vault.d/vault.key
  rm -rf "${TMPDIR}"
}


function grab_hashitool {
  echo "Tool: $1"
  echo "Version: $2"
  echo "OS: $3"


  wget  --progress=bar:force -O ./$1.zip https://releases.hashicorp.com/$1/$2/$1_$2_$3_amd64.zip
  wget  --progress=bar:force -O ./$1_$2_SHA256SUMS https://releases.hashicorp.com/$1/$2/$1_$2_SHA256SUMS
  wget  --progress=bar:force -O ./$1_$2_SHA256SUMS.sig https://releases.hashicorp.com/$1/$2/$1_$2_SHA256SUMS.sig
  keybase pgp verify -d ./$1_$2_SHA256SUMS.sig -i ./$1_$2_SHA256SUMS
  if [[ $? -eq 2 ]] ; then
    echo "Vault Validation Failed: Signature doesn't verify!"
    exit 2
  fi
  unzip ./$1.zip
  mv ./$1 /usr/local/bin/$1
  rm ./$1_$2_SHA256SUMS.sig
  rm ./$1_$2_SHA256SUMS
  rm ./$1.zip
}


function grab_plugin {
  echo "OS: $1"
  echo "Version: $2"

  wget --progress=bar:force -O ./$1.zip https://github.com/immutability-io/vault-ethereum/releases/download/v$2/vault-ethereum_$2_$1_amd64.zip
  wget --progress=bar:force -O ./SHA256SUMS https://github.com/immutability-io/vault-ethereum/releases/download/v$2/SHA256SUMS
  wget --progress=bar:force -O ./SHA256SUMS.sig https://github.com/immutability-io/vault-ethereum/releases/download/v$2/SHA256SUMS.sig
  keybase pgp verify -d ./SHA256SUMS.sig -i ./SHA256SUMS
  if [[ $? -eq 2 ]] ; then
    echo "Plugin Validation Failed: Signature doesn't verify!"
    exit 2
  fi
  rm ./SHA256SUMS.sig
  rm ./SHA256SUMS
}

function move_plugin {
  echo "OS: $1"
  unzip ./$1.zip
  rm ./$1.zip
  mv ./vault-ethereum $HOME/etc/vault.d/vault_plugins/vault-ethereum
}


if [ -n "`$SHELL -c 'echo $ZSH_VERSION'`" ]; then
    # assume Zsh
    shell_profile="zshrc"
elif [ -n "`$SHELL -c 'echo $BASH_VERSION'`" ]; then
    # assume Bash
    shell_profile="bashrc"
fi

if [ "$1" == "--darwin" ]; then
    PLUGIN_OS="darwin"
elif [ "$1" == "--linux" ]; then
    PLUGIN_OS="linux"
elif [ "$1" == "--help" ]; then
    print_help
    exit 0
else
    print_help
    exit 1
fi

if [ -d "$HOME/etc/vault.d" ]; then
    echo "The 'etc/vault.d' directories already exist. Exiting."
    exit 1
fi

mkdir -p $HOME/etc/vault.d/vault_plugins
mkdir -p $HOME/etc/vault.d/data

gencerts

grab_plugin $PLUGIN_OS $PLUGIN_VERSION
move_plugin $PLUGIN_OS
grab_hashitool vault $VAULT_VERSION $PLUGIN_OS

cat << EOF > $HOME/etc/vault.d/vault.hcl
"default_lease_ttl" = "24h"
"disable_mlock" = "true"
"max_lease_ttl" = "24h"

"backend" "file" {
  "path" = "$HOME/etc/vault.d/data"
}

"api_addr" = "https://localhost:8200"
"ui" = "true"
"listener" "tcp" {
  "address" = "localhost:8200"

  "tls_cert_file" = "$HOME/etc/vault.d/vault.crt"
  "tls_client_ca_file" = "$HOME/etc/vault.d/root.crt"
  "tls_key_file" = "$HOME/etc/vault.d/vault.key"
}

"plugin_directory" = "$HOME/etc/vault.d/vault_plugins"
EOF

touch "$HOME/.${shell_profile}"
{
    echo "# Vault"
    echo "export VAULT_ADDR=https://localhost:8200"
    echo "export VAULT_CACERT=$HOME/etc/vault.d/root.crt"
} >> "$HOME/.${shell_profile}"


echo -e "$HOME/.${shell_profile} has been modified."
echo "============================================="
echo "The following were set in your profile:"
echo "export VAULT_ADDR=https://localhost:8200"
echo "export VAULT_CACERT=$HOME/etc/vault.d/root.crt"
echo -e "=============================================\n"
