import Clients from '@/generated/core/clients'

export const state = function() {
    var data = Clients.getConfigEnvMap()
    return data
}

export const getters = {
  getConfig(state) {
    return state.config
  }
}

export const mutations = {
  setConfigParam(state, {name, value}) {
    state[name] = value
  },
  updateConfig(state, newConfig) {
    state = newConfig
  }
}
