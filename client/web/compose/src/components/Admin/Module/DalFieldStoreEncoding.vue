<template>
  <div class="grid grid-cols-12 gap-2 items-center mb-2">
    <!-- Checkbox or Label -->
    <div class="col-span-3">
      <div v-if="allowOmitStrategy" class="flex items-center gap-2">
        <Checkbox
          v-model="use"
          :binary="true"
          :disabled="disabled"
          :inputId="`chk-use-${field}`"
        />
        <label :for="`chk-use-${field}`" class="cursor-pointer select-none">
          {{ label }}
        </label>
      </div>
      <div v-else class="font-bold">
        {{ label }}
      </div>
    </div>

    <!-- Strategy Select -->
    <div class="col-span-3">
      <Select
        v-show="strategy !== 'omit'"
        v-model="strategy"
        :options="strategies"
        optionLabel="text"
        optionValue="value"
        :disabled="!use"
        size="small"
        class="w-full"
      />
    </div>

    <!-- Ident Input -->
    <div class="col-span-6">
      <InputText
        v-if="strategy === ''"
        :value="storeIdent"
        :placeholder="$t('module.edit.config.dal.ident.placeholder')"
        size="small"
        readonly
        class="w-full"
      />
      <InputText
        v-else-if="showIdentInput"
        v-model="draft.ident"
        :placeholder="$t('module.edit.config.dal.ident.placeholder')"
        :disabled="disableIdentInput"
        size="small"
        class="w-full"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { defaultConfigDraft, types } from './encoding-strategy'

const props = defineProps({
  config: {
    type: Object,
    required: true,
  },
  field: {
    type: String,
    required: true,
  },
  label: {
    type: String,
    required: true,
  },
  isMulti: {
    type: Boolean,
    default: false,
  },
  storeIdent: {
    type: String,
    required: true,
  },
  defaultStrategy: {
    type: String,
    default: types.Plain,
  },
  allowOmitStrategy: {
    type: Boolean,
    default: true,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['change'])
const { t } = useI18n()

// Holds working copy of strategy config
const draft = ref(defaultConfigDraft(props.config, props.storeIdent))
// Strategy before omit
const undoOmit = ref(props.defaultStrategy)

const strategies = computed(() => {
  return [
    { value: types.Plain, text: t('module.edit.config.dal.encoding-strategy.strategies.plain.label'), disabled: props.isMulti },
    { value: types.Alias, text: t('module.edit.config.dal.encoding-strategy.strategies.alias.label'), disabled: props.isMulti },
    { value: types.JSON, text: t('module.edit.config.dal.encoding-strategy.strategies.json.label') },
  ].filter(({ disabled }) => !disabled)
})

const showIdentInput = computed(() => {
  return [types.JSON, types.Alias, types.Plain].includes(strategy.value)
})

const disableIdentInput = computed(() => {
  return [types.Plain].includes(strategy.value)
})

const strategy = computed({
  get() {
    // iterate over all types and return the first one that matches
    for (const type of Object.values(types)) {
      if (props.config[type] === undefined) {
        continue
      }
      return type
    }
    return props.defaultStrategy
  },
  set(newStrategy) {
    emit('change', { strategy: newStrategy, config: draft.value })
  },
})

const use = computed({
  get() {
    return strategy.value !== types.Omit
  },
  set(newUse) {
    if (strategy.value !== types.Omit) {
      undoOmit.value = strategy.value
    }
    strategy.value = newUse ? undoOmit.value : types.Omit
  },
})

watch(
  draft,
  (newDraft) => {
    emit('change', { strategy: strategy.value, config: newDraft })
  },
  { deep: true },
)
</script>
