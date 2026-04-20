package automation

import (
	"context"
	"fmt"

	composeTypes "github.com/crusttech/human/server/compose/types"
	systemTypes "github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
)

type (
	notificationHandler struct {
		reg    notificationHandlerRegistry
		svc    notificationService
		uSvc   notificationUserService
		nSvc   notificationNamespaceService
		mSvc   notificationModuleService
		logger *zap.Logger
	}

	notificationService interface {
		Create(context.Context, *systemTypes.Notification) (*systemTypes.Notification, error)
	}

	// Interface for user service that we need for recipient lookups
	notificationUserService interface {
		FindByID(ctx context.Context, userID uint64) (*systemTypes.User, error)
		FindByHandle(ctx context.Context, handle string) (*systemTypes.User, error)
		FindByEmail(ctx context.Context, email string) (*systemTypes.User, error)
	}

	notificationNamespaceService interface {
		FindByHandle(ctx context.Context, handle string) (*composeTypes.Namespace, error)
	}

	notificationModuleService interface {
		FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*composeTypes.Module, error)
	}
)

// NotificationHandler initializes the notification handler
func NotificationHandler(reg notificationHandlerRegistry, ntfSvc notificationService, uSvc notificationUserService, nSvc notificationNamespaceService, mSvc notificationModuleService, log *zap.Logger) *notificationHandler {
	h := &notificationHandler{
		reg:    reg,
		svc:    ntfSvc,
		uSvc:   uSvc,
		nSvc:   nSvc,
		mSvc:   mSvc,
		logger: log.Named("notification"),
	}

	h.register()
	return h
}

// sendRecord creates and sends a record notification
func (h notificationHandler) sendRecord(ctx context.Context, args *notificationSendRecordArgs) error {
	// Get the recipient user ID from the input parameter
	var recipientID uint64

	// Check if we have a direct user object
	if args.recipientID > 0 {
		// Direct ID passed
		recipientID = args.recipientID
	} else {
		// Need to look up the user
		var user *systemTypes.User
		var err error

		if args.recipientHandle != "" {
			user, err = h.uSvc.FindByHandle(ctx, args.recipientHandle)
		} else if args.recipientEmail != "" {
			user, err = h.uSvc.FindByEmail(ctx, args.recipientEmail)
		} else {
			return fmt.Errorf("invalid recipient: unable to determine user ID")
		}

		if err != nil {
			return err
		}

		if user == nil {
			return fmt.Errorf("recipient not found")
		}

		recipientID = user.ID
	}

	if args.namespaceHandle != "" {
		namespace, err := h.nSvc.FindByHandle(ctx, args.namespaceHandle)
		if err != nil {
			return err
		}

		args.namespaceID = namespace.ID
	}

	if args.moduleHandle != "" {
		module, err := h.mSvc.FindByHandle(ctx, args.namespaceID, args.moduleHandle)
		if err != nil {
			return err
		}

		args.moduleID = module.ID
	}

	// Create notification config for a record notification
	config := systemTypes.NotificationConfig{
		Record: systemTypes.RecordNotificationConfig{
			Title:       args.Title,
			Description: args.Description,
			ModuleID:    args.moduleID,
			NamespaceID: args.namespaceID,
			RecordID:    args.Record,
			OpenMode:    systemTypes.OpenModeType(args.OpenMode),
			Edit:        args.Edit,
		},
	}

	// Create a record notification
	ntf := &systemTypes.Notification{
		Kind:      systemTypes.NotificationKindRecord,
		Config:    config,
		Recipient: recipientID,
	}

	// Create the notification
	_, err := h.svc.Create(ctx, ntf)
	return err
}
