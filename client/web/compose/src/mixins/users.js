export default {
  methods: {
    fetchUsers (fields = [], records = []) {
      if (records.length === 0 || fields.length === 0) {
        return
      }

      const list = [...new Set(records.map(r => {
        return fields
          .filter(c => c.kind === 'User')
          .map(f => {
            if (f.isSystem) {
              return [r[f.name]]
            } else {
              return f.isMulti ? r.values[f.name] : [r.values[f.name]]
            }
          })
      }).flat(Infinity))]

      if (list.length) {
        return this.$store.dispatch('user/resolveUsers', list)
      }
    },
  },
}
