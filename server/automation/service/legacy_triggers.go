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
			},
		},
		types.ConstructTrigger{
			ResourceType: "system",
			EventType:    "onInterval",
			Groups:       []string{"System"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Interval",
				Description: "Triggered on interval",
			},
		},
		types.ConstructTrigger{
			ResourceType: "system",
			EventType:    "onTimestamp",
			Groups:       []string{"System"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Timestamp",
				Description: "Triggered at timestamp",
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "beforeCreate",
			Groups:       []string{"Compose Record", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before record create",
				Description: "Triggered before record is created",
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "beforeUpdate",
			Groups:       []string{"Compose Record", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before record update",
				Description: "Triggered before record is updated",
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "beforeDelete",
			Groups:       []string{"Compose Record", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before record delete",
				Description: "Triggered before record is deleted",
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "beforeUndelete",
			Groups:       []string{"Compose Record", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before record undelete",
				Description: "Triggered before record is undeleted",
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterUndelete",
			Groups:       []string{"Compose Record", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After record undelete",
				Description: "Triggered after record is undeleted",
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "beforeCreate",
			Groups:       []string{"System User", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user create",
				Description: "Triggered before user is created",
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "afterCreate",
			Groups:       []string{"System User", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user create",
				Description: "Triggered after user is created",
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "beforeUpdate",
			Groups:       []string{"System User", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user update",
				Description: "Triggered before user is updated",
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "afterUpdate",
			Groups:       []string{"System User", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user update",
				Description: "Triggered after user is updated",
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "beforeDelete",
			Groups:       []string{"System User", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user delete",
				Description: "Triggered before user is deleted",
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "afterDelete",
			Groups:       []string{"System User", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user delete",
				Description: "Triggered after user is deleted",
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "beforeSuspend",
			Groups:       []string{"System User", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user suspend",
				Description: "Triggered before user is suspended",
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "afterSuspend",
			Groups:       []string{"System User", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user suspend",
				Description: "Triggered after user is suspended",
			},
		},
	)
}
