<template>
  <b-row no-gutters>
    <b-col>
      <b-form-group
        :label="$t('kind.select.optionType.label')"
        label-class="text-primary"
      >
        <b-form-radio-group
          v-model="f.options.selectType"
          :options="selectOptions"
          stacked
          @change="updateIsUniqueMultiValue"
        />
      </b-form-group>

      <b-form-group
        v-if="shouldAllowDuplicates"
      >
        <b-form-checkbox
          v-model="f.options.isUniqueMultiValue"
          :value="false"
          :unchecked-value="true"
        >
          {{ $t('kind.select.allow-duplicates') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        :label="$t('kind.select.displayType.label')"
        label-class="text-primary"
      >
        <b-form-radio-group
          v-model="f.options.displayType"
          :options="displayOptions"
          stacked
        />
      </b-form-group>

      <b-form-group
        :label="$t('kind.select.optionsLabel')"
        label-class="text-primary"
      >
        <c-form-table-wrapper
          :labels="{ addButton: $t('general:label.add') }"
          @add-item="handleAddOption"
        >
          <b-table-simple
            borderless
            small
            responsive
          >
            <b-thead>
              <b-tr>
                <b-th v-if="f.options.options.length > 0" />

                <b-th
                  class="text-primary"
                  style="min-width: 200px;"
                >
                  {{ $t('kind.select.options.value') }}
                </b-th>

                <b-th
                  class="text-primary"
                >
                  {{ $t('kind.select.options.label') }}
                </b-th>

                <b-th
                  v-if="f.options.displayType === 'badge'"
                  class="text-primary"
                >
                  {{ $t('kind.select.options.style.textColor') }}
                </b-th>

                <b-th
                  v-if="f.options.displayType === 'badge'"
                  class="text-primary"
                >
                  {{ $t('kind.select.options.style.backgroundColor') }}
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
                <b-td class="align-middle text-center">
                  <font-awesome-icon
                    :icon="['fas', 'bars']"
                    class="grab text-secondary"
                  />
                </b-td>
                <b-td
                  style="min-width: 200px;"
                >
                  <b-form-input
                    v-model.trim="f.options.options[index].value"
                    plain
                    :placeholder="$t('kind.select.options.value')"
                    :state="f.options.options[index].value ? null : false"
                  />
                </b-td>
                <b-td
                  style="min-width: 200px;"
                >
                  <b-input-group>
                    <b-form-input
                      v-model.trim="f.options.options[index].text"
                      plain
                      :placeholder="option.value || $t('kind.select.options.label')"
                    />

                    <b-input-group-append>
                      <field-select-translator
                        v-if="field"
                        :field="field"
                        :module="module"
                        :highlight-key="`meta.options.${option.value}.text`"
                        size="sm"
                        :disabled="isNew || option.new"
                      />
                    </b-input-group-append>
                  </b-input-group>
                </b-td>

                <b-td
                  v-if="f.options.displayType === 'badge'"
                  style="min-width: 120px;"
                >
                  <c-input-color-picker
                    v-model="f.options.options[index].style.textColor"
                    :default-value="defaultTextColor"
                    :theme-settings="themeSettings"
                    :translations="{
                      modalTitle: $t('kind.select.options.style.textColor'),
                      cancelBtnLabel: $t('general:label.cancel'),
                      saveBtnLabel: $t('general:label.saveAndClose')
                    }"
                  />
                </b-td>

                <b-td
                  v-if="f.options.displayType === 'badge'"
                  style="min-width: 130px;"
                >
                  <c-input-color-picker
                    v-model="f.options.options[index].style.backgroundColor"
                    :default-value="defaultBackgroundColor"
                    :theme-settings="themeSettings"
                    :translations="{
                      modalTitle: $t('kind.select.options.style.backgroundColor'),
                      defaultBtnLabel: $t('general:label.default'),
                      cancelBtnLabel: $t('general:label.cancel'),
                      saveBtnLabel: $t('general:label.saveAndClose')
                    }"
                  />
                </b-td>

                <b-td class="align-middle text-right">
                  <c-input-confirm
                    show-icon
                    @confirmed="f.options.options.splice(index, 1)"
                  />
                </b-td>
              </b-tr>
            </draggable>
          </b-table-simple>
        </c-form-table-wrapper>
      </b-form-group>
    </b-col>
  </b-row>
</template>

<script>
import base from './base'
import Draggable from 'vuedraggable'
import { NoID } from '@cortezaproject/corteza-js'
import FieldSelectTranslator from 'corteza-webapp-compose/src/components/Admin/Module/FieldSelectTranslator'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    FieldSelectTranslator,
    Draggable,
    CInputColorPicker,
  },

  extends: base,

  data () {
    return {
      newOption: { value: undefined, text: undefined, new: true },

      selectTypes: [
        { text: this.$t('kind.select.optionType.default'), value: 'default', allowDuplicates: true },
        { text: this.$t('kind.select.optionType.multiple'), value: 'multiple', onlyMulti: true },
        { text: this.$t('kind.select.optionType.each'), value: 'each', allowDuplicates: true, onlyMulti: true },
        { value: 'list' },
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

    selectOptions () {
      const selectOptions = this.selectTypes.map((o) => {
        if (o.value === 'list') {
          o.text = this.$t(`kind.select.optionType.${this.f.isMulti ? 'checkbox' : 'radio'}`)
        }

        return o
      })

      if (this.f.isMulti) {
        return selectOptions
      }

      return selectOptions.filter(({ onlyMulti }) => !onlyMulti)
    },

    displayOptions () {
      return [
        { text: this.$t('kind.select.displayType.text'), value: 'text' },
        { text: this.$t('kind.select.displayType.badge'), value: 'badge' },
      ]
    },

    shouldAllowDuplicates () {
      if (!this.f.isMulti) return false

      const { allowDuplicates } = this.selectTypes.find(({ value }) => value === this.f.options.selectType) || {}
      return !!allowDuplicates
    },

    themeSettings () {
      return this.$Settings.get('ui.studio.themes', [])
    },

    defaultTextColor () {
      return getComputedStyle(document.documentElement).getPropertyValue('--dark')
    },

    defaultBackgroundColor () {
      return getComputedStyle(document.documentElement).getPropertyValue('--extra-light')
    },
  },

  created () {
    if (!this.f) {
      this.f.options = { options: [] }
    } else if (!this.f.options.options) {
      this.f.options.options = []
    }
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    handleAddOption () {
      const option = this.f.createSelectOption()
      option.new = true

      this.f.options.options.push(option)
    },

    updateIsUniqueMultiValue (value) {
      const { allowDuplicates = false } = this.selectTypes.find(({ value: v }) => v === value) || {}
      if (!allowDuplicates) {
        this.f.options.isUniqueMultiValue = true
      }
    },

    setDefaultValues () {
      this.newOption = {}
      this.selectTypes = []
    },
  },
}
</script>
