// timing.js: Jablko Response Handling Timer Middleware
// Cale Overstreet
// August, 19, 2020
// Simple timing function that is triggered at the start of an Express response

const jablko = require("../jablko_interface.js");

module.exports.timing_middleware = async function(req, res, next) {
	const start_time = Date.now();
	await next()
	console.log(`${Date.now() - start_time} ms`);
}
