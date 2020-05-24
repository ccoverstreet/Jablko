const sqlite3 = require("sqlite3");
const rl = require("readline-sync");
const crypto = require("crypto");
const fs = require("fs");

while (true) {
	var mode_selection = rl.question(`
What do you want to do?
 1. Create New Jablko Database
$ `);
	switch (mode_selection) {
		case "1":
			console.log("Creation Mode Selected");
			create_jablko_database();
			break;
		default:
			console.log("Invalid Response. Response should be an integer value");
	}
}

function create_jablko_database() {
	const database_filename = rl.question("Enter a filename for the database: ");

	if (fs.existsSync(`./${database_filename}.db`)) {
		const overwrite_confirmation = rl.question("The file already exists. Do you wish to overwrite? ");
		if (overwrite_confirmation != "y") {
			console.log("Cancelling database creation... Returning to mode selection.");
			return;
		} else {
			const db = new sqlite3.Database("test.db", function(err) {
				if (err) {
					console.log(err.message);
				} else {
					console.log("Successfully created database file...\nCreating table for users");

					db.serialize(function() {
						db.run(`CREATE TABLE IF NOT EXISTS test(
							primary_key INTEGER PRIMARY KEY,
							username TEXT NOT NULL,
							password TEXT NOT NULL,
							first_name TEXT NOT NULL,
							phone_number TEXT NOT NULL,
							wakeup_time TEXT NOT NULL
						)`);
					});
				}
			});

			db.close();
		}
	}
}

