<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    modal
    :header="headerText"
    :style="{ width: '46rem' }"
    :pt="{ content: { class: '!pt-2' } }"
  >
    <!-- Reuse the per-category governance form; two-column layout with
         textareas spanning both columns. `model` is a plain object keyed by
         field.key and is reset every time the dialog opens. -->
    <GovernanceForm :schema="schema" v-model="model" :columns="2" />

    <template #footer>
      <div class="flex items-center justify-end gap-2 w-full">
        <Button :label="$t('general.label.cancel')" severity="secondary" text size="small" @click="close" />
        <Button :label="createLabel" size="small" @click="onCreate" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import GovernanceForm from '@/sections/project/components/wizard/GovernanceForm.vue'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  category: { type: String, required: true },
  schema: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:visible', 'create'])

const { t } = useI18n()

// Collected form values, keyed by field.key. Reset on each open.
const model = ref({})

// "New {type}" where {type} is the singular category label.
const createLabel = computed(() =>
  t('project.dashboard.newButton', { type: t(`project.dashboard.categorySingular.${props.category}`) }),
)
const headerText = createLabel

// Opening the dialog starts from a blank form so no stale values leak between
// categories or repeated creates.
watch(
  () => props.visible,
  v => {
    if (v) model.value = {}
  },
)

function close() {
  emit('update:visible', false)
}

function onCreate() {
  emit('create', { ...model.value })
  close()
}
</script>
