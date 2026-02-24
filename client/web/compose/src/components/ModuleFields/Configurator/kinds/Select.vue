<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center gap-2">
      <Checkbox
        v-model="field.options.multiple"
        inputId="multipleSelect"
        :binary="true"
        :disabled="field.isMulti"
      />
      <label for="multipleSelect" class="cursor-pointer">
        {{ $t('field.kind.select.optionType.multiple') }}
      </label>
    </div>

    <div class="flex flex-col gap-4">
      <h4 class="font-semibold text-md m-0">{{ $t('field.kind.select.optionsLabel') }}</h4>

      <div v-for="(opt, index) in optionsList" :key="index" class="flex gap-2 items-center">
        <InputText
          v-model="opt.value"
          :placeholder="$t('field.kind.select.options.value')"
          class="flex-1"
          size="small"
          @change="updateOptions"
        />
        <InputText
          v-model="opt.text"
          :placeholder="$t('field.kind.select.options.label')"
          class="flex-1"
          size="small"
          @change="updateOptions"
        />
        <Button icon="pi pi-trash" severity="danger" text rounded @click="removeOption(index)" />
      </div>

      <Button
        :label="$t('general.label.add')"
        icon="pi pi-plus"
        severity="secondary"
        size="small"
        outlined
        class="self-start mt-2"
        @click="addOption"
      />
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
})

const optionsList = ref([])

onMounted(() => {
  if (Array.isArray(props.field.options?.options)) {
    optionsList.value = JSON.parse(JSON.stringify(props.field.options.options))
  }
})

watch(
  () => props.field.options?.options,
  newOpts => {
    if (Array.isArray(newOpts)) {
      // Only update if length differs or something changed to avoid losing focus
      if (newOpts.length !== optionsList.value.length) {
        optionsList.value = JSON.parse(JSON.stringify(newOpts))
      }
    }
  },
  { deep: true },
)

function updateOptions() {
  if (!props.field.options) props.field.options = {}
  props.field.options.options = optionsList.value.map(o => ({
    value: o.value,
    text: o.text || o.value,
  }))
}

function addOption() {
  optionsList.value.push({ value: '', text: '' })
  updateOptions()
}

function removeOption(index) {
  optionsList.value.splice(index, 1)
  updateOptions()
}
</script>
