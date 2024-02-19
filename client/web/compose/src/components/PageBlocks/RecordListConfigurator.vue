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
            :lg="isInlineEditorAllowed ? 6 : 12"
          >
            <b-form-group
              :label="$t('general.module')"
              variant="primary"
              label-class="text-primary"
            >
              <c-input-select
                v-model="options.moduleID"
                :options="modules"
                label="name"
                :reduce="o => o.moduleID"
                :placeholder="$t('recordList.modulePlaceholder')"
                default-value="0"
                required
              />
            </b-form-group>
          </b-col>

          <b-col
            v-if="isInlineEditorAllowed"
            cols="12"
            lg="6"
          >
            <b-form-group
              label-class="d-flex align-items-center text-primary"
            >
              <template #label>
                {{ $t('recordList.record.inlineEditorAllow') }}
                <c-hint
                  :tooltip="$t('recordList.tooltip.performance.impact')"
                  icon-class="text-warning"
                />
              </template>

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
            <h5 class="d-flex align-items-center mb-1">
              {{ $t('module:general.fields') }}
              <c-hint
                :tooltip="$t('recordList.tooltip.performance.moduleFields')"
                icon-class="text-warning"
              />
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
                class="mb-3"
                style="height: 50vh;"
              />
            </b-col>

            <b-col
              cols="12"
              lg="6"
            >
              <b-form-group
                :label="$t('recordList.hideConfigureFieldsButton')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.hideConfigureFieldsButton"
                  switch
                  invert
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              v-if="onRecordPage"
              cols="12"
              lg="6"
            >
              <b-form-group
                :label="$t('recordList.refField.label')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="options.refField"
                  :options="parentFields"
                  :placeholder="$t('general.label.none')"
                  :reduce="f => f.name"
                />

                <b-form-text class="text-secondary small">
                  {{ $t('recordList.refField.footnote') }}
                </b-form-text>
              </b-form-group>
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
            label-class="text-primary"
            class="mb-0"
          >
            <field-picker
              :module="recordListModule"
              :fields.sync="options.editFields"
              :field-subset="options.fields"
              disable-system-fields
              style="height: 50vh;"
            />
          </b-form-group>

          <b-row
            class="mt-3"
          >
            <b-col
              cols="12"
              lg="6"
            >
              <b-form-group
                :label="$t('recordList.positionField.label')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="options.positionField"
                  :placeholder="$t('recordList.positionField.placeholder')"
                  :reduce="f => f.name"
                  label="label"
                />

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
              lg="6"
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
              lg="6"
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

                <i18next
                  path="recordList.record.prefilterFootnote"
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
          </b-row>

          <hr>

          <b-row>
            <b-col>
              <b-form-group
                :label="$t('recordList.record.setCustomFilterPresets')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.customFilterPresets"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>
          </b-row>

          <b-row>
            <b-col>
              <div class="list-background rounded border border-light p-3">
                <b-form-group
                  :label="$t('recordList.filter.presets')"
                  label-class="text-primary"
                  class="mb-0"
                >
                  <b-spinner v-if="fetchingRoles" />

                  <template v-else>
                    <b-table-simple
                      v-if="recordListModule"
                      borderless
                      small
                      responsive="lg"
                      class="mb-0"
                    >
                      <draggable
                        :list.sync="options.filterPresets"
                        group="sort"
                        handle=".grab"
                        tag="tbody"
                      >
                        <b-tr
                          v-for="(filter, index) in options.filterPresets"
                          :key="index"
                        >
                          <b-td
                            class="grab text-center align-middle"
                            style="width: 40px;"
                          >
                            <font-awesome-icon
                              :icon="['fas', 'bars']"
                              class="text-secondary"
                            />
                          </b-td>

                          <b-td
                            class="align-middle"
                            style="min-width: 150px;"
                          >
                            <b-input-group>
                              <b-form-input
                                v-model="filter.name"
                                :placeholder="$t('recordList.filter.name.placeholder')"
                              />

                              <b-input-group-append>
                                <record-list-filter
                                  class="d-print-none"
                                  :target="`record-filter-${index}`"
                                  :namespace="namespace"
                                  :module="recordListModule"
                                  :selected-field="recordListModule.fields[0]"
                                  :record-list-filter="filter.filter"
                                  variant="extra-light"
                                  inactive-icon-class="text-light"
                                  button-class="px-2 pt-2"
                                  button-style="border-top-left-radius: 0; border-bottom-left-radius: 0;"
                                  @filter="(filter) => onFilter(filter, index)"
                                />
                              </b-input-group-append>
                            </b-input-group>
                          </b-td>

                          <b-td
                            class="text-center align-middle"
                            style="min-width: 200px;"
                          >
                            <c-input-role
                              :value="getFilterRoles(filter)"
                              :placeholder="$t('recordList.filter.role.placeholder')"
                              :visible="isRoleVisible"
                              multiple
                              @input="onFilterRoleChange(filter, $event)"
                            />
                          </b-td>

                          <b-td
                            class="text-right align-middle"
                            style="min-width: 80px; width: 80px;"
                          >
                            <c-input-confirm
                              show-icon
                              @confirmed="options.filterPresets.splice(index, 1)"
                            />
                          </b-td>
                        </b-tr>
                      </draggable>
                    </b-table-simple>

                    <b-button
                      variant="primary"
                      size="sm"
                      class="mt-1"
                      @click="addFilterPreset"
                    >
                      <font-awesome-icon
                        :icon="['fas', 'plus']"
                        class="mr-1"
                      />
                      {{ $t('general:label.add') }}
                    </b-button>
                  </template>
                </b-form-group>
              </div>
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

          <b-row>
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
              <div class="list-background rounded border border-light p-3">
                <b-form-group
                  :label="$t('recordList.record.presortInputLabel')"
                  label-class="text-primary"
                  class="mb-0"
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
                  />
                </b-form-group>
              </div>
            </b-col>
          </b-row>
        </div>
        <hr v-if="!options.positionField">

        <div class="px-3">
          <h5 class="mb-3">
            {{ $t('recordList.record.pagingLabel') }}
          </h5>

          <b-row>
            <b-col
              cols="12"
              lg="6"
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
              lg="6"
            >
              <b-form-group
                label-class="d-flex align-items-center text-primary p-0"
              >
                <template #label>
                  {{ $t('recordList.record.fullPageNavigation') }}
                  <c-hint
                    :tooltip="$t('recordList.tooltip.performance.impact')"
                    icon-class="text-warning"
                  />
                </template>

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
              lg="6"
            >
              <b-form-group
                horizontal
                breakpoint="md"
                label-class="d-flex align-items-center text-primary"
              >
                <template #label>
                  {{ $t('recordList.record.perPage') }}
                  <c-hint
                    :tooltip="$t('recordList.tooltip.performance.perPage')"
                    icon-class="text-warning"
                  />
                </template>

                <b-form-input
                  v-model.number="options.perPage"
                  data-test-id="input-records-per-page"
                  type="number"
                  class="mb-2"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              lg="6"
            >
              <b-form-group
                :label="$t('recordList.record.showRecordPerPageOption')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.showRecordPerPageOption"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              lg="6"
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
              lg="6"
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
              lg="6"
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
              lg="6"
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
              lg="6"
            >
              <b-form-group
                :label="$t('recordList.record.addRecordOptions')"
                label-class="text-primary"
              >
                <b-form-select
                  v-model="options.addRecordDisplayOption"
                  :options="recordCreateOptions"
                  :disabled="options.hideAddButton"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              md="6"
            >
              <b-form-group
                :label="$t('recordList.record.editMode')"
                label-class="text-primary"
              >
                <c-input-checkbox
                  v-model="options.openRecordInEditMode"
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
              cols="12"
              lg="6"
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
              lg="6"
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
              lg="6"
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
              lg="6"
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
              lg="6"
            >
              <b-form-group
                label-class="d-flex align-items-center text-primary mb-0"
              >
                <template #label>
                  {{ $t('recordList.enableRecordPageNavigation') }}
                  <c-hint
                    :tooltip="$t('recordList.tooltip.performance.impact')"
                    icon-class="text-warning"
                  />
                </template>
                <c-input-checkbox
                  v-model="options.enableRecordPageNavigation"
                  switch
                  :labels="checkboxLabel"
                />
              </b-form-group>
            </b-col>

            <b-col
              cols="12"
              lg="6"
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
              lg="6"
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
import Draggable from 'vuedraggable'
import base from './base'
import AutomationTab from './Shared/AutomationTab'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'
import RecordListFilter from 'corteza-webapp-compose/src/components/Common/RecordListFilter'
import { components } from '@cortezaproject/corteza-vue'
const { CInputPresort, CInputRole } = components

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
    Draggable,
    CInputRole,
  },

  extends: base,

  data () {
    return {
      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },

      fetchingRoles: false,
      resolvedRoles: {},
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

    recordCreateOptions () {
      return [
        { value: 'sameTab', text: this.$t('recordList.record.createInSameTab') },
        { value: 'newTab', text: this.$t('recordList.record.createInNewTab') },
        { value: 'modal', text: this.$t('recordList.record.createInModal') },
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
            return {
              label: this.$t(`field:system.${sf.name}`),
              name: sf.name === 'recordID' ? 'ID' : sf.name,
            }
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

    isInlineEditorAllowed () {
      return this.recordListModule && (this.onRecordPage || this.options.editable)
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

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    fetchRoles () {
      if (!this.options.filterPresets.length) {
        return
      }

      this.fetchingRoles = true

      const rolesToResolve = this.options.filterPresets.reduce((acc, { roles }) => {
        return acc.concat(roles)
      }, [])

      this.$SystemAPI.roleList({ roleID: rolesToResolve }).then(({ set }) => {
        set.forEach(role => {
          this.resolvedRoles[role.roleID] = role
        })
      }).finally(() => {
        this.fetchingRoles = false
      })
    },

    onFilterRoleChange (filter, roles) {
      roles.forEach(r => {
        if (!this.resolvedRoles[r.roleID]) {
          this.resolvedRoles[r.roleID] = r
        }
      })

      filter.roles = roles.map(({ roleID }) => roleID)
    },

    getFilterRoles (filter) {
      return filter.roles.map(roleID => this.resolvedRoles[roleID])
    },

    isRoleVisible ({ meta }) {
      return !(meta.context && meta.context.resourceTypes)
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

    setDefaultValues () {
      this.checkboxLabel = {}
      this.resolvedRoles = {}
    },
  },
}
</script>

<style>
.w-fit {
  width: fit-content;
}
</style>

<style lang="scss" scoped>
.list-background {
  background-color: var(--body-bg);
}
</style>
