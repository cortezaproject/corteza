<template>
  <div v-if="!loaded" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else-if="!canGrant" class="flex items-center justify-center h-full">
    <Message severity="error" :closable="false">
      {{ $t('permissions.ui.not-allowed') }}
    </Message>
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="p-4 flex-1 flex flex-col min-h-0">
      <Card
        class="shadow flex-1 flex flex-col min-h-0"
        :pt="{
          root: { class: 'flex-1 flex flex-col min-h-0 overflow-hidden' },
          body: { class: 'flex-1 flex flex-col min-h-0', style: 'padding: 0' },
          content: { class: 'flex-1 flex flex-col min-h-0', style: 'padding: 0' },
        }"
      >
        <template #content>
          <!-- Permission grid -->
          <div class="flex-1 overflow-y-auto min-h-0">
            <!-- Header with roles (sticky inside scroll container) -->
            <div class="flex border-b bg-surface sticky top-0 z-20">
              <div class="w-1/3 p-3 text-sm text-muted-color">
                {{ $t('permissions.ui.click-on-cell-to-allow') }}
              </div>
              <div class="flex">
                <div
                  v-for="role in roles"
                  :key="role.ID"
                  class="w-48 flex flex-col items-center justify-center p-3 border-l cursor-pointer group"
                  @click="hideRole(role)"
                >
                  <span
                    v-for="(n, index) in role.name"
                    :key="index"
                    :title="n"
                    class="text-center text-primary text-sm font-medium truncate max-w-full"
                  >
                    {{ n }}
                  </span>
                  <span
                    class="text-xs mt-1"
                    :class="role.mode === 'edit' ? 'text-primary' : 'text-muted-color'"
                  >
                    {{ $t(`permissions.ui.${role.mode === 'edit' ? 'edit' : 'evaluate'}.title`) }}
                  </span>

                  <i
                    class="pi pi-times text-muted-color text-xs mt-1 opacity-0 group-hover:opacity-100 transition-opacity"
                  />
                </div>

                <div
                  v-if="roles.length < 8"
                  class="w-32 flex items-center justify-center gap-2 p-3 border-l cursor-pointer hover:bg-emphasis"
                  @click="addDialogVisible = true"
                >
                  <i class="pi pi-plus text-primary" />
                  <span class="text-primary text-sm font-medium">
                    {{ $t('permissions.ui.add.label') }}
                  </span>
                </div>
              </div>
            </div>

            <div v-for="type in sortedTypes" :key="type">
              <!-- Resource type header -->
              <div class="flex border-b bg-emphasis sticky top-0 z-10">
                <div class="w-1/3 p-3 text-sm font-semibold text-primary">
                  {{ getTranslation(type) }}
                </div>
                <div class="flex">
                  <div v-for="role in roles" :key="role.ID" class="w-48 border-l p-3" />
                  <div v-if="roles.length < 8" class="w-32 border-l p-3" />
                </div>
              </div>

              <!-- Operation rows -->
              <div
                v-for="operation in permissions[type].ops"
                :key="`${type}-${operation}`"
                class="flex border-b"
              >
                <div class="w-1/3 p-3 text-sm text-color" :title="getTranslation(type, operation)">
                  {{ getTranslation(type, operation) }}
                </div>
                <div class="flex">
                  <div
                    v-for="role in roles"
                    :key="role.ID"
                    class="w-48 flex items-center justify-center border-l p-3 text-lg"
                    :class="{
                      'cursor-pointer hover:bg-emphasis': role.mode === 'edit',
                      'cursor-not-allowed bg-emphasis': role.mode === 'eval',
                      'bg-highlight': checkChange(role.ID, permissions[type].any, operation),
                    }"
                    :title="
                      getRuleTooltip(
                        checkRule(role.ID, permissions[type].any, operation, 'unknown-context'),
                        !!role.userID,
                      )
                    "
                    @click="
                      role.mode === 'edit'
                        ? ruleChange(role.ID, permissions[type].any, operation)
                        : undefined
                    "
                  >
                    <i
                      v-if="checkRule(role.ID, permissions[type].any, operation, 'unknown-context')"
                      class="pi pi-question text-muted-color"
                    />
                    <i
                      v-else-if="checkRule(role.ID, permissions[type].any, operation, 'allow')"
                      class="pi pi-check text-primary"
                    />
                    <i v-else class="pi pi-times text-muted-color" />
                  </div>

                  <div v-if="roles.length < 8" class="w-32 border-l p-3" />
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-end">
        <Button
          :label="$t('permissions.ui.save')"
          icon="pi pi-save"
          :loading="saving"
          @click="onSubmit"
        />
      </div>
    </div>

    <!-- Add role dialog -->
    <Dialog
      v-model:visible="addDialogVisible"
      :header="$t('permissions.ui.edit-or-eval')"
      modal
      :style="{ width: '50rem' }"
      :breakpoints="{ '50rem': '94vw' }"
    >
      <div class="flex flex-col gap-4">
        <SelectButton
          v-model="add.mode"
          :options="modeOptions"
          option-label="label"
          option-value="value"
          fluid
        />

        <p class="text-sm text-muted-color">
          {{
            add.mode === 'edit'
              ? $t('permissions.ui.add.edit.description')
              : $t('permissions.ui.add.evaluate.description')
          }}
        </p>

        <div class="flex flex-col gap-1">
          <label class="font-medium text-sm text-primary">
            {{ $t('permissions.ui.add.role.label') }}
          </label>
          <CInputRole
            v-model="add.roleID"
            :placeholder="$t('permissions.ui.add.role.placeholder')"
            :disabled="add.mode === 'eval' && !!add.userID"
            :clear-on-select="add.mode === 'eval'"
            @select="onRoleSelected"
          />

          <!-- Show selected roles as chips in eval mode -->
          <div
            v-if="add.mode === 'eval' && add.selectedRoles.length"
            class="flex flex-wrap gap-1 mt-1"
          >
            <Chip
              v-for="role in add.selectedRoles"
              :key="role.roleID"
              :label="role.name || role.handle || role.roleID"
              removable
              @remove="removeSelectedRole(role.roleID)"
            />
          </div>
        </div>

        <div v-if="add.mode === 'eval'" class="flex flex-col gap-1">
          <label class="font-medium text-sm text-primary">
            {{ $t('permissions.ui.add.user.label') }}
          </label>
          <AutoComplete
            v-model="add.userID"
            :suggestions="userSuggestions"
            :placeholder="$t('permissions.ui.add.user.placeholder')"
            option-label="label"
            :disabled="add.selectedRoles.length > 0"
            class="w-full"
            dropdown
            @complete="searchUsers"
          />
        </div>
      </div>

      <template #footer>
        <Button :label="$t('permissions.ui.add.save')" :disabled="!addEnabled" @click="onAdd" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'
import { kebabCase } from 'lodash-es'

const { CInputRole } = components

const { t, te } = useI18n()

const props = defineProps({
  api: {
    type: Object,
    required: true,
  },
  component: {
    type: String,
    required: true,
  },
})

const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const LS_KEY = 'permissionList.roles'

// State
const loaded = ref(false)
const saving = ref(false)
const canGrant = ref(false)
const permissions = ref({})
const resources = ref(new Set())
const roles = ref([])
const rolePermissions = ref([])
const permissionChanges = ref([])

// Add dialog
const addDialogVisible = ref(false)
const add = ref({ mode: 'edit', roleID: null, userID: null, selectedRole: null, selectedRoles: [] })
const userSuggestions = ref([])
const fetchedUsers = ref({})

const modeOptions = [
  { label: t('permissions.ui.edit.title'), value: 'edit' },
  { label: t('permissions.ui.evaluate.description'), value: 'eval' },
]

// Computed
const sortedTypes = computed(() => Object.keys(permissions.value).sort())

const addEnabled = computed(() => {
  const { mode, selectedRole, selectedRoles, userID } = add.value
  if (mode === 'edit') return !!selectedRole
  if (mode === 'eval') return selectedRoles.length > 0 || !!userID
  return false
})

// Watch mode changes to reset
watch(
  () => add.value.mode,
  () => {
    add.value.roleID = null
    add.value.userID = null
    add.value.selectedRole = null
    add.value.selectedRoles = []
  },
)

function onRoleSelected(role) {
  if (add.value.mode === 'edit') {
    // In edit mode, store the full role object
    add.value.selectedRole = role
  } else {
    // In eval mode, accumulate roles
    if (!add.value.selectedRoles.some(r => r.roleID === role.roleID)) {
      add.value.selectedRoles.push(role)
    }
  }
}

function removeSelectedRole(roleID) {
  add.value.selectedRoles = add.value.selectedRoles.filter(r => r.roleID !== roleID)
}

// LocalStorage helpers
function getIncludedRoles() {
  try {
    return JSON.parse(localStorage.getItem(LS_KEY) || '[]')
  } catch {
    return []
  }
}

function setIncludedRoles(rr) {
  const filtered = rr.filter(r => !['1', '2'].includes(String(r.roleID)))
  localStorage.setItem(LS_KEY, JSON.stringify(filtered))
}

// Permission checks
function checkRule(ID, res, op, access) {
  const key = `${op}@${res}`
  const rp = rolePermissions.value.find(r => r.ID === ID)
  return rp ? rp.rules[key] === access : false
}

function checkChange(ID, res, op) {
  const key = `${op}@${res}`
  const current = (rolePermissions.value.find(r => r.ID === ID) || { rules: {} }).rules[key]
  const initial = (permissionChanges.value.find(r => r.ID === ID) || { rules: {} }).rules?.[key]
  return initial ? current !== initial : false
}

function ruleChange(ID, res, op) {
  const key = `${op}@${res}`
  const rp = rolePermissions.value.find(r => r.ID === ID)
  if (!rp) return

  let access = rp.rules[key]

  // Track initial value
  let pc = permissionChanges.value.find(r => r.ID === ID)
  if (!pc) {
    pc = { ID, rules: {} }
    permissionChanges.value.push(pc)
  }
  if (!pc.rules[key]) {
    pc.rules[key] = access || 'inherit'
  }

  // Toggle
  rp.rules[key] = access === 'allow' ? 'inherit' : 'allow'
}

function roleRules(rules, mode = 'edit') {
  return (rules || []).reduce((map, { resource, operation, access, resolution }) => {
    const [type] = resource.split('/', 2)
    const p = permissions.value[type]
    if (p && p.ops.indexOf(operation) > -1) {
      if (mode === 'eval') {
        if (resolution === 'unknown-context') access = 'unknown-context'
        else if (access === 'inherit') access = 'deny'
      }
      map[`${operation}@${resource}`] = access
    }
    return map
  }, {})
}

// Translation helpers
function humanize(str) {
  return str
    .replace(/-/g, ' ')
    .replace(/\.\w/g, m => ' ' + m[1])
    .replace(/^./, m => m.toUpperCase())
}

function getTranslation(resource, operation = '') {
  const key = kebabCase(resource.split(':')[3]) || 'component'
  if (operation) {
    const i18nKey = `permissions.resources.${props.component}.${key}.operations.${operation}.title`
    return te(i18nKey) ? t(i18nKey) : humanize(operation)
  }
  const i18nKey = `permissions.resources.${props.component}.${key}.label`
  return te(i18nKey) ? t(i18nKey) : humanize(key)
}

function getRuleTooltip(isUnknown = false, isUser = false) {
  if (!isUnknown) return ''
  return t(`permissions.ui.tooltip.unknown-context.${isUser ? 'user' : 'role'}`)
}

// API methods
async function readPermissions({ name, roleID }) {
  const resource = [...resources.value]
  try {
    const rr = await props.api.permissionsRead({ resource, roleID })
    const ID = `edit-${roleID}`
    rolePermissions.value.push({ resource: '', ID, rules: roleRules(rr, 'edit') })
    roles.value.push({ mode: 'edit', ID, roleID, name })
  } catch (e) {
    $toast.toastErrorHandler(t('permissions.ui.notification.save.failed'))(e)
  }
}

async function evaluatePermissions({ name, roleID, userID }) {
  const resource = [...resources.value]
  try {
    const rr = await props.api.permissionsTrace({ resource, roleID, userID })
    const ID = userID ? `eval-${userID}` : `eval-${roleID.join('-')}`
    rolePermissions.value.push({ resource: '', ID, rules: roleRules(rr, 'eval') })
    roles.value.push({ mode: 'eval', ID, roleID, userID, name })
  } catch (e) {
    $toast.toastErrorHandler(t('permissions.ui.notification.save.failed'))(e)
  }
}

async function prepareRoles() {
  rolePermissions.value = []
  const included = getIncludedRoles()
  await Promise.all(
    included.map(({ mode, name, roleID, userID }) => {
      if (mode === 'edit') return readPermissions({ name, roleID })
      return evaluatePermissions({ name, roleID, userID })
    }),
  )
}

async function fetchPermissions() {
  loaded.value = false
  try {
    const perms = await props.api.permissionsList()
    const res = new Set()
    permissions.value = (perms || []).reduce((map, { type, any, op }) => {
      res.add(any)
      if (!map[type]) map[type] = { any, ops: [] }
      map[type].ops.push(op)
      return map
    }, {})
    resources.value = res
    await prepareRoles()
  } catch (e) {
    $toast.toastErrorHandler(t('permissions.ui.notification.save.failed'))(e)
  } finally {
    loaded.value = true
  }
}

// Check grant permission
async function checkGrantPermission() {
  try {
    // Try to list permissions — if the API call succeeds, user can grant
    await props.api.permissionsList()
    canGrant.value = true
  } catch {
    canGrant.value = false
  }
}

// Submit
async function onSubmit() {
  saving.value = true
  try {
    const editRoles = rolePermissions.value.filter(({ ID }) => ID.includes('edit'))
    await Promise.all(
      editRoles.map(({ ID, rules }) => {
        const roleID = ID.split('-')[1]
        const externalRules = []
        Object.entries(rules).forEach(([key, value]) => {
          const [operation, resource] = key.split('@', 2)
          externalRules.push({ roleID, resource, operation, access: value })
        })
        return props.api.permissionsUpdate({ roleID, rules: externalRules })
      }),
    )

    // Re-evaluate eval roles after save
    const evalRoles = roles.value.filter(({ mode }) => mode === 'eval')
    await Promise.all(
      evalRoles.map(({ roleID, userID }) => {
        return props.api.permissionsTrace({ roleID, userID }).then(rr => {
          const ID = userID ? `eval-${userID}` : `eval-${roleID.join('-')}`
          rolePermissions.value = [
            ...rolePermissions.value.filter(rp => rp.ID !== ID),
            { resource: '', ID, rules: roleRules(rr, 'eval') },
          ]
        })
      }),
    )

    permissionChanges.value = []
    $toast.toastSuccess(t('permissions.ui.notification.save.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('permissions.ui.notification.save.failed'))(e)
  } finally {
    saving.value = false
  }
}

// Add role
function onAdd() {
  const { mode, selectedRole, selectedRoles, userID } = add.value

  if (mode === 'edit' && selectedRole) {
    const rid = selectedRole.roleID
    const name = [selectedRole.name || selectedRole.handle || rid]
    const ID = `edit-${rid}`
    if (roles.value.some(r => r.ID === ID)) {
      addDialogVisible.value = false
      return
    }
    readPermissions({ roleID: rid, name })
      .then(() => {
        setIncludedRoles(roles.value)
      })
      .catch($toast.toastErrorHandler(t('permissions.ui.notification.save.failed')))
  } else if (mode === 'eval') {
    let uid = null
    let rids = []
    let name = []

    if (userID && userID.userID) {
      uid = userID.userID
      name = [userID.label]
    } else if (selectedRoles.length) {
      rids = selectedRoles.map(r => r.roleID)
      name = selectedRoles.map(r => r.name || r.handle || r.roleID)
    }

    const ID = uid ? `eval-${uid}` : `eval-${rids.join('-')}`
    if (roles.value.some(r => r.ID === ID)) {
      addDialogVisible.value = false
      return
    }
    evaluatePermissions({ name, roleID: rids, userID: uid })
      .then(() => {
        setIncludedRoles(roles.value)
      })
      .catch($toast.toastErrorHandler(t('permissions.ui.notification.save.failed')))
  }

  add.value = { mode: 'edit', roleID: null, userID: null, selectedRole: null, selectedRoles: [] }
  addDialogVisible.value = false
}

// Hide role
function hideRole(role) {
  roles.value = roles.value.filter(r => r.ID !== role.ID)
  rolePermissions.value = rolePermissions.value.filter(r => r.ID !== role.ID)
  setIncludedRoles(roles.value)
}

// Search users for eval mode
async function searchUsers({ query }) {
  try {
    const { set } = await $SystemAPI.userList({ query: query || '', limit: 15 })
    userSuggestions.value = (set || []).map(({ userID, name, username, email }) => {
      const label = name || username || email || `<@${userID}>`
      fetchedUsers.value[userID] = label
      return { userID, label }
    })
  } catch {
    userSuggestions.value = []
  }
}

onMounted(async () => {
  await checkGrantPermission()
  if (canGrant.value) {
    await fetchPermissions()
  } else {
    loaded.value = true
  }
})
</script>
