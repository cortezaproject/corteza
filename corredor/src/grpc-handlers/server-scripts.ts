import * as grpc from '@grpc/grpc-js'
import pino, { Logger } from 'pino'
import { HandleException } from '../grpc-server'
import Service from '../services/server-scripts'
import { LogToArray } from '../scripts/log-to-array'
import { corredor as exec } from '@planetcrust/human-js'
import IsModifiedSince from './shared/is-modified-since'
import * as Sentry from '@sentry/node'

interface KV {
  [_: string]: string
}

interface ExecRequestWrap {
  request: ExecRequest
  metadata: grpc.Metadata
}

interface ExecRequest {
  name: string
  args?: KV
}

interface ExecResponse {
  result: object
}

interface ListRequestWrap {
  request: ListRequest
  metadata: grpc.Metadata
}

interface ListRequest {
  query?: string
  resourceType?: string
  eventTypes?: string[]
  security?: number
}

interface ListResponse {
  scripts: object[]
}

interface ExecArgsRaw {
  jwt?: string
  [_: string]: unknown
}

/**
 * Decodes exec arguments
 *
 * It assumes values of all properties in the
 * passed object are strings with encoded JSON
 *
 * @param {object} args
 * @returns {object} parsed arguments
 */
export function decodeExecArguments(args: KV | undefined): ExecArgsRaw {
  if (!args) {
    return {}
  }

  for (const k in args) {
    if (!Object.prototype.hasOwnProperty.call(args, k)) {
      continue
    }

    try {
      args[k] = JSON.parse(args[k])
    } catch (e) {
      throw new Error(`Could not parse argument ${k}: ${e}`)
    }
  }

  return args
}

export function encodeExecResult(args: object): KV {
  const enc: KV = {}

  for (const k in args) {
    if (!Object.prototype.hasOwnProperty.call(args, k)) {
      continue
    }

    enc[k] = JSON.stringify((args as { [_: string]: unknown })[k])
  }

  return enc
}

export default function Handler(h: Service, logger: Logger): object {
  logger = logger.child({ name: 'grpc.server-scripts' })

  return {
    Exec(
      { request, metadata }: ExecRequestWrap,
      done: grpc.sendUnaryData<ExecResponse | null>,
    ): void {
      const started = Date.now()

      // name of the script & encoded arguments
      const { name, args: eArgs } = request

      const [requestId] = metadata.get('x-request-id')
      const log = logger.child({ rpc: 'Exec', script: name, requestId })

      Sentry.withScope(scope => {
        if (requestId) {
          scope.setTag('requestId', requestId.toString())
        }

        scope.setTag('script', name)
      })

      let dArgs: ExecArgsRaw = {}

      try {
        // Decode arguments
        // passed in as keys with JSON-encoded values
        // log.debug({ eArgs }, 'encoded arguments')
        dArgs = decodeExecArguments(eArgs)

        log.debug('executing script %s', name)
      } catch (e) {
        HandleException(log, e, done, grpc.status.INVALID_ARGUMENT)
      }

      // global console replacement,
      // will allow us to catch console.* calls and return them to the caller
      const logBuffer = new LogToArray()
      // LogToArray is the write target that collects the script's log lines
      const scriptLogger = pino({}, logBuffer as unknown as pino.DestinationStream)

      // Cast some of the common argument types
      // from plain javascript object to proper classes
      const args = new exec.Args(dArgs)

      try {
        h.getExecutable(name)
      } catch (e) {
        HandleException(log, e, done, grpc.status.NOT_FOUND)
      }

      try {
        h.exec(name, args as exec.BaseArgs, scriptLogger)
          .then(result => {
            const meta = new grpc.Metadata()

            // Map each log line from the executed function to the metadata
            logBuffer.serialize().forEach((l: string) => {
              log.debug(`${name} emitted log: ${l}`)
              meta.add('log', l)
            })

            done(null, { result: encodeExecResult(result) }, meta)
            log.debug(
              {
                duration: Date.now() - started,
              },
              'done',
            )
          })
          .catch(e => {
            if (e instanceof Error && e.message === 'Aborted') {
              HandleException(log, e, done, grpc.status.ABORTED)
            } else {
              HandleException(log, e, done, grpc.status.UNKNOWN)
            }
          })
      } catch (e) {
        HandleException(log, e, done, grpc.status.INTERNAL)
      }
    },

    List(
      { request, metadata }: ListRequestWrap,
      done: grpc.sendUnaryData<ListResponse | null>,
    ): void {
      const { query, resourceType, eventTypes } = request
      const log = logger.child({ rpc: 'List' })

      const filter = {
        query,
        resourceType,
        eventTypes,
      }

      if (!IsModifiedSince(h.lastUpdated, metadata)) {
        done(null, { scripts: [] })
        return
      }

      const scripts = h.list(filter).map(s => ({
        ...s,
        updatedAt: s.updatedAt.toISOString(),
      }))

      log.debug({ filter, total: scripts.length }, 'returning list of server scripts')

      try {
        done(null, { scripts })
      } catch (e) {
        HandleException(log, e, done, grpc.status.INTERNAL)
      }
    },
  }
}
