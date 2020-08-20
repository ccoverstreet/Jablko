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
		// Get hash from SQLite
		const password_hash = (await jablko.user_db.get("SELECT password FROM users WHERE username=?", [req.body.username])).password;
		console.log(password_hash);
		if (await bcrypt.compare(req.body.password.toString(), password_hash)) {
			// Password hash matches, generate random string and put in sqlite database
			const random_string = await bcrypt.hash(Math.random().toString(36), 10);
			res.cookie("key_1", random_string);
			
			console.log(req.body);
			jablko.user_db.run("INSERT INTO login_sessions (username, session_cookie, creation_time) VALUES (?, ?, ?)", [req.body.username, random_string, Date.now()])
				.catch((error) => {
					console.log(error);
				});
			
			res.json({status: "good", message: "Logged In"});
			return
		} else {
			res.json({status: "fail", message: "Invalid Login"});
			return;
		}
	} else if (req.originalUrl == "/bot_callback") {
		console.log("Bot callback");
	} else if (req.cookies.key_1 == null) {
		console.log("No cookie, not logged in");
	} else {
		const session_id = await jablko.user_db.get("SELECT * from login_sessions WHERE session_cookie=?", [req.cookies.key_1]);
		if (session_id == undefined) {
			// Invalid session id, send login page
			res.sendFile(`${jablko.html_root}/login/login.html`);
			return;
		} 

		if (req.originalUrl == "/logout") {
			jablko.user_db.exec("DELETE FROM login_sessions WHERE session_cookie=?", [req.cookies.key_1]);
			res.sendFile(`${jablko.html_root}/login/login.html`);
		}  else {
			await next();
		}
	}
}
