// ip_logger.js: IP Logger
// Cale Overstreet
// September 15, 2020
// Logs IP addresses to ip_addresses.log in root of Jablko. Used for security monitoring

const fs = require("fs").promises;


var ip_addresses = {}

try {
	ip_addresses = require("../log/ip_addresses.json");
} catch(error) {
	console.log("Error getting ip addresses");
	console.debug(error);
}

console.log(ip_addresses);

var access_counter = 0;

const jablko = require("../jablko_interface.js");

module.exports.ip_logger_middleware = async (req, res, next) => {
	if (req.connection.remoteAddress in ip_addresses) {
		ip_addresses[req.connection.remoteAddress]++;
	} else {
		ip_addresses[req.connection.remoteAddress] = 1;

		jablko.messaging_system.send_message(`New access from ip "${req.connection.remoteAddress}"`);
		write_log();
		access_counter = 0;
	}

	if (access_counter > 10) {
		write_log();
		access_counter = 0;
	}

	access_counter++;

	await next();	
}

function write_log() {
	fs.writeFile("./log/ip_addresses.json", JSON.stringify(ip_addresses, null, 4))
		.catch((error) => {
			console.log("Unable to write to ip_addresses.log");
			console.debug(error);
		});
}
