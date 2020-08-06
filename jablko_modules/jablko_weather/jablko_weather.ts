// jablko_weather.ts: Jablko Modules Dashboard Weather
// Cale Overstreet
// July 7, 2020
// Dashboard implementation of Open Weather Map API interface.

import { readFileStr } from "https://deno.land/std@0.61.0/fs/mod.ts";

const weather = await import("../../src/weather.ts");

export async function generate_card() {
	return await readFileStr("./jablko_modules/jablko_weather/jablko_weather.html");
}

export function permission_level() {
	return 0;
}

export async function get_current_weather(context: any) {
	context.response.type = "json";
	context.response.body = await weather.get_current_weather();
}
