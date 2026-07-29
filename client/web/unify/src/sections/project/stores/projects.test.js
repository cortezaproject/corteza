import { createMockSystemAPI, createTestPinia } from '@planetcrust/human-test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useProjectsStore } from './projects'

describe('useProjectsStore', () => {
  let api

  beforeEach(() => {
    api = createMockSystemAPI()
    createTestPinia({
      $SystemAPI: api,
      $ComposeAPI: {},
      $AutomationAPI: {},
    })
  })

  describe('deploymentPlan()', () => {
    // The regression: Go marshals a nil slice as `null`, not `[]`, so a plan
    // with nothing in it arrives as {changes: null}. Spreading the response
    // over array defaults overwrote them with null, and the Publish tab then
    // died on `changes.filter(...)` and reported that it could not work out
    // what the revision changes — on every first revision, which is exactly
    // the case that has an empty plan.
    it('turns null slices into arrays', async () => {
      api.projectGetDeploymentPlan = vi.fn().mockResolvedValue({
        risk: 'safe',
        changes: null,
        suggestedMappings: null,
      })

      const plan = await useProjectsStore().deploymentPlan('1')

      expect(plan.changes).toEqual([])
      expect(plan.suggestedMappings).toEqual([])
      expect(plan.risk).toBe('safe')
    })

    it('survives an empty response body', async () => {
      api.projectGetDeploymentPlan = vi.fn().mockResolvedValue(undefined)

      const plan = await useProjectsStore().deploymentPlan('1')

      expect(plan).toEqual({ risk: 'safe', changes: [], suggestedMappings: [] })
    })

    it('passes a populated plan through unchanged', async () => {
      const changes = [
        { op: 'removed', kind: 'field', name: 'fax', risk: 'dangerous', records: 12 },
      ]
      const suggestedMappings = [{ module: 'customer', fields: [] }]
      api.projectGetDeploymentPlan = vi
        .fn()
        .mockResolvedValue({ risk: 'dangerous', changes, suggestedMappings })

      const plan = await useProjectsStore().deploymentPlan('1')

      expect(plan).toEqual({ risk: 'dangerous', changes, suggestedMappings })
      expect(api.projectGetDeploymentPlan).toHaveBeenCalledWith({ projectID: '1' })
    })
  })
})
