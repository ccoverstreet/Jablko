// jablko_sms_server.js: Jablko Server responsible for sending messages to smart home residents
// Cale Overstreet, Corinne Gerhold
// Uses a gmail account to send messages to residents

// Package requires
const express = require("express");

// Credentials
const credentials = require("../jablko_sms_config.json");

var server = express();

server.listen(10231, function() {
	console.log("Jablko SMS Interface started on port 10231");
	console.log(credentials);
});
