<template>
  <PageBlock :block="block" @refreshBlock="refreshAll">
    <div class="flex flex-row h-full min-h-0 overflow-hidden">
      <!-- Sessions list -->
      <div class="w-72 shrink-0 border-r border-surface flex flex-col min-h-0">
        <div class="px-3 py-2 border-b border-surface flex items-center gap-2">
          <span class="font-medium text-sm flex-1 truncate">
            {{ $t('block.chatbotInbox.inbox') }}
          </span>
          <Tag :value="String(sessions.length)" severity="secondary" />
        </div>

        <div v-if="loading && !sessions.length" class="flex items-center justify-center p-4">
          <ProgressSpinner style="width: 24px; height: 24px" strokeWidth="4" />
        </div>

        <div v-else-if="!chatbotIDs.length" class="flex-1 flex items-center justify-center p-4 text-muted-color text-sm text-center">
          {{ $t('block.chatbotInbox.emptyNoChatbots') }}
        </div>

        <div v-else-if="!sessions.length" class="flex-1 flex items-center justify-center p-4 text-muted-color text-sm text-center">
          {{ $t('block.chatbotInbox.empty') }}
        </div>

        <div v-else class="flex-1 overflow-y-auto">
          <button
            v-for="s in sessions"
            :key="s.id"
            type="button"
            :class="[
              'w-full text-left px-3 py-2 border-b border-surface hover:bg-emphasis transition-colors',
              s.id === selectedID ? 'bg-emphasis' : '',
            ]"
            @click="selectSession(s)"
          >
            <div class="flex items-center gap-2">
              <span class="flex-1 truncate text-sm font-medium">
                {{ chatbotName(s.chatbotID) }}
              </span>
              <Tag
                :value="statusLabel(s.status)"
                :severity="statusSeverity(s.status)"
                rounded
              />
            </div>
            <div class="text-xs text-muted-color truncate">
              {{ formatTime(s.createdAt) }} · #{{ s.id }}
            </div>
          </button>
        </div>
      </div>

      <!-- Takeover panel -->
      <div class="flex-1 flex flex-col min-h-0 overflow-hidden">
        <div v-if="!selected" class="flex-1 flex items-center justify-center text-muted-color text-sm text-center p-4">
          {{ $t('block.chatbotInbox.selectHint') }}
        </div>

        <template v-else>
          <!-- Panel header -->
          <div class="px-3 py-2 border-b border-surface flex items-center gap-2 shrink-0">
            <div class="flex flex-col min-w-0 flex-1">
              <span class="font-medium text-sm truncate">
                {{ chatbotName(selected.chatbotID) }}
              </span>
              <span class="text-xs text-muted-color truncate">
                {{ $t('block.chatbotInbox.session') }} #{{ selected.id }}
              </span>
            </div>

            <Tag
              :value="statusLabel(selected.status)"
              :severity="statusSeverity(selected.status)"
              rounded
            />

            <Button
              v-if="canAccept"
              :label="$t('block.chatbotInbox.accept')"
              icon="pi pi-check"
              severity="success"
              size="small"
              :loading="acting"
              @click="onAccept"
            />
            <Button
              v-if="canResolve"
              :label="$t('block.chatbotInbox.resolve')"
              icon="pi pi-times"
              severity="secondary"
              size="small"
              outlined
              :loading="acting"
              @click="onResolve"
            />
          </div>

          <!-- Messages + composer -->
          <div class="flex-1 flex flex-col min-h-0">
            <div v-if="loadingMessages" class="flex-1 flex items-center justify-center">
              <ProgressSpinner style="width: 24px; height: 24px" strokeWidth="4" />
            </div>
            <CChatMessages
              v-else
              :messages="messages"
              :readonly="!canSend"
              :placeholder="composerPlaceholder"
              class="flex-1 min-h-0"
              @send="onSend"
            />
          </div>
        </template>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { components } from '@planetcrust/human-vue'
import { computed, inject, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const { CChatMessages } = components
const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI', null)

// PageBlockChatbotInbox in lib/js owns the defaults; options always arrives
// fully shaped from PageBlockMaker. Auto-refresh is handled by PageBlock.vue
// from block.options.refreshRate (seconds), same as other blocks.
const options = computed(() => props.block.options)
const chatbotIDs = computed(() => options.value.chatbotIDs)
const statusFilter = computed(() => options.value.statusFilter)
const autoOpenFirst = computed(() => options.value.autoOpenFirst)

const loading = ref(false)
const sessions = ref([])
const chatbotsByID = ref({})

const selectedID = ref(null)
const selected = computed(() => sessions.value.find(s => s.id === selectedID.value) || null)

const messages = ref([])
const loadingMessages = ref(false)
const acting = ref(false)
let stream = null

const canAccept = computed(() => selected.value?.status === 'handoff_requested')
const canSend = computed(() => selected.value?.status === 'handoff_active')
const canResolve = computed(() => ['handoff_requested', 'handoff_active'].includes(selected.value?.status))

const composerPlaceholder = computed(() =>
  canSend.value
    ? t('block.chatbotInbox.composerPlaceholder')
    : t('block.chatbotInbox.composerDisabled'),
)

async function loadChatbots() {
  if (!$SystemAPI || !chatbotIDs.value.length) return
  // Single batch list — handoff inbox typically watches a handful of chatbots.
  try {
    const res = await $SystemAPI.chatbotList({ limit: 0, sort: 'name ASC' })
    const map = {}
    for (const cb of res?.set || []) {
      map[cb.chatbotID] = cb
    }
    chatbotsByID.value = map
  } catch (err) {
    console.warn('[chatbot-inbox] chatbotList failed', err)
  }
}

async function fetchOnce() {
  if (!$SystemAPI || !chatbotIDs.value.length) {
    sessions.value = []
    return
  }
  loading.value = true
  try {
    const results = await Promise.all(
      chatbotIDs.value.map(chatbotID =>
        $SystemAPI
          .chatbotSessionListByChatbot({
            chatbotID,
            status: statusFilter.value,
            sort: 'createdAt DESC',
            limit: 100,
          })
          .catch(err => {
            console.warn('[chatbot-inbox] list failed', chatbotID, err)
            return { set: [] }
          }),
      ),
    )
    const merged = results.flatMap(r => r?.set || [])
    merged.sort((a, b) => (a.createdAt < b.createdAt ? 1 : -1))
    sessions.value = merged

    if (autoOpenFirst.value && !selectedID.value && merged.length) {
      void selectSession(merged[0])
    }
    // Drop selection if the current row dropped out of the inbox.
    if (selectedID.value && !merged.some(s => s.id === selectedID.value)) {
      clearSelection()
    }
  } finally {
    loading.value = false
  }
}

function refreshAll() {
  void loadChatbots()
  void fetchOnce()
}

function chatbotName(id) {
  const cb = chatbotsByID.value[id]
  return cb?.name || cb?.handle || id
}

function statusLabel(status) {
  if (!status) return ''
  return t(`block.chatbotInbox.status.${status}`)
}

function statusSeverity(status) {
  switch (status) {
    case 'handoff_requested':
      return 'warn'
    case 'handoff_active':
      return 'success'
    case 'closed':
      return 'secondary'
    default:
      return 'secondary'
  }
}

function formatTime(iso) {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
}

async function selectSession(s) {
  selectedID.value = s.id
  await loadSessionMessages(s.id)
  attachStream(s.id)
}

function clearSelection() {
  selectedID.value = null
  messages.value = []
  closeStream()
}

async function loadSessionMessages(sessionID) {
  loadingMessages.value = true
  messages.value = []
  try {
    const res = await $SystemAPI.chatbotSessionRead({ sessionID })
    const steps = (res?.steps || []).sort(
      (a, b) => (a.scenarioIndex ?? 0) - (b.scenarioIndex ?? 0),
    )
    const all = []
    for (const step of steps) {
      if (!step.conversationID || step.conversationID === '0') continue
      try {
        const conv = await $SystemAPI.aiConversationRead({
          aiConversationID: step.conversationID,
        })
        all.push(...(conv?.messages || []))
      } catch {
        /* missing conversation — skip */
      }
    }
    messages.value = all
  } catch (err) {
    console.error('[chatbot-inbox] read session failed', err)
  } finally {
    loadingMessages.value = false
  }
}

function pushMessage(msg) {
  messages.value = [...messages.value, msg]
}

// Build a base URL pointing at /system/chatbots/sessions/{id}/... using
// SystemAPI.baseURL so the SSE EventSource hits the same origin as the rest of
// the admin API.
function sessionRoot(sessionID) {
  const base = ($SystemAPI?.baseURL || '').replace(/\/+$/, '')
  return `${base}/chatbots/sessions/${encodeURIComponent(sessionID)}`
}

function authHeaders() {
  const tok = $SystemAPI?.accessTokenFn?.() || ''
  const h = { 'Content-Type': 'application/json' }
  if (tok) h.Authorization = `Bearer ${tok}`
  return h
}

async function postAction(path, body) {
  const r = await fetch(`${sessionRoot(selectedID.value)}${path}`, {
    method: 'POST',
    credentials: 'include',
    headers: authHeaders(),
    body: JSON.stringify(body || {}),
  })
  if (!r.ok && r.status !== 204) {
    const txt = await r.text().catch(() => '')
    throw new Error(`${r.status}: ${txt || r.statusText}`)
  }
}

function attachStream(sessionID) {
  closeStream()
  // EventSource can't set headers — query-token fallback handled by the
  // upstream HTTP token validator (?jwt=).
  const tok = $SystemAPI?.accessTokenFn?.() || ''
  const qs = tok ? `?jwt=${encodeURIComponent(tok)}` : ''
  stream = new EventSource(`${sessionRoot(sessionID)}/stream${qs}`, {
    withCredentials: true,
  })

  stream.addEventListener('user_message', ev => {
    const p = safeJSON(ev.data)
    if (p?.content) pushMessage({ role: 'user', content: p.content })
  })
  stream.addEventListener('operator_message', ev => {
    const p = safeJSON(ev.data)
    if (p?.content) {
      pushMessage({ role: 'assistant', content: `[operator] ${p.content}` })
    }
  })
  stream.addEventListener('token', ev => {
    const p = safeJSON(ev.data)
    const text = p?.text || p?.token || ''
    if (!text) return
    // Coalesce streaming agent tokens into the last assistant bubble.
    const last = messages.value[messages.value.length - 1]
    if (last && last.role === 'assistant' && last.streaming) {
      messages.value = [
        ...messages.value.slice(0, -1),
        { ...last, content: last.content + text },
      ]
    } else {
      pushMessage({ role: 'assistant', content: text, streaming: true })
    }
  })
  stream.addEventListener('done', () => {
    const last = messages.value[messages.value.length - 1]
    if (last?.streaming) {
      messages.value = [
        ...messages.value.slice(0, -1),
        { ...last, streaming: false },
      ]
    }
  })
  stream.addEventListener('handoff_active', () => {
    bumpSessionStatus('handoff_active')
  })
  stream.addEventListener('handoff_complete', () => {
    bumpSessionStatus('active')
  })
  stream.addEventListener('session_closed', () => {
    bumpSessionStatus('closed')
  })
  stream.onerror = err => {
    console.warn('[chatbot-inbox] SSE error', err)
  }
}

function closeStream() {
  if (stream) {
    stream.close()
    stream = null
  }
}

// Patch the session row in `sessions` so the badge + button state updates
// without waiting for the next poll cycle.
function bumpSessionStatus(status) {
  if (!selectedID.value) return
  sessions.value = sessions.value.map(s =>
    s.id === selectedID.value ? { ...s, status } : s,
  )
}

function safeJSON(s) {
  try {
    return JSON.parse(s || '{}')
  } catch {
    return null
  }
}

async function onAccept() {
  if (!selected.value || acting.value) return
  acting.value = true
  try {
    await postAction('/handoff-accept', {})
    bumpSessionStatus('handoff_active')
  } catch (err) {
    console.error('[chatbot-inbox] accept failed', err)
  } finally {
    acting.value = false
  }
}

async function onSend(text) {
  if (!text || !canSend.value) return
  try {
    await postAction('/operator-message', { message: text })
    // Optimistic echo; SSE will not deliver our own operator_message back to
    // this client (server emits to the conversation bus — also reaches us, but
    // visual delay would be noticeable otherwise).
    pushMessage({ role: 'assistant', content: `[operator] ${text}` })
  } catch (err) {
    console.error('[chatbot-inbox] send failed', err)
  }
}

async function onResolve() {
  if (!selected.value || acting.value) return
  acting.value = true
  try {
    await postAction('/handoff-complete', {})
    bumpSessionStatus('active')
  } catch (err) {
    console.error('[chatbot-inbox] resolve failed', err)
  } finally {
    acting.value = false
  }
}

// Initial + reactive refetch on filter changes. Periodic refresh comes from
// PageBlock's `refreshRate` watcher, which emits `refreshBlock` on interval and
// is wired to `refreshAll` on the root <PageBlock>.
watch(
  () => [chatbotIDs.value.join(','), statusFilter.value.join(',')],
  () => { refreshAll() },
  { immediate: true },
)

onBeforeUnmount(() => {
  closeStream()
})
</script>
