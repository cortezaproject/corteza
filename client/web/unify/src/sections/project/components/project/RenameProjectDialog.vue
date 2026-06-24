<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    modal
    :header="$t('project.renameDialog.header')"
    :style="{ width: '28rem' }"
  >
    <CFormGroup :label="$t('general.label.name')" required>
      <InputText v-model="draft" fluid autofocus @keyup.enter="save" />
    </CFormGroup>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button :label="$t('general.label.cancel')" severity="secondary" text size="small" @click="$emit('update:visible', false)" />
        <Button :label="$t('general.label.save')" size="small" :disabled="!draft.trim()" @click="save" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  visible: { type: Boolean, default: false },
  project: { type: Object, default: null },
})
const emit = defineEmits(['update:visible', 'rename'])

const draft = ref('')

watch(
  () => props.visible,
  v => {
    if (v) draft.value = props.project?.name || ''
  },
)

function save() {
  const next = draft.value.trim()
  if (!next) return
  emit('rename', next)
  emit('update:visible', false)
}
</script>
