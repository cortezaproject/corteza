import MakeFilterFn from './shared/filter'
import fs from 'fs'
import path from 'path'
import * as config from '../config'
import { clientScripts as clientScriptsBundler } from '../bundler'
import { Logger } from 'pino'
import Loader, { CommonPath } from '../loader'
import watch from 'node-watch'
import { debounce } from 'lodash'
import GetLastUpdated from './shared/get-last-updated'
import { Script } from '../types'

interface ListFilter {
  query?: string
  resourceType?: string
  eventTypes?: string[]
  bundle?: string
  type?: string
}

interface CtorArgs {
  logger: Logger
  config: Config
  loader: Loader
}

interface Config {
  bundler: {
    outputPath: string
  }
}

/**
 *
 */
export default class ClientScripts {
  private scripts: Script[] = []
  private readonly config: Config
  protected readonly log: Logger
  protected readonly loader: Loader

  /**
   * Service constructor
   */
  constructor({ logger, config, loader }: CtorArgs) {
    this.config = config
    this.loader = loader
    this.log = logger.child({ name: 'services.client-scripts' })
    this.log.debug('initializing')
  }

  // Returns date of the most recently updated script from the set
  get lastUpdated(): Date {
    return GetLastUpdated(this.scripts)
  }

  /**
   * Loads scripts
   *
   * @return {void}
   */
  update(set: Script[]): void {
    // Scripts loaded, replace set
    this.scripts = set
  }

  getBundle(name: string, type: string): Buffer {
    return fs.readFileSync(path.join(this.config.bundler.outputPath, `${name}.${type}.js`))
  }

  /**
   * Returns list of scripts
   */
  list(f: ListFilter = {}): Script[] {
    return this.scripts.filter(MakeFilterFn(f))
  }

  /**
   * Reloads client scripts
   *
   * Note: NOOP when client scripts are disabled
   *
   * Function calls script loader (see loader.ts) and loads all available client scripts
   *
   * Valid scripts (w/o errors) and packed into browser bundles with webpack.
   *
   * It logs (warn) all errors on all scripts and (debug) valid scripts.
   *
   * Client scripts service is then updated with the new list of scripts.
   */
  async process(): Promise<void> {
    if (!this.loader) {
      this.log.debug('no loader: processing disabled')
    }

    this.log.info({ searchPaths: this.loader.searchPaths }, 'reloading client scripts')

    const scripts = await this.loader.scripts()
    const isValid = (s: Script): boolean => s.errors?.length === 0
    const vScripts = scripts.filter(isValid)

    // Make bundles out of all valid scripts
    const scriptListPerBundle = vScripts.reduce(
      (bi, s) => {
        // Extract bundle name from path -- expecting to be 1st subdirectory under 'client-scripts'
        const [, , bundle] = s.name.split(path.sep, 3)

        if (!bi[bundle]) {
          bi[bundle] = []
        }

        bi[bundle].push(s)

        return bi
      },
      {} as { [bundle: string]: Script[] },
    )

    const bootloaderPerBundle = clientScriptsBundler.BootLoader(
      config.bundler.outputPath,
      scriptListPerBundle,
    )
    for (const bnd in bootloaderPerBundle) {
      // Find longest common path for all scripts in the bundle
      // @see https://webpack.js.org/configuration/entry-context/#context
      const ctx = CommonPath(scriptListPerBundle[bnd].map(({ src }) => src))

      this.log.debug({ bundle: bnd }, 'bundling client scripts')

      clientScriptsBundler.Pack(
        bnd,
        bootloaderPerBundle[bnd],
        ctx,
        config.bundler.outputPath,
        this.log,
      )
    }

    // Log errors on all invalid scripts
    scripts
      .filter(s => !isValid(s))
      .forEach(({ src, errors }) => {
        errors?.forEach(error => {
          this.log.warn({ src }, error)
        })
      })

    // Let developer know about valid scripts loaded
    vScripts.forEach(({ src }) => this.log.debug({ src }, 'script ready'))

    // All scripts (even invalid ones) are given to client scripts service
    // we might want to look at errors
    this.update(scripts)

    // Summarize reloading stats
    this.log.info({ valid: vScripts.length, total: scripts.length }, 'processed')
  }

  watch(): void {
    this.log.info('initializing watcher')
    process.on(
      'SIGINT',
      watch(
        this.loader.basePaths(),
        {
          persistent: false,
          recursive: true,
          delay: 200,
          filter: /\.js$/,
        },
        debounce(() => {
          this.log.debug('change detected')
          this.process()
        }, 500),
      ).close,
    )
  }
}
