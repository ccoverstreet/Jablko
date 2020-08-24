// announcements.js: Used for sending announcements in GroupMe chat
// Cale Overstreet
// August 20, 2020
// Used for fun
// Exports: permission_level, async generate_card(), async send_announcement(req, res)

const fs = require("fs").promises;
const path = require("path");

const messaging_system = require("../../jablko_interface.js").messaging_system;

const module_name = path.basename(__dirname);


module.exports.permission_level = 1;

module.exports.generate_card = async () => {
	return (await fs.readFile("./jablko_modules/announcements/announcements.html", "utf8")).replace("$MODULE_NAME", module_name);
}

module.exports.send_announcement = async (req, res) => {
	await messaging_system.send_message(req.body.message)
		.then((data) => {
			console.log("Posted announcement to GroupMe");
			res.json({status: "good", message: "Sent announcement"});
		})
		.catch((error) => {
			console.log("Unable to post announcement to GroupMe");
			console.log(error);
			res.json({status: "fail", message: "Couldn't send announcement"});
		})
}
