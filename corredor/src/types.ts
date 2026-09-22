import { corredor as exec } from '@planetcrust/human-js'
import { Trigger } from './scripts/trigger'
import { Iterator } from './scripts/iterator'

export interface ScriptFile {
  // Script location
  src: string

  // File modification time
  updatedAt: Date

  // Errors detected when requiring the script
  errors?: string[]

  // Client script bundle
  bundle?: string
}

export interface Script extends ScriptFile {
  // Script reference, used by exec
  name: string

  // Display friendly script name
  label?: string

  // Description of a script, lorem ipsum
  description?: string

  // Security settings,
  // run-as settings and simplified RBAC (list of roles that are allowed/denied to execute the script)
  //
  // run-as:
  //   For client scripts: ignored
  //   For server scripts: enforced
  //
  // deny/allow:
  //   For client scripts: used for filtering (who can see what) but ignored
  //   For server scripts: enforced for manual scripts, ignored for implicit
  //
  security?: ScriptSecurity | null

  // When & what trigger this script
  triggers?: Trigger[]

  // Configured iterator
  iterator?: Iterator

  // Code (function) to be executed
  exec?: ScriptFn

  // Name of the exported symbol (default, ...)
  exportName?: string
}

export interface ScriptSecurity {
  runAs?: string
  deny: string[]
  allow: string[]
}

export interface ScriptFn {
  (args: exec.Args, ctx?: exec.Ctx): unknown
}

export interface Watcher {
  watch(): void
}
