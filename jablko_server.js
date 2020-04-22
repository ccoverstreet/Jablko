// jablko_web_interface.js: NodeJS server that runs the web interface for Jablko
// Cale Overstreet, Corinne Gerhold
// April 18, 2020
// Primary server for serving smart home web interface.

const jablko_root = __dirname; // Root of the Jablko Repo
const jablko_web_root = __dirname + "/public_html"; // Root of web assets

const express = require("express"); // Main express include
const axios = require("axios");

const server = express(); // Creating server object

server.use(express.static(`${jablko_web_root}`));


server.listen(10230, function() {
	console.log("Jablko Web Interface started on port 10230");
});

process.on("message", function(message) {
	console.log(message);
});
