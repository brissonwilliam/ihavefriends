const Server = require('./private/server')

console.log(`starting static server`)

const s = new Server();
s.start();