package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"

	composeTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The publish approval cycle: draft -> submitted -> approved | rejected.
//
// Server state, not the browser's: it survives a reload, a second user sees a
// submitted request, and the submitter cannot approve their own work. On a
// compliance product the approval is the point -- it is the record that a human
// other than the author agreed to put this revision in front of real users.
//
// Two capabilities, computed by ProjectMemberRole.Capabilities():
// CanRequestApproval submits, CanGrantApproval decides. They are read off the caller's membership of the
// chain ROOT (members are stored there -- see onAddMember), so a role granted
// on a project applies to every revision of it.

// RequestApproval submits a draft revision for review.
//
// Allowed from any state but "already waiting on somebody else": an approved
// revision that has since been edited has to come back through here (that is
// what the fingerprint mismatch in Publish tells its owner to do), and a
// rejected one obviously does. The one refusal is re-submitting a request that
// is already open -- unless the caller is the person who opened it, who is
// simply amending their own note rather than jumping a queue.
func (svc *project) RequestApproval(ctx context.Context, projectID uint64, note string) (p *types.Project, err error) {
	var aProps = &projectActionProps{project: &types.Project{ID: projectID}}

	err = func() error {
		if p, err = svc.loadApprovable(ctx, projectID, aProps); err != nil {
			return err
		}

		caps, err := svc.callerCapabilities(ctx, p)
		if err != nil {
			return err
		}
		if !caps.CanRequestApproval {
			return ProjectErrNotAllowedToRequestApproval()
		}

		me := a.GetIdentityFromContext(ctx).Identity()
		if p.ApprovalStatus == types.ProjectApprovalStatusSubmitted && p.ApprovalSubmittedBy != me {
			return ProjectErrApprovalAlreadyRequested()
		}

		p.ApprovalStatus = types.ProjectApprovalStatusSubmitted
		p.ApprovalNote = note
		p.ApprovalSubmittedBy = me
		p.ApprovalSubmittedAt = now()
		// The previous cycle's decision goes with the request that answered it.
		// Leaving the fingerprint behind would be the dangerous half: Publish
		// compares against it, so a stale one is a live approval nobody granted.
		p.ApprovalPlan = ""
		p.ApprovalDecidedBy = 0
		p.ApprovalDecidedAt = nil

		return store.UpdateProject(ctx, svc.store, p)
	}()

	return p, svc.recordAction(ctx, aProps, ProjectActionRequestApproval, err)
}

// GrantApproval approves a submitted revision for publishing, binding the
// approval to the plan it was granted against.
func (svc *project) GrantApproval(ctx context.Context, projectID uint64, note string) (p *types.Project, err error) {
	var aProps = &projectActionProps{project: &types.Project{ID: projectID}}

	err = func() error {
		if p, err = svc.decidable(ctx, projectID, aProps); err != nil {
			return err
		}

		// The fingerprint is taken HERE, not at publish time, because it is the
		// statement of what was reviewed. Publish recomputes it and compares.
		fp, err := svc.approvalFingerprint(ctx, p)
		if err != nil {
			return err
		}

		p.ApprovalStatus = types.ProjectApprovalStatusApproved
		p.ApprovalPlan = fp
		p.ApprovalNote = note
		p.ApprovalDecidedBy = a.GetIdentityFromContext(ctx).Identity()
		p.ApprovalDecidedAt = now()

		return store.UpdateProject(ctx, svc.store, p)
	}()

	return p, svc.recordAction(ctx, aProps, ProjectActionGrantApproval, err)
}

// RejectApproval sends a submitted revision back instead of approving it.
func (svc *project) RejectApproval(ctx context.Context, projectID uint64, note string) (p *types.Project, err error) {
	var aProps = &projectActionProps{project: &types.Project{ID: projectID}}

	err = func() error {
		if p, err = svc.decidable(ctx, projectID, aProps); err != nil {
			return err
		}

		p.ApprovalStatus = types.ProjectApprovalStatusRejected
		p.ApprovalNote = note
		p.ApprovalDecidedBy = a.GetIdentityFromContext(ctx).Identity()
		p.ApprovalDecidedAt = now()
		// No fingerprint on a rejection: there is nothing being held valid.
		p.ApprovalPlan = ""

		return store.UpdateProject(ctx, svc.store, p)
	}()

	return p, svc.recordAction(ctx, aProps, ProjectActionRejectApproval, err)
}

// loadApprovable loads the revision the approval cycle is about and refuses the
// states in which the question does not arise.
func (svc *project) loadApprovable(ctx context.Context, projectID uint64, aProps *projectActionProps) (*types.Project, error) {
	p, err := loadProject(ctx, svc.store, projectID)
	if err != nil {
		return nil, err
	}

	aProps.setProject(p)

	// Reading the project is the floor for taking part in its review at all --
	// the capability check below is about which part, not about whether the
	// caller may see the thing they are deciding on.
	if !svc.ac.CanReadProject(ctx, p) {
		return nil, ProjectErrNotAllowedToRead()
	}

	// loadProject returns deleted rows by design (see project.cue's lookup).
	if p.DeletedAt != nil {
		return nil, ProjectErrNotFound()
	}

	// Only a draft is a candidate for publishing, so only a draft has anything
	// to approve. Approving a live revision would leave an approval standing
	// against something already shipped.
	if p.Status != types.ProjectStatusDraft {
		return nil, ProjectErrInvalidStatus()
	}

	return p, nil
}

// decidable is the shared front half of granting and rejecting: the same
// capability, the same open-request requirement, and the same two-person rule.
func (svc *project) decidable(ctx context.Context, projectID uint64, aProps *projectActionProps) (*types.Project, error) {
	p, err := svc.loadApprovable(ctx, projectID, aProps)
	if err != nil {
		return nil, err
	}

	caps, err := svc.callerCapabilities(ctx, p)
	if err != nil {
		return nil, err
	}
	if !caps.CanGrantApproval {
		return nil, ProjectErrNotAllowedToGrantApproval()
	}

	// A decision always answers a request: without this, an Approve button
	// could act on a revision nobody put up for review, and the audit trail
	// would carry an approval with no request behind it.
	if p.ApprovalStatus != types.ProjectApprovalStatusSubmitted {
		return nil, ProjectErrApprovalNotRequested()
	}

	// Two people minimum for anything going live. This is the whole reason the
	// submitter is recorded on the row: the person who wrote the change cannot
	// also be the person who signs it off, however many capabilities they hold.
	// Rejecting is refused for the same reason -- a decision is a decision, and
	// a submitter who wants their request back simply re-submits it, which
	// resets the cycle without ever producing a review nobody performed.
	if p.ApprovalSubmittedBy == a.GetIdentityFromContext(ctx).Identity() {
		return nil, ProjectErrCannotApproveOwnRequest()
	}

	return p, nil
}

// callerCapabilities resolves what the caller may do in this project's review
// cycle, from their membership of the chain root.
//
// A non-member gets the zero value: no request, no grant. That is deliberate
// and it is not the same question as RBAC -- project.publish says you may ship
// a revision, membership says you are one of the people accountable for it.
func (svc *project) callerCapabilities(ctx context.Context, p *types.Project) (types.ProjectCapabilities, error) {
	userID := a.GetIdentityFromContext(ctx).Identity()
	if userID == 0 {
		return types.ProjectCapabilities{}, nil
	}

	m, err := store.LookupProjectMemberByProjectIDUserID(ctx, svc.store, p.RootProjectID(), userID)
	if err != nil {
		// Not a member of this project. Any other failure is a real one.
		if errors.IsNotFound(err) {
			return types.ProjectCapabilities{}, nil
		}
		return types.ProjectCapabilities{}, err
	}

	return m.RolePreset.Capabilities(), nil
}

// approvalFingerprint is the shape a publish approval is granted against.
//
// The invalidation rule is "approving states that THIS revision, as it stands,
// is fit to ship" -- so the approval has to die when the revision changes, and
// only then. That "only then" is the hard half, because the deployment plan is
// computed from LIVE data on both sides: the parent is the running project, and
// someone adding a record to it between approval and publish must not retire an
// approval that had nothing to do with them.
//
// So the fingerprint covers the DRAFT's contribution and nothing else:
//
//   - what the revision contains: every resource the project graph enumerates
//     (kind + handle, the same identity the diff matches on), plus every module
//     with its fields and their types. This is what makes the fingerprint work
//     on a FIRST publish too, where there is no parent and the plan is empty by
//     definition -- an empty plan says nothing about risk (which is why a first
//     publish needs approval at all), so something else has to carry it.
//   - what publishing would change: each plan change's op/kind/module/name/
//     detail/path/risk, and the suggested record mapping.
//
// Deliberately EXCLUDED: ProjectChange.Records. It is moduleRecordCount over
// the PARENT namespace, so it moves every time anyone uses the live app.
// Including it would have made the approval expire on its own, with an error
// telling the approver the revision had changed when nobody had touched it.
// Plan.Risk is excluded for the same reason it is redundant -- it is derived
// from the per-change risks, which are all here.
func (svc *project) approvalFingerprint(ctx context.Context, draft *types.Project) (string, error) {
	lines := make([]string, 0, 64)

	// -- what the revision contains -----------------------------------------
	_, src, err := ProjectGraph(svc.store).fetch(ctx, draft.ID)
	if err != nil {
		return "", err
	}
	for _, s := range src {
		// Modules are covered field-by-field below; users are membership rather
		// than something a publish deploys; knowledge bases are loaded unscoped
		// so they are identical on both sides by construction. Same exclusions
		// the diff makes, for the same reasons -- see resourceDiffSkip.
		if resourceDiffSkip[s.Kind] {
			continue
		}
		lines = append(lines, "res\x1f"+diffKey(s)+"\x1f"+firstNonEmpty(s.Name, s.Handle))
	}

	mods, _, err := store.SearchComposeModules(ctx, svc.store, composeTypes.ModuleFilter{
		NamespaceID: draft.Config.NamespaceID,
	})
	if err != nil {
		return "", err
	}
	for _, m := range mods {
		lines = append(lines, "mod\x1f"+m.Handle+"\x1f"+m.Name)
	}
	if len(mods) > 0 {
		ff, _, err := store.SearchComposeModuleFields(ctx, svc.store, composeTypes.ModuleFieldFilter{
			ModuleID: mods.IDs(),
		})
		if err != nil {
			return "", err
		}
		byID := make(map[uint64]string, len(mods))
		for _, m := range mods {
			byID[m.ID] = m.Handle
		}
		for _, f := range ff {
			lines = append(lines, "fld\x1f"+byID[f.ModuleID]+"\x1f"+f.Name+"\x1f"+f.Kind)
		}
	}

	// -- what publishing would change ---------------------------------------
	if draft.ParentRevisionID != 0 {
		parent, err := loadProject(ctx, svc.store, draft.ParentRevisionID)
		if err != nil {
			return "", err
		}

		plan, err := svc.computeDeploymentPlan(ctx, parent, draft)
		if err != nil {
			return "", err
		}

		for _, c := range plan.Changes {
			lines = append(lines, strings.Join([]string{
				"chg", string(c.Op), c.Kind, c.Module, c.Name, c.Detail, c.Path, string(c.Risk),
			}, "\x1f"))
		}
		for _, m := range plan.SuggestedMappings {
			for _, f := range m.Fields {
				lines = append(lines, strings.Join([]string{
					"map", m.Module, f.SourceField, f.TargetField, f.Op, f.Value, f.OnError,
				}, "\x1f"))
			}
			if len(m.Fields) == 0 {
				lines = append(lines, "map\x1f"+m.Module)
			}
		}
	}

	// Every source above is a store search with no promised order, and a hash
	// over an unordered list is a coin flip: publishing would refuse an
	// untouched revision roughly whenever the rows came back differently.
	sort.Strings(lines)

	h := sha256.New()
	// Length-prefixed so no rearrangement of the same bytes across lines can
	// collide with a different set of lines.
	for _, l := range lines {
		h.Write([]byte(strconv.Itoa(len(l))))
		h.Write([]byte{'\x1e'})
		h.Write([]byte(l))
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
