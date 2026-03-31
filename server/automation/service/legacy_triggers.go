package service

import "github.com/cortezaproject/corteza/server/automation/types"

func init() {
	ConstructLibrary().AddTriggers(
		types.ConstructTrigger{
			ResourceType: "system",
			EventType:    "onManual",
			Groups:       []string{"System"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Manual",
				Description: "Triggered manually",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "play"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system",
			EventType:    "onInterval",
			Groups:       []string{"System"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Interval",
				Description: "Triggered on interval",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "refresh"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system",
			EventType:    "onTimestamp",
			Groups:       []string{"System"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Timestamp",
				Description: "Triggered at timestamp",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "calendar"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "beforeCreate",
			Groups:       []string{"Records"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before Record Create",
				Description: "Triggered before record is created",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "database"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "beforeUpdate",
			Groups:       []string{"Records"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before Record Update",
				Description: "Triggered before record is updated",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "database"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "beforeDelete",
			Groups:       []string{"Records"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before Record Delete",
				Description: "Triggered before record is deleted",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "database"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "beforeUndelete",
			Groups:       []string{"Records"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before Record Undelete",
				Description: "Triggered before record is undeleted",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "database"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterUndelete",
			Groups:       []string{"Records"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After Record Undelete",
				Description: "Triggered after record is undeleted",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "database"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "beforeCreate",
			Groups:       []string{"Users"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before User Create",
				Description: "Triggered before user is created",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "users"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "afterCreate",
			Groups:       []string{"Users"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After User Create",
				Description: "Triggered after user is created",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "users"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "beforeUpdate",
			Groups:       []string{"Users"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before User Update",
				Description: "Triggered before user is updated",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "users"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "afterUpdate",
			Groups:       []string{"Users"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After User Update",
				Description: "Triggered after user is updated",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "users"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "beforeDelete",
			Groups:       []string{"Users"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before User Delete",
				Description: "Triggered before user is deleted",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "users"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "afterDelete",
			Groups:       []string{"Users"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After User Delete",
				Description: "Triggered after user is deleted",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "users"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "beforeSuspend",
			Groups:       []string{"Users"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before User Suspend",
				Description: "Triggered before user is suspended",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "users"},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "afterSuspend",
			Groups:       []string{"Users"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After User Suspend",
				Description: "Triggered after user is suspended",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "users"},
			},
		},
	)
}
