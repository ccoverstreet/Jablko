// jablko.js: Main entrypoint for the Jablko Smart Home System
// Cale Overstreet, Corinne Gerhold
// April 17, 2020
// Jablko.js serves as a manager for all Jablko related processes and is responsible for monitoring I/O output for all subprocesses, restarting processes, identifying problems internally.

const jablko_root = __dirname;
const child_options = {
	silent: true
};

function jablko_log(process_name, data, options={color: "normal"}) {
	var prefix = `${process_name} [${new Date().toLocaleString("sv-SE")}]: ` // Creating prefix for process log
	const split_data_string = data.toString().split("\n"); // Split at newline to create proper indentation later on

	var output_string = prefix;
	for (var i = 0; i < split_data_string.length; i++) {
		if (i == 0) {
			output_string += split_data_string[i];
		} 
		output_string += "\nASD";
		
	}

	console.log(output_string);
		
	console.log(options.color);
}

const fork = require("child_process").fork; // Used for forking submodules

var jablko_web_interface = fork(`${jablko_root}/web_interface/jablko_web_interface.js`, child_options);

jablko_web_interface.stdout.on("data", function(data) {
	jablko_log("Jablko Web Interface", data);
});

jablko_web_interface.stderr.on("data", function(data) {
	console.log(data.toString());
});

