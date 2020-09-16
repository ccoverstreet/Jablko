// ip_logger.js: IP Logger
// Cale Overstreet
// September 15, 2020
// Logs IP addresses to ip_addresses.log in root of Jablko. Used for security monitoring

const fs = require("fs").promises;

var ip_addresses = {}
var access_counter = 0;

module.exports.ip_logger_middleware = async (req, res, next) => {
	if (req.connection.remoteAddress in ip_addresses) {
		ip_addresses[req.connection.remoteAddress]++;
	} else {
		ip_addresses[req.connection.remoteAddress] = 1;
		write_log();
		access_counter = 0;
	}

	if (access_counter > 10) {
		write_log();
		access_counter = 0;
	}

	access_counter++;
	console.log(access_counter);
	console.log(ip_addresses);

	await next();	
}

function write_log() {
	fs.writeFile("./ip_addresses.json", JSON.stringify(ip_addresses, null, 4))
		.catch((error) => {
			console.log("Unable to write to ip_addresses.log");
			console.debug(error);
		});
}
