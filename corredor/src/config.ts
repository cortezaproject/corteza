import path from 'path'
import { NodeOptions as SentryOptions } from '@sentry/node'

// Read .env into process.ENV.*
require('dotenv').config({ quiet: true })

// Let's make our code a bit shorter
const e = process.env

const rootDir = path.normalize(path.join(__dirname, '../'))

function isTrue(input: string | undefined): boolean | undefined {
  if (undefined === input) {
    return
  }

  return /^s*(t(rue)?|y(es)?|1)s*$/i.test(input)
}

/**
 * Discovers/generates API baseURL from environmental variables & given service
 *
 * @param {string} service
 * @returns {string|undefined}
 */
function discoverApiBaseURL(service: string): string | undefined {
  const host =
    e.CORREDOR_EXEC_CSERVERS_API_HOST ??
    // DOMAIN will be present in the standard configuration
    e.DOMAIN ??
    // Other ways to get to the hostname
    e.HOSTNAME ??
    e.HOST

  if (host === undefined) {
    return undefined
  }

  const tpl =
    // will replace HOST and SERVICE
    e.CORREDOR_EXEC_CSERVERS_API_BASEURL_TEMPLATE ??
    // Assume standard "api." prefix for the API
    'https://api.{host}/{service}'

  return tpl.replace('{host}', host).replace('{service}', service)
}

/**
 * Discovers/generates frontend baseURL from environmental variables & given service
 *
 * @param {string} service
 * @returns {string|undefined}
 */
function discoverFrontendBaseURL(): string {
  return (
    // Try to get this from the backend provision variable
    e.PROVISION_SETTINGS_AUTH_FRONTEND_URL_BASE ??
    // Fallback
    `https://${e.DOMAIN || e.HOSTNAME || e.HOST || 'localhost'}`
      // If we've recycled hostname for the API somehow, let's cut it off
      .replace('api.', '')
  )
}

export const env = (e.CORREDOR_ENVIRONMENT ?? e.CORREDOR_ENV ?? e.NODE_ENV ?? 'prod')
  .trim()
  .toLowerCase()

export const isDevelopment = env.startsWith('dev')
export const isProduction = !isDevelopment

const certPath = e.CORREDOR_SERVER_CERTIFICATES_PATH ?? '/certs'

// Server settings
// CORREDOR_ADDR is used by the API as well to configure gRPC client connection
export const server = {
  addr: e.CORREDOR_ADDR || 'localhost:50051',

  certificates: {
    enabled: isTrue(e.CORREDOR_SERVER_CERTIFICATES_ENABLED) ?? isProduction,
    ca: e.CORREDOR_SERVER_CERTIFICATES_CA ?? path.join(certPath, 'ca.crt'),
    private: e.CORREDOR_SERVER_CERTIFICATES_PRIVATE ?? path.join(certPath, 'private.key'),
    public: e.CORREDOR_SERVER_CERTIFICATES_PUBLIC ?? path.join(certPath, 'public.crt'),
  },
}

export const logger = {
  // Enable/disable logging
  // CORREDOR_LOG_ENABLED is used by the API as well to configure gRPC client logging
  enabled: isTrue(e.CORREDOR_LOG_ENABLED) ?? true,

  // Enable/disable pretty logging
  //
  // if LOG_PRETTY is set or inherit from debug
  prettyPrint: isTrue(e.CORREDOR_LOG_PRETTY) ?? isDevelopment,

  // Log level
  level: e.CORREDOR_LOG_LEVEL ?? (isDevelopment ? 'trace' : 'info'),
}

// See https://docs.sentry.io/clients/node/config/ for details
export const sentry: SentryOptions = {
  enabled: !!e.SENTRY_DSN,
  dsn: e.SENTRY_DSN,
  serverName: e.SENTRY_SERVERNAME,
  release: e.SENTRY_RELEASE,
  dist: e.SENTRY_DIST,
  environment: e.SENTRY_ENVIRONMENT,
}

// Root of the protobuf definitions; CORREDOR_CORTEZA_PROTOBUF_PATH is the legacy alias
export const protobuf = {
  path: path.normalize(
    e.CORREDOR_PROTOBUF_PATH ??
      e.CORREDOR_CORTEZA_PROTOBUF_PATH ??
      path.join(rootDir, '../def/protobuf'),
  ),
}

export const services = {
  scriptRunner: {
    timeout: {
      max: 30 * 1000, // 30s
      min: 100, // 0.1s
      def: 2 * 1000, // 2s
    },
  },
}

export const execContext = {
  cortezaServers: {
    system: {
      apiBaseURL:
        e.CORREDOR_EXEC_CTX_CORTEZA_SERVERS_SYSTEM_API_BASEURL ?? discoverApiBaseURL('system'),
    },

    compose: {
      apiBaseURL:
        e.CORREDOR_EXEC_CTX_CORTEZA_SERVERS_COMPOSE_API_BASEURL ?? discoverApiBaseURL('compose'),
    },
  },

  frontend: {
    baseURL: e.CORREDOR_EXEC_CTX_FRONTEND_BASEURL ?? discoverFrontendBaseURL(),
  },
}

export const bundler = {
  outputPath: path.resolve(e.CORREDOR_BUNDLER_OUTPUT_PATH ?? '/tmp/corredor/bundler-dist'),
  enabled: isTrue(e.CORREDOR_BUNDLER_ENABLED) ?? true,
}

const extensionsSearchPaths = (
  e.CORREDOR_EXT_SEARCH_PATHS ?? [path.join('usr'), path.join('usr', '*')].join(':')
)
  .trim()
  .split(/[:]+/)
  .filter(p => p.length > 0)
  .map(p => {
    if (!p.startsWith('/')) {
      return path.join(rootDir, p)
    }

    return p
  })

export const extensions = {
  // Path to extensions
  searchPaths: extensionsSearchPaths,

  dependencies: {
    // do we automatically update deps?
    autoUpdate: isTrue(e.CORREDOR_EXT_DEPENDENCIES_AUTO_UPDATE) ?? true,

    // package manager used to install extension dependencies;
    // undefined means "yarn when available, npm otherwise"
    installer: e.CORREDOR_EXT_DEPENDENCIES_INSTALLER?.trim().toLowerCase(),
  },

  serverScripts: {
    enabled: isTrue(e.CORREDOR_EXT_SERVER_SCRIPTS_ENABLED) ?? true,
    watch: isTrue(e.CORREDOR_EXT_SERVER_SCRIPTS_WATCH) ?? true,
  },

  clientScripts: {
    enabled: isTrue(e.CORREDOR_EXT_CLIENT_SCRIPTS_ENABLED) ?? true,
    watch: isTrue(e.CORREDOR_EXT_CLIENT_SCRIPTS_WATCH) ?? true,
  },

  // vueComponents: {
  //   enabled: isTrue(e.CORREDOR_EXT_VUE_COMPONENTS_ENABLED) ?? true,
  //   watch: isTrue(e.CORREDOR_EXT_VUE_COMPONENTS_WATCH) ?? true,
  // },
  //
  // styleSheets: {
  //   enabled: isTrue(e.CORREDOR_EXT_STYLESHEETS_ENABLED) ?? true,
  //   watch: isTrue(e.CORREDOR_EXT_STYLESHEETS_WATCH) ?? true,
  // },
}
