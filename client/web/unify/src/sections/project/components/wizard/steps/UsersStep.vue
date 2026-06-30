<template>
  <div class="h-full overflow-auto p-4">
    <CFormGroup :label="$t('project.accessUsers.title')">
      <template #actions>
        <Button
          v-if="!disabled && roles.length"
          icon="pi pi-plus"
          :label="$t('project.accessUsers.add')"
          severity="secondary"
          size="small"
          @click="openAdd"
        />
      </template>

      <!-- Need roles before users can be assigned to them. -->
      <div
        v-if="!roles.length"
        class="rounded-lg border border-surface mt-1 px-4 py-8 text-center text-muted-color"
      >
        <i class="pi pi-id-card text-3xl mb-2" />
        <p class="text-sm">{{ $t('project.accessUsers.noRoles') }}</p>
      </div>

      <div v-else class="rounded-lg border border-surface overflow-x-auto mt-1">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-emphasis text-muted-color text-xs uppercase tracking-wider">
              <th class="text-left font-medium px-3 py-2">{{ $t('project.accessUsers.columns.user') }}</th>
              <th class="text-left font-medium px-3 py-2">{{ $t('project.accessUsers.columns.roles') }}</th>
              <th class="px-3 py-2" />
            </tr>
          </thead>
          <tbody class="divide-y divide-surface">
            <tr v-for="row in rows" :key="row.userId" class="group hover:bg-emphasis/50">
              <td class="px-3 py-2 align-top">
                <div class="leading-tight">{{ row.name }}</div>
                <div class="text-xs text-muted-color">{{ row.email }}</div>
              </td>
              <td class="px-3 py-2">
                <div class="flex flex-wrap gap-1.5">
                  <button
                    v-for="role in roles"
                    :key="role.id"
                    type="button"
                    :disabled="disabled || busy"
                    class="px-2 py-0.5 rounded-md text-xs font-medium border transition-colors"
                    :class="[
                      row.roleIds.includes(role.id)
                        ? 'bg-primary border-primary text-primary-contrast'
                        : 'bg-transparent border-surface text-muted-color hover:border-primary',
                      { 'opacity-60 cursor-not-allowed': disabled || busy },
                    ]"
                    @click="!disabled && !busy && toggleRole(row, role.id)"
                  >
                    {{ role.name }}
                  </button>
                </div>
              </td>
              <td class="text-right px-3 py-2 align-top">
                <Button
                  v-if="!disabled"
                  icon="pi pi-times"
                  severity="secondary"
                  text
                  rounded
                  size="small"
                  class="opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity"
                  @click="remove(row)"
                />
              </td>
            </tr>
            <tr v-if="!rows.length">
              <td colspan="3" class="px-3 py-4 text-center text-muted-color italic">
                {{ $t('project.accessUsers.empty') }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </CFormGroup>

    <Dialog
      v-model:visible="addOpen"
      modal
      :header="$t('project.accessUsers.addDialog.header')"
      :style="{ width: '32rem' }"
      :pt="{ footer: { class: 'flex justify-end gap-2' } }"
    >
      <div class="flex flex-col gap-4">
        <SelectButton
          v-model="draft.mode"
          :options="modeOptions"
          option-label="label"
          option-value="value"
          :allow-empty="false"
          fluid
        />

        <template v-if="draft.mode === 'existing'">
          <CFormGroup :label="$t('general.label.user.single')" required>
            <Select
              v-model="draft.userId"
              :options="availableUsers"
              option-label="name"
              option-value="id"
              filter
              fluid
              :placeholder="$t('project.accessUsers.addDialog.userPlaceholder')"
            >
              <template #option="{ option }">
                <div class="leading-tight">
                  <div>{{ option.name }}</div>
                  <div class="text-xs text-muted-color">{{ option.email }}</div>
                </div>
              </template>
            </Select>
          </CFormGroup>
        </template>
        <template v-else>
          <CFormGroup :label="$t('project.accessUsers.addDialog.name')">
            <InputText v-model="draft.name" fluid :placeholder="$t('project.accessUsers.addDialog.namePlaceholder')" />
          </CFormGroup>
          <CFormGroup :label="$t('project.accessUsers.addDialog.email')" required>
            <InputText v-model="draft.email" fluid :placeholder="$t('project.accessUsers.addDialog.emailPlaceholder')" />
          </CFormGroup>
        </template>

        <CFormGroup :label="$t('project.accessUsers.addDialog.roles')" required>
          <MultiSelect
            v-model="draft.roleIds"
            :options="roles"
            option-label="name"
            option-value="id"
            filter
            display="chip"
            fluid
            :placeholder="$t('project.accessUsers.addDialog.rolePlaceholder')"
          />
        </CFormGroup>
      </div>
      <template #footer>
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="addOpen = false"
        />
        <Button
          :label="$t('general.label.add')"
          size="small"
          :disabled="!addEnabled || busy"
          @click="add"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const { t } = useI18n()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

const roles = computed(() => store.rolesFor(props.project.id))
const projectUsers = computed(() => store.projectUsersFor(props.project.id))

const busy = ref(false)
const addOpen = ref(false)
const draft = reactive({ mode: 'existing', userId: null, name: '', email: '', roleIds: [] })

const modeOptions = [
  { label: t('project.accessUsers.mode.existing'), value: 'existing' },
  { label: t('project.accessUsers.mode.new'), value: 'new' },
]

const rows = computed(() =>
  projectUsers.value.map(u => {
    const dir = usersStore.findUser(u.userId)
    return {
      userId: u.userId,
      roleIds: u.roleIds,
      name: dir?.name || u.userId,
      email: dir?.email || '',
    }
  }),
)

const availableUsers = computed(() => {
  const taken = new Set(projectUsers.value.map(u => u.userId))
  return usersStore.users.filter(u => !taken.has(u.id))
})

const addEnabled = computed(() => {
  if (!draft.roleIds.length) return false
  return draft.mode === 'existing' ? !!draft.userId : !!draft.email.trim()
})

function openAdd() {
  draft.mode = 'existing'
  draft.userId = null
  draft.name = ''
  draft.email = ''
  draft.roleIds = []
  addOpen.value = true
}

async function toggleRole(row, roleId) {
  busy.value = true
  try {
    await store.setProjectUserRole(props.project.id, row.userId, roleId, !row.roleIds.includes(roleId))
  } catch (err) {
    $toast.toastErrorHandler(t('project.accessUsers.toastUpdateFailed'))(err)
  } finally {
    busy.value = false
  }
}

async function add() {
  if (!addEnabled.value) return
  busy.value = true
  try {
    let userId = draft.userId
    if (draft.mode === 'new') {
      userId = await store.addProjectUser(props.project.id, { email: draft.email, name: draft.name })
      await usersStore.reload()
    }
    await store.assignProjectUserRoles(props.project.id, userId, draft.roleIds)
    addOpen.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.accessUsers.toastAddFailed'))(err)
  } finally {
    busy.value = false
  }
}

function remove(row) {
  confirmDelete({
    header: t('project.accessUsers.removeConfirm.header'),
    message: t('project.accessUsers.removeConfirm.message', { name: row.name }),
    onConfirm: async () => {
      try {
        await store.removeProjectUser(props.project.id, row.userId, row.roleIds)
      } catch (err) {
        $toast.toastErrorHandler(t('project.accessUsers.toastRemoveFailed'))(err)
      }
    },
  })
}

async function refresh(id) {
  if (!id) return
  try {
    await Promise.all([usersStore.load(), store.loadProjectUsers(id)])
  } catch (err) {
    $toast.toastErrorHandler(t('project.accessUsers.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.id))
watch(() => props.project.id, refresh)
</script>
