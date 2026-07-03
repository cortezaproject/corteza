<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="$t('project.userCreate.title')"
    :style="{ width: '32rem' }"
    :pt="{ footer: { class: 'flex justify-end gap-2' } }"
  >
    <p class="text-sm text-muted-color mb-4">
      {{ $t('project.userCreate.blurb') }}
    </p>

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
            :placeholder="$t('project.userCreate.userPlaceholder')"
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
        <CFormGroup :label="$t('project.userCreate.name')">
          <InputText
            v-model="draft.name"
            fluid
            :placeholder="$t('project.userCreate.namePlaceholder')"
          />
        </CFormGroup>
        <CFormGroup :label="$t('project.userCreate.email')" required>
          <InputText
            v-model="draft.email"
            fluid
            :placeholder="$t('project.userCreate.emailPlaceholder')"
          />
        </CFormGroup>
      </template>

      <CFormGroup :label="$t('project.userCreate.roles')" required>
        <MultiSelect
          v-model="draft.roleIds"
          :options="roles"
          option-label="name"
          option-value="id"
          filter
          display="chip"
          fluid
          :placeholder="$t('project.userCreate.rolePlaceholder')"
        />
      </CFormGroup>
    </div>

    <template #footer>
      <Button
        :label="$t('general.label.cancel')"
        severity="secondary"
        text
        size="small"
        @click="visible = false"
      />
      <Button
        :label="$t('general.label.add')"
        size="small"
        :loading="saving"
        :disabled="!addEnabled || saving"
        @click="onCreate"
      />
    </template>
  </Dialog>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
})

const emit = defineEmits(['update:modelValue', 'created'])

const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const { t } = useI18n()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const roles = computed(() => store.rolesFor(props.project.id))
const projectUsers = computed(() => store.projectUsersFor(props.project.id))

const modeOptions = [
  { label: t('project.userCreate.modeExisting'), value: 'existing' },
  { label: t('project.userCreate.modeNew'), value: 'new' },
]

// Existing users already on the project can't be added again.
const availableUsers = computed(() => {
  const taken = new Set(projectUsers.value.map(u => u.userId))
  return usersStore.users.filter(u => !taken.has(u.id))
})

const draft = reactive({ mode: 'existing', userId: null, name: '', email: '', roleIds: [] })

// Reseed the (minimal) form on every open.
watch(
  () => props.modelValue,
  open => {
    if (!open) return
    draft.mode = 'existing'
    draft.userId = null
    draft.name = ''
    draft.email = ''
    draft.roleIds = []
    // Ensure the picker has the directory to choose from.
    usersStore.load()
  },
)

const addEnabled = computed(() => {
  if (!draft.roleIds.length) return false
  return draft.mode === 'existing' ? !!draft.userId : !!draft.email.trim()
})

const saving = ref(false)

async function onCreate() {
  if (!addEnabled.value || saving.value) return
  saving.value = true
  try {
    let userId = draft.userId
    if (draft.mode === 'new') {
      userId = await store.addProjectUser(props.project.id, { email: draft.email, name: draft.name })
      // Refresh the directory so the new user resolves to a name/email.
      await usersStore.reload()
    }
    await store.assignProjectUserRoles(props.project.id, userId, draft.roleIds)
    emit('created', userId)
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.userCreate.toastCreateFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
