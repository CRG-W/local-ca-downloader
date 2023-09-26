#!/bin/bash

CERT_PASSPHRASE=$1
CA_TTL_DAYS=$2
CA_SUBJECT=$3
CERT_CN="/CN=$4"
CERT_ALT_NAMES=$5
CERT_TTL_DAYS=$6

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

mkdir certs-backup
cp -r certs/* certs-backup/

if generate_certs; then
    echo "Certs generated successfully."
    rm -rf certs-backup
else
    echo "Certs generation failed."
    cp -r certs-backup/* certs/
    rm -rf certs-backup
    exit 1
fi
