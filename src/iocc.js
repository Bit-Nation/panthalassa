import wallet from "./ethereum/wallet";

const awilix = require('awilix');

import {factory} from "./database/db";
import ethUtils from './ethereum/utils';
import PanthalassaProvider from './ethereum/PanthalassaProvider';
import web3 from "./ethereum/web3";
import profile from "./profile/profile";

const iocc = awilix.createContainer();

iocc.register({
    'database:db' : awilix.asFunction(factory).singleton(),
    'ethereum:PanthalassaProvider': awilix.asFunction((c) => {
        return new PanthalassaProvider(c['ethereum:utils'], c['_ethNode'].rpcUrl)
    }).singleton(),
    'ethereum:utils': awilix.asFunction((c) => ethUtils(
        c['_secureStorage'],
        c['_eventEmitter']
    )).singleton(),
    'ethereum:wallet' : awilix.asFunction((c) => wallet(
        c['ethereum:ethUtils'],
        c['ethereum:web3'],
        c['database:db']
    )),
    'ethereum:web3' : awilix.asFunction((c) => web3(
        c['_ethNode'],
        c['_eventEmitter'],
        c['ethereum:utils']
    )),
    'profile:profile': awilix.asFunction((c) => profile(
        c['database:db'],
        c['ethereum:utils']
    ))
});

export default iocc;