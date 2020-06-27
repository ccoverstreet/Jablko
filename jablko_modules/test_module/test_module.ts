// test_module.ts: Test Module
// Cale Overstreet
// June 26, 2020
// This is a testbed for new Jablko Modules features

import { readFileStr } from 'https://deno.land/std/fs/mod.ts';

const smtp_client = (await import("../../jablko_interface.ts")).smtp_client;

export const info = {
	permissions: "All"
};

export async function generate_card() {
	return await readFileStr("jablko_modules/test_module/test_module.html");
}

export async function send_message(context: any) {
	const send_status = await smtp_client.send_message(context.json_data.username, context.json_data.message);

	if (send_status == false) {
		context.response.type = "json";
		context.response.body = {
			status: "fail"
		};
	} else {
		context.response.type = "json";
		context.response.body = {
			status: "good"
		};
	}
}
