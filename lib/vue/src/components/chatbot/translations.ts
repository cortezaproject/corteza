// Translation key catalog for CChatbotInbox.
//
// The component itself has English defaults baked in; consumers that want
// real i18n can build the full translations object by calling
// `makeChatbotInboxTranslations` with their vue-i18n `t` function and the
// path prefix at which the same set of keys exists in their locale tree.

export const CHATBOT_INBOX_TRANSLATION_KEYS = [
  'inbox',
  'sourcePreview',
  'draftChatbot',
  'authorUser',
  'authorAgent',
  'sectionHandoffRequested',
  'sectionHandoffActive',
  'sectionActive',
  'sectionClosed',
  'empty',
  'emptyNoChatbots',
  'selectHint',
  'accept',
  'resolve',
  'operator',
  'replyingAs',
  'changeAlias',
  'aliasLabel',
  'aliasHint',
  'aliasSave',
  'aliasReset',
  'composerPlaceholder',
  'composerDisabled',
  'confirmAcceptHeader',
  'confirmAcceptMessage',
  'confirmResolveHeader',
  'confirmResolveMessage',
  'adminActions',
  'adminAdvanceStep',
  'adminCloseSession',
  'confirmAdvanceStepHeader',
  'confirmAdvanceStepMessage',
  'confirmCloseSessionHeader',
  'confirmCloseSessionMessage',
  'statusHandoffRequested',
  'statusHandoffActive',
  'statusActive',
  'statusClosed',
  'filterTitle',
  'filterStatus',
  'filterChatbot',
  'filterSource',
  'filterSourceLive',
  'filterSourcePreview',
  'filterApply',
  'filterReset',
  'filterEmptyFiltered',
  'cancel',
] as const

export type ChatbotInboxTranslationKey = (typeof CHATBOT_INBOX_TRANSLATION_KEYS)[number]
export type ChatbotInboxTranslations = Partial<Record<ChatbotInboxTranslationKey, string>>

type Translator = (key: string) => string

// Build a translations object for CChatbotInbox by resolving each key under
// `prefix` against the supplied vue-i18n `t` function. The component has no
// built-in fallbacks, so the consumer's locale tree must define every key
// under `prefix` — anything missing renders as empty in the UI.
export function makeChatbotInboxTranslations(
  t: Translator,
  prefix = '',
): ChatbotInboxTranslations {
  const out: ChatbotInboxTranslations = {}
  for (const key of CHATBOT_INBOX_TRANSLATION_KEYS) {
    out[key] = t(`${prefix}${key}`)
  }
  return out
}
