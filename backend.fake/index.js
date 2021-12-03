const FakeBackend = require('./private/server')

console.log(`starting fake backend server`)

const s = new FakeBackend();
s.start();