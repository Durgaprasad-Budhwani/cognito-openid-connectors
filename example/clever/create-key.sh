#!/bin/bash -eu

KEY_FILE=./resources/jwtRS256.key

if [ ! -f "$KEY_FILE" ]; then
  echo "  --- Creating private key, as it does not exist ---"
  ssh-keygen -t rsa -b 4096 -m PEM -f "$KEY_FILE" -N ''
  openssl rsa -in "$KEY_FILE" -pubout -outform PEM -out "$KEY_FILE".pub
  chmod 777 "$KEY_FILE".pub
  chmod 777 "$KEY_FILE"
fi


