<template>
  <wrap
    v-if="recordListModule"
    v-bind="$props"
    :scrollable-body="false"
    v-on="$listeners"
    @refreshBlock="refresh(false, false)"
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
        class="py-2 d-print-none"
        fluid
      >
        <b-row
          no-gutters
          class="justify-content-between wrap-with-vertical-gutters"
        >
          <div class="text-nowrap flex-grow-1">
            <div
              class="wrap-with-vertical-gutters"
            >
              <template v-if="recordListModule.canCreateRecord">
                <template v-if="inlineEditing">
                  <b-btn
                    v-if="!options.hideAddButton"
                    data-test-id="button-add-record"
                    variant="primary"
                    size="lg"
                    class="float-left mr-1"
                    @click="addInline"
                  >
                    + {{ $t('recordList.addRecord') }}
                  </b-btn>
                </template>

                <template v-else-if="!inlineEditing && (recordPageID || options.allRecords)">
                  <router-link
                    v-if="!options.hideAddButton"
                    data-test-id="button-add-record"
                    class="btn btn-lg btn-primary float-left mr-1"
                    :to="newRecordRoute"
                  >
                    + {{ $t('recordList.addRecord') }}
                  </router-link>
                  <importer-modal
                    v-if="!options.hideImportButton"
                    :module="recordListModule"
                    :namespace="namespace"
                    class="mr-1 float-left"
                    @importSuccessful="onImportSuccessful"
                  />
                </template>
              </template>

              <exporter-modal
                v-if="options.allowExport && !inlineEditing"
                :module="recordListModule"
                :record-count="pagination.count"
                :query="query"
                :prefilter="prefilter"
                :selection="selected"
                :processing="processing"
                :preselected-fields="fields.map(({ moduleField }) => moduleField)"
                class="mr-1 float-left"
                @export="onExport"
              />

              <b-dropdown
                v-if="filterPresets.length"
                size="lg"
                variant="light"
                right
                :text="$t('recordList.filter.filters.label')"
              >
                <b-dropdown-item
                  v-for="(f, idx) in filterPresets"
                  :key="idx"
                  :disabled="activeFilters.includes(f.name)"
                  @click="updateFilter(f.filter, f.name)"
                >
                  {{ f.name }}
                </b-dropdown-item>
              </b-dropdown>

              <column-picker
                v-if="!options.hideConfigureFieldsButton"
                :module="recordListModule"
                :fields="fields"
                class="float-left"
                @updateFields="onUpdateFields"
              />
            </div>
          </div>
          <div
            v-if="!options.hideSearch"
            class="flex-grow-1 w-25"
          >
            <c-input-search
              v-model.trim="query"
              :placeholder="$t('general.label.search')"
            />
          </div>
        </b-row>

        <div
          v-if="activeFilters.length || drillDownFilter || options.showDeletedRecordsOption"
          class="d-flex mt-2"
        >
          <div
            v-if="activeFilters.length"
            class="d-flex align-items-center flex-wrap"
          >
            {{ $t('recordList.filter.filters.active') }}
            <b-form-tags
              v-model="activeFilters"
              tag-variant="secondary"
              tag-pills
              size="lg"
              input-class="d-none"
              tag-class="align-items-center"
              class="filter-tags border-0 p-0 ml-1"
              style="width: fit-content;"
              @input="removeFilter"
            />
          </div>

          <b-button
            v-if="options.showDeletedRecordsOption"
            variant="outline-light"
            size="sm"
            class="text-primary border-0 text-nowrap ml-auto"
            @click="handleShowDeleted()"
          >
            {{ showingDeletedRecords ? $t('recordList.showRecords.existing') : $t('recordList.showRecords.deleted') }}
          </b-button>
        </div>

        <div
          class="d-none flex-wrap align-items-center mt-1"
          :class="{ 'd-flex': options.selectable && selected.length }"
        >
          <div
            class="d-flex align-items-baseline my-auto pt-1 text-nowrap h-100"
          >
            {{ $t('recordList.selected', { count: selected.length, total: items.length }) }}
            <b-button
              variant="link"
              class="p-0 text-decoration-none"
              @click.prevent="handleSelectAllOnPage({ isChecked: false })"
            >
              ({{ $t('recordList.cancelSelection') }})
            </b-button>
          </div>

          <div class="d-flex align-items-center ml-auto">
            <automation-buttons
              class="d-inline m-0 mr-2"
              :buttons="options.selectionButtons"
              :module="recordListModule"
              :extra-event-args="{ selected, filter }"
              v-bind="$props"
              @refresh="refresh()"
            />

            <bulk-edit-modal
              v-show="options.bulkRecordEditEnabled && canUpdateSelectedRecords"
              :module="recordListModule"
              :namespace="namespace"
              :selected-records="selected"
              @save="onBulkUpdate()"
            />

            <template v-if="canDeleteSelectedRecords && !areAllRowsDeleted">
              <c-input-confirm
                v-if="!inlineEditing"
                :tooltip="$t('recordList.tooltip.deleteSelected')"
                @confirmed="handleDeleteSelectedRecords()"
              />
              <b-button
                v-else
                variant="link"
                size="md"
                :title="$t('recordList.tooltip.deleteSelected')"
                class="text-danger"
                @click.prevent="handleDeleteSelectedRecords()"
              >
                <font-awesome-icon
                  class="text-danger"
                  :icon="['far', 'trash-alt']"
                />
              </b-button>
            </template>

            <template v-if="canUndeleteSelectedRecords && areAllRowsDeleted">
              <c-input-confirm
                v-if="!inlineEditing"
                :tooltip="$t('recordList.tooltip.undeleteSelected')"
                @confirmed="handleRestoreSelectedRecords()"
              >
                <font-awesome-icon
                  :icon="['fa', 'trash-restore']"
                />
              </c-input-confirm>
              <b-button
                v-else
                variant="link"
                size="md"
                :title="$t('recordList.tooltip.undeleteSelected')"
                class="text-danger"
                @click.prevent="handleRestoreSelectedRecords()"
              >
                <font-awesome-icon
                  :icon="['fa', 'trash-restore']"
                />
              </b-button>
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
          class="record-list-table border-top mh-100 h-100 mb-0"
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
                  cursor: field.sortable ? 'pointer' : 'default',
                }"
                @click="handleSort(field)"
              >
                <div
                  class="d-flex align-self-center"
                >
                  <div
                    :class="{ required: field.required }"
                    class="d-flex align-self-center text-nowrap"
                  >
                    {{ field.label }}
                  </div>
                  <div
                    class="d-flex"
                  >
                    <record-list-filter
                      v-if="!options.hideFiltering && field.filterable"
                      :target="uniqueID"
                      :selected-field="field.moduleField"
                      :namespace="namespace"
                      :module="recordListModule"
                      :record-list-filter="recordListFilter"
                      class="d-print-none ml-1"
                      @filter="onFilter"
                      @reset="activeFilters = []"
                    />

                    <b-button
                      v-if="field.sortable"
                      variant="link p-0 ml-1"
                      :title="$t('recordList.sort.tooltip')"
                      class="d-flex align-items-center justify-content-center"
                    >
                      <font-awesome-layers
                        class="d-print-none"
                      >
                        <font-awesome-icon
                          :icon="['fas', 'angle-up']"
                          class="mb-1"
                          :style="{
                            color: 'gray',
                            ...sorterStyle(field, 'ASC'),
                          }"
                        />
                        <font-awesome-icon
                          :icon="['fas', 'angle-down']"
                          class="mt-1"
                          :style="{
                            color: 'gray',
                            ...sorterStyle(field, 'DESC'),
                          }"
                        />
                      </font-awesome-layers>
                    </b-button>
                  </div>
                </div>
              </b-th>
            </b-tr>
          </b-thead>

          <draggable
            v-if="items.length && !processing"
            v-model="items"
            :disabled="!inlineEditing || !options.draggable"
            group="items"
            tag="b-tbody"
            handle=".handle"
          >
            <b-tr
              v-for="(item, index) in items"
              :key="`${index}${item.r.recordID}`"
              :class="{ 'pointer': !(options.editable && editing) }"
              @click="handleRowClicked(item)"
            >
              <b-td
                v-if="options.draggable && inlineEditing"
                class="pr-0"
                @click.stop
              >
                <font-awesome-icon
                  v-b-tooltip.hover
                  :icon="['fas', 'bars']"
                  :title="$t('general.tooltip.dragAndDrop')"
                  class="handle text-light my-1"
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
                  style="min-width: 150px;"
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
                    v-if="options.inlineRecordEditEnabled && field.canEdit"
                    class="inline-actions"
                  >
                    <b-button
                      :title="$t('recordList.inlineEdit.button.title')"
                      variant="outline-light"
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
                  boundary="scrollParent"
                  variant="outline-light"
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
                        :icon="['fa', 'trash-restore']"
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
                      :to="{ name: options.rowViewUrl || 'page.record', params: { pageID: recordPageID, recordID: item.r.recordID }, query: null }"
                    >
                      <font-awesome-icon
                        :icon="['far', 'file-alt']"
                        class="text-primary"
                      />
                      {{ $t('recordList.record.tooltip.view') }}
                    </b-dropdown-item>

                    <b-dropdown-item
                      v-if="isEditRecordActionVisible(item.r)"
                      :to="{ name: options.rowEditUrl || 'page.record.edit', params: { pageID: recordPageID, recordID: item.r.recordID }, query: null }"
                    >
                      <font-awesome-icon
                        :icon="['far', 'edit']"
                        class="text-primary"
                      />
                      {{ $t('recordList.record.tooltip.edit') }}
                    </b-dropdown-item>

                    <b-dropdown-item
                      v-if="isCloneRecordActionVisible"
                      :to="{ name: options.rowCreateUrl || 'page.record.create', params: { pageID: recordPageID, values: item.r.values }, query: null }"
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
                        button-variant="link dropdown-item text-decoration-none text-dark regular-font rounded-0"
                      />
                    </b-dropdown-item>

                    <c-input-confirm
                      v-if="isDeleteActionVisible(item.r)"
                      borderless
                      variant="link"
                      size="md"
                      button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
                      class="w-100"
                      @confirmed="handleDeleteSelectedRecords(item.r.recordID)"
                    >
                      <font-awesome-icon
                        :icon="['far', 'trash-alt']"
                        class="text-danger"
                      />
                      {{ $t('recordList.record.tooltip.delete') }}
                    </c-input-confirm>
                  </template>
                </b-dropdown>
              </b-td>
            </b-tr>
          </draggable>

          <div
            v-else
            class="position-absolute text-center mt-5 d-print-none"
            style="left: 0; right: 0;"
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
        :selected-records="inlineEdit.recordIDs"
        :selected-fields="inlineEdit.fields"
        :initial-record="inlineEdit.record"
        :modal-title="$t('recordList.inlineEdit.modal.title')"
        open-on-select
        @save="onInlineEdit()"
      />
    </template>

    <template
      #footer
    >
      <div
        v-if="showFooter"
        class="d-flex align-items-center justify-content-between p-2"
      >
        <div class="text-truncate">
          <div
            v-if="options.showTotalCount"
            class="ml-2 text-nowrap my-1"
          >
            <span
              v-if="pagination.count > options.perPage"
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
        </div>

        <div
          v-if="showPageNavigation"
          class="d-flex align-items-center justify-content-end"
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

          <b-button-group v-else>
            <b-button
              :disabled="!hasPrevPage || processing"
              data-test-id="first-page"
              variant="outline-light"
              class="d-flex align-items-center justify-content-center text-primary border-0"
              @click="goToPage()"
            >
              <font-awesome-icon :icon="['fas', 'angle-double-left']" />
            </b-button>
            <b-button
              :disabled="!hasPrevPage || processing"
              data-test-id="previous-page"
              variant="outline-light"
              class="d-flex align-items-center justify-content-center text-primary border-0"
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
              variant="outline-light"
              class="d-flex align-items-center justify-content-center text-primary border-0"
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
import { debounce, isEqual } from 'lodash'
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
import { evaluatePrefilter, queryToFilter } from 'corteza-webapp-compose/src/lib/record-filter'
import { getItem, setItem, removeItem } from 'corteza-webapp-compose/src/lib/local-storage'
import { components, url } from '@cortezaproject/corteza-vue'
import draggable from 'vuedraggable'
import RecordListFilter from 'corteza-webapp-compose/src/components/Common/RecordListFilter'
import ColumnPicker from 'corteza-webapp-compose/src/components/Admin/Module/Records/ColumnPicker'
import BulkEditModal from 'corteza-webapp-compose/src/components/Public/Record/BulkEdit'

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
      return this.showPageNavigation || this.options.showTotalCount
    },

    getPagination () {
      const { page = 1, count = 0 } = this.pagination
      const { perPage = 10 } = this.options

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
      return this.items.length && !this.options.hidePaging
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

      if (this.options.fields.length > 0) {
        fields = this.recordListModule.filterFields(this.options.fields)
      } else {
        // Record list block does not have any configured fields
        // Use first five fields from the module.
        fields = this.recordListModule.fields.slice(0, 5)
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

    canUndeleteSelectedRecords () {
      return this.items.filter(({ id, r }) => this.selected.includes(id) && r.canUndeleteRecord).length
    },

    newRecordRoute () {
      const refRecord = this.options.linkToParent ? this.record : undefined
      const pageID = this.recordPageID

      if (pageID || this.options.rowCreateUrl) {
        return {
          name: this.options.rowCreateUrl || 'page.record.create',
          params: { pageID, refRecord },
          query: null,
        }
      }

      return undefined
    },

    isCloneRecordActionVisible () {
      return !this.options.hideRecordCloneButton && this.recordListModule.canCreateRecord && (this.options.rowCreateUrl || this.recordPageID || this.inlineEditing)
    },

    isReminderActionVisible () {
      return !this.options.hideRecordReminderButton
    },

    filterPresets () {
      return this.options.filterPresets.filter(({ name, roles }) => name && this.isUserRoleMember(roles))
    },

    authUserRoles () {
      return this.$auth.user.roles
    },
  },

  watch: {
    query: debounce(function (e) {
      this.refresh(true)
    }, 500),

    options: {
      deep: true,
      handler () {
        this.prepRecordList()
        this.refresh(true)
      },
    },

    'record.recordID': {
      immediate: true,
      handler () {
        this.createEvents()
        this.getStorageRecordListFilter()
        this.prepRecordList()
        this.refresh(true)
      },
    },
  },

  mounted () {
    this.createEvents()
  },

  beforeDestroy () {
    this.destroyEvents()
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
        this.$root.$off(`record-line:collect:${this.uniqueID}`)
        this.$root.$off(`page-block:validate:${this.uniqueID}`)
        this.$root.$off(`drill-down-recordList:${this.uniqueID}`)
        this.$root.$off(`refetch-non-record-blocks:${pageID}`)
      }

      this.uniqueID = [pageID, recordID, this.block.blockID].map(v => v || NoID).join('-')
      this.$root.$on(`record-line:collect:${this.uniqueID}`, this.resolveRecords)
      this.$root.$on(`page-block:validate:${this.uniqueID}`, this.validatePageBlock)
      this.$root.$on(`drill-down-recordList:${this.uniqueID}`, this.setDrillDownFilter)
      this.$root.$on(`refetch-non-record-blocks:${pageID}`, () => {
        this.refresh(true)
      })
    },

    destroyEvents () {
      this.$root.$off(`record-line:collect:${this.uniqueID}`)
      this.$root.$off(`page-block:validate:${this.uniqueID}`)
      this.$root.$off(`drill-down-recordList:${this.uniqueID}`)
      this.$root.$off(`refetch-non-record-blocks:${this.page.pageID}`)
    },

    onFilter (filter = []) {
      filter.forEach(f => {
        if (this.activeFilters.includes(f.name)) {
          this.filterPresets.find(p => p.name === f.name).filter.forEach((filterPreset) => {
            if (!isEqual(f.filter, filterPreset.filter)) {
              const filterIndex = this.activeFilters.indexOf(f.name)
              this.activeFilters.splice(filterIndex, 1)
            }
          })
        }
      })

      this.recordListFilter = filter
      this.setStorageRecordListFilter()
      this.refresh(true)
    },

    onUpdateFields (fields = []) {
      this.options.fields = [...fields]
      this.$emit('save-fields', this.options.fields)
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
      }
    },

    sorterStyle ({ key }, dir) {
      const { sort = '' } = this.filter

      const sortedFields = (sort.includes(',') ? sort.split(',') : [sort])

      const isSorted = sortedFields.map(v => v.trim()).some(value => {
        let valueDir = 'ASC'

        if (value.includes(' ')) {
          value = value.split(' ')[0]
          valueDir = 'DESC'
        }

        return valueDir === dir && value === key
      })

      return isSorted ? { color: 'black' } : {}
    },

    handleShowDeleted () {
      this.showingDeletedRecords = !this.showingDeletedRecords
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

    addInline () {
      const r = new compose.Record(this.recordListModule, {})

      // Set record values that should be prefilled
      if (this.record.recordID && this.options.linkToParent) {
        r.values[this.options.refField] = this.record.recordID
      }

      this.items.unshift(this.wrapRecord(r))
    },

    /**
     * Helper method to fetch all records available to this record list
     * at the given point in time.
     *
     * It:
     *    * assures that local records have a sequencial indexing
     *    * appends aditional meta fields
     *    * resolves payloadediting
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
      const { moduleID, presort, prefilter, editable, perPage, refField, positionField } = this.options

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

      const filter = []
      let sort = ''

      if (presort) {
        sort = presort
      }

      // Initial filter
      if (prefilter) {
        const pf = evaluatePrefilter(prefilter, {
          record: this.record,
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
        filter.push(`(${pf})`)
      }

      if (editable) {
        if (positionField) {
          sort = `${positionField}`
        }

        if (refField && this.record.recordID) {
          filter.push(`(${refField} = ${this.record.recordID})`)
        }
      }

      this.prefilter = filter.join(' AND ')
      const limit = perPage
      this.filter = {
        limit,
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
          filter: filter,
          jwt: this.$auth.accessToken,
          timezone: timezone ? timezone.tzCode : undefined,
        },
      })

      window.open(exportUrl)
      this.processing = false
    },

    handleRowClicked ({ r: { recordID } }) {
      if ((this.options.editable && this.editing) || (!this.recordPageID && !this.options.rowViewUrl)) {
        return
      }

      if (this.options.enableRecordPageNavigation) {
        this.loadPaginationRecords({
          filter: {
            ...this.filter,
            limit: Math.min(this.pagination.count, 100),
          },
        })
      }

      const pageID = this.recordPageID
      const route = {
        name: this.options.rowViewUrl || 'page.record',
        params: {
          pageID,
          recordID,
        },
        query: null,
      }

      if (this.options.recordDisplayOption === 'newTab') {
        window.open(this.$router.resolve(route).href)
      } else if (this.options.recordDisplayOption === 'modal') {
        this.$root.$emit('show-record-modal', {
          recordID,
          recordPageID: this.recordPageID,
        })
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
      }
    },

    handleRestoreSelectedRecords () {
      if (this.inlineEditing) {
        const sel = new Set(this.selected)
        this.items.forEach((item, index) => {
          if (sel.has(item.id)) {
            this.handleRestoreInline(item, index)
          }
        })
      } else {
        const { moduleID, namespaceID } = this.items[0].r

        // filter undeletable records from the selected list
        const recordIDs = this.items
          .filter(({ id, r }) => r.canUndeleteRecord && this.selected.includes(id))
          .map(({ id }) => id)

        this.processing = true

        this.$ComposeAPI
          .recordBulkUndelete({ moduleID, namespaceID, recordIDs })
          .then(() => {
            this.refresh(true)
            this.toastSuccess(this.$t('notification:record.undeleteBulkSuccess'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:record.undeleteBulkFailed')))
          .finally(() => {
            this.processing = false
          })
      }
    },

    handleDeleteSelectedRecords (selected = this.selected) {
      if (selected.length === 0) {
        return
      }

      if (this.inlineEditing) {
        const sel = new Set(selected)
        for (let i = 0; i < this.items.length; i++) {
          if (sel.has(this.items[i].id)) {
            this.handleDeleteInline(this.items[i], i)
          }
        }
      } else {
        // Pick module and namespace ID from the first record
        //
        // We are always showing list of records from the
        // same module so this should be safe to do.
        const { moduleID, namespaceID } = this.items[0].r

        // filter deletable records from the selected list
        const recordIDs = this.items
          .filter(({ id, r }) => r.canDeleteRecord && selected.includes(id))
          .map(({ id }) => id)

        this.processing = true

        this.$ComposeAPI
          .recordBulkDelete({ moduleID, namespaceID, recordIDs })
          .then(() => this.refresh(true))
          .then(() => {
            this.toastSuccess(this.$t('notification:record.deleteBulkSuccess'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:record.deleteBulkFailed')))
          .finally(() => {
            this.processing = false
          })
      }
    },

    refresh (resetPagination = false, checkSelected = false) {
      // Prevent refresh if records are selected or inline editing
      if (checkSelected && (this.selected.length || this.inlineEdit.recordIDs.length)) return

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
      if (this.filter.pageCursor) {
        this.filter.sort = ''
      }

      let paginationOptions = {}
      if (resetPagination) {
        this.filter.pageCursor = undefined
        const { fullPageNavigation = false, showTotalCount = false } = this.options
        paginationOptions = {
          incPageNavigation: fullPageNavigation,
          incTotal: showTotalCount,
        }
      }

      // Filter's out deleted records when filter.deleted is 2, and undeleted records when filter.deleted is 0
      this.showingDeletedRecords ? this.filter.deleted = 2 : this.filter.deleted = 0

      await this.$ComposeAPI.recordList({ ...this.filter, moduleID, namespaceID, query, ...paginationOptions })
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
                  count = ((pages.length - 1) * this.options.perPage) + lastPageCount
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
        .catch(this.toastErrorHandler(this.$t('notification:record.listLoadFailed')))
        .finally(() => {
          this.processing = false
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
        }
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
        currentListFilters = getItem(`record-list-filters-${this.uniqueID}`)
        currentListFilters = this.recordListFilter
        setItem(`record-list-filters-${this.uniqueID}`, currentListFilters)
      } catch (e) {
        console.warning(this.$t('notification:record-list.corrupted-filter'))
      }
    },

    onImportSuccessful () {
      this.refresh(true)
    },

    setDrillDownFilter (drillDownFilter) {
      if (drillDownFilter) {
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

    isDeleteActionVisible ({ canDeleteRecord }) {
      return canDeleteRecord
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
      ].some(v => v)
    },

    onBulkUpdate () {
      this.refresh(true)
    },

    editInlineField (record, field) {
      this.inlineEdit.fields = [field]
      this.inlineEdit.record = record.clone()
      this.inlineEdit.recordIDs = [record.recordID]
    },

    onInlineEdit () {
      this.refresh(true)
      this.inlineEdit.fields = []
      this.inlineEdit.recordIDs = []
      this.inlineEdit.record = {}
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

    updateFilter (filter, name) {
      const lastFilterIdx = this.recordListFilter.length - 1
      filter = filter.map((filter) => ({ ...filter, name }))

      if (this.recordListFilter.length) {
        this.recordListFilter[lastFilterIdx].groupCondition = 'AND'
      }

      this.recordListFilter = this.recordListFilter.concat(filter)
      this.activeFilters.push(name)
      this.refresh(true)
    },

    removeFilter (currentFilters) {
      if (this.drillDownFilter && !currentFilters.includes(this.$t('recordList.drillDown.filter.label'))) {
        this.setDrillDownFilter(undefined)
      }

      this.recordListFilter = this.recordListFilter.filter(({ name }) => !name || currentFilters.includes(name))
      this.refresh(true)
    },

    isUserRoleMember (roles) {
      if (!roles.length) return true

      return roles.some(roleID => this.authUserRoles.includes(roleID))
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
  color: $primary;
  vertical-align: sub;
  margin-left: 2px;
  width: 10px;
  height: 16px;
  overflow: hidden;
}

tr:hover td.actions {
  opacity: 1;
  background-color: $gray-200;
}

.inline-actions {
  min-width: 30px;
  margin-top: -2px;
  opacity: 0;
  transition: opacity 0.25s;
}

td:hover .inline-actions {
  opacity: 1;
  background-color: $gray-200;

  button:hover {
    color: $primary !important;
  }
}
</style>

<style lang="scss">
.record-list-table .actions {
  padding-top: 8px;
  position: sticky;
  right: 0;
  opacity: 0;
  transition: opacity 0.25s;
  width: 1%;

  .regular-font {
    font-family: $font-regular !important;
  }
}
</style>
