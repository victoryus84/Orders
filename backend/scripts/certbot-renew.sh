#!/bin/sh
set -e

WEBROOT=${WEBROOT:-/var/www/certbot}

echo "$(date): certbot renew --webroot -w ${WEBROOT}"
certbot renew --webroot -w "${WEBROOT}" --quiet

if [ $? -eq 0 ]; then
  echo "$(date): Renew succeeded."

  # Попробуем аккуратно перезагрузить nginx:
  # 1) Если доступен docker CLI (и /var/run/docker.sock смонтирован), перезапустим контейнер nginx
  if command -v docker >/dev/null 2>&1; then
    echo "$(date): Restarting nginx via docker CLI"
    docker restart nginx || true
  else
    # Альтернатива: создать файл-сигнал или вызвать reload через docker exec (если есть права)
    echo "$(date): docker CLI not available - ensure nginx is reloaded manually or mount docker.sock"
  fi
else
  echo "$(date): Renew failed"
fi
