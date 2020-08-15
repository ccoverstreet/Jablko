import { readFileStr } from "../../src/util.ts";

export function permission_level() {
	return 1;
}

export async function generate_card() {
	return await readFileStr("./jablko_modules/mantle_rgb_light/mantle_rgb_light.html");
}

export async function set_rgba(context: any) {
	const raw_response = await fetch("http://10.0.0.46/set_rgba", {
		method: "POST",
		headers: {
			"Accept": "application/json",
			"Content-Type": "application/json"
		},
		body: JSON.stringify(context.json_data)
	});
	//console.log(await raw_response.json());
	context.response.type = "json";
	context.response.body = {status: "good", message: "Tried to update RGB"};
}
