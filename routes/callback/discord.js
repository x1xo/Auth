const axios = require('axios');
const jwt = require('jsonwebtoken');
const User = require('../../models/User');
const { removePrivateData } = require('../../Utils');
const Router = require('express').Router();

Router.get('/', async (req, res) => {
    try {
        //dres -> discord-response
        let {data: dres} = await axios.post('https://discord.com/api/oauth2/token', {
            client_id: process.env.DISCORD_CLIENT_ID,
            client_secret: process.env.DISCORD_CLIENT_SECRET,
            grant_type: 'authorization_code',
            code: req.query.code,
            redirect_uri: `${process.env.NODE_ENV === 'production' ? process.env.GLOBAL_URL : process.env.LOCAL_URL}/callback/discord`,
        }, { headers: { 'Content-Type': 'application/x-www-form-urlencoded', 'Accept-Encoding': 'identity'}})
        //console.log('discord token data:', dres)
        //ures -> discords user response
        let {data: ures} = await axios.get('https://discord.com/api/users/@me', { headers: { Authorization: `Bearer ${dres.access_token}`, 'Accept-Encoding': 'identity' }})
        //console.log('discord user data:', ures)
        let { email } = ures

        if(!ures.verified) return res.send({ error: true, code: 401, message: "UNVERIFIED_DISCORD_EMAIL" })

        let mongoUserSearch = await User.findOne({$or: [
            {'github.email': email},  //search for user that has github strategy with this email
            {'google.email': email},  // search for user that has google strategy with this email
            {'discord.email': email}  //search for user that has discord strategy with this email
        ]})
        
        if(mongoUserSearch) {
            let updatedUser = await User.findOneAndUpdate({id: mongoUserSearch.id}, { discord: {
                linked: true, email, username: ures.username, 
                discriminator: ures.discriminator, avatar: ures.avatar,
                access_token: dres.access_token, refresh_token: dres.refresh_token }}, {new: true})

            const jwt_user = removePrivateData(updatedUser.toJSON());
            let token = jwt.sign(jwt_user, req.app.get('secret'), {algorithm: 'RS256', expiresIn: '1d'})
            res.cookie('token', token, {httpOnly: true, maxAge: Date.now()+24*60*60*1000})
            return res.send({token})
        }
        const user = new User({
            username: ures.username,
            email: email,
            avatar: `https://cdn.discordapp.com/avatars/${ures.id}/${ures.avatar}.png`, 
            discord: {
                username: ures.username,
                discriminator: ures.discriminator,
                email,
                avatar: `https://cdn.discordapp.com/avatars/${ures.id}/${ures.avatar}.png`, 
                access_token: dres.access_token,
                refresh_token: dres.refresh_token,
                linked: true
            }
        }); user.save()

        const jwt_user = removePrivateData(user.toJSON());
        let token = jwt.sign(jwt_user, req.app.get('secret'), {algorithm: 'RS256', expiresIn: '10m'})
        res.cookie('token', token)
        return res.send({token})
    } catch(e) {
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