const Redis = require('ioredis')
const mongoose = require('mongoose');
const { removeInvalidToken } = require('./Utils');
const redis = new Redis(process.env.REDIS_URL)

global.invalidTokenTimeout = setTimeout(() => {}, 0)
redis.on('connect', async () => {
    console.log('[Database] Connected to Redis')
    const lindex = await global.redis.lindex("invalid_tokens", -1);
    removeInvalidToken(lindex)
})
global.redis = redis;


mongoose.connect(process.env.MONGODB_URL, {}, (err) => {
    if(err) {
        console.log('[Database] [MongoDB]' + err)
        process.exit(1)
    }
    console.log('[Database] Connected to MongoDB')
})
