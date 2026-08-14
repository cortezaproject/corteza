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

  // The publish approval is the one part of the review cycle that is REAL
  // server state: it lives on the project row, and the backend refuses to
  // publish anything that is not approved, refuses a decision from whoever
  // submitted the request, and refuses an approval granted against a different
  // version of the revision.
  describe('publish approval', () => {
    const seed = async (api, project) => {
      api.projectList = vi.fn().mockResolvedValue({ set: [project] })
      const store = useProjectsStore()
      await store.load()
      return store
    }

    it('reads the approval standing off the project row', async () => {
      const store = await seed(api, { projectID: '7', approvalStatus: 'submitted' })

      expect(store.publishApprovalStatus('7')).toBe('submitted')
    })

    // The server says 'rejected'; every status chip and hint in this section
    // says 'changes-requested' for the same thing.
    it("speaks the wizard's word for a rejection", async () => {
      const store = await seed(api, {
        projectID: '7',
        approvalStatus: 'rejected',
        approvalNote: 'the retention rule is missing',
      })

      expect(store.publishApprovalStatus('7')).toBe('changes-requested')
      expect(store.publishApprovalNote('7')).toBe('the retention rule is missing')
    })

    // A row written before the column existed carries no value at all. Reading
    // that as anything but 'draft' would show an unreviewed revision as though
    // somebody had looked at it.
    it('reads a missing standing as unreviewed', async () => {
      const store = await seed(api, { projectID: '7' })

      expect(store.publishApprovalStatus('7')).toBe('draft')
    })

    it('sends each decision to its own endpoint and absorbs the answer', async () => {
      const store = await seed(api, { projectID: '7', approvalStatus: 'draft' })

      api.projectRequestApproval = vi
        .fn()
        .mockResolvedValue({ projectID: '7', approvalStatus: 'submitted' })
      await store.requestPublishApproval('7', 'ready')
      expect(api.projectRequestApproval).toHaveBeenCalledWith({ projectID: '7', note: 'ready' })
      expect(store.publishApprovalStatus('7')).toBe('submitted')

      api.projectGrantApproval = vi
        .fn()
        .mockResolvedValue({ projectID: '7', approvalStatus: 'approved' })
      await store.grantPublishApproval('7', 'reviewed')
      expect(api.projectGrantApproval).toHaveBeenCalledWith({ projectID: '7', note: 'reviewed' })
      expect(store.publishApprovalStatus('7')).toBe('approved')

      api.projectRejectApproval = vi
        .fn()
        .mockResolvedValue({ projectID: '7', approvalStatus: 'rejected' })
      await store.rejectPublishApproval('7', 'not yet')
      expect(api.projectRejectApproval).toHaveBeenCalledWith({ projectID: '7', note: 'not yet' })
      expect(store.publishApprovalStatus('7')).toBe('changes-requested')
    })
  })

  describe('deploymentPlan()', () => {
    // Go marshals a nil slice as `null`, not `[]`, so an empty plan arrives as
    // {changes: null}. Spreading that over array defaults overwrites them with
    // null and the Publish tab dies on `changes.filter(...)` — on every first
    // revision, which is exactly the case with an empty plan.
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

  // The endpoint's account of what it could NOT resolve — a chatbot pointing at
  // another revision's agent, a TAQ binding a branch copy dropped — has to reach
  // the canvas. Destructuring only `nodes` and `edges` renders a broken project
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
      // Approve is reachable only from 'submitted': nothing is approved before
      // anyone asked, and an approved step cannot be approved again. The button
      // is gated on the same rule (Wizard.vue canApproveStep).
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

      // The revision's own approval is no longer one of these steps — it lives
      // on the project row and the server enforces it — so a flagged step no
      // longer reaches in and changes it. What survives is the report the
      // Publish tab uses to refuse offering an approval over an open flag.
      it('reports a flagged step so the publish screen can refuse to approve over it', async () => {
        const store = useProjectsStore()
        expect(store.hasFlaggedSteps(P)).toBe(false)

        await store.transitionStep(P, 'agents', 'request-changes', 'no model set')

        expect(store.hasFlaggedSteps(P)).toBe(true)
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

      // The revision's OWN approval is deliberately NOT retired from here any
      // more: it is server state, and the server retires it by recomputing the
      // deployment plan's fingerprint when a publish is attempted. Nothing on
      // this side has to remember to call an invalidation hook, which is the
      // whole point — an edit path that forgot to would have kept a stale
      // approval alive.
      it("leaves the revision's own approval to the server", async () => {
        const store = useProjectsStore()

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
