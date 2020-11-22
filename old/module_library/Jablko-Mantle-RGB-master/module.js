// mantle_rgb_light.js: RGB Mantle Light Controller
// Cale Overstreet
// August 20, 2020
// Controls RGB lights through an ESP8266
// Exports: 

const fs = require("fs").promises;
const path = require("path");
const fetch = require("node-fetch");

const module_name = path.basename(__dirname);
const module_config = require(module.parent.filename).jablko_config.jablko_modules[module_name];

// Check if module_config is correct
if (module_config.controller_ip == null || module_config.controller_ip == undefined) {
	throw new Error("Incorrect Configuration");
}

module.exports.permission_level = 0;

module.exports.generate_card = async () => {
	return (await fs.readFile(`${__dirname}/mantle_rgb.html`, "utf8")).replace(/\$MODULE_NAME/g, module_name);
}

module.exports.status = async (req, res) => {
	await fetch(`http://${module_config.controller_ip}/status`)
		.then(async function(data) {
			const response = await data.json();
			if (data.status == 200) {
				res.json({status: "good", message: "Set RGBA", r: response.r, g: response.g, b: response.b, a: response.a});
			} else {
				throw new Error("Error communicating with controller");
			}
		})
		.catch(function(error) {
			console.log("Mantle RGB Light: Error communicating with controller");
			console.debug(error);

			res.json({status: "fail", message: `Error contacting controller (HTTP ERROR)`});
			return;
		});
}

module.exports.set_rgba = async (req, res) => {
	const raw_response = await fetch(`http://${module_config.controller_ip}/set_rgba`, {
		method: "POST",
		headers: {
			"Accept": "application/json",
			"Content-Type": "application/json"
		},
		body: JSON.stringify(req.body)
	})
		.then(async function(data) {
			if (data.status == 200) {
				res.json({status: "good", message: "Set RGBA"});
			} else {
				throw new Error("Error communicating with controller");
			}
		})
		.catch(function(error) {
			console.log("Mantle RGB Light: Error communicating with controller");
			console.debug(error);

			res.json({status: "fail", message: `Error contacting controller (HTTP ERROR)`});
			return;
		});
}
