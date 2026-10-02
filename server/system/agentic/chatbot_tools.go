package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// This file holds declarations only. Implementations are in chatbot_handler.go,
// in the same order.

// A chatbot's configuration is three nested objects. Each is one parameter that
// REPLACES that whole section when sent and leaves it untouched when omitted
// (CONVENTIONS.md §8.2).
const chatbotSectionDoc = "A chatbot's configuration is three objects — handoff, styling and scenarios. " +
	"Send only the ones you are changing; each REPLACES that whole section, so read the chatbot " +
	"first and send an edited copy rather than a fragment."

const (
	chatbotHandoffDoc = `JSON object: passing the conversation to a person. ` +
		`{"enabled":true,"automation":{"onRequested":{"taqID":"<id>"},"onAccepted":{"taqID":"<id>"}}}. ` +
		`The hooks fire when a visitor asks for a human and when one picks the conversation up.`

	chatbotStylingDoc = `JSON object: how the widget looks. ` +
		`{"logoAttachmentID":"<id>","logoURL":"https://…","fontFamily":"Inter","fontSizes":{},"colors":{},"launcher":{}}. ` +
		`An attachment must already exist — it is uploaded through the REST endpoint /chatbots/{id}/upload-asset, and MCP has no way to transfer a file.`

	chatbotScenariosDoc = `JSON array: what the chatbot can do, in order. ` +
		`[{"id":"triage","name":"Triage","type":"conversation","agentID":"<agentID>","automation":{"before":{},"after":{}}}]. ` +
		`A scenario of type "conversation" MUST carry an agentID — the whole save is rejected otherwise, not just that scenario. ` +
		`Find the agent with system_agent_lookup. This is a collection and always replaces: send every scenario the chatbot should have, not only the new one.`
)

func (h *chatbotHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_chatbot_lookup",
			mcp.WithDescription(
				"Read the chatbots configured on this instance — a chatbot is the embeddable widget and "+
					"the conversation around it, and it delegates the thinking to an agent through a "+
					"scenario. "+
					"Provide 'chatbot' to fetch one whole, with its handoff, styling and scenarios; omit "+
					"it to list them as {chatbotID, handle, name, enabled}. "+
					"Call this before system_chatbot_update: an update replaces whole sections. "+
					"A chatbot is not an agent — see system_agent_lookup for the model and prompt behind "+
					"it. Conversation transcripts are deliberately not reachable here.",
			),
			mcp.WithString("chatbot", mcp.Description("Chatbot handle, or ID as a string to prevent precision loss. Omit to list. A chatbot's name is not resolvable — only its handle is.")),
			mcp.WithString("query", mcp.Description("Free-text search across the chatbots.")),
			mcp.WithString("handle", mcp.Description("Filter the listing to an exact handle.")),
			mcp.WithBoolean("includeDeleted", mcp.Description("Include soft-deleted chatbots in the listing. They are excluded by default, which is why undelete needs an ID.")),
			mcp.WithString("limit", mcp.Description("Maximum chatbots to return when listing, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup chatbot",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_chatbot_create",
			mcp.WithDescription(
				"Create a chatbot. "+chatbotSectionDoc+" "+
					"Two things decide whether it does anything. It needs a 'conversation' scenario "+
					"naming an agentID (from system_agent_lookup), or it has nothing to answer with. And it needs "+
					"'allowedOrigins', because that list is what actually restricts where the widget may "+
					"run — a chatbot with none is usable from any site that has its widget key. "+
					"The widget key is generated here and returned; you cannot choose it.",
			),
			mcp.WithString("name", mcp.Required(), mcp.Description("Human-readable name, shown in the admin list.")),
			mcp.WithString("handle", mcp.Description("URL-friendly identifier, unique across chatbots. Use snake_case (lowercase letters, digits, underscores).")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the widget answers. A disabled chatbot keeps its configuration and serves nobody.")),
			mcp.WithString("allowedOrigins", mcp.Description(`JSON array of origins the widget may be embedded on, e.g. ["https://example.com"]. This is the access control on the widget, not decoration: an empty list restricts nothing.`)),
			mcp.WithString("sessionTTL", mcp.Description(`How long a visitor's session survives, as a duration string, e.g. "30m".`)),
			mcp.WithString("handoff", mcp.Description(chatbotHandoffDoc)),
			mcp.WithString("styling", mcp.Description(chatbotStylingDoc)),
			mcp.WithString("scenarios", mcp.Description(chatbotScenariosDoc)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
			hmcp.NeedsFullDocs(),
		),
		"Create chatbot",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_chatbot_update",
			mcp.WithDescription(
				"Change an existing chatbot. "+chatbotSectionDoc+" "+
					"An argument you omit is left unchanged, so this is safe to call with one section. "+
					"'scenarios' is a collection and always replaces: send every scenario the chatbot "+
					"should end up with, because sending only the new one deletes the rest. Call "+
					"system_chatbot_lookup with 'chatbot' first. "+
					"A 'conversation' scenario missing its agentID rejects the whole call. "+
					"The widget key cannot be set here — it is carried forward untouched. Rotate it with "+
					"system_chatbot_regenerate_key.",
			),
			mcp.WithString("chatbot", mcp.Required(), mcp.Description("Chatbot handle, or ID as a string to prevent precision loss.")),
			mcp.WithString("name", mcp.Description("New name. An empty string is ignored — a chatbot without one is unidentifiable in the admin list.")),
			mcp.WithString("handle", mcp.Description("New handle. Pass an empty string to clear it.")),
			mcp.WithBoolean("enabled", mcp.Description("Whether the widget answers.")),
			mcp.WithString("allowedOrigins", mcp.Description(`JSON array of origins the widget may be embedded on. Replaces the list wholesale; an empty array removes every restriction.`)),
			mcp.WithString("sessionTTL", mcp.Description(`How long a visitor's session survives, as a duration string, e.g. "30m".`)),
			mcp.WithString("handoff", mcp.Description(chatbotHandoffDoc)),
			mcp.WithString("styling", mcp.Description(chatbotStylingDoc)),
			mcp.WithString("scenarios", mcp.Description(chatbotScenariosDoc)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
			hmcp.NeedsFullDocs(),
		),
		"Update chatbot",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_chatbot_delete",
			mcp.WithDescription(
				"Delete a chatbot. The delete is soft: the configuration is retained and "+
					"system_chatbot_undelete brings it back, but the widget stops answering everywhere "+
					"it is embedded. "+
					"Note the chatbotID before calling — a deleted chatbot is excluded from every "+
					"listing, so its handle no longer resolves and undelete needs the ID. "+
					"The agent a scenario pointed at is untouched; deleting the widget does not delete "+
					"what answered through it.",
			),
			mcp.WithString("chatbot", mcp.Required(), mcp.Description("Chatbot handle, or ID as a string to prevent precision loss.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete chatbot",
		h.del,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_chatbot_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted chatbot, reversing system_chatbot_delete. Nothing was erased, so "+
					"it returns with its scenarios, styling and the same widget key — pages already "+
					"embedding it start working again without being changed. "+
					"Requires the numeric chatbotID: a deleted chatbot is excluded from the default "+
					"listing, so its handle resolves to nothing. Use the ID system_chatbot_delete "+
					"reported, or call system_chatbot_lookup with includeDeleted first. "+
					"Calling this on a chatbot that is not deleted is accepted and changes nothing.",
			),
			mcp.WithString("chatbotID", mcp.Required(), mcp.Description("ID of the deleted chatbot, as a string to prevent precision loss. A handle will not work.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete chatbot",
		h.undelete,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_chatbot_regenerate_key",
			mcp.WithDescription(
				"Issue the chatbot a new widget key. The previous key stops working the moment this "+
					"returns, so every page embedding the widget breaks until its embed code is updated "+
					"with the new key — which this returns. "+
					"Use it when a key has leaked somewhere it should not be, or when retiring an embed "+
					"you no longer control. It is not part of ordinary editing: nothing else about the "+
					"chatbot changes, and an update never rotates the key on its own. "+
					"The key identifies the widget rather than authorising it — 'allowedOrigins' is what "+
					"restricts where it may run, so rotating a key on a chatbot that allows every origin "+
					"buys very little.",
			),
			mcp.WithString("chatbot", mcp.Required(), mcp.Description("Chatbot handle, or ID as a string to prevent precision loss.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Regenerate chatbot widget key",
		h.regenerateKey,
	)
}
