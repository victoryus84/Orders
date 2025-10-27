#!/bin/sh
set -e

echo "$(date): Starting certificate renewal"

certbot renew --webroot -w /var/www/certbot

if [ $? -eq 0 ]; then
    echo "$(date): Certificates renewed successfully, restarting nginx"
    docker container restart nginx
    echo "$(date): Nginx restarted"
else
    echo "$(date): Certificate renewal failed"
    exit 1
fi