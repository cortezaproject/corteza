<template>
  <b-row>
    <b-col
      cols="12"
      lg="6"
    >
      <b-form-group
        :label="$t('calendar.recordFeed.moduleLabel')"
        label-class="text-primary"
      >
        <b-input-group>
          <c-input-select
            v-model="feed.options.moduleID"
            :options="modules"
            :reduce="o => o.moduleID"
            default-value="0"
            :placeholder="$t('calendar.recordFeed.modulePlaceholder')"
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
          :label="$t('calendar.recordFeed.titleLabel')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="feed.titleField"
            :options="titleFields"
            :get-option-key="getOptionEventFieldKey"
            :get-option-label="getOptionEventFieldLabel"
            :reduce="o => o.name"
            :placeholder="$t('calendar.recordFeed.titlePlaceholder')"
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('calendar.recordFeed.eventStartFieldLabel')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="feed.startField"
            :options="dateFields"
            :get-option-key="getOptionEventFieldKey"
            :get-option-label="getOptionEventFieldLabel"
            :reduce="o => o.name"
            :placeholder="$t('calendar.recordFeed.eventStartFieldPlaceholder')"
          />
        </b-form-group>
      </b-col>

      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('calendar.recordFeed.eventEndFieldLabel')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="feed.endField"
            :options="dateFields"
            :get-option-key="getOptionEventFieldKey"
            :get-option-label="getOptionEventFieldLabel"
            :reduce="o => o.name"
            :disabled="feed.allDay"
            :placeholder="$t('calendar.recordFeed.eventEndFieldPlaceholder')"
          />

          <b-form-checkbox
            v-model="feed.allDay"
            :value="true"
            :unchecked-value="false"
            class="mt-1"
          >
            {{ $t('calendar.recordFeed.eventAllDay') }}
          </b-form-checkbox>
        </b-form-group>
      </b-col>

      <b-col cols="12">
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
        lg="6"
      >
        <b-form-group
          :label="$t('calendar.recordFeed.colorLabel')"
          label-class="text-primary"
        >
          <c-input-color-picker
            v-model="feed.options.color"
            :translations="{
              modalTitle: $t('calendar.recordFeed.colorPicker'),
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
  </b-row>
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
import { compose, NoID } from '@cortezaproject/corteza-js'
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

  props: {
    record: {
      type: compose.Record,
      required: false,
      default: undefined,
    },

    page: {
      type: compose.Page,
      required: true,
    },
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

    moduleOptions () {
      return [
        { moduleID: '0', name: this.$t('calendar.recordFeed.modulePlaceholder') },
        ...this.modules,
      ]
    },

    /**
     * Determines available title fields based on the given module.
     * @returns {Array}
     */
    titleFields () {
      if (!this.module) {
        return []
      }

      return [...this.module.fields]
        .filter(f => ['String', 'Email', 'Url'].includes(f.kind))
        .sort((a, b) => a.label.localeCompare(b.label))
    },

    /**
     * Determines available date fields based on the given module.
     * Currently ignores multi-fields
     * @returns {Array}
     */
    dateFields () {
      if (!this.module) {
        return []
      }

      const moduleFields = this.module.fields.slice().sort((a, b) => a.label.localeCompare(b.label))

      return [
        ...moduleFields,
        ...this.module.systemFields().map(sf => {
          sf.label = this.$t(`field:system.${sf.name}`)
          return sf
        }),
      ].filter(f => f.kind === 'DateTime' && !f.isMulti)
    },

    themeSettings () {
      return this.$Settings.get('ui.studio.themes', [])
    },

    isRecordPage () {
      return this.page && this.page.moduleID !== NoID
    },

    recordAutoCompleteParams () {
      return this.processRecordAutoCompleteParams({ operators: true })
    },
  },

  methods: {
    onModuleChange () {
      this.feed.titleField = ''
      this.feed.startField = ''
      this.feed.endField = ''
    },

    getOptionEventFieldKey ({ name }) {
      return name
    },

    getOptionEventFieldLabel ({ name, label }) {
      return label || name
    },
  },
}
</script>
