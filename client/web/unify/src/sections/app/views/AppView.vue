<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ application?.name || $t('app.title') }}</span>
  </Teleport>

  <div class="h-full w-full overflow-auto">
    <div v-if="loading" class="h-full flex items-center justify-center">
      <ProgressSpinner style="width: 2.5rem; height: 2.5rem" />
    </div>

    <Message v-else-if="problem" severity="error" class="m-4">
      {{ problem }}
    </Message>

    <iframe
      v-else
      ref="frameRef"
      :srcdoc="outerDocument"
      :title="application?.name || $t('app.title')"
      class="block w-full border-0"
      :style="{ height: frameHeight }"
    />
  </div>
</template>

<script setup>
import { useApplicationsStore, useUserStore } from '@planetcrust/human-vue'
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { BRIDGE_SCRIPT } from '../bridge'
import { buildOuterDocument, dispatch, hostScriptSource } from '../host'

const { t } = useI18n()
const route = useRoute()

const $Auth = inject('$Auth')
const $SystemAPI = inject('$SystemAPI')
const $ComposeAPI = inject('$ComposeAPI')

const applicationsStore = useApplicationsStore()
const userStore = useUserStore()

const loading = ref(true)
const problem = ref('')
const application = ref(null)
const outerDocument = ref('')
// The app fills the content area unless it asks for a height of its own.
const frameHeight = ref('100%')

// What the source declares it may reach, resolved to the IDs the compose
// endpoints take.
const sourceMeta = ref({})
const namespaceID = ref('')
const moduleIDs = ref({})
// User-kind fields per module, so a listing's user IDs can be labelled.
const userFields = ref({})

const frameRef = ref(null)

function themeInfo() {
  const styles = getComputedStyle(document.documentElement)
  const read = name => styles.getPropertyValue(name).trim()

  return {
    dark: document.documentElement.classList.contains('dark'),
    colors: {
      primary: read('--p-primary-color'),
      'body-bg': read('--body-bg'),
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

// User IDs a listing carries, mapped to the name a person would see.
async function refs(records) {
  const ids = new Set()

  for (const record of records) {
    const fields = userFields.value[record.moduleID] || []
    for (const id of [record.ownedBy, record.createdBy]) {
      if (id && id !== '0') ids.add(id)
    }
    for (const { name, value } of record.values || []) {
      if (fields.includes(name) && value) ids.add(value)
    }
  }

  if (!ids.size) return {}

  await userStore.resolveUsers([...ids]).catch(() => {})

  const out = {}
  for (const id of ids) {
    const user = userStore.findByID(id)
    if (user) out[id] = user.name || user.username || user.email || id
  }
  return out
}

const ctx = {
  compose: $ComposeAPI,
  get meta() {
    return sourceMeta.value
  },
  get namespaceID() {
    return namespaceID.value
  },
  get moduleIDs() {
    return moduleIDs.value
  },
  refs,
  user: userInfo,
  theme: themeInfo,
  resize: height => {
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
    problem.value = t('app.state.navigated')
    return
  }

  if (data.type !== 'human:call') return

  try {
    reply({ id: data.id, result: await dispatch(data.op, data.args || {}, ctx) })
  } catch (error) {
    reply({ id: data.id, error: error?.message || String(error) })
  }
}

// Handles and slugs are what the source declares; the endpoints take IDs.
async function resolveDeclared(meta) {
  if (!meta.namespace) return

  const { set: namespaces = [] } = await $ComposeAPI.namespaceList({
    slug: meta.namespace,
    limit: 1,
  })
  const namespace = namespaces[0]
  if (!namespace) return

  namespaceID.value = namespace.namespaceID

  for (const handle of meta.modules || []) {
    const { set: modules = [] } = await $ComposeAPI.moduleList({
      namespaceID: namespace.namespaceID,
      handle,
      limit: 1,
    })
    const module = modules[0]
    if (!module) continue

    moduleIDs.value[handle] = module.moduleID
    userFields.value[module.moduleID] = (module.fields || [])
      .filter(field => field.kind === 'User')
      .map(field => field.name)
  }
}

// One application per route, and the route changes without the view being
// remounted when the user goes from one custom app straight to another.
async function load(applicationID) {
  loading.value = true
  problem.value = ''
  application.value = null
  outerDocument.value = ''
  frameHeight.value = '100%'
  sourceMeta.value = {}
  namespaceID.value = ''
  moduleIDs.value = {}
  userFields.value = {}

  try {
    application.value = await applicationsStore.findByID(applicationID)
  } catch {
    application.value = null
  }

  // The route guard has already refused what this catches; a plain message is
  // for the case where the view is reached some other way.
  if (
    !application.value ||
    application.value.unify?.kind !== 'custom' ||
    !application.value.canAccessApplication
  ) {
    problem.value = t('app.state.refused')
    loading.value = false
    return
  }

  try {
    const { source, sourceMeta: meta } = await $SystemAPI.applicationSourceRead({ applicationID })

    if (!source) {
      problem.value = t('app.state.empty')
      return
    }

    sourceMeta.value = meta || {}
    await resolveDeclared(sourceMeta.value)

    outerDocument.value = buildOuterDocument({
      source,
      bridgeScript: BRIDGE_SCRIPT,
      hostScript: hostScriptSource({ origin: window.location.origin }),
    })
  } catch (error) {
    problem.value = t('app.state.error', { reason: error?.message || String(error) })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  window.addEventListener('message', onMessage)
})

watch(
  () => String(route.params.applicationID || ''),
  applicationID => {
    if (applicationID) load(applicationID)
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  window.removeEventListener('message', onMessage)
})
</script>
