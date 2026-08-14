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

	// maxBody is one short sentence. Measured in characters rather than lines
	// because a wrapped sentence is still one sentence, while three tight lines
	// are already an essay.
	maxBody = 160
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
				"Commit staged or named files with this repo's conventions enforced rather than remembered. "+
					"The convention is stated in one place — CLAUDE.md, section 'Commit convention' — and this "+
					"tool enforces the half of it a machine can decide: a short imperative subject, no AI "+
					"attribution of any kind, formatting applied before staging. It refuses rather than fixing "+
					"a message for you, because the subject is the author's to write, and warns rather than "+
					"refuses where only judgement can tell (unrelated changes sharing a commit, an intent "+
					"change missing from one). The two halves it cannot check, from that same section: commit "+
					"without being asked but only once nothing on the issue is still open, and push nothing.",
			),
			mcp.WithString("subject", mcp.Required(), mcp.Description(
				"One imperative line, under 72 characters, no trailing period. \"Fix the parser\" not "+
					"\"Fixed the parser\" and not \"fix\".")),
			mcp.WithString("body", mcp.Description(
				"Optional and usually omitted: at most one short sentence naming what changed, for when the "+
					"subject cannot hold it. Not prose, not the reasoning — the diff shows what, the intent "+
					"doc holds why.")),
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
	out.Warnings = append(out.Warnings, checkIntentPaired(ctx, root, out.Files)...)

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
		out.Note += ". The warnings above are worth reading: what makes history hard to read later is " +
			"an unrelated change sharing a commit, or a related one missing from it"
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

	if trimmed := strings.TrimSpace(body); trimmed != "" {
		if len(trimmed) > maxBody {
			out = append(out, fmt.Sprintf("body is %d characters, over the %d limit — a body here is one "+
				"short sentence naming what changed, not prose and not the reasoning",
				len(trimmed), maxBody))
		}

		if strings.Contains(trimmed, "\n\n") {
			out = append(out, "body has more than one paragraph; it is one short sentence or nothing")
		}
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
		case isIntentPath(f):
			kinds["intent"] = true
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

	// Tests belong with the code they cover, generated files with the definition
	// that produced them, and an intent doc with the code it governs: the doc
	// states what that very change made true, so splitting them leaves history
	// with a commit whose doc contradicts its code. None of the three is a mixed
	// commit.
	delete(kinds, "tests")
	delete(kinds, "generated")
	delete(kinds, "intent")

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

// isIntentPath covers both halves of an intent change: the docs themselves and
// the lock the sync writes, which has to travel with them or the next check
// reports drift that was already reconciled.
func isIntentPath(f string) bool {
	return strings.HasSuffix(f, ".intent.md") || strings.HasPrefix(f, ".intent/")
}

// checkIntentPaired warns when a commit changes code an intent doc governs and
// leaves the intent side for later — a commit that looks complete and is not.
// A warning, not a refusal: reconciling a doc is the human's judgement call.
func checkIntentPaired(ctx context.Context, root string, files []string) []string {
	for _, f := range files {
		if isIntentPath(f) {
			return nil
		}
	}

	staged := make(map[string]bool, len(files))
	for _, f := range files {
		staged[f] = true
	}

	var governed []intentDrift
	for _, d := range intentDrifted(ctx, root) {
		if staged[d.File] {
			governed = append(governed, d)
		}
	}

	if len(governed) == 0 {
		return nil
	}

	docs := map[string]bool{}
	named := make([]string, 0, len(governed))
	for _, d := range governed {
		if d.Doc == "" || docs[d.Doc] {
			continue
		}
		docs[d.Doc] = true
		named = append(named, d.Doc)
	}

	return []string{fmt.Sprintf("this commit changes code governed by %s, but carries no intent change. "+
		"The doc and the code it describes belong in one commit — reconcile it, run "+
		"'node .intent/intent.mjs sync <files>', and include the doc and .intent/intent.lock.json here",
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
