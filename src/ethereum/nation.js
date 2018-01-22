// @flow

import type {NationType} from '../database/schemata';
import type {DBInterface} from '../database/db';
import type {TransactionQueueInterface} from '../queues/transaction';
const Web3 = require('web3');
const EventEmitter = require('eventemitter3');
const eachSeries = require('async/eachSeries');
const waterfall = require('async/waterfall');

/**
 * @typedef NationType
 * @property {boolean} created Created mean's if it is on the blockchain
 * @property {string} nationName
 * @property {string} nationDescription
 * @property {boolean} exists
 * @property {boolean} virtualNation
 * @property {string} nationCode
 * @property {string} lawEnforcementMechanism
 * @property {boolean} profit
 * @property {boolean} nonCitizenUse
 * @property {boolean} diplomaticRecognition
 * @property {string} decisionMakingProcess
 * @property {string} governanceService
 */
export type NationInputType = {
    nationName: string,
    nationDescription: string,
    exists: boolean,
    virtualNation: boolean,
    nationCode: string,
    lawEnforcementMechanism: string,
    profit: boolean,
    nonCitizenUse: boolean,
    diplomaticRecognition: boolean,
    decisionMakingProcess: string,
    governanceService: string
}

/**
 * @typedef NationInterface
 * @property {function(nationData:NationInputType)} create
 * @property {function()} all fetch all nations
 */
export interface NationInterface {
    create(nationData: NationInputType) : Promise<NationType>,
    all() : Promise<NationType>,
    index() : Promise<void>,
    joinNation(id: number) : Promise<void>,
    leaveNation(id: number) : Promise<void>
}

/**
 *
 * @param db
 * @param {TransactionQueueInterface} txQueue
 * @param {Web3} web3
 * @param {EventEmitter} ee
 * @param {object} nationContract
 * @return {NationInterface}
 */
export default function(db: DBInterface, txQueue: TransactionQueueInterface, web3: Web3, ee: EventEmitter, nationContract: {...any}) {
    const impl:NationInterface = {
        create: (nationData: NationInputType): Promise<NationType> => new Promise((res, rej) => {
            db
                .write(function(realm) {
                    // Persist nation data
                    // created is set to false, since the nation is not written
                    // to the blockchain
                    const nation = realm.create('Nation', {
                        id: realm.objects('Nation').length +1,
                        created: false,
                        nationName: nationData.nationName,
                        nationDescription: nationData.nationDescription,
                        exists: nationData.exists,
                        virtualNation: nationData.virtualNation,
                        nationCode: nationData.nationCode,
                        lawEnforcementMechanism: nationData.lawEnforcementMechanism,
                        profit: nationData.profit,
                        nonCitizenUse: nationData.nonCitizenUse,
                        diplomaticRecognition: nationData.diplomaticRecognition,
                        decisionMakingProcess: nationData.decisionMakingProcess,
                        governanceService: nationData.governanceService,
                    });

                    return nation;
                })
                .then((nation: NationType) => {
                    nationContract.createNation(
                        JSON.stringify(nationData),
                        function(err, txHash) {
                            if (err) {
                                return rej(err);
                            }

                            // Attach transaction hash to nation
                            db
                                .write((realm) => nation.txHash = txHash)
                                .then((_) => res(nation))
                                .catch(rej);
                        }
                    );
                })
                .catch(rej);
        }),
        all: () => db.query((realm) => realm.objects('Nation')),
        index: () => new Promise((res, rej) => {
            const nationCreatedEvent = nationContract.NationCreated({}, {fromBlock: 0, toBlock: 'latest'});

            nationCreatedEvent.get(function(err, logs) {
                if (err) {
                    return rej(err);
                }

                const joinedNations = [];

                waterfall(
                    [
                        (cb) => {
                            nationContract.getJoinedNations(function(err, res) {
                                if (err) {
                                    return cb(err);
                                }

                                res.map((nationBigNumber) => joinedNations.push(nationBigNumber.toNumber()));

                                cb();
                            });
                        },
                        (cb) => {
                            eachSeries(logs, function(log, cb) {
                                const nationId = log.args.nationId.toNumber();

                                nationContract.getNumCitizens(nationId, function(err, citizens) {
                                    if (err) {
                                        return cb(err);
                                    }

                                    citizens = citizens.toNumber();

                                    db
                                    // We query for txHash since we get the tx hash when submitting the nation to the blockchain
                                        .query((realm) => realm.objects('Nation').filtered(`txHash = "${log.transactionHash}"`))
                                        .then((nations) => {
                                            const nation = nations[0];

                                            if (nation) {
                                                return db.write((realm) => {
                                                    nation.idInSmartContract = nationId;
                                                    nation.created = true;
                                                    nation.joined = joinedNations.includes(nationId);
                                                    nation.citizens = citizens;
                                                });
                                            }

                                            return new Promise((res, rej) => {
                                                nationContract.getNationMetaData(nationId, function(err, result) {
                                                    if (err) {
                                                        return rej(err);
                                                    }

                                                    try {
                                                        result = JSON.parse(result);
                                                    } catch (e) {
                                                        return rej(e);
                                                    }

                                                    db
                                                        .write((realm) => {
                                                            const nationCount = realm.objects('Nation').length;

                                                            realm.create('Nation', {
                                                                id: nationCount+1,
                                                                idInSmartContract: nationId,
                                                                txHash: log.transactionHash,
                                                                nationName: result.nationName,
                                                                nationDescription: result.nationDescription,
                                                                created: true,
                                                                exists: result.exists,
                                                                virtualNation: result.virtualNation,
                                                                nationCode: result.nationCode,
                                                                lawEnforcementMechanism: result.lawEnforcementMechanism,
                                                                profit: result.profit,
                                                                nonCitizenUse: result.nonCitizenUse,
                                                                diplomaticRecognition: result.diplomaticRecognition,
                                                                decisionMakingProcess: result.decisionMakingProcess,
                                                                governanceService: result.governanceService,
                                                                joined: joinedNations.includes(nationId),
                                                                citizens: citizens,
                                                            });
                                                        })
                                                        .then((_) => res())
                                                        .catch(rej);
                                                });
                                            });
                                        })
                                        .then((_) => setTimeout(cb, 200))
                                        .catch(cb);
                                });
                            }, function(err) {
                                if (err) {
                                    return rej(err);
                                }

                                cb();
                            });
                        },
                    ],
                    function(err) {
                        if (err) {
                            return rej(err);
                        }

                        res();
                    }
                );
            });
        }),
        joinNation: (id: number): Promise<void> => new Promise((res, rej) => {
            nationContract.joinNation(id, function(err, txHash) {
                if (err) {
                    return rej(err);
                }

                res();
            });
        }),
        leaveNation: (id: number): Promise<void> => new Promise((res, rej) => {
            nationContract.leaveNation(id, function(err, txHash) {
                if (err) {
                    return rej(err);
                }

                return res();
            });
        }),
    };

    return impl;
}
