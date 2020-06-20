// timing.ts: Timing Module
// Cale Overstreet
// June 20, 2020
// This module contains the timing data and middleware.

export var request_handling_times = {
	current_average: 0,
	window_size: 10,
	size_counter: 0
};

export async function timing_middleware(context: any, next: any) {
	// Uses data from request_handling_times to display a windowed average for request handling times
	const request_start_time = new Date().getTime();	
	await next();

	if (request_handling_times.size_counter == request_handling_times.window_size) {
		request_handling_times.size_counter = 0;		
		request_handling_times.current_average = 0;
	} 
	
	request_handling_times.size_counter++;

	const size = request_handling_times.size_counter;
	const request_end_time = new Date().getTime();

	request_handling_times.current_average = request_handling_times.current_average * (size - 1) / size + (request_end_time - request_start_time) / size;
}
