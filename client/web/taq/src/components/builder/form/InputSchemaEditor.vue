<template>
  <div class="flex flex-col gap-3">
    <div class="flex items-center justify-between">
      <div class="text-sm font-medium text-color">
        {{ $t('builder.inputSchema.title') }}
      </div>
      <Button
        icon="pi pi-plus"
        size="small"
        severity="secondary"
        text
        :label="$t('builder.inputSchema.addParam')"
        @click="addRow"
      />
    </div>

    <p class="text-xs text-muted-color">
      {{ $t('builder.inputSchema.hint') }}
    </p>

    <div
      v-if="!rows.length"
      class="text-xs text-muted-color italic border border-dashed border-surface-300 dark:border-surface-700 rounded px-3 py-4 text-center"
    >
      {{ $t('builder.inputSchema.empty') }}
    </div>

    <div
      v-for="(row, index) in rows"
      :key="index"
      class="flex flex-col gap-2 border border-surface-200 dark:border-surface-700 rounded p-2"
    >
      <div class="flex items-start gap-2">
        <div class="flex-1 min-w-0">
          <label class="text-xs text-muted-color block mb-1">
            {{ $t('builder.inputSchema.name') }}
          </label>
          <InputText
            :model-value="row.name"
            size="small"
            class="w-full"
            :invalid="!!nameErrors[index]"
            :placeholder="$t('builder.inputSchema.namePlaceholder')"
            @update:model-value="updateRow(index, { name: $event })"
          />
          <div v-if="nameErrors[index]" class="text-xs text-red-500 mt-1">
            {{ nameErrors[index] }}
          </div>
        </div>

        <div class="w-32 shrink-0">
          <label class="text-xs text-muted-color block mb-1">
            {{ $t('builder.inputSchema.type') }}
          </label>
          <Select
            :model-value="row.type"
            :options="typeOptions"
            option-label="label"
            option-value="value"
            size="small"
            class="w-full"
            @update:model-value="updateRow(index, { type: $event })"
          />
        </div>

        <div class="shrink-0 flex flex-col items-center">
          <label class="text-xs text-muted-color block mb-1">
            {{ $t('builder.inputSchema.required') }}
          </label>
          <Checkbox
            :model-value="!!row.required"
            binary
            class="mt-1"
            @update:model-value="updateRow(index, { required: $event })"
          />
        </div>

        <Button
          icon="pi pi-trash"
          size="small"
          severity="danger"
          text
          class="mt-5 shrink-0"
          :aria-label="$t('builder.inputSchema.remove')"
          @click="removeRow(index)"
        />
      </div>

      <div>
        <label class="text-xs text-muted-color block mb-1">
          {{ $t('builder.inputSchema.description') }}
        </label>
        <InputText
          :model-value="row.description"
          size="small"
          class="w-full"
          :placeholder="
            $t('builder.inputSchema.descriptionPlaceholder')
          "
          @update:model-value="updateRow(index, { description: $event })"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue'])

const rows = computed(() => props.modelValue || [])

const typeOptions = computed(() => [
  { label: t('builder.inputSchema.types.string'), value: 'String' },
  { label: t('builder.inputSchema.types.number'), value: 'Number' },
  { label: t('builder.inputSchema.types.boolean'), value: 'Boolean' },
])

const nameErrors = computed(() => {
  const errors = {}
  const seen = new Map()
  rows.value.forEach((row, i) => {
    const name = (row.name || '').trim()
    if (!name) {
      errors[i] = t('builder.inputSchema.errors.nameRequired')
      return
    }
    if (!/^[A-Za-z_][A-Za-z0-9_]*$/.test(name)) {
      errors[i] = t('builder.inputSchema.errors.nameInvalid')
      return
    }
    if (seen.has(name)) {
      errors[i] = t('builder.inputSchema.errors.nameDuplicate')
      const first = seen.get(name)
      if (!errors[first]) {
        errors[first] = t('builder.inputSchema.errors.nameDuplicate')
      }
    } else {
      seen.set(name, i)
    }
  })
  return errors
})

function emitRows(next) {
  emit('update:modelValue', next)
}

function addRow() {
  emitRows([...rows.value, { name: '', type: 'String', required: false, description: '' }])
}

function removeRow(index) {
  const next = [...rows.value]
  next.splice(index, 1)
  emitRows(next)
}

function updateRow(index, patch) {
  const next = rows.value.map((r, i) => (i === index ? { ...r, ...patch } : r))
  emitRows(next)
}
</script>
