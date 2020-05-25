const sqlite3 = require("sqlite3");
const rl = require("readline-sync");
const bcrypt = require("bcrypt");
const crypto = require("crypto");
const fs = require("fs");

var mode_selection = rl.question(`
What do you want to do?
 1. Create New Jablko Database
 2. Add User to Jablko Database
 3. Modify user info inside Jablko Database
$ `);
switch (mode_selection) {
	case "1":
		console.log("Creation Mode Selected");
		create_jablko_database();
		break;
	case "2":
		console.log("User Creation Mode Selected");
		add_user();
		break;
	default:
		console.log("Invalid Response. Response should be an integer value");
}

function create_jablko_database() {
	const database_filename = rl.question("Enter a filename for the database: ").replace(".db", "");

	if (fs.existsSync(`./${database_filename}.db`)) {
		const overwrite_confirmation = rl.question("The file already exists. Do you wish to overwrite? (y/n)");
		if (overwrite_confirmation != "y") {
			console.log("Cancelling database creation... Returning to mode selection.");
			return;
		} 
	}

	return db = new sqlite3.Database(`./${database_filename}.db`, function(err) {
		if (err) {
			console.log(err.message);
			console.log("Unable to create database. See above error message");
			return;
		} 
		console.log("Successfully created database file...\nCreating table for users");

		db.serialize(function() {
			db.run(`CREATE TABLE IF NOT EXISTS users(
							primary_key INTEGER PRIMARY KEY,
							username TEXT NOT NULL,
							password TEXT NOT NULL,
							salt TEXT NOT NULL,
							first_name TEXT NOT NULL,
							phone_number TEXT NOT NULL,
							phone_carrier TEXT NOT NULL,
							wakeup_time TEXT NOT NULL
						)`);
		});

		const update_jablko_config = rl.question("Would you like to update Jablko's config file? (y/n) ")

		if (update_jablko_config == "y") {
			var current_config = fs.readFileSync(`../jablko_config.json`).toString().split("\n");

			// Search for row containing "database_name" and replace it with a new line that contains the new database line
			var new_config = "";

			for (var i = 0; i < current_config.length; i++) {
				console.log(current_config[i]);
				if (current_config[i].includes("database_name")) {
					current_config[i] = `\t"database_name": "${database_filename}.db"`
				}

				new_config += current_config[i] + "\n";
			}	

			// Write out to file
			console.log("Updating Jablko Config");
			fs.writeFileSync("../jablko_config.json", new_config);
		}

		db.close();
	});
}

function add_user() {
	const config_settings = require("../jablko_config.json");
	
	const db = new sqlite3.Database(config_settings.database_name, function(err) {
		if (err) {
			console.log(err);
			console.log("Error occured when connecting to the database. Check above error.")
			return;
		}

		const username = rl.question("Enter username: ");
		const salt = crypto.randomBytes(16).toString("hex");
		const password = bcrypt.hashSync(rl.question("Enter password: ") + salt, 10000)
		const first_name = rl.question("Enter First Name: ");
		const phone_number = rl.question("Enter Phone Number: ");
		const phone_carrier = rl.question("Enter Phone Carrier: ");
		const wakeup_time = rl.question("Entered Wakeup Time (hh:mm): ");

		db.run(`INSERT INTO users(username, salt, password, first_name, phone_number, wakeup_time, phone_carrier) VALUES (?, ?, ?, ?, ?, ?, ?)`, username, salt, password, first_name, phone_number, wakeup_time, phone_carrier);
		db.close();
	});
}
