//@flow

import type {NationType} from "../database/schemata";
import type {DBInterface} from "../database/db";
import type {TransactionQueueInterface} from "../queues/transaction";
import type {TransactionJobInputType} from "../queues/transaction";

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
 * @return {NationInterface}
 */
export default function (db:DBInterface, txQueue:TransactionQueueInterface) {

    const impl:NationInterface = {
        create: (nationData:NationInputType) : Promise<NationType> => new Promise((res, rej) => {

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

                    const txJob:TransactionJobInputType = {
                        timeout: 30,
                        processor: 'NATION',
                        data: {
                            dataBaseId: nation.id
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

                })
                .catch(rej)

        }),
        all: () => db.query((realm) => realm.objects('Nation'))
    };

    return impl;

}