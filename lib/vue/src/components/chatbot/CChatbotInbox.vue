<template>
  <div class="flex flex-row flex-1 min-h-0 overflow-hidden">
    <!-- Sessions list -->
    <div class="w-96 shrink-0 border-r border-surface flex flex-col min-h-0">
      <div class="p-3 border-b border-surface flex items-center gap-2 shrink-0">
        <span class="font-medium text-sm truncate">{{ l('inbox') }}</span>
        <Tag :value="String(filteredSessions.length)" severity="secondary" :pt="tagSmallPt" />
        <Button
          v-if="showFilter"
          v-tooltip.bottom="l('filterTitle')"
          icon="pi pi-filter"
          :class="['ml-auto', { 'text-primary': hasActiveFilter }]"
          severity="secondary"
          size="small"
          text
          rounded
          @click="filterMenu?.toggle($event)"
        />
      </div>

      <Popover ref="filterMenu">
        <div class="flex flex-col gap-3 w-72 p-1">
          <div class="flex flex-col gap-1">
            <label class="text-xs font-medium text-primary">{{ l('filterStatus') }}</label>
            <div class="flex flex-col gap-1">
              <div v-for="s in filterStatusChoices" :key="s.value" class="flex items-center gap-2">
                <Checkbox
                  v-model="userFilter.statuses"
                  :input-id="`cb-inbox-flt-st-${s.value}`"
                  :value="s.value"
                />
                <label :for="`cb-inbox-flt-st-${s.value}`" class="text-sm">{{ s.label }}</label>
              </div>
            </div>
          </div>

          <div v-if="filterChatbotChoices.length > 1" class="flex flex-col gap-1">
            <label class="text-xs font-medium text-primary">{{ l('filterChatbot') }}</label>
            <div class="flex flex-col gap-1 max-h-32 overflow-y-auto">
              <div v-for="c in filterChatbotChoices" :key="c.value" class="flex items-center gap-2">
                <Checkbox
                  v-model="userFilter.chatbotIDs"
                  :input-id="`cb-inbox-flt-cb-${c.value}`"
                  :value="c.value"
                />
                <label :for="`cb-inbox-flt-cb-${c.value}`" class="text-sm truncate">
                  {{ c.label }}
                </label>
              </div>
            </div>
          </div>

          <div class="flex flex-col gap-1">
            <label class="text-xs font-medium text-primary">{{ l('filterSource') }}</label>
            <div class="flex flex-wrap gap-3">
              <div v-for="s in filterSourceChoices" :key="s.value" class="flex items-center gap-2">
                <Checkbox
                  v-model="userFilter.sources"
                  :input-id="`cb-inbox-flt-src-${s.value}`"
                  :value="s.value"
                />
                <label :for="`cb-inbox-flt-src-${s.value}`" class="text-sm">{{ s.label }}</label>
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-1">
            <Button
              :label="l('filterReset')"
              size="small"
              severity="secondary"
              text
              :disabled="!hasActiveFilter"
              @click="resetFilter"
            />
            <Button :label="l('filterApply')" size="small" @click="filterMenu?.hide()" />
          </div>
        </div>
      </Popover>

      <div v-if="loading && !sessions.length" class="flex items-center justify-center p-4">
        <ProgressSpinner style="width: 24px; height: 24px" strokeWidth="4" />
      </div>

      <div
        v-else-if="!chatbotIDs.length"
        class="flex-1 flex items-center justify-center p-4 text-muted-color text-sm text-center"
      >
        {{ l('emptyNoChatbots') }}
      </div>

      <div
        v-else-if="!filteredSessions.length"
        class="flex-1 flex items-center justify-center p-4 text-muted-color text-sm text-center"
      >
        {{ hasActiveFilter ? l('filterEmptyFiltered') : l('empty') }}
      </div>

      <div v-else class="flex-1 overflow-y-auto p-2 flex flex-col gap-3">
        <div v-for="g in sessionGroups" :key="g.key" class="flex flex-col gap-2">
          <button
            type="button"
            class="w-full flex items-center gap-2 px-1 py-1 text-xs uppercase tracking-wide text-muted-color font-medium hover:text-primary transition-colors"
            @click="toggleSection(g.key)"
          >
            <i
              class="pi text-[10px]"
              :class="isCollapsed(g.key) ? 'pi-chevron-right' : 'pi-chevron-down'"
            />
            <span class="flex-1 text-left">{{ g.label }}</span>
            <span class="text-xs normal-case opacity-70">{{ g.items.length }}</span>
          </button>

          <div v-show="!isCollapsed(g.key)" class="flex flex-col gap-2">
            <button
              v-for="s in g.items"
              :key="s.id"
              type="button"
              :class="[
                'w-full text-left p-3 border border-surface rounded-border shadow-sm hover:bg-emphasis transition-colors',
                s.id === selectedID ? 'ring-2 ring-primary bg-emphasis' : '',
              ]"
              @click="selectSession(s)"
            >
              <div class="flex items-center gap-2">
                <span class="flex-1 truncate text-sm font-medium">{{ chatbotName(s) }}</span>
                <Tag
                  v-if="s.source === 'preview'"
                  :value="l('sourcePreview')"
                  severity="info"
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
      </div>
    </div>

    <!-- Takeover panel -->
    <div class="flex-1 flex flex-col min-h-0 overflow-hidden">
      <div
        v-if="!selected"
        class="flex-1 flex items-center justify-center text-muted-color text-sm text-center p-4"
      >
        {{ l('selectHint') }}
      </div>

      <template v-else>
        <div class="px-3 py-2 border-b border-surface flex items-center gap-2 shrink-0">
          <div class="flex flex-col gap-2 min-w-0 flex-1">
            <div class="flex items-center gap-2 min-w-0">
              <span class="font-medium text-sm truncate">{{ chatbotName(selected) }}</span>
              <Tag
                v-if="selected.source === 'preview'"
                :value="l('sourcePreview')"
                severity="info"
                :pt="tagSmallPt"
                rounded
                class="shrink-0"
              />
              <Tag
                v-if="selected.status !== 'active'"
                :value="statusLabel(selected.status)"
                :severity="statusSeverity(selected.status)"
                :pt="tagSmallPt"
                rounded
                class="shrink-0"
              />
            </div>
            <button
              v-if="canManageHandoff"
              type="button"
              class="text-xs text-muted-color hover:text-primary text-left flex items-center gap-1"
              :title="l('changeAlias')"
              @click="aliasMenu?.toggle($event)"
            >
              {{ l('replyingAs', { name: operatorName }) }}
              <i class="pi pi-pencil text-[10px]" />
            </button>
          </div>

          <Button
            v-if="canAccept"
            :label="l('accept')"
            icon="pi pi-check"
            severity="success"
            size="small"
            :loading="acting"
            @click="onAccept"
          />
          <Button
            v-if="canResolve"
            :label="l('resolve')"
            icon="pi pi-times"
            severity="secondary"
            size="small"
            outlined
            :loading="acting"
            @click="onResolve"
          />

          <Button
            v-if="canManage"
            v-tooltip.bottom="l('adminActions')"
            icon="pi pi-ellipsis-v"
            severity="secondary"
            size="small"
            text
            :disabled="acting"
            @click="adminMenu?.toggle($event)"
          />
          <Menu ref="adminMenu" :model="adminMenuItems" :popup="true" />
        </div>

        <Popover ref="aliasMenu">
          <div class="flex flex-col gap-2 w-64 p-1">
            <label class="text-xs font-medium text-primary">{{ l('aliasLabel') }}</label>
            <InputText
              v-model="aliasDraft"
              :placeholder="defaultOperatorName"
              size="small"
              @keydown.enter.prevent="saveAlias"
            />
            <small class="text-muted-color">{{ l('aliasHint') }}</small>
            <div class="flex justify-end gap-1">
              <Button
                :label="l('aliasReset')"
                size="small"
                severity="secondary"
                text
                @click="resetAlias"
              />
              <Button :label="l('aliasSave')" size="small" @click="saveAlias" />
            </div>
          </div>
        </Popover>

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
                :class="msg.isOutgoing ? 'items-end' : 'items-start'"
              >
                <span v-if="msg.authorName" class="text-xs font-semibold text-primary px-1">
                  {{ msg.authorName }}
                </span>
                <div
                  :class="[
                    'p-3 rounded-xl max-w-[85%] text-sm shadow-sm',
                    msg.isPrimary ? 'bg-primary text-primary-contrast' : 'bg-emphasis text-color',
                    msg.role === 'user' ? 'whitespace-pre-wrap' : '',
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

            <CChatComposer
              :placeholder="composerPlaceholder"
              :disabled="!canSend"
              @send="submitComposer"
            />
          </template>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
// CChatbotInbox: shared operator inbox panel. Consumed by the Compose
// ChatbotInboxBlock (with a PageBlock wrapper around it) and by the chatbot
// admin app's Sessions tab (rendered directly). The component is pure UX +
// SDK calls; it has no opinion on what wraps it.
//
// All user-facing strings come from the `translations` prop (with English
// fallbacks built in) so the component has zero coupling to vue-i18n or any
// locale file. Consumers can pass an object built from their own `t()` to
// localize, or omit the prop entirely to use defaults.
import { renderMarkdown } from '@planetcrust/human-js'
import { computed, inject, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import CChatComposer from '../agent/CChatComposer.vue'
import { useConfirm } from 'primevue/useconfirm'

const props = defineProps({
  // Chatbot IDs to watch. Empty list → empty state ("pick chatbots").
  chatbotIDs: { type: Array, default: () => [] },
  // Status scope passed to the BE list call. Empty list → all statuses.
  statusFilter: { type: Array, default: () => [] },
  // Auto-select the most recent session on first load.
  autoOpenFirst: { type: Boolean, default: false },
  // Show the per-row filter popover (status / chatbot / source narrowing).
  showFilter: { type: Boolean, default: true },
  // Polling cadence in seconds; 0 disables polling.
  refreshRate: { type: Number, default: 0 },
  // String overrides — keyed by the flat names listed in DEFAULT_TRANSLATIONS
  // below. Anything missing falls back to the English default.
  translations: { type: Object, default: () => ({}) },
})

const $SystemAPI = inject('$SystemAPI', null)
const $Auth = inject('$Auth', null)
const confirm = useConfirm()

// l(key, params) resolves a string from props.translations. No fallbacks are
// baked in — consumers own the labels. Missing keys render as empty so a
// forgotten translation is visible but doesn't leak the internal key name.
// {name}-style interpolation is supported for dynamic values.
function l(key, params) {
  const raw = props.translations[key]
  if (raw == null) return ''
  if (params && typeof raw === 'string') {
    return raw.replace(/\{(\w+)\}/g, (_, k) => (params[k] !== undefined ? String(params[k]) : ''))
  }
  return raw
}

const tagSmallPt = {
  root: { class: '!text-[11px] !leading-none !py-1 !px-2' },
}

const chatbotIDs = computed(() => props.chatbotIDs || [])
const statusFilter = computed(() => props.statusFilter || [])
const autoOpenFirst = computed(() => !!props.autoOpenFirst)
const showFilter = computed(() => !!props.showFilter)

const loading = ref(false)
const sessions = ref([])
const chatbotsByID = ref({})

const filterMenu = ref(null)
const userFilter = ref({
  statuses: [],
  chatbotIDs: [],
  sources: [],
})
const hasActiveFilter = computed(
  () =>
    userFilter.value.statuses.length > 0 ||
    userFilter.value.chatbotIDs.length > 0 ||
    userFilter.value.sources.length > 0,
)

// Map status enum → translation key for both the filter chips and the inline
// status badge. Used by `statusLabel` (badge) and `filterStatusChoices`.
const STATUS_LABEL_KEYS = {
  handoff_requested: 'statusHandoffRequested',
  handoff_active: 'statusHandoffActive',
  active: 'statusActive',
  closed: 'statusClosed',
}

const filterStatusChoices = computed(() =>
  (statusFilter.value || []).map(v => ({
    value: v,
    label: l(STATUS_LABEL_KEYS[v] || v),
  })),
)

const filterChatbotChoices = computed(() =>
  (chatbotIDs.value || []).map(id => {
    const cb = chatbotsByID.value[id]
    return { value: String(id), label: cb?.name || cb?.handle || String(id) }
  }),
)

const filterSourceChoices = computed(() => [
  { value: 'live', label: l('filterSourceLive') },
  { value: 'preview', label: l('filterSourcePreview') },
])

const filteredSessions = computed(() => {
  if (!hasActiveFilter.value) return sessions.value
  const { statuses, chatbotIDs: cbIDs, sources } = userFilter.value
  return sessions.value.filter(s => {
    if (statuses.length && !statuses.includes(s.status)) return false
    if (sources.length) {
      const src = s.source === 'preview' ? 'preview' : 'live'
      if (!sources.includes(src)) return false
    }
    if (cbIDs.length && !cbIDs.includes(String(s.chatbotID))) return false
    return true
  })
})

const sectionDefs = [
  { key: 'handoff_requested', labelKey: 'sectionHandoffRequested' },
  { key: 'handoff_active', labelKey: 'sectionHandoffActive' },
  { key: 'active', labelKey: 'sectionActive' },
  { key: 'closed', labelKey: 'sectionClosed' },
]

const sessionGroups = computed(() => {
  const buckets = Object.fromEntries(sectionDefs.map(d => [d.key, []]))
  for (const s of filteredSessions.value) {
    const k = buckets[s.status] ? s.status : 'active'
    buckets[k].push(s)
  }
  return sectionDefs
    .map(d => ({ ...d, label: l(d.labelKey), items: buckets[d.key] }))
    .filter(g => g.items.length)
})

const collapsedSections = ref({ closed: true })

function toggleSection(key) {
  collapsedSections.value = {
    ...collapsedSections.value,
    [key]: !collapsedSections.value[key],
  }
}

function isCollapsed(key) {
  return !!collapsedSections.value[key]
}

function resetFilter() {
  userFilter.value = { statuses: [], chatbotIDs: [], sources: [] }
}

const selectedID = ref(null)
const selected = computed(() => sessions.value.find(s => s.id === selectedID.value) || null)

const messages = ref([])
const loadingMessages = ref(false)
const acting = ref(false)
const messagesScrollRef = ref(null)
const aliasMenu = ref(null)
const aliasDraft = ref('')
let stream = null

const defaultOperatorName = computed(() => {
  const u = $Auth?.user
  const name = u?.name?.trim().split(/\s+/)[0]
  return name || u?.handle || u?.email || l('operator')
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
    /* private mode / quota: ignore */
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

const visibleMessages = computed(() =>
  messages.value.map(m => {
    if (m.role === 'user') {
      return {
        ...m,
        authorName: l('authorUser'),
        isOutgoing: false,
        isPrimary: false,
      }
    }
    if (m.role === 'assistant') {
      const isOperator = !!m.operator
      return {
        ...m,
        authorName: isOperator ? m.operator : l('authorAgent'),
        isOutgoing: true,
        isPrimary: isOperator,
      }
    }
    return { ...m, isOutgoing: false, isPrimary: false }
  }),
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

function submitComposer(input) {
  const txt = input.trim()
  if (!txt || !canSend.value) return
  void onSend(txt)
}

const handoffStatus = computed(() => selected.value?.handoff?.status || '')
const canManage = computed(() => !!selected.value?.canManage)
const canManageHandoff = computed(() => !!selected.value?.canManageHandoff)
const canAccept = computed(() => handoffStatus.value === 'requested' && canManageHandoff.value)
const canSend = computed(() => handoffStatus.value === 'active' && canManageHandoff.value)
const canResolve = computed(
  () => ['requested', 'active'].includes(handoffStatus.value) && canManageHandoff.value,
)

const canForceAdvance = computed(() => canManage.value && selected.value?.status !== 'closed')
const canForceClose = computed(() => canManage.value && selected.value?.status !== 'closed')

const adminMenu = ref(null)
const adminMenuItems = computed(() => [
  {
    label: l('adminAdvanceStep'),
    icon: 'pi pi-forward',
    disabled: !canForceAdvance.value || acting.value,
    command: () => onForceAdvance(),
  },
  {
    label: l('adminCloseSession'),
    icon: 'pi pi-stop-circle',
    disabled: !canForceClose.value || acting.value,
    command: () => onForceClose(),
  },
])

const composerPlaceholder = computed(() =>
  canSend.value ? l('composerPlaceholder') : l('composerDisabled'),
)

async function loadChatbots() {
  if (!$SystemAPI || !chatbotIDs.value.length) return
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

async function fetchOnce({ showSpinner = false } = {}) {
  if (!$SystemAPI || !chatbotIDs.value.length) {
    sessions.value = []
    return
  }
  if (showSpinner) loading.value = true
  try {
    const res = await $SystemAPI.chatbotSessionList({
      status: statusFilter.value,
      sort: 'createdAt DESC',
      limit: 100,
    })
    const watched = new Set(chatbotIDs.value.map(String))
    const merged = (res?.set || []).filter(
      s => s.source === 'preview' || watched.has(String(s.chatbotID)),
    )
    merged.sort((a, b) => (a.createdAt < b.createdAt ? 1 : -1))
    sessions.value = merged

    if (autoOpenFirst.value && !selectedID.value && merged.length) {
      void selectSession(merged[0])
    }
    if (selectedID.value && !merged.some(s => s.id === selectedID.value)) {
      clearSelection()
    }
  } catch (err) {
    console.warn('[chatbot-inbox] list failed', err)
    sessions.value = []
  } finally {
    if (showSpinner) loading.value = false
  }
}

function chatbotName(s) {
  const id = s?.chatbotID ?? s
  if (id === '0') return l('draftChatbot')
  const cb = chatbotsByID.value[id]
  return cb?.name || cb?.handle || id
}

function statusLabel(status) {
  if (!status) return ''
  return l(STATUS_LABEL_KEYS[status] || status)
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
  if (s.canView === false) {
    console.warn('[chatbot-inbox] select blocked by canView=false', s.id)
    return
  }
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

    const convIDs = []
    const seen = new Set()
    const addConvID = id => {
      if (!id || id === '0' || seen.has(id)) return
      seen.add(id)
      convIDs.push(id)
    }
    for (const step of steps) addConvID(step.conversationID)
    addConvID(res?.session?.conversationID)

    const all = []
    for (const id of convIDs) {
      try {
        const conv = await $SystemAPI.aiConversationRead({ aiConversationID: id })
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

function patchSelected(patch) {
  if (!selectedID.value) return
  sessions.value = sessions.value.map(s => {
    if (s.id !== selectedID.value) return s
    const next = { ...s }
    if (Object.prototype.hasOwnProperty.call(patch, 'status')) {
      next.status = patch.status
    }
    if (Object.prototype.hasOwnProperty.call(patch, 'handoff')) {
      next.handoff = patch.handoff ? { ...(s.handoff || {}), ...patch.handoff } : undefined
    }
    return next
  })
}

function attachStream(sessionID) {
  closeStream()
  const tok = $SystemAPI?.accessTokenFn?.() || ''
  const base = ($SystemAPI?.baseURL || '').replace(/\/+$/, '')
  const qs = tok ? `?jwt=${encodeURIComponent(tok)}` : ''
  const url = `${base}/chatbots/sessions/${encodeURIComponent(sessionID)}/stream${qs}`
  stream = new EventSource(url, { withCredentials: true })

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
  stream.addEventListener('handoff_requested', () => {
    patchSelected({ status: 'handoff_requested', handoff: { status: 'requested' } })
  })
  stream.addEventListener('handoff_active', () => {
    patchSelected({ status: 'handoff_active', handoff: { status: 'active' } })
  })
  stream.addEventListener('handoff_complete', () => {
    patchSelected({ status: 'active', handoff: { status: 'closed' } })
  })
  stream.addEventListener('session_closed', () => {
    patchSelected({ status: 'closed' })
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
    message: l('confirmAcceptMessage'),
    header: l('confirmAcceptHeader'),
    rejectProps: {
      label: l('cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: {
      label: l('accept'),
      severity: 'success',
      size: 'small',
    },
    accept: () => void doAccept(),
  })
}

async function doAccept() {
  if (!selectedID.value) return
  acting.value = true
  try {
    await $SystemAPI.chatbotSessionHandoffAccept({
      sessionID: selectedID.value,
      operator: operatorName.value,
    })
    patchSelected({ status: 'handoff_active', handoff: { status: 'active' } })
  } catch (err) {
    console.error('[chatbot-inbox] accept failed', err)
  } finally {
    acting.value = false
  }
}

async function onSend(text) {
  if (!text || !canSend.value || !selectedID.value) return
  try {
    await $SystemAPI.chatbotSessionOperatorMessage({
      sessionID: selectedID.value,
      message: text,
      operator: operatorName.value,
    })
  } catch (err) {
    console.error('[chatbot-inbox] send failed', err)
  }
}

function onResolve() {
  if (!selected.value || acting.value) return
  confirm.require({
    message: l('confirmResolveMessage'),
    header: l('confirmResolveHeader'),
    rejectProps: {
      label: l('cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: {
      label: l('resolve'),
      severity: 'warn',
      size: 'small',
    },
    accept: () => void doResolve(),
  })
}

async function doResolve() {
  if (!selectedID.value) return
  acting.value = true
  try {
    await $SystemAPI.chatbotSessionHandoffComplete({ sessionID: selectedID.value })
    patchSelected({ status: 'active', handoff: { status: 'closed' } })
  } catch (err) {
    console.error('[chatbot-inbox] resolve failed', err)
  } finally {
    acting.value = false
  }
}

function onForceAdvance() {
  if (!canForceAdvance.value || acting.value) return
  confirm.require({
    message: l('confirmAdvanceStepMessage'),
    header: l('confirmAdvanceStepHeader'),
    rejectProps: {
      label: l('cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: {
      label: l('adminAdvanceStep'),
      size: 'small',
    },
    accept: () => void doForceAdvance(),
  })
}

async function doForceAdvance() {
  if (!selectedID.value) return
  acting.value = true
  try {
    await $SystemAPI.chatbotSessionAdvanceStep({ sessionID: selectedID.value })
    void fetchOnce({ showSpinner: false })
  } catch (err) {
    console.error('[chatbot-inbox] advance step failed', err)
  } finally {
    acting.value = false
  }
}

function onForceClose() {
  if (!canForceClose.value || acting.value) return
  confirm.require({
    message: l('confirmCloseSessionMessage'),
    header: l('confirmCloseSessionHeader'),
    rejectProps: {
      label: l('cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: {
      label: l('adminCloseSession'),
      severity: 'danger',
      size: 'small',
    },
    accept: () => void doForceClose(),
  })
}

async function doForceClose() {
  if (!selectedID.value) return
  acting.value = true
  try {
    await $SystemAPI.chatbotSessionClose({ sessionID: selectedID.value })
    patchSelected({ status: 'closed', handoff: { status: 'closed' } })
  } catch (err) {
    console.error('[chatbot-inbox] close session failed', err)
  } finally {
    acting.value = false
  }
}

watch(
  () => [chatbotIDs.value.join(','), statusFilter.value.join(',')],
  () => {
    void loadChatbots()
    void fetchOnce({ showSpinner: !sessions.value.length })
  },
  { immediate: true },
)

// Local poll driven by the refreshRate prop. Wrappers that already poll
// (e.g. Compose PageBlock's @refreshBlock event) can leave refreshRate at 0
// and call refresh() themselves; standalone consumers pass a value.
let pollHandle = null
watch(
  () => Number(props.refreshRate) || 0,
  rate => {
    if (pollHandle) {
      clearInterval(pollHandle)
      pollHandle = null
    }
    if (rate > 0) {
      pollHandle = setInterval(() => {
        void fetchOnce({ showSpinner: false })
      }, rate * 1000)
    }
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  if (pollHandle) clearInterval(pollHandle)
  closeStream()
})

defineExpose({
  refresh: () => fetchOnce({ showSpinner: false }),
})
</script>
