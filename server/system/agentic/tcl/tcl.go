package tcl

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed tcl.json
var rawMasterList []byte

type (
	MasterList struct {
		Version     string      `json:"version"`
		Treaties    []Treaty    `json:"treaties"`
		Articles    []Group     `json:"articles"`
		SelectionUI SelectionUI `json:"selectionUi"`
	}

	Treaty struct {
		ID              string `json:"id"`
		Label           string `json:"label"`
		Citation        string `json:"citation"`
		Jurisdiction    string `json:"jurisdiction"`
		Required        bool   `json:"required"`
		CanDeselect     bool   `json:"canDeselect"`
		DefaultSelected bool   `json:"defaultSelected"`
		SourceURL       string `json:"sourceUrl"`
	}

	Group struct {
		TreatyID string    `json:"treatyId"`
		Items    []Article `json:"items"`
	}

	Article struct {
		ID              string `json:"id"`
		Label           string `json:"label"`
		Category        string `json:"category"`
		Hardwired       bool   `json:"hardwired"`
		DefaultSelected bool   `json:"defaultSelected"`
		ReferenceURL    string `json:"referenceUrl"`
		Interpretation  string `json:"interpretation"`
	}

	SelectionUI struct {
		MultiSelect MultiSelectConfig `json:"multiSelect"`
		LLMHelper   LLMHelperConfig   `json:"llmHelper"`
	}

	MultiSelectConfig struct {
		TreatyFieldID              string `json:"treatyFieldId"`
		ArticleFieldID             string `json:"articleFieldId"`
		AllowArticleFilterByTreaty bool   `json:"allowArticleFilterByTreaty"`
		ShowHardwiredBadge         bool   `json:"showHardwiredBadge"`
		DisallowDeselectHardwired  bool   `json:"disallowDeselectHardwired"`
	}

	LLMHelperConfig struct {
		PromptTemplate       string                `json:"promptTemplate"`
		AutoSelectHeuristics []AutoSelectHeuristic `json:"autoSelectHeuristics"`
	}

	AutoSelectHeuristic struct {
		Condition          string   `json:"condition"`
		AutoSelectArticles []string `json:"autoSelectArticles"`
	}
)

var (
	master       MasterList
	articleIndex map[string]Article
	hardwiredIDs []string
	defaultIDs   []string
)

func init() {
	if err := json.Unmarshal(rawMasterList, &master); err != nil {
		panic(fmt.Sprintf("tcl: failed to parse master list: %v", err))
	}

	articleIndex = make(map[string]Article)
	for _, g := range master.Articles {
		for _, a := range g.Items {
			articleIndex[a.ID] = a
			if a.Hardwired {
				hardwiredIDs = append(hardwiredIDs, a.ID)
			}
			if a.DefaultSelected {
				defaultIDs = append(defaultIDs, a.ID)
			}
		}
	}
}

// Master returns the full embedded TCL master list.
func Master() MasterList {
	return master
}

// DefaultArticleIDs returns the IDs of all defaultSelected articles.
func DefaultArticleIDs() []string {
	return defaultIDs
}

// MergeWithHardwired takes a user's article selection and ensures all hardwired
// articles are present. Order is preserved: user selections first, then any
// hardwired articles not already in the list.
func MergeWithHardwired(selected []string) []string {
	seen := make(map[string]bool, len(selected))
	result := make([]string, 0, len(selected)+len(hardwiredIDs))

	for _, id := range selected {
		if !seen[id] {
			result = append(result, id)
			seen[id] = true
		}
	}
	for _, id := range hardwiredIDs {
		if !seen[id] {
			result = append(result, id)
			seen[id] = true
		}
	}
	return result
}

// BuildPrompt assembles the compliance system prompt block from the given
// article IDs and TCL temperature (1–10).
func BuildPrompt(articleIDs []string, temperature int) string {
	if len(articleIDs) == 0 {
		return ""
	}

	var articles []Article
	seen := make(map[string]bool)
	for _, id := range articleIDs {
		if a, ok := articleIndex[id]; ok && !seen[id] {
			articles = append(articles, a)
			seen[id] = true
		}
	}
	if len(articles) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString("<compliance>\n")
	sb.WriteString("1. Any questions you ask or answer, or any conversation you engage in must be compliant with the list of treaties and their articles below.\n")
	sb.WriteString("2. If you have any doubts, then you should reference the source information URLs which are included in their list of treaties and articles below before exercising your discretion as to which response you should give.\n")
	sb.WriteString("3. You may choose to escalate to a person should you feel the scenario requires it. You should inform the interviewee/user that this will be the case for the question/scenario and continue the interview or conversation.\n")
	sb.WriteString("4. You do not need to tell the interviewee/user every compliance treaty related reason for your questions and answers unless they specifically ask for the reason. Otherwise, you may use your discretion when you think telling the interviewee is particularly important.\n\n")

	fmt.Fprintf(&sb, "Very important: You must take the Treaty Compliance Layer into account in every conversational exchange. %d is where you cite all references of Treaty Compliance in your conversation. 0 is where you cite zero references of Treaty Compliance in your conversation. However, under all circumstances your conversations should be compliant, regardless of whether you cite references or not.\n\n", temperature)

	sb.WriteString("Critical: If your definition or behaviour as an Agent falls into EU AI Act's high-risk classification, then you must notify clearly that a Fundamental Rights Impact Assessment or Data Protection Impact Assessment is required and guarantee that this notification is escalated to a human being.\n\n")

	sb.WriteString("Applicable compliance articles:\n")
	for _, a := range articles {
		fmt.Fprintf(&sb, "- %s: %s (Reference: %s)\n", a.Label, a.Interpretation, a.ReferenceURL)
	}

	sb.WriteString("</compliance>")

	return sb.String()
}
