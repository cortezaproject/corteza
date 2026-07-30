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

  // The regression: the store destructured only `nodes` and `edges`, so the
  // endpoint's whole account of what it could NOT resolve — a chatbot pointing
  // at another revision's agent, a TAQ binding a branch copy dropped — was
  // thrown away between the API and the canvas, and a broken project rendered
  // as a correct one with one fewer line on it.
  describe('graph()', () => {
    const call = async response => {
      api.projectGraph = vi.fn().mockResolvedValue(response)
      return useProjectsStore().graph('42')
    }

    it('carries missing references through, in the section kind vocabulary', async () => {
      const g = await call({
        nodes: [{ id: '1', kind: 'chart', name: 'Revenue' }],
        edges: [],
        missing: [
          {
            sourceID: '1',
            kind: 'corteza::compose:module',
            targetID: '9',
            reason: 'chart-module',
          },
        ],
      })

      expect(g.missing).toEqual([
        {
          sourceID: '1',
          kind: 'module',
          resourceType: 'corteza::compose:module',
          targetID: '9',
          targetIdent: '',
          reason: 'chart-module',
        },
      ])
    })

    // Both connection resource types draw as a "connection" node; the reverse
    // of the AI-system map can only carry one of them, which is why the graph
    // direction has its own table.
    it('resolves the tenant-level connection type to the connection kind', async () => {
      const g = await call({
        nodes: [],
        edges: [],
        warnings: [
          {
            sourceID: '1',
            kind: 'corteza::system:dal-connection',
            reason: 'step-connection',
            path: 'steps[2].args.connection',
          },
        ],
      })

      expect(g.warnings[0]).toMatchObject({
        kind: 'connection',
        reason: 'step-connection',
        path: 'steps[2].args.connection',
      })
    })

    // A resource type this build has never heard of keeps its raw string rather
    // than silently becoming some other kind's icon.
    it('leaves an unknown resource type unnamed', async () => {
      const g = await call({
        nodes: [],
        edges: [],
        missing: [{ sourceID: '1', kind: 'corteza::system:whatsit', reason: 'step-argument' }],
      })

      expect(g.missing[0]).toMatchObject({ kind: null, resourceType: 'corteza::system:whatsit' })
    })

    // Go marshals an empty slice as null, and both fields are omitempty on top
    // of that — the graph component iterates them either way.
    it('normalises absent and null problem lists to arrays', async () => {
      expect(await call({ nodes: [], edges: [], missing: null })).toMatchObject({
        missing: [],
        warnings: [],
      })
    })

    it('renames edge endpoints for the chart component', async () => {
      const g = await call({
        nodes: [],
        edges: [{ sourceID: '1', targetID: '2', reason: 'agent-module' }],
      })

      expect(g.edges).toEqual([{ source: '1', target: '2', reason: 'agent-module' }])
    })
  })

  // The governance cycle is session-local in-memory state (see the store's own
  // header comment), so it tests directly — no API in the way. What it is worth
  // testing is the rules themselves: every step, the well-known 'publish' one
  // included, runs draft -> submitted -> approved | changes-requested, and a
  // review does not survive a change to what it reviewed.
  describe('governance', () => {
    const P = '42'

    describe('transitionStep()', () => {
      // The regression: Approve used to be a direct action from ANY status, so
      // a step could be approved before anyone asked for approval, and an
      // already-approved step could be approved again to no effect. The button
      // that fired it is gated on the same rule now (Wizard.vue canApproveStep).
      it('refuses to approve a step that was never submitted', async () => {
        const store = useProjectsStore()

        await expect(store.transitionStep(P, 'data-model', 'approve')).rejects.toThrow(
          /Cannot approve "data-model" from status "draft"/,
        )
        expect(store.governanceStatus(P, 'data-model')).toBe('draft')
      })

      it('refuses to approve the same step twice', async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'data-model', 'submit')
        await store.transitionStep(P, 'data-model', 'approve')

        await expect(store.transitionStep(P, 'data-model', 'approve')).rejects.toThrow(
          /from status "approved"/,
        )
      })

      it('takes a step through submit -> approve', async () => {
        const store = useProjectsStore()

        await store.transitionStep(P, 'data-model', 'submit')
        expect(store.governanceStatus(P, 'data-model')).toBe('submitted')

        await store.transitionStep(P, 'data-model', 'approve')
        expect(store.governanceStatus(P, 'data-model')).toBe('approved')
      })

      it('refuses to submit a step already waiting for approval', async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'roles', 'submit')

        await expect(store.transitionStep(P, 'roles', 'submit')).rejects.toThrow(
          /Cannot submit "roles" from status "submitted"/,
        )
      })

      it('lets a sent-back step be submitted again', async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'roles', 'request-changes', 'name them properly')
        expect(store.governanceStatus(P, 'roles')).toBe('changes-requested')

        await store.transitionStep(P, 'roles', 'submit')
        expect(store.governanceStatus(P, 'roles')).toBe('submitted')
        expect(store.governanceNote(P, 'roles')).toBe('')
      })

      // A reviewer must be able to flag a step nobody has submitted, and to
      // reopen one already approved — the one action with no status guard.
      it('requests changes from any status, keeping the note', async () => {
        const store = useProjectsStore()

        await store.transitionStep(P, 'pages', 'request-changes', 'from draft')
        expect(store.governanceNote(P, 'pages')).toBe('from draft')

        await store.transitionStep(P, 'pages', 'submit')
        await store.transitionStep(P, 'pages', 'approve')
        await store.transitionStep(P, 'pages', 'request-changes', 'reopened')
        expect(store.governanceStatus(P, 'pages')).toBe('changes-requested')
        expect(store.governanceNote(P, 'pages')).toBe('reopened')
      })

      it('sends a submitted publish back when any step is flagged', async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'publish', 'submit')

        await store.transitionStep(P, 'agents', 'request-changes', 'no model set')

        expect(store.governanceStatus(P, 'publish')).toBe('changes-requested')
        expect(store.governanceNote(P, 'publish')).toContain('no model set')
      })

      it('refuses to approve the revision while a step is flagged', async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'agents', 'request-changes', 'no model set')
        await store.transitionStep(P, 'publish', 'submit')

        await expect(store.transitionStep(P, 'publish', 'approve')).rejects.toThrow(
          /one or more steps still have changes requested/,
        )
      })

      // The old backend cleared every flagged step when the publish was
      // resubmitted. With per-step submit that would wipe a reviewer's note
      // without the step being fixed, and walk the gate above around itself.
      it('leaves flagged steps alone when the revision is resubmitted', async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'agents', 'request-changes', 'no model set')

        await store.transitionStep(P, 'publish', 'submit')

        expect(store.governanceStatus(P, 'agents')).toBe('changes-requested')
        expect(store.governanceNote(P, 'agents')).toBe('no model set')
      })
    })

    describe('review invalidation', () => {
      it('drops an approved step back to draft when its contents change', async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'summary', 'submit')
        await store.transitionStep(P, 'summary', 'approve')

        await store.saveStepForm(P, 'summary', { name: 'Renamed' })

        expect(store.governanceStatus(P, 'summary')).toBe('draft')
      })

      it("drops the revision's own approval when any step changes", async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'publish', 'submit')
        await store.transitionStep(P, 'publish', 'approve')

        // Awaited: scenarios are persisted now, so the review is retired only
        // once the write lands — an approval must not be dropped for a save
        // that failed.
        await store.createFriaScenario(P, { id: 's1', title: 'Denied a loan' })

        expect(store.governanceStatus(P, 'publish')).toBe('draft')
      })

      it('keeps a flagged step flagged, so the reviewer note survives the fix', async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'summary', 'request-changes', 'state the purpose')

        await store.saveStepForm(P, 'summary', { purpose: 'Loan scoring' })

        expect(store.governanceStatus(P, 'summary')).toBe('changes-requested')
        expect(store.governanceNote(P, 'summary')).toBe('state the purpose')
      })

      it('leaves untouched steps alone', async () => {
        const store = useProjectsStore()
        await store.transitionStep(P, 'roles', 'submit')
        await store.transitionStep(P, 'roles', 'approve')

        await store.saveStepForm(P, 'summary', { name: 'Renamed' })

        expect(store.governanceStatus(P, 'roles')).toBe('approved')
      })
    })
  })
})
