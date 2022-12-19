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
/**
 * 
 * @returns If Invalid: true, If Valid: false
 */
async function checkForInvalidToken(result){
    //this will be the date of logout in miliseconds
    let lastLogout = await global.redis.zscore("logouts",result.id);
    if(lastLogout > parseInt(`${result.iat}000`))
        return true;
    return false;
}
module.exports = {removePrivateData, removeData, checkForInvalidToken}

