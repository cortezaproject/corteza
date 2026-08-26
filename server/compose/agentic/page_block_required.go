package agentic

import (
	"fmt"
	"sort"
	"strings"

	cmpTypes "github.com/crusttech/human/server/compose/types"
)

// requiredBlockOptions is what each block kind needs before it renders as
// anything but an empty panel.
//
// It mirrors PageBlock.validate() in lib/js, which is where the rule really
// lives: the page builder reads it, badges the block and refuses the save. The
// API had no equivalent, so a page built through these tools was accepted with
// a RecordList that names no module and a Chart that names no chart — silently,
// and looking exactly like a page that worked.
//
// Only the top-level requirement of each kind is mirrored. The nested ones a
// few kinds add (every Automation button needing a target, say) are where a
// mirror starts to drift from the thing it mirrors; the exception is a Metric's
// per-tile module, which is the mistake people actually make.
var requiredBlockOptions = map[string][]string{
	"AgentChat":       {"allowedAgentIDs"},
	"Automation":      {"buttons"},
	"Calendar":        {"feeds"},
	"Chart":           {"chartID"},
	"ChatbotInbox":    {"chatbotIDs"},
	"Comment":         {"moduleID", "contentField"},
	"Content":         {"body"},
	"Geometry":        {"feeds"},
	"IFrame":          {"src"},
	"Metric":          {"metrics"},
	"Navigation":      {"navigationItems"},
	"Progress":        {"value"},
	"RecordList":      {"moduleID"},
	"RecordOrganizer": {"moduleID"},
	"Tabs":            {"tabs"},
}

// blockReadinessNote names the blocks that will render empty, or "" when every
// block is configured.
//
// A note rather than a refusal: a caller may be laying a page out before
// filling it in, and a write that is rejected after the blocks were numbered
// is worse than one that says what is still missing.
func blockReadinessNote(blocks []cmpTypes.PageBlock) string {
	var unfinished []string

	for i, b := range blocks {
		missing := missingBlockOptions(b)
		if len(missing) == 0 {
			continue
		}
		unfinished = append(unfinished, fmt.Sprintf(
			"block %s (%s) needs %s",
			blockRef(b, i), b.Kind, strings.Join(missing, " and "),
		))
	}

	if len(unfinished) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"%d of %d blocks are not configured and will render empty — %s. The page builder refuses to save a page in this state; these tools do not, so fix them with compose_page_update.",
		len(unfinished), len(blocks), strings.Join(unfinished, "; "),
	)
}

// blockRef identifies a block the way the caller can act on it: by blockID once
// the server has assigned one, and by position while it has not.
func blockRef(b cmpTypes.PageBlock, i int) string {
	if b.BlockID > 0 {
		return fmt.Sprintf("%d", b.BlockID)
	}
	return fmt.Sprintf("#%d", i+1)
}

func missingBlockOptions(b cmpTypes.PageBlock) []string {
	var missing []string

	for _, opt := range requiredBlockOptions[b.Kind] {
		if optionIsUnset(b.Options[opt]) {
			missing = append(missing, opt)
		}
	}

	if b.Kind == "Metric" && len(missing) == 0 {
		if tiles, ok := b.Options["metrics"].([]any); ok {
			for i, t := range tiles {
				tile, ok := t.(map[string]any)
				if !ok || optionIsUnset(tile["moduleID"]) {
					missing = append(missing, fmt.Sprintf("metrics.%d.moduleID", i))
				}
			}
		}
	}

	sort.Strings(missing)
	return missing
}

// optionIsUnset covers the three shapes an unconfigured option takes: absent,
// blank, and the "0" that stands for an unset ID (lib/js NoID). A numeric zero
// counts too — a JSON number arrives here as float64 and an ID sent that way
// is still no ID.
func optionIsUnset(v any) bool {
	switch t := v.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(t) == "" || t == "0"
	case float64:
		return t == 0
	case []any:
		return len(t) == 0
	case map[string]any:
		return len(t) == 0
	default:
		return false
	}
}
