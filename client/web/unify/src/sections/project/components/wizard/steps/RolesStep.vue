<template>
  <div class="h-full overflow-auto p-4">
    <CFormGroup :label="$t('project.accessRoles.title')">
      <template #actions>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.accessRoles.add')"
          severity="secondary"
          size="small"
          @click="openCreate"
        />
      </template>

      <div class="mt-1">
        <CFormItemList
          :items="roles"
          item-key="id"
          reveal-on-hover
          :empty-message="$t('project.accessRoles.empty')"
          :hide-remove="disabled"
          :remove-label="$t('project.accessRoles.remove')"
          @remove="onRemove"
        >
          <template #default="{ item }">
            <div class="flex items-center gap-3 min-w-0">
              <span
                class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
                :class="[cfg.bg, cfg.ring]"
              >
                <i :class="[cfg.icon, cfg.text]" />
              </span>
              <div class="min-w-0">
                <div class="font-medium truncate">{{ item.name }}</div>
                <div v-if="item.description" class="text-xs text-muted-color truncate">
                  {{ item.description }}
                </div>
              </div>
            </div>
          </template>
        </CFormItemList>
      </div>
    </CFormGroup>

    <Dialog
      v-model:visible="dialogOpen"
      modal
      :header="$t('project.accessRoles.addDialog.header')"
      :style="{ width: '30rem' }"
      :pt="{ footer: { class: 'flex justify-end gap-2' } }"
    >
      <div class="flex flex-col gap-4">
        <CFormGroup :label="$t('project.accessRoles.addDialog.name')" required>
          <InputText v-model="draft.name" fluid :placeholder="$t('project.accessRoles.addDialog.namePlaceholder')" />
        </CFormGroup>
        <CFormGroup :label="$t('project.accessRoles.addDialog.description')">
          <Textarea
            v-model="draft.description"
            rows="3"
            auto-resize
            fluid
            :placeholder="$t('project.accessRoles.addDialog.descriptionPlaceholder')"
          />
        </CFormGroup>
      </div>
      <template #footer>
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="dialogOpen = false"
        />
        <Button
          :label="$t('general.label.add')"
          size="small"
          :disabled="!draft.name.trim() || saving"
          @click="add"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

const cfg = kindConfig('role')
const roles = computed(() => store.rolesFor(props.project.id))

const dialogOpen = ref(false)
const saving = ref(false)
const draft = reactive({ name: '', description: '' })

function openCreate() {
  draft.name = ''
  draft.description = ''
  dialogOpen.value = true
}

async function add() {
  if (!draft.name.trim()) return
  saving.value = true
  try {
    await store.addRole(props.project.id, { name: draft.name, description: draft.description })
    dialogOpen.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.accessRoles.toastAddFailed'))(err)
  } finally {
    saving.value = false
  }
}

function onRemove(r) {
  confirmDelete({
    header: t('project.accessRoles.removeConfirm.header'),
    message: t('project.accessRoles.removeConfirm.message', { name: r.name }),
    onConfirm: async () => {
      try {
        await store.removeRole(props.project.id, r.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.accessRoles.toastRemoveFailed'))(err)
      }
    },
  })
}

async function refresh(id) {
  try {
    await store.loadRoles(id)
  } catch (err) {
    $toast.toastErrorHandler(t('project.accessRoles.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.id))
watch(
  () => props.project.id,
  id => id && refresh(id),
)
</script>
