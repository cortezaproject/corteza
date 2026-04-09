// This mixin is used on View component of Records.

import { NoID, compose, validator } from '@cortezaproject/corteza-js'
import { throttle } from 'lodash'
import { mapGetters } from 'vuex'

export default {
  data () {
    return {
      inEditing: false,
      processing: false,
      processingAction: '',
      record: undefined,
      initialRecordState: undefined,
      errors: new validator.Validated(),
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
    }),

    validator () {
      if (!this.module) {
        throw new Error('can not initialize record validator without module')
      }

      return new compose.RecordValidator(this.module)
    },

    isValid () {
      return this.errors.valid()
    },

    /**
     * Tells if given record is deleted; If record not provided, returns undefined
     * @returns {Boolean}
     */
    isDeleted () {
      return this.record && this.record.deletedAt
    },
  },

  watch: {
    'record.valueErrors': {
      handler () {
        this.setWarnings()
      },
    },

    processing: {
      handler (processing) {
        // If processing is set to false we know that one of them is also true, so we reset all of them since we don't know which one is true
        if (!processing) {
          this.processingAction = ''
        }
      },
    },
  },

  methods: {
    /**
     * Handle form submit for record create & update
     *
     *  -> dispatch beforeFormSubmit (on ui:compose:record-page)
     *  -> validate record (see validateRecord())
     *     -> stop on errors
     *  -> send record to the API
     *  -> apply changes received from the API to current record
     *  -> dispatch afterFormSubmit
     *  -> redirect user to record viewer page
     *
     * @returns {Promise<void>}
     */
    handleFormSubmit: throttle(async function (route = 'page.record') {
      this.processingAction = 'submit'
      this.processing = true

      let record
      const isNew = this.record.recordID === NoID
      const queue = []

      // Collect records from all record lines
      this.blocks.forEach((b, index) => {
        if (b.kind === 'RecordList' && b.options.editable) {
          const p = new Promise((resolve) => {
            const recordListUniqueID = [this.page.pageID, (this.record || {}).recordID, b.blockID, false].map(v => v || NoID).join('-')
            this.$root.$emit(`record-line:collect:${recordListUniqueID}`, resolve)
          })

          queue.push(p)
        }
      })

      const pairs = await Promise.all(queue)

      for (const p of pairs) {
        if (p.positionField) {
          let i = 0
          for (const item of p.items) {
            if (!item.r.deletedAt) {
              item.r.values[p.positionField] = i++
            }
          }
        }
      }

      // Construct batch record payload
      const records = pairs.reduce((acc, cur) => {
        if (cur.idPrefix) {
          // If same module exists, use latest to avoid stale data
          const existingIndex = acc.findIndex(({ module }) => module.moduleID === cur.module.moduleID)
          if (existingIndex !== -1) {
            acc[existingIndex].set = cur.items.map(({ r }) => r).filter(({ deletedAt, recordID }) => recordID !== NoID || !deletedAt)
          } else {
            acc.push({
              refField: cur.refField,
              set: cur.items.map(({ r }) => r).filter(({ deletedAt, recordID }) => recordID !== NoID || !deletedAt),
              module: cur.module,
              idPrefix: cur.idPrefix,
            })
          }
        }

        return acc
      }, [])

      const { recordID = NoID } = this.record || {}

      // Append after the payload construction, so it is not presented as a
      // sub record.
      pairs.push({
        module: this.module,
        items: [{ r: this.record, id: recordID === NoID ? 'parent:0' : recordID }],
      })

      // Clone record so unsaved changes doesn't trigger false positive because of serialize values
      this.record = this.record.clone()

      return this
        .dispatchUiEvent('beforeFormSubmit', this.record, { $records: records })
        .then(() => {
          // Evaluate layout required fields before validation
          return this.evaluateLayoutRequiredFields()
        })
        .then(() => this.validateRecord(pairs))
        .then(() => {
          if (isNew) {
            return this.$ComposeAPI.recordCreate({ ...this.record, records })
          } else {
            return this.$ComposeAPI.recordUpdate({ ...this.record, records })
          }
        }).catch(err => {
          this.processing = false

          const { details = undefined } = err

          if (!!details && Array.isArray(details) && details.length > 0) {
            this.errors.push(...details.filter(d => !d.kind.includes('warning')))

            throw new Error(this.$t('notification:record.validationErrors', { fields: this.getValidationErrorFields() }))
          }

          throw err
        }).then(r => {
          record = new compose.Record(this.module, r)
        }).then(() => this.dispatchUiEvent('afterFormSubmit', record, { $records: records }))
        .then(() => {
          // Clear draft revision on successful save
          if (this.activeDraftKey && this.$store) {
            this.$store.dispatch('drafts/removeDraft', { changeID: this.activeDraftKey })
            this.activeDraftKey = null
          }

          if (record.valueErrors.set) {
            this.record = record.clone()
            this.initialRecordState = record.clone()
            this.setWarnings()
            this.$root.$emit('refetch-records', { recordID: record.recordID })
            this.toastWarning(this.$t('notification:record.validationWarnings', { errors: this.errors, fields: this.getValidationErrorFields({ includeWarnings: true }) }))

            this.processing = false
          } else {
            this.initialRecordState = this.record.clone()

            if (this.inModal) {
              this.$emit('handle-record-redirect', { recordID: record.recordID, recordPageID: this.page.pageID, edit: false })

              // If we are in a modal we need to refresh blocks/records not in modal
              this.$root.$emit('refetch-records', { recordID: record.recordID })
            } else {
              // Refresh blocks that are related to the updated records
              const relatedRecords = [
                this.module.moduleID,
                ...new Set(records.filter(r => r.module.moduleID !== this.module.moduleID).map(r => r.module.moduleID)),
              ]
              relatedRecords.forEach(moduleID => this.$root.$emit('module-records-updated', { moduleID }))

              this.$router.push({ name: route, params: { ...this.$route.params, recordID: record.recordID, edit: false } })
            }

            if (this.page.meta.notifications.enabled) {
              this.toastSuccess(this.$t(`notification:record.${isNew ? 'create' : 'update'}Success`))
            }
          }
        }).catch(e => {
          this.processing = false
          this.toastErrorHandler(this.$t(`notification:record.${isNew ? 'create' : 'update'}Failed`))(e)
        })
    }, 500),

    /**
     * Handle form submit for record browser
     * @returns {Promise<void>}
     */
    handleFormSubmitSimple: throttle(function (route = 'admin.modules.record.view') {
      this.processingAction = 'submit'
      this.processing = true

      let record
      const isNew = this.record.recordID === NoID

      // Clone record so unsaved changes doesn't trigger false positive because of serialize values
      this.record = this.record.clone()

      return this
        .dispatchUiEvent('beforeFormSubmit')
        .then(() => this.validateRecordSimple())
        .then(() => {
          if (isNew) {
            return this.$ComposeAPI.recordCreate(this.record)
          } else {
            return this.$ComposeAPI.recordUpdate(this.record)
          }
        }).catch(err => {
          this.processing = false

          const { details = undefined } = err

          if (!!details && Array.isArray(details) && details.length > 0) {
            this.errors.push(...details.filter(d => !d.kind.includes('warning')))

            throw new Error(this.$t('notification:record.validationErrors', { fields: this.getValidationErrorFields() }))
          }

          throw err
        }).then(r => {
          record = new compose.Record(this.module, r)
        }).then(() => this.dispatchUiEvent('afterFormSubmit', record))
        .then(() => {
          // Clear draft revision on successful save
          if (this.activeDraftKey && this.$store) {
            this.$store.dispatch('drafts/removeDraft', { changeID: this.activeDraftKey })
            this.activeDraftKey = null
          }

          if (record.valueErrors.set) {
            this.record = record.clone()
            this.initialRecordState = record.clone()
            this.setWarnings()
            this.toastWarning(this.$t('notification:record.validationWarnings', { fields: this.getValidationErrorFields({ includeWarnings: true }) }))
            this.processing = false
          } else {
            this.initialRecordState = this.record.clone()
            this.$router.push({ name: route, params: { ...this.$route.params, recordID: record.recordID, edit: false } })
            this.toastSuccess(this.$t(`notification:record.${isNew ? 'create' : 'update'}Success`))
          }
        }).catch(this.toastErrorHandler(this.$t(`notification:record.${isNew ? 'create' : 'update'}Failed`)))
    }, 500),

    /**
     * On delete, preserve user's view. Show a notification that the record
     * has been deleted.
     */
    handleDelete: throttle(function () {
      this.processing = true
      this.processingAction = 'delete'

      return this.dispatchUiEvent('beforeDelete')
        .then(() => this.$ComposeAPI.recordDelete(this.record))
        .then(this.dispatchUiEvent('afterDelete'))
        .then(() => {
          // Clear draft revision on successful delete
          if (this.activeDraftKey && this.$store) {
            this.$store.dispatch('drafts/removeDraft', { changeID: this.activeDraftKey })
            this.activeDraftKey = null
          }

          this.$root.$emit('refetch-records')
          this.toastSuccess(this.$t('notification:record.deleteSuccess'))
        }).catch(e => {
          this.processing = false
          this.toastErrorHandler(this.$t('notification:record.deleteFailed'))(e)
        })
    }, 500),

    handleUndelete: throttle(function () {
      this.processingAction = 'undelete'
      this.processing = true

      return this.dispatchUiEvent('beforeUndelete')
        .then(() => this.$ComposeAPI.recordUndelete(this.record))
        .then(this.dispatchUiEvent('afterUndelete'))
        .then(() => {
          this.$root.$emit('refetch-records')
          this.toastSuccess(this.$t('notification:record.restoreSuccess'))
        }).catch(e => {
          this.processing = false
          this.toastErrorHandler(this.$t('notification:record.restoreFailed'))(e)
        })
    }, 500),

    handleBulkUpdateSelectedRecords: throttle(function (query) {
      this.processing = true

      const values = []

      this.fields.forEach(f => {
        const { name, isMulti, isSystem } = this.getField(f)
        const value = isSystem ? this.record[name] : this.record.values[name]

        if (!isMulti) {
          // Handle single value case
          values.push({
            name,
            value: value?.toString() ?? '',
          })
        } else {
          // Handle multi-value case
          if (!Array.isArray(value) || value.length === 0) {
            values.push({ name })
            return
          }

          // Map non-undefined values to proper format
          const multiValues = value
            .filter(v => v !== undefined)
            .map(v => ({
              name,
              value: v?.toString() ?? '',
            }))

          values.push(...multiValues)
        }
      })

      const { moduleID, namespaceID } = this.module

      return this
        .$ComposeAPI.recordPatch({ moduleID, namespaceID, values, query })
        .catch(err => {
          const { details = undefined } = err

          if (!!details && Array.isArray(details) && details.length > 0) {
            this.errors = new validator.Validated()
            this.errors.push(...details)

            throw new Error(this.$t('notification:record.validationErrors', { fields: this.getValidationErrorFields() }))
          }

          throw err
        })
        .then(() => {
          this.toastSuccess(this.$t('notification:record.bulkRecordUpdateSuccess'))
          this.showModal = false
          this.$emit('save')

          this.$root.$emit('refetch-records', { stayOnPage: true })
        }).catch(this.toastErrorHandler(this.$t('notification:record.deleteBulkRecordUpdateFailed')))
        .finally(() => {
          this.processing = false
        })
    }, 500),

    getValidationErrorFields ({ includeWarnings = false, includeErrors = true } = {}) {
      const { set = [] } = this.errors || {}

      const fields = new Set(set.filter(({ meta = {} } = {}) => {
        if (includeWarnings) {
          return meta.isWarning
        } else if (includeErrors) {
          return !meta.isWarning
        }

        return true
      }).map(d => {
        const fieldName = d.meta.field
        // Try to resolve field label from module
        const mod = (d.meta.moduleID && this.getModuleByID(d.meta.moduleID)) || this.module
        if (mod) {
          const f = mod.fields.find(f => f.name === fieldName)
          if (f?.label) return f.label
        }
        return fieldName
      }))

      return Array.from(fields).join(', ')
    },

    /**
     * Validates record and dispatches onFormSubmitError
     *
     * onFormSubmitError is dispatched only if there are record errors,
     * if not, we continue with form submit handling
     *
     * After onFormSubmitError, record is re-validated and if errors
     * are still present, we stop form submit handing
     *
     * @returns {Promise<void>}
     */
    async validateRecord (pairs) {
      const layoutRequiredFields = this.$store.getters['ui/layoutRequiredFields'] || []

      // Create validators with modified modules (if layout required fields exist)
      const validators = {}
      for (const p of pairs) {
        let moduleForValidator = p.module

        // If we have layout required fields, create a modified copy of the module
        if (layoutRequiredFields.length > 0) {
          // Clone the module with modified fields
          const modifiedFields = p.module.fields.map(field => {
            const isLayoutRequired = layoutRequiredFields.includes(field.name) || layoutRequiredFields.includes(field.fieldID)
            if (isLayoutRequired && !field.isRequired) {
              // Create a copy of the field with isRequired set to true
              return new compose.ModuleField({ ...field, isRequired: true })
            }
            return field
          })

          // Create a new module with the modified fields
          moduleForValidator = new compose.Module({ ...p.module, fields: modifiedFields })
        }

        validators[p.module.resourceID] = validators[p.module.resourceID] || new compose.RecordValidator(moduleForValidator)
      }

      const vRunner = () => {
        // Reset errors
        this.errors = new validator.Validated()

        // validate
        for (const p of pairs) {
          const v = validators[p.module.resourceID]
          const errs = new validator.Validated()

          p.items.forEach(({ r, id }) => {
            if (r.deletedAt) {
              return
            }

            const fields = p.module.fields
              .filter(({ canReadRecordValue, canUpdateRecordValue }) => canReadRecordValue && canUpdateRecordValue)
              .map(({ name }) => name)

            // cover the edge case where all fields are not updatable
            if (fields.length) {
              const err = v.run(r, ...fields)
              if (!err.valid()) {
                err.applyMeta({ id })
                errs.push(...err.set)
              }
            }
          })

          this.errors.push(...errs.set)
        }
      }

      vRunner()
      if (this.errors.valid()) {
        return
      }

      await this.dispatchUiEvent('onFormSubmitError')
      vRunner()
      if (!this.errors.valid()) {
        throw new Error(this.$t('notification:record.validationErrors', { fields: this.getValidationErrorFields() }))
      }
    },

    /**
     * Validates record browser record
     *
     * @returns {Promise<void>}
     */
    async validateRecordSimple () {
      this.errors = this.validator.run(this.record)
      if (this.errors.valid()) {
        return
      }

      await this.dispatchUiEvent('onFormSubmitError')

      this.errors = this.validator.run(this.record)
      if (!this.errors.valid()) {
        throw new Error(this.$t('notification:record.validationErrors', { fields: this.getValidationErrorFields() }))
      }
    },

    setWarnings () {
      const { set = [] } = this.record.valueErrors || {}
      this.errors.push(...set.map(e => {
        e.meta.isWarning = true
        return e
      }))
    },

    resetErrors () {
      this.errors = new validator.Validated()
    },

    /**
     * Returns errors, filtered for a specific field
     * @param name
     * @returns {validator.Validated} filtered validation results
     */
    fieldErrors (name) {
      if (!this.errors) {
        return new validator.Validated()
      }

      return this.errors.filterByMeta('field', name)
    },

    /**
     * Generic event dispatcher for ui:compose:record-page resource type
     *
     * This is used when deleting, updating, creating
     * records and where validation errors occur
     *
     * @param eventType
     */
    dispatchUiEvent (eventType, record = this.record, args = {}) {
      const resourceType = `ui:compose:${this.getUiEventResourceType || 'record-page'}`

      const argsBase = {
        errors: this.errors,
        validator: this.validator,
        ...args,
      }

      return this
        .$EventBus
        .Dispatch(compose.RecordEvent(
          record, { eventType, resourceType, args: argsBase }))
    },
  },
}
