package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

page: schema.#resource & {
	parents: [
		{handle: "namespace"},
	]

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
		}
	}

	locale: {
		extended: true

		keys: {
			title: {}
			description: {}
			recordToolbarButtonNewLabel: {
				path: ["recordToolbar", "new", "label"]
				customHandler: true
			}
			recordToolbarButtonEditLabel: {
				path: ["recordToolbar", "edit", "label"]
				customHandler: true
			}
			recordToolbarButtonSubmitLabel: {
				path: ["recordToolbar", "submit", "label"]
				customHandler: true
			}
			recordToolbarButtonDeleteLabel: {
				path: ["recordToolbar", "delete", "label"]
				customHandler: true
			}
			recordToolbarButtonCloneLabel: {
				path: ["recordToolbar", "clone", "label"]
				customHandler: true
			}
			recordToolbarButtonBackLabel: {
				path: ["recordToolbar", "back", "label"]
				customHandler: true
			}
			blockTitle: {
				path: ["pageBlock", {part: "blockID", var: true}, "title"]
				customHandler: true
			}
			blockDescription: {
				path: ["pageBlock", {part: "blockID", var: true}, "description"]
				customHandler: true
			}
			blockAutomationButtonLabel: {
				path: ["pageBlock", {part: "blockID", var: true}, "button", {part: "buttonID", var: true}, "label"]
				customHandler: true
			}
			blockContentBody: {
				path: ["pageBlock", {part: "blockID", var: true}, "content", "body"]
				customHandler: true
			}
		}
	}
}
