export const actions = {
  async nuxtServerInit({ commit, state })  {
      commit('config/updateConfig', state.config)
  }
}
