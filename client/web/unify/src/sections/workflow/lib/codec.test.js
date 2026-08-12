import { describe, it, expect } from 'vitest'
import { decodeWorkflow, encodeWorkflow } from './codec'

// ─── Minimal factories ────────────────────────────────────────────────────────

function makeStep(overrides = {}) {
  return {
    stepID: '10',
    kind: 'function',
    ref: 'some-fn',
    arguments: [],
    results: [],
    meta: {
      label: 'My step',
      description: 'desc',
      visual: {
        id: '10',
        value: 'My step',
        xywh: [100, 200, 200, 80],
        parent: '1',
      },
    },
    ...overrides,
  }
}

function makeTrigger(overrides = {}) {
  return {
    triggerID: 'T1',
    workflowID: 'W1',
    resourceType: 'compose:record',
    eventType: 'afterCreate',
    enabled: true,
    constraints: [],
    stepID: '10',
    meta: {
      name: 'On create',
      description: 'trigger desc',
      visual: {
        id: 'TR1',
        value: 'On create',
        xywh: [0, 0, 200, 80],
        parent: '1',
        edges: [],
      },
    },
    ...overrides,
  }
}

function makeWorkflow(steps = [], paths = []) {
  return { steps, paths }
}

// ─── Handle ↔ mxGraph coord round-trip ───────────────────────────────────────

describe('handle ↔ mxGraph round-trip', () => {
  const HANDLE_PAIRS = [
    ['source-bottom', 'target-top'],
    ['source-bottom-left', 'target-top-right'],
    ['source-bottom-right', 'target-top-left'],
    ['source-right', 'target-left'],
    ['source-left', 'target-right'],
    ['source-top', 'target-bottom'],
  ]

  it.each(HANDLE_PAIRS)('%s → %s survives encode→decode', (src, tgt) => {
    const step = makeStep()
    const termination = makeStep({
      stepID: '20',
      kind: 'termination',
      meta: { label: 'End', visual: { id: '20', xywh: [300, 200, 200, 80], parent: '1' } },
    })
    const path = {
      parentID: '10',
      childID: '20',
      expr: '',
      meta: {
        visual: {
          id: 'E1',
          style: '', // will be populated by encode
          points: [],
          parent: '1',
        },
      },
    }

    const workflow = makeWorkflow([step, termination], [path])
    const { nodes, edges } = decodeWorkflow(workflow, [])

    // Manually set handles to simulate what VueFlow does after a user drags
    const edge = edges.find(e => e.source === '10' && e.target === '20')
    expect(edge).toBeDefined()
    edge.sourceHandle = src
    edge.targetHandle = tgt

    // Encode back
    const { paths } = encodeWorkflow(nodes, edges)
    expect(paths).toHaveLength(1)
    const encodedStyle = paths[0].meta.visual.style

    // Decode the encoded style
    const workflow2 = makeWorkflow(
      [step, termination],
      [{ ...path, meta: { visual: { id: 'E1', style: encodedStyle, points: [], parent: '1' } } }],
    )
    const { edges: edges2 } = decodeWorkflow(workflow2, [])
    const edge2 = edges2.find(e => e.source === '10' && e.target === '20')
    expect(edge2).toBeDefined()
    expect(edge2.sourceHandle).toBe(src)
    expect(edge2.targetHandle).toBe(tgt)
  })
})

// ─── decodeWorkflow ───────────────────────────────────────────────────────────

describe('decodeWorkflow', () => {
  it('converts steps to nodes', () => {
    const wf = makeWorkflow([makeStep()])
    const { nodes } = decodeWorkflow(wf, [])
    expect(nodes).toHaveLength(1)
    const n = nodes[0]
    expect(n.id).toBe('10')
    expect(n.type).toBe('workflow')
    expect(n.data.kind).toBe('function')
    expect(n.data.ref).toBe('some-fn')
    expect(n.data.label).toBe('My step')
    expect(n.position).toEqual({ x: 100, y: 200 })
  })

  it('termination step → termination node type', () => {
    const { nodes } = decodeWorkflow(makeWorkflow([makeStep({ kind: 'termination' })]), [])
    expect(nodes[0].type).toBe('termination')
  })

  it('visual step → visual node type with zIndex -1', () => {
    const step = makeStep({ kind: 'visual' })
    step.meta.visual.xywh = [0, 0, 400, 240]
    const { nodes } = decodeWorkflow(makeWorkflow([step]), [])
    expect(nodes[0].type).toBe('visual')
    expect(nodes[0].zIndex).toBe(-1)
  })

  it('upgrades legacy "workflow" kind to "exec-workflow"', () => {
    const { nodes } = decodeWorkflow(makeWorkflow([makeStep({ kind: 'workflow' })]), [])
    expect(nodes[0].data.kind).toBe('exec-workflow')
  })

  it('converts paths to edges', () => {
    const step1 = makeStep({
      stepID: '1',
      meta: { visual: { id: '1', xywh: [0, 0, 200, 80], parent: '1' } },
    })
    const step2 = makeStep({
      stepID: '2',
      meta: { visual: { id: '2', xywh: [300, 0, 200, 80], parent: '1' } },
    })
    const path = {
      parentID: '1',
      childID: '2',
      expr: 'x > 0',
      meta: {
        label: 'yes',
        visual: {
          id: 'P1',
          value: 'yes',
          style: 'exitX=0.5;exitY=1;exitDx=0;exitDy=0;entryX=0.5;entryY=0;entryDx=0;entryDy=0;',
          points: [],
          parent: '1',
        },
      },
    }
    const { edges } = decodeWorkflow(makeWorkflow([step1, step2], [path]), [])
    expect(edges).toHaveLength(1)
    const e = edges[0]
    expect(e.source).toBe('1')
    expect(e.target).toBe('2')
    expect(e.data.expr).toBe('x > 0')
    expect(e.sourceHandle).toBe('source-bottom')
    expect(e.targetHandle).toBe('target-top')
  })

  it('converts triggers to trigger nodes', () => {
    const wf = makeWorkflow([makeStep()])
    const { nodes } = decodeWorkflow(wf, [makeTrigger()])
    const triggerNode = nodes.find(n => n.type === 'trigger')
    expect(triggerNode).toBeDefined()
    expect(triggerNode.data.kind).toBe('trigger')
    expect(triggerNode.data.label).toBe('On create')
    expect(triggerNode.data.triggers.resourceType).toBe('compose:record')
  })

  it('trigger with stepID creates an edge to the target step', () => {
    const wf = makeWorkflow([makeStep()])
    const { edges } = decodeWorkflow(wf, [makeTrigger()])
    expect(edges.length).toBeGreaterThan(0)
    const te = edges[0]
    expect(te.source).toBe('TR1')
    expect(te.target).toBe('10')
  })

  it('converts child absolute positions to relative (swimlane)', () => {
    const parent = makeStep({
      stepID: '100',
      kind: 'visual',
      meta: { visual: { id: '100', xywh: [200, 300, 400, 240], parent: '1' } },
    })
    const child = makeStep({
      stepID: '101',
      meta: {
        visual: {
          id: '101',
          xywh: [250, 350, 200, 80], // absolute: parent offset + 50,50
          parent: '100',
        },
      },
    })
    const { nodes } = decodeWorkflow(makeWorkflow([parent, child]), [])
    const childNode = nodes.find(n => n.id === '101')
    expect(childNode.parentNode).toBe('100')
    // Relative position should be absolute - parent position
    expect(childNode.position).toEqual({ x: 50, y: 50 })
  })

  it('orphaned child (missing parent) is promoted to top-level', () => {
    const child = makeStep({
      stepID: '200',
      meta: { visual: { id: '200', xywh: [100, 100, 200, 80], parent: '999' } },
    })
    const { nodes } = decodeWorkflow(makeWorkflow([child]), [])
    expect(nodes[0].parentNode).toBeUndefined()
  })
})

// ─── encodeWorkflow ───────────────────────────────────────────────────────────

describe('encodeWorkflow', () => {
  it('converts nodes back to steps', () => {
    const wf = makeWorkflow([makeStep()])
    const { nodes, edges } = decodeWorkflow(wf, [])
    const { steps } = encodeWorkflow(nodes, edges)
    expect(steps).toHaveLength(1)
    expect(steps[0].kind).toBe('function')
    expect(steps[0].ref).toBe('some-fn')
    expect(steps[0].meta.label).toBe('My step')
  })

  it('encodes legacy "workflow" kind as "exec-workflow"', () => {
    const wf = makeWorkflow([makeStep()])
    const { nodes, edges } = decodeWorkflow(wf, [])
    nodes[0].data.kind = 'workflow'
    const { steps } = encodeWorkflow(nodes, edges)
    expect(steps[0].kind).toBe('exec-workflow')
  })

  it('encodes edges as paths with mxGraph style', () => {
    const step1 = makeStep({
      stepID: '1',
      meta: { visual: { id: '1', xywh: [0, 0, 200, 80], parent: '1' } },
    })
    const step2 = makeStep({
      stepID: '2',
      meta: { visual: { id: '2', xywh: [300, 0, 200, 80], parent: '1' } },
    })
    const { nodes, edges } = decodeWorkflow(
      makeWorkflow(
        [step1, step2],
        [
          {
            parentID: '1',
            childID: '2',
            expr: '',
            meta: { visual: { id: 'E1', style: '', points: [], parent: '1' } },
          },
        ],
      ),
      [],
    )
    edges[0].sourceHandle = 'source-right'
    edges[0].targetHandle = 'target-left'
    const { paths } = encodeWorkflow(nodes, edges)
    expect(paths[0].meta.visual.style).toContain('exitX=1')
    expect(paths[0].meta.visual.style).toContain('exitY=0.5')
    expect(paths[0].meta.visual.style).toContain('entryX=0')
    expect(paths[0].meta.visual.style).toContain('entryY=0.5')
  })

  it('trigger nodes become trigger entries, not steps', () => {
    const wf = makeWorkflow([makeStep()])
    const { nodes, edges } = decodeWorkflow(wf, [makeTrigger()])
    const { steps, triggers } = encodeWorkflow(nodes, edges)
    const triggerStepInSteps = steps.find(s => s.kind === 'trigger')
    expect(triggerStepInSteps).toBeUndefined()
    expect(triggers).toHaveLength(1)
    expect(triggers[0].resourceType).toBe('compose:record')
  })

  it('trigger edges are stored in trigger.meta.visual.edges, not paths', () => {
    const wf = makeWorkflow([makeStep()])
    const { nodes, edges } = decodeWorkflow(wf, [makeTrigger()])
    const { paths, triggers } = encodeWorkflow(nodes, edges)
    // No path should have a trigger node as source
    const triggerNode = nodes.find(n => n.type === 'trigger')
    const pathFromTrigger = paths.find(p => p.parentID === triggerNode.id)
    expect(pathFromTrigger).toBeUndefined()
    // Instead, the edge is in the trigger's visual edges
    expect(triggers[0].meta.visual.edges).toHaveLength(1)
  })

  it('resolves child absolute position from relative (swimlane encode)', () => {
    const parent = makeStep({
      stepID: '100',
      kind: 'visual',
      meta: { visual: { id: '100', xywh: [200, 300, 400, 240], parent: '1' } },
    })
    const child = makeStep({
      stepID: '101',
      meta: { visual: { id: '101', xywh: [250, 350, 200, 80], parent: '100' } },
    })
    const { nodes, edges } = decodeWorkflow(makeWorkflow([parent, child]), [])
    const { steps } = encodeWorkflow(nodes, edges)
    const childStep = steps.find(s => s.stepID === '101')
    // Encoded xywh should be absolute (parent pos + relative pos)
    expect(childStep.meta.visual.xywh[0]).toBe(250) // 200 + 50
    expect(childStep.meta.visual.xywh[1]).toBe(350) // 300 + 50
  })
})

// ─── Full round-trip ──────────────────────────────────────────────────────────

describe('full decode → encode round-trip', () => {
  it('steps survive without mutation', () => {
    const wf = makeWorkflow(
      [
        makeStep(),
        makeStep({
          stepID: '20',
          kind: 'termination',
          meta: { visual: { id: '20', xywh: [400, 200, 200, 80], parent: '1' } },
        }),
      ],
      [
        {
          parentID: '10',
          childID: '20',
          expr: 'ok',
          meta: {
            visual: {
              id: 'P1',
              style: 'exitX=0.5;exitY=1;exitDx=0;exitDy=0;entryX=0.5;entryY=0;entryDx=0;entryDy=0;',
              points: [],
              parent: '1',
            },
          },
        },
      ],
    )

    const { nodes, edges } = decodeWorkflow(wf, [])
    const { steps, paths } = encodeWorkflow(nodes, edges)

    expect(steps.find(s => s.stepID === '10')?.kind).toBe('function')
    expect(steps.find(s => s.stepID === '20')?.kind).toBe('termination')
    expect(paths[0].expr).toBe('ok')
    expect(paths[0].parentID).toBe('10')
    expect(paths[0].childID).toBe('20')
  })

  it('trigger round-trip preserves resourceType and eventType', () => {
    const wf = makeWorkflow([makeStep()])
    const trigger = makeTrigger()
    const { nodes, edges } = decodeWorkflow(wf, [trigger])
    const { triggers } = encodeWorkflow(nodes, edges)
    expect(triggers[0].resourceType).toBe('compose:record')
    expect(triggers[0].eventType).toBe('afterCreate')
  })

  it('snapCoord: off-by-small value snaps to nearest handle', () => {
    // Simulate an mxGraph-written style with a slightly off coordinate
    // (e.g., a user dragged to exitX=0.247 instead of 0.25)
    const step1 = makeStep({
      stepID: '1',
      meta: { visual: { id: '1', xywh: [0, 0, 200, 80], parent: '1' } },
    })
    const step2 = makeStep({
      stepID: '2',
      meta: { visual: { id: '2', xywh: [300, 0, 200, 80], parent: '1' } },
    })
    const wf = makeWorkflow(
      [step1, step2],
      [
        {
          parentID: '1',
          childID: '2',
          expr: '',
          meta: {
            visual: {
              id: 'E1',
              style:
                'exitX=0.247;exitY=1;exitDx=0;exitDy=0;entryX=0.5;entryY=0;entryDx=0;entryDy=0;',
              points: [],
              parent: '1',
            },
          },
        },
      ],
    )
    const { edges } = decodeWorkflow(wf, [])
    // 0.247 should snap to 0.25 → source-bottom-left
    expect(edges[0].sourceHandle).toBe('source-bottom-left')
  })

  it('unknown coords fall back to defaults', () => {
    const step1 = makeStep({
      stepID: '1',
      meta: { visual: { id: '1', xywh: [0, 0, 200, 80], parent: '1' } },
    })
    const step2 = makeStep({
      stepID: '2',
      meta: { visual: { id: '2', xywh: [300, 0, 200, 80], parent: '1' } },
    })
    const wf = makeWorkflow(
      [step1, step2],
      [
        {
          parentID: '1',
          childID: '2',
          expr: '',
          meta: { visual: { id: 'E1', style: '', points: [], parent: '1' } },
        },
      ],
    )
    const { edges } = decodeWorkflow(wf, [])
    // No style → no handle set (null)
    expect(edges[0].sourceHandle).toBeNull()
    expect(edges[0].targetHandle).toBeNull()
  })
})
