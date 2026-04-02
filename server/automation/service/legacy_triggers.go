package service

import "github.com/cortezaproject/corteza/server/automation/types"

func recordConstraints() []types.ConstructTriggerConstraint {
	return []types.ConstructTriggerConstraint{
		{Name: "namespace", Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true},
		{Name: "module", Types: []string{"ID", "Handle", "ComposeModule"}, Required: true},
		{Name: "record", Types: []string{"ID", "ComposeRecord"}},
	}
}

func recordProperties() []types.ConstructTriggerProperty {
	return []types.ConstructTriggerProperty{
		{Name: "record", Type: "ComposeRecord"},
		{Name: "oldRecord", Type: "ComposeRecord"},
		{Name: "module", Type: "ComposeModule"},
		{Name: "namespace", Type: "ComposeNamespace"},
		{Name: "recordValueErrors", Type: "ComposeRecordValueErrorSet"},
		{Name: "selected", Type: ""},
	}
}

func recordSegments() []types.ConstructSegment {
	return []types.ConstructSegment{{
		Sections: []types.ConstructSection{{
			Elements: []types.SectionElement{{
				Input: types.SectionElementInput{
					Type:     "NamespaceSelector",
					Label:    "Namespace",
					Argument: "namespace",
					Required: true,
				},
			}, {
				Input: types.SectionElementInput{
					Type:     "ModuleSelector",
					Label:    "Module",
					Argument: "module",
					Required: true,
					Context: types.SectionElementInputContext{
						DependsOn: map[string]string{
							"namespaceID": "namespace",
						},
					},
				},
			}, {
				Input: types.SectionElementInput{
					Type:     "RecordSelector",
					Label:    "Record",
					Argument: "record",
					Context: types.SectionElementInputContext{
						DependsOn: map[string]string{
							"namespaceID": "namespace",
							"moduleID":    "module",
						},
					},
				},
			}},
		}},
	}}
}

func userConstraints() []types.ConstructTriggerConstraint {
	return []types.ConstructTriggerConstraint{
		// {Name: "user", Types: []string{"ID", "Handle", "SystemUser"}},
	}
}

func userProperties() []types.ConstructTriggerProperty {
	return []types.ConstructTriggerProperty{
		{Name: "user", Type: "User"},
		{Name: "oldUser", Type: "User"},
	}
}

func userSegments() []types.ConstructSegment {
	return []types.ConstructSegment{{
		// Sections: []types.ConstructSection{{
		// 	Elements: []types.SectionElement{{
		// 		Input: types.SectionElementInput{
		// 			Type:     "UserSelector",
		// 			Label:    "User",
		// 			Argument: "user",
		// 		},
		// 	}},
		// }},
	}}
}

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
			Segments: []types.ConstructSegment{{
				Sections: []types.ConstructSection{{
					Elements: []types.SectionElement{{
						Input: types.SectionElementInput{
							Type:        "Interval",
							Label:       "Interval",
							Argument:    "interval",
							Placeholder: "* * * * *",
						},
					}},
				}},
			}},
			Constraints: []types.ConstructTriggerConstraint{
				{Name: "interval", Types: []string{"String"}},
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
			Segments: []types.ConstructSegment{{
				Sections: []types.ConstructSection{{
					Elements: []types.SectionElement{{
						Input: types.SectionElementInput{
							Type:        "DateTime",
							Label:       "Timestamp (RFC3339)",
							Argument:    "timestamp",
							Placeholder: "2026-01-02T15:04:05Z",
						},
					}},
				}},
			}},
			Constraints: []types.ConstructTriggerConstraint{
				{Name: "timestamp", Types: []string{"String"}},
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
			Properties:  recordProperties(),
			Segments:    recordSegments(),
			Constraints: recordConstraints(),
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterCreate",
			Groups:       []string{"Records"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After Record Create",
				Description: "Triggered after record is created",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "database"},
			},
			Properties:  recordProperties(),
			Segments:    recordSegments(),
			Constraints: recordConstraints(),
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
			Properties:  recordProperties(),
			Segments:    recordSegments(),
			Constraints: recordConstraints(),
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterUpdate",
			Groups:       []string{"Records"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After Record Update",
				Description: "Triggered after record is updated",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "database"},
			},
			Properties:  recordProperties(),
			Segments:    recordSegments(),
			Constraints: recordConstraints(),
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
			Properties:  recordProperties(),
			Segments:    recordSegments(),
			Constraints: recordConstraints(),
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterDelete",
			Groups:       []string{"Records"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After Record Delete",
				Description: "Triggered after record is deleted",
				Icon:        &types.NgAutomationIcon{Type: "name", Value: "database"},
			},
			Properties:  recordProperties(),
			Segments:    recordSegments(),
			Constraints: recordConstraints(),
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
			Properties:  recordProperties(),
			Segments:    recordSegments(),
			Constraints: recordConstraints(),
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
			Properties:  recordProperties(),
			Segments:    recordSegments(),
			Constraints: recordConstraints(),
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
			Properties:  userProperties(),
			Segments:    userSegments(),
			Constraints: userConstraints(),
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
			Properties:  userProperties(),
			Segments:    userSegments(),
			Constraints: userConstraints(),
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
			Properties:  userProperties(),
			Segments:    userSegments(),
			Constraints: userConstraints(),
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
			Properties:  userProperties(),
			Segments:    userSegments(),
			Constraints: userConstraints(),
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
			Properties:  userProperties(),
			Segments:    userSegments(),
			Constraints: userConstraints(),
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
			Properties:  userProperties(),
			Segments:    userSegments(),
			Constraints: userConstraints(),
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
			Properties:  userProperties(),
			Segments:    userSegments(),
			Constraints: userConstraints(),
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
			Properties:  userProperties(),
			Segments:    userSegments(),
			Constraints: userConstraints(),
		},
	)
}
