import nationsFactory from './nation';
import dbFactory from '../database/db';

const randomPath = () => 'database/'+Math.random();

describe('nation', () => {

    test('createNation', (done) => {

        const nations = nationsFactory(dbFactory(randomPath()));

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
                //Default value from realm
                expect(createdNation.idInSmartContract).toBe(-1);
                //Should be false since all nation's are only created locally
                expect(createdNation.created).toBe(false);
                expect(createdNation.nationName).toBe('Bitnation');
                expect(createdNation.nationDescription).toBe('We <3 cryptography');
                expect(createdNation.exists).toBe(true);
                expect(createdNation.virtualNation).toBe(false);
                expect(createdNation.nationCode).toBe('Code civil');
                expect(createdNation.lawEnforcementMechanism).toBe('xyz');
                expect(createdNation.profit).toBe(true);
                expect(createdNation.nonCitizenUse).toBe(false);
                expect(createdNation.decisionMakingProcess).toBe('dictatorship');
                expect(createdNation.governanceService).toBe('Security');

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

        const nations = nationsFactory(dbFactory(randomPath()));

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
