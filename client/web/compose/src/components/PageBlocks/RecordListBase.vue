<template>
  <wrap
    v-if="recordListModule"
    v-bind="$props"
    :scrollable-body="false"
    v-on="$listeners"
    @refreshBlock="refresh(true, false)"
  >
    <template
      v-if="isFederated"
      #title-badge
    >
      <b-badge
        variant="primary"
        class="d-inline-block mb-0 ml-2"
      >
        {{ $t('recordList.federated') }}
      </b-badge>
    </template>

    <template #toolbar>
      <b-container
        ref="toolbar"
        fluid
        class="d-flex flex-column gap-2 p-3 d-print-none"
      >
        <b-row
          no-gutters
          class="d-flex align-items-center justify-content-between gap-1"
        >
          <div class="d-flex align-items-center flex-grow-1 flex-wrap flex-fill-child gap-1">
            <template v-if="recordListModule.canCreateRecord">
              <template v-if="inlineEditing">
                <b-button
                  v-if="!options.hideAddButton"
                  data-test-id="button-add-record"
                  variant="primary"
                  size="lg"
                  @click="addInlineRecord()"
                >
                  + {{ $t('recordList.addRecord') }}
                </b-button>
              </template>

              <template v-else-if="!inlineEditing && (recordPageID || options.allRecords)">
                <b-button
                  v-if="!options.hideAddButton"
                  data-test-id="button-add-record"
                  variant="primary"
                  size="lg"
                  @click="handleAddRecord()"
                >
                  + {{ $t('recordList.addRecord') }}
                </b-button>

                <importer-modal
                  v-if="!options.hideImportButton"
                  :module="recordListModule"
                  :namespace="namespace"
                  @importSuccessful="onImportSuccessful"
                />
              </template>
            </template>

            <exporter-modal
              v-if="options.allowExport && !inlineEditing"
              :module="recordListModule"
              :filter="filter.query"
              :selection="selected"
              :selected-all-records="selectedAllRecords"
              :processing="processing"
              :preselected-fields="fields.map(({ moduleField }) => moduleField)"
              @export="onExport"
            />

            <b-dropdown
              v-if="filterPresets.length"
              ref="filterPresets"
              variant="light"
              size="lg"
              menu-class="shadow-sm"
              boundary="viewport"
              :text="$t('recordList.filter.filters.label')"
            >
              <li
                v-for="(f, idx) in filterPresets"
                :key="idx"
                class="d-flex align-items-center justify-content-between"
              >
                <button
                  class="dropdown-item"
                  :disabled="activeFilters.includes(f.name)"
                  @click="updateFilter(f.filter, f.name)"
                >
                  {{ f.name }}
                </button>

                <c-input-confirm
                  v-if="!f.roles"
                  show-icon
                  class="mr-1"
                  @confirmed="removeStorageRecordListFilterPreset(f.name)"
                />
              </li>
            </b-dropdown>

            <column-picker
              v-if="!options.hideConfigureFieldsButton"
              :module="recordListModule"
              :fields="fields"
              @updateFields="onUpdateFields"
            />
          </div>

          <div
            v-if="!options.hideSearch"
            class="flex-fill"
          >
            <c-input-search
              v-model.trim="query"
              :placeholder="$t('general.label.search')"
              :debounce="500"
            />
          </div>
        </b-row>

        <div
          v-if="activeFilters.length || drillDownFilter || options.showDeletedRecordsOption"
          class="d-flex"
        >
          <div
            v-if="activeFilters.length"
            class="d-flex align-items-center flex-wrap"
          >
            {{ $t('recordList.filter.filters.active') }}
            <b-form-tags
              size="lg"
              class="d-flex align-items-center border-0 p-0 bg-transparent"
              style="width: fit-content;"
            >
              <b-form-tag
                v-for="(title, i) in activeFilters"
                :key="i"
                :title="title"
                variant="light"
                pill
                class="align-items-center ml-2"
                @remove="removeFilter(i)"
              />
            </b-form-tags>
          </div>

          <b-button
            v-if="options.showDeletedRecordsOption"
            variant="outline-extra-light"
            size="sm"
            class="text-primary border-0 text-nowrap ml-auto"
            @click="handleShowDeleted()"
          >
            {{ showingDeletedRecords ? $t('recordList.showRecords.existing') : $t('recordList.showRecords.deleted') }}
          </b-button>
        </div>

        <div
          v-if="options.selectable && selected.length"
          class="d-flex align-items-center flex-wrap align-items-center"
        >
          <div class="mr-1">
            {{ selectedRecordsDisplayText }}
          </div>

          <b-button
            v-if="!inlineEditing"
            size="sm"
            variant="outline-extra-light"
            class="text-primary border-0"
            @click="selectAllRecords()"
          >
            {{ selectedAllRecords ? $t('recordList.unselectAllRecords') : $t('recordList.selectAllRecords') }}
          </b-button>

          <div class="d-flex align-items-center ml-auto gap-1">
            <automation-buttons
              class="d-inline m-0 mr-2"
              :buttons="options.selectionButtons"
              :module="recordListModule"
              :extra-event-args="{ selected, filter}"
              v-bind="$props"
              @refresh="refresh()"
            />

            <bulk-edit-modal
              v-show="options.bulkRecordEditEnabled && canUpdateSelectedRecords && !showingDeletedRecords"
              :module="recordListModule"
              :namespace="namespace"
              :query="bulkQuery"
              @save="onBulkUpdate()"
            />

            <template v-if="canDeleteSelectedRecords && !areAllRowsDeleted">
              <c-input-confirm
                show-icon
                :tooltip="$t('recordList.tooltip.deleteSelected')"
                @confirmed="handleDeleteSelectedRecords()"
              />
            </template>

            <template v-if="canRestoreSelectedRecords && areAllRowsDeleted">
              <c-input-confirm
                show-icon
                :icon="['fas', 'trash-restore']"
                :tooltip="$t('recordList.tooltip.restoreSelected')"
                variant="outline-warning"
                @confirmed="handleRestoreSelectedRecords()"
              />
            </template>
          </div>
        </div>
      </b-container>
    </template>

    <template #default>
      <div
        class="d-flex position-relative h-100"
        :class="{ 'overflow-hidden': !items.length || processing }"
      >
        <b-table-simple
          data-test-id="table-record-list"
          hover
          responsive
          sticky-header
          class="record-list-table mh-100 h-100 mb-0"
        >
          <b-thead>
            <b-tr :variant="showingDeletedRecords ? 'warning' : ''">
              <b-th v-if="options.draggable && inlineEditing" />

              <b-th
                v-if="options.selectable"
                style="width: 0%;"
                class="d-print-none"
              >
                <b-checkbox
                  :disabled="disableSelectAll"
                  :checked="areAllRowsSelected && !disableSelectAll"
                  class="ml-1"
                  @change="handleSelectAllOnPage({ isChecked: $event })"
                />
              </b-th>

              <b-th v-if="isFederated" />

              <b-th
                v-for="(field, fieldIndex) in fields"
                :key="field.key"
                sticky-column
                :colspan="fieldIndex === (fields.length - 1) ? 2 : 1"
                :style="{
                  'padding-right': fieldIndex === (fields.length - 1) ? '15px' : '',
                }"
              >
                <div class="d-flex align-items-center">
                  <div
                    :class="{ required: field.required }"
                    class="d-flex align-self-center text-nowrap"
                  >
                    {{ field.label }}
                  </div>

                  <b-button
                    v-if="field.sortable"
                    v-b-tooltip.noninteractive.hover="{ title: $t('recordList.sort.tooltip'), container: '#body' }"
                    variant="outline-extra-light"
                    class="d-flex align-items-center text-secondary d-print-none border-0 px-1 ml-1"
                    @click="handleSort(field)"
                  >
                    <font-awesome-layers
                      class="d-print-none"
                    >
                      <font-awesome-icon
                        :icon="['fas', 'angle-up']"
                        class="mb-1"
                        :class="{ 'text-primary': isSortedBy(field, 'ASC') }"
                      />
                      <font-awesome-icon
                        :icon="['fas', 'angle-down']"
                        class="mt-1"
                        :class="{ 'text-primary': isSortedBy(field, 'DESC') }"
                      />
                    </font-awesome-layers>
                  </b-button>

                  <record-list-filter
                    v-if="!options.hideFiltering && field.filterable"
                    :target="uniqueID"
                    :selected-field="field.moduleField"
                    :namespace="namespace"
                    :module="recordListModule"
                    variant="outline-extra-light"
                    :record-list-filter="recordListFilter"
                    :allow-filter-preset-save="options.customFilterPresets"
                    class="d-print-none ml-1"
                    @filter="onFilter"
                    @filter-preset="onSaveFilterPreset"
                  />
                </div>
              </b-th>
            </b-tr>
          </b-thead>

          <draggable
            v-if="items.length && !processing && !resizing"
            v-model="items"
            :disabled="!inlineEditing || !options.draggable"
            group="items"
            tag="b-tbody"
            handle=".handle"
          >
            <b-tr
              v-for="(item, index) in items"
              :key="`${index}${item.r.recordID}`"
              :class="{ 'pointer': !(options.editable && editing), }"
              :variant="inlineEditing && item.r.deletedAt ? 'warning' : ''"
              @click="handleRowClick(item)"
            >
              <b-td
                v-if="options.draggable && inlineEditing"
                class="pr-0"
                @click.stop
              >
                <font-awesome-icon
                  :icon="['fas', 'bars']"
                  class="handle text-secondary my-1"
                />
              </b-td>

              <b-td
                v-if="options.selectable"
                class="pr-0 d-print-none"
                @click.stop
              >
                <b-form-checkbox
                  class="ml-1"
                  :checked="selected.includes(item.id)"
                  @change="onSelectRow($event, item)"
                />
              </b-td>

              <b-td
                v-if="isFederated"
                class="align-middle pl-0"
              >
                <b-badge
                  v-if="Object.keys(item.r.labels || {}).includes('federation')"
                  variant="primary"
                  class="align-text-top"
                >
                  F
                </b-badge>
              </b-td>

              <b-td
                v-for="field in fields"
                :key="field.key"
              >
                <field-editor
                  v-if="field.moduleField.canUpdateRecordValue && field.editable"
                  :field="field.moduleField"
                  value-only
                  :record="item.r"
                  :module="module"
                  :namespace="namespace"
                  :errors="recordErrors(item, field)"
                  class="mb-0"
                  style="min-width: 250px;"
                  @click.stop
                />

                <div
                  v-else-if="field.moduleField.canReadRecordValue && !field.edit"
                  class="d-flex mb-0"
                  :class="{
                    'field-adjust-offset': inlineEditing,
                  }"
                >
                  <field-viewer
                    :field="field.moduleField"
                    value-only
                    :record="item.r"
                    :module="module"
                    :namespace="namespace"
                    :extra-options="options"
                  />
                  <div
                    v-if="options.inlineRecordEditEnabled && field.canEdit && !showingDeletedRecords"
                    class="inline-actions"
                  >
                    <b-button
                      v-b-tooltip.noninteractive.hover="{ title: $t('recordList.inlineEdit.button.title', { label: field.label }), container: '#body' }"
                      variant="outline-extra-light"
                      size="sm"
                      class="text-secondary border-0 ml-1"
                      @click.stop="editInlineField(item.r, field.key)"
                    >
                      <font-awesome-icon
                        :icon="['fas', 'pen']"
                      />
                    </b-button>
                  </div>
                </div>

                <i
                  v-else
                  class="text-primary"
                >
                  {{ $t('field.noPermission') }}
                </i>
              </b-td>

              <b-td
                class="actions px-2"
                @click.stop
              >
                <b-dropdown
                  v-if="areActionsVisible(item.r)"
                  boundary="viewport"
                  variant="outline-extra-light"
                  toggle-class="d-flex align-items-center justify-content-center text-primary border-0 py-2"
                  no-caret
                  dropleft
                  menu-class="m-0"
                >
                  <template #button-content>
                    <font-awesome-icon
                      :icon="['fas', 'ellipsis-v']"
                    />
                  </template>

                  <template v-if="inlineEditing">
                    <b-dropdown-item
                      v-if="isCloneRecordActionVisible"
                      @click="handleCloneInline(item.r)"
                    >
                      <font-awesome-icon
                        :icon="['far', 'clone']"
                        class="text-primary"
                      />
                      {{ $t('recordList.record.tooltip.clone') }}
                    </b-dropdown-item>

                    <b-dropdown-item
                      v-if="isInlineRestoreActionVisible(item.r)"
                      @click.prevent="handleRestoreInline(item, index)"
                    >
                      <font-awesome-icon
                        :icon="['fas', 'trash-restore']"
                        class="text-warning"
                      />
                      {{ $t('recordList.record.tooltip.restore') }}
                    </b-dropdown-item>

                    <!-- The user should be able to delete the record if it's not yet saved -->
                    <b-dropdown-item
                      v-else-if="isInlineDeleteActionVisible(item.r)"
                      @click.prevent="handleDeleteInline(item, index)"
                    >
                      <font-awesome-icon
                        :icon="['far', 'trash-alt']"
                        class="text-danger"
                      />
                      {{ $t('recordList.record.tooltip.delete') }}
                    </b-dropdown-item>
                  </template>

                  <template
                    v-else
                  >
                    <b-dropdown-item
                      v-if="isViewRecordActionVisible(item.r)"
                      @click="handleViewRecordAction(item.r.recordID)"
                    >
                      <font-awesome-icon
                        :icon="['far', 'file-alt']"
                        class="text-primary"
                      />
                      {{ $t('recordList.record.tooltip.view') }}
                    </b-dropdown-item>

                    <b-dropdown-item
                      v-if="isEditRecordActionVisible(item.r)"
                      @click="handleEditRecordAction(item.r.recordID)"
                    >
                      <font-awesome-icon
                        :icon="['far', 'edit']"
                        class="text-primary"
                      />
                      {{ $t('recordList.record.tooltip.edit') }}
                    </b-dropdown-item>

                    <b-dropdown-item
                      v-if="isCloneRecordActionVisible"
                      @click="handleCloneRecordAction(item.r.recordID, item.r.values)"
                    >
                      <font-awesome-icon
                        :icon="['far', 'clone']"
                        class="text-primary"
                      />
                      {{ $t('recordList.record.tooltip.clone') }}
                    </b-dropdown-item>

                    <b-dropdown-item
                      v-if="isReminderActionVisible"
                      @click.prevent="createReminder(item.r)"
                    >
                      <font-awesome-icon
                        :icon="['far', 'bell']"
                        class="text-primary"
                      />
                      {{ $t('recordList.record.tooltip.reminder') }}
                    </b-dropdown-item>

                    <b-dropdown-item
                      v-if="isRecordPermissionButtonVisible(item.r)"
                      link-class="p-0"
                      variant="light"
                    >
                      <c-permissions-button
                        :resource="`corteza::compose:record/${item.r.namespaceID}/${item.r.moduleID}/${item.r.recordID}`"
                        :target="item.r.recordID"
                        :title="item.r.recordID"
                        :button-label="$t('recordList.record.tooltip.permissions')"
                        button-variant="link dropdown-item text-decoration-none text-dark rounded-0"
                      />
                    </b-dropdown-item>

                    <c-input-confirm
                      v-if="isDeleteActionVisible(item.r)"
                      :text="$t('recordList.record.tooltip.delete')"
                      show-icon
                      borderless
                      variant="link"
                      size="md"
                      button-class="dropdown-item text-decoration-none text-dark rounded-0"
                      icon-class="text-danger"
                      class="w-100"
                      @confirmed="handleDeleteSelectedRecords(item.r.recordID)"
                    />

                    <c-input-confirm
                      v-else-if="isRestoreActionVisible(item.r)"
                      :text="$t('recordList.record.tooltip.restore')"
                      :icon="['fas', 'trash-restore']"
                      show-icon
                      borderless
                      variant="link"
                      size="md"
                      button-class="dropdown-item text-decoration-none text-dark rounded-0"
                      icon-class="text-warning"
                      class="w-100"
                      @confirmed="handleRestoreSelectedRecords(item.r.recordID)"
                    />
                  </template>
                </b-dropdown>
              </b-td>
            </b-tr>
          </draggable>

          <div
            v-else
            class="position-absolute text-center mt-5 d-print-none"
            style="left: 0; right: 0; bottom: calc(50% - 33px);"
          >
            <b-spinner
              v-if="processing"
            />

            <p
              v-else-if="!items.length"
              class="mb-0 mx-2"
            >
              {{ $t('recordList.noRecords') }}
            </p>
          </div>
        </b-table-simple>
      </div>

      <!-- Modal for inline editing -->
      <bulk-edit-modal
        v-if="options.inlineRecordEditEnabled"
        :namespace="namespace"
        :module="recordListModule"
        :selected-fields="inlineEdit.fields"
        :initial-record="inlineEdit.record"
        :query="inlineEdit.query"
        :modal-title="$t('recordList.inlineEdit.modal.title')"
        open-on-select
        @save="onInlineEdit()"
        @close="onInlineEditClose()"
      />

      <!-- Modal for naming custom filter -->
      <custom-filter-preset
        :visible="showCustomPresetFilterModal"
        @save="setStorageRecordListFilterPreset"
        @close="showCustomPresetFilterModal = false"
      />
    </template>

    <template
      v-if="showFooter"
      #footer
    >
      <div class="record-list-footer d-flex align-items-center flex-wrap justify-content-between px-3 py-2 gap-1">
        <div class="d-flex align-items-center flex-wrap gap-3 gap-col-3">
          <div
            v-if="options.showTotalCount"
            class="text-nowrap text-truncate"
          >
            <span
              v-if="pagination.count > recordsPerPage"
              data-test-id="pagination-range"
            >
              {{ $t('recordList.pagination.showing', getPagination) }}
            </span>

            <span
              v-else
              data-test-id="pagination-single-number"
            >
              {{ $t('recordList.pagination.single', getPagination) }}
            </span>
          </div>

          <div
            v-if="options.showRecordPerPageOption"
            class="d-flex align-items-center gap-1 text-nowrap"
          >
            <span>
              {{ $t('recordList.pagination.recordsPerPage') }}
            </span>

            <b-form-select
              v-model="recordsPerPage"
              :options="perPageOptions"
              size="sm"
              @change="handlePerPageChange"
            />
          </div>
        </div>

        <div
          v-if="showPageNavigation"
          class="d-flex align-items-center justify-content-end "
        >
          <b-pagination
            v-if="options.fullPageNavigation"
            data-test-id="pagination"
            align="right"
            aria-controls="record-list"
            class="m-0 d-print-none"
            pills
            :disabled="processing"
            :value="getPagination.page"
            :per-page="getPagination.perPage"
            :total-rows="getPagination.count"
            @change="goToPage"
          >
            <template #first-text>
              <font-awesome-icon :icon="['fas', 'angle-double-left']" />
            </template>

            <template #prev-text>
              <font-awesome-icon :icon="['fas', 'angle-left']" />
            </template>

            <template #next-text>
              <font-awesome-icon :icon="['fas', 'angle-right']" />
            </template>

            <template #last-text>
              <font-awesome-icon :icon="['fas', 'angle-double-right']" />
            </template>

            <template #elipsis-text>
              <font-awesome-icon :icon="['fas', 'ellipsis-h']" />
            </template>
          </b-pagination>

          <b-button-group
            v-else
            class="gap-1"
          >
            <b-button
              :disabled="!hasPrevPage || processing"
              data-test-id="first-page"
              variant="outline-extra-light"
              class="d-flex align-items-center justify-content-center text-dark border-0 p-1"
              @click="goToPage()"
            >
              <font-awesome-icon :icon="['fas', 'angle-double-left']" />
            </b-button>

            <b-button
              :disabled="!hasPrevPage || processing"
              data-test-id="previous-page"
              variant="outline-extra-light"
              class="d-flex align-items-center justify-content-center text-dark border-0 p-1"
              @click="goToPage('prevPage')"
            >
              <font-awesome-icon
                :icon="['fas', 'angle-left']"
                class="mr-1"
              />
              {{ $t('recordList.pagination.prev') }}
            </b-button>

            <b-button
              :disabled="!hasNextPage || processing"
              data-test-id="next-page"
              variant="outline-extra-light"
              class="d-flex align-items-center justify-content-center text-dark border-0 p-1"
              @click="goToPage('nextPage')"
            >
              {{ $t('recordList.pagination.next') }}
              <font-awesome-icon
                :icon="['fas', 'angle-right']"
                class="ml-1"
              />
            </b-button>
          </b-button-group>
        </div>
      </div>
    </template>
  </wrap>
</template>
<script>
import axios from 'axios'
import { mapGetters, mapActions } from 'vuex'
import base from './base'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import FieldEditor from 'corteza-webapp-compose/src/components/ModuleFields/Editor'
import ExporterModal from 'corteza-webapp-compose/src/components/Public/Record/Exporter'
import ImporterModal from 'corteza-webapp-compose/src/components/Public/Record/Importer'
import AutomationButtons from './Shared/AutomationButtons'
import { compose, validator, NoID } from '@cortezaproject/corteza-js'
import users from 'corteza-webapp-compose/src/mixins/users'
import records from 'corteza-webapp-compose/src/mixins/records'
import { evaluatePrefilter, queryToFilter, isFieldInFilter } from 'corteza-webapp-compose/src/lib/record-filter'
import { getItem, setItem, removeItem } from 'corteza-webapp-compose/src/lib/local-storage'
import { components, url } from '@cortezaproject/corteza-vue'
import draggable from 'vuedraggable'
import RecordListFilter from 'corteza-webapp-compose/src/components/Common/RecordListFilter'
import ColumnPicker from 'corteza-webapp-compose/src/components/Admin/Module/Records/ColumnPicker'
import BulkEditModal from 'corteza-webapp-compose/src/components/Public/Record/BulkEdit'
import CustomFilterPreset from 'corteza-webapp-compose/src/components/Public/Record/CustomFilterPreset'

const { CInputSearch } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    ExporterModal,
    ImporterModal,
    AutomationButtons,
    FieldViewer,
    FieldEditor,
    draggable,
    RecordListFilter,
    ColumnPicker,
    CInputSearch,
    BulkEditModal,
    CustomFilterPreset,
  },

  extends: base,

  mixins: [
    users,
    records,
  ],

  props: {
    errors: {
      type: validator.Validated,
      required: false,
      default: () => new validator.Validated(),
    },
  },

  data () {
    return {
      uniqueID: undefined,

      processing: false,
      // prefilter from block config
      prefilter: undefined,
      recordListFilter: [],
      drillDownFilter: undefined,

      // raw query string used to build final filter
      query: null,

      // used to construct request parameters
      // AND to store response params
      filter: {
        query: '',
        sort: '',
        limit: 10,
        pageCursor: '',
        prevPage: '',
        nextPage: '',
      },

      pagination: {
        pages: [],
        page: 1,
        count: 0,
      },

      selected: [],
      inlineEdit: {
        fields: [],
        recordIDs: [],
        initialRecord: {},
      },

      sortBy: undefined,
      sortDirecton: undefined,

      // This counter helps us generate unique ID's for the lifetime of this
      // component
      ctr: 0,
      items: [],
      showingDeletedRecords: false,
      activeFilters: [],
      customPresetFilters: [],
      currentCustomPresetFilter: undefined,
      showCustomPresetFilterModal: false,
      selectedAllRecords: false,

      abortableRequests: [],
      recordsPerPage: undefined,

      customConfiguredFields: [],
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      pages: 'page/set',
    }),

    isFederated () {
      return Object.keys(this.recordListModule.labels || {}).includes('federation')
    },

    showFooter () {
      return this.showPageNavigation || this.options.showTotalCount || this.options.showRecordPerPageOption
    },

    perPageOptions () {
      const defaultText = this.options.perPage === 0 ? this.$t('general:label.all') : this.options.perPage.toString()
      return [
        { text: defaultText, value: this.options.perPage },
        { text: '25', value: 25 },
        { text: '50', value: 50 },
        { text: '100', value: 100 },
      ].filter((v, i) => i === 0 || v.value !== this.options.perPage).sort((a, b) => {
        if (a.value === 0) return 1
        if (b.value === 0) return -1
        return a.value - b.value
      })
    },

    getPagination () {
      const { page = 1, count = 0 } = this.pagination
      const perPage = this.recordsPerPage

      return {
        from: ((page - 1) * perPage) + 1,
        to: perPage > 0 ? Math.min((page * perPage), count) : count,
        page,
        perPage,
        count,
      }
    },

    hasPrevPage () {
      return this.filter.prevPage
    },

    hasNextPage () {
      return this.filter.nextPage
    },

    editing () {
      return this.mode === 'editor'
    },

    showPageNavigation () {
      return !this.options.hidePaging
    },

    showPerPageSelector () {
      return this.options.showRecordPerPageOption
    },

    disableSelectAll () {
      if (this.options.hidePaging) {
        return !this.items.length
      }
      return this.items.length === 0
    },

    inlineEditing () {
      return !!this.options.editable && !!this.editing
    },

    /**
     * Check if all rows are selected
     */
    areAllRowsSelected () {
      return this.selected.length === this.items.length
    },

    areAllRowsDeleted () {
      const selItems = this.items.filter(({ id }) => this.selected.includes(id))
      return !!this.selected.length && !selItems.find(({ r }) => !r.deletedAt)
    },

    // Returns module, configured for this record list
    recordListModule () {
      if (this.options.moduleID) {
        return this.getModuleByID(this.options.moduleID)
      } else {
        return undefined
      }
    },

    // Tries to determine ID of the page we're supposed to redirect
    recordPageID () {
      // Relying on pages having unique moduleID,
      const { moduleID } = this.recordListModule || {}
      if (!moduleID) {
        return undefined
      }

      const { pageID } = this.pages.find(p => p.moduleID === moduleID) || {}
      if (!pageID) {
        return undefined
      }

      return pageID
    },

    fields () {
      let fields = []

      const editable = (!this.options.editable || !this.editing)
        ? []
        : this.options.editFields.map(({ name }) => name)

      if (!this.options.hideConfigureFieldsButton && this.customConfiguredFields.length > 0) {
        fields = this.recordListModule.filterFields(this.customConfiguredFields)
      } else if (this.options.fields.length > 0) {
        fields = this.recordListModule.filterFields(this.options.fields)
      } else {
        // Record list block does not have any configured fields
        // Use first five fields from the module.
        fields = [...this.recordListModule.fields.slice(0, 5), ...this.recordListModule.systemFields()]
      }

      const configured = fields.map(mf => ({
        key: mf.name,
        label: mf.isSystem ? this.$t(`field:system.${mf.name}`) : mf.label || mf.name,
        moduleField: mf,
        sortable: !this.options.hideSorting && !(this.options.editable && this.editing) && !mf.isMulti && mf.isSortable,
        filterable: mf.isFilterable,
        tdClass: 'record-value',
        editable: !!editable.find(f => mf.name === f),
        canEdit: this.isFieldEditable(mf),
        required: this.inlineEditing && mf.isRequired,
      }))

      const pre = []
      const post = []

      return [
        ...pre,
        ...configured,
        ...post,
      ]
    },

    canDeleteSelectedRecords () {
      return this.items.filter(({ id, r }) => this.selected.includes(id) && r.canDeleteRecord).length
    },

    canUpdateSelectedRecords () {
      return this.items.filter(({ id, r }) => this.selected.includes(id) && r.canUpdateRecord).length
    },

    canRestoreSelectedRecords () {
      return this.items.filter(({ id, r }) => this.selected.includes(id) && r.canUndeleteRecord).length
    },

    isCloneRecordActionVisible () {
      return !this.options.hideRecordCloneButton && this.recordListModule.canCreateRecord && (this.options.rowCreateUrl || this.recordPageID || this.inlineEditing)
    },

    isReminderActionVisible () {
      return !this.options.hideRecordReminderButton
    },

    filterPresets () {
      return [
        ...this.options.filterPresets.filter(({ name, roles }) => name && this.isUserRoleMember(roles)),
        ...this.customPresetFilters,
      ]
    },

    authUserRoles () {
      return this.$auth.user.roles
    },

    selectedRecordsDisplayText () {
      const count = this.selectedAllRecords ? (this.options.showTotalCount ? this.pagination.count : undefined) : this.selected.length
      const total = this.items.length
      const key = this.selectedAllRecords ? 'selectedFromAllPages' : 'selected'

      return this.$t(`recordList.${key}`, { count, total })
    },

    bulkQuery () {
      if (this.selectedAllRecords) {
        return this.filter.query
      }

      return this.selected.map(r => `recordID='${r}'`).join(' OR ')
    },

    isOnRecordPage () {
      return this.page && this.page.moduleID !== NoID
    },
  },

  watch: {
    query: {
      handler () {
        this.refresh(true)
      },
    },

    options: {
      deep: true,
      handler () {
        this.prepRecordList()
        this.refresh(true)
      },
    },

    'record.updatedAt': {
      immediate: true,
      handler () {
        this.createEvents()
        this.getStorageRecordListFilter()
        this.getStorageRecordListFilterPreset()
        this.getStorageRecordListConfiguredFields()
        this.prepRecordList()
        this.refresh(true)
      },
    },
  },

  beforeDestroy () {
    this.abortRequests()
    this.destroyEvents()
    this.setDefaultValues()
  },

  created () {
    if (!this.inlineEditing) {
      this.refreshBlock(this.refresh, false, true)
    }
  },

  methods: {
    ...mapActions({
      loadPaginationRecords: 'ui/loadPaginationRecords',
    }),

    createEvents () {
      const { pageID = NoID } = this.page
      const { recordID = NoID } = this.record || {}

      // Set uniqueID so that events dont mix
      if (this.uniqueID) {
        this.destroyEvents()
      }

      this.uniqueID = [pageID, recordID, this.block.blockID, this.magnified].map(v => v || NoID).join('-')
      this.$root.$on(`record-line:collect:${this.uniqueID}`, this.resolveRecords)
      this.$root.$on(`page-block:validate:${this.uniqueID}`, this.validatePageBlock)
      this.$root.$on(`drill-down-recordList:${this.uniqueID}`, this.setDrillDownFilter)
      this.$root.$on(`refetch-non-record-blocks:${pageID}`, this.refreshAndResetPagination)

      this.$root.$on('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$on('record-field-change', this.refetchOnPrefilterValueChange)
    },

    refetchOnPrefilterValueChange ({ fieldName }) {
      const { prefilter } = this.options

      if (isFieldInFilter(fieldName, prefilter)) {
        this.prepRecordList()
        this.refresh()
      }
    },

    refreshOnRelatedRecordsUpdate ({ moduleID, notPageID }) {
      if (this.page.pageID === notPageID) {
        return
      }

      if (this.recordListModule.moduleID === moduleID) {
        this.refresh(true)
      } else {
        const recordFields = this.fields.filter((f) => f.moduleField.kind === 'Record')
        const hasMatchingModule = recordFields.some((r) => {
          return r.moduleField.options.moduleID === moduleID
        })

        if (hasMatchingModule) {
          this.refresh(false)
        }
      }
    },

    onFilter (filter = []) {
      if (filter.length) {
        filter.forEach(f => {
          f.name = this.$t('recordList.customFilter')
        })

        this.activeFilters = [this.$t('recordList.customFilter')]
      } else {
        this.activeFilters = []
      }

      this.recordListFilter = filter
      this.setStorageRecordListFilter()
      this.refresh(true)
    },

    handlePerPageChange () {
      this.filter.limit = this.recordsPerPage
      this.refresh(true)
    },

    onSaveFilterPreset (filter = []) {
      this.currentCustomPresetFilter = {
        filter,
      }

      this.showCustomPresetFilterModal = true
      this.refresh(true)
    },

    onUpdateFields (fields = []) {
      this.options.fields = [...fields]
      this.customConfiguredFields = fields.map((f) => f.isSystem ? f.name : f.fieldID)
      this.setStorageRecordListConfiguredFields()

      this.$emit('save-fields', this.options.fields)
    },

    setStorageRecordListConfiguredFields () {
      try {
        // Get record list configured fields from localStorage
        setItem(`record-list-configured-columns-${this.uniqueID}`, this.customConfiguredFields)
      } catch (e) {
        console.warn(this.$t('notification:record-list.corrupted-configured-fields'))
      }
    },

    getStorageRecordListConfiguredFields () {
      try {
        // Get record list configured fields from localStorage
        this.customConfiguredFields = getItem(`record-list-configured-columns-${this.uniqueID}`)
      } catch (e) {
        // Land here if the configured fields is corrupted
        console.warn(this.$t('notification:record-list.corrupted-configured-fields'))
        // Remove filter from the local storage
        removeItem(`record-list-configured-columns-${this.uniqueID}`)
      }
    },

    onSelectRow (selected, item) {
      if (selected) {
        if (this.selected.includes(item.id)) {
          return
        }

        this.selected.push(item.id)
      } else {
        const i = this.selected.indexOf(item.id)
        if (i < 0) {
          return
        }
        this.selected.splice(i, 1)

        this.selectedAllRecords = false
      }
    },

    isSortedBy ({ key }, dir) {
      const { sort = '' } = this.filter

      const sortedFields = (sort.includes(',') ? sort.split(',') : [sort])

      return sortedFields.map(v => v.trim()).some(value => {
        let valueDir = 'ASC'

        if (value.includes(' ')) {
          value = value.split(' ')[0]
          valueDir = 'DESC'
        }

        return valueDir === dir && value === key
      })
    },

    handleShowDeleted () {
      this.showingDeletedRecords = !this.showingDeletedRecords
      this.selectedAllRecords = false
      this.refresh(true)
    },

    // Grabs errors specific to this record item
    recordErrors (item, field) {
      if (field) {
        return this.errors.filterByMeta('id', item.id)
          .filterByMeta('field', field.key)
      }
      return this.errors.filterByMeta('id', item.id)
    },

    wrapRecord (r, id) {
      if (r.id) {
        id = r.id
        r = r.r
      }

      return {
        r,
        id: id || (r.recordID !== NoID ? r.recordID : `${this.uniqueID}:${this.ctr++}`),
      }
    },

    addInlineRecord () {
      const r = new compose.Record(this.recordListModule, {})

      // Set record values that should be prefilled
      if ((this.record || {}).recordID && this.options.refField) {
        r.values[this.options.refField] = (this.record || {}).recordID
      }

      this.items.unshift(this.wrapRecord(r))
    },

    /**
     * Helper method to fetch all records available to this record list
     * at the given point in time.
     *
     * It:
     *    * assures that local records have a sequential indexing
     *    * appends additional meta fields
     *    * resolves payload editing
     */
    resolveRecords (resolve) {
      this.ctr = 0
      this.items = this.items.map(this.wrapRecord)

      resolve({
        items: this.items,
        module: this.recordListModule,
        refField: this.options.refField,
        positionField: this.options.positionField,
        idPrefix: this.uniqueID,
      })
    },

    validatePageBlock (resolve) {
      // For now, only record lines should be validated
      if (!this.options.editable) {
        resolve({ valid: true })
      }

      // Find all required fields
      const req = new Set(this.recordListModule.fields.filter(({ isRequired = false }) => isRequired).map(({ name }) => name))

      // Check if all required fields are there
      for (const f of this.options.editFields) {
        req.delete(f.name)
      }

      // If required fields are satisfied, then the validation passes
      resolve({ valid: !req.size })
      req.clear()
    },

    handleDeleteInline (item, i) {
      if (item.r.recordID !== NoID) {
        const r = new compose.Record(this.recordListModule, { ...item.r, deletedAt: new Date() })
        this.items.splice(i, 1, this.wrapRecord(r, item.id))
      } else {
        this.items.splice(i, 1)
      }
    },

    handleRestoreInline (item, i) {
      const r = new compose.Record(this.recordListModule, { ...item.r, deletedAt: undefined })
      this.items.splice(i, 1, this.wrapRecord(r, item.id))
    },

    handleCloneInline (r) {
      r = new compose.Record(r.module, { ...r.values })
      this.items.splice(0, 0, this.wrapRecord(r))
    },

    // Sanitizes record list config and
    // prepares prefilter
    prepRecordList () {
      this.recordsPerPage = this.options.perPage

      // Legacy support for linkToParent
      if (this.isOnRecordPage && this.options.linkToParent) {
        if (this.options.refField) {
          this.options.linkToParent = false
        } else {
          this.options.refField = (this.recordListModule.fields.find(f => f.kind === 'Record' && f.options.moduleID === this.page.moduleID) || {}).name
        }
      }

      const { moduleID, presort, prefilter, editable, refField, positionField } = this.options

      // Validate props
      if (!moduleID) {
        throw Error(this.$t('record.moduleOrPageNotSet'))
      }

      // If there is no current record and we are using recordID/ownerID variable in (pre)filter
      // we should disable the block
      /* eslint-disable no-template-curly-in-string */
      if (!this.record) {
        if ((prefilter || '').includes('${record')) {
          throw Error(this.$t('record.invalidRecordVar'))
        }

        if ((prefilter || '').includes('${ownerID}')) {
          throw Error(this.$t('record.invalidOwnerVar'))
        }
      }

      // Sorting
      let sort = presort || 'createdAt DESC'

      if (editable && positionField) {
        sort = `${positionField}`
      }

      const filter = []

      // Initial filter
      if (prefilter) {
        const pf = evaluatePrefilter(prefilter, {
          record: this.record,
          user: this.$auth.user || {},
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
        filter.push(`(${pf})`)
      }

      if (refField && (this.record || {}).recordID) {
        filter.push(`(${refField} = ${this.record.recordID})`)
      }

      this.prefilter = filter.join(' AND ')

      this.filter = {
        limit: this.recordsPerPage,
        sort,
      }
    },

    createReminder (record) {
      // Determine initial reminder title
      const { recordID, values = {} } = record
      const { name, isMulti } = (this.options.fields || []).find(({ name }) => !!values[name]) || {}
      const title = isMulti ? values[name].join(', ') : values[name]

      const resource = `compose:record:${recordID}`
      const payload = {
        title,
        link: {
          name: 'page.record',
          label: 'Record page',
          params: {
            slug: this.namespace.slug || this.namespace.namespaceID,
            pageID: this.recordPageID,
            recordID: recordID,
          },
        },
      }

      this.$root.$emit('reminder.create', { payload, resource })
      this.$root.$emit('rightPanel.toggle', true)
    },

    onExport (e) {
      this.processing = true

      const { namespaceID, moduleID } = this.filter || {}
      const { filter, filterRaw, timezone } = e
      e = {
        ...e,
        namespaceID,
        moduleID,
        filename: `${this.namespace.slug || namespaceID} - ${this.recordListModule.name}`,
      }

      if (filterRaw.rangeType === 'range') {
        e.filename += ` - ${filterRaw.date.start} - ${filterRaw.date.end}`
      } else {
        e.filename += ` - ${filterRaw.rangeType}`
      }

      if (timezone) {
        e.filename += ` - ${timezone.label}`
      }

      // Make sure the generated filename won't break the URL
      e.filename = encodeURIComponent(e.filename.replace(/\./g, '-'))

      const exportUrl = url.Make({
        url: `${this.$ComposeAPI.baseURL}${this.$ComposeAPI.recordExportEndpoint(e)}`,
        query: {
          fields: e.fields,
          // url.Make already URL encodes the the values, so the filter shouldn't be encoded
          filter: this.selectedAllRecords ? this.bulkQuery : filter,
          jwt: this.$auth.accessToken,
          timezone: timezone ? timezone.tzCode : undefined,
        },
      })

      window.open(exportUrl)
      this.processing = false
    },

    handleRowClick ({ r: { recordID } }) {
      if ((this.options.editable && this.editing) || (!this.recordPageID && !this.options.rowViewUrl)) {
        return
      }

      if (this.options.enableRecordPageNavigation) {
        this.loadPaginationRecords({
          filter: {
            ...this.filter,
            limit: Math.min(this.filter.limit, 100),
          },
        })
      }

      if (this.options.recordDisplayOption === 'modal' || this.inModal) {
        this.$root.$emit('show-record-modal', {
          recordID,
          recordPageID: this.recordPageID,
          edit: this.options.openRecordInEditMode,
        })
        return
      }

      const pageID = this.recordPageID
      const name = this.options.openRecordInEditMode ? this.options.rowEditUrl || 'page.record.edit' : this.options.rowViewUrl || 'page.record'
      const route = {
        name,
        params: {
          pageID,
          recordID,
        },
        query: null,
      }

      if (this.options.recordDisplayOption === 'newTab') {
        window.open(this.$router.resolve(route).href)
      } else {
        this.$router.push(route)
      }
    },

    handleSort ({ key, sortable }) {
      if (!sortable) {
        return
      }

      if (this.sortBy !== key) {
        this.filter.sort = `${key}`
        this.sortDirecton = 'ASC'
      } else {
        if (this.sortDirecton === 'ASC') {
          this.filter.sort = `${key} DESC`
          this.sortDirecton = 'DESC'
        } else {
          this.filter.sort = `${key}`
          this.sortDirecton = 'ASC'
        }
      }
      this.sortBy = key
      this.refresh(true)
    },

    goToPage (page) {
      if (page >= 1) {
        this.filter.pageCursor = (this.pagination.pages[page - 1] || {}).cursor
        this.pagination.page = page
      } else {
        this.filter.pageCursor = this.filter[page]
        if (this.filter.pageCursor) {
          this.pagination.page += page === 'nextPage' ? 1 : -1
        } else {
          this.pagination.page = 1
        }
      }
      this.refresh()
    },

    handleSelectAllOnPage ({ isChecked }) {
      if (isChecked) {
        this.selected = this.items.map(({ id }) => id)
      } else {
        this.selected = []
        this.selectedAllRecords = isChecked
      }
    },

    selectAllRecords () {
      this.selectedAllRecords = !this.selectedAllRecords
      this.handleSelectAllOnPage({ isChecked: this.selectedAllRecords })
    },

    handleRestoreSelectedRecords (recordID) {
      if (this.inlineEditing) {
        const sel = new Set(this.selected)
        this.items.forEach((item, index) => {
          if (sel.has(item.id)) {
            this.handleRestoreInline(item, index)
          }
        })
        sel.clear()
      } else {
        this.processing = true

        const query = recordID ? `recordID = ${recordID}` : this.bulkQuery
        const { moduleID, namespaceID } = this.filter

        const { response, cancel } = this.$ComposeAPI
          .recordBulkUndeleteCancellable({ moduleID, namespaceID, query })

        this.abortableRequests.push(cancel)

        response()
          .then(() => {
            this.refresh(true)
            this.toastSuccess(this.$t('notification:record.restoreBulkSuccess'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:record.restoreBulkFailed')))
          .finally(() => {
            setTimeout(() => {
              this.processing = false
              this.selectedAllRecords = false
            }, 300)
          })
      }
    },

    handleDeleteSelectedRecords (recordID) {
      if (this.inlineEditing) {
        const sel = new Set(this.selected)
        for (let i = 0; i < this.items.length; i++) {
          if (sel.has(this.items[i].id)) {
            this.handleDeleteInline(this.items[i], i)
          }
        }
        sel.clear()
      } else {
        this.processing = true

        const query = recordID ? `recordID = ${recordID}` : this.bulkQuery
        // Pick module and namespace ID from the filter
        const { moduleID, namespaceID } = this.filter

        const { response, cancel } = this.$ComposeAPI
          .recordBulkDeleteCancellable({ moduleID, namespaceID, query })

        this.abortableRequests.push(cancel)

        response()
          .then(() => this.refresh(true))
          .then(() => {
            this.toastSuccess(this.$t('notification:record.deleteBulkSuccess'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:record.deleteBulkFailed')))
          .finally(() => {
            setTimeout(() => {
              this.processing = false
              this.selectedAllRecords = false
            }, 300)
          })
      }
    },

    async refresh (resetPagination = false, checkSelected = false) {
      // Prevent refresh if records are selected or inline editing
      if (checkSelected && (this.selected.length || this.inlineEdit.recordIDs.length)) return

      await this.$nextTick()
      return this.pullRecords(resetPagination)
    },

    /**
     * Loader for b-table
     *
     * Will ignore b-tables input arguments for filter
     * and assemble them on our own
     */
    async pullRecords (resetPagination = false) {
      if (!this.recordListModule) {
        return
      }

      if (this.recordListModule.moduleID !== this.options.moduleID) {
        throw Error(this.$t('record.moduleMismatch'))
      }

      this.processing = true
      this.selected = []

      // Compute query based on query, prefilter and recordListFilter
      const query = queryToFilter(this.query, this.drillDownFilter || this.prefilter, this.fields.map(({ moduleField }) => moduleField), this.recordListFilter)

      const { moduleID, namespaceID } = this.recordListModule

      let paginationOptions = {}
      if (resetPagination) {
        this.filter.pageCursor = undefined
        const { fullPageNavigation = false, showTotalCount = false } = this.options
        paginationOptions = {
          incPageNavigation: fullPageNavigation,
          incTotal: showTotalCount,
        }
      } else if (this.filter.pageCursor) {
        this.filter.sort = ''
      }

      // Filter's out deleted records when filter.deleted is 2, and undeleted records when filter.deleted is 0
      this.showingDeletedRecords ? this.filter.deleted = 2 : this.filter.deleted = 0

      const { response, cancel } = this.$ComposeAPI.recordListCancellable({ ...this.filter, moduleID, namespaceID, query, ...paginationOptions })
      this.abortableRequests.push(cancel)

      return response()
        .then(({ set, filter }) => {
          const records = set.map(r => new compose.Record(r, this.recordListModule))

          this.filter = { ...this.filter, ...filter }
          this.filter.nextPage = filter.nextPage
          this.filter.prevPage = filter.prevPage

          if (resetPagination) {
            let count = this.pagination.count || 0

            if (paginationOptions.incTotal) {
              count = filter.total || 0
              this.filter.incTotal = false
            }

            if (paginationOptions.incPageNavigation) {
              const pages = filter.pageNavigation || []
              this.pagination.pages = pages

              if (!paginationOptions.incTotal) {
                if (pages.length > 1) {
                  const lastPageCount = pages[pages.length - 1].items
                  count = ((pages.length - 1) * this.recordsPerPage) + lastPageCount
                } else {
                  count = records.length
                }
              }

              this.filter.incPageNavigation = false
            }

            this.pagination.count = count
            this.pagination.page = 1
          }

          // Extract user IDs from record values and load all users
          const fields = this.fields.filter(f => f.moduleField).map(f => f.moduleField)

          return Promise.all([
            this.fetchUsers(fields, records),
            this.fetchRecords(namespaceID, fields, records),
          ]).then(() => {
            this.items = records.map(r => this.wrapRecord(r))
          })
        })
        .catch((e) => {
          if (!axios.isCancel(e)) {
            this.toastErrorHandler(this.$t('notification:record.listLoadFailed'))(e)
          }
        })
        .finally(() => {
          setTimeout(() => {
            this.processing = false
          }, 300)
        })
    },

    getStorageRecordListFilter () {
      try {
        // Get record list filters from localStorage
        const currentFilters = getItem(`record-list-filters-${this.uniqueID}`)

        // Check type of filter value
        if (!Array.isArray(currentFilters)) {
          console.warn(this.$t('notification:record-list.incorrect-filter-structure', { filterID: this.uniqueID }))
          // Remove the filter from the local storage if the type doesn't match
          removeItem(`record-list-filters-${this.uniqueID}`)
        } else {
          this.recordListFilter = currentFilters
          this.activeFilters = [...new Set(currentFilters.map(f => f.name).filter(f => !!f))]
        }
      } catch (e) {
        // Land here if the filter is corrupted
        console.warn(this.$t('notification:record-list.corrupted-filter'))
        // Remove filter from the local storage
        removeItem(`record-list-filters-${this.uniqueID}`)
      }
    },

    getStorageRecordListFilterPreset () {
      try {
        // Get record list filters from localStorage
        const currentFilterPresets = getItem(`record-list-preset-${this.uniqueID}`)

        // Set the custom preset filters
        this.customPresetFilters = currentFilterPresets
      } catch (e) {
        // Land here if the filter is corrupted
        console.warn(this.$t('notification:record-list.corrupted-filter'))
        // Remove filter from the local storage
        removeItem(`record-list-filters-${this.uniqueID}`)
      }
    },

    setStorageRecordListFilter () {
      let currentListFilters = []

      try {
        // Get record list filters from localStorage
        currentListFilters = this.recordListFilter
        setItem(`record-list-filters-${this.uniqueID}`, currentListFilters)
      } catch (e) {
        console.warn(this.$t('notification:record-list.corrupted-filter'))
      }
    },

    setStorageRecordListFilterPreset ({ name }) {
      this.showCustomPresetFilterModal = false

      const currentListFilters = [...this.customPresetFilters]
      currentListFilters.push({ ...this.currentCustomPresetFilter, name })

      this.customPresetFilters = currentListFilters
      this.updateFilter(this.currentCustomPresetFilter.filter, name)

      try {
        setItem(`record-list-preset-${this.uniqueID}`, currentListFilters)
      } catch (e) {
        console.warn(this.$t('notification:record-list.corrupted-filter'))
      }
    },

    removeStorageRecordListFilterPreset (name) {
      this.customPresetFilters = this.customPresetFilters.filter(f => f.name !== name)

      if (this.$refs.filterPresets) {
        this.$refs.filterPresets.hide(true)
      }

      try {
        setItem(`record-list-preset-${this.uniqueID}`, this.customPresetFilters)
      } catch (e) {
        console.warn(this.$t('notification:record-list.corrupted-filter'))
      }
    },

    onImportSuccessful () {
      this.refresh(true)
    },

    setDrillDownFilter (drillDownFilter) {
      if (!this.drillDownFilter) {
        this.activeFilters.push(this.$t('recordList.drillDown.filter.label'))
      }

      this.drillDownFilter = drillDownFilter
      this.pullRecords(true)
    },

    isInlineRestoreActionVisible ({ deletedAt }) {
      return !!deletedAt
    },

    isInlineDeleteActionVisible ({ recordID, canDeleteRecord, deletedAt }) {
      return !deletedAt && (canDeleteRecord || recordID === NoID)
    },

    isViewRecordActionVisible ({ canReadRecord }) {
      return !this.options.hideRecordViewButton && canReadRecord && (this.options.rowViewUrl || this.recordPageID)
    },

    isEditRecordActionVisible ({ canUpdateRecord }) {
      return !this.options.hideRecordEditButton && canUpdateRecord && (this.options.rowEditUrl || this.recordPageID)
    },

    isRecordPermissionButtonVisible ({ canGrant }) {
      return canGrant && !this.options.hideRecordPermissionsButton
    },

    isDeleteActionVisible ({ deletedAt, canDeleteRecord }) {
      return !deletedAt && canDeleteRecord
    },

    isRestoreActionVisible ({ canUndeleteRecord }) {
      return canUndeleteRecord
    },

    areActionsVisible (record) {
      if (this.inlineEditing) {
        return [
          this.isCloneRecordActionVisible,
          this.isInlineDeleteActionVisible(record),
          this.isInlineRestoreActionVisible(record),
        ].some(v => v)
      }

      return [
        this.isCloneRecordActionVisible,
        this.isReminderActionVisible,
        this.isViewRecordActionVisible(record),
        this.isEditRecordActionVisible(record),
        this.isRecordPermissionButtonVisible(record),
        this.isDeleteActionVisible(record),
        this.isRestoreActionVisible(record),
      ].some(v => v)
    },

    onBulkUpdate () {
      this.selectedAllRecords = false
    },

    editInlineField (record, field) {
      this.inlineEdit.fields = [field]
      this.inlineEdit.record = record.clone()
      this.inlineEdit.query = `recordID = ${record.recordID}`
    },

    onInlineEditClose () {
      this.inlineEdit.fields = []
      this.inlineEdit.record = {}
      this.inlineEdit.query = ''
    },

    onInlineEdit () {
      this.onInlineEditClose()
    },

    isFieldEditable (field) {
      if (!field) return false

      const { canCreateOwnedRecord } = this.recordListModule || {}
      const { createdAt, canManageOwnerOnRecord } = this.record || {}
      const { name, canUpdateRecordValue, isSystem, expressions = {} } = field || {}

      if (!canUpdateRecordValue) return false

      if (isSystem) {
        // Make ownedBy field editable if correct permissions
        if (name === 'ownedBy') {
          // If not created we check module permissions, otherwise the canManageOwnerOnRecord
          return createdAt ? canManageOwnerOnRecord : canCreateOwnedRecord
        }

        return false
      }

      return !expressions.value
    },

    updateFilter (filter = [], name) {
      const lastFilterIdx = this.recordListFilter.length - 1
      filter = filter.map((filter) => ({ ...filter, name }))

      if (this.recordListFilter.length) {
        this.recordListFilter[lastFilterIdx].groupCondition = 'AND'
      }

      this.recordListFilter = this.recordListFilter.concat(filter)
      this.activeFilters.push(name)
      this.refresh(true)

      if (this.$refs.filterPresets) {
        this.$refs.filterPresets.hide(true)
      }
    },

    removeFilter (filterIndex) {
      this.activeFilters.splice(filterIndex, 1)

      if (this.drillDownFilter && !this.activeFilters.includes(this.$t('recordList.drillDown.filter.label'))) {
        this.setDrillDownFilter(undefined)
      }

      this.recordListFilter = this.recordListFilter.filter(({ name }) => !name || this.activeFilters.includes(name))

      this.setStorageRecordListFilter()
      this.refresh(true)
    },

    isUserRoleMember (roles) {
      if (!roles.length) return true

      return roles.some(roleID => this.authUserRoles.includes(roleID))
    },

    setDefaultValues () {
      this.uniqueID = undefined
      this.processing = false
      this.prefilter = undefined
      this.recordListFilter = []
      this.drillDownFilter = undefined
      this.query = null
      this.filter = {}
      this.pagination = {}
      this.selected = []
      this.inlineEdit = {}
      this.sortBy = undefined
      this.sortDirecton = undefined
      this.ctr = 0
      this.items = []
      this.showingDeletedRecords = false
      this.activeFilters = []
      this.customPresetFilters = []
      this.currentCustomPresetFilter = undefined
      this.showCustomPresetFilterModal = false
      this.selectedAllRecords = false
      this.abortableRequests = []
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },

    refreshAndResetPagination () {
      this.refresh(true)
    },

    destroyEvents () {
      const { pageID = NoID } = this.page

      this.$root.$off(`record-line:collect:${this.uniqueID}`, this.resolveRecords)
      this.$root.$off(`page-block:validate:${this.uniqueID}`, this.validatePageBlock)
      this.$root.$off(`drill-down-recordList:${this.uniqueID}`, this.setDrillDownFilter)
      this.$root.$off(`refetch-non-record-blocks:${pageID}`, this.refreshAndResetPagination)
      this.$root.$off('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$off('record-field-change', this.refetchOnPrefilterValueChange)
    },

    handleAddRecord () {
      const refRecord = this.options.refField && this.recordID !== NoID ? this.record : undefined
      const pageID = this.recordPageID

      if (!(pageID || this.options.rowCreateUrl)) return

      if (this.inModal) {
        this.$root.$emit('show-record-modal', {
          recordID: NoID,
          recordPageID: this.recordPageID,
          refRecord,
          edit: true,
        })
      } else {
        this.$router.push({
          name: this.options.rowCreateUrl || 'page.record.create',
          params: { pageID, refRecord },
          query: null,
          edit: true,
        })
      }
    },

    handleViewRecordAction (recordID) {
      if (this.inModal) {
        this.$root.$emit('show-record-modal', {
          recordID: recordID,
          recordPageID: this.recordPageID,
          edit: false,
        })
      } else {
        this.$router.push({
          name: this.options.rowViewUrl || 'page.record',
          params: { pageID: this.recordPageID, recordID },
          query: null,
          edit: false,
        })
      }
    },

    handleEditRecordAction (recordID) {
      if (this.inModal) {
        this.$root.$emit('show-record-modal', {
          recordID: recordID,
          recordPageID: this.recordPageID,
          edit: true,
        })
      } else {
        this.$router.push({
          name: this.options.rowEditUrl || 'page.record.edit',
          params: { pageID: this.recordPageID, recordID },
          query: null,
          edit: true,
        })
      }
    },

    handleCloneRecordAction (recordID, values) {
      if (this.inModal) {
        this.$root.$emit('show-record-modal', {
          recordID: recordID,
          recordPageID: this.recordPageID,
          values,
          edit: true,
        })
      } else {
        this.$router.push({
          name: this.options.rowCreateUrl || 'page.record.create',
          params: { pageID: this.recordPageID, values },
          query: null,
          edit: true,
        })
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.handle {
  cursor: grab;
}

.pointer {
  cursor: pointer;
}

th .required::after {
  content: "*";
  display: inline-block;
  color: var(--primary);
  vertical-align: sub;
  margin-left: 2px;
  width: 10px;
  height: 16px;
  overflow: hidden;
}

tr:hover td.actions {
  opacity: 1;
  z-index: 1;
  background-color: var(--light);
}

.inline-actions {
  min-width: 30px;
  margin-top: -2px;
  opacity: 0;
  transition: opacity 0.25s;
}

td:hover .inline-actions {
  opacity: 1;

  button:hover {
    color: var(--primary) !important;
  }
}
</style>

<style lang="scss">
.record-list-table {
  .actions {
    padding-top: 8px;
    position: sticky;
    right: -1px;
    opacity: 0;
    transition: opacity 0.25s;
    width: 1%;
    font-family: var(--font-regular) !important;
  }
}

.record-list-footer {
  font-family: var(--font-medium);
}
</style>
