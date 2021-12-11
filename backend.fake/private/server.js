const express = require("express");
const cookieParser = require('cookie-parser')
const app = express();
const cors = require('cors');

const path = require('path');
const { userInfo } = require("os");

const port = process.env.port || 6900;

var c = 200;

class FakeBackend {
    start() {
        app.use(express.json())
        app.use(cors())
        app.use(cookieParser())

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
            let username = req.body.username
            let password = req.body.password

            if (!(username === "fake" && password === "backend")) {
                res.statusCode = 401;
                res.send("invalid credentials");
                return
            }

            let userInfo = {
                exp: this.getJwtExpirationDate(),
                jwt: "token123",
                username: "FakeBackend"
            };

            // re-enable if we want to backend to serve cookies rather than let the client handle them
            // let cookieOptions = {
            //     maxAge: userInfo.exp,
            //     httpOnly: true,
            //     signed: false,
            // };
            // res.cookie('user', userInfo, cookieOptions);
            

            res.setHeader('Content-Type', 'application/json');
            res.end(JSON.stringify(userInfo));
        })

        app.listen(port, err => {
            if (err) {
                console.log("ERROR", err);
            }
            console.log(`listening on port ${port}`);
        });
    }

    getJwtExpirationDate() {
        var jwtExpiration = new Date(Date.now());
        jwtExpiration.setMinutes(jwtExpiration.getUTCMinutes() + 10);
        return jwtExpiration
    }
}

module.exports = FakeBackend;