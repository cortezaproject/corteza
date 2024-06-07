export function getStyleFromKind ({ kind = '', ref = '' }) {
  let kindRef = kind

  if (['visual', 'gateway'].includes(kind)) {
    kindRef = `${kind}${ref ? (ref[0].toUpperCase() + ref.slice(1).toLowerCase()) : ''}`
  }

  return kindToStyle[kindRef] || {}
}

// The style property tells mxGraph what internal style to use for displaying the specific step
const kindToStyle = {
  visualSwimlane: {
    width: 320,
    height: 160,
    icon: 'swimlane',
    style: 'swimlane',
  },

  expressions: {
    width: 200,
    height: 80,
    icon: 'expressions',
    style: 'expressions',
  },

  function: {
    width: 200,
    height: 80,
    icon: 'function',
    style: 'function',
  },

  iterator: {
    width: 200,
    height: 80,
    icon: 'iterator',
    style: 'iterator',
  },

  'exec-workflow': {
    width: 200,
    height: 80,
    icon: 'exec-workflow',
    style: 'exec-workflow',
  },

  break: {
    width: 200,
    height: 80,
    icon: 'break',
    style: 'break',
  },

  continue: {
    width: 200,
    height: 80,
    icon: 'continue',
    style: 'continue',
  },

  trigger: {
    width: 200,
    height: 80,
    icon: 'trigger',
    style: 'trigger',
  },

  'error-handler': {
    width: 200,
    height: 80,
    icon: 'error-handler',
    style: 'error-handler',
  },

  error: {
    width: 200,
    height: 80,
    icon: 'error',
    style: 'error',
  },

  termination: {
    width: 200,
    height: 80,
    icon: 'termination',
    style: 'termination',
  },

  gatewayExcl: {
    width: 200,
    height: 80,
    icon: 'gateway-exclusive',
    style: 'gatewayExclusive',
  },

  gatewayIncl: {
    width: 200,
    height: 80,
    icon: 'gateway-inclusive',
    style: 'gatewayInclusive',
  },

  gatewayFork: {
    width: 200,
    height: 80,
    icon: 'gateway-parallel',
    style: 'gatewayParallel',
  },

  gatewayJoin: {
    width: 200,
    height: 80,
    icon: 'gateway-parallel',
    style: 'gatewayParallel',
  },

  prompt: {
    width: 200,
    height: 80,
    icon: 'prompt',
    style: 'prompt',
  },

  delay: {
    width: 200,
    height: 80,
    icon: 'delay',
    style: 'delay',
  },

  debug: {
    width: 200,
    height: 80,
    icon: 'debug',
    style: 'debug',
  },

  visualContent: {
    width: 400,
    height: 240,
    icon: 'content',
    style: 'content',
  },
}

// When adding & or copy/pasting a new cell, this is used to determine the kind & ref
export function getKindFromStyle (vertex) {
  const { style } = vertex
  if (!style) {
    return {}
  }

  const kind = style.split(';')[0]

  if (kind.includes('gateway')) {
    if (gatewayKinds[kind]) {
      return gatewayKinds[kind]
    } else {
      // Determine if fork or join
      let inEdgeCount = 0
      let outEdgeCount = 0
      const edges = vertex.edges || []

      edges.forEach(({ source, target }) => {
        if (source.id === vertex.id) {
          outEdgeCount++
        } else if (target.id === vertex.id) {
          inEdgeCount++
        }
      })

      return { kind: 'gateway', ref: (inEdgeCount > outEdgeCount ? 'join' : 'fork') }
    }
  } else if (kind === 'swimlane') {
    return { kind: 'visual', ref: 'swimlane' }
  } else if (kind === 'content') {
    return { kind: 'visual', ref: 'content' }
  } else {
    return { kind }
  }
}

const gatewayKinds = {
  gatewayExclusive: { kind: 'gateway', ref: 'excl' },
  gatewayInclusive: { kind: 'gateway', ref: 'incl' },
}
