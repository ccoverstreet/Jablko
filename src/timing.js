// timing.js: Jablko Response Handling Timer Middleware
// Cale Overstreet
// August, 19, 2020
// Simple timing function that is triggered at the start of an Express response

const jablko = require("../jablko_interface.js");

var handling_times = {
	window_size: 100,
	n: 0,
	current_average: 0 
};

module.exports.timing_middleware = async function(req, res, next) {
	const start_time = Date.now();
	await next()
	if (handling_times.n == handling_times.window_size) {
		handling_times.n = 0;
		handling_times.current_average = 0;
	} 

	const interval_time = Date.now() - start_time;
	console.debug(`Request took ${interval_time} ms`);

	handling_times.current_average = (handling_times.current_average * handling_times.n + interval_time) / (handling_times.n + 1);
	handling_times.n++;	
}

module.exports.get_handling_time = () => {
	return handling_times.current_average;
}
