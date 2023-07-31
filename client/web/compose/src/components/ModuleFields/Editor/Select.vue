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

    <template v-if="field.isMulti">
      <template v-if="field.options.selectType === 'list'">
        <b-form-checkbox-group
          v-model="value"
          :options="selectOptions"
          stacked
        />
      </template>

      <multi
        v-else
        :value.sync="value"
        :errors="errors"
        :single-input="field.options.selectType !== 'each'"
      >
        <template v-slot:single>
          <c-input-select
            v-if="field.options.selectType === 'default'"
            ref="singleSelect"
            :options="selectOptions"
            :placeholder="$t('kind.select.placeholder')"
            :reduce="o => o.value"
            label="text"
            @input="selectChange"
          />

          <c-input-select
            v-if="field.options.selectType === 'multiple'"
            v-model="value"
            :options="selectOptions"
            label="text"
            :reduce="o => o.value"
            multiple
          />
        </template>

        <template v-slot:default="ctx">
          <c-input-select
            v-if="field.options.selectType === 'each'"
            v-model="value[ctx.index]"
            :options="selectOptions"
            :reduce="o => o.value"
            :placeholder="$t('kind.select.placeholder')"
            label="text"
          />

          <span v-else>{{ findLabel(value[ctx.index]) }}</span>
        </template>
      </multi>
    </template>

    <template
      v-else
    >
      <c-input-select
        v-if="field.options.selectType === 'default'"
        v-model="value"
        :options="selectOptions"
        :reduce="o => o.value"
        label="text"
        :placeholder="$t('kind.select.optionNotSelected')"
      />

      <b-form-radio-group
        v-else
        v-model="value"
        :options="selectOptions"
        stacked
      />

      <errors :errors="errors" />
    </template>
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
    selectOptions () {
      return this.field.options.options.map(o => {
        const disabled = o.value && this.field.isMulti && !this.field.options.isUniqueMultiValue
          ? this.value === o.value
          : (this.value || []).includes(o.value)
        return { ...o, disabled: this.field.options.selectType !== 'list' && disabled }
      }).filter(({ value = '', text = '' }) => value && text)
    },
  },

  methods: {
    selectChange (value) {
      this.value.push(value)
      // Reset select
      this.$refs.singleSelect.localValue = undefined
    },

    /**
     * Helper to resolve a label for a given value
     * @param {String} v Value in question
     * @returns {String}
     */
    findLabel (v) {
      return (this.selectOptions.find(({ value }) => value === v) || {}).text || v
    },
  },
}
</script>
