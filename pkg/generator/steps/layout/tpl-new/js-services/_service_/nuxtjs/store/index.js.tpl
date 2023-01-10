// vim: set ft=javascript:
import { newContext } from "./mifycontext"
export const actions = {
  async nuxtServerInit({ commit })  {
      commit('mifycontext/update', newContext())
  }
}
