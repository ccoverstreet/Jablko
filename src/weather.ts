// weather.ts: Open Weather Map Wrapper
// Cale Overstreet
// July 31, 2020
// This wrapper is used by Jablko components to interface with the OWM API through a unified module.
// Exports: get_current_weather()

const jablko_config = (await import("../jablko_interface.ts")).jablko_config;

console.log("Getting lattitude and longitude for weather data...");
const raw_location_fetch = await fetch("http://ip-api.com/json");
const location = await raw_location_fetch.json();

const owm_prefix = "https://api.openweathermap.org/data/2.5/onecall?";


export async function get_current_weather() {
	const raw_weather_data = await fetch(`${owm_prefix}lat=${location.lat}&lon=${location.lon}&appid=${jablko_config.weather.key}&exclude=minutely,daily,hourly`);
	return await raw_weather_data.json();
}

export async function get_hourly_weather() {
	const raw_weather_data = await fetch(`${owm_prefix}lat=${location.lat}&lon=${location.lon}&appid=${jablko_config.weather.key}&exclude=minutely,current,daily`);
	return await raw_weather_data.json();
}

export async function get_daily_weather() {
	const raw_weather_data = await fetch(`${owm_prefix}lat=${location.lat}&lon=${location.lon}&appid=${jablko_config.weather.key}&exclude=minutely,current,hourly`);
	return await raw_weather_data.json();
}

export async function get_all_weather() {
	const raw_weather_data = await fetch(`${owm_prefix}lat=${location.lat}&lon=${location.lon}&appid=${jablko_config.weather.key}&exclude=minutely`);
	return await raw_weather_data.json();
}

