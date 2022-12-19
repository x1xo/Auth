const jwt = require('jsonwebtoken')
const { checkForInvalidToken } = require('../../Utils')
const Router = require('express').Router()

Router.post('/', async (req, res) => {
    try {
        if(!req.body.token && !req.cookies.token) return res.send({error: true, code: 401, message: "UNAUTHORIZED"})
        const result = jwt.verify(req.body.token || req.cookies.token, req.app.get("secret"), {algorithms: ["RS256"]})
        if(parseInt(`${result.exp}000`) < Date.now())
            return res.send({error: true, code: 401, message: "JWT_EXPIRED"})
        if(await checkForInvalidToken(result)) return res.send({error: true, code: 401, message: "INVALID_TOKEN"})
        await global.redis.zadd("logouts", Date.now(), `${result.id}`)
        res.clearCookie('token');
        if(req.body.redirect_url)
            return res.redirect(req.body.redirect_url)
        else
            return res.redirect("/")
        
    } catch (e) {
        if(e.message.includes('jwt expired'))
            return res.send({error: true, message: "JWT_EXPIRED"})
        console.log(e)
        return res.send({
            error: true,
            code: 500,
            message: "SOMETHING_WENT_WRONG",
            errorStack: `${e.message} ${e.stack}`
        })
    }
})

module.exports = Router;