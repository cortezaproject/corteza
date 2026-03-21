<template>
  <wrap
    v-bind="$props"
    class="collapsible-record-list-wrap"
    body-class="pt-0 px-0"
    v-on="$listeners"
    @refreshBlock="refresh(true, false)"
  >
    <!-- Loading state -->
    <div
      v-if="isProcessing"
      class="d-flex align-items-center justify-content-center h-100 p-3"
    >
      <b-spinner />
    </div>

    <div
      v-else-if="recordListModule"
      class="collapsible-record-list"
    >
      <!-- Iterate over each record and display in collapsible section -->
      <div
        v-for="(item, index) in items"
        :key="item.r.recordID || index"
        class="collapsible-record-item"
      >
        <!-- Collapsible Header -->
        <div
          class="collapsible-header d-flex align-items-center p-3"
          :class="getHeaderClass(item.r)"
          role="button"
          :aria-expanded="!isCollapsed(item.r.recordID)"
          :aria-controls="`collapsible-body-${blockIndex}-${item.r.recordID}`"
          tabindex="0"
          @click="toggleCollapse(item.r.recordID)"
          @keydown.enter="toggleCollapse(item.r.recordID)"
          @keydown.space.prevent="toggleCollapse(item.r.recordID)"
        >
          <!-- Title and Subtitle Row -->
          <div class="collapsible-titles">
            <!-- Title with inline alignment OR regular alignment -->
            <template v-if="useInlineTitleAlignment(item.r)">
              <div class="collapsible-title-inline d-flex justify-content-between w-100">
                <span
                  v-if="getTitleAlignmentData(item.r).left"
                  class="text-left"
                  style="flex: 1;"
                  v-html="getTitleAlignmentData(item.r).left"
                />
                <span
                  v-if="getTitleAlignmentData(item.r).center"
                  class="text-center"
                  style="flex: 1;"
                  v-html="getTitleAlignmentData(item.r).center"
                />
                <span
                  v-if="getTitleAlignmentData(item.r).right"
                  class="text-right"
                  style="flex: 1;"
                  v-html="getTitleAlignmentData(item.r).right"
                />
              </div>
            </template>
            <template v-else>
              <div
                class="collapsible-title"
                :class="getTitleAlignmentClass(item.r)"
              >
                {{ getTitleText(item.r) || options.titleExpression || 'Untitled' }}
              </div>
            </template>

            <!-- Subtitle -->
            <template v-if="displaySubtitle(item.r)">
              <template v-if="useInlineSubtitleAlignment(item.r)">
                <div class="collapsible-subtitle-inline d-flex justify-content-between w-100 mt-1">
                  <span
                    v-if="getSubtitleAlignmentData(item.r).left"
                    class="text-left"
                    style="flex: 1;"
                    v-html="getSubtitleAlignmentData(item.r).left"
                  />
                  <span
                    v-if="getSubtitleAlignmentData(item.r).center"
                    class="text-center"
                    style="flex: 1;"
                    v-html="getSubtitleAlignmentData(item.r).center"
                  />
                  <span
                    v-if="getSubtitleAlignmentData(item.r).right"
                    class="text-right"
                    style="flex: 1;"
                    v-html="getSubtitleAlignmentData(item.r).right"
                  />
                </div>
              </template>
              <template v-else>
                <div
                  class="collapsible-subtitle"
                  :class="getSubtitleAlignmentClass(item.r)"
                >
                  {{ getSubtitleText(item.r) || options.subtitleExpression || '' }}
                </div>
              </template>
            </template>
          </div>

          <!-- Chevron -->
          <div class="collapsible-chevron ml-2">
            <font-awesome-icon
              :icon="isCollapsed(item.r.recordID) ? ['fas', 'chevron-down'] : ['fas', 'chevron-up']"
              class="text-muted"
            />
          </div>
        </div>

        <!-- Collapsible Body -->
        <b-collapse
          :id="`collapsible-body-${blockIndex}-${item.r.recordID}`"
          :visible="!isCollapsed(item.r.recordID)"
          class="collapsible-content"
        >
          <!-- Body Field (Full Width) -->
          <div
            v-if="bodyField && bodyFieldConfig"
            class="body-field-container p-3"
          >
            <div
              v-if="bodyFieldConfig.canReadRecordValue"
              class="body-field-content"
              v-html="getBodyFieldValue(item.r)"
            />
            <i
              v-else
              class="text-muted"
            >
              {{ $t('field.noPermission') }}
            </i>
          </div>

          <!-- Other Fields Section - Above Body -->
          <div
            v-if="options.otherFieldsPosition === 'above'"
            class="other-fields-section px-3 pb-3"
          >
            <div :class="fieldLayoutClass">
              <template v-for="field in otherFields">
                <b-form-group
                  v-if="canDisplay(field, item.r)"
                  :key="`${field.fieldID}-${field.name}-${item.r.recordID}`"
                  :label-cols-md="options.horizontalFieldLayoutEnabled && '6'"
                  :label-cols-xl="options.horizontalFieldLayoutEnabled && '5'"
                  :content-cols-md="options.horizontalFieldLayoutEnabled && '6'"
                  :content-cols-xl="options.horizontalFieldLayoutEnabled && '7'"
                  :class="columnWrapClass"
                  :style="fieldWidth"
                  class="field-container"
                >
                  <template #label>
                    <div class="d-flex align-items-center text-primary mb-0">
                      <span
                        class="d-flex"
                        style="margin-top: 0.1rem;"
                      >
                        {{ field.label || field.name }}
                      </span>
                      <c-hint :tooltip="((field.options.hint || {}).view || '')" />
                    </div>
                    <div
                      class="small text-muted"
                      :class="{ 'mb-1': !!(field.options.description || {}).view }"
                    >
                      {{ (field.options.description || {}).view }}
                    </div>
                  </template>
                  <div
                    v-if="field.canReadRecordValue"
                    class="value align-self-center"
                  >
                    <field-viewer
                      v-bind="{ ...$props, record: item.r }"
                      :module="recordListModule"
                      :field="field"
                      :extra-options="options"
                    />
                  </div>
                  <i
                    v-else
                    class="text-muted"
                  >{{ $t('field.noPermission') }}</i>
                </b-form-group>
              </template>
            </div>
          </div>

          <!-- Default Position (after body) -->
          <div
            v-if="options.otherFieldsPosition === 'default' || !options.otherFieldsPosition"
            class="other-fields-section px-3 pb-3"
          >
            <div :class="fieldLayoutClass">
              <template v-for="field in otherFields">
                <b-form-group
                  v-if="canDisplay(field, item.r)"
                  :key="`${field.fieldID}-${field.name}-${item.r.recordID}`"
                  :label-cols-md="options.horizontalFieldLayoutEnabled && '6'"
                  :label-cols-xl="options.horizontalFieldLayoutEnabled && '5'"
                  :content-cols-md="options.horizontalFieldLayoutEnabled && '6'"
                  :content-cols-xl="options.horizontalFieldLayoutEnabled && '7'"
                  :class="columnWrapClass"
                  :style="fieldWidth"
                  class="field-container"
                >
                  <template #label>
                    <div class="d-flex align-items-center text-primary mb-0">
                      <span
                        class="d-flex"
                        style="margin-top: 0.1rem;"
                      >
                        {{ field.label || field.name }}
                      </span>
                      <c-hint :tooltip="((field.options.hint || {}).view || '')" />
                    </div>
                    <div
                      class="small text-muted"
                      :class="{ 'mb-1': !!(field.options.description || {}).view }"
                    >
                      {{ (field.options.description || {}).view }}
                    </div>
                  </template>
                  <div
                    v-if="field.canReadRecordValue"
                    class="value align-self-center"
                  >
                    <field-viewer
                      v-bind="{ ...$props, record: item.r }"
                      :module="recordListModule"
                      :field="field"
                      :extra-options="options"
                    />
                  </div>
                  <i
                    v-else
                    class="text-muted"
                  >{{ $t('field.noPermission') }}</i>
                </b-form-group>
              </template>
            </div>
          </div>

          <!-- Below Body -->
          <div
            v-if="options.otherFieldsPosition === 'below'"
            class="other-fields-section px-3 pb-3"
          >
            <div :class="fieldLayoutClass">
              <template v-for="field in otherFields">
                <b-form-group
                  v-if="canDisplay(field, item.r)"
                  :key="`${field.fieldID}-${field.name}-${item.r.recordID}`"
                  :label-cols-md="options.horizontalFieldLayoutEnabled && '6'"
                  :label-cols-xl="options.horizontalFieldLayoutEnabled && '5'"
                  :content-cols-md="options.horizontalFieldLayoutEnabled && '6'"
                  :content-cols-xl="options.horizontalFieldLayoutEnabled && '7'"
                  :class="columnWrapClass"
                  :style="fieldWidth"
                  class="field-container"
                >
                  <template #label>
                    <div class="d-flex align-items-center text-primary mb-0">
                      <span
                        class="d-flex"
                        style="margin-top: 0.1rem;"
                      >
                        {{ field.label || field.name }}
                      </span>
                      <c-hint :tooltip="((field.options.hint || {}).view || '')" />
                    </div>
                    <div
                      class="small text-muted"
                      :class="{ 'mb-1': !!(field.options.description || {}).view }"
                    >
                      {{ (field.options.description || {}).view }}
                    </div>
                  </template>
                  <div
                    v-if="field.canReadRecordValue"
                    class="value align-self-center"
                  >
                    <field-viewer
                      v-bind="{ ...$props, record: item.r }"
                      :module="recordListModule"
                      :field="field"
                      :extra-options="options"
                    />
                  </div>
                  <i
                    v-else
                    class="text-muted"
                  >{{ $t('field.noPermission') }}</i>
                </b-form-group>
              </template>
            </div>
          </div>
        </b-collapse>
      </div>

      <!-- Empty state -->
      <div
        v-if="!items.length && !isProcessing"
        class="text-center p-3 text-muted"
      >
        {{ $t('recordList.noRecords') }}
      </div>
    </div>

    <template
      v-if="recordListModule"
    >
      <!-- Pagination Footer -->
      <div
        v-if="showFooter"
        class="d-flex align-items-center justify-content-between px-3 py-2 border-top mt-2"
      >
        <div
          v-if="options.showTotalCount"
          class="text-muted small"
        >
          {{ pagination.count }} records
        </div>
        <b-button-group class="gap-1 ml-auto">
          <b-button
            :disabled="!hasPrevPage || isProcessing"
            variant="outline-extra-light"
            class="d-flex align-items-center text-dark border-0 p-1"
            @click="goToPage('first')"
          >
            <font-awesome-icon :icon="['fas', 'angle-double-left']" />
          </b-button>
          <b-button
            :disabled="!hasPrevPage || isProcessing"
            variant="outline-extra-light"
            class="d-flex align-items-center text-dark border-0 p-1"
            @click="goToPage('prev')"
          >
            <font-awesome-icon
              :icon="['fas', 'angle-left']"
              class="mr-1"
            />
            {{ $t('recordList.pagination.prev') }}
          </b-button>
          <b-button
            :disabled="!hasNextPage || isProcessing"
            variant="outline-extra-light"
            class="d-flex align-items-center text-dark border-0 p-1"
            @click="goToPage('next')"
          >
            {{ $t('recordList.pagination.next') }}
            <font-awesome-icon
              :icon="['fas', 'angle-right']"
              class="ml-1"
            />
          </b-button>
        </b-button-group>
      </div>
    </template>
  </wrap>
</template>

<script>
import { NoID, compose } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'
import base from './base'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import users from 'corteza-webapp-compose/src/mixins/users'
import records from 'corteza-webapp-compose/src/mixins/records'
import conditionalFields from 'corteza-webapp-compose/src/mixins/conditionalFields'
import recordLayout from 'corteza-webapp-compose/src/mixins/recordLayout'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldViewer,
  },

  extends: base,

  mixins: [
    users,
    records,
    conditionalFields,
    recordLayout,
  ],

  data () {
    return {
      // Collapse state per record (stored locally in component)
      collapsedRecords: {},

      // Record list data
      processing: false,
      items: [],

      // Filter and pagination
      prefilter: undefined,
      filter: {
        query: '',
        sort: 'createdAt DESC',
        limit: 10,
        pageCursor: '',
      },

      pagination: {
        count: 0,
      },

      // Track if records are loaded
      recordsLoaded: false,
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
    }),

    // Returns the module configured for this record list
    recordListModule () {
      if (this.options.moduleID) {
        return this.getModuleByID(this.options.moduleID)
      } else {
        return undefined
      }
    },

    isProcessing () {
      return this.processing
    },

    fieldLayoutClass () {
      const classes = {
        default: 'd-flex flex-column',
        noWrap: 'd-flex gap-2',
        wrap: 'row no-gutters',
      }
      return classes[this.options.recordFieldLayoutOption] || classes.default
    },

    columnWrapClass () {
      if (this.options.recordFieldLayoutOption !== 'wrap') {
        return ''
      }
      return 'field-col'
    },

    fieldWidth () {
      if (this.options.recordFieldLayoutOption !== 'noWrap') {
        return {}
      }
      return { 'min-width': '13rem' }
    },

    bodyFieldConfig () {
      if (!this.options.bodyField || !this.recordListModule) {
        return null
      }
      const field = this.recordListModule.fields.find(f => f.name === this.options.bodyField || f.fieldID === this.options.bodyField)
      if (!field) {
        return null
      }
      return field
    },

    bodyField () {
      return this.options.bodyField
    },

    otherFields () {
      if (!this.recordListModule) {
        return []
      }

      const bodyFieldName = this.options.bodyField
      const selectedFields = this.options.fields || []

      let fields = []

      if (selectedFields.length === 0) {
        // No fields selected, show all except body field
        fields = this.recordListModule.fields.filter(f => f.name !== bodyFieldName && f.fieldID !== bodyFieldName)
      } else {
        // Show selected fields except body field
        fields = this.recordListModule.filterFields(selectedFields).filter(f => f.name !== bodyFieldName && f.fieldID !== bodyFieldName)
      }

      return fields.map(f => {
        f.label = f.isSystem ? this.$t(`field:system.${f.name}`) : f.label || f.name
        return f
      })
    },

    hasPrevPage () {
      return !!this.filter.prevPage
    },

    hasNextPage () {
      return !!this.filter.nextPage
    },

    showPageNavigation () {
      return !this.options.hidePaging
    },

    showFooter () {
      return this.recordListModule && !this.options.hidePaging
    },
  },

  watch: {
    'options.moduleID': {
      immediate: true,
      handler () {
        if (this.options.moduleID && !this.recordsLoaded) {
          this.loadRecords()
        }
      },
    },
  },

  created () {
    if (this.options.moduleID) {
      this.loadRecords()
    }
  },

  beforeDestroy () {
    this.collapsedRecords = {}
  },

  methods: {
    // Record list loading methods
    async loadRecords () {
      if (!this.recordListModule) {
        return
      }

      this.processing = true

      try {
        // Build filter
        const filter = this.buildFilter()

        const { moduleID, namespaceID } = this.recordListModule

        const { set: records, filter: responseFilter } = await this.$ComposeAPI.recordList({
          moduleID,
          namespaceID,
          ...filter,
          pageCursor: this.filter.pageCursor,
        })

        // Store pagination cursors
        this.filter.nextPage = responseFilter.nextPage
        this.filter.prevPage = responseFilter.prevPage
        this.pagination.count = responseFilter.total || records.length

        // Convert to Record objects
        this.items = records.map(r => ({
          r: new compose.Record(r, this.recordListModule),
          id: r.recordID,
        }))

        this.recordsLoaded = true

        // Initialize collapse state for each record
        this.items.forEach(item => {
          if (this.collapsedRecords[item.r.recordID] === undefined) {
            this.$set(this.collapsedRecords, item.r.recordID, true)
          }
        })
      } catch (e) {
        console.error('Failed to load records:', e)
      } finally {
        this.processing = false
      }
    },

    buildFilter () {
      const { presort, prefilter, perPage, refField } = this.options

      const sort = presort || 'createdAt DESC'
      const filterExpressions = []

      // Initial prefilter
      if (prefilter) {
        try {
          const pf = evaluatePrefilter(prefilter, {
            record: this.record,
            user: this.$auth.user || {},
            recordID: (this.record || {}).recordID || NoID,
            ownerID: (this.record || {}).ownedBy || NoID,
            userID: (this.$auth.user || {}).userID || NoID,
          })
          filterExpressions.push(`(${pf})`)
        } catch (e) {
          console.warn('Failed to evaluate prefilter:', e)
        }
      }

      // Link to parent record via refField
      if (refField && this.record && this.record.recordID && this.record.recordID !== NoID) {
        const refFieldObj = this.recordListModule.fields.find(f => f.name === refField)
        if (refFieldObj && refFieldObj.isMulti) {
          filterExpressions.push(`${refField} IN ('${this.record.recordID}')`)
        } else {
          filterExpressions.push(`${refField} = '${this.record.recordID}'`)
        }
      }

      return {
        sort,
        limit: perPage || 10,
        filter: filterExpressions.join(' AND '),
      }
    },

    async refresh (resetPagination = false) {
      if (resetPagination) {
        this.filter.pageCursor = undefined
      }
      await this.loadRecords()
    },

    async goToPage (direction) {
      if (direction === 'next') {
        this.filter.pageCursor = this.filter.nextPage
      } else if (direction === 'prev') {
        this.filter.pageCursor = this.filter.prevPage
      } else {
        this.filter.pageCursor = undefined
      }
      await this.loadRecords()
    },

    // Collapse/expand methods
    toggleCollapse (recordID) {
      this.$set(this.collapsedRecords, recordID, !this.collapsedRecords[recordID])
    },

    isCollapsed (recordID) {
      return this.collapsedRecords[recordID] !== false
    },

    getHeaderClass (record) {
      return 'cursor-pointer'
    },

    // Per-record title/subtitle methods
    getTitleAlignmentData (record) {
      return this.parseAlignmentExpression(this.options.titleExpression || '', record)
    },

    getSubtitleAlignmentData (record) {
      return this.parseAlignmentExpression(this.options.subtitleExpression || '', record)
    },

    useInlineTitleAlignment (record) {
      return this.getTitleAlignmentData(record).hasValidMarkers
    },

    useInlineSubtitleAlignment (record) {
      return this.getSubtitleAlignmentData(record).hasValidMarkers
    },

    getTitleText (record) {
      if (this.useInlineTitleAlignment(record)) {
        return ''
      }

      try {
        return evaluatePrefilter(this.options.titleExpression || '', {
          record,
          user: this.$auth.user || {},
          recordID: record.recordID || NoID,
          ownerID: record.ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      } catch (e) {
        return this.options.titleExpression || ''
      }
    },

    getSubtitleText (record) {
      if (this.useInlineSubtitleAlignment(record)) {
        return ''
      }

      try {
        return evaluatePrefilter(this.options.subtitleExpression || '', {
          record,
          user: this.$auth.user || {},
          recordID: record.recordID || NoID,
          ownerID: record.ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      } catch (e) {
        return this.options.subtitleExpression || ''
      }
    },

    displaySubtitle (record) {
      // If option to show when collapsed is NOT checked (false) and block is collapsed, hide subtitle
      if (!this.options.subtitleShowWhenCollapsed && this.isCollapsed(record.recordID)) {
        return false
      }
      // Check both regular text and inline alignment data
      if (this.useInlineSubtitleAlignment(record)) {
        const data = this.getSubtitleAlignmentData(record)
        return !!(data.left || data.center || data.right)
      }
      return !!this.getSubtitleText(record)
    },

    getTitleAlignmentClass (record) {
      return {
        'text-left': !this.useInlineTitleAlignment(record) || this.options.titleAlignment === 'left',
        'text-center': this.options.titleAlignment === 'center' && !this.useInlineTitleAlignment(record),
        'text-right': this.options.titleAlignment === 'right' && !this.useInlineTitleAlignment(record),
      }
    },

    getSubtitleAlignmentClass (record) {
      return {
        'text-left': !this.useInlineSubtitleAlignment(record) || this.options.subtitleAlignment === 'left',
        'text-center': this.options.subtitleAlignment === 'center' && !this.useInlineSubtitleAlignment(record),
        'text-right': this.options.subtitleAlignment === 'right' && !this.useInlineSubtitleAlignment(record),
      }
    },

    getBodyFieldValue (record) {
      if (!this.bodyFieldConfig || !record) {
        return ''
      }
      const fieldName = this.bodyFieldConfig.name
      return record.values[fieldName] || ''
    },

    canDisplay (field, record) {
      // Check conditional field visibility
      if (!this.options.fieldConditions || !this.options.fieldConditions.length) {
        return true
      }

      const conditions = this.options.fieldConditions.filter(c => c.field === field.name)
      if (!conditions.length) {
        return true
      }

      for (const condition of conditions) {
        if (condition.condition) {
          try {
            const result = evaluatePrefilter(condition.condition, {
              record,
              user: this.$auth.user || {},
              recordID: record.recordID || NoID,
              ownerID: record.ownedBy || NoID,
              userID: (this.$auth.user || {}).userID || NoID,
            })

            if (!result) {
              return false
            }
          } catch (e) {
            console.warn('Failed to evaluate condition:', e)
          }
        }
      }

      return true
    },

    /**
     * Parse alignment expression and return structured data
     * Supports: L[...] for left, C[...] for center, R[...] for right
     */
    parseAlignmentExpression (expression, record) {
      if (!expression) {
        return { hasValidMarkers: false, left: '', center: '', right: '' }
      }

      const ctx = {
        record: record || {},
        user: this.$auth.user || {},
        recordID: (record || {}).recordID || NoID,
        ownerID: (record || {}).ownedBy || NoID,
        userID: (this.$auth.user || {}).userID || NoID,
      }

      // Find all markers
      const leftMatch = expression.match(/^L\[([^\]]*)\]/)
      const centerMatch = expression.match(/C\[([^\]]*)\]/)
      const rightMatch = expression.match(/R\[([^\]]*)\]/)

      // Check if there are ANY markers
      const hasAnyMarkers = leftMatch || centerMatch || rightMatch
      if (!hasAnyMarkers) {
        return { hasValidMarkers: false, left: '', center: '', right: '' }
      }

      // Check for duplicate markers
      const allMatches = expression.match(/([LCR])\[/g) || []
      const uniqueMarkers = [...new Set(allMatches)]
      if (allMatches.length !== uniqueMarkers.length) {
        return { hasValidMarkers: false, left: '', center: '', right: '' }
      }

      // Validate order: L must come before C, C must come before R
      const lIndex = expression.indexOf('L[')
      const cIndex = expression.indexOf('C[')
      const rIndex = expression.indexOf('R[')

      const hasL = lIndex !== -1
      const hasC = cIndex !== -1
      const hasR = rIndex !== -1

      if (hasL && hasC && lIndex > cIndex) return { hasValidMarkers: false, left: '', center: '', right: '' }
      if (hasC && hasR && cIndex > rIndex) return { hasValidMarkers: false, left: '', center: '', right: '' }
      if (hasL && hasR && lIndex > rIndex) return { hasValidMarkers: false, left: '', center: '', right: '' }

      // Extract and evaluate content
      let left = ''
      let center = ''
      let right = ''

      if (leftMatch) {
        try {
          left = evaluatePrefilter(leftMatch[1], ctx) || ''
        } catch (e) {
          left = leftMatch[1]
        }
      }

      if (centerMatch) {
        try {
          center = evaluatePrefilter(centerMatch[1], ctx) || ''
        } catch (e) {
          center = centerMatch[1]
        }
      }

      if (rightMatch) {
        try {
          right = evaluatePrefilter(rightMatch[1], ctx) || ''
        } catch (e) {
          right = rightMatch[1]
        }
      }

      return {
        hasValidMarkers: true,
        left,
        center,
        right,
      }
    },
  },
}
</script>

<style scoped>
.collapsible-record {
  border: 1px solid #dee2e6;
  border-radius: 0.25rem;
}

.collapsible-header {
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
  cursor: pointer;
  user-select: none;
}

.collapsible-header:hover {
  background-color: #e9ecef;
}

.collapsible-titles {
  flex-grow: 1;
}

.collapsible-title {
  font-weight: 600;
  font-size: 1.1rem;
}

.collapsible-subtitle {
  font-size: 0.9rem;
  color: #6c757d;
}

.collapsible-chevron {
  flex-shrink: 0;
}

.body-field-container {
  border-bottom: 1px solid #dee2e6;
}

.field-col > * {
  margin-left: 1rem;
  margin-right: 1rem;
}
</style>
