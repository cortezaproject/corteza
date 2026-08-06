package tools

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// The rules CLAUDE.md states, enforced rather than remembered.
const (
	// maxSubject is measured from this repo's own history: median 44, longest
	// 77. Seventy-two is the conventional ceiling and comfortably above what
	// anyone here actually writes.
	maxSubject = 72

	// minSubject catches "fix", "wip" and "update" — subjects that pass every
	// other rule and say nothing.
	minSubject = 12
)

// aiTrailer matches the attribution this repo does not use. Checked over the
// whole message, not just its end, because a model asked for a commit message
// will happily put it anywhere.
var aiTrailer = regexp.MustCompile(`(?i)co-authored-by:.*(claude|anthropic|copilot|gpt)|generated with \[?claude|🤖`)

// pastTense catches the most common non-imperative openings. Deliberately a
// short list of frequent offenders rather than a grammar: a wrong rejection
// costs more than a missed one, since the author can always see the subject.
var pastTense = regexp.MustCompile(`(?i)^(added|fixed|updated|removed|changed|created|deleted|renamed|moved|refactored|implemented|adds|fixes|updates|removes|changes|creates|deletes|renames|moves|refactors|implements)\b`)

type commitReport struct {
	Committed bool     `json:"committed"`
	Commit    string   `json:"commit,omitempty"`
	Subject   string   `json:"subject,omitempty"`
	Files     []string `json:"files,omitempty"`
	Refusals  []string `json:"refusals,omitempty"`
	Warnings  []string `json:"warnings,omitempty"`
	Note      string   `json:"note"`
}

func registerCommit(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_commit_create",
			mcp.WithDescription(
				"Commit staged or named files with this repo's conventions enforced rather than remembered: "+
					"a short imperative subject, no AI attribution of any kind, and formatting applied before "+
					"staging. It refuses rather than fixing a message for you, because the subject is the "+
					"author's to write. Commit only when the human has asked you to — this tool enforces how "+
					"a commit is made, never whether one should be.",
			),
			mcp.WithString("subject", mcp.Required(), mcp.Description(
				"One imperative line, under 72 characters, no trailing period. \"Fix the parser\" not "+
					"\"Fixed the parser\" and not \"fix\".")),
			mcp.WithString("body", mcp.Description(
				"Optional, two or three lines, only when the why is not obvious from the diff. Most commits "+
					"here have no body.")),
			mcp.WithString("files", mcp.Description(
				"Space-separated repo-relative paths to stage. Omit to commit what is already staged, which "+
					"is how you keep a commit atomic when the working tree holds unrelated work.")),
			mcp.WithBoolean("skipFormat", mcp.Description(
				"Skip formatting before staging. Off by default; formatting after staging is how an "+
					"unformatted file reaches history.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskWrite),
		),
		"Create a commit",
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, err := toolkit.Args(req)
			if err != nil {
				return nil, err
			}

			subject, err := toolkit.ReqStr(args, "subject")
			if err != nil {
				return nil, err
			}

			return commit(ctx, root, commitInput{
				Subject:    subject,
				Body:       toolkit.Str(args, "body"),
				Files:      strings.Fields(toolkit.Str(args, "files")),
				SkipFormat: toolkit.Bool(args, "skipFormat"),
			})
		},
	)
}

type commitInput struct {
	Subject    string
	Body       string
	Files      []string
	SkipFormat bool
}

func commit(ctx context.Context, root string, in commitInput) (*mcp.CallToolResult, error) {
	git := gitRunner(root)
	out := commitReport{Subject: in.Subject}

	out.Refusals = checkMessage(in.Subject, in.Body)
	if len(out.Refusals) > 0 {
		out.Note = "nothing was committed: fix the message and call again"
		return toolkit.JSONResult(out)
	}

	// Format before staging, not after: the whole point of the rule is that
	// what gets staged is what was formatted.
	if !in.SkipFormat && len(in.Files) > 0 {
		formatFiles(ctx, root, in.Files)
	}

	if len(in.Files) > 0 {
		// -A so a deleted path still stages: plain `git add -- <path>` fails
		// with "pathspec did not match any files" when the file is gone, which
		// makes naming a file you just removed an error rather than a commit.
		if _, err := git(ctx, append([]string{"add", "-A", "--"}, in.Files...)...); err != nil {
			return nil, err
		}
	}

	staged, err := git(ctx, "diff", "--cached", "--name-only")
	if err != nil {
		return nil, err
	}

	out.Files = nonEmptyLines(staged)
	if len(out.Files) == 0 {
		out.Refusals = append(out.Refusals, "nothing is staged: pass files, or stage them first")
		out.Note = "nothing was committed"
		return toolkit.JSONResult(out)
	}

	out.Warnings = append(out.Warnings, checkAtomic(out.Files)...)

	message := in.Subject
	if in.Body != "" {
		message += "\n\n" + in.Body
	}

	if _, err := git(ctx, "commit", "-m", message); err != nil {
		return nil, err
	}

	sha, err := git(ctx, "log", "-1", "--pretty=%h")
	if err != nil {
		return nil, err
	}

	out.Committed = true
	out.Commit = sha

	out.Note = fmt.Sprintf("committed %s with %d file(s)", sha, len(out.Files))
	if len(out.Warnings) > 0 {
		out.Note += ". The warnings above are worth reading: an unrelated change in the same commit is " +
			"the thing that makes history hard to read later"
	}

	return toolkit.JSONResult(out)
}

// checkMessage returns the reasons this message may not be committed.
func checkMessage(subject, body string) []string {
	var out []string

	full := subject + "\n" + body

	if aiTrailer.MatchString(full) {
		out = append(out, "the message carries AI attribution. Commits here are authored by the human, "+
			"with no co-author trailer, no generated-with line and no robot emoji")
	}

	switch {
	case len(subject) > maxSubject:
		out = append(out, fmt.Sprintf("subject is %d characters, over the %d limit — say less, or move the "+
			"detail into the body", len(subject), maxSubject))
	case len(subject) < minSubject:
		out = append(out, fmt.Sprintf("subject %q is too short to mean anything on its own", subject))
	}

	if strings.HasSuffix(subject, ".") {
		out = append(out, "subject ends with a period; this repo's subjects do not")
	}

	if strings.Contains(subject, "\n") {
		out = append(out, "subject must be one line — put the rest in body")
	}

	if pastTense.MatchString(subject) {
		out = append(out, fmt.Sprintf("subject %q is not imperative: write what the commit does, as an "+
			"instruction — \"Fix the parser\", not \"Fixed the parser\"", subject))
	}

	if first := strings.Fields(subject); len(first) > 0 {
		if r := []rune(first[0]); len(r) > 0 && r[0] >= 'a' && r[0] <= 'z' {
			out = append(out, "subject should start with a capital letter, matching the rest of the history")
		}
	}

	if strings.Count(body, "\n") > 4 {
		out = append(out, "body is longer than three or four lines; this repo keeps the why short, and the "+
			"diff says the what")
	}

	return out
}

// checkAtomic warns when one commit spans work that reads as unrelated.
//
// A warning, not a refusal: "docs, bugfix and cleanup separately" is a rule
// about intent, and no file-path heuristic can tell a documented fix from a
// doc change that happens to sit beside one. Refusing on a guess would make the
// tool something to work around.
func checkAtomic(files []string) []string {
	kinds := map[string]bool{}

	for _, f := range files {
		switch {
		case strings.HasSuffix(f, ".intent.md"):
			kinds["intent doc"] = true
		case strings.HasSuffix(f, ".md"):
			kinds["docs"] = true
		case strings.Contains(f, "/locale/") || strings.HasPrefix(f, "locale/"):
			kinds["translations"] = true
		case strings.HasSuffix(f, "_test.go") || strings.Contains(f, ".test.") || strings.Contains(f, ".spec."):
			kinds["tests"] = true
		case strings.Contains(f, ".gen."):
			kinds["generated"] = true
		default:
			kinds["code"] = true
		}
	}

	// Tests belong with the code they cover, and generated files belong with
	// the definition that produced them; neither pairing is a mixed commit.
	delete(kinds, "tests")
	delete(kinds, "generated")

	if len(kinds) < 2 {
		return nil
	}

	named := make([]string, 0, len(kinds))
	for k := range kinds {
		named = append(named, k)
	}

	return []string{fmt.Sprintf("this commit mixes %s. The convention here is separate atomic commits — "+
		"docs, bugfix and cleanup apart — so split it unless these genuinely belong together",
		strings.Join(sortedStrings(named), " and "))}
}

func nonEmptyLines(s string) []string {
	out := make([]string, 0, 8)
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out
}
