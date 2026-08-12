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

    <InputText v-else v-model="value" :type="fieldType" :disabled="loading" @keyup.enter="submit" />

    <div class="flex justify-end">
      <Button
        :disabled="loading"
        :label="pVal('buttonLabel', tF('general.label.submit', 'Submit'))"
        @click="submit"
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
  methods: {
    submit() {
      if (this.loading) return
      this.$emit('submit', { value: { '@value': this.value, '@type': 'String' } })
    },
  },
}
</script>
