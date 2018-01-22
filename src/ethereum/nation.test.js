import nationsFactory from './nation';
import dbFactory from '../database/db';
import {NATION_CREATE} from '../events';
import {NATION_CONTRACT_ABI} from '../constants';
const EventEmitter = require('eventemitter3');
const Web3 = require('web3');

const randomPath = () => 'database/'+Math.random();

describe('nation', () => {
    test('create nation with tx queue job', (done) => {
        const txQueue = {
            addJob: jest.fn((data) => {
                expect(data.timeout).toBe(60);
                expect(data.processor).toBe('NATION');
                expect(data.data).toEqual({
                    dataBaseId: 1,
                    from: '0x85c725a18b09907e874229fcaf36f4e16792214d',
                    gasPrice: '30000000000',
                    steps: {
                        createNationCore: {
                            done: false,
                            txHash: '',
                        },
                        setNationPolicy: {
                            done: false,
                            txHash: '',
                        },
                        setNationGovernance: {
                            done: false,
                            txHash: '',
                        },
                    },
                });
                expect(data.successHeading).toBe('Nation created');
                expect(data.successBody).toBe('Your nation: Bitnation was created successfully');
                expect(data.failHeading).toBe('Failed to create nation');
                expect(data.failBody).toBe('');

                return new Promise((res, rej) => res());
            }),
        };

        const web3 = new Web3(new Web3.providers.HttpProvider('https://rinkeby.infura.io/development'));
        web3.eth.defaultAccount = '0x85c725a18b09907e874229fcaf36f4e16792214d';

        const ee = new EventEmitter();

        const nationData = {
            nationName: 'Bitnation',
            nationDescription: 'We <3 cryptography',
            exists: true,
            virtualNation: false,
            nationCode: 'Code civil',
            lawEnforcementMechanism: 'xyz',
            profit: true,
            nonCitizenUse: false,
            diplomaticRecognition: false,
            decisionMakingProcess: 'dictatorship',
            governanceService: 'Security',
        };

        // Make sure to confirm nation creation
        ee.on(NATION_CREATE, function(eventData) {
            eventData.confirm();
        });

        const nationContractMock = {
            createNation: jest.fn(function(data, cb) {
                expect(JSON.parse(data)).toEqual(nationData);
                cb(null, 'I_AM_A_TRANSACTION_HASH');
            }),
        };

        const nations = nationsFactory(dbFactory(randomPath()), txQueue, web3, ee, nationContractMock);

        nations
            .create(nationData)
            .then((nationData) => {
                // Make sure that the smart contract was called
                expect(nationContractMock.createNation).toHaveBeenCalledTimes(1);

                // Created by write action
                expect(nationData.id).toBe(1);
                // Default value from realm
                expect(nationData.idInSmartContract).toBe(-1);
                // Should be false since all nation's are only created locally
                expect(nationData.created).toBe(false);
                expect(nationData.nationName).toBe('Bitnation');
                expect(nationData.nationDescription).toBe('We <3 cryptography');
                expect(nationData.exists).toBe(true);
                expect(nationData.virtualNation).toBe(false);
                expect(nationData.nationCode).toBe('Code civil');
                expect(nationData.lawEnforcementMechanism).toBe('xyz');
                expect(nationData.profit).toBe(true);
                expect(nationData.nonCitizenUse).toBe(false);
                expect(nationData.decisionMakingProcess).toBe('dictatorship');
                expect(nationData.governanceService).toBe('Security');
                expect(nationData.txHash).toBe('I_AM_A_TRANSACTION_HASH');

                done();
            })
            .catch(done.fail);
    });

    test('increment database id', (done) => {
        const txQueue = {
            addJob: () => new Promise((res, rej) => res()),
        };

        const nationContractMock = {
            createNation: jest.fn(function(data, cb) {
                cb(null, 'I_AM_A_TRANSACTION_HASH');
            }),
        };

        const ee = new EventEmitter();

        // Make sure to confirm nation creation
        ee.on(NATION_CREATE, function(eventData) {
            eventData.confirm();
        });

        const web3 = new Web3(new Web3.providers.HttpProvider('https://rinkeby.infura.io/d'));
        web3.eth.defaultAccount = '0x85c725a18b09907e874229fcaf36f4e16792214d';

        const nations = nationsFactory(dbFactory(randomPath()), txQueue, web3, ee, nationContractMock);

        const nationData = {
            nationName: 'Bitnation',
            nationDescription: 'We <3 cryptography',
            exists: true,
            virtualNation: false,
            nationCode: 'Code civil',
            lawEnforcementMechanism: 'xyz',
            profit: true,
            nonCitizenUse: false,
            diplomaticRecognition: false,
            decisionMakingProcess: 'dictatorship',
            governanceService: 'Security',
        };

        nations
            .create(nationData)
            .then((createdNation) => {
                // Created by write action
                expect(createdNation.id).toBe(1);

                // create a second nation with same data
                return nations.create(nationData);
            })
            .then((nationData) => {
                // Make sure that the id is incremented by one
                expect(nationData.id).toBe(2);

                done();
            })
            .catch(done.fail);
    });

    test('nations', (done) => {
        const txQueue = {
            addJob: () => new Promise((res, rej) => res()),
        };

        const ee = new EventEmitter();

        // Make sure to confirm nation creation
        ee.on(NATION_CREATE, function(eventData) {
            eventData.confirm();
        });

        const nationContractMock = {
            createNation: jest.fn(function(data, cb) {
                cb(null, 'I_AM_A_TRANSACTION_HASH');
            }),
        };

        const web3 = new Web3(new Web3.providers.HttpProvider('https://rinkeby.infura.io/d'));
        web3.eth.defaultAccount = '0x85c725a18b09907e874229fcaf36f4e16792214d';

        const nations = nationsFactory(dbFactory(randomPath()), txQueue, web3, ee, nationContractMock);

        const nationData = {
            nationName: 'Bitnation',
            nationDescription: 'We <3 cryptography',
            exists: true,
            virtualNation: false,
            nationCode: 'Code civil',
            lawEnforcementMechanism: 'xyz',
            profit: true,
            nonCitizenUse: false,
            diplomaticRecognition: false,
            decisionMakingProcess: 'dictatorship',
            governanceService: 'Security',
        };

        nations
            .create(nationData)
            .then((_) => nations.create(nationData))
            .then((_) => nations.create(nationData))
            .then((_) => nations.all())
            .then((nations) => {
                expect(nations.length).toBe(3);

                done();
            })
            .catch(done.fail);
    });

    describe('joinNation', () => {

        test('success', (done) => {

            const nationContractMock = {
                joinNation: jest.fn(function (nationId, cb) {
                    expect(nationId).toEqual(4);
                    cb();
                })
            };

            const web3 = new Web3(new Web3.providers.HttpProvider('https://rinkeby.infura.io/d'));
            web3.eth.defaultAccount = '0x85c725a18b09907e874229fcaf36f4e16792214d';

            const nations = nationsFactory(null, null, null, null, nationContractMock);

            nations
                .joinNation(4)
                .then(_ => {
                    expect(_).toBeUndefined();
                    done();
                })
                .catch(done.fail)

        });

        test('fail', (done) => {

            const nationContractMock = {
                joinNation: jest.fn(function (nationId, cb) {
                    expect(nationId).toEqual(4);
                    cb('i_am_a_error');
                    done();
                })
            };

            const web3 = new Web3(new Web3.providers.HttpProvider('https://rinkeby.infura.io/d'));
            web3.eth.defaultAccount = '0x85c725a18b09907e874229fcaf36f4e16792214d';

            const nations = nationsFactory(null, null, null, null, nationContractMock);

            nations
                .joinNation(4)
                .then(_ => {
                    done.fail('should be rejected');
                })
                .catch(error => {
                    expect(error).toBe('i_am_a_error');
                    done();
                })

        })

    });

    describe('leaveNation', () => {

        test('success', (done) => {

            const nationContractMock = {
                leaveNation: jest.fn(function (nationId, cb) {
                    expect(nationId).toEqual(4);
                    cb();
                })
            };

            const web3 = new Web3(new Web3.providers.HttpProvider('https://rinkeby.infura.io/d'));
            web3.eth.defaultAccount = '0x85c725a18b09907e874229fcaf36f4e16792214d';

            const nations = nationsFactory(null, null, null, null, nationContractMock);

            nations
                .leaveNation(4)
                .then(_ => {
                    expect(_).toBeUndefined();
                    done();
                })
                .catch(done.fail)

        });

        test('fail', (done) => {

            const nationContractMock = {
                leaveNation: jest.fn(function (nationId, cb) {
                    expect(nationId).toEqual(4);
                    cb('i_am_a_error');
                    done();
                })
            };

            const web3 = new Web3(new Web3.providers.HttpProvider('https://rinkeby.infura.io/d'));
            web3.eth.defaultAccount = '0x85c725a18b09907e874229fcaf36f4e16792214d';

            const nations = nationsFactory(null, null, null, null, nationContractMock);

            nations
                .leaveNation(4)
                .then(_ => {
                    done.fail('should be rejected');
                })
                .catch(error => {
                    expect(error).toBe('i_am_a_error');
                    done();
                })

        })

    });

    describe('index', () => {

        test('NationCreated event filter', (done) => {

            const nationContract = {
                NationCreated: (filter, blockFilter) => {

                    expect(filter).toEqual({});

                    expect(blockFilter).toEqual({
                        fromBlock: 0,
                        toBlock:  'latest'
                    });


                    return {
                        get: () => {

                            done();

                        }
                    }

                }
            };

            const nations = nationsFactory(null, null, null, null, nationContract);

            nations
                .index()
                .then()

        });

        test('NationCreated event error handling', (done) => {

            const nationContract = {
                NationCreated: (filter, blockFilter) => {

                    expect(filter).toEqual({});

                    expect(blockFilter).toEqual({
                        fromBlock: 0,
                        toBlock:  'latest'
                    });

                    return {
                        get: (cb) => {
                            cb('i_am_an_error');
                        }
                    }

                }
            };

            const nations = nationsFactory(null, null, null, null, nationContract);

            nations
                .index()
                .then(_ => done.fail('should be rejected'))
                .catch(error => {

                    expect(error).toBe('i_am_an_error');
                    done();

                })

        });

    })

});
