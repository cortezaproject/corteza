// This mixin is used on View component of Records.

import { compose, validator, NoID } from '@cortezaproject/corteza-js'
import { mapGetters, mapActions } from 'vuex'
import { throttle } from 'lodash'

export default {
  data () {
    return {
      inEditing: false,
      processing: false,
      processingDelete: false,
      processingUndelete: false,
      processingSubmit: false,
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
      handler ({ set = [] } = {}) {
        this.errors.push(...set)
      },
    },
  },

  methods: {
    ...mapActions({
      updatePrompts: 'wfPrompts/update',
    }),

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
      this.processingSubmit = true
      this.processing = true

      let record
      const isNew = this.record.recordID === NoID
      const queue = []

      // Collect records from all record lines
      this.page.blocks.forEach((b, index) => {
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

      return this
        .dispatchUiEvent('beforeFormSubmit', this.record, { $records: records })
        .then(() => this.validateRecord(pairs))
        .then(() => {
          if (isNew) {
            return this.$ComposeAPI.recordCreate({ ...this.record, records })
          } else {
            return this.$ComposeAPI.recordUpdate({ ...this.record, records })
          }
        })
        .catch(err => {
          const { details = undefined } = err
          if (!!details && Array.isArray(details) && details.length > 0) {
            this.errors = new validator.Validated()
            this.errors.push(...details)

            throw new Error(this.$t('notification:record.validationErrors'))
          }

          throw err
        })
        .then(r => {
          record = new compose.Record(this.module, r)
        })
        .then(() => this.dispatchUiEvent('afterFormSubmit', this.record, { $records: records }))
        .then(() => this.updatePrompts())
        .then(() => {
          if (this.record.valueErrors.set) {
            this.toastWarning(this.$t('notification:record.validationWarnings'))
          } else {
            this.inCreating = false
            this.inEditing = false
            // reset the record initial state in cases where the record edit page is redirected to the record view page
            this.initialRecordState = this.record.clone()

            if (this.showRecordModal) {
              this.$emit('handle-record-redirect', { recordID: record.recordID, recordPageID: this.page.pageID })

              // If we are in a modal we need to refresh blocks not in modal
              this.$root.$emit('module-records-updated', {
                moduleID: this.module.moduleID,
                notPageID: this.page.pageID,
              })
            } else {
              this.$router.push({ name: route, params: { ...this.$route.params, recordID: record.recordID } })
            }

            if (!isNew) {
              this.record = record
              // reset the record initial state in cases where the record edit page is opened on a modal
              this.initialRecordState = this.record.clone()

              this.determineLayout().then(() => {
                this.$root.$emit(`refetch-non-record-blocks:${this.page.pageID}`)
              })
            }
          }

          this.toastSuccess(this.$t(`notification:record.${isNew ? 'create' : 'update'}Success`))
        })
        .catch(this.toastErrorHandler(this.$t(`notification:record.${isNew ? 'create' : 'update'}Failed`)))
        .finally(() => {
          this.processingSubmit = false
          this.processing = false
        })
    }, 500),

    /**
     * Handle form submit for record browser
     * @returns {Promise<void>}
     */
    handleFormSubmitSimple: throttle(function (route = 'page.record') {
      this.processingSubmit = true
      this.processing = true

      let record
      const isNew = this.record.recordID === NoID

      return this
        .dispatchUiEvent('beforeFormSubmit')
        .then(() => this.validateRecordSimple())
        .then(() => {
          if (isNew) {
            return this.$ComposeAPI.recordCreate(this.record)
          } else {
            return this.$ComposeAPI.recordUpdate(this.record)
          }
        })
        .catch(err => {
          const { details = undefined } = err
          if (!!details && Array.isArray(details) && details.length > 0) {
            this.errors.push(...details)

            throw new Error(this.$t('notification:record.validationErrors'))
          }

          throw err
        })
        .then(r => {
          record = new compose.Record(this.module, r)
        })
        .then(() => this.dispatchUiEvent('afterFormSubmit', record))
        .then(() => this.updatePrompts())
        .then(() => {
          if (this.record.valueErrors.set) {
            this.toastWarning(this.$t('notification:record.validationWarnings'))
          } else {
            this.inCreating = false
            this.inEditing = false
            this.record = record
            this.initialRecordState = this.record.clone()

            this.$router.push({ name: route, params: { ...this.$route.params, recordID: record.recordID } })
          }
        })
        .catch(this.toastErrorHandler(this.$t(
          isNew
            ? 'notification:record.createFailed'
            : 'notification:record.updateFailed',
        )))
        .finally(() => {
          this.processingSubmit = false
          this.processing = false
        })
    }, 500),

    /**
     * On delete, preserve user's view. Show a notification that the record
     * has been deleted.
     */
    handleDelete: throttle(function () {
      this.processing = true
      this.processingDelete = true

      return this
        .dispatchUiEvent('beforeDelete')
        .then(this.$ComposeAPI.recordDelete(this.record))
        .then(() => {
          this.record.deletedAt = (new Date()).toISOString()
        })
        .then(this.dispatchUiEvent('afterDelete'))
        .then(this.updatePrompts())
        .then(this.loadRecord)
        .then(this.toastSuccess(this.$t('notification:record.deleteSuccess')))
        .catch(this.toastErrorHandler(this.$t('notification:record.deleteFailed')))
        .finally(() => {
          this.processing = false
          this.processingDelete = false
        })
    }, 500),

    handleUndelete: throttle(function () {
      this.processingUndelete = true
      this.processing = true

      return this
        .dispatchUiEvent('beforeUndelete')
        .then(this.$ComposeAPI.recordUndelete(this.record))
        .then(this.dispatchUiEvent('afterUndelete'))
        .then(this.updatePrompts())
        .then(this.loadRecord)
        .then(this.toastSuccess(this.$t('notification:record.restoreSuccess')))
        .catch(this.toastErrorHandler(this.$t('notification:record.restoreFailed')))
        .finally(() => {
          this.processingUndelete = false
          this.processing = false
        })
    }, 500),

    handleBulkUpdateSelectedRecords: throttle(function (query) {
      this.processing = true

      const values = []
      this.fields.forEach(f => {
        const { name, isMulti } = this.getField(f)
        const value = this.record.values[name] || this.record[name]

        if (!isMulti) {
          values.push({ name, value: value ? value.toString() : value })
        } else {
          value.forEach(v => {
            values.push({ name, value: v ? v.toString() : v })
          })
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

            throw new Error(this.$t('notification:record.validationErrors'))
          }

          throw err
        })
        .then(this.updatePrompts())
        .then(() => {
          this.toastSuccess(this.$t('notification:record.bulkRecordUpdateSuccess'))
          this.onModalHide()
          this.fields = []
          this.record = new compose.Record(this.module, {})
          this.initialRecordState = this.record.clone()
          this.$emit('save')

          this.$root.$emit('module-records-updated', { moduleID })
        })
        .catch(this.toastErrorHandler(this.$t('notification:record.deleteBulkRecordUpdateFailed')))
        .finally(() => {
          this.processing = false
        })
    }, 500),

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
      // Cache validators for later use
      const validators = {}
      for (const p of pairs) {
        validators[p.module.resourceID] = validators[p.module.resourceID] || new compose.RecordValidator(p.module)
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
              .filter(({ canUpdateRecordValue }) => canUpdateRecordValue)
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
        throw new Error(this.$t('notification:record.validationErrors'))
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
        throw new Error(this.$t('notification:record.validationErrors'))
      }
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
