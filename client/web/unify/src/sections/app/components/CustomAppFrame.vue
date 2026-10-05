<template>
  <Message
    v-if="failure"
    severity="warn"
    :closable="false"
    class="m-4"
    data-test-id="custom-app-failure"
  >
    {{ failure }}
  </Message>

  <iframe
    v-if="outerDocument"
    ref="frameRef"
    :srcdoc="outerDocument"
    :title="name || $t('app.title')"
    class="block w-full border-0"
    :style="{ height: frameHeight }"
  />
</template>

<script setup>
// One custom app's HTML in the two-frame sandbox, with the host end of its
// bridge. Shared by the app view and the compose Custom block; see
// app.intent.md for the boundary it keeps.
import { getThemeVariables, useUserStore } from '@planetcrust/human-vue'
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useConfirm } from 'primevue/useconfirm'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { mountChatbot } from 'human-webapp-chatbot-widget/mount'
import { BRIDGE_SCRIPT } from '../bridge'
import { bridgeVersion, buildOuterDocument, cspInner, dispatch, hostScriptSource } from '../host'

const props = defineProps({
  source: { type: String, required: true },
  // What the source declares: namespace, modules, writes, and the IDs they
  // were resolved to when stored, where they were.
  sourceMeta: { type: Object, default: () => ({}) },
  // Named to the viewer when the app asks to change records.
  name: { type: String, default: '' },
  // What `human.context()` answers: where the app is shown.
  context: { type: Object, default: null },
  // Whether `human.resize()` sets the frame height; a host with a height of
  // its own leaves the frame filling it.
  resizable: { type: Boolean, default: true },
})

// `changed`: the app changed data, or asked for the page to catch up.
const emit = defineEmits(['navigated', 'changed'])

const { t } = useI18n()
const confirm = useConfirm()
const router = useRouter()

const $Auth = inject('$Auth')
const $ComposeAPI = inject('$ComposeAPI')
const $AutomationAPI = inject('$AutomationAPI')
const $SystemAPI = inject('$SystemAPI')

const userStore = useUserStore()

// What the app itself said went wrong, shown over the page it half drew.
const failure = ref('')
const outerDocument = ref('')
// The app fills the content area unless it asks for a height of its own.
const frameHeight = ref('100%')

// What the source declares it may reach, resolved to the IDs the compose
// endpoints take.
const sourceMeta = ref({})
const namespaceID = ref('')
const moduleIDs = ref({})
// Each declared module's fields, keyed by module ID: they type the values and
// say which of them are user or record references.
const moduleFields = ref({})
// The bridge contract the page was written against.
const version = ref(1)

const frameRef = ref(null)

function themeInfo() {
  const styles = getComputedStyle(document.documentElement)
  const read = name => styles.getPropertyValue(name).trim()
  const dark = document.documentElement.classList.contains('dark')

  // The instance's own palette, under the names it is read by everywhere else:
  // an app told to look the theme up through the API is shown these, and has
  // to get the same ones back here. `surface` is a set of shades rather than a
  // colour, so it is left out of what an app is handed.
  const palette = {}
  for (const [name, value] of Object.entries(getThemeVariables(dark ? 'dark' : 'light'))) {
    if (typeof value === 'string') palette[name] = value
  }

  return {
    dark,
    colors: {
      ...palette,
      primary: read('--p-primary-color') || palette.primary,
      'body-bg': read('--body-bg') || palette['body-bg'],
      'content-bg': read('--p-content-background'),
      text: read('--p-text-color'),
      'text-muted': read('--p-text-muted-color'),
      border: read('--p-content-border-color'),
    },
  }
}

function userInfo() {
  const user = $Auth.user || {}
  return {
    userID: user.userID,
    name: user.name || user.username || user.email || '',
    email: user.email || '',
  }
}

// The names behind a set of user IDs.
async function userLabels(ids) {
  const wanted = [...new Set(ids)].filter(id => id && id !== '0')
  if (!wanted.length) return {}

  await userStore.resolveUsers(wanted).catch(() => {})

  const out = {}
  for (const id of wanted) {
    const user = userStore.findByID(id)
    if (user) out[id] = user.name || user.username || user.email || id
  }
  return out
}

// User IDs a listing carries, mapped to the name a person would see.
async function refs(records) {
  const ids = new Set()

  for (const record of records) {
    const fields = (moduleFields.value[record.moduleID] || [])
      .filter(field => field.kind === 'User')
      .map(field => field.name)
    for (const id of [record.ownedBy, record.createdBy]) {
      if (id && id !== '0') ids.add(id)
    }
    for (const { name, value } of record.values || []) {
      if (fields.includes(name) && value) ids.add(value)
    }
  }

  return userLabels([...ids])
}

// What the viewer has been asked about this app, for as long as it is open.
// An app may not ask again after a refusal, and is asked once after an answer.
let consented = null

function askToChangeRecords() {
  if (consented !== null) return Promise.resolve(consented)

  return new Promise(resolve => {
    const answer = value => {
      consented = value
      resolve(value)
    }

    confirm.require({
      header: t('app.consent.header'),
      message: t('app.consent.message', {
        name: props.name,
        modules: (sourceMeta.value.writes || []).join(', '),
      }),
      icon: 'pi pi-pencil',
      rejectProps: { label: t('app.consent.reject'), severity: 'secondary', text: true },
      acceptProps: { label: t('app.consent.accept') },
      accept: () => answer(true),
      reject: () => answer(false),
      onHide: () => answer(consented === true),
    })
  })
}

// The same, for running what the app declared; asked on its own, because a
// viewer who lets an app edit records has not agreed to it starting workflows.
let consentedToRun = null

function askToRun() {
  if (consentedToRun !== null) return Promise.resolve(consentedToRun)

  return new Promise(resolve => {
    const answer = value => {
      consentedToRun = value
      resolve(value)
    }

    confirm.require({
      header: t('app.consent.runHeader'),
      message: t('app.consent.runMessage', {
        name: props.name,
        automations: (sourceMeta.value.automations || []).join(', '),
      }),
      icon: 'pi pi-bolt',
      rejectProps: { label: t('app.consent.reject'), severity: 'secondary', text: true },
      acceptProps: { label: t('app.consent.runAccept') },
      accept: () => answer(true),
      reject: () => answer(false),
      onHide: () => answer(consentedToRun === true),
    })
  })
}

// Runs a declared workflow or TAQ by handle, as the viewer, through the same
// endpoints a page's automation button uses. The record the app is shown on
// goes along, the way a button on a record page sends it.
async function runAutomation(handle, input) {
  const where = props.context || {}
  if (where.recordID && where.moduleID) {
    const record = await $ComposeAPI
      .recordRead({
        namespaceID: namespaceID.value,
        moduleID: where.moduleID,
        recordID: where.recordID,
      })
      .catch(() => null)
    if (record) input = { record: { '@type': 'ComposeRecord', '@value': record }, ...input }
  }

  const listed = await $AutomationAPI
    .ngAutomationList({ query: handle, limit: 50 })
    .catch(() => ({}))
  const taqs = Array.isArray(listed) ? listed : listed?.set || []
  const taq = taqs.find(a => a.handle === handle)
  if (taq) {
    await $AutomationAPI.ngAutomationExec({ automationID: taq.automationID, input, wait: true })
    return { ok: true }
  }

  const { set: workflows = [] } = await $AutomationAPI
    .workflowList({ query: handle, limit: 50 })
    .catch(() => ({}))
  const workflow = workflows.find(w => w.handle === handle)
  if (!workflow) throw new Error(`there is no automation "${handle}" this viewer can run`)

  const { set: triggers = [] } = await $AutomationAPI
    .triggerList({ workflowID: workflow.workflowID, eventType: 'onManual' })
    .catch(() => ({}))
  if (!triggers.length) throw new Error(`workflow "${handle}" has no manual trigger to start it`)

  await $AutomationAPI.workflowExec({
    workflowID: workflow.workflowID,
    stepID: triggers[0].stepID || '0',
    input,
    wait: true,
  })
  return { ok: true }
}

// Chatbots the app opened, by handle. Each is Human's own widget, mounted in
// the shell rather than the sandbox, under the chatbot's own settings — it
// must be enabled and allow this site's origin. They go when the app does.
const chatbots = new Map()

async function openChatbot(handle) {
  if (chatbots.has(handle)) {
    chatbots.get(handle).open()
    return
  }

  const { set = [] } = await $SystemAPI.chatbotList({ handle, limit: 1 }).catch(() => ({}))
  const chatbot = set.find(c => c.handle === handle)
  if (!chatbot?.widgetKey) throw new Error(`there is no chatbot "${handle}" this viewer can open`)
  if (!chatbot.enabled) throw new Error(`chatbot "${handle}" is switched off`)

  const apiOrigin = new URL($SystemAPI.baseURL, window.location.href).origin
  const mounted = await mountChatbot({
    scriptSrc: `${apiOrigin}/chatbot/widget.js`,
    widgetKey: chatbot.widgetKey,
    startOpen: true,
  })
  if (!mounted) {
    throw new Error(
      `chatbot "${handle}" could not start; check that it allows ${window.location.origin}`,
    )
  }
  chatbots.set(handle, mounted)
}

function closeChatbot(handle) {
  chatbots.get(handle)?.close()
}

const ctx = {
  compose: $ComposeAPI,
  consent: askToChangeRecords,
  get meta() {
    return sourceMeta.value
  },
  get namespaceID() {
    return namespaceID.value
  },
  get moduleIDs() {
    return moduleIDs.value
  },
  get version() {
    return version.value
  },
  fields: moduleID => moduleFields.value[moduleID] || [],
  refs,
  userLabels,
  user: userInfo,
  theme: themeInfo,
  // A plain copy: what the host hands over is posted to the frame, and a
  // reactive object cannot be cloned across it.
  context: () => JSON.parse(JSON.stringify(props.context || {})),
  navigate,
  consentToRun: askToRun,
  openChatbot,
  closeChatbot,
  changed: () => emit('changed'),
  runAutomation,
  attachment: attachmentID =>
    $ComposeAPI.attachmentRead({ kind: 'record', namespaceID: namespaceID.value, attachmentID }),
  fileData,
  // The sandbox blocks a download of its own, so the shell hands the file over.
  download: (name, text) => {
    const url = URL.createObjectURL(new Blob([text], { type: 'text/plain;charset=utf-8' }))
    const link = document.createElement('a')
    link.href = url
    link.download = name
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(url)
  },
  resize: height => {
    if (!props.resizable) return
    const asked = Number(height)
    frameHeight.value = Number.isFinite(asked) && asked > 0 ? `${Math.floor(asked)}px` : '100%'
  },
}

function reply(payload) {
  frameRef.value?.contentWindow?.postMessage(
    { type: 'human:result', ...payload },
    window.location.origin,
  )
}

async function onMessage(event) {
  if (event.origin !== window.location.origin) return
  if (event.source !== frameRef.value?.contentWindow) return

  const data = event.data || {}

  if (data.type === 'human:navigated') {
    emit('navigated')
    return
  }

  if (data.type === 'human:failed') {
    failure.value = t('app.state.failed', { message: data.message || '' })
    return
  }

  if (data.type !== 'human:call') return

  try {
    reply({ id: data.id, result: await dispatch(data.op, data.args || {}, ctx) })
  } catch (error) {
    reply({ id: data.id, error: error?.message || String(error) })
  }
}

// The IDs a page is stored with, and the fields of each module it declares.
//
// A page stored before the declaration was resolved carries handles only, and
// finding the namespace by slug then needs a search the viewer may not be
// allowed; those pages keep working, the rest ask nothing extra of a viewer.
async function resolveDeclared(meta) {
  if (!meta.namespace && !meta.namespaceID) return

  namespaceID.value = await namespaceIDFor(meta)
  if (!namespaceID.value) return

  for (const handle of meta.modules || []) {
    const module = await moduleFor(handle, meta.moduleIDs?.[handle])
    if (!module) continue

    moduleIDs.value[handle] = module.moduleID
    moduleFields.value[module.moduleID] = module.fields || []
  }
}

// The stored ID first, then the slug or handle: a page brought in from another
// instance carries IDs that mean nothing here, while the names it declared
// still point at the same things.
async function namespaceIDFor(meta) {
  if (meta.namespaceID) {
    const stored = await $ComposeAPI
      .namespaceRead({ namespaceID: meta.namespaceID })
      .catch(() => null)
    if (stored) return stored.namespaceID
  }
  return namespaceIDBySlug(meta.namespace)
}

async function moduleFor(handle, storedID) {
  const read = moduleID =>
    $ComposeAPI.moduleRead({ namespaceID: namespaceID.value, moduleID }).catch(() => null)

  const stored = storedID ? await read(storedID) : null
  if (stored) return stored

  const moduleID = await moduleIDByHandle(handle)
  return moduleID ? read(moduleID) : null
}

async function namespaceIDBySlug(slug) {
  const { set = [] } = await $ComposeAPI.namespaceList({ slug, limit: 1 }).catch(() => ({}))
  return set[0]?.namespaceID || ''
}

async function moduleIDByHandle(handle) {
  const { set = [] } = await $ComposeAPI
    .moduleList({
      namespaceID: namespaceID.value,
      handle,
      limit: 1,
    })
    .catch(() => ({}))
  return set[0]?.moduleID || ''
}

// A record attachment as a data: URL, fetched by the shell through the
// attachment's signed address, so the sandbox needs no network for it.
async function fileData(attachmentID) {
  const { url } = await $ComposeAPI.attachmentRead({
    kind: 'record',
    namespaceID: namespaceID.value,
    attachmentID,
  })
  const response = await fetch(`${$ComposeAPI.baseURL}${url}`)
  if (!response.ok) throw new Error(`the file could not be read (HTTP ${response.status})`)
  const blob = await response.blob()
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result)
    reader.onerror = () => reject(new Error('the file could not be read'))
    reader.readAsDataURL(blob)
  })
}

// Opens a page of the namespace the app reads, in the shell, as the viewer:
// the page itself still asks its own permissions.
async function navigate({ page, module, recordID }) {
  const slug = sourceMeta.value.namespace
  if (!namespaceID.value || !slug) {
    throw new Error('this app reads no namespace, so there is no page for it to open')
  }

  const first = async query =>
    (
      await $ComposeAPI
        .pageList({ namespaceID: namespaceID.value, limit: 1, ...query })
        .catch(() => ({}))
    ).set?.[0]

  let target = null
  if (page && /^[0-9]+$/.test(page)) {
    target = await $ComposeAPI
      .pageRead({ namespaceID: namespaceID.value, pageID: page })
      .catch(() => null)
  } else if (page) {
    target = await first({ handle: page })
  } else {
    const moduleID = moduleIDs.value[module] || (await moduleIDByHandle(module))
    target = moduleID ? await first({ moduleID }) : null
  }

  if (!target) throw new Error(`there is no page "${page || module}" in namespace "${slug}"`)

  const recordPage = target.moduleID && target.moduleID !== '0'
  if (recordPage && !recordID)
    throw new Error(`page "${page}" shows one record; name it with recordID`)
  if (!recordPage && recordID) throw new Error(`page "${page}" is not a record page`)

  const path = `/compose/namespace/${slug}/pages/${target.pageID}`
  await router.push(recordID ? `${path}/records/${recordID}` : path)
}

async function build() {
  failure.value = ''
  outerDocument.value = ''
  frameHeight.value = '100%'
  sourceMeta.value = props.sourceMeta || {}
  consented = null
  consentedToRun = null
  namespaceID.value = ''
  moduleIDs.value = {}
  moduleFields.value = {}
  version.value = bridgeVersion(props.source)

  await resolveDeclared(sourceMeta.value)

  outerDocument.value = buildOuterDocument({
    source: props.source,
    bridgeScript: BRIDGE_SCRIPT,
    hostScript: hostScriptSource({ origin: window.location.origin }),
    cspInner: cspInner(sourceMeta.value.origins),
  })
}

onMounted(() => {
  window.addEventListener('message', onMessage)
})

watch(() => [props.source, props.sourceMeta], build, { immediate: true })

onBeforeUnmount(() => {
  window.removeEventListener('message', onMessage)
  for (const chatbot of chatbots.values()) chatbot.destroy()
  chatbots.clear()
})
</script>
