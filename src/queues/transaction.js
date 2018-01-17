//@flow

import type {TransactionJobType} from "../database/schemata";
import type {DBInterface} from "../database/db";
import {TRANSACTION_QUEUE_JOB_ADDED} from '../events'
import {MessagingQueueInterface} from "./messaging";
const EventEmitter = require('eventemitter3');
const each = require('async/each');
const Realm = require('realm');

/**
 * @typedef TransactionQueueInterface
 * @property {function(job:TransactionJobType) : Promise<void>} addJob Add's an job to the queue and emit's the "transaction_queue:job:added" event
 */
export interface TransactionQueueInterface {
    addJob(job:TransactionJobType) : Promise<void>,
    registerProcessor(name:string, processor: (done: () => void, data:TransactionJobType) => void) : void,
    process() : Promise<void>
}

export default function (db:DBInterface, ee:EventEmitter, messagingQueue:MessagingQueueInterface) : TransactionQueueInterface {

    const impl = {
        processors: {},
        addJob: (job:TransactionJobType) : Promise<void> => new Promise((res, rej) => {
            db
                .write((realm) => {

                    const countOfJobs = realm.objects('TransactionJob').length;

                    //Set job default data
                    job.id = countOfJobs +1;
                    job.data = JSON.stringify(job.data);
                    job.version = 1;
                    job.status = 'WAITING';

                    realm.create('TransactionJob', job)

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

                        function done() {

                            db
                                .write(function (realm:Realm) : void {

                                    //update the job (true mean's the it's an update)
                                    const job = realm.create('TransactionJob', TXJob, true);

                                    //if the job is done we will remove it and add a message for the user
                                    if(job.status === 'DONE'){
                                        messagingQueue
                                            .addJob(
                                                job.messages.successHeading,
                                                job.messages.successBody
                                            ).then(_ => {
                                                realm.delete(job);
                                                cb()
                                            })
                                            .catch(cb);

                                        return;
                                    }

                                    if(job.status === 'FAILED'){
                                        messagingQueue
                                            .addJob(
                                                job.messages.failHeading,
                                                job.messages.failBody
                                            )
                                            .then(_ => cb())
                                            .catch(cb)
                                    }

                                })
                                .then(cb)
                                .catch(cb)

                        }

                        //
                        processor(done)

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
