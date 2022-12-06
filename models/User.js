const { customAlphabet } = require('nanoid')
const nanoid = customAlphabet('0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz', 16)

const mongoose = require('mongoose');
module.exports = mongoose.model('Users', mongoose.Schema({
    id: { type: String, default: () => nanoid() },
    avatar: { type: String, default: '' },
    /* local: {
        username: { type: String },
        password: { type: String, default: "" },
        email: { type: String },
    }, */
    github: {
        username: { type: String },
        email: { type: String },
        avatar: { type: String } ,
        access_token: { type: String, default: "" },
        linked: { type: Boolean, default: false }, 
    },
    google: { 
        username: { type: String }, 
        family_name: { type: String },
        email: { type: String },
        avatar: { type: String },
        access_token: { type: String },
        refresh_token: { type: String },
        linked: { type: Boolean, default: false }       
    },
    discord: { 
        username: { type: String },
        discriminator: { type: String },
        email: { type: String },
        avatar: { type: String },
        linked: { type: Boolean, default: false },
        access_token: { type: String },
        refresh_token: { type: String }
    }
}, { timestamps: true }))