const axios = require('axios');
const jwt = require('jsonwebtoken');
const User = require('../../models/User');
const { removePrivateData } = require('../../Utils');
const Router = require('express').Router();


Router.get('/', async (req, res) => {
    try {
        //gres -> google-response
        const {data: gres} = await axios.post('https://oauth2.googleapis.com/token', {
            client_id: process.env.GOOGLE_CLIENT_ID,
            client_secret: process.env.GOOGLE_CLIENT_SECRET,
            code: req.query.code,
            grant_type: 'authorization_code',
            redirect_uri: `${process.env.NODE_ENV === 'production' ? process.env.GLOBAL_URL : process.env.LOCAL_URL}/callback/google`,
        }, { headers: {'Content-Type': 'application/x-www-form-urlencoded', 'Accept-Encoding': 'identity'}})
        //console.log('google token data:', gres)
        //ures -> google user response
        const {data: ures} = await axios.get(`https://www.googleapis.com/oauth2/v2/userinfo?access_token=${gres.access_token}`, { headers: {'Accept-Encoding': 'identity'}})
        //console.log('google user data:', ures)
        let { email } = ures;
        let username = ures.name;
        if(!username)
            username = ures.email.split('@')[0]
            
        let mongoUserSearch = await User.findOne({$or: [
            {'github.email': email},  //search for user that has github strategy with this email
            {'google.email': email},  // search for user that has google strategy with this email
            {'discord.email': email}  //search for user that has discord strategy with this email
        ]})

        if(mongoUserSearch) {
            const updatedUser = await User.findOneAndUpdate({id: mongoUserSearch.id}, {google: { 
                username, email, avatar: ures.picture, access_token: gres.access_token,
                refresh_token: gres.refresh_token, linked: true}}, {new: true})
            
            const jwt_user = removePrivateData(updatedUser.toJSON());
            let token = jwt.sign(jwt_user, req.app.get('secret'), {algorithm: 'RS256', expiresIn: '10m'})
            res.cookie('token', token, {httpOnly: true, maxAge: Date.now()+24*60*60*1000})
            return res.send({token})
        }

        const user = new User({
            username,
            email,
            avatar: ures.picture,
            google: {
                username,
                familiy_name: ures.name,
                email,
                avatar: ures.picture,
                access_token: gres.access_token,
                refresh_token: gres.refresh_token,
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