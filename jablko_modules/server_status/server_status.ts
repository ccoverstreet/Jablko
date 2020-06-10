// Jablko Modules: Server Status Module
// Cale Overstreet
// May 3, 2020
// This module serves as an example for creating custom modules and provides information on the smart home dashboard related to the run condition, uptime, temperature, and memory usage

import { Context } from "https://deno.land/x/oak/mod.ts";
import { request_handling_times, server_start_time } from "../../jablko_interface.ts"

export const info = {
	permissions: "all"
}

export function generate_card() {
	return `
<script> 
	server_status = {
		check_status: async function() {
			const interface_values = document.querySelectorAll("#server_status_card>div>.value");
			const response = await (await fetch("/jablko_modules/server_status/check_status", {method: "POST"})
			.catch(error => {
				console.log(error);
				for (var i = 0; i < interface_values.length; i++) {
					interface_values[i].style.color = "#c23616";
					interface_values[i].textContent = "N/A";
				}
				return;
			})).json();

			const response_keys = Object.keys(response);
			console.log(interface_values);

			console.log(response);

			for (var i = 0; i < interface_values.length; i++) {
				interface_values[i].style.color = "#44bd32";
				interface_values[i].textContent = response[response_keys[i]];
			}
		}

	}	
	
	setTimeout(server_status.check_status, 1000);
	setInterval(server_status.check_status, 15000);
</script>
<div id="server_status_card" class="jablko_module_card">
	<div class="card_title" style="background: url('/icons/server_status_icon.svg') right; background-size: contain; background-repeat: no-repeat;">Server Status</div>
	<hr>
	<div class="label_value_pair">
		<div class="label">Interface Status:</div>
		<div id="interface_status" class="value">N/A</div>
	</div>
	<div class="label_value_pair">
		<div class="label">Interface Uptime:</div>
		<div id="interface_uptime" class="value">N/A</div>
	</div>
	<div class="label_value_pair">
		<div class="label">Response Time:</div>
		<div id="interface_response_time" class="value">N/A</div>
	</div>
	<div class="label_value_pair">
		<div class="label">CPU Temperature:</div>
		<div id="cpu_temperature" class="value">N/A</div>
	</div>
	<div class="label_value_pair">
		<div class="label">Memory Usage:</div>
		<div id="memory_usage" class="value">N/A</div>
	</div>
</div>
`
}


export async function check_status(context: Context) {
	// This function reads necessary data and sends results back to client
	
	// Read CPU temp
	const temperature = parseFloat(new TextDecoder().decode(await Deno.readFile("/sys/class/thermal/thermal_zone0/temp"))) / 1000;

	const meminfo = new TextDecoder().decode(await Deno.readFile("/proc/meminfo")).split("\n");
	const total_mem = (parseFloat(meminfo[0].split(/[ ]+/)[1]) / 1000000);
	const meminfo_summary = `${(total_mem - parseFloat(meminfo[1].split(/[ ]+/)[1]) / 1000000).toFixed(2)} / ${total_mem.toFixed(2)} GB`;

	const raw_uptime = (new Date().getTime() - server_start_time) / 1000;
	const hours = Math.floor(raw_uptime / 3600);
	const minutes = Math.floor((raw_uptime - hours) / 60);
	const seconds = Math.floor(raw_uptime - hours - minutes);
	const formatted_uptime = `${hours} h ${minutes} m ${seconds}`;
   
	context.response.type = "json";
	context.response.body = {
		interface_status: "good",
		interface_uptime: formatted_uptime,
		interface_response_time: request_handling_times.current_average.toFixed(3) + " ms",
		cpu_temperature: temperature + " C",
		memory_usage: meminfo_summary
	};
}
