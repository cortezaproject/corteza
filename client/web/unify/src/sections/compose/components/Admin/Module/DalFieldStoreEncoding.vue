<template>
  <div class="grid grid-cols-12 gap-2 items-center mb-2">
    <div class="col-span-3">
      <div v-if="allowOmitStrategy" class="flex items-center gap-2">
        <Checkbox
          :modelValue="use"
          :binary="true"
          :disabled="disabled"
          :inputId="checkboxId"
          @update:modelValue="onUseToggle"
        />
        <label
          :for="checkboxId"
          :class="['select-none', disabled ? 'cursor-not-allowed text-muted-color' : 'cursor-pointer']"
        >
          {{ label }}
        </label>
      </div>
      <div v-else class="font-bold">
        {{ label }}
      </div>
    </div>

    <div class="col-span-3">
      <Select
        v-show="strategy !== types.Omit"
        :modelValue="strategy"
        :options="strategies"
        optionLabel="text"
        optionValue="value"
        :disabled="disabled || !use"
        size="small"
        class="w-full"
        @update:modelValue="onStrategyChange"
      />
    </div>

    <div class="col-span-6">
      <InputText
        v-if="strategy === types.Plain"
        :value="storeIdent"
        :placeholder="$t('module.edit.config.dal.encoding-strategy.ident.placeholder')"
        size="small"
        readonly
        class="w-full"
      />
      <InputText
        v-else-if="showIdentInput"
        v-model="draft.ident"
        :placeholder="$t('module.edit.config.dal.encoding-strategy.ident.placeholder')"
        :disabled="disabled"
        size="small"
        class="w-full"
        @update:modelValue="onIdentInput"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { defaultConfigDraft, types } from './encoding-strategy'

const props = defineProps({
  config: { type: Object, required: true },
  field: { type: String, required: true },
  label: { type: String, required: true },
  isMulti: { type: Boolean, default: false },
  storeIdent: { type: String, required: true },
  defaultStrategy: { type: String, default: types.Plain },
  allowOmitStrategy: { type: Boolean, default: true },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['change'])
const { t } = useI18n()

const draft = ref(defaultConfigDraft(props.config, props.storeIdent))
const undoOmit = ref(props.defaultStrategy)

const checkboxId = computed(() => `dal-fse-${props.field}`)

const strategies = computed(() => [
  { value: types.Plain, text: t('module.edit.config.dal.encoding-strategy.strategies.plain.label'), disabled: props.isMulti },
  { value: types.Alias, text: t('module.edit.config.dal.encoding-strategy.strategies.alias.label'), disabled: props.isMulti },
  { value: types.JSON, text: t('module.edit.config.dal.encoding-strategy.strategies.json.label') },
].filter(({ disabled }) => !disabled))

const strategy = computed(() => {
  for (const type of Object.values(types)) {
    if (props.config[type] !== undefined) return type
  }
  return props.defaultStrategy
})

const use = computed(() => strategy.value !== types.Omit)

const showIdentInput = computed(() => [types.JSON, types.Alias, types.Plain].includes(strategy.value))

function onUseToggle(newUse) {
  if (strategy.value !== types.Omit) {
    undoOmit.value = strategy.value
  }
  const next = newUse ? undoOmit.value : types.Omit
  emit('change', { strategy: next, config: draft.value })
}

function onStrategyChange(newStrategy) {
  if (newStrategy == null || newStrategy === strategy.value) return
  emit('change', { strategy: newStrategy, config: draft.value })
}

function onIdentInput() {
  emit('change', { strategy: strategy.value, config: draft.value })
}

// Keep draft in sync if parent feeds a different config (e.g. preset reset)
watch(
  () => props.config,
  (cfg) => {
    const next = defaultConfigDraft(cfg, props.storeIdent)
    if (next.ident !== draft.value.ident) {
      draft.value = next
    }
  },
  { deep: true },
)
</script>
