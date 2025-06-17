#!/bin/sh

hashed_password=$(openssl passwd -apr1 "$BASIC_PASSWORD")

echo "$BASIC_USERNAME:$hashed_password" > /etc/nginx/.htpasswd 2>/dev/null
