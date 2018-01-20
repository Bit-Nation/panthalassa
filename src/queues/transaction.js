//@flow

import type {TransactionJobType} from "../database/schemata";
import type {DBInterface} from "../database/db";
import {TRANSACTION_QUEUE_JOB_ADDED} from '../events'
const EventEmitter = require('eventemitter3');
const each = require('async/each');
const Realm = require('realm');

/**
 * @typedef TransactionQueueInterface
 * @property {function(job:TransactionJobInputType) : Promise<void>} addJob Add's an job to the queue and emit's the "transaction_queue:job:added" event
 */
export interface TransactionQueueInterface {
    addJob(job:TransactionJobInputType) : Promise<void>,
    registerProcessor(name:string, processor: (done: () => void, data:TransactionJobType) => void) : void,
    process() : Promise<void>
}

/**
 * @typedef TransactionJobInputType
 * @property {number} timeout
 * @property {string} processor
 * @property {object} data
 * @property {string} successHeading
 * @property {string} successBody
 * @property {string} failHeading
 * @property {string} failBody
 */
export type TransactionJobInputType = {
    timeout: number,
    processor: string,
    data: {...mixed},
    successHeading: string,
    successBody: string,
    failHeading: string,
    failBody: string,
}

/**
 *
 * @param {DBInterface} db
 * @param {EventEmitter} ee
 * @return {TransactionQueueInterface}
 */
export default function (db:DBInterface, ee:EventEmitter) : TransactionQueueInterface {

    const impl = {
        processors: {},
        addJob: (job:TransactionJobInputType) : Promise<void> => new Promise((res, rej) => {
            db
                .write((realm) => {

                    realm.create('TransactionJob', {
                        timeout: job.timeout,
                        processor: job.processor,
                        data: JSON.stringify(job.data),
                        id: realm.objects('TransactionJob').length +1,
                        version: 1,
                        successHeading: job.successHeading,
                        successBody: job.successBody,
                        failHeading: job.failHeading,
                        failBody: job.failBody,
                        status: 'WAITING'
                    })

                })
                .then(_ => {

                    ee.emit(TRANSACTION_QUEUE_JOB_ADDED);

                    res();

                })
                .catch(rej)
        }),
        registerProcessor: (name:string, processor: (done: () => void, data:TransactionJobType) => void) : void => {

            impl.processors[name] = processor;

        },
        process: () : Promise<void> => new Promise((res, rej) => {

            db
                .query(function (realm) {
                    return realm.objects('TransactionJob');
                })
                .then((jobs) => {

                    each(jobs, function (TXJob:TransactionJobType, cb) {

                        //find processor
                        const processor = impl.processors[TXJob.processor];

                        if(typeof processor !== 'function'){
                            return rej(new Error(`Couldn't find processor for ${TXJob.processor}`));
                        }

                        /**
                         * @desc This should be called in the processor to end the job
                         * @param data
                         */
                        function done(data) {

                            if(typeof data !== 'object'){
                                return rej('data is not an object');
                            }

                            db
                                .write((realm:Realm) => realm.create('TransactionJob', data, true))
                                .then(_ => cb())
                                .catch(error => cb(error));

                        }

                        //JOSN.parse / stringify is used to remove the realm context from the object
                        //there might be a better solution for it
                        processor(done, JSON.parse(JSON.stringify(TXJob)));

                    }, (error) => {

                        if(error){
                            return rej(error)
                        }

                        res();

                    });

                })
                .catch(rej)

        })
    };

    return impl;

}
