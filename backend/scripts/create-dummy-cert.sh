: "${DOMAIN:?Need to set DOMAIN env var}"

CERT_DIR="/etc/letsencrypt/live/${DOMAIN}"

if [ -f "${CERT_DIR}/fullchain.pem" ] && [ -f "${CERT_DIR}/privkey.pem" ]; then
  echo "$(date): Real certs exist for ${DOMAIN}"
  exit 0
fi

echo "$(date): Creating dummy cert for ${DOMAIN}"
mkdir -p "${CERT_DIR}"
chmod 700 "${CERT_DIR}"

openssl req -x509 -nodes -newkey rsa:2048 -days 1 \
  -keyout "${CERT_DIR}/privkey.pem" \
  -out "${CERT_DIR}/fullchain.pem" \
  -subj "/CN=${DOMAIN}"

chmod 600 "${CERT_DIR}/privkey.pem" "${CERT_DIR}/fullchain.pem"
echo "$(date): Dummy cert created at ${CERT_DIR}"