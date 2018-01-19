//@flow

import type {NationType} from "../database/schemata";
import type {DBInterface} from "../database/db";

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
 * @return {NationInterface}
 */
export default function (db:DBInterface) {

    const impl:NationInterface = {
        create: (nationData:NationInputType) => new Promise((res, rej) => {

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
                .then((nationDBId:NationType) => {

                    //@todo place queue job

                    res(nationDBId);

                })
                .catch(rej)

        }),
        all: () => db.query((realm) => realm.objects('Nation'))
    };

    return impl;

}