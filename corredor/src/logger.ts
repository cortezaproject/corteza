import pino from 'pino'
import { logger as loggerConfig } from './config'

export default pino(
  {
    enabled: loggerConfig.enabled,
    base: null,
    level: loggerConfig.level,
  },
  loggerConfig.prettyPrint
    ? pino.transport({
        target: 'pino-pretty',
        options: {
          colorize: true,
        },
      })
    : undefined,
)
