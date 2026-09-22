import * as ts from 'typescript'
import vm from 'vm'
import { Script, ScriptSecurity } from '../types'
import MakeIterator from './iterator'
import MakeTrigger from './trigger'

interface RawScriptSecurity {
  runAs?: string
  allow: string | string[]
  deny: string | string[]
}

export default class ScriptParser {
  source: ts.SourceFile

  constructor(
    source: Buffer | string,
    { src = 'buffer', languageVersion = ts.ScriptTarget.Latest } = {},
  ) {
    this.source = ts.createSourceFile(src, source.toString(), languageVersion)
  }

  async parse(): Promise<Partial<Script>> {
    return this.findDefaultExport(this.source)
      .then(n => this.parseDefaultExport(n))
      .then(def => this.parseExportedObject(def))
  }

  /**
   * Parses first level of AST
   *
   * Tries to find 'export default { ... }'
   *
   * @param n
   */
  protected findDefaultExport(n: ts.Node): Promise<ts.Node> {
    return new Promise((resolve, reject) => {
      n.forEachChild(n => {
        if (ts.isExportAssignment(n)) {
          n.forEachChild(o => {
            if (ts.isObjectLiteralExpression(o)) {
              resolve(n)
            }
          })

          reject(new SyntaxError('expecting object as default export'))
        }
      })

      reject(new SyntaxError('expecting default export'))
    })
  }

  /**
   * Evaluates extracted default export and returns result
   *
   * This is needed so we can safely parse all scripts, even if they are
   * completely broken or import large modules.
   *
   * @param n
   */
  protected parseDefaultExport(n: ts.Node): { [_: string]: unknown } {
    const source = n.getText(this.source)

    const tm = ts.transpileModule(source, {
      compilerOptions: {
        allowJs: true,
        checkJs: true,
        allowSyntheticDefaultImports: true,
        target: ts.ScriptTarget.Latest,
        module: ts.ModuleKind.CommonJS,
      },
    })

    const snippet = `var exports = {};${tm.outputText}`
    const def = new vm.Script(snippet).runInNewContext({})
    if (typeof def !== 'object') {
      throw new SyntaxError('expecting object as default export')
    }

    return def
  }

  protected parseExportedObject(o: { [_: string]: unknown }): Partial<Script> {
    const s: Partial<Script> = {}

    if (Object.prototype.hasOwnProperty.call(o, 'label') && typeof o.label === 'string') {
      s.label = o.label
    }

    if (
      Object.prototype.hasOwnProperty.call(o, 'description') &&
      typeof o.description === 'string'
    ) {
      s.description = o.description
    }

    if (Object.prototype.hasOwnProperty.call(o, 'iterator') && o.iterator) {
      s.iterator = MakeIterator(o.iterator) ?? undefined
      if (!s.iterator) {
        throw new SyntaxError('iterator not defined')
      }
    } else if (Object.prototype.hasOwnProperty.call(o, 'triggers') && o.triggers) {
      s.triggers = MakeTrigger(o.triggers) ?? []
      if (!s.triggers || s.triggers.length === 0) {
        throw new SyntaxError('triggers not defined')
      }
    } else {
      throw new SyntaxError('triggers or iterator definition missing')
    }

    if (!Object.prototype.hasOwnProperty.call(o, 'exec')) {
      throw new SyntaxError('exec function missing')
    } else if (typeof o.exec !== 'function') {
      throw new SyntaxError('exec is not a function')
    }

    if (Object.prototype.hasOwnProperty.call(o, 'security')) {
      s.security = this.parseSecurity(o.security as RawScriptSecurity)
    }

    // Merge resolved & the rest
    return s
  }

  protected parseSecurity(sec: unknown): ScriptSecurity | null {
    if (typeof sec === 'undefined') {
      return null
    }

    if (typeof sec === 'string') {
      return { runAs: sec, deny: [], allow: [] }
    }

    if (typeof sec === 'object') {
      const { allow, deny, runAs } = sec as RawScriptSecurity

      if (!allow && !deny && !runAs) {
        return null
      }

      const out: ScriptSecurity = { runAs, allow: [], deny: [] }

      if (allow) {
        out.allow = Array.isArray(allow) ? allow : [allow]
      }

      if (deny) {
        out.deny = Array.isArray(deny) ? deny : [deny]
      }

      return out
    }

    throw new SyntaxError('unknown security definition format')
  }
}
