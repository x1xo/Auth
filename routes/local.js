const axios = require('axios');
const bcrypt = require('bcrypt')
const User = require('../models/User');
const Router = require('express').Router();

Router.get('/', async (req, res) => {
    try {
        const {username, email, password} = req.body

        //check if user exists
        let mongoUserSearch = await User.findOne({$or: [
            {local: { email }},     //search for user that has local strategy with this email 
            {local: { username }},  //search for user that has local strategy with this username
            {github: { email }},    //search for user that has github strategy with this email
            {google: { email }},  // search for user that has google strategy with this email
            {discord: { email }}  //search for user that has discord strategy with this email
        ]})
        if(mongoUserSearch) {
            if(mongoUserSearch.local.password) {
                bcrypt.compare(password, mongoUserSearch.local.password, (err, same) => {
                    if(err) throw new Error(...err) //This will be handled with the catch at the end
                    if(!same) return res.send({
                        error: true,
                        code: 403,
                        message: "INVALID_PASSWORD"
                    })

                    req.session.user = { ...mongoUserSearch.toObject() }
                    return res.send({
                        error: false,
                        code: 200, 
                        message: "SUCCESS"
                    })
                })
            } else { 
                return res.send({
                    error: true,
                    code: 403,
                    message: "WRONG_AUTH_PROVIDER"
                })
            }            
        }
        return res.send({error: true, code: 404, message: "USER_NOT_FOUND"})
    } catch(e) {
        return res.send({
            error: true,
            code: 500,
            message: "SOMETHING_WENT_WRONG",
            errorStack: `${e.message} ${e.stack}`
        })
    }
    
})

Router.post('/', async (req, res) => {
    try { 
        const {username, email, password} = req.body

        //check if user exists
        let mongoUserSearch = await User.findOne({$or: [
            {local: { email }},     //search for user that has local strategy with this email 
            {local: { username }},  //search for user that has local strategy with this username
            {github: { email }},    //search for user that has github strategy with this email
            {google: { email }},  // search for user that has google strategy with this email
            {discord: { email }}  //search for user that has discord strategy with this email
        ]})
        if(mongoUserSearch){
            if(mongoUserSearch.local.password)
                return res.send({error: true, code: 403, message: "USER_ALREADY_EXISTS"})
            if(mongoUserSearch.github.linked || mongoUserSearch.discord.linked || mongoUserSearch.google.linked) 
                return res.send({error: true, code: 401, message: "WRONG_AUTH_PROVIDER"})
        }

        let passwordHash;
        bcrypt.hash(myPlaintextPassword, 12, function(err, hash) {
            if(err || !hash) throw new Error('Something went wrong while caching.')
            if(hash) hash = passwordHash
        });

        let user = new User({
            local: {
                username, 
                email,
                password: passwordHash
            }  
        });
        await user.save()

        req.session.user = { ...user.toObject() }
        return res.send({
            error: false,
            code: 200,
            message: "SUCCESS"  
        })
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