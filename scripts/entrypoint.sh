#!/bin/bash

export CERT_PASSPHRASE="defaultPassphase" \
  CERT_CN="/CN=HomeLab" \
  CERT_ALT_NAMES="DNS:localhost,IP:127.0.0.1" \
  CA_TTL_DAYS=18250 \
  CERT_TTL_DAYS=18250 \
  CA_SUBJECT="/C=US/ST=AZ/L=Phoenix/O=Home/CN=HomeLab"

generate_certs() {
  # Generate the CA RSA
  echo "Generating CA RSA Private Key..."
  openssl genrsa \
  -aes256 \
  -passout env:CERT_PASSPHRASE \
  -out ./certs/ca-key.pem \
  4096

  # Generate the Public CA
  echo "Generating Public CA..."
  openssl req \
  -new \
  -x509 \
  -sha256 \
  -days "${CA_TTL_DAYS}" \
  -passin env:CERT_PASSPHRASE \
  -key ./certs/ca-key.pem \
  -out ./certs/public-ca.pem \
  -subj "${CA_SUBJECT}"

  # Generate the Cert RSA
  echo "Generating Cert RSA Private Key..."
  openssl genrsa \
  -out ./certs/cert-key.pem \
  4096

  # Generate the CSR
  echo "Generating CSR..."
  openssl req \
  -new \
  -sha256 \
  -subj "${CERT_CN}" \
  -key ./certs/cert-key.pem \
  -out ./certs/cert.csr

  # Set the alt names
  echo "Setting alt names..."
  echo "subjectAltName=${CERT_ALT_NAMES}" >> ./certs/extfile.cnf

  # Generate a Cert from the CA
  echo "Generating Public Cert from CA..."
  openssl x509 \
  -req \
  -sha256 \
  -days "${CERT_TTL_DAYS}" \
  -passin env:CERT_PASSPHRASE \
  -in ./certs/cert.csr \
  -CA ./certs/public-ca.pem \
  -CAkey ./certs/ca-key.pem \
  -out ./certs/cert.pem \
  -extfile ./certs/extfile.cnf \
  -CAcreateserial
}


if [ ! -d "certs" ]; then
  echo "Creating certs directory..."
  mkdir "certs"
else
  echo "Certs directory already exists."
fi

if [ -z "$(ls -A certs)" ]; then
  echo "Certs directory is empty, generating certs..."
  generate_certs
else
  echo "Certs directory is not empty, skipping cert generation."
fi


if [ -x "/app" ]; then
  echo "Attempting to run app executable..."
  /app
else
  echo "Attempting to run app via Go file..."
  go run ./cmd/server/main.go
fi
