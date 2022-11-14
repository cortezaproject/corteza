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

    <template v-if="field.isMulti">
      <multi
        :value.sync="value"
        :errors="errors"
        :single-input="field.options.selectType !== 'each'"
      >
        <template v-slot:single>
          <b-form-select
            v-if="field.options.selectType === 'default'"
            ref="singleSelect"
            :options="selectOptions"
            @change="selectChange"
          >
            <template slot="first">
              <option
                :value="undefined"
                disabled
              >
                {{ $t('kind.select.placeholder') }}
              </option>
            </template>
          </b-form-select>
          <b-form-select
            v-if="field.options.selectType === 'multiple'"
            v-model="value"
            :options="selectOptions"
            :select-size="6"
            multiple
          />
        </template>
        <template v-slot:default="ctx">
          <b-form-select
            v-if="field.options.selectType === 'each'"
            v-model="value[ctx.index]"
            :options="selectOptions"
          >
            <template slot="first">
              <option
                :value="undefined"
                disabled
              >
                {{ $t('kind.select.placeholder') }}
              </option>
            </template>
          </b-form-select>
          <span v-else>{{ findLabel(value[ctx.index]) }}</span>
        </template>
      </multi>
    </template>

    <template
      v-else
    >
      <b-form-select
        v-model="value"
        :options="selectOptions"
      >
        <template slot="first">
          <option :value="undefined">
            {{ $t('kind.select.optionNotSelected') }}
          </option>
        </template>
      </b-form-select>
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
        const disabled = o.value && this.field.isMulti ? (this.value || []).includes(o.value) : this.value === o.value
        return { ...o, disabled }
      }).filter(({ value = '' }) => value)
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
