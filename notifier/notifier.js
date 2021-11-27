const express = require("express");
const app = express();
const bodyParser = require("body-parser");
const port = process.env.PORT || 8000;
const notifier = require("node-notifier");
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.get("/health", (_, res) => {
    res.status(200).send();
});
app.post("/notify", (req, res) => {
    notify(req.body, (reply) => res.send(reply));
});
app.listen(port, () => console.log(`Listening on port ${port}`));

const notify = ({ title, message }, cb) => {
    notifier.notify(
        {
            title: title || "Notification",
            message: message || "Unknown Message",
            sound: true,
            wait: true,
            reply: true,
            closeLabel: "Completed?",
            timeout: 15,
        },
        (_, __, reply) => {
            cb(reply);
        }
    );
};
