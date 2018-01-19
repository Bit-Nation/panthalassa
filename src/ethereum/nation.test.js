import nationsFactory from './nation';
import dbFactory from '../database/db';

const randomPath = () => 'database/'+Math.random();

describe('nation', () => {

    test('create nation with tx queue job', (done) => {

        const txQueue = {
            addJob: jest.fn((data) => {

                expect(data.timeout).toBe(30);
                expect(data.processor).toBe('NATION');
                expect(data.data).toEqual({
                    dataBaseId: 1
                });
                expect(data.successHeading).toBe('Nation created');
                expect(data.successBody).toBe('Your nation: Bitnation was created successfully');
                expect(data.failHeading).toBe('Failed to create nation');
                expect(data.failBody).toBe('');

                return new Promise((res, rej) => res());

            })
        };

        const nations = nationsFactory(dbFactory(randomPath()), txQueue);

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
            governanceService: 'Security'
        };

        nations
            .create(nationData)
            .then(nationData => {

                //Make sure that the id is incremented by one
                expect(txQueue.addJob).toHaveBeenCalledTimes(1);

                //Created by write action
                expect(nationData.id).toBe(1);
                //Default value from realm
                expect(nationData.idInSmartContract).toBe(-1);
                //Should be false since all nation's are only created locally
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

                done();

            })
            .catch(done.fail);

    });

    test('increment database id', (done) => {

        const txQueue = {
            addJob: () => new Promise((res, rej) => res())
        };

        const nations = nationsFactory(dbFactory(randomPath()), txQueue);

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
            governanceService: 'Security'
        };

        nations
            .create(nationData)
            .then(createdNation => {

                //Created by write action
                expect(createdNation.id).toBe(1);

                //create a second nation with same data
                return nations.create(nationData);
            })
            .then(nationData => {

                //Make sure that the id is incremented by one
                expect(nationData.id).toBe(2);
                
                done();

            })
            .catch(done.fail);

    });

    test('nations', (done) => {

        const txQueue = {
            addJob: () => new Promise((res, rej) => res())
        };

        const nations = nationsFactory(dbFactory(randomPath()), txQueue);

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
            governanceService: 'Security'
        };

        nations
            .create(nationData)
            .then(_ => nations.create(nationData))
            .then(_ => nations.create(nationData))
            .then(_ => nations.all())
            .then(nations => {

                expect(nations.length).toBe(3);

                done();

            })
            .catch(done.fail);

    })

});
