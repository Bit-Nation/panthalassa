import transactionQueue from './transaction';
import messagingQueue from './messaging';
import type {TransactionJobType} from "../database/schemata";
const EventEmitter = require('eventemitter3');
import {TRANSACTION_QUEUE_JOB_ADDED} from '../events';
import dbFactory from "../database/db";

const dbPath = () => 'database/'+Math.random();

describe('transaction', () => {

    describe('addJob', () => {

        test('emit "transaction_queue:job:added" event on save', (done) => {

            const db = {
                write: () => new Promise((res, rej) => res())
            };

            const ee = new EventEmitter();

            ee.on(TRANSACTION_QUEUE_JOB_ADDED, () => {
                done();
            });

            const txQueue = transactionQueue(db, ee);

            txQueue
                .addJob({})
                .then()
                .catch()

        });

        test('write job to database', (done) => {

            const db = dbFactory(dbPath());

            const txQueue = transactionQueue(db, new EventEmitter());

            const job:TransactionJobType = {
                timeout: 30,
                processor: 'my_processor',
                data: {
                    key_one: 'data_one',
                    key_two: 'data_two'
                },
                successHeading: 'success_heading',
                successBody: 'success_body',
                failHeading: 'fail_heading',
                failBody: 'fail_body'
            };

            txQueue
                .addJob(job)
                .then(result => {
                    expect(result).toBeUndefined();

                    db
                        .query((realm) => realm.objects('TransactionJob'))
                        .then(jobs => {

                            expect(jobs[0].timeout).toBe(30);
                            expect(jobs[0].id).toBe(1);
                            expect(jobs[0].processor).toBe('my_processor');
                            expect(jobs[0].data).toBe(JSON.stringify({key_one: 'data_one', key_two: 'data_two'}));
                            expect(jobs[0].successHeading).toBe('success_heading');
                            expect(jobs[0].successBody).toBe('success_body');
                            expect(jobs[0].failHeading).toBe('fail_heading');
                            expect(jobs[0].failBody).toBe('fail_body');
                            expect(jobs[0].version).toBe(1);
                            expect(jobs[0].status).toBe('WAITING');

                            done();

                        })
                        .catch(done.fail);

                })
                .catch(done.fail);

        })

    });

    test('registerProcessor', () => {

        const db = dbFactory(dbPath());

        const txQueue = transactionQueue(db, new EventEmitter());

        expect(txQueue.processors['my_processor']).toBeUndefined();

        txQueue.registerProcessor('my_processor', function () {

        });

        expect(txQueue.processors['my_processor']).toBeDefined();

    });

    describe('process', () => {

        test('processor not found', (done) => {

            const db = dbFactory(dbPath());

            const job:TransactionJobType = {
                timeout: 30,
                processor: 'my_processor',
                data: {
                    key_one: 'data_one',
                    key_two: 'data_two'
                },
                successHeading: 'success_heading',
                successBody: 'success_body',
                failHeading: 'fail_heading',
                failBody: 'fail_body'
            };

            const txQueue = transactionQueue(db, new EventEmitter());

            txQueue
                .addJob(job)
                .then(_ => txQueue.process())
                .then(_ => {
                    done.fail('The promise should be rejected');
                })
                .catch(e => {
                    expect(e.message).toBe('Couldn\'t find processor for my_processor');
                    done();
                })

        });

        test('added success message when job is completed', (done) => {

            const db = dbFactory(dbPath());

            const job:TransactionJobType = {
                timeout: 30,
                processor: 'my_processor',
                data: {
                    key_one: 'data_one',
                    key_two: 'data_two'
                },
                successHeading: 'success_heading',
                successBody: 'success_body',
                failHeading: 'fail_heading',
                failBody: 'fail_body'
            };

            const msgQueue = {
                addJob: jest.fn(function (title, body) {
                    expect(title).toBe('success_heading');
                    expect(body).toBe('success_body');
                    return new Promise((res, rej) => res());
                })
            };

            const txQueue = transactionQueue(db, new EventEmitter(), msgQueue);

            txQueue.registerProcessor('my_processor', function (done, txJob) {

                expect(typeof txJob).toBe('object');

                //Set job it to fake id. The transaction queue should fix that
                //by setting the id of the job passed by the "done" callback
                //to it's correct value
                txJob.id = 1111111;

                txJob.status = 'DONE';

                done(txJob);

            });

            txQueue
                .addJob(job)
                .then(_ => db.query((realm) => realm.objects('TransactionJob')))
                .then(jobs => {

                    expect(jobs[0].status).toBe('WAITING');
                    expect(jobs[0].id).toBe(1);

                })
                .then(_ => txQueue.process())
                .then(_ => db.query((realm) => realm.objects('TransactionJob')))
                .then(jobs => {

                    expect(jobs[0].status).toBe('DONE');
                    expect(jobs[0].id).toBe(1);
                    done();

                })

                .catch(done.fail);

        })

    })

});