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
	assert.Equal(t, "app/42", sysService.CustomApplicationPath(42))
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
	// The last two are what reach the skill before a page exists: a session
	// modelling for an app it is about to build needs the rules that say the
	// model comes first.
	assert.ElementsMatch(t, []string{
		"system_application_create",
		"system_application_update",
		"system_application_source_set",
		"system_application_source_get",
		"compose_module_create",
		"compose_namespace_create",
	}, skill.Triggers)

	snippet := indentedBlockContaining(skill.Body, "human:hello")
	require.NotEmpty(t, snippet,
		"the bridge snippet was not found in the custom_app skill as an indented block; "+
			"this test is asserting nothing until it is")

	require.NoError(t, sysService.CheckApplicationSource("<script>\n"+snippet+"\n</script>"))
}
