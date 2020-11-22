// user_authentication.js: Jablko User Authentication Middleware
// Cale Overstreet
// August 19, 2020
// Reads from SQLite database and checks if request has the proper authentication.
// Exports: user_authentication_middleware

const fs = require("fs").promises;
const bcrypt = require("bcrypt");
const fetch = require("node-fetch");

const jablko = require("../jablko_interface.js");

var login_history = {};
try {
	login_history = require("../log/login_history.json");
} catch (error) {
	console.log("./log/login_history.json not found");
	console.debug(error);
}

function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

module.exports.user_authentication_middleware = async function(req, res, next) {
	console.debug(req.originalUrl);
	if (req.originalUrl == "/login") {
		try {
			await sleep(1000);

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

				// Add to login history
				if (req.connection.remoteAddress in login_history) {
					// Connection is already in history
					if (req.body.username in login_history[req.connection.remoteAddress].access) {
						login_history[req.connection.remoteAddress].access[req.body.username].count++
					} else {
						login_history[req.connection.remoteAddress].access[req.body.username].count = 1;
					}

					login_history[req.connection.remoteAddress].access[req.body.username].last_access = new Date().toLocaleString("sv-SE");
				} else {
					login_history[req.connection.remoteAddress] = {
						access: {}
					}

					login_history[req.connection.remoteAddress].access[req.body.username] = {
						count: 1,
						last_access: new Date().toLocaleString("sv-SE")

					}

					const remote_address = req.connection.remoteAddress;

					await fetch(`https://tools.keycdn.com/geo.json?host=${remote_address}`)
						.then(async (data) => {
							const response = await data.json();
							jablko.messaging_system.send_message(`New login from ip "${remote_address}"
User "${req.body.username}"
At ${response.data.geo.region_name}, ${response.data.geo.country_name}
Provider: "${response.data.geo.isp}"
Request happened at ${new Date().toLocaleString("sv-SE")}`);
						})
						.catch((error) => {
							console.log("Unable to reverse-lookup new-ip");
							console.debug(error);
						})

				}

				write_login_history();
				return
			} else {
				throw new Error("Invalid credentials");
			}
		} catch (err) {
			// Log invalid login and store in file.
			console.log(`Invalid login user "${req.body.username}"`);
			console.debug(err);
			res.json({status: "fail", message: "Invalid Login"});

			fs.appendFile("./log/failed_logins.log", `${req.connection.remoteAddress},${req.body.username},${new Date().toLocaleString("sv-SE")}\n`);

			return;
		}
	} else if (req.originalUrl == "/bot_callback") {
		await next();
	} else if (req.originalUrl.startsWith("/module_callback")) {
		// Callback for Jablko Modules on local network. Checks if IPv4 matches pattern "10.0.0.*". Need to add configuration later
		// There might be a better option to this.

		const split_ip = req.ip.split(":");
		const IPv4 = split_ip[split_ip.length - 1];
		console.log(IPv4);

		if (!IPv4.startsWith("10.0.0.") && IPv4 != "1") {
			res.send("Not a valid IP");
			return;
		}

		req.url = req.url.replace("/module_callback", ""); // Remove the module_callback part of the request
		await next()
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
			// Add user data to req object and pass to route handlers
			const user_data = await jablko.user_db.get("SELECT * FROM users WHERE username=(?)", [session_id.username]);

			req.user_data = {
				username: session_id.username,
				first_name: user_data.first_name,
				wakeup_time: user_data.wakeup_time,
				permission_level: user_data.permission_level
			}

			await next();
		}
	}
}

function write_login_history() {
	fs.writeFile("./log/login_history.json", JSON.stringify(login_history, null, 4))
		.catch((error) => {
			console.log("Unable to write to ip_addresses.log");
			console.debug(error);
		});
}
