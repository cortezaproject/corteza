<template>
  <div class="flex flex-col gap-4">
    <p v-if="message" class="m-0 whitespace-pre-wrap" v-html="message" />

    <label v-if="label" class="text-sm font-medium text-color">{{ label }}</label>

    <CInputDateTime
      v-if="fieldType === 'date' || fieldType === 'time' || fieldType === 'datetime'"
      v-model="value"
      :disabled="loading"
      :time-only="fieldType === 'time'"
      :only-date="fieldType === 'date'"
      :show-time="fieldType === 'datetime'"
    />

    <InputText
      v-else
      v-model="value"
      :type="fieldType"
      :disabled="loading"
    />

    <div class="flex justify-end">
      <Button
        :disabled="loading"
        :label="pVal('buttonLabel', 'Submit')"
        @click="$emit('submit', { value: { '@value': value, '@type': 'String' } })"
      />
    </div>
  </div>
</template>

<script>
import { CInputDateTime } from '../../input'
import base from './base.vue'

const validTypes = ['text', 'number', 'email', 'password', 'search', 'date', 'time', 'datetime']

export default {
  name: 'CPromptInput',
  components: {
    CInputDateTime,
  },
  extends: base,
  emits: ['submit'],
  data() {
    return {
      value: undefined,
    }
  },
  computed: {
    fieldType() {
      const type = this.pVal('type', 'text')
      return validTypes.includes(type) ? type : 'text'
    },
  },
  beforeMount() {
    this.value = this.pVal('inputValue', '')
  },
}
</script>
