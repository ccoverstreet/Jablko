// test_module.ts: Test Module
// Cale Overstreet
// June 26, 2020
// This is a testbed for new Jablko Modules features

import { readFileStr } from 'https://deno.land/std@0.61.0/fs/mod.ts';

const interface_exports = (await import("../../jablko_interface.ts"));
const messaging_system = interface_exports.messaging_system;


export function permission_level() {
	return 2;
}

export async function generate_card() {
	return await readFileStr("jablko_modules/announcements/announcements.html");
}

export async function send_message(context: any) {
	//const raw_response = await fetch(`https://api.groupme.com/v3/groups/60780309?token=${groupme_config.access_token}`);
	// const response_json = await raw_response.json();
	//console.log(response_json);
	messaging_system.send_message(`Announcement @all: ${context.json_data.message}`);
	context.response.header = "json";
	context.response.body = {status: "good", message: "sent message"};
}
