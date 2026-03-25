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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:        "String",
										Label:       "Interval (Cron)",
										Argument:    "interval",
										Placeholder: "* * * * *",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{{
				Name:  "interval",
				Types: []string{"String"},
				Meta:  types.ConstructTriggerConstraintMeta{},
			}},
			Properties: []types.ConstructTriggerProperty{},
		},
		types.ConstructTrigger{
			ResourceType: "system",
			EventType:    "onTimestamp",
			Groups:       []string{"System"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Timestamp",
				Description: "Triggered at timestamp",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:        "DateTime",
										Label:       "Timestamp (RFC3339)",
										Argument:    "timestamp",
										Placeholder: "2026-01-02T15:04:05Z",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{{
				Name:  "timestamp",
				Types: []string{"String"},
				Meta:  types.ConstructTriggerConstraintMeta{},
			}},
			Properties: []types.ConstructTriggerProperty{},
		},
		types.ConstructTrigger{
			ResourceType: "compose:namespace",
			EventType:    "beforeCreate",
			Groups:       []string{"Compose Namespace", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before namespace create",
				Description: "Triggered before namespace is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:namespace",
			EventType:    "afterCreate",
			Groups:       []string{"Compose Namespace", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After namespace create",
				Description: "Triggered after namespace is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:namespace",
			EventType:    "beforeUpdate",
			Groups:       []string{"Compose Namespace", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before namespace update",
				Description: "Triggered before namespace is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:namespace",
			EventType:    "afterUpdate",
			Groups:       []string{"Compose Namespace", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After namespace update",
				Description: "Triggered after namespace is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:namespace",
			EventType:    "beforeDelete",
			Groups:       []string{"Compose Namespace", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before namespace delete",
				Description: "Triggered before namespace is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:namespace",
			EventType:    "afterDelete",
			Groups:       []string{"Compose Namespace", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After namespace delete",
				Description: "Triggered after namespace is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page",
			EventType:    "beforeCreate",
			Groups:       []string{"Compose Page", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before page create",
				Description: "Triggered before page is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Handle",
										Argument: "page.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Name",
										Argument: "page.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "page.handle",
					Types: []string{"String"},
				},
				{
					Name:  "page.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page",
			EventType:    "afterCreate",
			Groups:       []string{"Compose Page", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After page create",
				Description: "Triggered after page is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Handle",
										Argument: "page.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Name",
										Argument: "page.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "page.handle",
					Types: []string{"String"},
				},
				{
					Name:  "page.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page",
			EventType:    "beforeUpdate",
			Groups:       []string{"Compose Page", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before page update",
				Description: "Triggered before page is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Handle",
										Argument: "page.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Name",
										Argument: "page.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "page.handle",
					Types: []string{"String"},
				},
				{
					Name:  "page.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page",
			EventType:    "afterUpdate",
			Groups:       []string{"Compose Page", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After page update",
				Description: "Triggered after page is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Handle",
										Argument: "page.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Name",
										Argument: "page.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "page.handle",
					Types: []string{"String"},
				},
				{
					Name:  "page.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page",
			EventType:    "beforeDelete",
			Groups:       []string{"Compose Page", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before page delete",
				Description: "Triggered before page is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Handle",
										Argument: "page.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Name",
										Argument: "page.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "page.handle",
					Types: []string{"String"},
				},
				{
					Name:  "page.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page",
			EventType:    "afterDelete",
			Groups:       []string{"Compose Page", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After page delete",
				Description: "Triggered after page is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Handle",
										Argument: "page.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Page Name",
										Argument: "page.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "page.handle",
					Types: []string{"String"},
				},
				{
					Name:  "page.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page-layout",
			EventType:    "beforeCreate",
			Groups:       []string{"Compose Page Layout", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before page layout create",
				Description: "Triggered before page layout is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Handle",
										Argument: "pageLayout.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Name",
										Argument: "pageLayout.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.handle",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page-layout",
			EventType:    "afterCreate",
			Groups:       []string{"Compose Page Layout", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After page layout create",
				Description: "Triggered after page layout is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Handle",
										Argument: "pageLayout.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Name",
										Argument: "pageLayout.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.handle",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page-layout",
			EventType:    "beforeUpdate",
			Groups:       []string{"Compose Page Layout", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before page layout update",
				Description: "Triggered before page layout is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Handle",
										Argument: "pageLayout.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Name",
										Argument: "pageLayout.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.handle",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page-layout",
			EventType:    "afterUpdate",
			Groups:       []string{"Compose Page Layout", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After page layout update",
				Description: "Triggered after page layout is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Handle",
										Argument: "pageLayout.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Name",
										Argument: "pageLayout.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.handle",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page-layout",
			EventType:    "beforeDelete",
			Groups:       []string{"Compose Page Layout", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before page layout delete",
				Description: "Triggered before page layout is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Handle",
										Argument: "pageLayout.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Name",
										Argument: "pageLayout.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.handle",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:page-layout",
			EventType:    "afterDelete",
			Groups:       []string{"Compose Page Layout", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After page layout delete",
				Description: "Triggered after page layout is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Handle",
										Argument: "pageLayout.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Pagelayout Name",
										Argument: "pageLayout.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.handle",
					Types: []string{"String"},
				},
				{
					Name:  "pageLayout.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:module",
			EventType:    "beforeCreate",
			Groups:       []string{"Compose Module", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before module create",
				Description: "Triggered before module is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:module",
			EventType:    "afterCreate",
			Groups:       []string{"Compose Module", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After module create",
				Description: "Triggered after module is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:module",
			EventType:    "beforeUpdate",
			Groups:       []string{"Compose Module", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before module update",
				Description: "Triggered before module is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:module",
			EventType:    "afterUpdate",
			Groups:       []string{"Compose Module", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After module update",
				Description: "Triggered after module is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:module",
			EventType:    "beforeDelete",
			Groups:       []string{"Compose Module", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before module delete",
				Description: "Triggered before module is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:module",
			EventType:    "afterDelete",
			Groups:       []string{"Compose Module", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After module delete",
				Description: "Triggered after module is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "onIteration",
			Groups:       []string{"Compose Record"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "On record iteration",
				Description: "Triggered on record iteration",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Created-at",
										Argument: "record.created-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Updated-at",
										Argument: "record.updated-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Deleted-at",
										Argument: "record.deleted-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Values Any",
										Argument: "record.values.*",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
				{
					Name:  "record.created-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.updated-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.deleted-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.values.*",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Created-at",
										Argument: "record.created-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Updated-at",
										Argument: "record.updated-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Deleted-at",
										Argument: "record.deleted-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Values Any",
										Argument: "record.values.*",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
				{
					Name:  "record.created-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.updated-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.deleted-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.values.*",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Created-at",
										Argument: "record.created-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Updated-at",
										Argument: "record.updated-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Deleted-at",
										Argument: "record.deleted-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Values Any",
										Argument: "record.values.*",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
				{
					Name:  "record.created-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.updated-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.deleted-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.values.*",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Created-at",
										Argument: "record.created-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Updated-at",
										Argument: "record.updated-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Deleted-at",
										Argument: "record.deleted-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Values Any",
										Argument: "record.values.*",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
				{
					Name:  "record.created-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.updated-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.deleted-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.values.*",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Created-at",
										Argument: "record.created-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Updated-at",
										Argument: "record.updated-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Deleted-at",
										Argument: "record.deleted-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Values Any",
										Argument: "record.values.*",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
				{
					Name:  "record.created-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.updated-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.deleted-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.values.*",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Created-at",
										Argument: "record.created-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Updated-at",
										Argument: "record.updated-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Deleted-at",
										Argument: "record.deleted-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Values Any",
										Argument: "record.values.*",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
				{
					Name:  "record.created-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.updated-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.deleted-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.values.*",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "beforeOrganize",
			Groups:       []string{"Compose Record", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before record organize",
				Description: "Triggered before record is organized",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Created-at",
										Argument: "record.created-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Updated-at",
										Argument: "record.updated-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Deleted-at",
										Argument: "record.deleted-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Values Any",
										Argument: "record.values.*",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
				{
					Name:  "record.created-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.updated-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.deleted-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.values.*",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "compose:record",
			EventType:    "afterOrganize",
			Groups:       []string{"Compose Record", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After record organize",
				Description: "Triggered after record is organized",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Handle",
										Argument: "namespace.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Namespace Name",
										Argument: "namespace.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Handle",
										Argument: "module.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Module Name",
										Argument: "module.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Created-at",
										Argument: "record.created-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Updated-at",
										Argument: "record.updated-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Deleted-at",
										Argument: "record.deleted-at",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Record Values Any",
										Argument: "record.values.*",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "namespace.handle",
					Types: []string{"String"},
				},
				{
					Name:  "namespace.name",
					Types: []string{"String"},
				},
				{
					Name:  "module.handle",
					Types: []string{"String"},
				},
				{
					Name:  "module.name",
					Types: []string{"String"},
				},
				{
					Name:  "record.created-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.updated-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.deleted-at",
					Types: []string{"String"},
				},
				{
					Name:  "record.values.*",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:sink",
			EventType:    "onRequest",
			Groups:       []string{"System Sink"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "On sink request",
				Description: "Triggered on sink is requestd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Host",
										Argument: "request.host",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Remote-address",
										Argument: "request.remote-address",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Method",
										Argument: "request.method",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Path",
										Argument: "request.path",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Username",
										Argument: "request.username",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Password",
										Argument: "request.password",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Content-type",
										Argument: "request.content-type",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Get Any",
										Argument: "request.get.*",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Post Any",
										Argument: "request.post.*",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Request Header Any",
										Argument: "request.header.*",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "request.host",
					Types: []string{"String"},
				},
				{
					Name:  "request.remote-address",
					Types: []string{"String"},
				},
				{
					Name:  "request.method",
					Types: []string{"String"},
				},
				{
					Name:  "request.path",
					Types: []string{"String"},
				},
				{
					Name:  "request.username",
					Types: []string{"String"},
				},
				{
					Name:  "request.password",
					Types: []string{"String"},
				},
				{
					Name:  "request.content-type",
					Types: []string{"String"},
				},
				{
					Name:  "request.get.*",
					Types: []string{"String"},
				},
				{
					Name:  "request.post.*",
					Types: []string{"String"},
				},
				{
					Name:  "request.header.*",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:mail",
			EventType:    "onReceive",
			Groups:       []string{"System Mail"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "On mail receive",
				Description: "Triggered on mail is received",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header Subject",
										Argument: "message.header.subject",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header From",
										Argument: "message.header.from",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header To",
										Argument: "message.header.to",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header Reply-to",
										Argument: "message.header.reply-to",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header Cc",
										Argument: "message.header.cc",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header Bcc",
										Argument: "message.header.bcc",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "message.header.subject",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.from",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.to",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.reply-to",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.cc",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.bcc",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:mail",
			EventType:    "onSend",
			Groups:       []string{"System Mail"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "On mail send",
				Description: "Triggered on mail is sendd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header Subject",
										Argument: "message.header.subject",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header From",
										Argument: "message.header.from",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header To",
										Argument: "message.header.to",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header Reply-to",
										Argument: "message.header.reply-to",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header Cc",
										Argument: "message.header.cc",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Message Header Bcc",
										Argument: "message.header.bcc",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "message.header.subject",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.from",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.to",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.reply-to",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.cc",
					Types: []string{"String"},
				},
				{
					Name:  "message.header.bcc",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth",
			EventType:    "beforeLogin",
			Groups:       []string{"System Auth", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before auth login",
				Description: "Triggered before auth is logind",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth",
			EventType:    "afterLogin",
			Groups:       []string{"System Auth", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After auth login",
				Description: "Triggered after auth is logind",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth",
			EventType:    "beforeSignup",
			Groups:       []string{"System Auth", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before auth signup",
				Description: "Triggered before auth is signupd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth",
			EventType:    "afterSignup",
			Groups:       []string{"System Auth", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After auth signup",
				Description: "Triggered after auth is signupd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth-client",
			EventType:    "beforeCreate",
			Groups:       []string{"System Auth Client", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before auth client create",
				Description: "Triggered before auth client is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Auth-client Handle",
										Argument: "auth-client.handle",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "auth-client.handle",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth-client",
			EventType:    "afterCreate",
			Groups:       []string{"System Auth Client", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After auth client create",
				Description: "Triggered after auth client is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Auth-client Handle",
										Argument: "auth-client.handle",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "auth-client.handle",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth-client",
			EventType:    "beforeUpdate",
			Groups:       []string{"System Auth Client", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before auth client update",
				Description: "Triggered before auth client is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Auth-client Handle",
										Argument: "auth-client.handle",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "auth-client.handle",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth-client",
			EventType:    "afterUpdate",
			Groups:       []string{"System Auth Client", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After auth client update",
				Description: "Triggered after auth client is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Auth-client Handle",
										Argument: "auth-client.handle",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "auth-client.handle",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth-client",
			EventType:    "beforeDelete",
			Groups:       []string{"System Auth Client", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before auth client delete",
				Description: "Triggered before auth client is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Auth-client Handle",
										Argument: "auth-client.handle",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "auth-client.handle",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:auth-client",
			EventType:    "afterDelete",
			Groups:       []string{"System Auth Client", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After auth client delete",
				Description: "Triggered after auth client is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Auth-client Handle",
										Argument: "auth-client.handle",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "auth-client.handle",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
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
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "beforeSuspend",
			Groups:       []string{"System User", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user suspend",
				Description: "Triggered before user is suspendd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user",
			EventType:    "afterSuspend",
			Groups:       []string{"System User", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user suspend",
				Description: "Triggered after user is suspendd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role",
			EventType:    "beforeCreate",
			Groups:       []string{"System Role", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before role create",
				Description: "Triggered before role is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role",
			EventType:    "afterCreate",
			Groups:       []string{"System Role", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After role create",
				Description: "Triggered after role is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role",
			EventType:    "beforeUpdate",
			Groups:       []string{"System Role", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before role update",
				Description: "Triggered before role is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role",
			EventType:    "afterUpdate",
			Groups:       []string{"System Role", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After role update",
				Description: "Triggered after role is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role",
			EventType:    "beforeDelete",
			Groups:       []string{"System Role", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before role delete",
				Description: "Triggered before role is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role",
			EventType:    "afterDelete",
			Groups:       []string{"System Role", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After role delete",
				Description: "Triggered after role is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role:member",
			EventType:    "beforeAdd",
			Groups:       []string{"System Role Member", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before role add",
				Description: "Triggered before role is addd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role:member",
			EventType:    "afterAdd",
			Groups:       []string{"System Role Member", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After role add",
				Description: "Triggered after role is addd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role:member",
			EventType:    "beforeRemove",
			Groups:       []string{"System Role Member", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before role remove",
				Description: "Triggered before role is removed",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:role:member",
			EventType:    "afterRemove",
			Groups:       []string{"System Role Member", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After role remove",
				Description: "Triggered after role is removed",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Handle",
										Argument: "role.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Handle",
										Argument: "user.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "User Email",
										Argument: "user.email",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.handle",
					Types: []string{"String"},
				},
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
				{
					Name:  "user.handle",
					Types: []string{"String"},
				},
				{
					Name:  "user.email",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "beforeCreate",
			Groups:       []string{"System User Group", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user group create",
				Description: "Triggered before user group is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "afterCreate",
			Groups:       []string{"System User Group", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user group create",
				Description: "Triggered after user group is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "beforeUpdate",
			Groups:       []string{"System User Group", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user group update",
				Description: "Triggered before user group is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "afterUpdate",
			Groups:       []string{"System User Group", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user group update",
				Description: "Triggered after user group is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "beforeDelete",
			Groups:       []string{"System User Group", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user group delete",
				Description: "Triggered before user group is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "afterDelete",
			Groups:       []string{"System User Group", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user group delete",
				Description: "Triggered after user group is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "beforeMemberAdd",
			Groups:       []string{"System User Group", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user group memberadd",
				Description: "Triggered before user group is memberaddd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "afterMemberAdd",
			Groups:       []string{"System User Group", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user group memberadd",
				Description: "Triggered after user group is memberaddd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "beforeMemberRemove",
			Groups:       []string{"System User Group", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before user group memberremove",
				Description: "Triggered before user group is memberremoved",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:user-group",
			EventType:    "afterMemberRemove",
			Groups:       []string{"System User Group", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After user group memberremove",
				Description: "Triggered after user group is memberremoved",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Handle",
										Argument: "userGroup.handle",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Usergroup Name",
										Argument: "userGroup.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "userGroup.handle",
					Types: []string{"String"},
				},
				{
					Name:  "userGroup.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:application",
			EventType:    "beforeCreate",
			Groups:       []string{"System Application", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before application create",
				Description: "Triggered before application is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Application Name",
										Argument: "application.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "application.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:application",
			EventType:    "afterCreate",
			Groups:       []string{"System Application", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After application create",
				Description: "Triggered after application is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Application Name",
										Argument: "application.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "application.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:application",
			EventType:    "beforeUpdate",
			Groups:       []string{"System Application", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before application update",
				Description: "Triggered before application is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Application Name",
										Argument: "application.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "application.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:application",
			EventType:    "afterUpdate",
			Groups:       []string{"System Application", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After application update",
				Description: "Triggered after application is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Application Name",
										Argument: "application.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "application.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:application",
			EventType:    "beforeDelete",
			Groups:       []string{"System Application", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before application delete",
				Description: "Triggered before application is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Application Name",
										Argument: "application.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "application.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:application",
			EventType:    "afterDelete",
			Groups:       []string{"System Application", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After application delete",
				Description: "Triggered after application is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Application Name",
										Argument: "application.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "application.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:queue",
			EventType:    "onMessage",
			Groups:       []string{"System Queue"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "On queue message",
				Description: "Triggered on queue is messaged",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Payload Queue",
										Argument: "payload.queue",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "payload.queue",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:data-privacy-request",
			EventType:    "beforeCreate",
			Groups:       []string{"System Data Privacy Request", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before data privacy request create",
				Description: "Triggered before data privacy request is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:data-privacy-request",
			EventType:    "afterCreate",
			Groups:       []string{"System Data Privacy Request", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After data privacy request create",
				Description: "Triggered after data privacy request is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:data-privacy-request",
			EventType:    "beforeUpdate",
			Groups:       []string{"System Data Privacy Request", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before data privacy request update",
				Description: "Triggered before data privacy request is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:data-privacy-request",
			EventType:    "afterUpdate",
			Groups:       []string{"System Data Privacy Request", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After data privacy request update",
				Description: "Triggered after data privacy request is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:data-privacy-request",
			EventType:    "beforeDelete",
			Groups:       []string{"System Data Privacy Request", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before data privacy request delete",
				Description: "Triggered before data privacy request is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:data-privacy-request",
			EventType:    "afterDelete",
			Groups:       []string{"System Data Privacy Request", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After data privacy request delete",
				Description: "Triggered after data privacy request is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Role Name",
										Argument: "role.name",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "role.name",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "beforeCreate",
			Groups:       []string{"System Reminder", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before reminder create",
				Description: "Triggered before reminder is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "afterCreate",
			Groups:       []string{"System Reminder", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After reminder create",
				Description: "Triggered after reminder is created",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "beforeUpdate",
			Groups:       []string{"System Reminder", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before reminder update",
				Description: "Triggered before reminder is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "afterUpdate",
			Groups:       []string{"System Reminder", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After reminder update",
				Description: "Triggered after reminder is updated",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "beforeDelete",
			Groups:       []string{"System Reminder", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before reminder delete",
				Description: "Triggered before reminder is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "afterDelete",
			Groups:       []string{"System Reminder", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After reminder delete",
				Description: "Triggered after reminder is deleted",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "beforeDismiss",
			Groups:       []string{"System Reminder", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before reminder dismiss",
				Description: "Triggered before reminder is dismissd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "afterDismiss",
			Groups:       []string{"System Reminder", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After reminder dismiss",
				Description: "Triggered after reminder is dismissd",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "beforeSnooze",
			Groups:       []string{"System Reminder", "Before"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "Before reminder snooze",
				Description: "Triggered before reminder is snoozed",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
		types.ConstructTrigger{
			ResourceType: "system:reminder",
			EventType:    "afterSnooze",
			Groups:       []string{"System Reminder", "After"},
			Meta: &types.ConstructTriggerMeta{
				Short:       "After reminder snooze",
				Description: "Triggered after reminder is snoozed",
			},
			Segments: []types.ConstructSegment{
				{
					Sections: []types.ConstructSection{
						{
							Elements: []types.SectionElement{
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Resource",
										Argument: "reminder.resource",
									},
								},
								{
									Input: types.SectionElementInput{
										Type:     "String",
										Label:    "Reminder Assigned-to",
										Argument: "reminder.assigned-to",
									},
								},
							},
						},
					},
				},
			},
			Constraints: []types.ConstructTriggerConstraint{
				{
					Name:  "reminder.resource",
					Types: []string{"String"},
				},
				{
					Name:  "reminder.assigned-to",
					Types: []string{"String"},
				},
			},
		},
	)
}
