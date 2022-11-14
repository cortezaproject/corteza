import {
  LOADER_INC,
  LOADER_DEC,
  LOADER_HIDE,
} from './types'

export default {
  [LOADER_INC] (state) {
    state.loader++
  },

  [LOADER_DEC] (state) {
    if (state.loader > 0) {
      state.loader--
    }
  },

  [LOADER_HIDE] (state) {
    state.loader = 0
  },
}
