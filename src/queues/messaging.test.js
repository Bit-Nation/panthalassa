import dbFactory from '../database/db';
import messagingQueueFactory from './messaging';
import type {DBInterface} from "../database/db";
import {MESSAGING_QUEUE_JOB_ADDED} from '../events';
const EventEmitter = require('eventemitter3');

function createDbPath() {
    return 'database/'+Math.random();
}

describe('messaging', () => {

    describe('addJob', () => {

        test('emit event', (done) => {

            const db = {
                write: () => new Promise((res, rej) => res())
            };

            const eventEmitter = new EventEmitter();

            eventEmitter.on(MESSAGING_QUEUE_JOB_ADDED, done);

            const queue = messagingQueueFactory(eventEmitter, db);

            queue
                .addJob('Nation', 'Your nation ABC was created successfully')
                .then()
                .catch(done.fail)

        });

        test('save to database', () => {

        })

    });

    test('removeJob', () => {

    });

    test('messages', () => {

    })

});