import type { Ctx } from '@planetcrust/human-js/src/corredor/ctx'

/**
 * The logger type a Corredor exec context is built with.
 */
export type CtxLogger = ConstructorParameters<typeof Ctx>[1]

const levels = ['trace', 'debug', 'info', 'warn', 'error', 'fatal'] as const

const consoleMethod: Record<string, 'debug' | 'info' | 'warn' | 'error'> = {
  trace: 'debug',
  debug: 'debug',
  info: 'info',
  warn: 'warn',
  error: 'error',
  fatal: 'error',
}

/**
 * Browser console in the shape of the logger a Corredor context expects.
 */
export function consoleLogger(): CtxLogger {
  const logger: Record<string, unknown> = {
    level: 'info',
    silent: (): void => undefined,
    child: (): CtxLogger => consoleLogger(),
  }

  levels.forEach(level => {
    logger[level] = (...args: unknown[]): void => {
      console[consoleMethod[level]](...args)
    }
  })

  return logger as unknown as CtxLogger
}
