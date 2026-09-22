import * as grpc from '@grpc/grpc-js'
import * as Sentry from '@sentry/node'
import { Logger } from 'pino'

// Defines the structure of our API stack trace entries
interface APIStackTrace {
  file: string
  func: string
  line: number
}

interface APIError extends Omit<Error, 'stack'> {
  stack?: string | Array<APIStackTrace>
}

/**
 * Parses stacktrace string into array of strings
 *
 * @param stack
 */
export function ParseStackStr(stack: string): string[] {
  return (
    stack
      .split('\n')
      // Remove first line
      .slice(1)
      // Trim all lines and remove 'at ' prefix
      .map(l => l.trim().substring(3))
  )
}

/**
 * Parses stacktrace array into array of strings
 *
 * @param stack
 */
export function ParseStackArr(stack: Array<APIStackTrace>): string[] {
  return stack.map(cr => `${cr.func} (${cr.file}:${cr.line}:1)`)
}

/**
 * Handle exceptions and prepare grpc error payload
 */
export function HandleException(
  log: Logger,
  err: unknown,
  done: grpc.sendUnaryData<null>,
  code: grpc.status = grpc.status.ABORTED,
): void {
  const { name, message, stack } = err as APIError
  const grpcErr: grpc.ServiceError = {
    code,
    name,
    message,
    details: '',
    metadata: new grpc.Metadata(),
  }

  if (code !== grpc.status.ABORTED) {
    // Capture exception with sentry
    // but only for non-aborted codes
    Sentry.captureException(err)

    // Log exception
    log.debug({ stack, code }, message)

    if (stack) {
      const metadata = new grpc.Metadata()
      if (typeof stack === 'string') {
        ParseStackStr(stack).forEach(f => metadata.add('stack', f))
      } else if (Array.isArray(stack)) {
        ParseStackArr(stack as Array<APIStackTrace>).forEach(f => metadata.add('stack', f))
      }

      grpcErr.metadata = metadata
    }
  }

  done(grpcErr, null)
}
