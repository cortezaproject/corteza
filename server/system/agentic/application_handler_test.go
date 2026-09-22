package agentic

import (
	"strings"
	"testing"

	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/system/agentic/skills"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleApp = `<!doctype html>
<html><body><div id="app"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/chart.js/4.4.1/chart.umd.min.js"></script>
<script>
window.SAMPLE = { 'records.list': () => ({ records: [], refs: {} }) }
window.human = window.human || (() => { return {} })()
</script>
</body></html>`

func TestApplicationSourcePatchReplacesTheOneMatch(t *testing.T) {
	out, err := applyApplicationSourcePatch("<h1>Leads</h1><p>Leads by owner</p>", "<h1>Leads</h1>", "<h1>Deals</h1>")
	require.NoError(t, err)
	assert.Equal(t, "<h1>Deals</h1><p>Leads by owner</p>", out)
}

func TestApplicationSourcePatchDeletesOnAnEmptyReplacement(t *testing.T) {
	out, err := applyApplicationSourcePatch("<h1>Leads</h1><hr>", "<hr>", "")
	require.NoError(t, err)
	assert.Equal(t, "<h1>Leads</h1>", out)
}

func TestApplicationSourcePatchRefusesNoMatch(t *testing.T) {
	_, err := applyApplicationSourcePatch(sampleApp, "<h1>Missing</h1>", "x")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not appear")
	assert.Contains(t, err.Error(), "system_application_source_get")
}

// The count is in the message because it is the only thing that tells the
// caller how much more surrounding text to include.
func TestApplicationSourcePatchRefusesSeveralMatchesAndSaysHowMany(t *testing.T) {
	_, err := applyApplicationSourcePatch("<td>-</td><td>-</td><td>-</td>", "<td>-</td>", "<td>0</td>")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "matches 3 times")
}

func TestApplicationSourcePatchRefusesAnEmptyOldString(t *testing.T) {
	_, err := applyApplicationSourcePatch(sampleApp, "", "x")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "old_string is empty")
}

func TestApplicationSourcePatchRefusesAnApplicationWithNoSource(t *testing.T) {
	_, err := applyApplicationSourcePatch("", "<h1>Leads</h1>", "x")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no source stored yet")
}

func TestApplicationSourceAcceptsAPlainDocument(t *testing.T) {
	require.NoError(t, checkApplicationSource(sampleApp))
}

func TestApplicationSourceRefusesWhatTheSandboxCannotRun(t *testing.T) {
	cases := []struct {
		name    string
		source  string
		message string
	}{
		{"module script", `<script type="module">const a = 1</script>`, "script-src"},
		{"single-quoted module script", "<script type='module'>const a = 1</script>", "script-src"},
		{"import statement", "<script>\nimport { h } from 'https://esm.sh/preact'\n</script>", "line 2"},
		{"export statement", "<script>\nexport const render = () => {}\n</script>", "ES module statement"},
		{"react import", `<script>const { useState } = from "react"</script>`, "React"},
		{"fetch", `<script>fetch('/api/compose').then(r => r.json())</script>`, "connect-src"},
		{"window fetch", `<script>window.fetch('/api/compose')</script>`, "connect-src"},
		{"xhr", `<script>const r = new XMLHttpRequest()</script>`, "connect-src"},
		{"websocket", `<script>const s = new WebSocket('wss://example.com')</script>`, "connect-src"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := checkApplicationSource(c.source)
			require.Error(t, err)
			assert.Contains(t, err.Error(), c.message)
		})
	}
}

// A name that merely ends in one of the refused ones is somebody's own
// function, and refusing it would send the caller looking for a call it never
// wrote.
func TestApplicationSourceAllowsIdentifiersThatOnlyEndInARefusedName(t *testing.T) {
	require.NoError(t, checkApplicationSource(`<script>const prefetch = () => {}; prefetch()</script>`))
	require.NoError(t, checkApplicationSource(`<script>function importRows () {}</script>`))
}

func TestApplicationSourceRefusesAnEmptyDocument(t *testing.T) {
	err := checkApplicationSource("   \n\t ")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

// The service caps the source too; failing here is what keeps a 300 KB argument
// from being carried through the whole write before it is refused.
func TestApplicationSourceRefusesOverTheSizeCap(t *testing.T) {
	oversize := "<html>" + strings.Repeat("x", sysService.ApplicationSourceMaxSize) + "</html>"

	err := checkApplicationSource(oversize)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "the limit is")

	require.NoError(t, checkApplicationSource("<html>"+strings.Repeat("x", sysService.ApplicationSourceMaxSize-20)+"</html>"))
}

func TestCustomApplicationNeedsURL(t *testing.T) {
	custom := func(url string) *sysTypes.Application {
		return &sysTypes.Application{
			ID:    42,
			Unify: &sysTypes.ApplicationUnify{Kind: sysService.ApplicationKindCustom, Url: url},
		}
	}

	assert.True(t, customApplicationNeedsURL(custom("")))
	assert.False(t, customApplicationNeedsURL(custom("app/7")), "an explicit url is the caller's")
	assert.False(t, customApplicationNeedsURL(&sysTypes.Application{ID: 42, Unify: &sysTypes.ApplicationUnify{}}))
	assert.False(t, customApplicationNeedsURL(&sysTypes.Application{ID: 42}))
	assert.Equal(t, "app/42", customApplicationPath(42))
}

func TestApplicationHandlerRegistersTheSourceTools(t *testing.T) {
	reg := &stubRegistrar{}
	ApplicationHandler(reg)

	get, ok := reg.byName("system_application_source_get")
	require.True(t, ok)
	assert.Equal(t, hmcp.RiskRead, hmcp.RiskOf(get))
	assert.Equal(t, []hmcp.Group{hmcp.GroupConfiguring}, hmcp.GroupsOf(get))

	set, ok := reg.byName("system_application_source_set")
	require.True(t, ok)
	assert.Equal(t, hmcp.RiskWrite, hmcp.RiskOf(set))
	assert.Equal(t, []hmcp.Group{hmcp.GroupConfiguring}, hmcp.GroupsOf(set))

	// Every argument the handler reads has to be declared, or
	// rejectUnknownArguments refuses the call before it arrives.
	for _, arg := range []string{"application", "source", "old_string", "new_string", "namespace", "modules"} {
		assert.Containsf(t, set.InputSchema.Properties, arg, "system_application_source_set must declare %q", arg)
	}
}

// The skill hands out a bridge snippet that every custom app pastes in. If the
// guard refused it, the tool would reject the first app written to the skill.
func TestTheSkillsBridgeSnippetSurvivesTheSourceGuard(t *testing.T) {
	lib, err := skills.LoadLibrary()
	require.NoError(t, err)

	var skill *skills.Skill
	for _, s := range lib.All() {
		if s.Name == "custom_app" {
			skill = s
		}
	}
	require.NotNil(t, skill, "the custom_app skill is missing from the library")
	assert.ElementsMatch(t, []string{
		"system_application_create",
		"system_application_update",
		"system_application_source_set",
		"system_application_source_get",
	}, skill.Triggers)

	snippet := indentedBlockContaining(skill.Body, "human:hello")
	require.NotEmpty(t, snippet,
		"the bridge snippet was not found in the custom_app skill as an indented block; "+
			"this test is asserting nothing until it is")

	require.NoError(t, checkApplicationSource("<script>\n"+snippet+"\n</script>"))
}

// indentedBlockContaining returns the run of indented lines holding marker.
func indentedBlockContaining(body, marker string) string {
	var block []string
	for _, line := range strings.Split(body, "\n") {
		switch {
		case strings.HasPrefix(line, "    "):
			block = append(block, strings.TrimPrefix(line, "    "))
		case strings.TrimSpace(line) == "" && len(block) > 0:
			block = append(block, "")
		default:
			if joined := strings.Join(block, "\n"); strings.Contains(joined, marker) {
				return joined
			}
			block = nil
		}
	}
	if joined := strings.Join(block, "\n"); strings.Contains(joined, marker) {
		return joined
	}
	return ""
}
