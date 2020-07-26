// bot_brain.ts: GroupMe Bot message callback handler
// Cale Overstreet
// July 24, 2020
// This file describes the behavior of the GroupMe bot. In the future it could use language comprehension to determine actions, but I'll just use some general rules for now.
// Exports: handle_message

const messaging_system = (await import("../jablko_interface.ts")).messaging_system;

export async function handle_message(context: any) {
	context.json_content.text = context.json_content.text.toLowerCase();
	console.log(context.json_content.text);

	if (context.json_content.text.includes("@jablko")) {
		await determine_response(context);
	}

	context.response.type = "html"
	context.response.body = "";
}

async function determine_response(context: any) {
	context.json_content.split_text = context.json_content.text.split(" ");
	for (var i = 0; i < context.json_content.split_text.length; i++) {
		if (context.json_content.split_text[i] == "hello") {
			messaging_system.send_message(`What's up ${context.json_content.name}`);
		}
	}
}
