package inputguard

// injectionMarkers are literal tokens used by various LLM prompt formats.
// These are format-level markers and are language-agnostic.
var injectionMarkers = []string{
	// ChatML
	//
	// The bare two-character fragments "<|" and "|>" used to be listed here as
	// well, and they blocked ordinary text: "|>" is the pipe operator in
	// Elixir, F# and OCaml, and both appear in ASCII art and diagrams. They
	// also added no coverage, since a real ChatML delimiter contains one of the
	// full tokens below.
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
