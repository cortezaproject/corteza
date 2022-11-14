import { camelCase } from 'lodash'
import fs from 'fs'
import yaml from 'js-yaml'
import handlebars from 'handlebars'
import esformatter from 'esformatter'
import esFormatterOptions from 'esformatter/lib/preset/default'
import { template } from './template'

let path
if (process.argv.length >= 3) {
  path = process.argv[2]
} else {
  // Assume "standard" dev environment
  // where corteza server source could be found
  // next to this lib
  path = '../corteza-server'
}

const dst = `${__dirname}/../../src/api-clients`

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

esFormatterOptions.plugins = ['esformatter-add-trailing-commas']
esFormatterOptions.whiteSpace.after = {
  MethodDefinitionName: 1,
  MethodDefinitionOpeningBrace: 0,
  MethodDefinition: 0,
}

namespaces.forEach(({ path, namespace, className }) => {
  console.log(`Generating '${className}' from specs file '${path}'`)

  let spec

  try {
    spec = yaml.safeLoad(fs.readFileSync(path)).endpoints
  } catch (err) {
    switch (err.code) {
      case 'ENOENT':
        console.error('Could not find specs file')
        return
    }

    throw err
  }

  if (!spec) {
    console.error('Endpoints are undefined')
    return
  }

  const endpoints = [].concat.apply([], spec.map(e => {
    const { get = [], post = [], path = [] } = e.parameters || {}
    const parentGet = get
    const parentPost = post
    const parentPath = path

    return e.apis.map(a => {
      let { get = [], post = [], path = [] } = a.parameters || {}

      path = [...parentPath, ...path]
      get = [...parentGet, ...get]
      post = [...parentPost, ...post]

      const allvars = [...path, ...get, ...post]

      return {
        title: a.title,
        description: a.description,

        fname: camelCase(e.entrypoint + ' ' + a.name),
        fargs: allvars.map(v => v.name),

        pathParams: path.map(v => v.name),

        required: allvars.filter(v => v.required).map(v => v.name),

        method: a.method.toLowerCase(),
        path: (e.path + a.path).replace(/\{/g, '${'),

        hasParams: get.length > 0,
        params: get ? get.map(p => p.name) : [],

        hasData: post.length > 0,
        data: post ? post.map(p => p.name) : [],
      }
    })
  }))

  try {
    const tpl = handlebars.compile(template.trimStart())
    let gen = tpl({ endpoints, className, namespace })
    // gen = esformatter.format(gen, esFormatterOptions)
    gen = gen.replace(/[^\S\n]+$/gm, '')

    fs.writeFileSync(`${dst}/${namespace}.ts`, gen)
  } catch (err) {
    console.error(err)
  }
})
