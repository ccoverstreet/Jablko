// Jablko: Database Maintanence Tool
// Cale Overstreet
// May 30, 2020
// Used for creating the database used by Jablko. Can be used to add users and modify user data.

import { DB } from "https://deno.land/x/sqlite/mod.ts"; // SQLite3 module
import { readLines } from "https://deno.land/std/io/bufio.ts";


async function mainloop() {
	console.log("Mode Selection:\n\t1. Create Database\n\t2. Add User");

	for await (const line of readLines(Deno.stdin)) {
		switch (line.trim()) {
			case "1":
				console.log("Database Creation Selected.");
				await create_database();
				break;
			default:
				console.log("Invalid Input");
		}
	}
}

async function create_database() {
	console.log("Enter Database Name:");

	for await (const database_name of readLines(Deno.stdin)) {
		console.log(`Is the name "${database_name}" correct? <y/n>:`)
		for await (const line of readLines(Deno.stdin)) { 
			if (line == "y") {
				const db = new DB(`${await database_name.replace(".db", "")}.db`);
				db.query(`CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT, 
					username TEXT NOT NULL,
					password TEXT NOT NULL,
					salt TEXT NOT NULL,
					first_name TEXT NOT NULL,
					phone_number TEXT NOT NULL,
					phone_carrier TEXT NOT NULL,
					wakeup_time TEXT NOT NULL
				)`);
				return;
			} else {
				break;
			}
	   	}
		console.log("Enter Database Name:");
	}
}

mainloop();
