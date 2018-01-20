//@flow

import type {NationType} from "../database/schemata";
import type {DBInterface} from "../database/db";
import type {TransactionQueueInterface} from "../queues/transaction";
import type {TransactionJobInputType} from "../queues/transaction";
import {NATION_CREATE} from '../events';
import {NATION_CONTRACT_ABI} from '../constants'
const Web3 = require('web3');
const waterfall = require('async/waterfall');
const BigNumber = require('bignumber.js');
const EventEmitter = require('eventemitter3');

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
    create(nationData:NationInputType) : Promise<NationType>,
    all() : Promise<NationType>
}

/**
 *
 * @param db
 * @param {TransactionQueueInterface} txQueue
 * @param {Web3} web3
 * @param {EventEmitter} ee
 * @param {string} nationContractAddress
 * @return {NationInterface}
 */
export default function (db:DBInterface, txQueue:TransactionQueueInterface, web3:Web3, ee:EventEmitter, nationContractAddress: string) {

    const impl:NationInterface = {
        create: (nationData:NationInputType) : Promise<NationType> => new Promise((res, rej) => {

            //@todo check if address is valid
            const address = web3.eth.defaultAccount;

            db
                .write(function (realm) {

                    //Persist nation data
                    //created is set to false, since the nation is not written
                    //to the blockchain
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
                        governanceService: nationData.governanceService
                    });

                    return nation;
                })
                .then((nation:NationType) => {

                    const nationContract = web3.eth.contract(NATION_CONTRACT_ABI).at(nationContractAddress);

                    /**
                     * Calculate gas for the whole nation creation process
                     */
                    const gasSummary = {};
                    waterfall([
                        (cb) => {
                            //@todo change the to address to real address of contract
                            gasSummary.nationCore = web3.eth.estimateGas({
                                data: nationContract.createNationCore.getData(
                                    nation.nationName,
                                    nation.nationDescription,
                                    nation.exists,
                                    nation.virtualNation
                                ),
                                to: '0x627306090abaB3A6e1400e9345bC60c78a8BEf57'
                            });

                            cb()
                        },
                        (cb) => {

                            gasSummary.nationPolicy = web3.eth.estimateGas({
                                data: nationContract.SetNationPolicy.getData(
                                    999999999999999999,
                                    nation.nationCode,
                                    "",
                                    nation.lawEnforcementMechanism,
                                    nation.profit
                                ),
                                to: "0x627306090abaB3A6e1400e9345bC60c78a8BEf57"
                            });

                            cb()

                        },
                        (cb) => {

                            gasSummary.nationGovernance = web3.eth.estimateGas({
                                data: nationContract.SetNationGovernance.getData(
                                    999999999999999999,
                                    nation.decisionMakingProcess,
                                    nation.diplomaticRecognition,
                                    nation.governanceService,
                                    nation.nonCitizenUse
                                ),
                                to: "0x627306090abaB3A6e1400e9345bC60c78a8BEf57"
                            });

                            cb();
                        }
                    ], function (err) {

                        if (err) {
                            return rej(err);
                        }

                        //Gas price for nation creation (in wei). ethToSpend = gas * gasPrice
                        const gasPrice = web3.toWei(30, 'gwei');

                        //price for all three transactions (nation creation must be submitted in multiple transactions)
                        let totalPrice = new BigNumber(gasSummary.nationCore + gasSummary.nationPolicy + gasSummary.nationGovernance);
                        totalPrice = web3.fromWei(totalPrice.times(gasPrice), 'ether').toString(10);

                        /**
                         * submit data to job queue
                         */
                        const createNation = () => {

                            const txJob:TransactionJobInputType = {
                                timeout: 30,
                                processor: 'NATION',
                                data: {
                                    dataBaseId: nation.id,
                                    gasSummary: gasSummary,
                                    gasPrice: gasPrice,
                                    from: address
                                },
                                successHeading: `Nation created`,
                                successBody: `Your nation: ${nation.nationName} was created successfully`,
                                failHeading: 'Failed to create nation',
                                failBody: ''
                            };

                            txQueue
                                .addJob(txJob)
                                .then(_ => res(nation))
                                .catch(rej);

                        };

                        ee.emit(NATION_CREATE, {
                            heading: `Confirm nation creation`,
                            msg: `In order to create your nation you have to pay ${totalPrice} ETH. Please confirm or abort it.`,
                            confirm: createNation,
                            abort: async () => await db.write((realm) => realm.delete(nation))
                        });

                    });

                })
                .catch(rej)

        }),
        all: () => db.query((realm) => realm.objects('Nation'))
    };

    return impl;

}