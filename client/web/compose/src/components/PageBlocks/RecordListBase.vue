<template>
  <wrap
    v-if="recordListModule"
    v-bind="$props"
    :scrollable-body="false"
    v-on="$listeners"
  >
    <template
      v-if="showHeader"
      #header
    >
      <h5
        class="d-flex align-items-center text-truncate mb-0"
      >
        {{ block.title }}
        <b-badge
          v-if="isFederated"
          variant="primary"
          class="d-inline-block mb-0 ml-2"
        >
          {{ $t('recordList.federated') }}
        </b-badge>
      </h5>

      <b-card-text
        v-if="block.description"
        class="text-dark text-truncate mt-1"
      >
        {{ block.description }}
      </b-card-text>
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
                class="mr-1 float-left"
                @export="onExport"
              />

              <column-picker
                v-if="options.allRecords"
                :module="recordListModule"
                :fields="fields"
                class="float-left"
                @updateFields="onUpdateFields"
              />
            </div>
          </div>
          <div
            v-if="!options.hideSearch && !inlineEditing"
            class="flex-grow-1 w-25"
          >
            <c-input-search
              v-model.trim="query"
              :placeholder="$t('general.label.search')"
            />
          </div>
        </b-row>
        <b-row
          v-if="options.selectable"
          v-show="selected.length > 0"
          class="mt-2 no-gutters"
        >
          <b-col
            cols="4"
            class="pt-1 text-nowrap font-weight-bold"
          >
            {{ $t('recordList.selected', { count: selected.length, total: items.length }) }}
            <a
              href="#"
              @click.prevent="handleSelectAllOnPage({ isChecked: false })"
            >
              ({{ $t('recordList.cancelSelection') }})
            </a>
          </b-col>
          <b-col
            class="text-right"
            cols="8"
          >
            <automation-buttons
              class="d-inline m-0"
              :buttons="options.selectionButtons"
              :module="recordListModule"
              :extra-event-args="{ selected, filter }"
              v-bind="$props"
              @refresh="refresh()"
            />
            <span v-if="canDeleteSelectedRecords">
              <c-input-confirm
                v-if="!inlineEditing"
                class="ml-2"
                variant="link-light"
                @confirmed="handleDeleteSelectedRecords()"
              />
              <b-button
                v-else-if="!areAllRowsDeleted"
                variant="link"
                size="md"
                @click.prevent="handleDeleteSelectedRecords()"
              >
                <font-awesome-icon
                  class="text-danger"
                  :icon="['far', 'trash-alt']"
                />
              </b-button>
              <b-button
                v-else
                variant="link"
                size="md"
                class="text-danger"
                @click.prevent="handleRestoreSelectedRecords()"
              >
                <font-awesome-icon
                  :icon="['fa', 'trash-restore']"
                />
              </b-button>
            </span>
          </b-col>
        </b-row>
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
          class="border-top mh-100 h-100 mb-0"
        >
          <b-thead>
            <b-tr>
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
                v-for="field in fields"
                :key="field.key"
                sticky-column
                class="pr-0"
                :style="{
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
                      v-if="field.filterable"
                      class="d-print-none"
                      :target="block.blockID"
                      :selected-field="field.moduleField"
                      :namespace="namespace"
                      :module="recordListModule"
                      :record-list-filter="recordListFilter"
                      @filter="onFilter"
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

              <b-th />
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
              :variant="!!item.r.deletedAt ? 'danger' : undefined"
              :class="{ 'pointer': !(options.editable && editing) }"
              @click="handleRowClicked(item)"
            >
              <b-td
                v-if="options.draggable && inlineEditing"
                class="align-middle pr-0"
                @click.stop
              >
                <font-awesome-icon
                  v-b-tooltip.hover
                  :icon="['fas', 'bars']"
                  :title="$t('general.tooltip.dragAndDrop')"
                  class="handle text-light"
                />
              </b-td>

              <b-td
                v-if="options.selectable"
                class="align-middle pr-0 d-print-none"
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
                  class="mb-0"
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
                  />
                </div>
                <i
                  v-else
                  class="text-primary"
                >
                  {{ $t('field.noPermission') }}
                </i>
              </b-td>

              <b-td
                class="d-flex justify-content-end align-items-top"
                @click.stop
              >
                <template v-if="inlineEditing">
                  <b-button
                    v-if="item.r.deletedAt"
                    variant="link"
                    size="md"
                    :title="$t('recordList.record.tooltip.undelete')"
                    class="border-0 text-dark mt-1 d-print-none"
                    @click.prevent="handleRestoreInline(item, index)"
                  >
                    <font-awesome-icon
                      :icon="['fa', 'trash-restore']"
                    />
                  </b-button>
                  <!-- The user should be able to delete the record if it's not yet saved -->
                  <b-button
                    v-else-if="item.r.canDeleteRecord && !item.r.deletedAt"
                    variant="link"
                    size="md"
                    class="border-0 show-when-hovered text-danger mt-1 d-print-none"
                    @click.prevent="handleDeleteInline(item, index)"
                  >
                    <font-awesome-icon
                      :icon="['far', 'trash-alt']"
                    />
                  </b-button>
                </template>

                <b-button
                  v-if="!inlineEditing && !options.hideRecordReminderButton"
                  variant="link"
                  :title="$t('recordList.record.tooltip.reminder')"
                  class="p-0 m-0 ml-2 text-primary d-print-none"
                  @click.prevent="createReminder(item.r)"
                >
                  <font-awesome-icon
                    :icon="['far', 'bell']"
                  />
                </b-button>

                <template v-if="!options.hideRecordCloneButton && recordListModule.canCreateRecord && (options.rowCreateUrl || recordPageID)">
                  <b-button
                    v-if="!inlineEditing"
                    variant="link"
                    :title="$t('recordList.record.tooltip.clone')"
                    class="p-0 m-0 ml-2 text-secondary d-print-none"
                    :to="{ name: options.rowCreateUrl || 'page.record.create', params: { pageID: recordPageID, values: item.r.values }, query: null }"
                  >
                    <font-awesome-icon
                      :icon="['far', 'clone']"
                    />
                  </b-button>
                  <b-button
                    v-else
                    variant="link"
                    :title="$t('recordList.record.tooltip.clone')"
                    class="ml-2 text-primary d-print-none"
                    @click="handleCloneInline(item.r)"
                  >
                    <font-awesome-icon
                      :icon="['far', 'clone']"
                    />
                  </b-button>
                </template>

                <template v-if="!inlineEditing">
                  <b-button
                    v-if="!options.hideRecordEditButton && item.r.canUpdateRecord && (options.rowEditUrl || recordPageID)"
                    variant="link"
                    :title="$t('recordList.record.tooltip.edit')"
                    class="p-0 m-0 ml-2 text-primary d-print-none"
                    :to="{ name: options.rowEditUrl || 'page.record.edit', params: { pageID: recordPageID, recordID: item.r.recordID }, query: null }"
                  >
                    <font-awesome-icon
                      :icon="['far', 'edit']"
                    />
                  </b-button>
                  <b-button
                    v-if="!options.hideRecordViewButton && item.r.canReadRecord && (options.rowViewUrl || recordPageID)"
                    variant="link"
                    :title="$t('recordList.record.tooltip.view')"
                    class="p-0 m-0 ml-2 text-primary d-print-none"
                    :to="{ name: options.rowViewUrl || 'page.record', params: { pageID: recordPageID, recordID: item.r.recordID }, query: null }"
                  >
                    <font-awesome-icon
                      :icon="['far', 'eye']"
                    />
                  </b-button>

                  <c-permissions-button
                    v-if="item.r.canGrant && !options.hideRecordPermissionsButton"
                    :resource="`corteza::compose:record/${item.r.namespaceID}/${item.r.moduleID}/${item.r.recordID}`"
                    :target="item.r.recordID"
                    :title="item.r.recordID"
                    :tooltip="$t('permissions:resources.compose.record.tooltip')"
                    button-variant="link"
                    class="text-dark m-0 p-0 ml-2 d-print-none"
                  />
                </template>
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
    </template>

    <template
      v-if="showFooter"
      #footer
    >
      <b-container
        ref="footer"
        fluid
        class="m-0 p-2"
      >
        <b-row no-gutters>
          <b-col class="d-flex justify-content-between align-items-center">
            <div class="text-truncate">
              <div
                v-if="options.showTotalCount"
                class="ml-2 text-nowrap font-weight-bold"
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
              v-if="!options.hidePaging && !inlineEditing"
            >
              <b-pagination
                v-if="options.fullPageNavigation"
                data-test-id="pagination"
                align="right"
                aria-controls="record-list"
                class="m-0 d-print-none"
                pills
                variant="link"
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
                  :disabled="!hasPrevPage"
                  data-test-id="first-page"
                  variant="link"
                  class="text-dark"
                  @click="goToPage()"
                >
                  <font-awesome-icon :icon="['fas', 'angle-double-left']" />
                </b-button>
                <b-button
                  :disabled="!hasPrevPage"
                  data-test-id="previous-page"
                  variant="link"
                  class="text-dark"
                  @click="goToPage('prevPage')"
                >
                  <font-awesome-icon :icon="['fas', 'angle-left']" />
                  {{ $t('recordList.pagination.prev') }}
                </b-button>
                <b-button
                  :disabled="!hasNextPage"
                  data-test-id="next-page"
                  variant="link"
                  class="text-dark"
                  @click="goToPage('nextPage')"
                >
                  {{ $t('recordList.pagination.next') }}
                  <font-awesome-icon :icon="['fas', 'angle-right']" />
                </b-button>
              </b-button-group>
            </div>
          </b-col>
        </b-row>
      </b-container>
    </template>
  </wrap>
</template>
<script>
import { debounce } from 'lodash'
import { mapGetters } from 'vuex'
import base from './base'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import FieldEditor from 'corteza-webapp-compose/src/components/ModuleFields/Editor'
import ExporterModal from 'corteza-webapp-compose/src/components/Public/Record/Exporter'
import ImporterModal from 'corteza-webapp-compose/src/components/Public/Record/Importer'
import AutomationButtons from './Shared/AutomationButtons'
import { compose, validator, NoID } from '@cortezaproject/corteza-js'
import users from 'corteza-webapp-compose/src/mixins/users'
import { evaluatePrefilter, queryToFilter } from 'corteza-webapp-compose/src/lib/record-filter'
import { getItem, setItem, removeItem } from 'corteza-webapp-compose/src/lib/local-storage'
import { components, url } from '@cortezaproject/corteza-vue'
import draggable from 'vuedraggable'
import RecordListFilter from 'corteza-webapp-compose/src/components/Common/RecordListFilter'
import ColumnPicker from 'corteza-webapp-compose/src/components/Admin/Module/Records/ColumnPicker'
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
  },

  extends: base,

  mixins: [
    users,
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
      prefilter: null,

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

      sortBy: undefined,
      sortDirecton: undefined,

      // This counter helps us generate unique ID's for the lifetime of this
      // component
      ctr: 0,
      items: [],
      idPrefix: `rl:${this.blockIndex}`,
      recordListFilter: [],
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      pages: 'page/set',
    }),

    loaderCollSpan () {
      // 2 for right side actions
      let base = this.fields.length + 2
      if (this.options.draggable && this.inlineEditing) base++
      if (this.options.selectable) base++
      return base
    },

    isFederated () {
      return Object.keys(this.recordListModule.labels || {}).includes('federation')
    },

    showHeader () {
      return !!(this.block.title || this.block.description || this.isFederated)
    },

    showFooter () {
      return !this.options.hidePaging && !this.inlineEditing
    },

    hasRightActions () {
      return this.editing
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
      handler (recordID = NoID) {
        const { pageID = NoID } = this.page

        // Set uniqueID so that events dont mix
        if (this.uniqueID) {
          this.$root.$off(`record-line:collect:${this.uniqueID}`)
          this.$root.$off(`page-block:validate:${this.uniqueID}`)
          this.$root.$off(`refetch-non-record-blocks:${pageID}`)
        }

        this.uniqueID = `${pageID}-${recordID}-${this.blockIndex}`
        this.$root.$on(`record-line:collect:${this.uniqueID}`, this.resolveRecords)
        this.$root.$on(`page-block:validate:${this.uniqueID}`, this.validatePageBlock)
        this.$root.$on(`refetch-non-record-blocks:${pageID}`, () => {
          this.refresh(true)
        })

        this.getStorageRecordListFilter()
        this.prepRecordList()
        this.refresh(true)
      },
    },
  },

  beforeDestroy () {
    this.$root.$off(`record-line:collect:${this.uniqueID}`)
    this.$root.$off(`page-block:validate:${this.uniqueID}`)
    this.$root.$off(`refetch-non-record-blocks:${this.page.pageID}`)
  },

  methods: {
    onFilter (filter = []) {
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
        id: id || (r.recordID !== NoID ? r.recordID : `${this.idPrefix}:${this.ctr++}`),
        _rowVariant: r.deletedAt ? 'danger' : undefined,
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
        idPrefix: this.idPrefix,
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
      const r = new compose.Record(this.recordListModule, { ...item.r, deletedAt: new Date() })
      this.items.splice(i, 1, this.wrapRecord(r, item.id))
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
      const limit = editable ? 0 : perPage
      this.filter = {
        limit,
        sort,
      }
    },

    createReminder (record) {
      // Determine initial reminder title
      const { recordID, values = {} } = record
      const tField = ((this.options.fields || []).find(({ name }) => !!values[name]) || {}).name
      const title = values[tField]

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

      const pageID = this.recordPageID
      const route = {
        name: this.options.rowViewUrl || 'page.record',
        params: {
          pageID,
          recordID,
        },
        query: null,
      }

      if (this.options.openInNewTab) {
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
      }
    },

    handleRestoreSelectedRecords () {
      const sel = new Set(this.selected)
      for (let i = 0; i < this.items.length; i++) {
        if (sel.has(this.items[i].id)) {
          this.handleRestoreInline(this.items[i], i)
        }
      }
    },

    handleDeleteSelectedRecords () {
      if (this.selected.length === 0) {
        return
      }

      if (this.inlineEditing) {
        const sel = new Set(this.selected)
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
          .filter(({ id, r }) => r.canDeleteRecord && this.selected.includes(id))
          .map(({ id }) => id)

        this.processing = true

        this.$ComposeAPI
          .recordBulkDelete({ moduleID, namespaceID, recordIDs })
          .then(() => {
            this.refresh(true)
            this.toastSuccess(this.$t('notification:record.deleteBulkSuccess'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:record.deleteBulkFailed')))
          .finally(() => {
            this.processing = false
          })
      }
    },

    refresh (resetPagination = false) {
      this.pullRecords(resetPagination)
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
      const query = queryToFilter(this.query, this.prefilter, this.recordListModule.filterFields(this.options.fields), this.recordListFilter)

      const { moduleID, namespaceID } = this.recordListModule
      if (this.filter.pageCursor) {
        this.filter.sort = ''
      }

      let paginationOptions = {}
      if (resetPagination) {
        const { fullPageNavigation = false, showTotalCount = false } = this.options
        paginationOptions = {
          incPageNavigation: fullPageNavigation,
          incTotal: showTotalCount,
        }
      }

      await this.$ComposeAPI.recordList({ ...this.filter, moduleID, namespaceID, query, ...paginationOptions })
        .then(({ set, filter }) => {
          const records = set.map(r => new compose.Record(r, this.recordListModule))

          this.filter = { ...this.filter, ...filter }
          this.filter.pageCursor = undefined
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
          this.fetchUsers(fields, records)

          this.items = records.map(r => this.wrapRecord(r))
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
</style>
