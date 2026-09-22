// Package weburl builds links into the Human webapp.
//
// It exists because a tool result that says what changed, and cannot say where
// to go and look at it, leaves the caller guessing at a path. The guesses were
// already visible in the codebase: the discovery indexer emitted
// /compose/ns/{slug}/pages/{id}, a prefix the unified webapp does not route.
//
// The paths here mirror client/web/unify/src/sections/*/routes; routes_test.go
// reads those files and fails when one stops matching, because a link table
// nobody checks is a link table that rots.
package weburl

import (
	"fmt"
	"strings"

	"github.com/crusttech/human/server/pkg/options"
)

// Base is where the webapp is served from, computed the way app/boot_levels.go
// computes `webapp.base-url`: the API host when the server serves the webapp
// itself, the separately configured webapp host when it does not.
func Base() string {
	base := options.EnvString("HTTP_BASE_URL", "")
	webapp := options.EnvString("HTTP_WEBAPP_BASE_URL", "/")

	if options.EnvBool("HTTP_WEBAPP_ENABLED", true) {
		return options.FullURL(base, webapp)
	}

	return options.FullWebappURL(base, webapp)
}

// abs turns a webapp path into a link. An empty path means the resource has no
// screen, and stays empty rather than becoming a link to the root.
func abs(path string) string {
	if path == "" {
		return ""
	}
	return strings.TrimRight(Base(), "/") + path
}

// Resources with no screen of their own return "" from every builder here:
// reminders, themes, skills, event types and triggers are edited elsewhere or
// not at all, and inventing a link for them would be worse than none.

// ComposeNamespace is a namespace's page list — where opening it lands.
func ComposeNamespace(slug string) string {
	if slug == "" {
		return ""
	}
	return abs("/compose/namespace/" + slug)
}

// ComposeNamespaceEdit is the namespace's own settings form.
func ComposeNamespaceEdit(slug string) string {
	if slug == "" {
		return ""
	}
	return abs("/compose/namespaces/edit/" + slug)
}

// ComposePage is the page as a user sees it.
func ComposePage(slug string, pageID uint64) string {
	return composePath(slug, pageID, "/compose/namespace/%s/pages/%d")
}

// ComposePageBuilder is the block builder — where a layout is arranged. Page
// layouts have no screen of their own; they are edited here.
func ComposePageBuilder(slug string, pageID uint64) string {
	return composePath(slug, pageID, "/compose/namespace/%s/admin/pages/%d/builder")
}

// ComposeModuleEdit is the module's field editor.
func ComposeModuleEdit(slug string, moduleID uint64) string {
	return composePath(slug, moduleID, "/compose/namespace/%s/admin/modules/%d/edit")
}

// ComposeChartEdit is the chart configurator.
func ComposeChartEdit(slug string, chartID uint64) string {
	return composePath(slug, chartID, "/compose/namespace/%s/admin/charts/%d/edit")
}

// ComposeRecord is the record as the module admin shows it, which is the one
// view that exists for every record. A record also opens on any page bound to
// its module, and that link needs the pageID the caller chose.
func ComposeRecord(slug string, moduleID, recordID uint64) string {
	if slug == "" || moduleID == 0 || recordID == 0 {
		return ""
	}
	return abs(fmt.Sprintf("/compose/namespace/%s/admin/modules/%d/records/%d", slug, moduleID, recordID))
}

// ComposeRecordOnPage is the record opened on a specific page.
func ComposeRecordOnPage(slug string, pageID, recordID uint64) string {
	if slug == "" || pageID == 0 || recordID == 0 {
		return ""
	}
	return abs(fmt.Sprintf("/compose/namespace/%s/pages/%d/records/%d", slug, pageID, recordID))
}

// TAQ is the Trigger Action Query builder.
func TAQ(taqID uint64) string {
	return idPath(taqID, "/taq/builder/%d")
}

// Workflow is the workflow editor. A trigger has no screen of its own — it is
// edited inside the workflow it fires, so a trigger links here.
func Workflow(workflowID uint64) string {
	return idPath(workflowID, "/workflow/%d/edit")
}

// Agent is the agent editor.
func Agent(agentID uint64) string {
	return idPath(agentID, "/agentic/%d/edit")
}

// Chatbot is the chatbot editor.
func Chatbot(chatbotID uint64) string {
	return idPath(chatbotID, "/chatbot/%d/edit")
}

// User is the user's admin record.
func User(userID uint64) string {
	return idPath(userID, "/admin/system/users/%d")
}

// UserGroup is the group's admin record.
func UserGroup(userGroupID uint64) string {
	return idPath(userGroupID, "/admin/system/user-groups/%d")
}

// Role is the role's admin record.
func Role(roleID uint64) string {
	return idPath(roleID, "/admin/system/roles/%d")
}

// Application is the application's admin record.
func Application(applicationID uint64) string {
	return idPath(applicationID, "/admin/system/applications/%d")
}

// CustomApplication is a custom application as the app view shows it.
func CustomApplication(applicationID uint64) string {
	return idPath(applicationID, "/app/%d")
}

// AuthClient is the auth client's admin record.
func AuthClient(authClientID uint64) string {
	return idPath(authClientID, "/admin/system/auth-clients/%d")
}

func composePath(slug string, id uint64, format string) string {
	if slug == "" || id == 0 {
		return ""
	}
	return abs(fmt.Sprintf(format, slug, id))
}

func idPath(id uint64, format string) string {
	if id == 0 {
		return ""
	}
	return abs(fmt.Sprintf(format, id))
}
