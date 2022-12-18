const mongoose = require('mongoose');
module.exports = mongoose.model('Users', mongoose.Schema({
    id: { type: String, default: () => require('crypto').randomBytes(8).toString('hex') },
    username: { type: String }, 
    email: { type: String },
    avatar: { type: String, default: '' },
    admin: { type: Boolean, default: false },
    /* local: {
        username: { type: String },
        password: { type: String, default: "" },
        email: { type: String },
    }, */
    github: {
        username: { type: String },
        email: { type: String },
        avatar: { type: String } ,
        access_token: { type: String },
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