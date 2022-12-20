const mongoose = require('mongoose');
module.exports=mongoose.model('apps', mongoose.Schema({
    id: {type: String, required: true, default: require('crypto').randomBytes(8).toString('hex')},
    name: {type: String, required: true},
    domain: {type: String, required:true},
    redirect_uri: {type: String, required: true}
}))