package automation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

workflow: schema.#resource & {
	rbac: {
		operations: {
		  "read": description: "Read workflow"
		  "update": description: "Update workflow"
		  "delete": description: "Delete workflow"
  		"undelete": description: "Undelete workflow"
		  "execute": description: "Execute workflow"
	  	"triggers.manage": description: "Manage workflow triggers"
  		"sessions.manage": description: "Manage workflow sessions"
		}
	}
}
