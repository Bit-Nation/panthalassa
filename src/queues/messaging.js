// @flow

import type {MessageJobType} from '../database/schemata';
import type {DBInterface} from '../database/db';
const EventEmitter = require('eventemitter3');
import {MESSAGING_QUEUE_JOB_ADDED} from './../events';

/**
 * @typedef MessagingQueueInterface
 * @property {function(heading:string, text:string)} addJob Add's an job to the messaging queue
 * @property {function(id:number)} removeJob removes an job by it's id
 * @property {function()} messages fetch all message job's
 */
export interface MessagingQueueInterface {
    addJob(heading: string, text: string) : Promise<void>,
    removeJob(id: number) : Promise<void>,
    messages() : Promise<Array<MessageJobType>>
}

/**
 * @desc Factory for the messagingJob queue
 * @param {EventEmitter} eventEmitter
 * @param {DBInterface} db
 * @return {MessagingQueueInterface}
 */
export default function(eventEmitter: EventEmitter, db: DBInterface): MessagingQueueInterface {
    const impl:MessagingQueueInterface = {
        addJob: (heading: string, text: string) => new Promise((res, rej) => {
            // Writes job to database
            const writeAction = (realm) => {
                const numberOfJobs = realm.objects('MessageJob').length;

                const job:MessageJobType = {
                    id: numberOfJobs +1,
                    heading: heading,
                    text: text,
                    version: 1,
                    created_at: new Date(),
                };

                realm.create('MessageJob', job);
            };

            db
                // Write to db
                .write(writeAction)
                .then((_) => {
                    eventEmitter.emit(MESSAGING_QUEUE_JOB_ADDED);

                    res();
                })
                .catch(rej);
        }),
        removeJob: (id: number) => new Promise((res, rej) => {
            const writeAction = (realm) => {
                const messageJob = realm.objects('MessageJob').filtered(`id = "${id}"`);

                realm.delete(messageJob);
            };

            db
                .write(writeAction)
                .then((_) => res())
                .catch(rej);
        }),
        messages: () => db.query((realm) => realm.objects('MessageJob')),
    };

    return impl;
}
