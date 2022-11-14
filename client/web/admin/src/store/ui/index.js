import actions from './actions'
import mutations from './mutations'

export default {
  namespaced: true,

  state: {
    loader: 0,
  },

  getters: {
    isLoading (state) {
      return state.loader > 0
    },
  },

  actions,
  mutations,
}
