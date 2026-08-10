package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// themeSetting is where the webapp keeps its palettes. Writing it is what
// triggers the stylesheet recompile — the settings service watches this name
// and regenerates the CSS, so nothing here has to.
const themeSetting = "ui.studio.themes"

// Declarations for these handlers are in theme_tools.go, in the same order.
type themeHandler struct {
	reg toolRegistrar
}

func ThemeHandler(reg toolRegistrar) *themeHandler {
	h := &themeHandler{reg: reg}
	h.register()
	return h
}

var (
	// The stored palette is a JSON string nested inside the setting's JSON, so
	// a colour value survives two rounds of encoding. Anything malformed is
	// therefore easy to write and hard to see, which is why the value is
	// validated here rather than passed through.
	hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

	themeIDs = []string{"general", "light", "dark"}
)

func (h *themeHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	themes, err := loadThemes(ctx)
	if err != nil {
		return nil, err
	}

	wanted := strings.ToLower(strings.TrimSpace(toolkit.Str(args, "theme")))
	if wanted != "" && !knownTheme(wanted) {
		return nil, fmt.Errorf("no such theme %q — it is one of %s", wanted, strings.Join(themeIDs, ", "))
	}

	out := make([]map[string]any, 0, len(themes))
	for _, t := range themes {
		id, _ := t["id"].(string)
		if wanted != "" && id != wanted {
			continue
		}

		colors, err := decodeColors(t)
		if err != nil {
			return nil, err
		}

		out = append(out, map[string]any{"theme": id, "colors": colors})
	}

	return toolkit.JSONResult(out)
}

func (h *themeHandler) update(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	wanted, err := toolkit.ReqStr(args, "theme")
	if err != nil {
		return nil, err
	}

	wanted = strings.ToLower(strings.TrimSpace(wanted))
	if !knownTheme(wanted) {
		return nil, fmt.Errorf("no such theme %q — it is one of %s", wanted, strings.Join(themeIDs, ", "))
	}

	var changes map[string]string
	present, err := toolkit.JSONArg(args, "colors", "colors", &changes)
	if err != nil {
		return nil, err
	}
	if !present || len(changes) == 0 {
		return nil, fmt.Errorf("colors is required: a JSON object of colour name to hex value")
	}

	for name, value := range changes {
		if !hexColor.MatchString(value) {
			return nil, fmt.Errorf("colour %q is %q, which is not a six-digit hex value like \"#09344E\"", name, value)
		}
	}

	themes, err := loadThemes(ctx)
	if err != nil {
		return nil, err
	}

	var (
		updated map[string]string
		found   bool
	)

	for _, t := range themes {
		if id, _ := t["id"].(string); id != wanted {
			continue
		}

		colors, err := decodeColors(t)
		if err != nil {
			return nil, err
		}

		// Merge rather than replace: a caller changing the primary colour is
		// not saying the other eleven should go back to their defaults.
		for name, value := range changes {
			colors[name] = value
		}

		encoded, err := json.Marshal(colors)
		if err != nil {
			return nil, toolkit.Errf("theme colours", err)
		}

		t["values"] = string(encoded)
		updated, found = colors, true
		break
	}

	if !found {
		return nil, fmt.Errorf("theme %q is not in %s on this instance", wanted, themeSetting)
	}

	// Marshalling the whole array back, rather than the one theme, is what
	// keeps the fields this tool does not model — "title" among them — instead
	// of dropping them by round-tripping through a narrower struct.
	value, err := json.Marshal(themes)
	if err != nil {
		return nil, toolkit.Errf("theme settings", err)
	}

	err = sysService.DefaultSettings.Set(ctx, &sysTypes.SettingValue{
		Name:      themeSetting,
		Value:     value,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, toolkit.Errf("theme settings", err)
	}

	changed := make([]string, 0, len(changes))
	for name := range changes {
		changed = append(changed, name)
	}
	sort.Strings(changed)

	return toolkit.JSONResult(map[string]any{
		"theme":   wanted,
		"changed": changed,
		"colors":  updated,
		"note":    "Applies to every user of this instance on their next page load.",
	})
}

// loadThemes reads the setting as free-form maps.
//
// types.Theme exists and would be the obvious target, but it models only id and
// values — the stored objects also carry "title", and unmarshalling into the
// struct to marshal it back would silently drop it.
func loadThemes(ctx context.Context) ([]map[string]any, error) {
	v, err := sysService.DefaultSettings.Get(ctx, themeSetting, 0)
	if err != nil {
		return nil, toolkit.Errf("theme settings", err)
	}
	if v == nil || v.IsNull() {
		return nil, fmt.Errorf("%s is not set on this instance, so there is no palette to read or change", themeSetting)
	}

	var themes []map[string]any
	if err = json.Unmarshal(v.Value, &themes); err != nil {
		return nil, toolkit.Errf("theme settings", err)
	}

	return themes, nil
}

// decodeColors unwraps the palette from the JSON string it is stored in.
func decodeColors(theme map[string]any) (map[string]string, error) {
	raw, _ := theme["values"].(string)
	if strings.TrimSpace(raw) == "" {
		return map[string]string{}, nil
	}

	colors := map[string]string{}
	if err := json.Unmarshal([]byte(raw), &colors); err != nil {
		id, _ := theme["id"].(string)
		return nil, fmt.Errorf("theme %q holds a palette that is not a JSON object: %w", id, err)
	}

	return colors, nil
}

func knownTheme(id string) bool {
	for _, known := range themeIDs {
		if known == id {
			return true
		}
	}
	return false
}
