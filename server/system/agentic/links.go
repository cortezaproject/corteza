package agentic

import (
	"github.com/crusttech/human/server/pkg/weburl"
	sysTypes "github.com/crusttech/human/server/system/types"
)

// Links a single-item result carries, so a caller can say where to go and look
// at what it changed. A builder that cannot resolve an ID returns "", and
// toolkit.JSONResultWith drops the key rather than emitting a blank one.
//
// Reminders and themes are absent on purpose: neither has a screen of its own
// in the webapp, and a link to a list that does not show the thing reads as an
// answer when it is not one.

func userLinks(u *sysTypes.User) map[string]string {
	if u == nil {
		return nil
	}
	return map[string]string{"url": weburl.User(u.ID)}
}

func userGroupLinks(g *sysTypes.UserGroup) map[string]string {
	if g == nil {
		return nil
	}
	return map[string]string{"url": weburl.UserGroup(g.ID)}
}

func roleLinks(r *sysTypes.Role) map[string]string {
	if r == nil {
		return nil
	}
	return map[string]string{"url": weburl.Role(r.ID)}
}

func applicationLinks(app *sysTypes.Application) map[string]string {
	if app == nil {
		return nil
	}
	return map[string]string{"url": weburl.Application(app.ID)}
}

func agentLinks(a *sysTypes.Agent) map[string]string {
	if a == nil {
		return nil
	}
	return map[string]string{"url": weburl.Agent(a.ID)}
}

func chatbotLinks(c *sysTypes.Chatbot) map[string]string {
	if c == nil {
		return nil
	}
	return map[string]string{"url": weburl.Chatbot(c.ID)}
}

func authClientLinks(c *sysTypes.AuthClient) map[string]string {
	if c == nil {
		return nil
	}
	return map[string]string{"url": weburl.AuthClient(c.ID)}
}

// withNote adds a note to a result's extra fields, leaving them alone when
// there is nothing to say.
func withNote(extra map[string]string, note string) map[string]string {
	if note == "" {
		return extra
	}
	if extra == nil {
		extra = make(map[string]string, 1)
	}
	extra["note"] = note
	return extra
}
