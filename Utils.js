function removePrivateData(user) {
    Object.keys(user).forEach(key => {
        if(typeof user[key] === 'object')
            removePrivateData(user[key])
        if(['access_token', 'refresh_token',
            'password', 'local', 'email', 'avatar',
            'createdAt', 'updatedAt', 'linked', 'discord', 'google', 'github',
            '__v', '_id'].includes(key))
            delete user[key]
    })
    return user;
}
function removeData(user, keys) {
    if(!(keys instanceof Array)) throw new Error("keys argument should be Array. Recieved: " + typeof keys)
    Object.keys(user).forEach(key => {
        if(typeof user[key] === 'object')
            removeData(user[key], keys)
        if(keys.includes(key))
            delete user[key]
    })
    return user;
}
module.exports = {removePrivateData, removeData}

