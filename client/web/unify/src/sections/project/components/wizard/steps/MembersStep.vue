<template>
  <div class="flex flex-col gap-6">
    <p class="text-sm text-muted-color leading-relaxed">
      Project Members are human beings with roles to either build part or all of the project and its
      AI systems, or manage project governance and make approvals at governance gates. Each role's
      responsibilities form the accountability framework for the project and are set out as Role
      Descriptions below.
    </p>

    <CFormGroup
      label="Members"
      description="Assign people to a role. Each role carries fixed permissions; tab access (Build vs Governance) follows the approval flags."
    >
      <template #actions>
        <Button
          v-if="!disabled"
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
              <th class="text-left font-medium px-3 py-2">Backup</th>
              <th class="px-3 py-2" />
            </tr>
          </thead>
          <tbody class="divide-y divide-surface">
            <tr v-for="row in rows" :key="row.id" class="hover:bg-emphasis/50">
              <td class="px-3 py-2">
                <div class="leading-tight">{{ row.user.name }}</div>
                <div class="text-xs text-muted-color">{{ row.user.email }}</div>
              </td>
              <td class="px-3 py-2 font-medium whitespace-nowrap">{{ row.preset.label }}</td>
              <td class="text-center px-3 py-2"><YesNo :on="row.preset.read" /></td>
              <td class="text-center px-3 py-2"><YesNo :on="row.preset.write" /></td>
              <td class="text-center px-3 py-2"><YesNo :on="row.preset.requestApproval" /></td>
              <td class="text-center px-3 py-2"><YesNo :on="row.preset.grantApproval" /></td>
              <td class="px-3 py-2 text-muted-color">{{ row.preset.resources }}</td>
              <td class="px-3 py-2 text-muted-color">{{ row.backup?.name || '—' }}</td>
              <td class="text-right px-3 py-2">
                <Button
                  v-if="!disabled"
                  icon="pi pi-times"
                  severity="secondary"
                  text
                  rounded
                  size="small"
                  @click="store.removeMember(project.id, row.id)"
                />
              </td>
            </tr>
            <tr v-if="!rows.length">
              <td colspan="9" class="px-3 py-4 text-center text-muted-color italic">
                No members yet.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </CFormGroup>

    <!-- Role descriptions — the accountability framework (Article 17(m)). -->
    <CFormGroup label="Role descriptions">
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
import { MOCK_USERS, findUser } from '@/sections/project/mock/users'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, reactive, ref } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()

const members = computed(() => props.project.members || [])
const rows = computed(() =>
  members.value.map(m => ({
    id: m.id,
    preset: rolePreset(m.role),
    user: findUser(m.userId) || { name: m.userId, email: '' },
    backup: m.backupUserId ? findUser(m.backupUserId) : null,
  })),
)

const addOpen = ref(false)
const draft = reactive({ role: 'developer', userId: null })

const availableUsers = computed(() => {
  const taken = new Set(members.value.map(m => m.userId))
  return MOCK_USERS.filter(u => !taken.has(u.id))
})

function openAdd() {
  draft.role = 'developer'
  draft.userId = null
  addOpen.value = true
}

function add() {
  if (!draft.role || !draft.userId) return
  store.addMember(props.project.id, { userId: draft.userId, role: draft.role })
  addOpen.value = false
}
</script>
