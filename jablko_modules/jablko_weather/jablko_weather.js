// jablko_weather.js: Jablko Weather
// Cale Overstreet
// August 20, 2020
// Interface for the Jablko Weater utility
// Exports: permission_level, async generate_card(), async get_current_weather(req, res), async get_hourly_weather(req, res), async get_daily_weather(req, res), async get_all_weather(req, res)

const fs = require("fs").promises;
const path = require("path");

const weather = require("../../jablko_interface.js").weather;

const module_name = path.basename(__dirname);


module.exports.permission_level = 0;

module.exports.generate_card = async () => {
	return (await fs.readFile("./jablko_modules/jablko_weather/jablko_weather.html", "utf8")).replace("$MODULE_NAME", module_name);
}

module.exports.get_current_weather = async (req, res) => {
	res.json(await weather.get_current_weather());
}

module.exports.get_hourly_weather = async (req, res) => {
	res.json(await weather.get_hourly_weather());
}

module.exports.get_daily_weather = async (req, res) => {
	res.json(await weather.get_daily_weather());
}

module.exports.get_all_weather = async (req, res) => {
	res.json(await weather.get_all_weather());
}
