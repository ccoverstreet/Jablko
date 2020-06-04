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
		check_status: function() {
			$.post("/jablko_modules/server_status/check_status", function(data) {
				console.log(data);
				const interface_status_element = document.getElementById("interface_status");
				const interface_uptime_element = document.getElementById("interface_uptime");
				const interface_response_time_element = document.getElementById("interface_response_time");

				interface_status_element.textContent = data.interface_status;
				const hour = Math.floor(data.interface_uptime / 3600);
				const minute = Math.floor((data.interface_uptime - hour * 3600) / 60);
				const second = Math.floor((data.interface_uptime - hour * 3600 - minute * 60));
				interface_uptime_element.textContent = hour + " h " + minute + " m " + second + " s";
				interface_response_time_element.textContent = data.interface_response_time.toFixed(3) + " ms";

				switch (data.interface_status) {
					case "good":
						interface_status_element.style.color = "#44bd32";
						interface_uptime_element.style.color = "#44bd32";
						interface_response_time.style.color = "#44bd32";
 						break;
					case "fail":
						interface_status_element.style.color = "#c23616";
						break;
					default:
						interface_status_element.style.color = "#e1b12c";
				}
			})
				.fail(function() {
					const interface_status_element = document.getElementById("interface_status")
					const interface_uptime_element = document.getElementById("interface_uptime")
					const interface_response_time_element = document.getElementById("interface_response_time");

					interface_status_element.textContent = "fail"
					interface_status_element.style.color = "#c23616"

					interface_uptime_element.textContent = "Server Down"
					interface_uptime_element.style.color = "#c23616"

					interface_response_time_element.textContent = "Server Down";
					interface_response_time_element.style.color = "#c23616";
				});
		}
	}	
	
	server_status.check_status();
	setInterval(server_status.check_status, 15000);
</script>
<div class="jablko_module_card">
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
</div>
`
}

export function check_status(context: Context) {
	context.response.type = "json";
	context.response.body = {
		interface_status: "good",
		interface_uptime: (new Date().getTime() - server_start_time) / 1000,
		interface_response_time: request_handling_times.current_average
	};
}
