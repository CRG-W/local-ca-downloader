#!/bin/bash

# Generate the CA RSA
openssl genrsa \
-aes256 \
-passout env:CERT_PASSPHRASE \
-out /certs/ca-key.pem \
4096

# Generate the Public CA
openssl req \
-new \
-x509 \
-sha256 \
-days "${CA_TTL_DAYS}" \
-passin env:CERT_PASSPHRASE \
-key /certs/ca-key.pem \
-out /certs/public-ca.pem \
-subj "${CA_SUBJECT}"

# Generate the Cert RSA
openssl genrsa \
-out /certs/cert-key.pem \
4096

# Generate the CSR
openssl req \
-new \
-sha256 \
-subj "${CERT_CN}" \
-key /certs/cert-key.pem \
-out /certs/cert.csr

# Set the alt names
echo "subjectAltName=${CERT_ALT_NAMES}" >> extfile.cnf

# Generate a Cert from the CA
openssl x509 \
-req \
-sha256 \
-days "${CERT_TTL_DAYS}" \
-passin env:CERT_PASSPHRASE \
-in /certs/cert.csr \
-CA /certs/public-ca.pem \
-CAkey /certs/ca-key.pem \
-out /certs/cert.pem \
-extfile extfile.cnf \
-CAcreateserial

# Run the web app
/app
