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

        test('save to database', (done) => {

            const db = dbFactory(createDbPath());

            const queue = messagingQueueFactory(new EventEmitter(), db);

            queue
                .addJob('Nation', 'Your nation ABC was created successfully')
                .then(_ => {

                    //Query the database for the job

                    db
                        .query(function (realm) {
                            return realm.objects('MessageJob')
                        })
                        .then(result => {

                            expect(result.length).toBe(1);

                            expect(result[0].id).toBe(1);
                            expect(result[0].heading).toBe('Nation');
                            expect(result[0].text).toBe('Your nation ABC was created successfully');
                            expect(result[0].version).toBe(1);

                            done();

                        })
                        .catch(done.fail)

                })
                .catch(done.fail)


        })

    });

    test('removeJob', () => {

    });

    test('messages', () => {

    })

});