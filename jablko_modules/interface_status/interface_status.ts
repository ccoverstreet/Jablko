// Jablko Modules: Server Status Module // Cale Overstreet
// May 3, 2020
// This module serves as an example for creating custom modules and provides information on the smart home dashboard related to the run condition, uptime, temperature, and memory usage

import { Context } from "https://deno.land/x/oak/mod.ts";
import { request_handling_times } from "../../src/timing.ts";
import { server_start_time } from "../../jablko_interface.ts";
import { readFileStr } from "https://deno.land/std/fs/mod.ts";

export const info = {
	permissions: "all"
}

export async function generate_card() {
	return await readFileStr("jablko_modules/interface_status/interface_status.html");
}

export async function check_status(context: Context) {
	// This function reads necessary data and sends results back to client
	
	// Read CPU temp
	const temperature = parseFloat(new TextDecoder().decode(await Deno.readFile("/sys/class/thermal/thermal_zone0/temp"))) / 1000;

	const meminfo = new TextDecoder().decode(await Deno.readFile("/proc/meminfo")).split("\n");
	const total_mem = (parseFloat(meminfo[0].split(/[ ]+/)[1]) / 1000000);
	const meminfo_summary = `${(total_mem - parseFloat(meminfo[2].split(/[ ]+/)[1]) / 1000000).toFixed(2)} / ${total_mem.toFixed(2)} GB`;

	const raw_uptime = (new Date().getTime() - server_start_time) / 1000;
	const hours = Math.floor(raw_uptime / 3600);
	const minutes = Math.floor((raw_uptime - hours) / 60);
	const seconds = Math.floor(raw_uptime - hours - minutes);
	const formatted_uptime = `${hours} h ${minutes} m ${seconds}s`;
   
	context.response.type = "json";
	context.response.body = {
		interface_status: "good",
		interface_uptime: formatted_uptime,
		interface_response_time: request_handling_times.current_average.toFixed(3) + " ms",
		cpu_temperature: temperature + " C",
		memory_usage: meminfo_summary
	};
}
