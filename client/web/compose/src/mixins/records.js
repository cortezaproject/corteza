export default {
  methods: {
    fetchRecords (namespaceID, fields = [], records = []) {
      if (records.length === 0 || fields.length === 0) {
        return
      }

      const moduleRecords = {}

      fields.filter(c => c.kind === 'Record').forEach(f => {
        const { moduleID } = f.options || {}
        if (!moduleRecords[moduleID]) {
          moduleRecords[moduleID] = new Set()
        }

        records.forEach(r => {
          let recordIDs = []

          if (f.isSystem) {
            recordIDs = [r[f.name]]
          } else {
            recordIDs = f.isMulti ? r.values[f.name] : [r.values[f.name]]
          }

          recordIDs.forEach(recordID => {
            if (!recordID) return
            moduleRecords[moduleID].add(recordID)
          })
        })
      })

      // Dispatch resolution per module
      return Promise.all(Object.entries(moduleRecords).map(([moduleID, recordIDs]) => {
        recordIDs = [...recordIDs]

        if (recordIDs.length) {
          return this.$store.dispatch('record/resolveRecords', { namespaceID, moduleID, recordIDs })
        }
      }))
    },
  },
}
