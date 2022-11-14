#!/bin/sh

set -eu

if [ ! -z "${1:-}" ]; then
	exec "$@"
else

  function config() {
    prefix=$1

    # check if config.js is present (via volume)
    # or if it's missing
    if [ ! -f "$prefix/config.js" ]; then
      # config.js missing, generate it
      if [ ! -z "${CONFIGJS:-}" ]; then
        # use $CONFIGJS variable that might be passed to the container:
        # --env CONFIGJS="$(cat public/config.example.js)"
        echo "${CONFIGJS}" > "$prefix/config.js"
      else
        # Try to guess where the API is located by using DOMAIN or VIRTUAL_HOST and prefix it with "api."
        API_HOST=${API_HOST:-"api.${VIRTUAL_HOST:-"${DOMAIN:-"local.cortezaproject.org"}"}"}
        API_BASEURL=${API_FULL_URL:-"//${API_HOST}"}
        DISCOVERY_BASEURL=${DISCOVERY_BASE_URL:-"//discovery.local.cortezaproject.org"}

        echo "window.CortezaAPI = '${API_BASEURL}'" > "$prefix/config.js"
        echo "window.CortezaDiscoveryAPI = '${DISCOVERY_BASEURL}'" >> "$prefix/config.js"
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

    sed -i "s|<base href=/ >|<base href=\"${BASE_PATH}$(echo "$prefix/" | cut -c 3-)\">|g" "$prefix/index.html"
    sed -i "s|<base href=\"/\">|<base href=\"${BASE_PATH}$(echo "$prefix/" | cut -c 3-)\">|g" "$prefix/index.html"
    sed -i "s|{{BASE_PATH}}|$BASE_PATH|g" /etc/nginx/nginx.conf
  }

  config './admin'
  config './reporter'
  config './compose'
  config './workflow'
  config './discovery'
  config './privacy'
  config '.'

  nginx -g "daemon off;"
fi
