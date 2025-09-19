#!/usr/bin/env bash

# Set paths
SECURITY_DIR="./security"
PRIVATE_KEY_FILE="$SECURITY_DIR/private.pem"
PUBLIC_KEY_FILE="$SECURITY_DIR/public.pem"

# Make sure the security directory exists
mkdir -p "$SECURITY_DIR"

# Generate a 2048-bit RSA private key
openssl genpkey -algorithm RSA -out "$PRIVATE_KEY_FILE" -pkeyopt rsa_keygen_bits:2048

# Extract the public key from the private key
openssl rsa -pubout -in "$PRIVATE_KEY_FILE" -out "$PUBLIC_KEY_FILE"

echo "Private key saved to $PRIVATE_KEY_FILE"
echo "Public key saved to $PUBLIC_KEY_FILE"
