import { readFileStr } from "../../src/util.ts";
import { module_config } from "./config.ts";

export function permission_level() {
	return 1;
}

export async function generate_card() {
	return await readFileStr("./jablko_modules/mantle_rgb_light/mantle_rgb_light.html");
}

export async function set_rgba(context: any) {
	const raw_response = await fetch(`http://${module_config.controller_ip}/set_rgba`, {
		method: "POST",
		headers: {
			"Accept": "application/json",
			"Content-Type": "application/json"
		},
		body: JSON.stringify(context.json_data)
	})
		.catch(function(error) {

		});
	
		/*
	if (raw_response.status < 200 && raw_response.status >= 300) {
		// Error in contacting controller
		context.response.type = "json";
		context.response.body = {status: "fail", message: `Error contacting controller (HTTP ERROR ${raw_response.status})`}
	}
   */
	//console.log(await raw_response.json());
	context.response.type = "json";
	context.response.body = {status: "good", message: "Tried to update RGB"};
}
