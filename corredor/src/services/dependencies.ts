import { Logger } from 'pino'
import watch from 'node-watch'
import { glob } from 'glob'
import path from 'path'
import { spawnSync } from 'child_process'

interface CtorArgs {
  logger: Logger
  searchPaths: string[]
  installer?: string
}

interface WatchCallback {
  (path: string): void
}

/**
 * Package managers that can install extension dependencies
 */
type Installer = 'yarn' | 'npm' | 'pnpm'

const installers: Installer[] = ['yarn', 'npm', 'pnpm']

/**
 * Install arguments per package manager
 *
 * Only yarn emits the newline-delimited JSON that install() parses for progress steps.
 */
const installArgs: Record<Installer, string[]> = {
  yarn: ['install', '--json', '--force', '--silent', '--emoji', 'off', '--no-progress'],
  npm: ['install', '--force', '--silent', '--no-audit', '--no-fund', '--no-progress'],
  pnpm: ['install', '--force', '--reporter', 'silent'],
}

/**
 * Utility function for flatting w/ Array.reduce
 */
const flatten = (r: string[], p: string[]): string[] => r.concat(p)

/**
 * True when the binary can be executed
 */
function isExecutable(bin: string): boolean {
  const probe = spawnSync(bin, ['--version'], { shell: process.platform === 'win32' })
  return !probe.error && probe.status === 0
}

export default class Dependencies {
  protected searchPaths: string[] = []
  protected readonly log: Logger
  protected readonly installer: Installer

  constructor({ logger, searchPaths, installer }: CtorArgs) {
    this.searchPaths = searchPaths
    this.log = logger.child({ name: 'services.dependencies' })
    this.installer = this.resolveInstaller(installer)
    this.log.debug({ installer: this.installer }, 'initializing')
  }

  /**
   * Picks the package manager used for extension dependencies
   *
   * Corteza extensions ship a yarn.lock, so yarn is preferred when present.
   */
  protected resolveInstaller(requested?: string): Installer {
    if (requested) {
      if (!installers.includes(requested as Installer)) {
        this.log.warn(
          { requested, supported: installers },
          'unsupported dependency installer, falling back to auto-detection',
        )
      } else if (isExecutable(requested)) {
        return requested as Installer
      } else {
        this.log.warn(
          { requested },
          'requested dependency installer not found on PATH, falling back to auto-detection',
        )
      }
    }

    const found = installers.find(isExecutable)
    if (!found) {
      this.log.error({ supported: installers }, 'no dependency installer found on PATH')
      return 'npm'
    }

    return found
  }

  /**
   * Find all package.json files under all search paths
   */
  getPackageJsonFiles(): string[] {
    const opt = {
      // Ignore node_modules
      ignore: '**/node_modules/**',

      // Only interested in files
      nodir: true,
    }

    return (
      this.searchPaths
        // run all paths through glob
        .map(sp => glob.sync(path.join(sp, '**', 'package.json'), opt))

        // flatten glob results (expanding search paths) of each search path
        .reduce(flatten, [])
    )
  }

  /**
   * Installs packages
   *
   * @todo check what happens with require cache after new yarn install
   *       it might be a problem
   *
   * @param pkgJsonPath - path to package.json
   */
  install(pkgJsonPath: string): void {
    const installer = this.installer

    const opts = {
      cwd: path.dirname(pkgJsonPath),
      shell: process.platform === 'win32',
    }

    const log = this.log.child({ path: pkgJsonPath, installer })
    log.info('installing extension dependencies')

    const proc = spawnSync(installer, installArgs[installer], opts)
    const stderr = proc.stderr?.toString() ?? ''
    const stdout = proc.stdout?.toString() ?? ''
    const nmdir = path.join(path.dirname(pkgJsonPath), 'node_modules')

    if (stderr.length > 0) {
      log.error('err' + stderr)
    } else if (stdout.length > 0) {
      stdout
        .split('\n')
        .filter(line => line.length > 0)
        .forEach(line => {
          if (installer !== 'yarn') {
            log.debug(line)
            return
          }

          try {
            const { type, data } = JSON.parse(line)
            if (type === 'step') {
              log.debug(data.message)
            }
          } catch {
            log.debug(line)
          }
        })
    }

    // Purge require cache -- remove all files that are in the
    //
    Object.getOwnPropertyNames(require.cache)
      .filter(path => path.startsWith(nmdir))
      .forEach(filename => {
        delete require.cache[filename]
      })
  }

  /**
   * Watches all loaded package.json files
   *
   * Function installs dependencies on change and resolves
   *
   * @return path to changed file
   */
  watch(callback: WatchCallback): void {
    this.log.info('initializing watcher')
    process.on(
      'SIGINT',
      watch(
        this.getPackageJsonFiles(),
        {
          persistent: false,
          recursive: false,
          delay: 1000,
          filter: /\/package\.json$/,
        },
        (eventType, filename) => {
          switch (eventType) {
            case 'update':
              delete require.cache[require.resolve(filename)]
              this.install(filename)
              callback(filename)
              break
          }
        },
      ).close,
    )
  }
}
