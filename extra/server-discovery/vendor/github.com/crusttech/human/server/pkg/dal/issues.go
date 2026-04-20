package dal

import "go.uber.org/zap"

type (
	Issue struct {
		err   error
		Issue string `json:"issue,omitempty"`

		Kind issueKind      `json:"kind,omitempty"`
		Meta map[string]any `json:"meta,omitempty"`
	}
	issueSet []Issue

	issueHelper struct {
		// these two will be used to help clear out unneeded errors
		connections []uint64
		models      []uint64

		connectionIssues dalIssueIndex
		modelIssues      dalIssueIndex
	}

	issueKind     string
	dalIssueIndex map[uint64]issueSet
)

const (
	connectionIssue issueKind = "connection"
	modelIssue      issueKind = "model"
)

func newIssueHelper() *issueHelper {
	return &issueHelper{
		connectionIssues: make(dalIssueIndex),
		modelIssues:      make(dalIssueIndex),
	}
}

func (svc *service) SearchConnectionIssues(connectionID uint64) (out []Issue) {
	return svc.connectionIssues[connectionID]
}

func (svc *service) SearchModelIssues(resourceID uint64) (out []Issue) {
	return svc.modelIssues[resourceID]
}

func (svc *service) SearchResourceIssues(resourceType, resource string) (out []Issue) {
	var m *Model
	for _, ax := range svc.models {
		m = ax.FindByResourceIdent(resourceType, resource)
	}

	if m == nil {
		return
	}

	return svc.modelIssues[m.ResourceID]
}

func (svc *service) hasConnectionIssues(connectionID uint64) bool {
	return len(svc.SearchConnectionIssues(connectionID)) > 0
}

func (svc *service) hasModelIssues(modelID uint64) bool {
	return len(svc.SearchModelIssues(modelID)) > 0
}

func (svc *service) updateIssues(issues *issueHelper) {
	for _, connectionID := range issues.connections {
		delete(svc.connectionIssues, connectionID)
	}
	for connectionID, issues := range issues.connectionIssues {
		svc.connectionIssues[connectionID] = issues
	}

	for _, modelID := range issues.models {
		delete(svc.modelIssues, modelID)
	}
	for modelID, issues := range issues.modelIssues {
		svc.modelIssues[modelID] = issues
	}
}

func (svc *service) clearModelIssues() {
	svc.modelIssues = make(dalIssueIndex)
}

func (rd *issueHelper) addConnection(connectionID uint64) *issueHelper {
	rd.connections = append(rd.connections, connectionID)
	return rd
}

func (rd *issueHelper) addModel(modelID uint64) *issueHelper {
	rd.models = append(rd.models, modelID)
	return rd
}

func (rd *issueHelper) addConnectionIssue(connectionID uint64, i Issue) {
	i.Kind = connectionIssue
	i.Issue = i.err.Error()
	rd.connectionIssues[connectionID] = append(rd.connectionIssues[connectionID], i)
}

func (rd *issueHelper) addModelIssue(resourceID uint64, i Issue) {
	i.Kind = modelIssue
	i.Issue = i.err.Error()
	rd.modelIssues[resourceID] = append(rd.modelIssues[resourceID], i)
}

func (rd *issueHelper) hasConnectionIssues() bool {
	return len(rd.connectionIssues) > 0
}

func (rd *issueHelper) hasModelIssues() bool {
	return len(rd.modelIssues) > 0
}

func (a *issueHelper) mergeWith(b *issueHelper) {
	if b == nil {
		return
	}

	for connectionID, issues := range b.connectionIssues {
		a.connectionIssues[connectionID] = append(a.connectionIssues[connectionID], issues...)
	}

	for modelID, issues := range b.modelIssues {
		a.modelIssues[modelID] = append(a.modelIssues[modelID], issues...)
	}
}

// Op check utils
func (svc *service) canOpData(ref ModelRef) (err error) {
	if svc.hasConnectionIssues(ref.ConnectionID) {
		for _, i := range svc.connectionIssues[ref.ConnectionID] {
			svc.logger.Debug(
				"can not perform data operation due to connection issue: "+i.err.Error(),
				zap.Any("ref", ref.ResourceID),
			)
		}

		return errRecordOpProblematicConnection(ref.ConnectionID)
	}

	var mod = svc.FindModelByRef(ref)
	if mod == nil {
		return errModelNotFound(ref.ResourceID)
	}

	if svc.hasModelIssues(mod.ResourceID) {
		for _, i := range svc.modelIssues[mod.ResourceID] {
			svc.logger.Debug(
				"can not perform data operation due to model issue: "+i.err.Error(),
				zap.Any("ref", ref.ResourceID),
			)
		}

		return errRecordOpProblematicModel(mod.ResourceID)
	}

	return nil
}
