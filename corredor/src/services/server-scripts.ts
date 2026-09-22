import MakeFilterFn from './shared/filter'
import { corredor as exec } from '@planetcrust/human-js'
import { Logger } from 'pino'
import watch from 'node-watch'
import { debounce } from 'lodash'
import { Script } from '../types'
import GetLastUpdated from './shared/get-last-updated'
import Loader from '../loader'
import * as Sentry from '@sentry/node'

interface ListFilter {
  query?: string
  resourceType?: string
  eventTypes?: string[]
}

interface CtorArgs {
  logger: Logger
  config: exec.Config
  loader?: Loader
}

/**
 *
 */
export default class ServerScripts {
  protected scripts: Script[] = []
  protected readonly config: exec.Config
  protected readonly log: Logger
  protected readonly loader?: Loader

  /**
   * Service constructor
   */
  constructor({ logger, config, loader }: CtorArgs) {
    this.config = config
    this.loader = loader
    this.log = logger.child({ name: 'services.server-scripts' })
    this.log.debug('initializing')
  }

  // Returns date of the most recently updated script from the set
  get lastUpdated(): Date {
    return GetLastUpdated(this.scripts)
  }

  /**
   * Loads scripts
   */
  update(set: Script[]): void {
    // Scripts loaded, replace set
    this.scripts = set
  }

  /**
   * Returns executable script (or throws error)
   *
   * @param name
   */
  getExecutable(name: string): exec.ExecutableScript {
    const script: Script | undefined = this.scripts.find(s => s.name === name)

    if (script === undefined) {
      throw new Error('script not found')
    }

    if (script.errors && script.errors.length > 0) {
      throw new Error('can not run script with initialization errors')
    }

    if (!script.exec || !(script.exec as exec.ScriptExecFn)) {
      throw new Error('can not run uninitialized script')
    }

    return script as exec.ExecutableScript
  }

  /**
   * Finds and executes the script using current configuration, passed arguments and logger
   *
   * @param name Name of the script
   * @param args Arguments for the script
   * @param log Exec logger to capture and proxy all log.* and console.* calls
   */
  async exec(name: string, args: exec.BaseArgs, log: Logger): Promise<object> {
    return exec.Exec(
      this.getExecutable(name),
      args,
      new exec.Ctx(args, log, { config: this.config }),
    )
  }

  /**
   * Returns list of scripts
   */
  list(f: ListFilter = {}): Script[] {
    return this.scripts.filter(MakeFilterFn(f))
  }

  /**
   * Processes server scripts from loader
   *
   * Note: NOOP when client scripts are disabled
   *
   * Function calls script loader and loads all available server scripts
   * It logs (warn) all errors on all scripts and (debug) valid scripts
   *
   * Server scripts service is then updated with the new list of scripts.
   */
  async process(): Promise<void> {
    if (!this.loader) {
      this.log.debug('no loader: processing disabled')
      return
    }

    this.log.info({ searchPaths: this.loader.searchPaths }, 'reloading server scripts')
    const isValid = (s: Script): boolean => s.errors?.length === 0

    this.loader
      .scripts()
      .then(scripts => {
        this.update(scripts)

        // Log errors on all invalid scripts
        scripts
          .filter(s => !isValid(s))
          .forEach(({ src, errors }) => {
            errors?.forEach(error => {
              this.log.warn({ src }, 'script error: %s', error)
            })
          })

        return scripts.filter(isValid)
      })
      .then(() => {
        // Only a valid script gets its exec function; the raw export of an
        // invalid one may carry functions where the list expects data
        // (an iterator that failed to parse), which no gRPC response can hold.
        for (let s = 0; s < this.scripts.length; s++) {
          if (!isValid(this.scripts[s])) {
            continue
          }

          const { src } = this.scripts[s]
          this.scripts[s] = { ...this.require(src), ...this.scripts[s] }
        }
      })
      .finally(() => {
        this.log.info(
          {
            valid: this.scripts.filter(isValid).length,
            total: this.scripts.length,
          },
          'processed',
        )
      })
  }

  watch(): void {
    this.log.info('initializing watcher')
    process.on(
      'SIGINT',
      watch(
        this.loader?.basePaths() ?? [],
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

  protected require(src: string): Script | undefined {
    delete require.cache[require.resolve(src)]
    try {
      return require(src).default
    } catch (error) {
      Sentry.captureException(error)
      this.log.warn({ src }, 'could not require() script: %s', (error as Error).stack)
      return undefined
    }
  }
}
