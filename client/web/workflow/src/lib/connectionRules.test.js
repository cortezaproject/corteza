import { describe, it, expect } from 'vitest'
import { getMaxOutbound, isValidConnection } from './connectionRules'

const node = (id, type, kind, ref) => ({ id, type, data: { kind, ref } })
const edge = (id, source, target, sourceHandle = null, targetHandle = null) => ({
  id, source, target, sourceHandle, targetHandle,
})

describe('getMaxOutbound', () => {
  it.each([
    ['gateway', 'excl', Infinity],
    ['gateway', 'incl', Infinity],
    ['gateway', 'fork', Infinity],
    ['gateway', 'join', 1],
    ['iterator', undefined, 2],
    ['error-handler', undefined, 2],
    ['function', undefined, 1],
    ['expressions', undefined, 1],
    ['delay', undefined, 1],
    ['prompt', undefined, 1],
  ])('kind=%s ref=%s → %s', (kind, ref, want) => {
    expect(getMaxOutbound({ data: { kind, ref } })).toBe(want)
  })

  it('returns 0 for missing node', () => {
    expect(getMaxOutbound(null)).toBe(0)
    expect(getMaxOutbound(undefined)).toBe(0)
  })
})

describe('isValidConnection', () => {
  it('rejects connections involving a visual node on either side', () => {
    const nodes = [node('v', 'visual'), node('w', 'workflow', 'function')]
    expect(isValidConnection({ source: 'v', target: 'w' }, { nodes, edges: [] })).toBe(false)
    expect(isValidConnection({ source: 'w', target: 'v' }, { nodes, edges: [] })).toBe(false)
  })

  it('rejects targeting a trigger', () => {
    const nodes = [node('a', 'workflow', 'function'), node('t', 'trigger')]
    expect(isValidConnection({ source: 'a', target: 't' }, { nodes, edges: [] })).toBe(false)
  })

  it('rejects sourcing from a termination', () => {
    const nodes = [node('term', 'termination'), node('a', 'workflow', 'function')]
    expect(isValidConnection({ source: 'term', target: 'a' }, { nodes, edges: [] })).toBe(false)
  })

  it('caps termination inbound at 1', () => {
    const nodes = [
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
      node('term', 'termination'),
    ]
    const edges = [edge('e1', 'a', 'term')]
    expect(isValidConnection({ id: 'e2', source: 'b', target: 'term' }, { nodes, edges })).toBe(false)
  })

  it('caps trigger outbound at 1', () => {
    const nodes = [
      node('t', 'trigger'),
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
    ]
    const edges = [edge('e1', 't', 'a')]
    expect(isValidConnection({ id: 'e2', source: 't', target: 'b' }, { nodes, edges })).toBe(false)
  })

  it('caps a plain workflow step at 1 outbound', () => {
    const nodes = [
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
      node('c', 'workflow', 'function'),
    ]
    const edges = [edge('e1', 'a', 'b', 'sh1')]
    expect(isValidConnection(
      { id: 'e2', source: 'a', target: 'c', sourceHandle: 'sh2' },
      { nodes, edges },
    )).toBe(false)
  })

  it('allows >1 outbound from an exclusive gateway', () => {
    const nodes = [
      node('g', 'workflow', 'gateway', 'excl'),
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
      node('c', 'workflow', 'function'),
      node('d', 'workflow', 'function'),
    ]
    const edges = [
      edge('e1', 'g', 'a', 'sh1'),
      edge('e2', 'g', 'b', 'sh2'),
      edge('e3', 'g', 'c', 'sh3'),
    ]
    expect(isValidConnection(
      { id: 'e4', source: 'g', target: 'd', sourceHandle: 'sh4' },
      { nodes, edges },
    )).toBe(true)
  })

  it('allows >1 outbound from inclusive and fork gateways', () => {
    for (const ref of ['incl', 'fork']) {
      const nodes = [
        node('g', 'workflow', 'gateway', ref),
        node('a', 'workflow', 'function'),
        node('b', 'workflow', 'function'),
      ]
      const edges = [edge('e1', 'g', 'a', 'sh1')]
      expect(isValidConnection(
        { id: 'e2', source: 'g', target: 'b', sourceHandle: 'sh2' },
        { nodes, edges },
      )).toBe(true)
    }
  })

  it('caps a join gateway outbound at 1', () => {
    const nodes = [
      node('g', 'workflow', 'gateway', 'join'),
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
    ]
    const edges = [edge('e1', 'g', 'a', 'sh1')]
    expect(isValidConnection(
      { id: 'e2', source: 'g', target: 'b', sourceHandle: 'sh2' },
      { nodes, edges },
    )).toBe(false)
  })

  it('caps iterator outbound at 2', () => {
    const nodes = [
      node('it', 'workflow', 'iterator'),
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
      node('c', 'workflow', 'function'),
    ]
    const oneEdge = [edge('e1', 'it', 'a', 'sh1')]
    expect(isValidConnection(
      { id: 'e2', source: 'it', target: 'b', sourceHandle: 'sh2' },
      { nodes, edges: oneEdge },
    )).toBe(true)

    const twoEdges = [edge('e1', 'it', 'a', 'sh1'), edge('e2', 'it', 'b', 'sh2')]
    expect(isValidConnection(
      { id: 'e3', source: 'it', target: 'c', sourceHandle: 'sh3' },
      { nodes, edges: twoEdges },
    )).toBe(false)
  })

  it('caps error-handler outbound at 2', () => {
    const nodes = [
      node('eh', 'workflow', 'error-handler'),
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
      node('c', 'workflow', 'function'),
    ]
    const twoEdges = [edge('e1', 'eh', 'a', 'sh1'), edge('e2', 'eh', 'b', 'sh2')]
    expect(isValidConnection(
      { id: 'e3', source: 'eh', target: 'c', sourceHandle: 'sh3' },
      { nodes, edges: twoEdges },
    )).toBe(false)
  })

  it('rejects self-loops', () => {
    const nodes = [node('a', 'workflow', 'function')]
    expect(isValidConnection({ source: 'a', target: 'a' }, { nodes, edges: [] })).toBe(false)
  })

  it('rejects a second edge from the same source handle', () => {
    const nodes = [
      node('g', 'workflow', 'gateway', 'excl'),
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
    ]
    const edges = [edge('e1', 'g', 'a', 'source-bottom')]
    expect(isValidConnection(
      { id: 'e2', source: 'g', target: 'b', sourceHandle: 'source-bottom' },
      { nodes, edges },
    )).toBe(false)
  })

  it('rejects a second edge into the same target handle', () => {
    const nodes = [
      node('g', 'workflow', 'gateway', 'excl'),
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
    ]
    const edges = [edge('e1', 'a', 'b', null, 'target-top')]
    expect(isValidConnection(
      { id: 'e2', source: 'g', target: 'b', sourceHandle: 'sh', targetHandle: 'target-top' },
      { nodes, edges },
    )).toBe(false)
  })

  it('excludes the edge being reconnected when counting caps', () => {
    const nodes = [
      node('a', 'workflow', 'function'),
      node('b', 'workflow', 'function'),
      node('c', 'workflow', 'function'),
    ]
    const edges = [edge('e1', 'a', 'b', 'source-bottom')]
    // Reconnecting e1 from a→b to a→c with the same handle: cap=1 should
    // see 0 edges (the one being moved is excluded) and allow.
    expect(isValidConnection(
      { id: 'e1', source: 'a', target: 'c', sourceHandle: 'source-bottom' },
      { nodes, edges, edgeUpdatingId: 'e1' },
    )).toBe(true)
  })
})
