#!/bin/sh

set -eu

if [ ! -z "${1:-}" ]; then
	exec "$@"
else
  # check if config.js is present (via volume)
  # or if it's missing
  if [ ! -f "./config.js" ]; then
    # config.js missing, generate it
    if [ ! -z "${CONFIGJS:-}" ]; then
      # use $CONFIGJS variable that might be passed to the container:
      # --env CONFIGJS="$(cat public/config.example.js)"
      echo "${CONFIGJS}" > ./config.js
    else
      # Try to guess where the API is located by using DOMAIN or VIRTUAL_HOST and prefix it with "api."
      API_HOST=${API_HOST:-"api.${VIRTUAL_HOST:-"${DOMAIN:-"local.cortezaproject.org"}"}"}
      API_BASEURL=${API_FULL_URL:-"//${API_HOST}"}

      echo "window.CortezaAPI = '${API_BASEURL}'" > ./config.js
    fi
  fi

  BASE_PATH=${BASE_PATH:-"/"}
  if [ $BASE_PATH != "/" ]; then
    BASE_PATH_LENGTH=${#BASE_PATH}
    BASE_PATH_LAST_CHAR=${BASE_PATH:BASE_PATH_LENGTH-1:1}

    if [ $BASE_PATH_LAST_CHAR != "/" ]; then
      BASE_PATH="$BASE_PATH/"
    fi
  fi

  sed -i "s|<base href=/ >|<base href=\"$BASE_PATH\">|g" ./index.html
  sed -i "s|<base href=\"/\">|<base href=\"$BASE_PATH\">|g" ./index.html

  sed -i "s|{{BASE_PATH}}|$BASE_PATH|g" /etc/nginx/nginx.conf


  nginx -g "daemon off;"
fi
