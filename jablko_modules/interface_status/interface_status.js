// interface_status.js: Jablko Modules - Interface Status
// Cale Overstreet
// August 20, 2020
// Returns server process status and performance information to client
// Exports: permission_level, generate_card(), check_status(req, res)

const fs = require("fs").promises;

const jablko = require("../../jablko_interface.js");
const timing = require("../../src/timing.js")

module.exports.permission_level = 0

module.exports.generate_card = async () => {
	return await fs.readFile("./jablko_modules/interface_status/interface_status.html");
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


