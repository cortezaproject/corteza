import { Logger } from 'pino'
import watch from 'node-watch'
import { glob } from 'glob'
import fs from 'fs'
import path from 'path'
import { spawnSync, SpawnSyncReturns } from 'child_process'

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
 * Marker in an extension's node_modules, touched after each successful install
 */
const installedMarker = '.corredor-installed'

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

    const files = this.searchPaths
      // run all paths through glob
      .map(sp => glob.sync(path.join(sp, '**', 'package.json'), opt))

      // flatten glob results (expanding search paths) of each search path
      .reduce(flatten, [])

    // search paths overlap (usr and usr/*), so the same file can match twice
    return [...new Set(files)]
  }

  /**
   * Installs dependencies of every extension whose package.json declares any
   * and changed after its last successful install
   */
  installOutdated(): void {
    this.getPackageJsonFiles()
      .filter(pkgJsonPath => this.isOutdated(pkgJsonPath))
      .forEach(pkgJsonPath => this.install(pkgJsonPath))
  }

  protected isOutdated(pkgJsonPath: string): boolean {
    let deps: Record<string, string> | undefined

    try {
      deps = JSON.parse(fs.readFileSync(pkgJsonPath, 'utf8')).dependencies
    } catch (e) {
      this.log.warn({ path: pkgJsonPath, err: e }, 'could not read package.json')
      return false
    }

    if (!deps || Object.keys(deps).length === 0) {
      return false
    }

    const marker = path.join(path.dirname(pkgJsonPath), 'node_modules', installedMarker)
    if (!fs.existsSync(marker)) {
      return true
    }

    return fs.statSync(pkgJsonPath).mtimeMs > fs.statSync(marker).mtimeMs
  }

  protected spawnInstaller(cwd: string): SpawnSyncReturns<Buffer> {
    return spawnSync(this.installer, installArgs[this.installer], {
      cwd,
      shell: process.platform === 'win32',
    })
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

    const log = this.log.child({ path: pkgJsonPath, installer })
    log.info('installing extension dependencies')

    const proc = this.spawnInstaller(path.dirname(pkgJsonPath))
    const stderr = proc.stderr?.toString() ?? ''
    const stdout = proc.stdout?.toString() ?? ''
    const nmdir = path.join(path.dirname(pkgJsonPath), 'node_modules')

    if (!proc.error && proc.status === 0) {
      fs.mkdirSync(nmdir, { recursive: true })
      fs.writeFileSync(path.join(nmdir, installedMarker), '')
    } else {
      log.error(
        { status: proc.status, err: proc.error },
        'could not install extension dependencies',
      )
    }

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
