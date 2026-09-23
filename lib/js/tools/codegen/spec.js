import lodash from 'lodash-es'
import fs from 'fs'
import yaml from 'js-yaml'

/**
 * Flattens a service's rest.yaml into one entry per endpoint.
 *
 * Both generators that read the REST specs — the API clients and the Corredor
 * helper layer — read them through here, so the two agree on what an endpoint
 * is and on what its method is called.
 *
 * Returns null when the specs file is not there.
 */
export function loadEndpoints(path) {
  let spec

  try {
    spec = yaml.load(fs.readFileSync(path)).endpoints
  } catch (err) {
    if (err.code === 'ENOENT') {
      return null
    }

    throw err
  }

  if (!spec) {
    return undefined
  }

  return [].concat.apply(
    [],
    spec.map(e => {
      const { get = [], post = [], path = [] } = e.parameters || {}
      const parentGet = get
      const parentPost = post
      const parentPath = path

      return e.apis.map(a => {
        let { get = [], post = [], path = [] } = a.parameters || {}

        const currentPathTemplate = e.path + a.path
        // Filter out path parameters that are defined in yaml but not actually present in the path template
        path = [...parentPath, ...path].filter(v =>
          currentPathTemplate.includes('{' + v.name + '}'),
        )
        get = [...parentGet, ...get]
        post = [...parentPost, ...post]

        const allvars = [...path, ...get, ...post]

        return {
          title: a.title,
          description: a.description,

          entrypoint: e.entrypoint,
          action: a.name,

          fname: lodash.camelCase(e.entrypoint + ' ' + a.name),
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
    }),
  )
}
