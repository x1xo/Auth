const Redis = require('ioredis')
const mongoose = require('mongoose');
const { removeInvalidToken } = require('./Utils');
const redis = new Redis(process.env.REDIS_URL)

redis.on('connect', async () => {
    console.log('[Database] Connected to Redis')
    setInterval(async () => {
        const invalid_tokens = await redis.zrangebyscore("logouts", 0, Date.now()-24*60*60*1000)
        await redis.zrem("logouts", invalid_tokens)
    }, 24*60*60*1000)
})

global.redis = redis;


mongoose.connect(process.env.MONGODB_URL, {}, (err) => {
    if(err) {
        console.log('[Database] [MongoDB]' + err)
        process.exit(1)
    }
    console.log('[Database] Connected to MongoDB')
})
