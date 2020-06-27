import { Context, helpers } from "https://deno.land/x/oak/mod.ts";
import { DB } from "https://deno.land/x/sqlite/mod.ts" ;
import * as bcrypt from "https://deno.land/x/bcrypt/mod.ts";
import Random from "https://deno.land/x/random/Random.js";

/***
 *	@description Checks if request is authenticated and handles accordingly 
 *	@parameter context: Oak Context
 *	@parameter next: Handle for next function in Oak middleware
 */
export async function check_authentication(context: any, next: any) {
	// Create SQLite database connection
	const db = new DB("database/primary.db");

	if (context.request.url.pathname == "/login") {
		const login_data = (await context.request.body()).value;

		// Query database for user data to compare hash and add info to context
		const user_data = [...db.query("SELECT * FROM users WHERE username=(?)", [login_data.username])];

		// Check if any results were found
		if (user_data[0] == undefined) {
			context.response.type = "json";
			context.response.body = {status: "fail", message: "Invalid Login"};
			return	
		}

		// Use bcrypt to compare submitted password to database hash
		if (await bcrypt.compare(login_data.password, user_data[0][2]) == true) {

			// Create cookie string and add to login_sessions table
			const cookie_string = new Random().string(64);
			db.query("INSERT INTO login_sessions (session_cookie, username, creation_time) VALUES (?, ?, ?)", [cookie_string, login_data.username, new Date().getTime()]);

			context.cookies.set("key_1", cookie_string)			

			context.response.type = "json";
			context.response.body = {status: "good", message: `Welcome ${user_data[0][4]}`};
		} else {
			context.response.type = "json";
			context.response.body = {status: "fail", message: "Invalid Login"};
		}

		return;
	} else if (context.cookies.get("key_1") == null) {
		// Client has no corresponding cookies. Prevents from erroring out
		const decoder = new TextDecoder("utf-8");
		const data = decoder.decode(await Deno.readFile("./public_html/login/login.html"));

		context.response.type = "html";
		context.response.body = data;
	} else {
		// Check if user wishes to logout
		if (context.request.url.pathname == "/logout") {
			db.query("DELETE FROM login_sessions WHERE session_cookie=(?)", [context.cookies.get("key_1")]);
			context.response.type = "json";
			context.response.body = {status: "good", message: "You have logged out"};
			return;
		}

		// Query login session database to see if user session exists
		const session_data = [...db.query("SELECT session_cookie, username, creation_time FROM login_sessions WHERE session_cookie=(?)", [context.cookies.get("key_1")])];


		// Check if session was found or exists
		if (session_data.length === 0) {
			// User is not authenticated
			const decoder = new TextDecoder("utf-8");
			const data = decoder.decode(await Deno.readFile("./public_html/login/login.html"));
			context.response.type = "html";
			context.response.body = data;
			return;
		} else {
			// Cookie is in login_sessions table
			const min_time = new Date().getTime() - 259200000;
			if (min_time > session_data[0][2]) {
				// User is no longer authenticated, cookie expired
				db.query("DELETE FROM login_sessions WHERE creation_time<(?)", [min_time]);
				const decoder = new TextDecoder("utf-8");
				const data = decoder.decode(await Deno.readFile("./public_html/login/login.html"));

				context.response.type = "html";
				context.response.body = data;
			} else {
				// User is authenticated, add to context and pass to next()
				context.json_data = (await context.request.body()).value;
				context.session_data = session_data;
				context.user_data = await get_user_data(session_data[0][1]);
				await next();
			}
		}
	}

	db.close();
}

async function get_user_data(username: any) {
	const db = new DB("database/primary.db");
	const raw_user_data = [...db.query("SELECT username, first_name, phone_number, phone_carrier, wakeup_time, permissions FROM users WHERE username=(?)", [username])];

	return {
		username: raw_user_data[0][0],
		first_name: raw_user_data[0][1],
		phone_number: raw_user_data[0][2],
		phone_carrier: raw_user_data[0][3],
		wakeup_time: raw_user_data[0][4],
		permissions: raw_user_data[0][5]
	};

	db.close();
}
