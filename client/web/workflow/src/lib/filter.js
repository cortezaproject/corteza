import { filters } from '@cortezaproject/corteza-vue-next'

export function objectSearchMaker (field, ...fields) {
  return function (opts, search) {
    if (!search) return opts
    const s = search.toLowerCase()
    return opts.filter(o => {
      const allFields = [field, ...fields]
      return allFields.some(f => {
        const val = typeof o === 'object' ? o[f] : o
        return val && String(val).toLowerCase().includes(s)
      })
    })
  }
}

export function stringSearchMaker () {
  return function (opts, search) {
    if (!search) return opts
    const s = search.toLowerCase()
    return opts.filter(o => String(o).toLowerCase().includes(s))
  }
}
