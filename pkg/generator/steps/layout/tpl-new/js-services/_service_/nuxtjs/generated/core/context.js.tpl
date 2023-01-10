// vim: set ft=javascript:
import Clients from './clients'

class MifyContext {
    constructor() {
        this._config = Clients.getConfigEnvMap()
        this._clients = new Clients(this._config)
    }

    get config() {
        return this._config
    }

    get clients() {
        return this._clients
    }
}

export default MifyContext;
