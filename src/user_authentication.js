// user_authentication.js: Jablko User Authentication Middleware
// Cale Overstreet
// August 19, 2020
// Reads from SQLite database and checks if request has the proper authentication.
// Exports: user_authentication_middleware

const fs = require("fs").promises;

const jablko = require("../jablko_interface.js");

module.exports.user_authentication_middleware = async function(req, res, next) {
	console.log("Authenticating request");
	await next();
}
