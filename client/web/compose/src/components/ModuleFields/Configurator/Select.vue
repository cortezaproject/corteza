<template>
  <b-row no-gutters>
    <b-col>
      <b-form-group
        :label="$t('kind.select.optionsLabel')"
      >
        <b-table-simple
          borderless
          small
          responsive="lg"
        >
          <b-thead>
            <b-tr>
              <b-th />

              <b-th
                class="text-primary"
              >
                {{ $t('kind.select.optionValuePlaceholder') }}
              </b-th>

              <b-th
                class="text-primary"
              >
                {{ $t('kind.select.optionLabelPlaceholder') }}
              </b-th>

              <b-th />
            </b-tr>
          </b-thead>

          <draggable
            v-model="f.options.options"
            group="sort"
            handle=".grab"
            tag="tbody"
          >
            <b-tr
              v-for="(option, index) in f.options.options"
              :key="index"
            >
              <b-td class="d-flex align-items-center justify-content-center">
                <font-awesome-icon
                  :icon="['fas', 'bars']"
                  class="grab text-light"
                />
              </b-td>
              <b-td
                style="min-width: 200px;"
              >
                <b-form-input
                  v-model.trim="f.options.options[index].value"
                  plain
                  size="sm"
                  :placeholder="$t('kind.select.optionValuePlaceholder')"
                />
              </b-td>
              <b-td
                style="min-width: 200px;"
              >
                <b-input-group>
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
                  </b-input-group-append>
                </b-input-group>
              </b-td>

              <b-td class="d-flex align-items-center justify-content-center">
                <b-button
                  variant="outline-danger"
                  class="border-0"
                  @click.prevent="f.options.options.splice(index, 1)"
                >
                  <font-awesome-icon :icon="['far', 'trash-alt']" />
                </b-button>
              </b-td>
            </b-tr>

            <b-tr>
              <b-td />
              <b-td>
                <b-button
                  variant="primary px-3"
                  size="md"
                  @click.prevent="handleAddOption"
                >
                  + {{ $t('kind.select.optionAdd') }}
                </b-button>
              </b-td>
            </b-tr>
          </draggable>
        </b-table-simple>
      </b-form-group>

      <b-form-group
        v-if="f.isMulti"
      >
        <label class="d-block">{{ $t('kind.select.optionType.label') }}</label>
        <b-form-radio-group
          v-model="f.options.selectType"
          :options="selectOptions"
          stacked
          @change="onUpdateIsUniqueMultiValue"
        />
        <b-form-checkbox
          v-if="f.options.selectType !== 'multiple'"
          v-model="f.options.isUniqueMultiValue"
          :value="false"
          :unchecked-value="true"
          class="mt-2"
        >
          {{ $t('kind.select.allow-duplicates') }}
        </b-form-checkbox>
      </b-form-group>
    </b-col>
  </b-row>
</template>

<script>
import base from './base'
import Draggable from 'vuedraggable'
import { NoID } from '@cortezaproject/corteza-js'
import FieldSelectTranslator from 'corteza-webapp-compose/src/components/Admin/Module/FieldSelectTranslator'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    FieldSelectTranslator,
    Draggable,
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
      this.f.options.options.push({
        value: undefined,
        text: undefined,
        new: true,
      })
    },

    onUpdateIsUniqueMultiValue () {
      if (this.f.options.selectType === 'multiple') {
        this.f.options.isUniqueMultiValue = true
      }
    },
  },
}
</script>
