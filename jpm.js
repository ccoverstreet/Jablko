// jpm.js: Jablko Package Manager
// Cale Overstreet
// August 24, 2020
// Used to install Jablko Modules by first downloading module sources to the module_library directory and then copying contents to the correct directory in jablko_modules

const reader = require("readline-sync");
const extract = require("extract-zip");
const fetch = require("node-fetch");
const fs = require("fs");
const { execSync } = require("child_process");

const jablko_config = require("./jablko_config.json");

async function main() {
	console.log("Jablko Package Manager");
	if (process.argv[2] == "init") {
		await init(process.argv.slice(3));	
	} else if (process.argv[2] == "install") {
		await install(process.argv.slice(3));
	}
}

async function init(arguments) {
	console.log(arguments);
	execSync("mkdir -p jablko_modules");

	if (arguments.length == 0) {
		console.log("Installing all modules specified in jablko_config.json\n");

		const module_keys = Object.keys(jablko_config.jablko_modules);
		for (module in jablko_config.jablko_modules) {
			await install_module(jablko_config.jablko_modules[module].repo_archive, module);
			console.log();
		}
	}
}

async function install(arguments) {
	console.log(arguments);
	await install_module(arguments[0], arguments[1]);
}

async function install_module(repository_url, module_target_name) {
	// Will need to be updated to handle more tag/naming conventions
	const split_repo_url = repository_url.split("/");
	const author = split_repo_url[3];
	const repo = split_repo_url[4];
	const tag = split_repo_url[6].split(".zip")[0].replace("v", "");
	const extracted_zip = `${repo}-${tag}`;


	console.log(`Installing ${repo} by ${author} from ${tag}`);

	// Should emulate synchronous behavior
	const data = await fetch(repository_url);

	await execSync("mkdir -p module_library");
	const buffer = await data.buffer();

	// Check if files already exist
	if (fs.existsSync(`./module_library/${module_target_name}.zip`)) {
		const answer = reader.question(`Do you want to replace the module "${module_target_name}"?`).trim();
		if (!/[yn]/.test(answer)) {
			console.log("Invalid response... Aborting install.");;
			process.exit(1);
		} 
		
		if (answer == "y") {
			if (fs.existsSync(`./jablko_modules/${module_target_name}`)) {
				execSync(`rm -r ./jablko_modules/${module_target_name}`);
			}
		} else {
			console.log("Not re-installing module")
			return;
		}
	}
	fs.writeFileSync(`./module_library/${module_target_name}.zip`, buffer);

	await extract(`./module_library/${module_target_name}.zip`, {dir: `${process.cwd()}/module_library`});


	execSync(`mkdir -p ./jablko_modules && mkdir -p ./jablko_modules/${module_target_name} && cp ./module_library/${extracted_zip}/* ./jablko_modules/${module_target_name}`);
}

main();
