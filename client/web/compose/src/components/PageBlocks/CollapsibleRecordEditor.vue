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
        <div class="collapsible-titles">
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
      >
        <!-- Other Fields - Above Body -->
        <div
          v-if="showOtherFieldsAbove && otherFields.length"
          ref="fieldContainerAbove"
          class="other-fields-section px-3 pt-3 pb-3"
          :class="fieldLayoutClass"
        >
          <template v-for="field in otherFields">
            <div
              v-if="canDisplay(field)"
              :key="`${field.fieldID}-${field.name}-above`"
              :class="`field-container ${columnWrapClass}`"
              :style="fieldWidth"
            >
              <field-editor
                v-if="isFieldEditable(field)"
                v-bind="{ ...$props, errors: fieldErrors(field.name) }"
                :horizontal="horizontal"
                :field="field"
                :extra-options="options"
                @change="onFieldChange(field)"
              />

              <b-form-group
                v-else
                :label-cols-md="horizontal && '5'"
                :label-cols-xl="horizontal && '4'"
                :content-cols-md="horizontal && '7'"
                :content-cols-xl="horizontal && '8'"
              >
                <template #label>
                  <div class="d-flex align-items-center text-primary mb-0">
                    <span class="d-flex">
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
                    :field="field"
                    v-bind="{ ...$props, errors: fieldErrors(field.name) }"
                    value-only
                  />
                </div>
                <i
                  v-else
                  class="text-muted"
                >
                  {{ $t('field.noPermission') }}
                </i>
              </b-form-group>
            </div>
          </template>
        </div>

        <!-- Body Field (Full Width) -->
        <div
          v-if="bodyFieldConfig"
          class="body-field-container px-3 pt-3 pb-3"
        >
          <field-editor
            v-if="isFieldEditable(bodyFieldConfig)"
            v-bind="{ ...$props, errors: fieldErrors(bodyFieldConfig.name) }"
            :horizontal="horizontal"
            :field="bodyFieldConfig"
            :extra-options="options"
            @change="onFieldChange(bodyFieldConfig)"
          />
          <div
            v-else-if="bodyFieldConfig.canReadRecordValue"
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

        <!-- Other Fields - Default or Below Body -->
        <div
          v-if="showOtherFieldsAfter && otherFields.length"
          ref="fieldContainerAfter"
          class="other-fields-section px-3 pb-3"
          :class="fieldLayoutClass"
        >
          <template v-for="field in otherFields">
            <div
              v-if="canDisplay(field)"
              :key="`${field.fieldID}-${field.name}-after`"
              :class="`field-container ${columnWrapClass}`"
              :style="fieldWidth"
            >
              <field-editor
                v-if="isFieldEditable(field)"
                v-bind="{ ...$props, errors: fieldErrors(field.name) }"
                :horizontal="horizontal"
                :field="field"
                :extra-options="options"
                @change="onFieldChange(field)"
              />

              <b-form-group
                v-else
                :label-cols-md="horizontal && '5'"
                :label-cols-xl="horizontal && '4'"
                :content-cols-md="horizontal && '7'"
                :content-cols-xl="horizontal && '8'"
              >
                <template #label>
                  <div class="d-flex align-items-center text-primary mb-0">
                    <span class="d-flex">
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
                    :field="field"
                    v-bind="{ ...$props, errors: fieldErrors(field.name) }"
                    value-only
                  />
                </div>
                <i
                  v-else
                  class="text-muted"
                >
                  {{ $t('field.noPermission') }}
                </i>
              </b-form-group>
            </div>
          </template>
        </div>
      </b-collapse>
    </div>
  </wrap>
</template>

<script>
import { NoID, compose } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'
import base from './base'
import FieldEditor from 'corteza-webapp-compose/src/components/ModuleFields/Editor'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import users from 'corteza-webapp-compose/src/mixins/users'
import records from 'corteza-webapp-compose/src/mixins/records'
import conditionalFields from 'corteza-webapp-compose/src/mixins/conditionalFields'
import recordLayout from 'corteza-webapp-compose/src/mixins/recordLayout'
import alignment from 'corteza-webapp-compose/src/mixins/alignment'
import { mapGetters } from 'vuex'
import { debounce } from 'lodash'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldEditor,
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

  data () {
    return {
      collapsed: false,
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
      if (this.useInlineTitleAlignment) return ''

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
      if (this.useInlineSubtitleAlignment) return ''

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
      if (!this.options.subtitleShowWhenCollapsed && this.collapsed) {
        return false
      }
      if (this.useInlineSubtitleAlignment) {
        return !!(this.subtitleAlignmentData.left || this.subtitleAlignmentData.center || this.subtitleAlignmentData.right)
      }
      return !!this.subtitleText
    },

    titleAlignmentClass () {
      return {
        'text-left': !this.useInlineTitleAlignment || this.options.titleAlignment === 'left',
        'text-center': this.options.titleAlignment === 'center' && !this.useInlineTitleAlignment,
        'text-right': this.options.titleAlignment === 'right' && !this.useInlineTitleAlignment,
      }
    },

    subtitleAlignmentClass () {
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
      if (!this.options.bodyField || !this.module) return null
      const field = this.module.fields.find(f => f.name === this.options.bodyField || f.fieldID === this.options.bodyField)
      return field || null
    },

    bodyFieldValue () {
      if (!this.bodyFieldConfig || !this.record) return ''
      const fieldName = this.bodyFieldConfig.name
      return this.record.values[fieldName] || ''
    },

    otherFields () {
      if (!this.module) return []

      const selectedFields = this.options.fields || []
      if (selectedFields.length === 0) return []

      const bodyFieldName = this.options.bodyField
      const fields = this.module.filterFields(selectedFields).filter(f => f.name !== bodyFieldName && f.fieldID !== bodyFieldName)

      return fields.map(f => {
        f.label = f.isSystem ? this.$t(`field:system.${f.name}`) : f.label || f.name
        return f
      })
    },

    showOtherFieldsAbove () {
      return this.options.otherFieldsPosition === 'above'
    },

    showOtherFieldsAfter () {
      const pos = this.options.otherFieldsPosition
      return pos === 'below' || pos === 'default' || !pos
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
      if (this.options.recordFieldLayoutOption !== 'wrap') return ''
      return 'field-col'
    },

    fieldWidth () {
      if (this.options.recordFieldLayoutOption !== 'noWrap') return {}
      return { 'min-width': '20rem' }
    },

    isProcessing () {
      return this.loadingRecord || !this.record || this.evaluating
    },

    horizontal () {
      return this.block.options.horizontalFieldLayoutEnabled
    },

    isNew () {
      return this.record && this.record.recordID === NoID
    },
  },

  watch: {
    'options.defaultCollapsed': {
      immediate: true,
      handler (collapsed) {
        this.collapsed = !!collapsed
      },
    },

    loadingRecord: {
      immediate: true,
      handler (loadingRecord) {
        const { recordID } = this.record || {}

        if (!recordID || loadingRecord) return

        let resolutions = []

        if (recordID !== NoID) {
          resolutions = [
            this.fetchUsers(this.module.fields, [this.record]),
            this.fetchRecords(this.namespace.namespaceID, this.module.fields, [this.record]),
          ]
        }

        this.evaluating = true

        Promise.all([
          ...resolutions,
          this.evaluateExpressions(),
        ]).finally(() => {
          this.evaluating = false
        })
      },
    },
  },

  mounted () {
    this.createEvents()
  },

  beforeDestroy () {
    this.destroyEvents()
  },

  methods: {
    toggleCollapse () {
      this.collapsed = !this.collapsed
    },

    onFieldChange: debounce(function (field) {
      this.$root.$emit('record-field-change', {
        fieldName: field.name,
      })
    }, 500),

    createEvents () {
      this.$root.$on('record-field-change', this.evaluateExpressions)
    },

    destroyEvents () {
      this.$root.$off('record-field-change', this.evaluateExpressions)
    },
  },
}
</script>

<style scoped>
.collapsible-record-wrap {
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

.body-field-container {
  border-bottom: 1px solid #dee2e6;
}

.field-col > * {
  margin-left: 1rem;
  margin-right: 1rem;
}
</style>
</content>
</invoke>