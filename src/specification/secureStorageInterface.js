//@flow

/**
 * All method of the interface take a amount of input data and return
 * the callable methods. They are working like mini factorys.
 */
interface SecureStorage {

    /**
     * Set a key and a value. Return's a void promise that get's
     * resolved after the record is written to the key value / secure storage
     */
    set(...any) : ((key:string, value:any) =>  Promise<void>);

    /**
     * Get a value by it's key. And return a promise that can be resolved
     * with any value.
     */
    get(...any) : ((key:string) => Promise<any>);

    /**
     * Prove if a value exist by the key. It return's a
     * promise which will be resolved with true or false
     */
    has(...any) : ((key:string) => Promise<boolean>);

    /**
     * Remove a key, value pair based on the key. The promise
     * will be with void resolved.
     */
    remove(...any) : ((key:string) => Promise<void>);

    /**
     * Loops over all elements in the secure store and filter based on the given callback.
     * the callback need's to return true / false. When true is returned the dataset will be
     * added to the return list
     */
    fetchItems(...any) : ((filter: (key:string, value:string) => boolean) => Promise<Array<{key:string, value:any}>>);

    /**
     * Destroys the whole storage
     */
    destroyStorage(...any) : (() => Promise<any>);

}
