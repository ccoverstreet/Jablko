// jablko_sms_server.js: Jablko Server responsible for sending messages to smart home residents
// Cale Overstreet, Corinne Gerhold
// Uses a gmail account to send messages to residents

// Package requires
const express = require("express");
const bodyParser = require("body-parser");
const nodemailer = require("nodemailer");

// Read credentials from config file
const credentials = require("../jablko_sms_config.json");

// Domains for certain phone providers
carriers = {
    att:    'mms.att.net',
    tmobile: 'tmomail.net',
    verizon:  'vtext.com',
    sprint:   'page.nextel.com'
};

// Create Nodemailer transporter
const transporter = nodemailer.createTransport({
	service: "gmail",
	auth: {
		user: credentials.email,
		pass: credentials.password
	}
});

var server = express();
server.use(bodyParser.json());

server.post("/send_message", function(req, res) {
	// Change this to just get username and message from req. SQLite query will be done in here
	const to_address = `${req.body.number}@${carriers[req.body.carrier]}`;

	const mail_options=  {
		from: credentials.email,
		to: to_address,
		subject: "",
		text: req.body.message
	};

	transporter.sendMail(mail_options, function(err, info) {
		if (err) {
			console.log(`Unable to send message to ${to_address}`);
			res.json({status: "fail", message: "Unable to send message to ${to_address}"});
		} else {
			res.json({status: "good", message: "Sent message"});
		}
	});
});

server.get("/status", function(req, res) {
	res.json({status: "good", message: "SMS Server Operational"});
});

server.listen(10231, function() {
	console.log("Jablko SMS Interface started on port 10231");
});
