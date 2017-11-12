import {getAccounts, signTx} from './PanthalassaProvider'

describe('getAccounts', () => {

    test('success', done => {

        const address_one = "0x465868366a0f45748f24d8979a98c2118e71b2bc";

        const address_two = "0x26e75307fc0c021472feb8f727839531f112f317";

        const ethUtils = {
            allKeyPairs: () => {
                return new Promise((res, rej) => {

                    res([
                        {
                            key: address_one
                        },
                        {
                            key: address_two
                        }
                    ])

                })
            }
        };

        const cb = (error, addresses) => {

            expect(error).toBeNull();

            expect(addresses).toEqual([
                address_one,
                address_two
            ]);

            //This will mark the test as done
            done();

        };

        getAccounts(ethUtils)(cb)

    });

    test('error', done => {

        class TestError extends Error{}

        const ethUtils = {
            allKeyPairs: () => {
                return new Promise((res, rej) => {

                    //Reject promise with test error
                    rej(new TestError())

                })
            }
        };

        const cb = (error, addresses) => {

            expect(error).toEqual(new TestError());

            expect(addresses).toBeNull();

            //This will mark the test as done
            done();

        };

        getAccounts(ethUtils)(cb)

    });

});