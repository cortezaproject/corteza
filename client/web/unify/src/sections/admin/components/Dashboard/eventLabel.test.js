import { describe, expect, it } from 'vitest'
import { describeEvent, humanizeAction, humanizeType, parseResource } from './eventLabel'

const i18n = {
  te: k =>
    k === 'dashboard.events.resources.compose_namespace' || k === 'dashboard.events.actions.lookup',
  t: k =>
    ({
      'dashboard.events.resources.compose_namespace': 'Namespace',
      'dashboard.events.actions.lookup': 'lookup',
    })[k],
}

describe('event labels', () => {
  it('parses stored resources', () => {
    expect(parseResource('corteza::compose:namespace/123')).toEqual({
      type: 'compose:namespace',
      id: '123',
    })
    expect(parseResource('compose:namespace/*')).toEqual({ type: 'compose:namespace', id: '' })
    expect(parseResource('system:user')).toEqual({ type: 'system:user', id: '' })
  })

  it('spells identifiers out', () => {
    expect(humanizeType('system:user-group')).toBe('user group')
    expect(humanizeAction('execAndWait')).toBe('exec and wait')
    expect(humanizeAction('markAllAsRead')).toBe('mark all as read')
  })

  it('prefers translations and never shows raw identifiers', () => {
    expect(describeEvent(i18n, 'corteza::compose:namespace/*', 'lookup')).toBe('Namespace lookup')
    expect(describeEvent(i18n, 'system:dal-sensitivity-level', 'reloadSensitivityLevels')).toBe(
      'Dal sensitivity level reload sensitivity levels',
    )
  })
})
