<template>
  <div class="flex flex-col gap-6">
    <!-- Select type -->
    <div>
      <label class="font-medium text-muted-color text-sm block mb-3">
        {{ $t('field.kind.select.optionType.label') }}
      </label>
      <div class="flex flex-col gap-2">
        <div
          v-for="opt in selectTypeOptions"
          :key="opt.value"
          class="flex items-center gap-2"
        >
          <RadioButton
            :input-id="`selectType-${opt.value}`"
            v-model="field.options.selectType"
            :value="opt.value"
          />
          <label :for="`selectType-${opt.value}`" class="cursor-pointer">{{ opt.label }}</label>
        </div>
      </div>
    </div>

    <!-- Allow duplicates (only for types that allow it) -->
    <div v-if="showAllowDuplicates" class="flex items-center gap-2">
      <Checkbox
        :model-value="!field.options.isUniqueMultiValue"
        inputId="allowDuplicates"
        :binary="true"
        @update:model-value="field.options.isUniqueMultiValue = !$event"
      />
      <label for="allowDuplicates" class="cursor-pointer">{{ $t('field.kind.select.allow-duplicates') }}</label>
    </div>

    <!-- Display type -->
    <div>
      <label class="font-medium text-muted-color text-sm block mb-3">
        {{ $t('field.kind.select.displayType.label') }}
      </label>
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <RadioButton inputId="displayText" v-model="field.options.displayType" value="text" />
          <label for="displayText" class="cursor-pointer">{{ $t('field.kind.select.displayType.text') }}</label>
        </div>
        <div class="flex items-center gap-2">
          <RadioButton inputId="displayBadge" v-model="field.options.displayType" value="badge" />
          <label for="displayBadge" class="cursor-pointer">{{ $t('field.kind.select.displayType.badge') }}</label>
        </div>
      </div>
    </div>

    <!-- Options list -->
    <div class="flex flex-col gap-4">
      <h4 class="font-semibold text-md m-0">{{ $t('field.kind.select.optionsLabel') }}</h4>

      <div
        v-for="(opt, index) in optionsList"
        :key="index"
        class="flex flex-col gap-2 pb-4 border-b border-surface last:border-0 last:pb-0"
      >
        <div class="flex gap-2 items-center">
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

        <!-- Badge color styling -->
        <div v-if="field.options.displayType === 'badge'" class="flex gap-4 pl-1">
          <div class="flex items-center gap-2">
            <label class="text-xs text-muted-color">{{ $t('field.kind.select.options.style.textColor') }}</label>
            <CInputColor
              v-model="opt.style.textColor"
              @update:model-value="updateOptions"
            />
          </div>
          <div class="flex items-center gap-2">
            <label class="text-xs text-muted-color">{{ $t('field.kind.select.options.style.backgroundColor') }}</label>
            <CInputColor
              v-model="opt.style.backgroundColor"
              @update:model-value="updateOptions"
            />
          </div>
        </div>
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
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import CInputColor from '../CInputColor.vue'

const { t } = useI18n()

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
})

const selectTypeOptions = computed(() => [
  { value: 'default', label: t('field.kind.select.optionType.default') },
  { value: 'multiple', label: t('field.kind.select.optionType.multiple') },
  { value: 'each', label: t('field.kind.select.optionType.each') },
])

// Types where duplicates can be allowed
const duplicatesAllowedTypes = ['default', 'each']
const showAllowDuplicates = computed(() =>
  duplicatesAllowedTypes.includes(props.field.options.selectType),
)

const optionsList = ref([])

function makeOption(o = {}) {
  return {
    value: o.value || '',
    text: o.text || '',
    style: {
      textColor: o.style?.textColor || '',
      backgroundColor: o.style?.backgroundColor || '',
    },
  }
}

onMounted(() => {
  if (!props.field.options.selectType) {
    props.field.options.selectType = 'default'
  }
  if (!props.field.options.displayType) {
    props.field.options.displayType = 'text'
  }
  if (Array.isArray(props.field.options?.options)) {
    optionsList.value = props.field.options.options.map(makeOption)
  }
})

watch(
  () => props.field.options?.options,
  newOpts => {
    if (Array.isArray(newOpts) && newOpts.length !== optionsList.value.length) {
      optionsList.value = newOpts.map(makeOption)
    }
  },
  { deep: true },
)

function updateOptions() {
  if (!props.field.options) props.field.options = {}
  props.field.options.options = optionsList.value.map(o => ({
    value: o.value,
    text: o.text || o.value,
    style: {
      textColor: o.style?.textColor || '',
      backgroundColor: o.style?.backgroundColor || '',
    },
  }))
}

function addOption() {
  optionsList.value.push(makeOption())
  updateOptions()
}

function removeOption(index) {
  optionsList.value.splice(index, 1)
  updateOptions()
}
</script>
