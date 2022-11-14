<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>
    <div
      v-else-if="module"
      class="mt-3"
    >
      <div
        v-for="(field, index) in fields"
        :key="index"
        class="d-flex flex-column mb-3 px-3"
      >
        <label
          class="text-primary mb-0"
          :class="{ 'mb-0': !!(field.options.description || {}).view || false }"
        >
          {{ field.label || field.name }}
          <hint
            :id="field.fieldID"
            :text="(field.options.hint || {}).view || ''"
            class="d-inline-block"
          />
        </label>

        <small
          class="text-muted"
        >
          {{ (field.options.description || {}).view }}
        </small>
        <div
          v-if="field.canReadRecordValue"
          class="value mt-2"
        >
          <field-viewer
            v-bind="{ ...$props, field }"
          />
        </div>
        <i
          v-else
          class="text-primary"
        >
          {{ $t('field.noPermission') }}
        </i>
      </div>
    </div>
  </wrap>
</template>
<script>
import { NoID } from '@cortezaproject/corteza-js'
import base from './base'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import Hint from 'corteza-webapp-compose/src/components/Common/Hint.vue'
import users from 'corteza-webapp-compose/src/mixins/users'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldViewer,
    Hint,
  },

  extends: base,

  mixins: [
    users,
  ],

  computed: {
    fields () {
      if (!this.module) {
        // No module, no fields
        return []
      }

      if (!this.options.fields || this.options.fields.length === 0) {
        // No fields defined in the options, show all (buy system)
        return this.module.fields.slice().sort((a, b) => a.label.localeCompare(b.label))
      }

      // Show filtered & ordered list of fields
      return this.module.filterFields(this.options.fields).map(f => {
        f.label = f.isSystem ? this.$t(`field:system.${f.name}`) : f.label || f.name
        return f
      })
    },

    processing () {
      return !this.record
    },
  },

  watch: {
    'record.recordID': {
      immediate: true,
      handler (recordID) {
        if (recordID && recordID !== NoID) {
          this.fetchUsers(this.fields, [this.record])
        }
      },
    },
  },

}
</script>
<style lang="scss">
.value {
  min-height: 1.2rem;
}
</style>
