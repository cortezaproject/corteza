<template>
  <b-form-group
    label-class="text-primary"
    :class="formGroupStyleClasses"
  >
    <template
      v-if="!valueOnly"
      #label
    >
      <div
        class="d-flex align-items-top"
      >
        <label
          class="mb-0"
        >
          {{ label }}
        </label>

        <hint
          :id="field.fieldID"
          :text="hint"
        />
      </div>
      <small
        class="form-text font-weight-light text-muted"
      >
        {{ description }}
      </small>
    </template>

    <multi
      v-if="field.isMulti"
      v-slot="ctx"
      :value.sync="value"
      :errors="errors"
    >
      <c-input-date-time
        v-model="value[ctx.index]"
        :no-date="field.options.onlyTime"
        :no-time="field.options.onlyDate"
        :only-future="field.options.onlyFutureValues"
        :only-past="field.options.onlyPastValues"
        :labels="{
          clear: $t('general:label.clear'),
          none: $t('general:label.none'),
          now: $t('general:label.now'),
          today: $t('general:label.today'),
        }"
      />
    </multi>

    <template
      v-else
    >
      <c-input-date-time
        v-model="value"
        :no-date="field.options.onlyTime"
        :no-time="field.options.onlyDate"
        :only-future="field.options.onlyFutureValues"
        :only-past="field.options.onlyPastValues"
        :labels="{
          clear: $t('general:label.clear'),
          none: $t('general:label.none'),
          now: $t('general:label.now'),
          today: $t('general:label.today'),
        }"
      />
      <errors :errors="errors" />
    </template>
  </b-form-group>
</template>
<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
const { CInputDateTime } = components

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    CInputDateTime,
  },

  extends: base,
}
</script>
