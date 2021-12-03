const express = require("express");
const app = express();

const path = require('path')

const port = process.env.port || 6900;

class FakeBackend {
    start() {
        app.get("/api/bonneFeteRaph", (req, res) => {
            res.setHeader('Content-Type', 'application/json')
            res.end(JSON.stringify({nbBonneFete: 200}))
        });

        app.listen(port, err => {
            if (err) {
                console.log("ERROR", err);
            }
            console.log(`listening on port ${port}`)
        });
    }
}

module.exports = FakeBackend;