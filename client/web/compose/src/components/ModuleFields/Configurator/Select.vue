<template>
  <b-row no-gutters>
    <b-col>
      <b-form-group
        :label="$t('kind.select.optionsLabel')"
      >
        <b-input-group
          v-for="(option, index) in f.options.options"
          :key="index"
          class="mb-1"
        >
          <b-form-input
            v-model.trim="f.options.options[index].value"
            plain
            size="sm"
            :placeholder="$t('kind.select.optionValuePlaceholder')"
          />

          <b-form-input
            v-model.trim="f.options.options[index].text"
            plain
            size="sm"
            :placeholder="$t('kind.select.optionLabelPlaceholder')"
          />

          <b-input-group-append>
            <field-select-translator
              v-if="field"
              :field="field"
              :module="module"
              :highlight-key="`meta.options.${option.value}.text`"
              size="sm"
              :disabled="isNew || option.new"
              button-variant="light"
            />
            <b-button
              variant="outline-danger"
              class="border-0"
              @click.prevent="f.options.options.splice(index, 1)"
            >
              <font-awesome-icon :icon="['far', 'trash-alt']" />
            </b-button>
          </b-input-group-append>
        </b-input-group>

        <b-input-group>
          <b-form-input
            v-model.trim="newOption.value"
            plain
            size="sm"
            :placeholder="$t('kind.select.optionValuePlaceholder')"
            :state="newOptState"
            @keypress.enter.prevent="handleAddOption"
          />

          <b-form-input
            v-model.trim="newOption.text"
            plain
            size="sm"
            :placeholder="$t('kind.select.optionLabelPlaceholder')"
            :state="newOptState"
            @keypress.enter.prevent="handleAddOption"
          />

          <b-input-group-append>
            <b-button
              variant="primary"
              size="sm"
              :disabled="newOptState === false || newEmpty"
              @click.prevent="handleAddOption"
            >
              + {{ $t('kind.select.optionAdd') }}
            </b-button>
          </b-input-group-append>
        </b-input-group>
      </b-form-group>

      <b-form-group
        v-if="f.isMulti"
      >
        <label class="d-block">{{ $t('kind.select.optionType.label') }}</label>
        <b-form-radio-group
          v-model="f.options.selectType"
          :options="selectOptions"
          stacked
        />
      </b-form-group>
    </b-col>
  </b-row>
</template>

<script>
import base from './base'
import { NoID } from '@cortezaproject/corteza-js'
import FieldSelectTranslator from 'corteza-webapp-compose/src/components/Admin/Module/FieldSelectTranslator'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    FieldSelectTranslator,
  },

  extends: base,

  data () {
    return {
      newOption: { value: undefined, text: undefined, new: true },
      selectOptions: [
        { text: this.$t('kind.select.optionType.default'), value: 'default' },
        { text: this.$t('kind.select.optionType.multiple'), value: 'multiple' },
        { text: this.$t('kind.select.optionType.each'), value: 'each' },
      ],
    }
  },

  computed: {
    /**
     * Determines if newly entered option is empty
     * @returns {Boolean}
     */
    newEmpty () {
      return !this.newOption.text || !this.newOption.value
    },

    /**
     * Determines the state of new select option
     * @returns {Boolean|null}
     */
    newOptState () {
      // No duplicates
      if (this.f.options.options.find(({ text, value }) => text === this.newOption.text || value === this.newOption.value)) {
        return false
      }
      return null
    },

    isNew () {
      return this.module.moduleID === NoID || this.field.fieldID === NoID
    },
  },

  created () {
    if (!this.f) {
      this.f.options = { options: [] }
    } else if (!this.f.options.options) {
      this.f.options.options = []
    }
  },

  methods: {
    handleAddOption () {
      if (this.newOption.value) {
        this.f.options.options.push(this.newOption)
        this.newOption = { value: undefined, text: undefined, new: true }
      }
    },
  },
}
</script>
