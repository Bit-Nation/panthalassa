export interface JsonRpcNodeInterface {
    name: string,
    url: string,
    start() : Promise<void>,
    stop() : Promise<void>,
}