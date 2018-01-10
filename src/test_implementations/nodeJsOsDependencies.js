// @flow

// @todo maybe move this in a node implementations folder

import {Crypto, OsDependenciesInterface} from '../specification/osDependencies';

const crypto = require('crypto');

const cryptoImplementation:Crypto = {
    randomBytes: (length: number) => new Promise((res, rej) => {
        crypto.randomBytes(length, (err, buffer) => {
            if (err) {
                return rej(err);
            }

            res(buffer.toString('hex'));
        });
    }),
};

const osDepsImplementation:OsDependenciesInterface = {
    crypto: cryptoImplementation,
};

export default osDepsImplementation;
