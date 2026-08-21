# Brief: `system_chatbot`

Tool file: `server/system/agentic/chatbot_tools.go`
Handler:   `server/system/agentic/chatbot_handler.go`
Group: `configuring`

## 1. Service

`system/service/chatbot.go` + `chatbot.gen.go` — `sysService.DefaultChatbot`.

| Method | Signature | Becomes |
|---|---|---|
| `FindByID` | `(ctx, ID uint64) (*types.Chatbot, error)` | single-fetch in `lookup` |
| `Search` | `(ctx, types.ChatbotFilter) (ChatbotSet, ChatbotFilter, error)` | list mode, and handle resolution |
| `Create` / `Update` | `(ctx, *types.Chatbot) (*types.Chatbot, error)` | `create` / `update` |
| `DeleteByID` | `(ctx, ID uint64) error` | `delete` |
| `UndeleteByID` | `(ctx, ID uint64) error` | `undelete` — generated, `chatbot.gen.go:181` |
| `RegenerateWidgetKey` | `(ctx, ID uint64) (*types.Chatbot, error)` | domain op, write |
| `UploadAsset` | attachment upload | **not a tool** — no file transport over MCP |

Risks: `lookup` read; `create`/`update`/`undelete`/`regenerate_key` write;
`delete` destructive.

## 2. Authorization (§8.6) — PASSES

`CanReadChatbot` (`chatbot.go:37`, and in `onSearch`'s `filter.Check` at
`:170`), `CanUpdateChatbot` (`:68`, and `:102` for the key regeneration),
`CanDeleteChatbot` (`:123`, `:160`).

**`chatbotSession` and `chatbotPreview` stay on the deny-list.**
`chatbot_session.go` `Search` goes straight to the store with no check of any
kind — its authorization lives in the REST controller, which a tool does not go
through. Do not write session tools.

## 3. Identifier strategy

**No `FindByAny`.** Same shape as `system_agent`: all-digits → `FindByID`,
otherwise `Search` with `ChatbotFilter{Handle: ref}` and exactly one hit.

A chatbot has both `Name` and `Handle`; only the handle is resolvable, because
`ChatbotFilter` has no name field.

`undelete` takes a required `chatbotID`.

## 4. Filter fields → params

`ChatbotFilter` (`system/types/chatbot.go`):

| Field | Param? | Note |
|---|---|---|
| `Query` | yes | |
| `Handle` | yes | |
| `WidgetKey` | **no** | see traps — it is a lookup BY embed key, and exposing it invites key-guessing |
| `ChatbotID` | no | the `chatbot` ref covers it |
| `ProjectID` | no | as with agents |
| `Deleted` | yes, as `includeDeleted` | |
| `Labels` | no | no label tools yet |

## 5. Undelete

`UndeleteByID` is generated and exists. Ships as `system_chatbot_undelete`.

## 6. Domain ops

`system_chatbot_regenerate_key` — rotates the widget key. Risk `write`. Its
description must say the previous key stops working the moment this returns, so
every page embedding the widget has to be updated.

## 7. Traps

**Update replaces the whole record**, exactly as for agents. Load first, merge
what the caller sent, send the whole thing back.

**Three nested config sections**: `handoff`, `styling`, `scenarios`. Each is
one JSON param replacing that section wholesale (§8.2). `scenarios` is a
collection, so it always replaces — there is no per-scenario merge.

**A `conversation` scenario without an `agentID` is rejected.**
`validateChatbotScenarios` (`chatbot.go:235`) enforces it on both create and
update, so the error arrives on save rather than on the offending scenario. Say
so in the description and name the field.

**The widget key is generated when absent and preserved when omitted.**
`prepareChatbotOnCreate` mints one at create; `prepareChatbotOnUpdate` copies
the existing one forward when the payload has none. It is therefore never a
create or update param — a caller cannot set it, only rotate it.

It is also not a secret in the §8.6b sense: it is embedded in the HTML of
whatever page hosts the widget, and `allowedOrigins` is what actually gates
use. It is returned by a single-chatbot lookup, matching REST, and left out of
the list projection — a listing has no reason to spray embed keys.

**`allowedOrigins` is the real access control** on a widget. A chatbot created
with none is reachable from any origin that knows the key, so the create
description must say to set it.

**`styling.logoAttachmentID`** refers to an attachment uploaded through REST
(`/chatbots/{id}/upload-asset`). MCP has no file transport, so a caller can
point at an existing attachment ID but cannot create one.
