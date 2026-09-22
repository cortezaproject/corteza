import path from 'path'
import logger from './logger'
import * as config from './config'
import * as gRPCServer from './grpc-server'
import ClientScriptsHandler from './grpc-handlers/client-scripts'
import ServerScriptsHandler from './grpc-handlers/server-scripts'
import ClientScriptsService from './services/client-scripts'
import ServerScriptsService from './services/server-scripts'
import DependenciesService from './services/dependencies'
import * as Sentry from '@sentry/node'
import EnvCheck from './check'
import ScriptLoader from './loader'
import { Watcher } from './types'

Sentry.init(config.sentry)

EnvCheck()

const { searchPaths } = config.extensions

// <search-path>/client-scripts/<bundle>/<path-to-script>/*.js
const clientScriptsLoader = new ScriptLoader(
  searchPaths,
  'client-scripts',
  path.join('*', '**', '*.js'),
)
// <search-path>/server-scripts/<path-to-script>/*.js
const serverScriptsLoader = new ScriptLoader(searchPaths, 'server-scripts', path.join('**', '*.js'))

let clientScriptsService: ClientScriptsService
let serverScriptsService: ServerScriptsService
let dependenciesService: DependenciesService

const watchers: Array<Watcher> = []

if (config.extensions.clientScripts.enabled) {
  clientScriptsService = new ClientScriptsService({
    logger,
    loader: clientScriptsLoader,
    config: { bundler: config.bundler },
  })

  if (config.extensions.clientScripts.watch) {
    watchers.push(clientScriptsService)
  }
}

if (config.extensions.serverScripts.enabled) {
  const { cortezaServers: cServers, frontend } = config.execContext

  serverScriptsService = new ServerScriptsService({
    logger,
    loader: serverScriptsLoader,
    config: { cServers, frontend },
  })

  if (config.extensions.serverScripts.watch) {
    watchers.push(serverScriptsService)
  }
}

if (config.extensions.dependencies.autoUpdate) {
  dependenciesService = new DependenciesService({
    logger,
    searchPaths,
    installer: config.extensions.dependencies.installer,
  })

  dependenciesService.watch(() => {
    // @todo can we be more selective about what should be reloaded
    try {
      serverScriptsService.process()
      clientScriptsService.process()
    } catch (e) {
      Sentry.captureException(e)
    }
  })
}

gRPCServer
  .LoadDefinitions(path.join(config.protobuf.path, '/service-corredor.proto'))
  .then(def => gRPCServer.VerifyDefinitions(def))
  .then(corredor => {
    const serviceDefinitions: gRPCServer.ServiceDefinition = new Map()
    serviceDefinitions.set(
      // @ts-expect-error service is defined
      corredor.ServerScripts.service,
      ServerScriptsHandler(serverScriptsService, logger.child({ name: 'gRPC.ServerScripts' })),
    )

    serviceDefinitions.set(
      // @ts-expect-error service is defined
      corredor.ClientScripts.service,
      ClientScriptsHandler(clientScriptsService, logger.child({ name: 'gRPC.ClientScripts' })),
    )

    gRPCServer.Start(config.server, logger, serviceDefinitions)
  })
  .catch(e => {
    Sentry.captureException(e)
    logger.warn('could not start gRPC server:', e.message)
  })

try {
  Promise.all([
    // @ts-expect-error service is defined
    serverScriptsService.process(),
    // @ts-expect-error service is defined
    clientScriptsService.process(),
  ])

  watchers.forEach(w => w.watch())
} catch (e) {
  Sentry.captureException(e)
}
