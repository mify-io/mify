import Clients from '@/generated/core/clients'

export default function(ctx, inject) {
  const $clients = new Clients(ctx)
  inject('clients', $clients)
  ctx.$clients = $clients
}
