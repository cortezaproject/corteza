package automation

import (
	"context"
	"fmt"

	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

type (
	notificationHandler struct {
		reg    notificationHandlerRegistry
		svc    notificationService
		uSvc   notificationUserService
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
)

// NotificationHandler initializes the notification handler
func NotificationHandler(reg notificationHandlerRegistry, ntfSvc notificationService, uSvc notificationUserService, log *zap.Logger) *notificationHandler {
	h := &notificationHandler{
		reg:    reg,
		svc:    ntfSvc,
		uSvc:   uSvc,
		logger: log.Named("notification"),
	}

	h.register()
	return h
}

// send creates and sends a notification
func (h notificationHandler) send(ctx context.Context, args *notificationSendArgs) error {
	// Get the recipient user ID from the input parameter
	var recipientID uint64

	// Check if we have a direct user object
	if args.recipientRes != nil {
		recipientID = args.recipientRes.ID
	} else if args.recipientID > 0 {
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

	// Create notification config directly with the expected fields for a simple notification
	config := systemTypes.NotificationConfig{
		Simple: systemTypes.SimpleNotificationConfig{
			Title:       args.Title,
			Description: args.Description,
		},
	}

	// Create a simple notification
	ntf := &systemTypes.Notification{
		Kind:      systemTypes.NotificationKindSimple,
		Config:    config,
		Recipient: recipientID,
	}

	// Create the notification
	_, err := h.svc.Create(ctx, ntf)
	return err
}

// sendRecord creates and sends a record notification
func (h notificationHandler) sendRecord(ctx context.Context, args *notificationSendRecordArgs) error {
	// Get the recipient user ID from the input parameter
	var recipientID uint64

	// Check if we have a direct user object
	if args.recipientRes != nil {
		recipientID = args.recipientRes.ID
	} else if args.recipientID > 0 {
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

	// Create notification config for a record notification
	config := systemTypes.NotificationConfig{
		Record: systemTypes.RecordNotificationConfig{
			Title:       args.Title,
			Description: args.Description,
			ModuleID:    args.ModuleID,
			NamespaceID: args.NamespaceID,
			RecordID:    args.RecordID,
			OpenMode:    args.OpenMode,
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
