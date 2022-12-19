const jwt = require('jsonwebtoken')
const User = require('../../models/User')
const { removeData, checkForInvalidToken } = require('../../Utils')

const Router = require('express').Router()

Router.post('/', async (req, res) => {
    try {
        if(!req.body.token) return res.send({error: true, code: 401, message: "UNAUTHORIZED"})
        const result = jwt.verify(req.body.token, req.app.get("secret"), {algorithms: ["RS256"]})
        if(parseInt(`${result.exp}000`) < Date.now())
            return res.send({error: true, code: 401, message: "JWT_EXPIRED"})
        
        if(await checkForInvalidToken(result))
            return res.send({erorr: true, code: 401, message: "INVALID_TOKEN"})

        let user = await User.findOne({id: result.id})
        if(!user) return res.send({error: true, code: 404, message: "NOT_FOUND"})
        
        user = removeData(user, ['access_token', 'refresh_token', 'password', 'local', 'family_name', '__v', '_id'])
        return res.send(user)

    } catch (e) {
        console.log(e)
        return res.send({
            error: true,
            code: 500,
            message: "SOMETHING_WENT_WRONG",
            errorStack: `${e.message} ${e.stack}`
        })
    }
})


module.exports=Router