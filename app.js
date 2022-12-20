require('dotenv').config();
const cors = require('cors')
const express = require('express');
const fs = require('fs');
const Apps = require('./models/Apps');
const Certs = require('./models/Certs');

const app = express();

require('./databases')
app.use(require('cookie-parser')())
app.use(cors())
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(express.static('public'))


app.get('/', (req, res) => res.send({message: "Xixo Auth Server v1.0"}))
app.get('/keys', (req, res) => { res.send(require('./certs/keys.json')) })

app.all('/login', async (req, res) => {
  let oAuthLinks = {
    "github": `https://github.com/login/oauth/authorize?client_id=<client_id>&scope=user%20user:email%20repo%20repo_deployment&redirect_uri=<redirect_uri>`,
    "discord": `https://discord.com/oauth2/authorize?client_id=<client_id>&response_type=code&scope=identify%20guilds.join%20email&prompt=consent&redirect_uri=<redirect_uri>`,
    "google": `https://accounts.google.com/o/oauth2/v2/auth?client_id=<client_id>&prompt=consent&response_type=code&scope=profile+email&access_type=offline&redirect_uri=<redirect_uri>` }

  if(!req.query.type || !oAuthLinks[req.query.type]) 
    return res.send({error: true, message: "INVALID_LOGIN_TYPE"})

  if(!req.query.app)
    return res.send({error: true, message: "INVALID_REDIRECT_URI"})
  
  const app = await Apps.findOne({id: req.query.app})
  if(!app) return res.send({error: true, message: "INVALID_REDIRECT_URI"})

  let credentials = {
    "github": {client_id: process.env.GITHUB_CLIENT_ID},
    "discord": {client_id: process.env.DISCORD_CLIENT_ID},
    "google": {client_id: process.env.GOOGLE_CLIENT_ID}
  }

  let link = oAuthLinks[req.query.type]
  link = link.replace("<client_id>", credentials[req.query.type].client_id)
  link = link.replace("<redirect_uri>", 
      `${process.env.NODE_ENV === "production" ? process.env.GLOBAL_URL : process.env.LOCAL_URL}/callback/${req.query.type}?app=${req.query.app}`)

  if(req.query.res === "json")
    return res.send({link})

  return res.redirect(link)
})

app.use('/callback/github', require('./routes/callback/github'))
app.use('/callback/discord', require('./routes/callback/discord'))
app.use('/callback/google', require('./routes/callback/google'))
app.use('/user/info', require('./routes/user/userinfo'))
app.use('/user/logout', require('./routes/user/logout'))

app.listen(process.env.PORT || 3000, async () => { 
  console.log(`🚀 @ http://localhost:${process.env.PORT || 3000}`)
  if(!fs.existsSync('./certs')) fs.mkdirSync('./certs');
  if(!fs.existsSync('./certs/private.pem') || !fs.existsSync('./certs/public.pem')){
    const cert = await Certs.findOne({});
    fs.writeFileSync('./certs/private.pem', cert.private);
    fs.writeFileSync('./certs/public.pem', cert.public)
    app.set('secret', cert.private);
  } else {
    app.set('secret', fs.readFileSync('./certs/private.pem').toString());
  }
  if(!fs.existsSync('./certs/keys.json')) fs.writeFileSync('./certs/keys.json', JSON.stringify({}))
  const Keys = require('./certs/keys.json')
  Keys.keys=[]
  const Rasha = require('rasha')
  Rasha.import({pem: app.get('secret'), public: true}).then((val) => {
    val.use = 'sig'
    Keys.keys.push(val)
    fs.writeFileSync('./certs/keys.json', JSON.stringify(Keys))
  })
  
});
