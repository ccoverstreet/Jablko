// Jablko: Messaging System
// Cale Overstreet
// June 19th, 2020
// This file contains the setup and functions for sending sms messages to Smart Home users. 
// Uses the Deno-smtp module. Be wary of updates to the Deno-smtp repo as development is 
// still ongoing and the only branch currently working is master, which can change at any 
// time. As soon as possible, use a standard branch tag for Deno-smtp to avoid system 
// breaking changes

import { SmtpClient } from "https://deno.land/x/smtp/mod.ts";

export async function Jablko_Smtp_Initialize() {
	const client = new Jablko_Smtp();

	const environment_vars = Deno.env.toObject();

	if (environment_vars.JABLKO_SMTP_USERNAME == undefined || environment_vars.JABLKO_SMTP_PASSWORD == undefined) {
		console.log("ERROR: Environment variables for username and password not set. Please set by prepending \"JABLKO_SMTP_USERNAME=username JABLKP_SMTP_PASSWORD=password\" before ./start_jablko.sh");
		Deno.exit(1);
	}

	const connect_config: any = {
		hostname: "smtp.gmail.com",
		port: 465,
		username: environment_vars.JABLKO_SMTP_USERNAME,
		password: environment_vars.JABLKO_SMTP_PASSWORD
	};

	await client.client.connectTLS(connect_config);

	console.log("Created Client");

	return client;
}

const carrier_list: any = {
	"verizon": "vtext.com",
	"att": "txt.att.net"
};

class Jablko_Smtp {
	client: SmtpClient;

	constructor() {
		this.client = new SmtpClient();
	}

	async send_message(user_data: any, message: string) {
		await this.client.send({
			from: "jablkohome@gmail.com",
			to: `${user_data.phone_number}@${carrier_list[user_data.phone_carrier]}`,
			subject: "",
			content: message
		});	
	}
}

