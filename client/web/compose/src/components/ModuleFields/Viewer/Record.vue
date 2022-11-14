<template>
  <b-spinner
    v-if="processing"
    variant="primary"
    small
  />

  <div v-else>
    <span
      v-for="(v, index) of formattedValue"
      :key="index"
      class="text-nowrap"
      :class="{ 'd-block': field.options.multiDelimiter === '\n' }"
      @click.stop
    >
      <router-link
        :to="v.to"
        :class="{ 'text-decoration-none default-cursor': !v.to}"
      >
        {{ v.value }}{{ index !== formattedValue.length - 1 ? field.options.multiDelimiter : '' }}
      </router-link>
    </span>
  </div>
</template>

<script>
import base from './base'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { mapActions, mapGetters } from 'vuex'

export default {
  extends: base,

  data () {
    return {
      processing: false,

      recordValues: {},

      relRecords: [],
    }
  },

  computed: {
    ...mapGetters({
      pages: 'page/set',
      findUserByID: 'user/findByID',
    }),

    formattedValue () {
      const value = Array.isArray(this.value) ? this.value : [this.value].filter(v => v) || []
      return value.map(recordID => {
        return {
          to: this.linkToRecord(recordID),
          value: this.recordValues[recordID] || recordID,
        }
      })
    },

    recordPage () {
      return this.pages.find(p => p.moduleID === this.field.options.moduleID)
    },
  },

  watch: {
    value: {
      immediate: true,
      handler (value) {
        this.formatRecordValues(value)
      },
    },
  },

  methods: {
    ...mapActions({
      findModuleByID: 'module/findByID',
      resolveUsers: 'user/fetchUsers',
    }),

    linkToRecord (recordID) {
      if (!this.recordPage || !recordID) {
        return ''
      }

      return {
        name: 'page.record',
        params: {
          pageID: this.recordPage.pageID,
          recordID: recordID,
        },
      }
    },

    async formatRecordValues (value) {
      value = Array.isArray(value) ? value : [value].filter(v => v) || []
      const { namespaceID = NoID } = this.namespace
      const { moduleID = NoID, labelField, recordLabelField } = this.field.options

      if (!value.length || [moduleID, namespaceID].includes(NoID) || !labelField) {
        return
      }

      this.processing = true

      // Get configured module/field
      return this.findModuleByID({ namespace: this.namespace, moduleID }).then(module => {
        let relatedField = module.fields.find(({ name }) => name === labelField)
        const query = value.map(recordID => `recordID = ${recordID}`).join(' OR ')

        return this.$ComposeAPI.recordList({ namespaceID, moduleID, query, deleted: 1 }).then(async ({ set = [] }) => {
          if (recordLabelField) {
            set = await this.findModuleByID({ namespaceID, moduleID: relatedField.options.moduleID }).then(relatedModule => {
              const mappedIDs = {}
              const queryIDs = []

              set.forEach(r => {
                r = new compose.Record(module, r)
                mappedIDs[r.values[labelField]] = r.recordID
                queryIDs.push(`recordID = ${r.values[labelField]}`)
              })

              return this.$ComposeAPI.recordList({ namespaceID, moduleID: relatedField.options.moduleID, query: queryIDs.join(' OR '), deleted: 1 }).then(({ set = [] }) => {
                relatedField = relatedModule.fields.find(({ name }) => name === this.field.options.recordLabelField)
                return set.map(r => {
                  r.recordID = mappedIDs[r.recordID]
                  return new compose.Record(relatedModule, r)
                })
              })
            })
          } else {
            set = set.map(r => new compose.Record(module, r))
          }

          for (const record of set) {
            let recordValue = relatedField.isMulti ? record.values[relatedField.name] : [record.values[relatedField.name]]

            if (recordValue.length && relatedField.kind === 'User') {
              recordValue = await Promise.all(recordValue.map(async v => {
                if (!this.findUserByID(v)) {
                  await this.resolveUsers(v)
                }

                return relatedField.formatter(this.findUserByID(v))
              }))
            }

            this.$set(this.recordValues, record.recordID, recordValue.join(relatedField.options.multiDelimiter))
          }
        })
      }).finally(() => {
        setTimeout(() => {
          this.processing = false
        }, 300)
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.default-cursor {
  cursor: default;
}
</style>
