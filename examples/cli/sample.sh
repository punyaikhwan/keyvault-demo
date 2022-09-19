#!/bin/bash

# create key vault
az keyvault create --location southeastasia --resource-group demo --name civitas

# create key
az keyvault key create --vault-name civitas --name student

# encrypt with key
$plaintext = cat plaintext.txt
az keyvault key encrypt --name student --vault-name civitas --algorithm RSA-OAEP --value $plaintext --data-type plaintext

# decrypt with key
az keyvault key decrypt --name student --vault-name civitas --algorithm RSA-OAEP --data-type base64 --version ... --value ...

# rotate automatically
az keyvault key rotation-policy update --name student --vault-name civitas --value rotation-policy.json

# rotate manually
az keyvault key rotate --name student --vault-name demo-kms
