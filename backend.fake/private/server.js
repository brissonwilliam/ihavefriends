const express = require("express");
const { STATUS_CODES } = require("http");
const app = express();

const path = require('path')

const port = process.env.port || 6900;

var c = 200;

class FakeBackend {
    start() {
        app.use(express.json())

        app.get("/api/bonneFeteRaph", (req, res) => {
            res.setHeader('Content-Type', 'application/json')
            res.end(JSON.stringify({nbBonneFete: c}))
        });

        app.post("/api/bonneFeteRaph", (req, res) => {
            c++;
            res.setHeader('Content-Type', 'application/json')
            res.end(JSON.stringify({nbBonneFete: c}))
        })


        app.post("/api/auth", (req, res) => {            
            let username = req.body.user
            let password = req.body.password

            if (username != "fake" || password != "backend") {
                res.sendStatus(403)
                return
            }

        
            res.setHeader('Content-Type', 'application/json')
            res.end(JSON.stringify({username: "FakeBackend", jwt: "abc123", jwtExpiration: this.getJwtExpirationDate()}))
        })

        app.listen(port, err => {
            if (err) {
                console.log("ERROR", err);
            }
            console.log(`listening on port ${port}`)
        });
    }

    getJwtExpirationDate() {
        let jwtExpiration = new Date(date);
        jwtExpiration.setDate(jwt.getDate() + 6);
        return jwtExpiration
    }
}

module.exports = FakeBackend;