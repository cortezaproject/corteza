import pino from 'pino'
import { logger as loggerConfig } from './config'

const transport = pino.transport({
  target: 'pino-pretty',
  options: {
    colorize: true,
  },
})

export default pino(
  {
    enabled: loggerConfig.enabled,
    base: null,
    level: loggerConfig.level,
  },
  loggerConfig.prettyPrint ? transport : undefined,
)
