package inputguard

// injectionMarkers are literal tokens used by various LLM prompt formats.
// These are format-level markers and are language-agnostic.
var injectionMarkers = []string{
	// ChatML
	"<|im_start|>",
	"<|im_end|>",
	"<|",
	"|>",

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
