import { filter } from '@cortezaproject/corteza-vue'

export function objectSearchMaker (field, ...fields) {
  return function (opts, search) {
    return opts.filter(o => filter.Assert(o, search, field, ...fields))
  }
}

export function stringSearchMaker () {
  return function (opts, search) {
    return opts.filter(o => filter.Assert(o, search))
  }
}
