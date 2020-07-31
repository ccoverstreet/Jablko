// weather.ts: Open Weather Map Wrapper
// Cale Overstreet
// July 31, 2020
// This wrapper is used by Jablko components to interface with the OWM API through a unified module.
// Exports: get_current_weather()

const jablko_config = (await import("../jablko_interface.ts")).jablko_config;

export async function get_current_weather() {
	const raw_weather_data = await fetch(`https://api.openweathermap.org/data/2.5/weather?q=${jablko_config.weather.city}&appid=${jablko_config.weather.key}`);
	return await raw_weather_data.json();
}
