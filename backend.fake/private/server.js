const express = require("express");
const cookieParser = require('cookie-parser')
const app = express();
const cors = require('cors');

const path = require('path');

const port = process.env.port || 6900;

const increment = 1;

var analytics = {
    total: 0,
    totalByUsers: [
        /* just as a reference of what objects should look like
        {
            name: "",
            count: 0
        }
        */
    ]
}

const userSelection = ["willow", "jowanne", "rose", "lauree", "stayo", "boeuf", "raph", "goulat"];

class FakeBackend {
    start() {
        app.use(express.json())
        app.use(cors())
        app.use(cookieParser())

        app.get("/api/bonneFeteRaph", (req, res) => {
            res.setHeader('Content-Type', 'application/json')
            res.end(JSON.stringify(analytics))
        });

        app.post("/api/bonneFeteRaph", (req, res) => {
            this.updateReport(req.body.userId, 1); // in the real backend, instead of reading body, we should just read from the jwt or user creds
            console.log(analytics)

            res.setHeader('Content-Type', 'application/json');
            res.end(JSON.stringify(analytics))
        })


        app.post("/api/auth", (req, res) => {       
            let username = req.body.username
            let password = req.body.password

            if (!(userSelection.includes(username) && password === "corsaire")) {
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
        let jwtExpiration = new Date(Date.now());
        jwtExpiration.setMinutes(jwtExpiration.getUTCMinutes() + 10);
        return jwtExpiration
    }

    updateReport(userId) {
        if (userId == undefined) {
            userId = "unknown";
        }

        analytics.total += increment;

        let found = false;
        analytics.totalByUsers.every(function(userTotal, i) {
            if (userTotal.name === userId) {
                analytics.totalByUsers[i].count += increment;
                found = true;

                return false; // break
            }

            // remember that we didn't find anything so we can create a new object
            if (i === analytics.totalByUsers.length - 1) {
                found = false;
            }
            return true; // continue
        })

        if (!found) {
            analytics.totalByUsers.push({
                name: userId,
                count: increment
            })
        }
    }
}

module.exports = FakeBackend;