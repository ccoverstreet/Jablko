// interface_status.js: Jablko Modules - Interface Status
// Cale Overstreet
// August 20, 2020
// Returns server process status and performance information to client
// Exports: permission_level, generate_card(), check_status(req, res)

const fs = require("fs").promises;
const path = require("path");

const module_name = path.basename(__dirname);
const jablko = require(module.parent.filename);
const module_config = jablko.jablko_config.jablko_modules[module_name];

const timing = require(`${module.parent.path}/src/timing.js`);

// Check if config is valid
if (module_config.update_interval == undefined || module_config.update_interval == null) {
	throw new Error("Incorrect Configuration");
}

module.exports.permission_level = 0

module.exports.generate_card = async () => {
	var data = (await fs.readFile(`${__dirname}/interface_status.html`, "utf8")).replace(/\$MODULE_NAME/g, module_name).replace(/\$UPDATE_INTERVAL/g, module_config.update_interval);
	return data;
}

module.exports.check_status = async (req, res) => {
	const cpu_temp = parseInt(await fs.readFile("/sys/class/thermal/thermal_zone0/temp", "utf8")) / 1000;
	const meminfo = (await fs.readFile("/proc/meminfo", "utf8")).split("\n");
	const total_mem = (parseInt(meminfo[0].split(/[ ]+/)[1]) / 1000000);
	const meminfo_summary = `${(total_mem - parseFloat(meminfo[2].split(/[ ]+/)[1]) / 1000000).toFixed(2)} GB`
	
	const raw_uptime = (new Date().getTime() - jablko.server_start_time) / 1000;
	const hours = Math.floor(raw_uptime / 3600);
	const minutes = Math.floor((raw_uptime - 3600 * hours) / 60);
	const seconds = Math.floor(raw_uptime - 3600 * hours - 60 * minutes);
	const formatted_uptime = `${hours} h ${minutes} m ${seconds}s`;
	
	res.json({interface_status: "good", interface_uptime: formatted_uptime, interface_response_time: timing.get_handling_time().toFixed(2) + " ms", cpu_temperature: cpu_temp + " C", memory_usage: meminfo_summary});
}

module.exports.chatbot_uptime = async () => {
	const responses = [
		"I've been up for ",
		"I've been watching you for ",
		"This version of me has been alive for "
	];

	const raw_uptime = (new Date().getTime() - jablko.server_start_time) / 1000;
	const hours = Math.floor(raw_uptime / 3600);
	const minutes = Math.floor((raw_uptime - 3600 * hours) / 60);
	const seconds = Math.floor(raw_uptime - 3600 * hours - 60 * minutes);
	const formatted_uptime = `${hours} h ${minutes} m ${seconds}s.`;

	return responses[Math.floor(Math.random() * responses.length)] + formatted_uptime;
}
