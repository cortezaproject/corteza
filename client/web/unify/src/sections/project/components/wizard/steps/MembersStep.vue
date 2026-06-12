<template>
  <div class="flex flex-col gap-6">
    <CFormGroup
      label="Members"
      description="Assign people to a role. Each role carries fixed permissions."
    >
      <template #actions>
        <Button
          v-if="canManage"
          icon="pi pi-plus"
          label="Add member"
          severity="secondary"
          size="small"
          @click="openAdd"
        />
      </template>

      <div class="rounded-lg border border-surface overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-emphasis text-muted-color text-xs uppercase tracking-wider">
              <th class="text-left font-medium px-3 py-2">Member</th>
              <th class="text-left font-medium px-3 py-2">Role</th>
              <th class="text-center font-medium px-3 py-2">Read</th>
              <th class="text-center font-medium px-3 py-2">Write</th>
              <th class="text-center font-medium px-3 py-2">Request approval</th>
              <th class="text-center font-medium px-3 py-2">Grant approval</th>
              <th class="text-left font-medium px-3 py-2">Resources</th>
              <th class="px-3 py-2" />
            </tr>
          </thead>
          <tbody class="divide-y divide-surface">
            <tr v-for="row in rows" :key="row.id" class="hover:bg-emphasis/50">
              <td class="px-3 py-2">
                <div class="leading-tight">{{ row.user.name }}</div>
                <div class="text-xs text-muted-color">{{ row.user.email }}</div>
              </td>
              <td class="px-3 py-2 whitespace-nowrap">
                <Select
                  v-if="canManage"
                  :model-value="row.role"
                  :options="ROLE_PRESETS"
                  option-label="label"
                  option-value="id"
                  size="small"
                  @update:model-value="setRole(row, $event)"
                />
                <span v-else class="font-medium">{{ row.preset.label }}</span>
              </td>
              <td class="text-center px-3 py-2"><YesNo :on="row.capabilities.canRead" /></td>
              <td class="text-center px-3 py-2"><YesNo :on="row.capabilities.canWrite" /></td>
              <td class="text-center px-3 py-2">
                <YesNo :on="row.capabilities.canRequestApproval" />
              </td>
              <td class="text-center px-3 py-2">
                <YesNo :on="row.capabilities.canGrantApproval" />
              </td>
              <td class="px-3 py-2 text-muted-color">{{ row.preset.resources }}</td>
              <td class="text-right px-3 py-2">
                <Button
                  v-if="canManage"
                  icon="pi pi-times"
                  severity="secondary"
                  text
                  rounded
                  size="small"
                  @click="remove(row)"
                />
              </td>
            </tr>
            <tr v-if="!rows.length">
              <td colspan="8" class="px-3 py-4 text-center text-muted-color italic">
                No members yet.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </CFormGroup>

    <!-- Role descriptions — the accountability framework (Article 17(m)). -->
    <CFormGroup>
      <div class="rounded-lg border border-surface divide-y divide-surface">
        <div v-for="r in ROLE_PRESETS" :key="r.id" class="px-4 py-3">
          <div class="flex items-center gap-2">
            <span class="font-medium">{{ r.label }}</span>
            <span class="text-xs text-muted-color">{{ r.resources }}</span>
          </div>
          <p class="text-sm text-muted-color leading-relaxed mt-1">{{ r.description }}</p>
        </div>
      </div>
    </CFormGroup>

    <Dialog
      v-model:visible="addOpen"
      modal
      header="Add member"
      :style="{ width: '30rem' }"
      :pt="{ footer: { class: 'flex justify-end gap-2' } }"
    >
      <div class="flex flex-col gap-4">
        <CFormGroup label="Role" required>
          <Select
            v-model="draft.role"
            :options="ROLE_PRESETS"
            option-label="label"
            option-value="id"
            fluid
          />
        </CFormGroup>
        <CFormGroup label="User" required>
          <Select
            v-model="draft.userId"
            :options="availableUsers"
            option-label="name"
            option-value="id"
            filter
            fluid
            placeholder="Select a system user"
          >
            <template #option="{ option }">
              <div class="leading-tight">
                <div>{{ option.name }}</div>
                <div class="text-xs text-muted-color">{{ option.email }}</div>
              </div>
            </template>
          </Select>
        </CFormGroup>
      </div>
      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          outlined
          size="small"
          @click="addOpen = false"
        />
        <Button label="Add" size="small" :disabled="!draft.role || !draft.userId" @click="add" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import YesNo from '@/sections/project/components/wizard/YesNo.vue'
import { ROLE_PRESETS, rolePreset } from '@/sections/project/config/roles'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { useToast } from 'primevue/usetoast'
import { computed, reactive, ref } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const toast = useToast()
usersStore.load()

// Mutations need the members.manage RBAC permission (the `disabled` prop adds
// the gate lock on top); the backend enforces the same rule.
const canManage = computed(() => !props.disabled && props.project.canManageMembers)

const members = computed(() => props.project.members || [])
// Capability flags come from the backend (derived from the role preset);
// `resources` is FE-only descriptive text.
const rows = computed(() =>
  members.value.map(m => ({
    ...m,
    preset: rolePreset(m.role),
    user: usersStore.findUser(m.userId) || { name: m.userId, email: '' },
  })),
)

const addOpen = ref(false)
const draft = reactive({ role: 'developer', userId: null })

const availableUsers = computed(() => {
  const taken = new Set(members.value.map(m => m.userId))
  return usersStore.users.filter(u => !taken.has(u.id))
})

function openAdd() {
  draft.role = 'developer'
  draft.userId = null
  addOpen.value = true
}

function fail(summary, err) {
  toast.add({ severity: 'error', summary, detail: err.message, life: 4000 })
}

async function add() {
  if (!draft.role || !draft.userId) return
  try {
    await store.addMember(props.project.id, { userId: draft.userId, role: draft.role })
    addOpen.value = false
  } catch (err) {
    fail('Could not add member', err)
  }
}

async function setRole(row, role) {
  if (!role || role === row.role) return
  try {
    await store.updateMember(props.project.id, row.userId, role)
  } catch (err) {
    fail('Could not change role', err)
  }
}

async function remove(row) {
  try {
    await store.removeMember(props.project.id, row.userId)
  } catch (err) {
    fail('Could not remove member', err)
  }
}
</script>
