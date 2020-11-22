// database_maintanence.js: Database Maintanence Tool
// Cale Overstreet
// August 24, 2020
// Used for creating the sqlite database Jablko uses for user credentials/info and session storage

const sqlite = require("sqlite-async");
const bcrypt = require("bcrypt");
const reader = require("readline-sync");

async function main() {
	console.log("Jablko Database Maintanence Tool");

	while (true) {
		try {
			const mode = parseInt(reader.question("Select Mode"));
			switch (mode) {
				case 1:
					await create_database();
					break;
				case 2:
					await create_user();
					break;
				case 3:
					break;
				default:
					console.log("Not a valid mode");
					break;
			}
		} catch (error) {
			console.log("Invalid Input");
			console.log(error);
		}
	}
}

async function create_database() {
	console.log("Creationg databsae");

	const database_name_1 = reader.question("Enter Database Name: ").trim().replace(".db", "");
	const database_name_2 = reader.question("Confirm Database Name: ").trim().replace(".db", "");

	if (database_name_1 != database_name_2) {
		console.log("Database names do not match.");
		return;
	}

	const database = await sqlite.open(`database/${database_name_1}.db`)
		.catch((error) => {
			console.log(error);
		});

	database.exec(`CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT, 
					username TEXT NOT NULL,
					password TEXT NOT NULL,
					first_name TEXT NOT NULL,
					wakeup_time TEXT NOT NULL,
					permission_level INTEGER NOT NULL)`);

	database.exec(`CREATE TABLE IF NOT EXISTS login_sessions (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					session_cookie TEXT NOT NULL,
					username TEXT NOT NULL,
					creation_time INTEGER NOT NULL)`);
}

async function create_user() {
	const database_name_1 = reader.question("Enter Database Name: ").trim().replace(".db", "");
	const database_name_2 = reader.question("Confirm Database Name: ").trim().replace(".db", "");

	if (database_name_1 != database_name_2) {
		console.log("Database names do not match.");
		return;
	}

	var user_data = {
		username: "",
		password: "",
		first_name: "",
		wakeup_time: "",
		permission_level: 0
	}

	user_data.username = reader.question("Enter Username: ").trim();

	const password_1 = reader.question("Enter Password: ").trim();
	const password_2 = reader.question("Confirm Password: ").trim();

	if (password_1 != password_2) {
		console.log("Passwords do not match.");
		return;
	}

	user_data.password = password_1;
	user_data.first_name = reader.question("Enter First Name: ").trim();

	const wakeup_time_input = reader.question("Enter Wakeup Time (hh:mm): ");
	const wakeup_regex = /[0-9]{2}:[0-9]{2}/;
	if (!wakeup_regex.test(wakeup_time_input)) {
		console.log("Invalid wakeup time format.");
		return;
	}

	user_data.wakeup_time = wakeup_time_input;

	const permission_level_input = reader.question("Enter permission level for user (0: guest, 1: family, 2: admin): ").trim();
	if (/[0-9]/.test(permission_level_input)) {
		user_data.permission_level = parseInt(permission_level_input, 10);
	}

	// Create password hash
	const password_hash = await bcrypt.hash(user_data.password, 10);
	console.log(user_data);
	
	const database = await sqlite.open(`database/${database_name_1}.db`)
		.catch((error) => {
			console.log(error);
			return;
		});


	await database.run(`INSERT INTO users (username, password, first_name, wakeup_time, permission_level) VALUES (?, ?, ?, ?, ?)`, [user_data.username, password_hash, user_data.first_name, user_data.wakeup_time, user_data.permission_level])
		.catch((error) => {
			console.log(error);
			return;
		});

	database.close();
}

(async () => {
	await main();
})();
