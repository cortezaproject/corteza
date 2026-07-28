<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    modal
    :header="$t('project.newDialog.headerDetails')"
    :style="{ width: '32rem' }"
    :pt="{ content: { class: '!pt-2' } }"
  >
    <div class="flex flex-col gap-4 pt-1">
      <CFormGroup :label="$t('general.label.name')" required>
        <InputText
          v-model="name"
          :placeholder="$t('project.newDialog.namePlaceholder')"
          fluid
          autofocus
        />
      </CFormGroup>

      <CFormGroup :label="$t('general.label.description')">
        <Textarea
          v-model="description"
          rows="2"
          auto-resize
          fluid
          :placeholder="$t('project.newDialog.descriptionPlaceholder')"
        />
      </CFormGroup>
    </div>

    <template #footer>
      <div class="flex items-center justify-end gap-2 w-full">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="close"
        />
        <Button
          :label="$t('project.newDialog.createProject')"
          size="small"
          :disabled="!canCreate"
          :loading="creating"
          @click="onCreate"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
})
const emit = defineEmits(['update:visible', 'created'])

const { t } = useI18n()
const store = useProjectsStore()
const $toast = inject('$toast')

const name = ref('')
const description = ref('')

const canCreate = computed(() => !!name.value.trim())

function reset() {
  name.value = ''
  description.value = ''
}

watch(
  () => props.visible,
  v => {
    if (v) reset()
  },
)

function close() {
  emit('update:visible', false)
}

const creating = ref(false)

async function onCreate() {
  if (!canCreate.value || creating.value) return
  creating.value = true
  try {
    const project = await store.create({
      name: name.value,
      description: description.value,
    })
    emit('created', project)
    close()
  } catch (err) {
    // Dialog stays open so nothing typed is lost.
    $toast.toastErrorHandler(t('project.newDialog.toastCreateFailed'))(err)
  } finally {
    creating.value = false
  }
}
</script>
