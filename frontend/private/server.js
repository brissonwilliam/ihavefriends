const express = require("express");
const app = express();
app.set('view engine', 'html')
app.engine('html', require('ejs').renderFile)

const path = require('path')

const port = process.env.port || 80;

let staticOptions = {
    dotfiles: "ignore",
    index: false, // to disable directory indexing
    setHeaders: function(res, path, stat) {
        // add this header to all static responses
        res.set("x-timestamp", Date.now())
    }
};

class Server {
    start() {
        const publicPath = path.join(__dirname, "..", "public")

        app.use(express.static(publicPath, staticOptions));

        app.get("/", (req, res) => {
            res.render(path.join(publicPath, "index.html"))
        });

        app.listen(port, err => {
            if (err) {
                console.log("ERROR", err);
            }
            console.log(`listening on port ${port}`)
        });
    }
}

module.exports = Server;