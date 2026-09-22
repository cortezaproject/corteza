import fs from 'fs'
import { bundler } from '../config'
import { canAccessPath } from './fs'
import logger from '../logger'

/**
 * Checks for bundle output dir, tries to create it and verify if writable
 */
export default function (): void {
  const output = bundler.outputPath
  const log = logger.child({ name: 'check' })

  // Make directory if it does not exists

  if (!fs.existsSync(output)) {
    fs.mkdirSync(output, { recursive: true })
    log.info({ output }, 'making bundler output directory')
  }

  if (!canAccessPath(output, fs.constants.W_OK)) {
    bundler.enabled = false
    // extensions.vueComponents.enabled = false
    log.warn({ output }, 'output location not writable, disabling bundler')
    return
  }

  log.debug({ output }, 'bundler enabled')
}
