import * as fs from 'fs'
import * as grpc from '@grpc/grpc-js'
import { Logger } from 'pino'

interface ServerConfig {
  addr: string
  certificates: {
    enabled: boolean
    ca: string
    private: string
    public: string
  }
}

export type ServiceDefinition = Map<grpc.ServiceDefinition<unknown>, unknown>

/**
 * Initializes the server
 */
export function Start(
  { addr, certificates }: ServerConfig,
  logger: Logger,
  services: ServiceDefinition,
): void {
  const server = new grpc.Server({
    // setting this to 16mB
    // @todo should be configurable
    'grpc.max_receive_message_length': 2 << 23,
    'grpc.max_send_message_length': 2 << 23,
  })
  const log = logger.child({ name: 'gRPC' })

  log.debug('starting server')

  const handle = (): void => {
    // Override signal handler with more severe approach
    process.on('SIGINT', () => {
      log.warn('forcing server to stop')
      server.forceShutdown()
    })

    log.debug('trying to stop the server')
    server.tryShutdown(() => {
      log.info('server stopped')
    })
  }

  process.on('SIGINT', handle)
  process.on('SIGTERM', handle)

  // Allow registration of servies
  services.forEach((implementation, service) =>
    server.addService(service, implementation as grpc.UntypedServiceImplementation),
  )

  let security = grpc.ServerCredentials.createInsecure()

  if (certificates.enabled) {
    security = grpc.ServerCredentials.createSsl(
      fs.readFileSync(certificates.ca),
      [
        {
          cert_chain: fs.readFileSync(certificates.public),
          private_key: fs.readFileSync(certificates.private),
        },
      ],
      true,
    )
  }

  server.bindAsync(addr, security, (err, port) => {
    if (err) {
      log.error(`could not bind to ${addr}: ${err}`)
      return
    }

    log.info(`server running at ${addr} (port ${port})`)
  })
}
