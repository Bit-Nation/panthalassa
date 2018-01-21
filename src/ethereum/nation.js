//@flow

import type {NationType} from "../database/schemata";
import type {DBInterface} from "../database/db";
import type {TransactionQueueInterface} from "../queues/transaction";
import type {TransactionJobInputType} from "../queues/transaction";
import {NATION_CREATE} from '../events';
const Web3 = require('web3');
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

                    /**
                     * submit data to job queue
                     */
                    const createNation = () => {

                        const txJob:TransactionJobInputType = {
                            timeout: 60,
                            processor: 'NATION',
                            data: {
                                dataBaseId: nation.id,
                                gasPrice: web3.toWei(30, 'gwei'),
                                from: address,
                                steps: {
                                    createNationCore: {
                                        done: false,
                                        txHash: ""
                                    },
                                    setNationPolicy: {
                                        done: false,
                                        txHash: ""
                                    },
                                    setNationGovernance: {
                                        done: false,
                                        txHash: ""
                                    }
                                }
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
                        msg: `In order to create your nation you have to pay 0.05 ETH. Please confirm or abort it.`,
                        confirm: createNation,
                        abort: async () => await db.write((realm) => realm.delete(nation))
                    });

                })
                .catch(rej)

        }),
        all: () => db.query((realm) => realm.objects('Nation'))
    };

    return impl;

}