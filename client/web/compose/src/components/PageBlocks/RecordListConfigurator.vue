<template>
  <div>
    <b-tab
      data-test-id="record-list-configurator"
      :title="$t('recordList.label')"
      no-body
    >
      <div
        class="px-3 pt-3"
      >
        <h5 class="mb-3">
          {{ $t('recordList.record.generalLabel') }}
        </h5>
        <b-row>
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              class="form-group text-primary"
              variant="primary"
              :label="$t('general.module')"
            >
              <b-form-select
                v-model="options.moduleID"
                :options="moduleOptions"
                text-field="name"
                value-field="moduleID"
                required
              />
            </b-form-group>
          </b-col>

          <b-col
            v-if="onRecordPage || options.editable"
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('recordList.record.inlineEditorAllow')"
              label-class="text-primary"
            >
              <c-input-checkbox
                v-model="options.editable"
                switch
                :labels="checkboxLabel"
              />
            </b-form-group>
          </b-col>
        </b-row>
      </div>

      <template v-if="recordListModule">
        <hr>

        <div class="px-3">
          <div class="mb-3">
            <h5 class="mb-1">
              {{ $t('module:general.fields') }}
            </h5>
            <small class="text-muted">
              {{ $t('recordList.moduleFieldsFootnote') }}
            </small>
          </div>

          <b-row>
            <b-col
              cols="12"
            >
              <field-picker
                :module="recordListModule"
                :fields.sync="options.fields"
                style="height: 40vh;"
              />
            </b-col>
          </b-row>
        </div>

        <hr>

        <div
          v-if="options.editable"
          class="px-3"
        >
          <h5 class="mb-3">
            {{ $t('recordList.record.inlineEditor') }}
          </h5>

          <b-form-group
            v-if="recordListModule && options.editable"
            :label="$t('recordList.editFields')"
            class="mb-0 text-primary"
          >
            <field-picker
              :module="recordListModule"
              :fields.sync="options.editFields"
              :field-subset="options.fields"
              disable-system-fields
              style="height: 40vh;"
            />
          </b-form-group>

          <b-row
            class="mt-3"
          >
            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.refField.label')"
                label-class="text-primary"
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
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.positionField.label')"
                label-class="text-primary"
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
            </b-col>

            <b-col
              cols="12"
            >
              <b-form-group
                v-if="options.positionField"
                :label="$t('recordList.record.draggable')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.draggable"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>
          </b-row>
        </div>

        <hr v-if="options.editable">

        <div
          class="px-3"
        >
          <h5 class="mb-3">
            {{ $t('recordList.record.prefilterLabel') }}
          </h5>

          <b-row>
            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.filterHide')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.hideFiltering"
                  switch
                  invert
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.prefilterHideSearch')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.hideSearch"
                  switch
                  invert
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col>
              <b-form-group
                :label="$t('recordList.record.prefilterCommand')"
                label-class="text-primary"
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
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col>
              <b-form-group
                :label="$t('recordList.filter.presets')"
              >
                <b-table-simple
                  v-if="recordListModule"
                  borderless
                  small
                  responsive="lg"
                  class="mb-1"
                >
                  <b-thead>
                    <b-tr>
                      <b-th
                        class="text-primary"
                        style="min-width: 300px;"
                      >
                        {{ $t('recordList.filter.name.label') }}
                      </b-th>
                      <b-th
                        class="text-primary"
                        style="width: 45%; min-width: 250px;"
                      >
                        {{ $t('recordList.filter.role.label') }}
                      </b-th>

                      <b-th
                        style="width: 100px;"
                      />
                    </b-tr>
                  </b-thead>
                  <b-tbody>
                    <b-tr
                      v-for="(filter, index) in options.filterPresets"
                      :key="index"
                    >
                      <b-td>
                        <b-input-group>
                          <b-form-input
                            v-model="filter.name"
                            :placeholder="$t('recordList.filter.name.placeholder')"
                          />

                          <b-input-group-append class="border-0">
                            <record-list-filter
                              class="d-print-none"
                              :target="`record-filter-${index}`"
                              :namespace="namespace"
                              :module="recordListModule"
                              :selected-field="recordListModule.fields[0]"
                              :record-list-filter="filter.filter"
                              variant="primary"
                              button-class="px-2 pt-2 text-white"
                              button-style="padding-bottom: calc(0.5rem - 2px);"
                              @filter="(filter) => onFilter(filter, index)"
                            />
                          </b-input-group-append>
                        </b-input-group>
                      </b-td>

                      <b-td>
                        <vue-select
                          v-model="filter.roles"
                          :options="roleOptions"
                          :get-option-label="getRoleLabel"
                          :get-option-key="getOptionKey"
                          :placeholder="$t('recordList.filter.role.placeholder')"
                          :reduce="role => role.roleID"
                          :calculate-position="calculateDropdownPosition"
                          append-to-body
                          multiple
                          class="bg-white"
                        />
                      </b-td>

                      <b-td
                        class="text-center align-middle pr-2"
                      >
                        <c-input-confirm
                          @confirmed="options.filterPresets.splice(index, 1)"
                        />
                      </b-td>
                    </b-tr>
                  </b-tbody>
                </b-table-simple>

                <b-button
                  variant="primary"
                  size="sm"
                  class="ml-1"
                  @click="addFilterPreset"
                >
                  {{ $t('recordList.filter.addFilter') }}
                </b-button>
              </b-form-group>
            </b-col>
          </b-row>
        </div>
        <hr>

        <div
          v-if="!options.positionField"
          class="px-3"
        >
          <h5 class="mb-3">
            {{ $t('recordList.record.presortLabel') }}
          </h5>
          <b-row class="mb-3">
            <b-col>
              <b-form-group
                :label="$t('recordList.record.presortHideSort')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.hideSorting"
                  switch
                  invert
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col>
              <b-form-group
                :label="$t('recordList.record.presortInputLabel')"
                label-class="text-primary"
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
              </b-form-group>
            </b-col>
          </b-row>
        </div>
        <hr v-if="!options.positionField">

        <div class="px-3">
          <h5 class="mb-3">
            {{ $t('recordList.record.pagingLabel') }}
          </h5>

          <b-row
            class="mb-3"
          >
            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.hidePaging')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.hidePaging"
                  switch
                  invert
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                horizontal
                breakpoint="md"
                :label="$t('recordList.record.perPage')"
                label-class="text-primary"
              >
                <b-form-input
                  v-model.number="options.perPage"
                  data-test-id="input-records-per-page"
                  type="number"
                  class="mb-2"
                />
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.fullPageNavigation')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.fullPageNavigation"
                  switch
                  :labels="checkboxLabel"
                  data-test-id="hide-page-navigation"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.showTotalCount')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.showTotalCount"
                  data-test-id="show-total-record-count"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>
          </b-row>
        </div>

        <hr>

        <div
          class="px-3"
        >
          <h5 class="mb-3">
            {{ $t('recordList.record.recordsLabel') }}
          </h5>

          <b-row>
            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.recordDisplayOptions')"
                label-class="text-primary"
              >
                <b-form-select
                  v-model="options.recordDisplayOption"
                  :options="recordDisplayOptions"
                />
              </b-form-group>
            </b-col>
            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.recordSelectorDisplayOptions')"
                label-class="text-primary"
              >
                <b-form-select
                  v-model="options.recordSelectorDisplayOption"
                  :options="recordDisplayOptions"
                />
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.hideAddButton')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.hideAddButton"
                  switch
                  invert
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.selectable')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.selectable"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              v-if="onRecordPage"
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.linkToParent')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.linkToParent"
                  switch
                  invert
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.hideImportButton')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.hideImportButton"
                  switch
                  invert
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.export.allow')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.allowExport"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.inlineEdit.enabled')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.inlineRecordEditEnabled"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.enableBulkRecordEdit')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.bulkRecordEditEnabled"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.enableRecordPageNavigation')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.enableRecordPageNavigation"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.showDeletedRecordsOption')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.showDeletedRecordsOption"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.buttons')"
                label-class="text-primary"
              >
                <b-form-checkbox
                  v-model="options.hideRecordViewButton"
                >
                  {{ $t('recordList.hideRecordViewButton') }}
                </b-form-checkbox>

                <b-form-checkbox
                  v-model="options.hideRecordEditButton"
                >
                  {{ $t('recordList.hideRecordEditButton') }}
                </b-form-checkbox>

                <b-form-checkbox v-model="options.hideRecordCloneButton">
                  {{ $t('recordList.hideRecordCloneButton') }}
                </b-form-checkbox>

                <b-form-checkbox
                  v-model="options.hideRecordPermissionsButton"
                >
                  {{ $t('recordList.hideRecordPermissionsButton') }}
                </b-form-checkbox>

                <b-form-checkbox v-model="options.hideRecordReminderButton">
                  {{ $t('recordList.hideRecordReminderButton') }}
                </b-form-checkbox>
              </b-form-group>
            </b-col>
          </b-row>
        </div>
      </template>
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
import { VueSelect } from 'vue-select'
import base from './base'
import AutomationTab from './Shared/AutomationTab'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'
import RecordListFilter from 'corteza-webapp-compose/src/components/Common/RecordListFilter'
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
    RecordListFilter,
    VueSelect,
  },

  extends: base,

  data () {
    return {
      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
      roleOptions: [],
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      modules: 'module/set',
      pages: 'page/set',
    }),

    recordDisplayOptions () {
      return [
        { value: 'sameTab', text: this.$t('recordList.record.openInSameTab') },
        { value: 'newTab', text: this.$t('recordList.record.openInNewTab') },
        { value: 'modal', text: this.$t('recordList.record.openInModal') },
      ]
    },

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

      // Finds another inline editor block with the same recordListModule as this one
      const otherInlineWithSameModule = this.blocks.some(({ kind, options }, index) => {
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

  mounted () {
    this.fetchRoles()
  },

  methods: {
    getRoleLabel ({ name }) {
      return name
    },

    async fetchRoles () {
      this.$SystemAPI.roleList().then(({ set: roles = [] }) => {
        this.roleOptions = roles.filter(({ meta }) => !(meta.context && meta.context.resourceTypes))
      })
    },

    onFilter (filter = [], index) {
      this.options.filterPresets[index].filter = filter
    },

    addFilterPreset () {
      this.options.filterPresets.push({
        name: '',
        filter: [],
        roles: [],
      })
    },

    getOptionKey ({ roleID }) {
      return roleID
    },
  },
}
</script>

<style>
.w-fit {
  width: fit-content;
}
</style>
