import { readFileStr } from "../../src/util.ts";
import { module_config } from "./config.ts";
/* Should contain:
 * export const module_config = {
 *     controller_ip: 10.0.0.2
 * }
 */

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
		.then(async function(data) {
			if (data.status == 200) {
				context.response.type = "json";
				context.response.body = {status: "good", message: "Set RGBA"};
			} else {
				throw new Error("Error communicating with controller");
			}
		})
		.catch(function(error) {
			console.log("Mantle RGB: Error communicating with controller");
			console.log(error);

			context.response.type = "json";
			context.response.body = {status: "fail", message: `Error contacting controller (HTTP ERROR)`}
			return;
		});
}
