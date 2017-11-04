// Profile
const ProfileSchema  = {
    name: 'Profile',
    primaryKey: 'id',
    properties: {
        id: 'int',
        pseudo: 'string',
        description: 'string',
        image: 'string'
    },
};

module.exports = {
    ProfileSchema
};
