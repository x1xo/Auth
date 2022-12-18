const jwt = require('jsonwebtoken')
const { removeInvalidToken } = require('../../Utils')
const Router = require('express').Router()

Router.get('/', (req, res) => {
    res.send(`<html><script>window.onload = (e) => {fetch('/user/logout', {method: "POST"}).then(res => window.location==="/")}</script></html>`)
})

Router.post('/', async (req, res) => {
    //logout - invalidate jwt
    try {
        if(!req.body.token && !req.cookies.token) return res.send({error: true, code: 401, message: "UNAUTHORIZED"})
        const result = jwt.verify(req.body.token || req.cookies.token, req.app.get("secret"), {algorithms: ["RS256"]})
        if(parseInt(`${result.exp}000`) < Date.now())
            return res.send({error: true, code: 401, message: "JWT_EXPIRED"})
        await global.redis.lpush("invalid_tokens", `${req.body.token},${result.exp}000`)
        res.clearCookie('token');
        res.send({error: false, status: "SUCCESS"});
        const lindex = await global.redis.lindex("invalid_tokens", -1);
        removeInvalidToken(lindex);   
    } catch (e) {
        if(e.message.includes('jwt expired'))
            return res.send({error: true, message: "JWT_EXPIRED"})

        return res.send({error: true,
            code: 500,
            message: "SOMETHING_WENT_WRONG",
            errorStack: `${e.message} ${e.stack}`})
    }
})

module.exports = Router;