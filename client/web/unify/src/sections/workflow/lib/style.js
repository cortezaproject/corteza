export function getStyleFromKind({ kind = '', ref = '' }) {
  let kindRef = kind

  if (['visual', 'gateway'].includes(kind)) {
    kindRef = `${kind}${ref ? ref[0].toUpperCase() + ref.slice(1).toLowerCase() : ''}`
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
    width: 180,
    height: 64,
    icon: 'expressions',
    style: 'expressions',
  },

  function: {
    width: 180,
    height: 64,
    icon: 'function',
    style: 'function',
  },

  iterator: {
    width: 180,
    height: 64,
    icon: 'iterator',
    style: 'iterator',
  },

  'exec-workflow': {
    width: 180,
    height: 64,
    icon: 'exec-workflow',
    style: 'exec-workflow',
  },

  break: {
    width: 180,
    height: 64,
    icon: 'break',
    style: 'break',
  },

  continue: {
    width: 180,
    height: 64,
    icon: 'continue',
    style: 'continue',
  },

  trigger: {
    width: 180,
    height: 64,
    icon: 'trigger',
    style: 'trigger',
  },

  'error-handler': {
    width: 180,
    height: 64,
    icon: 'error-handler',
    style: 'error-handler',
  },

  error: {
    width: 180,
    height: 64,
    icon: 'error',
    style: 'error',
  },

  termination: {
    width: 180,
    height: 64,
    icon: 'termination',
    style: 'termination',
  },

  gatewayExcl: {
    width: 180,
    height: 64,
    icon: 'gateway-exclusive',
    style: 'gatewayExclusive',
  },

  gatewayIncl: {
    width: 180,
    height: 64,
    icon: 'gateway-inclusive',
    style: 'gatewayInclusive',
  },

  gatewayFork: {
    width: 180,
    height: 64,
    icon: 'gateway-parallel',
    style: 'gatewayParallel',
  },

  gatewayJoin: {
    width: 180,
    height: 64,
    icon: 'gateway-parallel',
    style: 'gatewayParallel',
  },

  prompt: {
    width: 180,
    height: 64,
    icon: 'prompt',
    style: 'prompt',
  },

  delay: {
    width: 180,
    height: 64,
    icon: 'delay',
    style: 'delay',
  },

  debug: {
    width: 180,
    height: 64,
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
export function getKindFromStyle(vertex) {
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

      return { kind: 'gateway', ref: inEdgeCount > outEdgeCount ? 'join' : 'fork' }
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
