<template>
  <PageBlock :block="block" @refreshBlock="refreshAll">
    <div class="flex flex-row flex-1 min-h-0 overflow-hidden">
      <!-- Sessions list -->
      <div class="w-96 shrink-0 border-r border-surface flex flex-col min-h-0">
        <div class="p-3 border-b border-surface flex items-center gap-2 shrink-0">
          <span class="font-medium text-sm flex-1 truncate">
            {{ $t('block.chatbotInbox.inbox') }}
          </span>
          <Tag :value="String(sessions.length)" severity="secondary" :pt="tagSmallPt" />
        </div>

        <div v-if="loading && !sessions.length" class="flex items-center justify-center p-4">
          <ProgressSpinner style="width: 24px; height: 24px" strokeWidth="4" />
        </div>

        <div
          v-else-if="!chatbotIDs.length"
          class="flex-1 flex items-center justify-center p-4 text-muted-color text-sm text-center"
        >
          {{ $t('block.chatbotInbox.emptyNoChatbots') }}
        </div>

        <div
          v-else-if="!sessions.length"
          class="flex-1 flex items-center justify-center p-4 text-muted-color text-sm text-center"
        >
          {{ $t('block.chatbotInbox.empty') }}
        </div>

        <div v-else class="flex-1 overflow-y-auto p-2 flex flex-col gap-2">
          <button
            v-for="s in sessions"
            :key="s.id"
            type="button"
            :class="[
              'w-full text-left p-3 border border-surface rounded-border shadow-sm hover:bg-emphasis transition-colors',
              s.id === selectedID ? 'ring-2 ring-primary bg-emphasis' : '',
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
                :pt="tagSmallPt"
                rounded
              />
            </div>
            <div class="text-xs text-muted-color truncate mt-1 text-right">
              {{ formatTime(s.createdAt) }}
            </div>
          </button>
        </div>
      </div>

      <!-- Takeover panel -->
      <div class="flex-1 flex flex-col min-h-0 overflow-hidden">
        <div
          v-if="!selected"
          class="flex-1 flex items-center justify-center text-muted-color text-sm text-center p-4"
        >
          {{ $t('block.chatbotInbox.selectHint') }}
        </div>

        <template v-else>
          <!-- Panel header -->
          <div class="px-3 py-2 border-b border-surface flex items-center gap-2 shrink-0">
            <div class="flex flex-col min-w-0 flex-1">
              <div class="flex items-center gap-2 min-w-0">
                <span class="font-medium text-sm truncate">
                  {{ chatbotName(selected.chatbotID) }}
                </span>
                <Tag
                  :value="statusLabel(selected.status)"
                  :severity="statusSeverity(selected.status)"
                  :pt="tagSmallPt"
                  rounded
                  class="shrink-0"
                />
              </div>
              <button
                type="button"
                class="text-xs text-muted-color hover:text-primary text-left flex items-center gap-1"
                :title="$t('block.chatbotInbox.changeAlias')"
                @click="aliasMenu?.toggle($event)"
              >
                {{ $t('block.chatbotInbox.replyingAs', { name: operatorName }) }}
                <i class="pi pi-pencil text-[10px]" />
              </button>
            </div>

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

          <Popover ref="aliasMenu">
            <div class="flex flex-col gap-2 w-64 p-1">
              <label class="text-xs font-medium text-primary">
                {{ $t('block.chatbotInbox.aliasLabel') }}
              </label>
              <InputText
                v-model="aliasDraft"
                :placeholder="defaultOperatorName"
                size="small"
                @keydown.enter.prevent="saveAlias"
              />
              <small class="text-muted-color">
                {{ $t('block.chatbotInbox.aliasHint') }}
              </small>
              <div class="flex justify-end gap-1">
                <Button
                  :label="$t('block.chatbotInbox.aliasReset')"
                  size="small"
                  severity="secondary"
                  text
                  @click="resetAlias"
                />
                <Button
                  :label="$t('block.chatbotInbox.aliasSave')"
                  size="small"
                  @click="saveAlias"
                />
              </div>
            </div>
          </Popover>

          <!-- Messages + composer -->
          <div class="flex-1 flex flex-col min-h-0">
            <div v-if="loadingMessages" class="flex-1 flex items-center justify-center">
              <ProgressSpinner style="width: 24px; height: 24px" strokeWidth="4" />
            </div>
            <template v-else>
              <div
                ref="messagesScrollRef"
                class="flex-1 p-4 overflow-y-auto flex flex-col gap-4 min-h-0"
              >
                <div
                  v-for="(msg, idx) in visibleMessages"
                  :key="idx"
                  class="flex flex-col gap-1"
                  :class="msg.role === 'user' ? 'items-end' : 'items-start'"
                >
                  <span v-if="msg.authorName" class="text-xs font-semibold text-primary px-1">
                    {{ msg.authorName }}
                  </span>
                  <div
                    :class="[
                      'p-3 rounded-xl max-w-[85%] text-sm shadow-sm',
                      msg.role === 'user'
                        ? 'bg-primary text-primary-contrast whitespace-pre-wrap'
                        : 'bg-emphasis text-color',
                    ]"
                  >
                    <div
                      v-if="msg.role !== 'user'"
                      class="rt-content"
                      v-html="renderMarkdown(msg.content)"
                    />
                    <template v-else>{{ msg.content }}</template>
                  </div>
                </div>
              </div>

              <div class="p-3 border-t border-surface flex gap-2 shrink-0 bg-surface">
                <InputText
                  v-model="composerInput"
                  :placeholder="composerPlaceholder"
                  :disabled="!canSend"
                  class="flex-1"
                  @keydown.enter.exact.prevent="submitComposer"
                />
                <Button
                  icon="pi pi-send"
                  :disabled="!canSend || !composerInput.trim()"
                  @click="submitComposer"
                />
              </div>
            </template>
          </div>
        </template>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { renderMarkdown } from '@planetcrust/human-vue'
import { computed, inject, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { useConfirm } from 'primevue/useconfirm'
import { useI18n } from 'vue-i18n'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI', null)
const $Auth = inject('$Auth', null)
const confirm = useConfirm()

// Surface the `operator` field on assistant messages so the renderer can show
// it as the speaker label above the bubble.
function decorateMessages(list) {
  return list.map(m =>
    m.role === 'assistant' && m.operator ? { ...m, authorName: m.operator } : m,
  )
}

const visibleMessages = computed(() => decorateMessages(messages.value))

// Compact tag styling. PrimeVue Tag has no `size` prop, so override padding +
// font size via passthrough on the root element.
const tagSmallPt = {
  root: { class: '!text-[11px] !leading-none !py-1 !px-2' },
}

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
const composerInput = ref('')
const messagesScrollRef = ref(null)
const aliasMenu = ref(null)
const aliasDraft = ref('')
let stream = null

// Operator alias persists per signed-in user via localStorage. Server can be
// told a display label via the optional `operator` body field on actions;
// when unset the server falls back to the user's account name.
// Default: first token of the full name (handle/email used as-is when name is
// missing). Visitors typically see "Maja" rather than "Maja Novak".
const defaultOperatorName = computed(() => {
  const u = $Auth?.user
  const name = u?.name?.trim().split(/\s+/)[0]
  return name || u?.handle || u?.email || t('block.chatbotInbox.operator')
})

const aliasStorageKey = computed(() => {
  const u = $Auth?.user
  return u?.userID ? `chatbot-inbox:alias:${u.userID}` : ''
})

const operatorAlias = ref(readAlias())

const operatorName = computed(() => operatorAlias.value?.trim() || defaultOperatorName.value)

function readAlias() {
  const key = aliasStorageKey.value
  if (!key) return ''
  try {
    return localStorage.getItem(key) || ''
  } catch {
    return ''
  }
}

function persistAlias(v) {
  const key = aliasStorageKey.value
  if (!key) return
  try {
    if (v) localStorage.setItem(key, v)
    else localStorage.removeItem(key)
  } catch {
    /* private mode / quota: ignore — alias only survives this session */
  }
}

function saveAlias() {
  const v = aliasDraft.value.trim()
  operatorAlias.value = v
  persistAlias(v)
  aliasMenu.value?.hide()
}

function resetAlias() {
  aliasDraft.value = ''
  operatorAlias.value = ''
  persistAlias('')
  aliasMenu.value?.hide()
}

watch(
  operatorAlias,
  v => {
    aliasDraft.value = v
  },
  { immediate: true },
)

watch(
  () => visibleMessages.value.length,
  () => {
    void nextTick(() => {
      const el = messagesScrollRef.value
      if (el) el.scrollTop = el.scrollHeight
    })
  },
)

function submitComposer() {
  const txt = composerInput.value.trim()
  if (!txt || !canSend.value) return
  composerInput.value = ''
  void onSend(txt)
}

const canAccept = computed(() => selected.value?.status === 'handoff_requested')
const canSend = computed(() => selected.value?.status === 'handoff_active')
const canResolve = computed(() =>
  ['handoff_requested', 'handoff_active'].includes(selected.value?.status),
)

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
    const steps = (res?.steps || []).sort((a, b) => (a.scenarioIndex ?? 0) - (b.scenarioIndex ?? 0))
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
      pushMessage({
        role: 'assistant',
        operator: p.operator || '',
        content: p.content,
      })
    }
  })
  stream.addEventListener('token', ev => {
    const p = safeJSON(ev.data)
    const text = p?.text || p?.token || ''
    if (!text) return
    // Coalesce streaming agent tokens into the last assistant bubble.
    const last = messages.value[messages.value.length - 1]
    if (last && last.role === 'assistant' && last.streaming) {
      messages.value = [...messages.value.slice(0, -1), { ...last, content: last.content + text }]
    } else {
      pushMessage({ role: 'assistant', content: text, streaming: true })
    }
  })
  stream.addEventListener('done', () => {
    const last = messages.value[messages.value.length - 1]
    if (last?.streaming) {
      messages.value = [...messages.value.slice(0, -1), { ...last, streaming: false }]
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
  sessions.value = sessions.value.map(s => (s.id === selectedID.value ? { ...s, status } : s))
}

function safeJSON(s) {
  try {
    return JSON.parse(s || '{}')
  } catch {
    return null
  }
}

function onAccept() {
  if (!selected.value || acting.value) return
  confirm.require({
    message: t('block.chatbotInbox.confirmAccept.message'),
    header: t('block.chatbotInbox.confirmAccept.header'),
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      outlined: true,
      size: 'small',
    },
    acceptProps: {
      label: t('block.chatbotInbox.accept'),
      severity: 'success',
      size: 'small',
    },
    accept: () => void doAccept(),
  })
}

async function doAccept() {
  acting.value = true
  try {
    await postAction('/handoff-accept', { operator: operatorName.value })
    bumpSessionStatus('handoff_active')
  } catch (err) {
    console.error('[chatbot-inbox] accept failed', err)
  } finally {
    acting.value = false
  }
}

async function onSend(text) {
  if (!text || !canSend.value) return
  const opName = operatorName.value
  try {
    await postAction('/operator-message', { message: text, operator: opName })
    // Optimistic echo; SSE will eventually deliver our own operator_message
    // back to this client, but waiting for the round-trip looks laggy.
    pushMessage({ role: 'assistant', operator: opName, content: text })
  } catch (err) {
    console.error('[chatbot-inbox] send failed', err)
  }
}

function onResolve() {
  if (!selected.value || acting.value) return
  confirm.require({
    message: t('block.chatbotInbox.confirmResolve.message'),
    header: t('block.chatbotInbox.confirmResolve.header'),
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      outlined: true,
      size: 'small',
    },
    acceptProps: {
      label: t('block.chatbotInbox.resolve'),
      severity: 'warn',
      size: 'small',
    },
    accept: () => void doResolve(),
  })
}

async function doResolve() {
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
  () => {
    refreshAll()
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  closeStream()
})
</script>
