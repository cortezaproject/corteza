import { describe, expect, it } from 'vitest'
import { describeEvent, humanizeAction } from './eventLabel'

describe('event labels', () => {
  it('reads the resource type and action from the action-log vocabulary', () => {
    expect(describeEvent('corteza::compose:namespace/*', 'lookup')).toBe('Namespace lookup')
    expect(describeEvent('system:user-group', 'members')).toMatch(/^User group /)
  })

  it('spells unknown actions out, keeping acronyms', () => {
    expect(humanizeAction('reloadDALModels')).toBe('reload DAL Models')
    expect(describeEvent('corteza::compose:module', 'reloadDALModels')).toBe(
      'Module reload DAL models',
    )
    expect(describeEvent('system:user', 'execAndWait')).toBe('User exec and wait')
  })

  it('reads an unknown resource type from its last segment', () => {
    expect(describeEvent('system:auth', 'authenticate')).toBe('Auth authenticate')
  })

  it('never shows a raw resource path', () => {
    expect(describeEvent('corteza::compose:record/1/2/3', 'update')).not.toContain('/')
  })
})
