<template>
  <b-button
    v-if="list.length > 0"
    data-test-id="button-export"
    variant="light"
    size="lg"
    @click="jsonExport(list, type)"
  >
    {{ $t('label.export') }}
  </b-button>
</template>

<script>
import { saveAs } from 'file-saver'
import { mapActions } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  props: {
    list: {
      type: Array,
      required: true,
    },
    type: {
      type: String,
      required: true,
    },
  },

  methods: {
    ...mapActions({
      findModuleByID: 'module/findByID',
    }),

    jsonExport (list, type) {
      Promise.all(list.map(i => i.export(this.findModuleByID))).then(list => {
        const blob = new Blob([JSON.stringify({ type, list }, null, 2)], { type: 'application/json' })
        saveAs(blob, `${type}-export.json`)
      })
    },
  },
}
</script>
