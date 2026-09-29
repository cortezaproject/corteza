import { describe, it } from 'mocha'
import { expect } from 'chai'

import pino from 'pino'
import fs from 'fs'
import os from 'os'
import path from 'path'
import { SpawnSyncReturns } from 'child_process'
import Dependencies from './dependencies'

// Records every install instead of running a package manager
class RecordingDependencies extends Dependencies {
  installed: string[] = []
  status = 0

  protected spawnInstaller(cwd: string): SpawnSyncReturns<Buffer> {
    this.installed.push(cwd)
    fs.mkdirSync(path.join(cwd, 'node_modules'), { recursive: true })
    return {
      pid: 0,
      output: [],
      stdout: Buffer.from(''),
      stderr: Buffer.from(''),
      status: this.status,
      signal: null,
    }
  }
}

function extension(root: string, name: string, pkg: object): string {
  const dir = path.join(root, name)
  fs.mkdirSync(dir, { recursive: true })
  fs.writeFileSync(path.join(dir, 'package.json'), JSON.stringify(pkg))
  return dir
}

describe('extension dependencies on start', () => {
  let root: string
  let withDeps: string
  let svc: RecordingDependencies

  beforeEach(() => {
    root = fs.mkdtempSync(path.join(os.tmpdir(), 'corredor-deps-'))
    withDeps = extension(root, 'with-deps', { dependencies: { lodash: '^4.0.0' } })
    extension(root, 'without-deps', { name: 'without-deps' })
    extension(root, 'dev-deps-only', { devDependencies: { eslint: '^9.0.0' } })

    // the two default search paths overlap: the root and every folder in it
    svc = new RecordingDependencies({
      logger: pino({ level: 'silent' }),
      searchPaths: [path.join(root, '*'), root],
      installer: 'npm',
    })
  })

  afterEach(() => {
    fs.rmSync(root, { recursive: true, force: true })
  })

  it('installs each extension that declares dependencies, once', () => {
    svc.installOutdated()
    expect(svc.installed).to.deep.equal([withDeps])
  })

  it('skips an extension installed since its package.json last changed', () => {
    svc.installOutdated()
    svc.installOutdated()
    expect(svc.installed).to.have.lengthOf(1)
  })

  it('installs again once its package.json changes', () => {
    svc.installOutdated()

    const later = new Date(Date.now() + 60_000)
    fs.utimesSync(path.join(withDeps, 'package.json'), later, later)

    svc.installOutdated()
    expect(svc.installed).to.have.lengthOf(2)
  })

  it('retries an install that failed', () => {
    svc.status = 1
    svc.installOutdated()
    svc.status = 0
    svc.installOutdated()
    expect(svc.installed).to.have.lengthOf(2)
  })
})
