<template>
  <b-form-group
    :label-cols-md="horizontal && '5'"
    :label-cols-xl="horizontal && '4'"
    :content-cols-md="horizontal && '7'"
    :content-cols-xl="horizontal && '8'"
    :class="formGroupStyleClasses"
  >
    <template
      #label
    >
      <div
        v-if="!valueOnly"
        class="d-flex align-items-center text-primary p-0"
      >
        <span
          :title="label"
          class="d-inline-block mw-100 py-1"
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
