package dal

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/errors"
)

// Generic errors

func errModelNotFound(modelID uint64) error {
	return errors.NotFound("model %d does not exist", modelID)
}

func errConnectionNotFound(connectionID uint64) error {
	return errors.NotFound("connection %d does not exist", connectionID)
}

// Sensitivity level errors
// - remove
func errSensitivityLevelRemoveNotFound(sensitivityLevelID uint64) error {
	return fmt.Errorf("cannot remove sensitivity level %d: sensitivity level does not exist", sensitivityLevelID)
}

// Connection errors
// - create
func errConnectionCreateMissingSensitivityLevel(connectionID, sensitivityLevelID uint64) error {
	return fmt.Errorf("cannot create connection %d: sensitivity level does not exist %d", connectionID, sensitivityLevelID)
}
func errConnectionCreateConnectionFailed(connectionID uint64, err error) error {
	return fmt.Errorf("cannot create connection %d: connection failed: %v", connectionID, err)
}
func errConnectionDeleteNotFound(connectionID uint64) error {
	return fmt.Errorf("cannot delete connection %d: connection does not exist", connectionID)
}
func errConnectionDeleteCloserFailed(connectionID uint64, err error) error {
	return fmt.Errorf("cannot delete connection %d: connection's driver failed to close: %v", connectionID, err)
}

// - update
func errConnectionUpdateNotFound(connectionID uint64) error {
	return fmt.Errorf("cannot update connection %d: connection does not exist", connectionID)
}
func errConnectionUpdateMissingSensitivityLevel(connectionID, sensitivityLevelID uint64) error {
	return fmt.Errorf("cannot update connection %d: sensitivity level %d does not exist", connectionID, sensitivityLevelID)
}

// Model errors
// - create
func errModelCreateProblematicConnection(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot create model %d on connection %d: connection has issues", modelID, connectionID)
}
func errModelCreateMissingConnection(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot create model %d on connection %d: connection does not exist", modelID, connectionID)
}
func errModelCreateMissingSensitivityLevel(connectionID, modelID, sensitivityLevelID uint64) error {
	return fmt.Errorf("cannot create model %d on connection %d: sensitivity level %d does not exist", modelID, connectionID, sensitivityLevelID)
}
func errModelCreateGreaterSensitivityLevel(connectionID, modelID, modelSensitivityLevelID, connSensitivityLevelID uint64) error {
	return fmt.Errorf("cannot create model %d on connection %d: sensitivity level %d exceeds connection supported sensitivity level %d", modelID, connectionID, modelSensitivityLevelID, connSensitivityLevelID)
}
func errModelCreateMissingAttributeSensitivityLevel(connectionID, modelID, sensitivityLevelID uint64) error {
	return fmt.Errorf("cannot create model %d on connection %d: attribute sensitivity level %d does not exist", modelID, connectionID, sensitivityLevelID)
}
func errModelCreateGreaterAttributeSensitivityLevel(connectionID, modelID, attrSensitivityLevelID, modelSensitivityLevelID uint64) error {
	return fmt.Errorf("cannot create model %d on connection %d: attribute sensitivity level %d exceeds model supported sensitivity level %d", modelID, connectionID, attrSensitivityLevelID, modelSensitivityLevelID)
}
func errModelCreateConnectionModelUnsupported(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot create model %d on connection %d: model already exists for connection but is not compatible with provided definition", modelID, connectionID)
}
func errModelCreateInvalidIdent(connectionID, modelID uint64, ident string) error {
	return fmt.Errorf("cannot create model %d on connection %d: malformed model ident %s", modelID, connectionID, ident)
}

// - update
func errModelUpdateProblematicConnection(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot update model %d on connection %d: connection has issues", modelID, connectionID)
}
func errModelUpdateMissingConnection(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot update model %d on connection %d: connection does not exist", modelID, connectionID)
}
func errModelUpdateConnectionModelUnsupported(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot update model %d on connection %d: model already exists for connection but is not compatible with provided definition", modelID, connectionID)
}
func errModelUpdateMissingOldModel(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot update model %d on connection %d: model does not exist", modelID, connectionID)
}
func errModelUpdateDuplicate(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot update model %d on connection %d: model already exists", modelID, connectionID)
}
func errModelUpdateConnectionMissmatch(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot update model %d on connection %d: cannot change model connection", modelID, connectionID)
}
func errModelUpdateMissingSensitivityLevel(connectionID, modelID, sensitivityLevelID uint64) error {
	return fmt.Errorf("cannot update model %d on connection %d: sensitivity level %d does not exist", modelID, connectionID, sensitivityLevelID)
}
func errModelUpdateGreaterSensitivityLevel(connectionID, modelID, modelSensitivityLevelID, connSensitivityLevelID uint64) error {
	return fmt.Errorf("cannot update model %d on connection %d: sensitivity level %d exceeds connection supported sensitivity level %d", modelID, connectionID, modelSensitivityLevelID, connSensitivityLevelID)
}

// - alterations
func errModelRequiresAlteration(connectionID, modelID, batchID uint64) error {
	return fmt.Errorf("model %d on connection %d requires schema alterations: alteration batchID %d", modelID, connectionID, batchID)
}

// Attribute errors
// - Update
func errAttributeUpdateProblematicConnection(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot update attribute for model %d on connection %d: connection has issues", modelID, connectionID)
}
func errAttributeUpdateMissingModel(connectionID, modelID uint64) error {
	return fmt.Errorf("cannot update attribute for model %d on connection %d: model does not exist", modelID, connectionID)
}
func errAttributeUpdateMissingSensitivityLevel(connectionID, modelID, sensitivityLevelID uint64) error {
	return fmt.Errorf("cannot update attribute for model %d on connection %d: sensitivity level %d does not exist", modelID, connectionID, sensitivityLevelID)
}
func errAttributeUpdateGreaterSensitivityLevel(connectionID, modelID, attrSensitivityLevelID, modelSensitivityLevelID uint64) error {
	return fmt.Errorf("cannot update attribute for model %d on connection %d: sensitivity level %d exceeds model supported sensitivity level %d", modelID, connectionID, attrSensitivityLevelID, modelSensitivityLevelID)
}

// Record errors

func errRecordOpProblematicConnection(connectionID uint64) error {
	return fmt.Errorf("cannot perform record operation: connection %d has issues", connectionID)
}
func errRecordOpProblematicModel(modelID uint64) error {
	return fmt.Errorf("cannot perform record operation: model %d has issues", modelID)
}

// func errModelHigherSensitivity(model, connection string) error {
// 	return errors.New(
// 		errors.KindSensitiveData,

// 		"model sensitivity surpasses connection sensitivity",

// 		errors.Meta("type", "invalid sensitivity"),

// 		// Translation namespace & key
// 		errors.Meta(locale.ErrorMetaNamespace{}, "internal"),
// 		errors.Meta(locale.ErrorMetaKey{}, "dal.sensitivity.model-exceeds-connection"),
// 		errors.Meta("model", model),
// 		errors.Meta("connection", connection),

// 		errors.StackSkip(1),
// 		errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
// 	)
// }

// func errAttributeHigherSensitivity(model, attribute string) error {
// 	return errors.New(
// 		errors.KindSensitiveData,

// 		"attribute sensitivity surpasses model sensitivity",

// 		errors.Meta("type", "invalid sensitivity"),

// 		// Translation namespace & key
// 		errors.Meta(locale.ErrorMetaNamespace{}, "internal"),
// 		errors.Meta(locale.ErrorMetaKey{}, "dal.sensitivity.attribute-exceeds-model"),
// 		errors.Meta("model", model),
// 		errors.Meta("attribute", attribute),

// 		errors.StackSkip(1),
// 		errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
// 	)
// }
