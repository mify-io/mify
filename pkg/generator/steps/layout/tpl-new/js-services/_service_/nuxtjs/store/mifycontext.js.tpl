// vim: set ft=javascript:
import MifyContext from '../generated/core/context'

export function newContext() {
    return {
        "config": MifyContext.getConfig(),
    }
}

export const state = function() {
    return newContext()
}

export const getters = {
  getContext(state) {
    return state.context
  }
}

export const mutations = {
  update(state, newContext) {
    state = newContext
  }
}
