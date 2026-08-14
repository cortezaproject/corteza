package inputguard

// injectionMarkers are literal tokens used by various LLM prompt formats.
// These are format-level markers and are language-agnostic.
var injectionMarkers = []string{
	// ChatML
	//
	// The bare fragments "<|" and "|>" are deliberately not listed: they are
	// substrings of the full tokens below, so they add no coverage, while "|>" is
	// the pipe operator in Elixir, F# and OCaml and both appear in ASCII art.
	"<|im_start|>",
	"<|im_end|>",

	// Llama / Llama 2
	"[INST]",
	"[/INST]",
	"<<SYS>>",
	"<</SYS>>",

	// Alpaca
	"### Instruction:",
	"### Response:",

	// Mistral
	"[AVAILABLE_TOOLS]",
	"[/AVAILABLE_TOOLS]",
}
