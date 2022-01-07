package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

component: schema.#component & {
	handle: "system"

	resources: {
		"apigw-route": apigwRoute
		"application": application
		"auth-client": authClient
		"queue":       queue
		"report":      report
		"role":        role
		"template":    template
		"user":        user
		"gig":         gig
	}

	rbac: operations: {
		"action-log.read": description: "Access to action log"

		"settings.read": description:       "Read system settings"
		"settings.manage": description:     "Manage system settings"
		"auth-client.create": description:  "Create auth clients"
		"auth-clients.search": description: "List, search or filter auth clients"

		"role.create": description:  "Create roles"
		"roles.search": description: "List, search or filter roles"

		"user.create": description:  "Create users"
		"users.search": description: "List, search or filter users"

		"application.create": description:      "Create applications"
		"applications.search": description:     "List, search or filter auth clients"
		"application.flag.self": description:   "Manage private flags for applications"
		"application.flag.global": description: "Manage global flags for applications"

		"template.create": description:  "Create template"
		"templates.search": description: "List, search or filter templates"

		"report.create": description:  "Create report"
		"reports.search": description: "List, search or filter reports"

		"reminder.assign": description: " Assign reminders"

		"queue.create": description:  "Create messagebus queues"
		"queues.search": description: "List, search or filter messagebus queues"

		"apigw-route.create": description:  "Create API gateway route"
		"apigw-routes.search": description: "List search or filter API gateway routes"

		"resource-translations.manage": description: "List, search, create, or update resource translations"

		"gig.create": description:  "Create gigs"
		"gigs.search": description: "List, search or filter gigs"
	}
}
