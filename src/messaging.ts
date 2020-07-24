// messaging.ts: Messaging System for Jablko
// Cale Overstreet
// July 18, 2020
// This module uses the GroupMe public API to send messages in a home groupchat. This represents a move away from SMS as different carriers handle SMS by email differently, resulting in wildly unpredictable behavior.

const groupme_config = (await import("../jablko_interface.ts")).jablko_config.GroupMe;

export async function send_message(message: String) {
	const attachments = await create_attachments(message);
	const raw_response = fetch("https://api.groupme.com/v3/bots/post", {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({"bot_id": groupme_config.bot_id, "text": message, "attachments": attachments})
	});
}

async function create_attachments(message: String) {
	const raw_response = await fetch(`https://api.groupme.com/v3/groups/60780309?token=${groupme_config.access_token}`);
	const response_json = await raw_response.json();

	const split_message = await message.split(" ");
	var attachments = []
	for (var i = 0; i < split_message.length; i++) {
		if (split_message.length > 0 && split_message[i][0] == "@")	{
			// Check if valid user id and add to mentions
			for (var j = 0; j < response_json.response.members.length; j++) {
				if ("@" + response_json.response.members[j].nickname == await split_message[i].replace(":", "")) {
					attachments.push({"type": "mentions", "user_ids": [response_json.response.members[j].user_id], "loci": [[0, 0]]});
				} else if (split_message[i].replace(":", "") == "@all") {
					for (var k = 0; k < response_json.response.members.length; k++) {
						attachments.push({"type": "mentions", "user_ids": [response_json.response.members[k].user_id], "loci": [[0, 0]]});
					}		
				}
			}
		}
	}

	return attachments;
}
