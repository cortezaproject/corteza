<template>
  <module-records
    v-slot="{ value }"
    :modules="payloadValues"
  >
    <p
      v-for="(v, vi) in value.value"
      :key="vi"
      class="mb-0"
      :class="{ 'mt-1': vi > 0 }"
    >
      {{ v }}
    </p>
  </module-records>
</template>

<script>
import base from './base'
import ModuleRecords from 'corteza-webapp-privacy/src/components/Common/ModuleRecords'

export default {
  i18nOptions: {
    namespaces: 'request',
    keyPrefix: 'view.delete',
  },

  components: {
    ModuleRecords,
  },

  extends: base,

  computed: {
    payloadValues () {
      const { modules = {} } = this.payload || {}

      return Object.entries(modules).map(([moduleID, { module, namespace, records = {} }]) => {
        records = Object.entries(records).map(([recordID, { values = {} }]) => {
          values = Object.entries(values).map(([name, value = []]) => {
            return { name, value }
          })
          return { recordID, values }
        })
        return { module, namespace, moduleID, records }
      })
    },
  },
}
</script>
