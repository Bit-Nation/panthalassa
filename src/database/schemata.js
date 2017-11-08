//@flow

///////////////////////////////////////////////////////////
// ATTENTION !!! Everyime you update the schema,         //
//               update the relating interfaces as well. //
///////////////////////////////////////////////////////////

//Profile
export interface ProfileObject {
    id: number,
    pseudo: string,
    description: string,
    image: string,
    version: string
}

export const ProfileSchema = {
    name: 'Profile',
    primaryKey: 'id',
    properties: {
        id: 'int',
        pseudo: 'string',
        description: 'string',
        image: 'string',
        version: 'string'
    },
};
