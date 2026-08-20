<template>
  <Dialog
    v-model:visible="dialogVisible"
    :header="dialogTitle"
    modal
    :style="{ width: '80rem' }"
    :breakpoints="{ '80rem': '96vw' }"
    :pt="{
      header: { class: 'p-3 pr-2' },
      content: { class: 'flex flex-col min-h-0 p-0' },
      footer: { class: 'border-t border-surface p-3' },
    }"
    @hide="onHide"
  >
    <!-- Loading state -->
    <div
      v-if="processing"
      class="flex flex-col items-center justify-center gap-3"
      style="min-height: 60vh"
    >
      <ProgressSpinner style="width: 2rem; height: 2rem" />
      <span class="text-sm text-muted-color">{{ $t('permissions.ui.loading') }}</span>
    </div>

    <!-- Content: Two-column layout -->
    <div v-else class="flex flex-col min-h-0">
      <!-- Single scrollable container with sticky headers -->
      <div class="flex-1 overflow-y-auto min-h-0" style="max-height: 80vh; min-height: 60vh">
        <!-- Description row (sticky) -->
        <div class="flex bg-emphasis sticky top-0 z-10">
          <div class="perm-col-left p-3 text-sm text-muted-color">
            {{ $t('permissions.ui.edit.description') }}
          </div>
          <div class="flex-1 p-3 text-sm text-muted-color border-l hidden lg:block">
            {{ $t('permissions.ui.evaluate.description') }}
          </div>
        </div>

        <!-- Selector row (sticky) -->
        <div class="flex border-y bg-emphasis sticky top-[37px] z-10">
          <!-- Left: Role picker -->
          <div class="perm-col-left p-3">
            <CFormGroup
              :label="$t('permissions.ui.edit.label')"
              input-id="permissions-role-selector"
            >
              <CInputRole
                id="permissions-role-selector"
                v-model="currentRoleID"
                :placeholder="$t('permissions.ui.edit.label')"
                :show-clear="false"
                :exclude-roles="['super-admin']"
                @select="onRoleChange"
              />
            </CFormGroup>
          </div>

          <!-- Right: Evaluation columns headers -->
          <div class="flex-1 min-w-0 hidden lg:flex">
            <div
              v-for="(e, i) in evaluate"
              :key="i"
              :style="evalColumnStyle(e)"
              class="flex flex-col items-center justify-center min-w-0 p-3 border-l cursor-pointer hover:bg-surface-hover transition-colors overflow-hidden"
              @click="removeEvalColumn(i)"
            >
              <span
                v-for="(n, ni) in getEvalName(e)"
                :key="ni"
                class="text-primary text-center truncate w-full"
                :title="n"
              >
                {{ n }}
              </span>
              <i class="pi pi-times text-muted-color mt-1 eval-remove-icon" />
            </div>

            <!-- Add column button -->
            <div
              v-if="evaluate.length < MAX_EVAL_COLUMNS"
              :style="{ width: ADD_COLUMN_WIDTH }"
              class="flex flex-col items-center justify-center p-3 border-l cursor-pointer hover:bg-surface-hover transition-colors"
              @click="showAddEval = true"
            >
              <span class="text-primary">{{ $t('permissions.ui.add.label') }}</span>
              <i class="pi pi-plus text-green-500 mt-1" />
            </div>

            <!-- Whatever the columns did not need -->
            <div class="flex-1 border-l" />
          </div>
        </div>

        <!-- Rules rows -->
        <div
          v-for="rule in rules"
          :key="rule.operation"
          class="flex border-b hover:bg-emphasis transition-colors"
        >
          <!-- Left: Rule with toggles -->
          <div class="perm-col-left p-3">
            <div class="font-medium text-sm mb-1">{{ rule.title }}</div>
            <div class="flex gap-1">
              <Button
                :severity="rule.access === 'allow' ? 'success' : 'secondary'"
                :outlined="rule.access !== 'allow'"
                size="small"
                :label="$t('permissions.ui.access.allow')"
                @click="setAccess(rule, 'allow')"
              />
              <Button
                :severity="
                  rule.access !== rule.initial && rule.access === 'inherit' ? 'warn' : 'secondary'
                "
                :outlined="rule.access !== 'inherit'"
                size="small"
                :label="$t('permissions.ui.access.inherit')"
                @click="setAccess(rule, 'inherit')"
              />
              <Button
                :severity="rule.access === 'deny' ? 'danger' : 'secondary'"
                :outlined="rule.access !== 'deny'"
                size="small"
                :label="$t('permissions.ui.access.deny')"
                @click="setAccess(rule, 'deny')"
              />
            </div>
          </div>

          <!-- Right: Evaluation results -->
          <div class="flex-1 min-w-0 hidden lg:flex">
            <div
              v-for="(e, i) in evaluate"
              :key="i"
              :style="evalColumnStyle(e)"
              class="flex items-center justify-center min-w-0 border-l bg-emphasis cursor-not-allowed"
            >
              <i
                v-if="getEvalAccess(e, rule.operation) === 'unknown-context'"
                v-tooltip.top="$t('permissions.ui.tooltip.unknown-context.role')"
                class="pi pi-question-circle text-muted-color"
              />
              <i
                v-else-if="getEvalAccess(e, rule.operation) === 'allow'"
                class="pi pi-check text-green-500"
              />
              <i v-else class="pi pi-times text-red-500" />
            </div>

            <!-- Empty column for the "Add" slot, and for the slack after it -->
            <div
              v-if="evaluate.length < MAX_EVAL_COLUMNS"
              :style="{ width: ADD_COLUMN_WIDTH }"
              class="border-l"
            />
            <div class="flex-1 border-l" />
          </div>
        </div>

        <div
          v-if="rules.length === 0 && currentRoleID"
          class="flex items-center justify-center py-8 text-sm text-muted-color"
        >
          {{ $t('permissions.ui.no-rules') }}
        </div>

        <div
          v-if="!currentRoleID"
          class="flex items-center justify-center py-8 text-sm text-muted-color"
        >
          {{ $t('permissions.ui.select-role') }}
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('permissions.ui.cancel')"
          severity="secondary"
          text
          size="small"
          @click="onHide"
        />
        <Button
          :label="$t('permissions.ui.save')"
          :disabled="!dirty || submitting"
          :loading="submitting"
          size="small"
          @click="onSubmit"
        />
      </div>
    </template>
  </Dialog>

  <!-- Add Evaluation Column Dialog -->
  <Dialog
    v-model:visible="showAddEval"
    :header="$t('permissions.ui.add.title')"
    modal
    :style="{ width: '28rem' }"
  >
    <div class="flex flex-col gap-4">
      <CFormGroup :label="$t('permissions.ui.add.role.label')" input-id="eval-role-selector">
        <CInputRole
          id="eval-role-selector"
          v-model="addEval.roleIDs"
          :placeholder="$t('permissions.ui.add.role.placeholder')"
          multiple
          :disabled="!!addEval.userID"
        />
      </CFormGroup>

      <CFormGroup :label="$t('permissions.ui.add.user.label')" input-id="eval-user-selector">
        <CInputUser
          id="eval-user-selector"
          v-model="addEval.userID"
          :placeholder="$t('permissions.ui.add.user.placeholder')"
          :disabled="addEval.roleIDs?.length > 0"
        />
      </CFormGroup>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('permissions.ui.cancel')"
          severity="secondary"
          text
          size="small"
          @click="showAddEval = false"
        />
        <Button
          :label="$t('permissions.ui.add.save')"
          :disabled="!addEvalEnabled"
          size="small"
          @click="onAddEvalColumn"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePermissions } from '../../composables/usePermissions'
import CInputRole from '../input/CInputRole.vue'
import CInputUser from '../input/CInputUser.vue'

const { t, te } = useI18n()
const { visible, options, close } = usePermissions()

const $SystemAPI = inject('$SystemAPI', null)
const $ComposeAPI = inject('$ComposeAPI', null)
const $AutomationAPI = inject('$AutomationAPI', null)
const $toast = inject('$toast')

// State
const processing = ref(false)
const submitting = ref(false)
const currentRoleID = ref(null)
const permissions = ref([])
const rules = ref([])
const initialRules = ref({})

// Evaluation columns state
const MAX_EVAL_COLUMNS = 4

// A column is a name over a column of ticks, so it needs only the width of the
// name. The header and the rule rows are separate rows, though, and each would
// size itself independently — so the width is decided once, from the longest
// name, and handed to both.
const EVAL_COLUMN_MIN = '6rem'
const EVAL_COLUMN_MAX = '14rem'
const ADD_COLUMN_WIDTH = '6rem'
const evaluate = ref([])
const showAddEval = ref(false)
const addEval = ref({ roleIDs: [], userID: null })

// The role the user last picked in this dialog, sticky across dialog opens and
// page reloads. Only an explicit pick in the role selector is recorded — a
// roleID preselected by the opener stays contextual and doesn't overwrite it.
const LAST_ROLE_KEY = 'permissionsDialog.lastRoleID'

function readLastRoleID() {
  try {
    return localStorage.getItem(LAST_ROLE_KEY) || null
  } catch {
    return null
  }
}

function storeLastRoleID(roleID) {
  try {
    localStorage.setItem(LAST_ROLE_KEY, String(roleID))
  } catch {
    // localStorage not available
  }
}

// The evaluation columns, sticky the same way and for the same reason. What is
// kept is who a column compares against, never what it showed: the accesses are
// a property of the resource the dialog was opened for, and are re-traced
// against whichever one that is.
const EVAL_COLUMNS_KEY = 'permissionsDialog.evalColumns'

function readStoredEvalColumns() {
  try {
    // Anything that is not a list of columns reaches the catch below — a
    // non-array has no map, and a bad parse throws before that.
    const stored = JSON.parse(localStorage.getItem(EVAL_COLUMNS_KEY) || '[]')
    return stored
      .map(c => ({
        roleIDs: Array.isArray(c?.roleIDs) ? c.roleIDs.map(String) : [],
        userID: c?.userID ? String(c.userID) : null,
      }))
      .filter(c => c.roleIDs.length || c.userID)
      .slice(0, MAX_EVAL_COLUMNS)
  } catch {
    return []
  }
}

function storeEvalColumns() {
  try {
    const identities = evaluate.value.map(({ roleIDs, userID }) => ({ roleIDs, userID }))
    localStorage.setItem(EVAL_COLUMNS_KEY, JSON.stringify(identities))
  } catch {
    // localStorage not available
  }
}

// Dialog visibility is controlled by the composable
const dialogVisible = computed({
  get: () => visible.value,
  set: val => {
    if (!val) close()
  },
})

// Determine which API to use based on resource string
const api = computed(() => {
  if (!options.value?.resource) return null
  const resource = options.value.resource
  const match = resource.match(/^corteza::(\w+):/)
  if (!match) return null
  const component = match[1]
  switch (component) {
    case 'system':
      return $SystemAPI
    case 'compose':
      return $ComposeAPI
    case 'automation':
      return $AutomationAPI
    default:
      return $SystemAPI
  }
})

// Dialog title
const dialogTitle = computed(() => {
  if (!options.value) return ''
  const { resource, title, allSpecific } = options.value

  const i18nPrefix = getI18nPrefix(resource)

  let target
  if (allSpecific && title) {
    const key = `permissions.resources.${i18nPrefix}.all-specific`
    target = te(key) ? t(key, { target: title }) : title
  } else if (title) {
    const key = `permissions.resources.${i18nPrefix}.specific`
    target = te(key) ? t(key, { target: title }) : title
  } else {
    const key = `permissions.resources.${i18nPrefix}.all`
    target = te(key) ? t(key) : i18nPrefix
  }

  return t('permissions.ui.set-for', { target })
})

// Dirty check
const dirty = computed(() => {
  return rules.value.some(r => r.access !== r.initial)
})

// Add eval enabled check
const addEvalEnabled = computed(() => {
  return addEval.value.roleIDs?.length > 0 || addEval.value.userID
})

// Re-run on every open — options is a fresh object each time, so opening the
// same resource for a different preselected role still refreshes the editor.
watch(
  () => options.value,
  async opts => {
    if (opts?.resource && api.value) {
      processing.value = true
      currentRoleID.value = null
      rules.value = []
      initialRules.value = {}
      evaluate.value = []
      try {
        // The operation catalog first: both the editor's rows and a restored
        // column's accesses are built from it.
        await fetchPermissions()
        if (opts.roleID) {
          currentRoleID.value = String(opts.roleID)
        }
        await Promise.all([
          opts.roleID ? fetchRules(currentRoleID.value) : selectInitialRole(),
          restoreEvalColumns(),
        ])
      } finally {
        processing.value = false
      }
    }
  },
)

// Helper functions
function getI18nPrefix(resource) {
  if (!resource) return ''
  const [tmp = ''] = resource.split('/', 2)
  const parts = tmp.split(':').filter(Boolean)
  const component = parts[1] || ''
  const resourceType = parts[2] || 'component'
  return `${component}.${resourceType}`
}

function getResourceType(resource) {
  if (!resource) return ''
  const [type] = resource.split('/', 2)
  return type
}

function humanize(str) {
  return str
    .replace(/-/g, ' ')
    .replace(/\.\w/g, m => ' ' + m[1])
    .replace(/^./, m => m.toUpperCase())
}

function describePermission(operation) {
  const resource = options.value?.resource
  const i18nPrefix = getI18nPrefix(resource)
  const target = options.value?.target || options.value?.title || ''
  const allSpecific = options.value?.allSpecific || false

  const titleKey = `permissions.resources.${i18nPrefix}.operations.${operation}.title`
  const descKey = `permissions.resources.${i18nPrefix}.operations.${operation}.description`

  let title
  if (allSpecific && target) {
    const specificKey = `permissions.resources.${i18nPrefix}.operations.${operation}.all-specific`
    title = te(specificKey)
      ? t(specificKey, { target })
      : te(titleKey)
        ? t(titleKey)
        : humanize(operation)
  } else if (target) {
    const specificKey = `permissions.resources.${i18nPrefix}.operations.${operation}.specific`
    title = te(specificKey)
      ? t(specificKey, { target })
      : te(titleKey)
        ? t(titleKey)
        : humanize(operation)
  } else {
    title = te(titleKey) ? t(titleKey) : humanize(operation)
  }

  const description = te(descKey) ? t(descKey) : ''

  return { title, description }
}

// Restore the last-picked role; fall back to the first available role when
// nothing is remembered or the remembered role no longer resolves (a successful
// fetch always yields one rule per operation, so empty rules means the read
// failed — e.g. the role was deleted).
async function selectInitialRole() {
  const remembered = readLastRoleID()
  if (remembered) {
    currentRoleID.value = remembered
    await fetchRules(remembered)
    if (rules.value.length) return
    currentRoleID.value = null
  }
  await autoSelectFirstRole()
}

// Auto-select first available role
async function autoSelectFirstRole() {
  if (!$SystemAPI) return
  try {
    const result = await $SystemAPI.roleList({ query: '', limit: 10 })
    const excluded = ['super-admin', 'Super administrator']
    const firstRole = (result.set || []).find(
      r => !excluded.includes(r.handle) && !excluded.includes(r.name),
    )
    if (firstRole) {
      currentRoleID.value = firstRole.roleID
      await fetchRules(firstRole.roleID)
    }
  } catch {
    // silent
  }
}

// API calls
async function fetchPermissions() {
  if (!api.value) return
  try {
    const pp = await api.value.permissionsList()
    const resourceType = getResourceType(options.value.resource)
    permissions.value = (pp || []).filter(({ type }) => resourceType === type)
  } catch {
    permissions.value = []
  }
}

async function fetchRules(roleID) {
  if (!api.value || !roleID) return

  try {
    const resource = options.value.resource
    const rr = await api.value.permissionsRead({ roleID, resource })

    rules.value = permissions.value.map(({ op: operation }) => {
      const found = (rr || []).find(r => r.operation === operation)
      const access = found?.access || 'inherit'
      const { title, description } = describePermission(operation)

      return {
        operation,
        resource: options.value.resource,
        access,
        initial: access,
        title,
        description,
      }
    })

    initialRules.value = Object.fromEntries(rules.value.map(r => [r.operation, r.access]))
  } catch {
    rules.value = []
  }
}

// Evaluation functions
async function evaluatePermissions({ roleID, userID }) {
  if (!api.value) return []
  try {
    const resource = options.value.resource
    const result = await api.value.permissionsTrace({ resource, roleID, userID })
    return normalizeEvalRules(result)
  } catch {
    return []
  }
}

function normalizeEvalRules(rr) {
  return permissions.value.map(({ op: operation }) => {
    const found = (rr || []).find(r => r.operation === operation)
    let access = 'deny'
    if (found) {
      if (found.resolution === 'unknown-context') {
        access = 'unknown-context'
      } else {
        access = found.access || 'deny'
      }
    }
    return { operation, access }
  })
}

function getEvalAccess(evalCol, operation) {
  const rule = (evalCol.rules || []).find(r => r.operation === operation)
  return rule?.access || 'deny'
}

// `ch` is the width of a "0", so this reads a touch wide for narrow text and a
// touch narrow for wide — near enough between a floor and a ceiling, and the
// names truncate with the full one on hover either way. The 2rem covers the
// cell's own padding.
function evalColumnStyle(evalCol) {
  const longest = getEvalName(evalCol).reduce((n, name) => Math.max(n, String(name).length), 0)
  return { width: `clamp(${EVAL_COLUMN_MIN}, calc(${longest}ch + 2rem), ${EVAL_COLUMN_MAX})` }
}

function getEvalName(evalCol) {
  if (evalCol.userName) return [evalCol.userName]
  if (evalCol.roleNames?.length) return evalCol.roleNames
  // A name that would not resolve leaves the column headerless; the raw ID at
  // least says which one it is.
  if (evalCol.userID) return [evalCol.userID]
  return evalCol.roleIDs?.length ? evalCol.roleIDs : ['Unknown']
}

// A column, traced against the resource the dialog is currently open for.
// Shared by adding one and by restoring a stored one, so both resolve names
// and accesses the same way.
async function buildEvalColumn({ roleIDs = [], userID = null }) {
  const roleIDList = Array.isArray(roleIDs)
    ? roleIDs.map(r => (typeof r === 'object' ? r.roleID : r))
    : []

  const [rules, roleNames, userName] = await Promise.all([
    evaluatePermissions({ roleID: roleIDList, userID }),
    readRoleNames(roleIDList),
    readUserName(userID),
  ])

  return { roleIDs: roleIDList, userID, roleNames, userName, rules }
}

async function readRoleNames(roleIDList) {
  if (!roleIDList.length || !$SystemAPI) return []
  const results = await Promise.all(
    roleIDList.map(id => $SystemAPI.roleRead({ roleID: id }).catch(() => null)),
  )
  return results.filter(Boolean).map(r => r.name || r.handle || r.roleID)
}

async function readUserName(userID) {
  if (!userID || !$SystemAPI) return null
  const user = await $SystemAPI.userRead({ userID }).catch(() => null)
  return user ? user.name || user.username || user.email || user.userID : null
}

// Whether the people a column compares against still exist.
function evalColumnResolves(col) {
  return col.userID ? !!col.userName : col.roleNames.length > 0
}

async function restoreEvalColumns() {
  const stored = readStoredEvalColumns()
  if (!stored.length) return

  const built = await Promise.all(stored.map(c => buildEvalColumn(c).catch(() => null)))
  evaluate.value = built.filter(c => c && evalColumnResolves(c))

  // A role or user that has since been deleted shouldn't keep being asked for.
  if (evaluate.value.length !== stored.length) storeEvalColumns()
}

function removeEvalColumn(index) {
  evaluate.value.splice(index, 1)
  storeEvalColumns()
}

async function onAddEvalColumn() {
  evaluate.value.push(await buildEvalColumn(addEval.value))
  storeEvalColumns()

  // Reset add form
  addEval.value = { roleIDs: [], userID: null }
  showAddEval.value = false
}

function onRoleChange(role) {
  if (role?.roleID) {
    storeLastRoleID(role.roleID)
    fetchRules(role.roleID)
  }
}

function setAccess(rule, access) {
  rule.access = access
}

async function onSubmit() {
  if (!api.value || !currentRoleID.value) return

  submitting.value = true
  try {
    const roleID = currentRoleID.value
    const changedRules = rules.value
      .filter(r => r.access !== r.initial)
      .map(({ resource, operation, access }) => ({
        resource,
        operation,
        access,
      }))

    if (changedRules.length > 0) {
      await api.value.permissionsUpdate({ roleID, rules: changedRules })
    }

    $toast.toastSuccess(t('permissions.ui.notification.save.success'))

    // Update initial state to clear dirty flag
    rules.value.forEach(r => {
      r.initial = r.access
    })

    // Re-evaluate all columns after save
    await reEvaluateAll()

    // Notify the opener so dependent views (e.g. a resource graph) can refresh.
    options.value?.onSaved?.()
  } catch (e) {
    $toast.toastErrorHandler(t('permissions.ui.notification.save.failed'))(e)
  } finally {
    submitting.value = false
  }
}

async function reEvaluateAll() {
  const updated = await Promise.all(
    evaluate.value.map(async e => {
      const evalRules = await evaluatePermissions({
        roleID: e.roleIDs || [],
        userID: e.userID,
      })
      return { ...e, rules: evalRules }
    }),
  )
  evaluate.value = updated
}

function onHide() {
  close()
  currentRoleID.value = null
  rules.value = []
  permissions.value = []
  initialRules.value = {}
  evaluate.value = []
}
</script>

<style scoped>
.perm-col-left {
  flex: 0 0 42%;
  max-width: 42%;
  min-width: 42%;
}

.eval-remove-icon {
  opacity: 0;
  transition: opacity 0.2s;
}
div:hover > .eval-remove-icon {
  opacity: 1;
}
</style>
