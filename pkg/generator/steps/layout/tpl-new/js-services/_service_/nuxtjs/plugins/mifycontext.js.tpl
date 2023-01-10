// vim: set ft=javascript:
import MifyContext from '../generated/core/context'

export default function(ctx, inject) {
    inject('mifyContext', new MifyContext(ctx.store.state.mifycontext.config))
}
