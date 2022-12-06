const axios = require('axios');
const jwt = require('jsonwebtoken');
const User = require('../../models/User');
const { removePrivateData } = require('../../Utils');
const Router = require('express').Router();

Router.get('/', async (req, res) => {
    try {
        //gres -> github-response
        let {data: gres} = await axios.post("https://github.com/login/oauth/access_token", {
            client_id: process.env.GITHUB_CLIENT_ID,
            client_secret: process.env.GITHUB_CLIENT_SECRET,
            code: req.query.code }, { headers: { 'Accept': 'application/json', 'Accept-Encoding': 'identity' }})
        //ures -> user-response
        //eres -> email-response
        let {data: ures} = await axios.get('https://api.github.com/user', { headers: {Authorization: `Bearer ${gres.access_token}`, 'Accept-Encoding': 'identity'}})
        let {data: eres} = await axios.get('https://api.github.com/user/emails', { headers: {Authorization: `Bearer ${gres.access_token}`, 'Accept-Encoding': 'identity'}})
        
        let { email } = eres.filter(email => email.primary)[0]

        //check if user exists
        let mongoUserSearch = await User.findOne({$or: [
            {'local.email': email},   //search for user that has local strategy with this email 
            {'github.email': email},  //search for user that has github strategy with this email
            {'google.email': email},  // search for user that has google strategy with this email
            {'discord.email': email}  //search for user that has discord strategy with this email
        ]})

        if(mongoUserSearch) {
            const updatedUser = await User.findOneAndUpdate({id: mongoUserSearch.id}, {github: { 
                username: ures.login, email, avatar: ures.avatar_url,
                access_token: gres.access_token, linked: true}}, {new: true})

            const jwt_user = removePrivateData(updatedUser.toJSON());
            let token = jwt.sign(jwt_user, req.app.get('secret'), {algorithm: 'RS256', expiresIn: '10m'})
            res.cookie('token', token)
            return res.send({token})
        }

        const user = new User({
            github: {
                username: ures.login,
                email,
                avatar: ures.avatar_url,
                access_token: gres.access_token,
                linked: true
            }
        }); user.save()

        const jwt_user = removePrivateData(user.toJSON());
        let token = jwt.sign(jwt_user, req.app.get('secret'), {algorithm: 'RS256', expiresIn: '10m'})
        res.cookie('token', token)
        return res.send({token})
    } catch(e) {
        return res.send({
            error: true,
            code: 500,
            message: "SOMETHING_WENT_WRONG",
            errorStack: `${e.message} ${e.stack}`
        })
    }
    
})


module.exports = Router;