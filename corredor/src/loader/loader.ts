import fs from 'fs'
import path from 'path'
import { glob } from 'glob'
import { File, FileType } from './types'
import { Script } from '../types'
import ScriptParser from '../scripts/parser'

/**
 * Utility function for flatting w/ Array.reduce
 */
const flatten = <T>(r: T[], p: T[]): T[] => r.concat(p)

export default class Loader {
  readonly searchPaths: string[]
  readonly required: string
  readonly pattern: string
  readonly exclude: RegExp

  /**
   * @param searchPaths - list of paths to search for scripts
   * @param req - valid search-path MUST have this subdir
   * @param pattern - directory structure and filename pattern
   * @param exclude - filenames to skip
   */
  constructor(
    searchPaths: string[],
    req: FileType,
    pattern: string,
    exclude = /\/.+(\.test)\.(js|ts)$/,
  ) {
    this.searchPaths = searchPaths
    this.required = req
    this.pattern = pattern
    this.exclude = exclude
  }

  basePaths(): string[] {
    return (
      this.searchPaths
        // run all paths through glob
        // by appending path.sep at the end of the search string, we tell
        // glob we're only interested in directories
        .map(sp => glob.sync(path.join(sp, this.required) + path.sep))

        // flatten glob results (expanding search paths) of each search path
        .reduce(flatten, [])
    )
  }

  /**
   * Finds files that match the pattern on all search paths
   *
   * This complex chain of map/reduce/filter/glob helps
   * collect all extensions scripts and
   */
  files(): Array<File> {
    // Using map to remove duplicates by overriding
    // existing entries with scripts from the later search-paths
    const uniq = new Map<string, File>()

    const opt = {
      // do not descend into node_modules. ever.
      ignore: '**/node_modules/**',

      // and we're only interested in files
      nodir: true,
    }

    this.basePaths()
      .map(base =>
        glob
          // glob over base + pattern strings
          .sync(path.join(base, this.pattern), opt)

          // exclude out all files we do not want (eg: tests)
          .filter(filepath => !this.exclude.test(filepath))

          // and prepare the final output (ScriptFile
          .map(src => ({
            src,

            // filename w/o search-path prefix
            // we're using this to remove duplicates and for script identification
            ref: path.normalize(src.substring(base.length - this.required.length - 1)),

            updatedAt: fs.statSync(src).mtime,
          })),
      )

      // again, flatten results of the glob (valid script files) on each search path
      .reduce(flatten, [])
      .forEach(sf => uniq.set(sf.ref, sf))

    // Return simple array of ScriptFile
    return [...uniq.values()]
  }

  scripts(): Promise<Array<Script>> {
    return Promise.all(
      this.files().map(async ({ src, ref, updatedAt }): Promise<Script> => {
        let script: Partial<Script>
        try {
          script = await new ScriptParser(fs.readFileSync(src), { src }).parse()
        } catch (e) {
          script = { errors: [(e as Error).toString()] }
        }

        return {
          src,
          name: `${ref}:default`,
          updatedAt,
          triggers: [],
          errors: [],
          ...script,
        }
      }),
    )
  }
}
