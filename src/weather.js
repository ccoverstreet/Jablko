// weather.js: Open Weather Map Wrapper
// Cale Overstreet
// August 20, 2020
// Used by Jablko and its modules to interface with the OWM API.
// Exports: async get_current_weather()

const fetch = require("node-fetch");
const jablko_config = require("../jablko_interface.js").jablko_config;

console.log("Getting lattitude and longitude for weather data...");
var location = undefined;
(async () => {
	const raw_location_fetch = await fetch("http://ip-api.com/json");
	location = await raw_location_fetch.json();
	console.log(`Location: ${location.lat}, ${location.lon}`);
})();

const owm_prefix = "https://api.openweathermap.org/data/2.5/onecall?";

module.exports.get_current_weather = async () => {
	const raw_weather_data = await fetch(`${owm_prefix}lat=${location.lat}&lon=${location.lon}&appid=${jablko_config.weather.key}&exclude=minutely,daily,hourly`);
	return await raw_weather_data.json();
}

module.exports.get_hourly_weather =  async () => {
	const raw_weather_data = await fetch(`${owm_prefix}lat=${location.lat}&lon=${location.lon}&appid=${jablko_config.weather.key}&exclude=minutely,current,daily`);
	return await raw_weather_data.json();
}

module.exports.get_daily_weather = async () => {
	const raw_weather_data = await fetch(`${owm_prefix}lat=${location.lat}&lon=${location.lon}&appid=${jablko_config.weather.key}&exclude=minutely,current,hourly`);
	return await raw_weather_data.json();
}

module.exports.get_all_weather = async () => {
	const raw_weather_data = await fetch(`${owm_prefix}lat=${location.lat}&lon=${location.lon}&appid=${jablko_config.weather.key}&exclude=minutely`);
	return await raw_weather_data.json();
}

