function removePrivateData(user) {
    Object.keys(user).forEach(key => {
        if(typeof user[key] === 'object')
            removePrivateData(user[key])
        if(['access_token', 'refresh_token',
            'password', 'local', 'email',
            '__v', '_id'].includes(key))
            delete user[key]
    })
    return user;
}
module.exports = {removePrivateData}

