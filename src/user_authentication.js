// user_authentication.js: Jablko User Authentication Middleware
// Cale Overstreet
// August 19, 2020
// Reads from SQLite database and checks if request has the proper authentication.
// Exports: user_authentication_middleware

const fs = require("fs").promises;
const bcrypt = require("bcrypt");

const jablko = require("../jablko_interface.js");

module.exports.user_authentication_middleware = async function(req, res, next) {
	if (req.originalUrl == "/login") {
		try {
			// Get hash from SQLite
			const password_hash = (await jablko.user_db.get("SELECT password FROM users WHERE username=?", [req.body.username])).password;

			// Compare hashes and handle correctly
			if (await bcrypt.compare(req.body.password.toString(), password_hash)) {
				// Password hash matches, generate random string and put in sqlite database
				const random_string = await bcrypt.hash(Math.random().toString(36), 10);
				res.cookie("key_1", random_string);
				
				jablko.user_db.run("INSERT INTO login_sessions (username, session_cookie, creation_time) VALUES (?, ?, ?)", [req.body.username, random_string, Date.now()])
					.catch((error) => {
						console.log(error);
					});
				
				console.log(`User "${req.body.username}" has logged in`);
				res.json({status: "good", message: "Logged In"});
				return
			} else {
				throw new Error("Invalid credentials");
			}
		} catch (err) {
			console.log(err);
			res.json({status: "fail", message: "Invalid Login"});
			return;
		}
		
	} else if (req.originalUrl == "/bot_callback") {
		await next();
	} else if (req.cookies.key_1 == null) {
		res.sendFile(`${jablko.html_root}/login/login.html`);
		return;
	} else {
		const session_id = await jablko.user_db.get("SELECT * from login_sessions WHERE session_cookie=?", [req.cookies.key_1]);

		if (session_id == undefined || Date.now() - session_id.creation_time > jablko.jablko_config.database.session_lifetime) {
			// Invalid session id or cookie is expired, send login page
			res.sendFile(`${jablko.html_root}/login/login.html`);
			return;
		} 

		if (req.originalUrl == "/logout") {
			console.log(`User "${session_id.username}" has logged out`);
			jablko.user_db.run("DELETE FROM login_sessions WHERE session_cookie=?", [req.cookies.key_1]);
			res.json({status: "good", message: "Logged out"});
		}  else {
			req.username = session_id.username;
			await next();
		}
	}
}
