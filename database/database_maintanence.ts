// Jablko: Database Maintanence Tool
// Cale Overstreet
// May 30, 2020
// Used for creating the database used by Jablko. Can be used to add users and modify user data.

import { DB } from "https://deno.land/x/sqlite/mod.ts"; // SQLite3 module
import { readLines } from "https://deno.land/std/io/bufio.ts";
import * as bcrypt from "https://deno.land/x/bcrypt/mod.ts";


async function mainloop() {
	console.log("Mode Selection:\n\t1. Create Database\n\t2. Add User");

	for await (const line of readLines(Deno.stdin)) {
		switch (line.trim()) {
			case "1":
				console.log("Database Creation Selected.");
			await create_database();
			break;
			case "2":
				console.log("User Creation Selected");
			await create_user();
			break;
			default:
				console.log("Invalid Input");
		}

		console.log("\nMode Selection:\n\t1. Create Database\n\t2. Add User");
	}
}

async function create_database() {
	console.log("Enter Database Name:");

	for await (const database_name of readLines(Deno.stdin)) {
		console.log(`Is the name "${database_name}" correct? <y/n>:`)
		for await (const line of readLines(Deno.stdin)) { 
			if (line == "y") {
				const db = new DB(`database/${await database_name.replace(".db", "")}.db`);
				db.query(`CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT, 
					username TEXT NOT NULL,
					password TEXT NOT NULL,
					first_name TEXT NOT NULL,
					phone_number TEXT NOT NULL,
					phone_carrier TEXT NOT NULL,
					wakeup_time TEXT NOT NULL,
					permissions TEXT NOT NULL
				)`);

				db.query(`CREATE TABLE IF NOT EXISTS login_sessions (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					session_cookie TEXT NOT NULL,
					username TEXT NOT NULL,
					creation_time INTEGER NOT NULL
				)`);
				return;
			} else {
				break;
			}
		}
		console.log("Enter Database Name:");
	}
}

async function create_user() {
	var database_name = "";
	console.log("Enter Database Name");
	for await (const line of readLines(Deno.stdin)) {
		database_name = line.trim();
		break;
	}

	var	user_data = {
		username: "",
		password_1: "",
		password_2: "",
		first_name: "",
		phone_number: "",
		phone_carrier: "",
		wakeup_time: "",
		permissions: ""
	};

	console.log("Enter Username:");
	for await (const line of readLines(Deno.stdin)) {
		user_data.username = line.trim();
		break;
	}

	console.log("Enter Password:");
	for await (const line of readLines(Deno.stdin)) {
		user_data.password_1 = line.trim();
		break;
	}

	console.log("Confirm Password:");
	for await (const line of readLines(Deno.stdin)) {
		user_data.password_2 = line.trim();
		break;
	}

	// Check if typed passwords match
	if (user_data.password_2 != user_data.password_1) {
		console.log("Passwords do not match. Returning to mode selection");
	}

	console.log("Enter First Name:");
	for await (const line of readLines(Deno.stdin)) {
		user_data.first_name = line.trim();
		break;
	}

	console.log("Enter Phone Number ##########:");
	for await (const line of readLines(Deno.stdin)) {
		const regex = /[0-9]{10}/;
		if (regex.test(line.trim())) {
			user_data.phone_number = line.trim();
			break;
		}

		console.log("Incorrect Format");
	}

	console.log("Enter Phone Carrier (eg. verizon, att, tmobile):");
	for await (const line of readLines(Deno.stdin)) {
		user_data.phone_carrier = line.trim();
		break;
	}

	console.log("Enter Preferred Wakeup Time 24 hour format (hh:mm):");
	for await (const line of readLines(Deno.stdin)) {
		const regex = /[0-9]{2}:[0-9]{2}/;
		if (regex.test(line.trim())) {
			user_data.wakeup_time = line.trim();
			break;
		}

		console.log("Incorrect Format");
	}

	console.log("Enter permissions for user (admin/guest):");
	for await (const line of readLines(Deno.stdin)) {
		user_data.permissions = line.trim();

		break;
		console.log("Incorrect Format");
	}



	// Create Password Hash	
	const hash = bcrypt.hashSync(user_data.password_1);

	// Create Database connection
	const db = new DB(`database/${await database_name.replace(".db", "")}.db`);
	db.query(`INSERT INTO users (username, password, first_name, phone_number, phone_carrier, wakeup_time, permissions) VALUES (
		?, ?, ?, ?, ?, ?, ?
	)`, [user_data.username, hash, user_data.first_name, user_data.phone_number, user_data.phone_carrier, user_data.wakeup_time, user_data.permissions]);
}

// Start Mainloop
mainloop();
