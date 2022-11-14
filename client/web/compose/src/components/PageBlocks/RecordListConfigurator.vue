<template>
  <div>
    <b-tab
      :title="$t('recordList.label')"
    >
      <b-form-group
        class="form-group"
        :label="$t('general.module')"
      >
        <b-form-select
          v-model="options.moduleID"
          :options="moduleOptions"
          text-field="name"
          value-field="moduleID"
          required
        />
        <b-form-text class="text-secondary small">
          <i18next
            path="recordList.moduleFootnote"
            tag="label"
          >
            <router-link :to="{ name: 'admin.pages'}">
              {{ $t('recordList.recordPage') }}
            </router-link>
          </i18next>
        </b-form-text>
      </b-form-group>

      <b-form-group
        v-if="recordListModule"
        :label="$t('module:general.fields')"
      >
        <field-picker
          :module="recordListModule"
          :fields.sync="options.fields"
          style="max-height: 40vh;"
        />
      </b-form-group>

      <b-form-group
        v-if="recordListModule"
        horizontal
        :label-cols="3"
        breakpoint="md"
      >
        <b-form-checkbox
          v-model="options.editable"
          :disabled="disableInlineEditor"
        >
          {{ $t('recordList.record.inlineEditorAllow') }}
        </b-form-checkbox>
      </b-form-group>

      <div
        v-if="options.editable"
      >
        <b-form-group
          v-if="recordListModule && options.editable"
          :label="$t('recordList.editFields')"
          label-size="lg"
          class="mb-0"
        >
          <field-picker
            :module="recordListModule"
            :fields.sync="options.editFields"
            :field-subset="options.fields"
            disable-system-fields
            style="max-height: 40vh;"
          />
        </b-form-group>

        <b-form-group
          horizontal
          :label-cols="3"
          breakpoint="md"
          :label="$t('recordList.refField.label')"
        >
          <b-form-select
            v-model="options.refField"
            required
          >
            <option :value="undefined">
              {{ $t('general.label.none') }}
            </option>
            <option
              v-for="field in parentFields"
              :key="field.fieldID"
              :value="field.name"
            >
              {{ field.name }}
            </option>
          </b-form-select>
          <b-form-text class="text-secondary small">
            {{ $t('recordList.refField.footnote') }}
          </b-form-text>
        </b-form-group>

        <b-form-group
          horizontal
          :label-cols="3"
          breakpoint="md"
          :label="$t('recordList.positionField.label')"
        >
          <b-form-select v-model="options.positionField">
            <option :value="undefined">
              {{ $t('general.label.none') }}
            </option>
            <option
              v-for="field in positionFields"
              :key="field.fieldID"
              :value="field.name"
            >
              {{ field.label || field.name }}
            </option>
          </b-form-select>
          <b-form-text class="text-secondary small">
            {{ $t('recordList.positionField.footnote') }}
          </b-form-text>
        </b-form-group>

        <b-form-group
          v-if="options.positionField"
          horizontal
          :label-cols="3"
          breakpoint="md"
        >
          <b-form-checkbox v-model="options.draggable">
            {{ $t('recordList.record.draggable') }}
          </b-form-checkbox>
        </b-form-group>
      </div>

      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        :label="$t('recordList.record.newLabel')"
      >
        <b-form-checkbox v-model="options.hideAddButton">
          {{ $t('recordList.record.hideAddButton') }}
        </b-form-checkbox>
        <b-form-checkbox v-model="options.hideImportButton">
          {{ $t('recordList.record.hideImportButton') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-if="onRecordPage"
          v-model="options.linkToParent"
        >
          {{ $t('recordList.record.linkToParent') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        :label="$t('recordList.record.prefilterLabel')"
      >
        <b-form-textarea
          v-model="options.prefilter"
          :placeholder="$t('recordList.record.prefilterPlaceholder')"
        />
        <b-form-text>
          <i18next
            path="recordList.record.prefilterFootnote"
            tag="label"
          >
            <code>${recordID}</code>
            <code>${ownerID}</code>
            <code>${userID}</code>
          </i18next>
        </b-form-text>
        <b-form-checkbox v-model="options.hideSearch">
          {{ $t('recordList.record.prefilterHideSearch') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        v-if="!options.positionField"
        horizontal
        :label-cols="3"
        breakpoint="md"
        :label="$t('recordList.record.presortLabel')"
      >
        <c-input-presort
          v-model="options.presort"
          :fields="recordListModuleFields"
          :labels="{
            add: $t('general:label.add'),
            ascending: $t('general:label.ascending'),
            descending: $t('general:label.descending'),
            none: $t('general:label.none'),
            placeholder: $t('recordList.record.presortPlaceholder'),
            footnote: $t('recordList.record.presortFootnote'),
            toggleInput: $t('recordList.record.presortToggleInput'),
          }"
          allow-text-input
          class="mb-2"
        />
        <b-form-checkbox v-model="options.hideSorting">
          {{ $t('recordList.record.presortHideSort') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        v-if="!options.editable"
        horizontal
        :label-cols="3"
        breakpoint="md"
        :label="$t('recordList.record.perPage')"
      >
        <b-form-input
          v-model.number="options.perPage"
          type="number"
          class="mb-2"
        />
        <b-form-checkbox v-model="options.hidePaging">
          {{ $t('recordList.record.hidePaging') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-if="!options.hidePaging"
          v-model="options.fullPageNavigation"
        >
          {{ $t('recordList.record.fullPageNavigation') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-if="!options.hidePaging"
          v-model="options.showTotalCount"
        >
          {{ $t('recordList.record.showTotalCount') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        class="mt-4"
      >
        <b-form-checkbox v-model="options.allowExport">
          {{ $t('recordList.export.allow') }}
        </b-form-checkbox>
      </b-form-group>
      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        class="mt-4"
      >
        <b-form-checkbox v-model="options.selectable">
          {{ $t('recordList.selectable') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.openInNewTab"
        >
          {{ $t('recordList.record.openInNewTab') }}
        </b-form-checkbox>
      </b-form-group>
      <b-form-group
        horizontal
        :label-cols="3"
        breakpoint="md"
        class="mt-4"
      >
        <b-form-checkbox v-model="options.hideRecordReminderButton">
          {{ $t('recordList.hideRecordReminderButton') }}
        </b-form-checkbox>
        <b-form-checkbox v-model="options.hideRecordCloneButton">
          {{ $t('recordList.hideRecordCloneButton') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.hideRecordEditButton"
        >
          {{ $t('recordList.hideRecordEditButton') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.hideRecordViewButton"
        >
          {{ $t('recordList.hideRecordViewButton') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.hideRecordPermissionsButton"
        >
          {{ $t('recordList.hideRecordPermissionsButton') }}
        </b-form-checkbox>
      </b-form-group>
    </b-tab>

    <automation-tab
      v-bind="$props"
      :module="recordListModule"
      :buttons.sync="options.selectionButtons"
    />
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import base from './base'
import AutomationTab from './Shared/AutomationTab'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'
import { components } from '@cortezaproject/corteza-vue'
const { CInputPresort } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'RecordList',

  components: {
    AutomationTab,
    FieldPicker,
    CInputPresort,
  },

  extends: base,

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      modules: 'module/set',
      pages: 'page/set',
    }),

    moduleOptions () {
      return [
        { moduleID: NoID, name: this.$t('general.label.none') },
        ...this.modules,
      ]
    },

    recordListModule () {
      if (this.options.moduleID !== NoID) {
        return this.getModuleByID(this.options.moduleID)
      } else {
        return undefined
      }
    },

    recordListModuleFields () {
      if (this.recordListModule) {
        return [
          ...this.recordListModule.fields,
          ...this.recordListModule.systemFields().map(sf => {
            sf.label = this.$t(`field:system.${sf.name}`)
            return sf
          }),
        ].map(({ name, label }) => ({ name, label }))
      }

      return []
    },

    onRecordPage () {
      return this.page && this.page.moduleID !== NoID
    },

    recordListModuleRecordPage () {
      // Relying on pages having unique moduleID,
      if (this.options.moduleID !== NoID) {
        return this.pages.find(p => p.moduleID === this.options.moduleID)
      } else {
        return undefined
      }
    },

    parentFields () {
      if (this.recordListModule) {
        return this.recordListModule.fields.filter(({ kind, isMulti, options }) => {
          if (kind === 'Record' && !isMulti && this.record) {
            return options.moduleID === this.record.moduleID
          }
        })
      }
      return []
    },

    positionFields () {
      if (this.recordListModule) {
        return this.recordListModule.fields.filter(({ kind, isMulti }) => kind === 'Number' && !isMulti)
      }
      return []
    },

    /*
     Inline record editor is disabled if:
      - An inline record editor for the same module already exists
      - Record list module doesn't have record page (inline record autoselected and disabled)
    */
    disableInlineEditor () {
      const thisModuleID = this.options.moduleID

      // Finds another inline editor block with the same recordListModulea as this one
      const otherInlineWithSameModule = !!this.page.blocks.find(({ kind, options }, index) => {
        if (this.blockIndex !== index) {
          return kind === 'RecordList' && options.editable && options.moduleID === thisModuleID
        }
      })

      return otherInlineWithSameModule || !this.recordListModuleRecordPage
    },
  },

  watch: {
    'options.moduleID' (newModuleID) {
      // Every time moduleID changes
      this.options.fields = []
      this.options.editable = false

      // If recordListModule doesn't have record page, auto check inline record editor
      if (newModuleID !== NoID) {
        if (!this.recordListModuleRecordPage) {
          this.options.editable = true
        }
      }
    },

    'options.editable' (value) {
      this.options.editFields = []
      this.options.positionField = undefined

      if (value) {
        this.options.hideRecordEditButton = true
        this.options.hideRecordViewButton = true
        this.options.hidePaging = true
        let f = null
        if (this.module && this.module.moduleID) f = this.recordListModule.fields.find(({ options: { moduleID } }) => moduleID === this.module.moduleID)
        this.options.refField = f ? f.name : undefined
      } else {
        this.options.refField = undefined
      }
    },

    'options.positionField' (v) {
      if (!v) {
        this.options.draggable = false
      }

      this.options.hideSorting = true
      this.options.presort = ''
    },

    'options.fields' (fields) {
      this.options.editFields = this.options.editFields.filter(a => fields.some(b => a.name === b.name))
    },
  },
}
</script>
