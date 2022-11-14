import {
  LOADER_INC,
  LOADER_DEC,
  LOADER_HIDE,
} from './types'

export default {
  incLoader ({ commit }) {
    commit(LOADER_INC)
  },

  decLoader ({ commit }) {
    commit(LOADER_DEC)
  },

  hideLoader ({ commit }) {
    commit(LOADER_HIDE)
  },
}
