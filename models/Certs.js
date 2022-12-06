const mongoose = require('mongoose');
module.exports=mongoose.model('certs', new mongoose.Schema({
    private: {type:String},
    public: {type:String}
}))