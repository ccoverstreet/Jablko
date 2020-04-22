// jablko.js: Main entrypoint for the Jablko Smart Home System
// Cale Overstreet, Corinne Gerhold
// April 17, 2020
// Jablko.js serves as a manager for all Jablko related processes and is responsible for monitoring I/O output for all subprocesses, restarting processes, identifying problems internally.

const jablko_root = __dirname;
const child_options = {
	silent: true
};

const fork = require("child_process").fork; // Used for forking submodules

console_colors = {
	"normal": "\x1b[0m",
	"reset": "\x1b[0m",
	"red": "\x1b[31m"
};

function jablko_log(process_name, data, options={color: "normal"}) {
	var prefix = `${process_name} [${new Date().toLocaleString("sv-SE")}]: ` // Creating prefix for process log
	const split_data_string = data.split("\n"); // Split at newline to create proper indentation later on

	var output_string = `${console_colors[options.color]}${prefix}`;

	// Append each line with proper indentation
	for (var i = 0; i < split_data_string.length; i++) {
		if (split_data_string[i].length != 0) {
			if (i == 0) {
				output_string += split_data_string[i];
			} else {
				output_string += "\n" + " ".repeat(prefix.length) + split_data_string[i];
			}	
		}
	}

	output_string += `${console_colors["reset"]}`;

	console.log(output_string);
}

function jablko_fork(fork_name, program_location, args) {
	// Fork name is somehting like "Jablko Web Interface"
	// Program location is relative to jablko_root and should be formatted like /path/to/file
	var new_fork = fork(`${jablko_root}${program_location}`, args, {silent: true});

	new_fork.stdout.on("data", function(data) {
		jablko_log(fork_name, data.toString());
	});

	new_fork.stderr.on("data", function(data) {
		jablko_log(fork_name, data.toString(), {color: "red"});
	});

	new_fork.on("exit", function(code) {
		if (code == 0) {
			jablko_log(fork_name, "Exited with no error code");
			return;
		} else {
			jablko_log("Jablko", `Process ${fork_name} exited with code ${code}`);
			jablko_fork(fork_name, program_location);
		}
	});

	return new_fork;
}

var jablko_server= jablko_fork("Jablko Server", "/jablko_server.js");

var jablko_sms_server = jablko_fork("Jablko SMS Server", "/sms_server/jablko_sms_server.js", process.argv.slice(2, 4));
