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
    icon: 'icons/swimlane.svg',
    style: 'swimlane',
  },

  expressions: {
    width: 200,
    height: 80,
    icon: 'icons/expressions.svg',
    style: 'expressions',
  },

  function: {
    width: 200,
    height: 80,
    icon: 'icons/function.svg',
    style: 'function',
  },

  iterator: {
    width: 200,
    height: 80,
    icon: 'icons/iterator.svg',
    style: 'iterator',
  },

  'exec-workflow': {
    width: 200,
    height: 80,
    icon: 'icons/exec-workflow.svg',
    style: 'exec-workflow',
  },

  break: {
    width: 200,
    height: 80,
    icon: 'icons/break.svg',
    style: 'break',
  },

  continue: {
    width: 200,
    height: 80,
    icon: 'icons/continue.svg',
    style: 'continue',
  },

  trigger: {
    width: 200,
    height: 80,
    icon: 'icons/trigger.svg',
    style: 'trigger',
  },

  'error-handler': {
    width: 200,
    height: 80,
    icon: 'icons/error-handler.svg',
    style: 'error-handler',
  },

  error: {
    width: 200,
    height: 80,
    icon: 'icons/error.svg',
    style: 'error',
  },

  termination: {
    width: 200,
    height: 80,
    icon: 'icons/termination.svg',
    style: 'termination',
  },

  gatewayExcl: {
    width: 200,
    height: 80,
    icon: 'icons/gateway-exclusive.svg',
    style: 'gatewayExclusive',
  },

  gatewayIncl: {
    width: 200,
    height: 80,
    icon: 'icons/gateway-inclusive.svg',
    style: 'gatewayInclusive',
  },

  gatewayFork: {
    width: 200,
    height: 80,
    icon: 'icons/gateway-parallel.svg',
    style: 'gatewayParallel',
  },

  gatewayJoin: {
    width: 200,
    height: 80,
    icon: 'icons/gateway-parallel.svg',
    style: 'gatewayParallel',
  },

  prompt: {
    width: 200,
    height: 80,
    icon: 'icons/prompt.svg',
    style: 'prompt',
  },

  delay: {
    width: 200,
    height: 80,
    icon: 'icons/delay.svg',
    style: 'delay',
  },

  debug: {
    width: 200,
    height: 80,
    icon: 'icons/debug.svg',
    style: 'debug',
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
  } else {
    return { kind }
  }
}

const gatewayKinds = {
  gatewayExclusive: { kind: 'gateway', ref: 'excl' },
  gatewayInclusive: { kind: 'gateway', ref: 'incl' },
}
