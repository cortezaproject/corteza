import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

vi.mock('@vue-flow/core', () => ({
  Handle: { name: 'Handle', template: '<div class="handle-stub" />' },
  Position: { Top: 'top', Bottom: 'bottom', Left: 'left', Right: 'right' },
  useVueFlow: () => ({ connectionStartHandle: { value: null } }),
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k, te: () => false }),
}))

vi.mock('../../lib/style', () => ({
  getStyleFromKind: () => ({ icon: 'icon', style: 'plain' }),
}))

vi.mock('../../lib/icon', () => ({
  getIcon: () => 'icon-src',
}))

import WorkflowNode from './WorkflowNode.vue'

const baseProps = {
  id: 'n1',
  data: { kind: 'function' },
  selected: false,
  outCount: 0,
  usedSourceHandles: [],
  usedTargetHandles: [],
  currentTheme: 'light',
  functionTypes: [],
  eventTypes: [],
  issues: {},
}

const mountNode = (overrides = {}) =>
  mount(WorkflowNode, {
    props: { ...baseProps, ...overrides },
    global: {
      mocks: { $t: k => k },
      directives: { tooltip: {} },
    },
  })

describe('WorkflowNode --hoverable class binding', () => {
  it('plain step is hoverable with 0 outbound edges', () => {
    expect(mountNode().classes()).toContain('workflow-node--hoverable')
  })

  it('plain step drops --hoverable after the first edge', () => {
    expect(mountNode({ outCount: 1 }).classes()).not.toContain('workflow-node--hoverable')
  })

  it('exclusive gateway stays hoverable past the first edge', () => {
    expect(
      mountNode({ data: { kind: 'gateway', ref: 'excl' }, outCount: 5 }).classes(),
    ).toContain('workflow-node--hoverable')
  })

  it('inclusive gateway stays hoverable past the first edge', () => {
    expect(
      mountNode({ data: { kind: 'gateway', ref: 'incl' }, outCount: 3 }).classes(),
    ).toContain('workflow-node--hoverable')
  })

  it('fork gateway stays hoverable past the first edge', () => {
    expect(
      mountNode({ data: { kind: 'gateway', ref: 'fork' }, outCount: 4 }).classes(),
    ).toContain('workflow-node--hoverable')
  })

  it('iterator is hoverable with 1 edge, drops at 2', () => {
    expect(
      mountNode({ data: { kind: 'iterator' }, outCount: 1 }).classes(),
    ).toContain('workflow-node--hoverable')
    expect(
      mountNode({ data: { kind: 'iterator' }, outCount: 2 }).classes(),
    ).not.toContain('workflow-node--hoverable')
  })

  it('error-handler is hoverable with 1 edge, drops at 2', () => {
    expect(
      mountNode({ data: { kind: 'error-handler' }, outCount: 1 }).classes(),
    ).toContain('workflow-node--hoverable')
    expect(
      mountNode({ data: { kind: 'error-handler' }, outCount: 2 }).classes(),
    ).not.toContain('workflow-node--hoverable')
  })

  it('join gateway is hoverable with 0 edges, drops at 1 (cap 1)', () => {
    expect(
      mountNode({ data: { kind: 'gateway', ref: 'join' }, outCount: 0 }).classes(),
    ).toContain('workflow-node--hoverable')
    expect(
      mountNode({ data: { kind: 'gateway', ref: 'join' }, outCount: 1 }).classes(),
    ).not.toContain('workflow-node--hoverable')
  })
})
