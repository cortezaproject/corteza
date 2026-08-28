// Editor canvas grid: the background dot spacing and the snap increment for
// every node placement, so a step always lands on a dot. Half the step node's
// height, as in the TAQ builder.
export const CANVAS_GRID = 32

// Node boxes are whole numbers of grid cells, so every edge lands on a dot
// whenever the corner does.
const cells = n => n * CANVAS_GRID

// One box for every step kind: six cells by two.
export const STEP_NODE = { width: cells(6), height: cells(2) }

// The two resizable kinds, and the smallest a resize may leave one.
export const VISUAL_NODE = {
  swimlane: { width: cells(10), height: cells(5) },
  content: { width: cells(12), height: cells(7) },
  min: { width: cells(5), height: cells(2) },
}

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
    ...VISUAL_NODE.swimlane,
    icon: 'swimlane',
    style: 'swimlane',
  },

  expressions: {
    ...STEP_NODE,
    icon: 'expressions',
    style: 'expressions',
  },

  function: {
    ...STEP_NODE,
    icon: 'function',
    style: 'function',
  },

  iterator: {
    ...STEP_NODE,
    icon: 'iterator',
    style: 'iterator',
  },

  'exec-workflow': {
    ...STEP_NODE,
    icon: 'exec-workflow',
    style: 'exec-workflow',
  },

  break: {
    ...STEP_NODE,
    icon: 'break',
    style: 'break',
  },

  continue: {
    ...STEP_NODE,
    icon: 'continue',
    style: 'continue',
  },

  trigger: {
    ...STEP_NODE,
    icon: 'trigger',
    style: 'trigger',
  },

  'error-handler': {
    ...STEP_NODE,
    icon: 'error-handler',
    style: 'error-handler',
  },

  error: {
    ...STEP_NODE,
    icon: 'error',
    style: 'error',
  },

  termination: {
    ...STEP_NODE,
    icon: 'termination',
    style: 'termination',
  },

  gatewayExcl: {
    ...STEP_NODE,
    icon: 'gateway-exclusive',
    style: 'gatewayExclusive',
  },

  gatewayIncl: {
    ...STEP_NODE,
    icon: 'gateway-inclusive',
    style: 'gatewayInclusive',
  },

  gatewayFork: {
    ...STEP_NODE,
    icon: 'gateway-parallel',
    style: 'gatewayParallel',
  },

  gatewayJoin: {
    ...STEP_NODE,
    icon: 'gateway-parallel',
    style: 'gatewayParallel',
  },

  prompt: {
    ...STEP_NODE,
    icon: 'prompt',
    style: 'prompt',
  },

  delay: {
    ...STEP_NODE,
    icon: 'delay',
    style: 'delay',
  },

  debug: {
    ...STEP_NODE,
    icon: 'debug',
    style: 'debug',
  },

  visualContent: {
    ...VISUAL_NODE.content,
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
