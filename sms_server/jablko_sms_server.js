// jablko_sms_server.js: Jablko Server responsible for sending messages to smart home residents
// Cale Overstreet, Corinne Gerhold
// Uses a gmail account to send messages to residents

// Package requires
const express = require("express");
const nodemailer = require("nodemailer");

// Credentials
const credentials = require("../jablko_sms_config.json");

// Create Nodemailer transporter
const transporter = nodemailer.createTransport({
	service: "gmail",
	auth: {
		user: credentials.email,
		pass: credentials.password
	}
});

const mail_options = {
	from: "jablkohome@gmail.com",
	to: "6156892522@mms.att.net",
	subject: "",
	text: "I Love You!"
};

transporter.sendMail(mail_options, function(err, info) {
	console.log(info);
});

var server = express();

server.listen(10231, function() {
	console.log("Jablko SMS Interface started on port 10231");
	console.log(credentials);
});
