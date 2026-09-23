import fs from 'fs'
import handlebars from 'handlebars'
import { template } from './template.js'
import { loadEndpoints } from './spec.js'
import { fileURLToPath } from 'url'
import { dirname, join } from 'path'

// Get equivalent of __dirname in ES modules
const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

let path
if (process.argv.length >= 3) {
  path = process.argv[2]
} else {
  // Assume "standard" dev environment
  // where human server source could be found
  // next to this lib
  path = '../../server'
}

const dst = join(__dirname, '../../src/api-clients')

const namespaces = [
  {
    path: `${path}/system/rest.yaml`,
    namespace: 'system',
    className: 'System',
  },
  {
    path: `${path}/compose/rest.yaml`,
    namespace: 'compose',
    className: 'Compose',
  },
  {
    path: `${path}/federation/rest.yaml`,
    namespace: 'federation',
    className: 'Federation',
  },
  {
    path: `${path}/automation/rest.yaml`,
    namespace: 'automation',
    className: 'Automation',
  },
]

namespaces.forEach(({ path, namespace, className }) => {
  console.log(`Generating '${className}' from specs file '${path}'`)

  const endpoints = loadEndpoints(path)

  if (endpoints === null) {
    console.error('Could not find specs file')
    return
  }

  if (!endpoints) {
    console.error('Endpoints are undefined')
    return
  }

  try {
    const tpl = handlebars.compile(template.trimStart())
    let gen = tpl({ endpoints, className, namespace })
    // Remove trailing whitespace from lines while preserving newlines
    gen = gen.replace(/[^\S\n]+$/gm, '')

    fs.writeFileSync(`${dst}/${namespace}.ts`, gen)
  } catch (err) {
    console.error(err)
  }
})
