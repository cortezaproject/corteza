<template>
  <wrap
    v-bind="$props"
    body-class="pt-0 px-0"
    class="collapsible-record-wrap"
    v-on="$listeners"
  >
    <div
      v-if="isProcessing"
      class="d-flex align-items-center justify-content-center h-100 p-3"
    >
      <b-spinner />
    </div>

    <div
      v-else-if="module"
      ref="collapsibleRecordContainer"
      class="collapsible-record"
    >
      <!-- Collapsible Header -->
      <div
        ref="collapsibleHeader"
        class="collapsible-header d-flex align-items-center p-3"
        :class="headerClass"
        role="button"
        :aria-expanded="!collapsed"
        :aria-controls="`collapsible-body-${blockIndex}`"
        tabindex="0"
        @click="toggleCollapse"
        @keydown.enter="toggleCollapse"
        @keydown.space.prevent="toggleCollapse"
      >
        <!-- Title and Subtitle Row - stacked vertically -->
        <div class="collapsible-titles">
          <!-- Title with inline alignment OR regular alignment -->
          <template v-if="useInlineTitleAlignment">
            <div class="collapsible-title d-flex justify-content-between w-100">
              <span
                v-if="titleAlignmentData.left"
                class="text-left"
                style="flex: 1;"
                v-html="titleAlignmentData.left"
              />
              <span
                v-if="titleAlignmentData.center"
                class="text-center"
                style="flex: 1;"
                v-html="titleAlignmentData.center"
              />
              <span
                v-if="titleAlignmentData.right"
                class="text-right"
                style="flex: 1;"
                v-html="titleAlignmentData.right"
              />
            </div>
          </template>
          <template v-else>
            <div
              class="collapsible-title"
              :class="titleAlignmentClass"
            >
              {{ titleText || options.titleExpression || 'Untitled' }}
            </div>
          </template>

          <!-- Subtitle (shown below title) -->
          <template v-if="displaySubtitle">
            <template v-if="useInlineSubtitleAlignment">
              <div class="collapsible-subtitle d-flex justify-content-between w-100 mt-1">
                <span
                  v-if="subtitleAlignmentData.left"
                  class="text-left"
                  style="flex: 1;"
                  v-html="subtitleAlignmentData.left"
                />
                <span
                  v-if="subtitleAlignmentData.center"
                  class="text-center"
                  style="flex: 1;"
                  v-html="subtitleAlignmentData.center"
                />
                <span
                  v-if="subtitleAlignmentData.right"
                  class="text-right"
                  style="flex: 1;"
                  v-html="subtitleAlignmentData.right"
                />
              </div>
            </template>
            <template v-else>
              <div
                class="collapsible-subtitle"
                :class="subtitleAlignmentClass"
              >
                {{ subtitleText || options.subtitleExpression || '' }}
              </div>
            </template>
          </template>
        </div>

        <!-- Chevron -->
        <div class="collapsible-chevron ml-2">
          <font-awesome-icon
            :icon="collapsed ? ['fas', 'chevron-down'] : ['fas', 'chevron-up']"
            class="text-muted"
          />
        </div>
      </div>

      <!-- Collapsible Body -->
      <b-collapse
        :id="`collapsible-body-${blockIndex}`"
        :visible="!collapsed"
        class="collapsible-content"
        :style="collapseMaxHeightStyle"
        @shown="handleCollapseShown"
        @hidden="handleCollapseHidden"
      >
        <!-- Other Fields Section - Above Body -->
        <div
          v-if="options.otherFieldsPosition === 'above' && otherFields.length"
          class="other-fields-section px-3 pb-3"
        >
          <div
            :class="fieldLayoutClass"
          >
            <template v-for="field in otherFields">
              <b-form-group
                v-if="canDisplay(field)"
                :key="`${field.fieldID}-${field.name}`"
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
                    v-bind="{ ...$props, field }"
                    :extra-options="options"
                  />
                </div>
                <i
                  v-else
                  class="text-muted"
                >
                  {{ $t('field.noPermission') }}
                </i>
              </b-form-group>
            </template>
          </div>
        </div>

        <!-- Body Field (Full Width - Rich Text) -->
        <div
          v-if="bodyField && bodyFieldConfig"
          class="body-field-container p-3"
        >
          <div
            v-if="bodyFieldConfig.canReadRecordValue"
            class="body-field-content"
            v-html="bodyFieldValue"
          />
          <i
            v-else
            class="text-muted"
          >
            {{ $t('field.noPermission') }}
          </i>
        </div>

        <!-- Default Position (after body) -->
        <div
          v-if="(options.otherFieldsPosition === 'default' || !options.otherFieldsPosition) && otherFields.length"
          class="other-fields-section px-3 pb-3"
        >
          <div
            :class="fieldLayoutClass"
          >
            <template v-for="field in otherFields">
              <b-form-group
                v-if="canDisplay(field)"
                :key="`${field.fieldID}-${field.name}`"
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
                    v-bind="{ ...$props, field }"
                    :extra-options="options"
                  />
                </div>
                <i
                  v-else
                  class="text-muted"
                >
                  {{ $t('field.noPermission') }}
                </i>
              </b-form-group>
            </template>
          </div>
        </div>

        <!-- Below Body -->
        <div
          v-if="options.otherFieldsPosition === 'below' && otherFields.length"
          class="other-fields-section px-3 pb-3"
        >
          <div
            :class="fieldLayoutClass"
          >
            <template v-for="field in otherFields">
              <b-form-group
                v-if="canDisplay(field)"
                :key="`${field.fieldID}-${field.name}`"
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
                    v-bind="{ ...$props, field }"
                    :extra-options="options"
                  />
                </div>
                <i
                  v-else
                  class="text-muted"
                >
                  {{ $t('field.noPermission') }}
                </i>
              </b-form-group>
            </template>
          </div>
        </div>
      </b-collapse>
    </div>
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
import alignment from 'corteza-webapp-compose/src/mixins/alignment'
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
    alignment,
  ],

  props: {
    onBlockHeightChange: {
      type: Function,
      default: null,
    },
  },

  data () {
    return {
      collapsed: false,
      originalHeight: null,
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      findRecordByID: 'record/findByID',
    }),

    resolvedRecord () {
      if (!this.record || !this.module) return this.record

      const resolvedValues = { ...this.record.values }

      this.module.fields
        .filter(f => f.kind === 'Record' && f.options && f.options.labelField)
        .forEach(field => {
          const rawValue = resolvedValues[field.name]
          if (!rawValue) return
          const rawResolved = this.findRecordByID(rawValue)
          if (!rawResolved) return
          const refModule = this.getModuleByID(field.options.moduleID)
          if (!refModule) return
          const refRecord = new compose.Record(refModule, rawResolved)
          resolvedValues[field.name] = refRecord.values[field.options.labelField] || rawValue
        })

      return { ...this.record, values: resolvedValues }
    },

    titleAlignmentData () {
      return this.parseAlignmentExpression(this.options.titleExpression || '', this.resolvedRecord)
    },

    subtitleAlignmentData () {
      return this.parseAlignmentExpression(this.options.subtitleExpression || '', this.resolvedRecord)
    },

    useInlineTitleAlignment () {
      return this.titleAlignmentData.hasValidMarkers
    },

    useInlineSubtitleAlignment () {
      return this.subtitleAlignmentData.hasValidMarkers
    },

    titleText () {
      // If using inline alignment, don't return plain text
      if (this.useInlineTitleAlignment) {
        return ''
      }

      try {
        return evaluatePrefilter(this.options.titleExpression || '', {
          record: this.resolvedRecord,
          user: this.$auth.user || {},
          recordID: (this.resolvedRecord || {}).recordID || NoID,
          ownerID: (this.resolvedRecord || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      } catch (e) {
        return this.options.titleExpression || ''
      }
    },

    subtitleText () {
      // If using inline alignment, don't return plain text
      if (this.useInlineSubtitleAlignment) {
        return ''
      }

      try {
        return evaluatePrefilter(this.options.subtitleExpression || '', {
          record: this.resolvedRecord,
          user: this.$auth.user || {},
          recordID: (this.resolvedRecord || {}).recordID || NoID,
          ownerID: (this.resolvedRecord || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      } catch (e) {
        return this.options.subtitleExpression || ''
      }
    },

    displaySubtitle () {
      // If option to show when collapsed is NOT checked (false) and block is collapsed, hide subtitle
      if (!this.options.subtitleShowWhenCollapsed && this.collapsed) {
        return false
      }
      // Check both regular text and inline alignment data
      if (this.useInlineSubtitleAlignment) {
        return !!(this.subtitleAlignmentData.left || this.subtitleAlignmentData.center || this.subtitleAlignmentData.right)
      }
      return !!this.subtitleText
    },

    titleAlignmentClass () {
      // Default to left alignment when no inline markers are present
      return {
        'text-left': !this.useInlineTitleAlignment || this.options.titleAlignment === 'left',
        'text-center': this.options.titleAlignment === 'center' && !this.useInlineTitleAlignment,
        'text-right': this.options.titleAlignment === 'right' && !this.useInlineTitleAlignment,
      }
    },

    subtitleAlignmentClass () {
      // Default to left alignment when no inline markers are present
      return {
        'text-left': !this.useInlineSubtitleAlignment || this.options.subtitleAlignment === 'left',
        'text-center': this.options.subtitleAlignment === 'center' && !this.useInlineSubtitleAlignment,
        'text-right': this.options.subtitleAlignment === 'right' && !this.useInlineSubtitleAlignment,
      }
    },

    headerClass () {
      return 'cursor-pointer'
    },

    bodyFieldConfig () {
      if (!this.options.bodyField || !this.module) {
        return null
      }
      const field = this.module.fields.find(f => f.name === this.options.bodyField || f.fieldID === this.options.bodyField)
      if (!field) {
        return null
      }
      return field
    },

    bodyField () {
      return this.options.bodyField
    },

    bodyFieldValue () {
      if (!this.bodyFieldConfig || !this.record) {
        return ''
      }
      // Use field name to look up value in record.values (record.values uses field names as keys)
      const fieldName = this.bodyFieldConfig.name
      return this.record.values[fieldName] || ''
    },

    otherFields () {
      if (!this.module) {
        return []
      }

      const selectedFields = this.options.fields || []

      if (selectedFields.length === 0) {
        return []
      }

      const bodyFieldName = this.options.bodyField
      const fields = this.module.filterFields(selectedFields).filter(f => f.name !== bodyFieldName && f.fieldID !== bodyFieldName)

      return fields.map(f => {
        f.label = f.isSystem ? this.$t(`field:system.${f.name}`) : f.label || f.name
        return f
      })
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

    isProcessing () {
      return this.loadingRecord || this.evaluating
    },

    collapsedMinHeight () {
      // Simple fixed values: 7 for title only, 9 for title + subtitle
      if (this.displaySubtitle) {
        return 9.3
      }
      return 7.3
    },

    collapseMaxHeightStyle () {
      // When collapsed, set max-height to 0 to ensure content is fully hidden
      if (this.collapsed) {
        return { maxHeight: '0px', overflow: 'hidden' }
      }
      return { maxHeight: 'none', overflow: 'visible' }
    },
  },

  watch: {
    'options.defaultCollapsed': {
      immediate: true,
      handler (collapsed) {
        this.collapsed = collapsed
        // Set initial height if collapsed on load
        if (collapsed && !this.editable) {
          this.$nextTick(() => {
            this.notifyHeightChange(true)
          })
        }
      },
    },

    loadingRecord: {
      immediate: true,
      handler (loadingRecord) {
        const { recordID } = this.record || {}

        if (!recordID || loadingRecord) return

        this.evaluating = true

        this.evaluateExpressions().finally(() => {
          this.evaluating = false
        })

        // Fetch related records so Record-type field labels resolve in title/subtitle expressions
        if (this.module && this.namespace) {
          this.fetchRecords(this.namespace.namespaceID, this.module.fields, [this.record])
        }
      },
    },

    collapsed: {
      immediate: true,
      handler (collapsed) {
        this.$emit('collapse-change', collapsed)
      },
    },
  },

  mounted () {
    if (this.block && this.block.xywh) {
      this.originalHeight = this.block.xywh[3]
    }
  },

  methods: {
    toggleCollapse () {
      this.collapsed = !this.collapsed
    },

    handleCollapseShown () {
      // Animation complete - content is now visible (expanded)
      this.$emit('collapse-change', false)
      this.notifyHeightChange(false)
    },

    handleCollapseHidden () {
      // Animation complete - content is now hidden (collapsed)
      this.$emit('collapse-change', true)
      this.notifyHeightChange(true)
    },

    notifyHeightChange (isCollapsed) {
      // Only notify in view mode (not in page builder/editor)
      if (this.editable) return
      if (!this.block || !this.block.xywh) return

      // Use originalHeight if available, otherwise get current height from block
      let currentHeight = this.originalHeight
      if (!currentHeight) {
        currentHeight = this.block.xywh[3]
      }
      if (!currentHeight) return

      // When collapsed (content hidden), use collapsedMinHeight
      // When expanded (content visible), use currentHeight
      const newHeight = isCollapsed ? this.collapsedMinHeight : currentHeight

      this.onBlockHeightChange({
        blockIndex: this.blockIndex,
        newHeight,
      })
    },
  },
}
</script>

<style scoped>
.collapsible-record-wrap {
  transition: height 0.3s ease;
  overflow: hidden;
}

.collapsible-record {
  border: 1px solid #dee2e6;
  border-radius: 0.25rem;
  overflow: hidden;
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

.collapsible-content {
  transition: height 0.3s ease;
}

.body-field-container {
  border-bottom: 1px solid #dee2e6;
}

.field-col > * {
  margin-left: 1rem;
  margin-right: 1rem;
}
</style>
