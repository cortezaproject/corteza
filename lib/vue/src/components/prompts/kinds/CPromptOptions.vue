<template>
  <div class="flex flex-col gap-4">
    <p v-if="message" class="m-0 whitespace-pre-wrap" v-html="message" />

    <label v-if="label" class="text-sm font-medium text-color">{{ label }}</label>

    <MultiSelect
      v-if="inputType === 'select' && multiple"
      v-model="value"
      :options="itemOptions"
      option-label="text"
      option-value="value"
      :placeholder="placeholder"
      class="w-full"
      :disabled="loading"
      display="chip"
      append-to="self"
    />

    <Select
      v-else-if="inputType === 'select'"
      v-model="value"
      :options="itemOptions"
      option-label="text"
      option-value="value"
      :placeholder="placeholder"
      class="w-full"
      :disabled="loading"
      append-to="self"
    />

    <div v-else class="flex flex-col gap-2">
      <div v-for="option in itemOptions" :key="option.value" class="flex items-center gap-2">
        <RadioButton v-model="value" :input-id="`prompt-option-${option.value}`" :value="option.value" />
        <label :for="`prompt-option-${option.value}`">{{ option.text }}</label>
      </div>
    </div>

    <div class="flex justify-end">
      <Button
        :disabled="loading"
        :label="pVal('buttonLabel', tF('general.label.submit', 'Submit'))"
        @click="$emit('submit', { value: encodeValue() })"
      />
    </div>
  </div>
</template>

<script>
import base from './base.vue'

const validTypes = ['select', 'radio']

export default {
  name: 'CPromptOptions',
  extends: base,
  emits: ['submit'],
  data() {
    return {
      value: undefined,
    }
  },
  computed: {
    itemOptions() {
      const options = this.pVal('options', {})
      return Object.keys(options || {}).map(value => ({
        value,
        text: options[value],
      }))
    },
    inputType() {
      const type = this.pVal('type', 'select')
      return validTypes.includes(type) ? type : 'select'
    },
    multiple() {
      return !!this.pVal('multiselect', false)
    },
    placeholder() {
      return this.pVal('placeholder', this.tF('prompt.options.placeholder', 'Select an option'))
    },
  },
  beforeMount() {
    let value = this.pVal('value')

    if (this.multiple && this.inputType === 'select') {
      if (Array.isArray(value)) {
        value = value.map(v => (v && typeof v === 'object' && '@value' in v) ? v['@value'] : v)
      } else {
        value = value ? [value] : []
      }
    }

    this.value = value
  },
  methods: {
    encodeValue() {
      if (Array.isArray(this.value)) {
        return { '@type': 'Array', '@value': this.value || [] }
      }
      return { '@type': 'String', '@value': this.value }
    },
  },
}
</script>
