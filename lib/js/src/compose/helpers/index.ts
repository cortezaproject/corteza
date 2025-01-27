import yaml from 'js-yaml'
import fs from 'fs'
import { Module } from '../types/module'

export function getModuleFromYaml (moduleName: string, yamlPath: string): Module|undefined {
  const data = yaml.loadAll(fs.readFileSync(yamlPath, 'utf8')) as Array<{ modules: { [key: string]: any } }>
  const mod = data[0].modules[moduleName]
  if (mod) {
    // Convert fields from object to array
    mod.fields = Object.keys(mod.fields).map((k) => mod.fields[k])
    return new Module(mod)
  } else {
    return undefined
  }
}
