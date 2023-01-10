// vim: set ft=javascript:
import Clients from './clients'

class MifyContext {
    constructor(config) {
        this._config = config
        this._clients = new Clients(this._config)
    }

    get config() {
        return this._config
    }

    get clients() {
        return this._clients
    }

    static getConfig() {
        return Clients.getConfigEnvMap()
    }
}

export default MifyContext;
