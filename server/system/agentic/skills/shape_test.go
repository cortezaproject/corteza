package skills

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLibraryShape holds every skill to the shape a reader can rely on: its
// name is its file name, its description is one sentence a listing can show,
// its body has a title and sections and is long enough to carry rules rather
// than a slogan. Content is reviewed by hand against the tools it names; this
// is the part a test can hold.
func TestLibraryShape(t *testing.T) {
	entries, err := fs.ReadDir(libraryFS, "library")
	require.NoError(t, err)

	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), ".md")
		t.Run(name, func(t *testing.T) {
			data, err := fs.ReadFile(libraryFS, "library/"+e.Name())
			require.NoError(t, err)
			s, err := parseSkill(data)
			require.NoError(t, err)

			require.Equal(t, name, s.Name, "the skill's name is its file name, so system_skill_lookup and the file agree")
			require.GreaterOrEqual(t, len(s.Description), 40, "the description is what a listing shows; one line that says what the skill decides")
			require.LessOrEqual(t, len(s.Description), 260, "a description longer than a sentence is body text")
			require.True(t, strings.HasPrefix(s.Body, "# "), "the body opens with its title")
			require.Contains(t, s.Body, "\n## ", "a skill has sections; one paragraph is a tool description, not a skill")
			require.GreaterOrEqual(t, len(s.Body), 600, "a skill carries working rules; under 600 characters it is a hint")
			require.NotContains(t, s.Body, "TODO")
			for _, tr := range s.Triggers {
				require.Regexp(t, `^[a-z][a-z0-9]*(_[a-z0-9]+)+$`, tr, "triggers are tool names")
			}
		})
	}
}

// TestParseSkillRejectsUnknownFrontmatter: a field the loader does not read
// fails the load rather than being carried as decoration.
func TestParseSkillRejectsUnknownFrontmatter(t *testing.T) {
	_, err := parseSkill([]byte("---\nname: x\nimportance: high\ntriggers: [a_b]\n---\n# x\n\n## y\n"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "importance")
}
