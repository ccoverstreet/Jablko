// messaging.js: Messaging System for Jablko
// Cale Overstreet
// August 20, 2020
// Uses GroupMe API to send and respond to messages in home groupchat. Replaced previous SMS system
// Exports: async send_message(message)

const fetch = require("node-fetch");

const groupme_config = require("../jablko_interface.js").jablko_config.GroupMe;

module.exports.send_message = async (message) => {
	const raw_response = await fetch("https://api.groupme.com/v3/bots/post", {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({"bot_id": groupme_config.bot_id, "text": message, "attachments": null})
	});
}
