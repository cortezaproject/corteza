<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    modal
    header="Rename project"
    :style="{ width: '28rem' }"
  >
    <CFormGroup label="Name" required>
      <InputText v-model="draft" fluid autofocus @keyup.enter="save" />
    </CFormGroup>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button label="Cancel" severity="secondary" outlined size="small" @click="$emit('update:visible', false)" />
        <Button label="Save" size="small" :disabled="!draft.trim()" @click="save" />
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
