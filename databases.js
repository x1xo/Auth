const ioredis = require('ioredis')
const mongoose = require('mongoose');
const redis = new ioredis(process.env.REDIS_URL)

redis.on('connect', () => console.log('[Database] Connected to Redis'))
global.redis = redis;

mongoose.connect(process.env.MONGODB_URL, {}, (err) => {
    if(err) {
        console.log('[Database] [MongoDB]' + err)
        process.exit(1)
    }
    console.log('[Database] Connected to MongoDB')
})
