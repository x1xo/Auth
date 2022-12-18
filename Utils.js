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

function removeInvalidToken(lindex){
    if(global.invalidTokenTimeout._destroyed && lindex){
        global.invalidTokenTimeout = setTimeout(async () => {
            await global.redis.rpop("invalid_tokens")
            const lindex = await global.redis.lindex("invalid_tokens", -1);
            removeInvalidToken(lindex)
        }, parseInt(lindex.split(",")[1])-Date.now())
    }
}

module.exports = {removePrivateData, removeData, removeInvalidToken}

