// jablko_weather.js: Jablko Weather
// Cale Overstreet
// August 20, 2020
// Interface for the Jablko Weater utility
// Exports: permission_level, async generate_card(), async get_current_weather(req, res), async get_hourly_weather(req, res), async get_daily_weather(req, res), async get_all_weather(req, res)

const fs = require("fs").promises;
const path = require("path");

const module_name = path.basename(__dirname);
const weather = require(module.parent.filename).weather;

module.exports.permission_level = 0;

module.exports.generate_card = async () => {
	return (await fs.readFile(`${__dirname}/jablko_weather.html`, "utf8"))
		.replace(/\$MODULE_NAME/g, module_name);
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

// -------------------- Chatbot Functions --------------------

function random_int(max) {
	return Math.floor(Math.random() * max);
}

module.exports.chatbot_current_weather = async (req, res) => {
	const responses = [
		"Here's the weather.",
		"Just look outside.",
		"Coming up with weather puns is a breeze!",
		"What falls but never hits the ground? The temperature."
	];

	const weather_data = await weather.get_current_weather();

	const temp_in_c = weather_data.current.temp - 273.15;
	const feels_in_c = weather_data.current.feels_like - 273.15;

	var message = responses[random_int(responses.length)];
	message += `
Temp: ${((9 / 5 * temp_in_c) + 32).toFixed(1).toString() + " °F (" + temp_in_c.toFixed(1) + " °C)"}
Feels: ${((9 / 5 * feels_in_c) + 32).toFixed(1).toString() + " °F (" + feels_in_c.toFixed(1) + " °C)"}
Humidity: ${weather_data.current.humidity}%,
Desc: ${weather_data.current.weather[0].description}
`;

	return message;
}
