import { describe, it, expect } from 'vitest'
import { CInputCorredorScript } from '@planetcrust/human-vue/src/components/input'
import { resolveInputComponent } from './registry'
import { resolveViewerComponent } from '../viewers/registry'
import CViewText from '../viewers/CViewText.vue'

// The corredorExec construct declares its `script` input as a
// CorredorScriptSelector; a script name may also arrive typed as CorredorScript.
const TYPES = ['CorredorScriptSelector', 'CorredorScript']

describe('TAQ registries — Corredor script', () => {
  it.each(TYPES)('edits a %s input with the Corredor script selector', type => {
    expect(resolveInputComponent(type)).toBe(CInputCorredorScript)
  })

  it.each(TYPES)('shows a %s value as the script name', type => {
    expect(resolveViewerComponent(type)).toBe(CViewText)
  })
})
