<template>
  <b-row>
    <template v-if="feed.options">
      <b-col cols="12">
        <b-form-group
          :label="$t('geometry.recordFeed.moduleLabel')"
          label-class="text-primary"
        >
          <b-input-group>
            <c-input-select
              v-model="feed.options.moduleID"
              :options="modules"
              :reduce="o => o.moduleID"
              :placeholder="$t('calendar.recordFeed.modulePlaceholder')"
              default-value="0"
              label="name"
              @input="onModuleChange"
            />
          </b-input-group>
        </b-form-group>
      </b-col>

      <template v-if="module">
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('geometry.recordFeed.geometryFieldLabel')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="feed.geometryField"
              :options="geometryFields"
              :get-option-key="getOptionGeometryAndTitleFieldKey"
              :get-option-label="getOptionGeometryAndTitleFieldLabel"
              :reduce="o => o.name"
              :placeholder="$t('geometry.recordFeed.geometryFieldPlaceholder')"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('geometry.recordFeed.titleLabel')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="feed.titleField"
              :options="titleFields"
              :get-option-key="getOptionGeometryAndTitleFieldKey"
              :get-option-label="getOptionGeometryAndTitleFieldLabel"
              :reduce="o => o.name"
              :placeholder="$t('geometry.recordFeed.titlePlaceholder')"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
        >
          <b-form-group
            :label="$t('calendar.recordFeed.prefilterLabel')"
            label-class="text-primary"
          >
            <c-input-expression
              v-model="feed.options.prefilter"
              height="3.688rem"
              lang="javascript"
              :suggestion-params="recordAutoCompleteParams"
              :placeholder="$t('calendar.recordFeed.prefilterPlaceholder')"
            />

            <i18next
              path="interpolationFootnote"
              tag="small"
              class="text-muted"
            >
              <code>${record.values.fieldName}</code>
              <code>${recordID}</code>
              <code>${ownerID}</code>
              <span><code>${userID}</code>, <code>${user.name}</code></span>
            </i18next>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="4"
        >
          <b-form-group
            :label="$t('geometry.recordFeed.displayMarker')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="feed.displayMarker"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="4"
        >
          <b-form-group
            :label="$t('geometry.recordFeed.displayPolygon')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="feed.displayPolygon"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="4"
        >
          <b-form-group
            :label="$t('geometry.recordFeed.colorLabel')"
            label-class="text-primary"
          >
            <c-input-color-picker
              v-model="feed.options.color"
              :translations="{
                modalTitle: $t('geometry.recordFeed.colorPicker'),
                light: $t('general:themes.labels.light'),
                dark: $t('general:themes.labels.dark'),
                cancelBtnLabel: $t('general:label.cancel'),
                saveBtnLabel: $t('general:label.saveAndClose')
              }"
              :theme-settings="themeSettings"
            />
          </b-form-group>
        </b-col>
      </template>
    </template>
  </b-row>
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
import autocomplete from 'corteza-webapp-compose/src/mixins/autocomplete.js'

const { CInputColorPicker, CInputExpression } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CInputColorPicker,
    CInputExpression,
  },

  extends: base,

  mixins: [autocomplete],

  data () {
    return {
      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
    }
  },

  computed: {
    /**
     * Finds the module, this feed configurator should use
     * @returns {Module|undefined}
     */
    module () {
      if (!(this.feed.options || {}).moduleID) {
        return
      }

      return this.modules.find(({ moduleID }) => moduleID === this.feed.options.moduleID)
    },

    /**
     * Determines available title fields based on the given module.
     * @returns {Array}
     */
    titleFields () {
      if (!this.module) {
        return []
      }

      return this.module.fields
        .filter(f => [
          'DateTime',
          'Select',
          'Number',
          'Bool',
          'String',
          'Record',
          'User',
        ].includes(f.kind) && f.label)
        .toSorted((a, b) => a.label.localeCompare(b.label))
    },

    /**
     * Determines available geometry fields based on the given module.
     * Currently ignores multi-fields
     * @returns {Array}
     */
    geometryFields () {
      if (!this.module) {
        return []
      }

      return [
        ...this.module.fields,
        ...this.module.systemFields().map(sf => {
          sf.label = this.$t(`field:system.${sf.name}`)
          return sf
        })]
        .filter(f => f.kind === 'Geometry')
        .toSorted((a, b) => a.label.localeCompare(b.label))
    },

    themeSettings () {
      return this.$Settings.get('ui.studio.themes', [])
    },
  },

  methods: {
    onModuleChange () {
      this.feed.geometryField = ''
      this.feed.titleField = ''
    },

    getOptionGeometryAndTitleFieldKey ({ name }) {
      return name
    },

    getOptionGeometryAndTitleFieldLabel ({ name, label }) {
      return label || name
    },
  },
}
</script>
