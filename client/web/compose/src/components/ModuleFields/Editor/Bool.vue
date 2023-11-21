<template>
  <b-form-group
    :label-cols-md="horizontal && '5'"
    :label-cols-xl="horizontal && '4'"
    :content-cols-md="horizontal && '7'"
    :content-cols-xl="horizontal && '8'"
  >
    <template
      v-if="field.options.switch"
      #label
    >
      <div
        v-if="!valueOnly"
        class="d-flex align-items-center text-primary p-0"
      >
        <span
          :title="label"
          class="d-inline-block mw-100 pt-0"
          :class="{ 'py-1': !horizontal }"
        >
          {{ label }}
        </span>

        <c-hint :tooltip="hint" />

        <slot name="tools" />
      </div>
      <div
        class="small text-muted"
        :class="{ 'mb-1': description }"
      >
        {{ description }}
      </div>
    </template>

    <c-input-checkbox
      v-model="value"
      :switch="field.options.switch"
      :labels="field.options.switch ? checkboxLabel : {}"
    >
      <div
        v-if="!field.options.switch"
        class="d-flex align-items-center text-primary"
      >
        {{ label }}

        <c-hint :tooltip="hint" />
      </div>
    </c-input-checkbox>

    <div
      v-if="!field.options.switch"
      class="small text-muted"
    >
      {{ description }}
    </div>
    <errors :errors="errors" />
  </b-form-group>
</template>
<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  computed: {
    value: {
      get () {
        return this.record.values[this.field.name] === '1'
      },

      set (value) {
        this.record.values[this.field.name] = value ? '1' : '0'
        this.$emit('change', value)
      },
    },

    checkboxLabel () {
      return {
        on: this.field.options.trueLabel || this.$t('general:label.yes'),
        off: this.field.options.falseLabel || this.$t('general:label.no'),
      }
    },
  },
}
</script>
