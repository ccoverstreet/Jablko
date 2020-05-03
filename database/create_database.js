const sqlite3 = require("sqlite3");

const db = new sqlite3.Database("test.db", function(err) {
	if (err) {
		console.log(err.message);
	} else {
		console.log("Connected to database");
	}
});

db.serialize(function() {
	db.run(`CREATE TABLE IF NOT EXISTS test(
	primary_key INTEGER PRIMARY KEY,
	date TEXT)`);
	
	db.run(`INSERT INTO test(date) VALUES("2020-04-03 00:00:00"), ("ASDASD")`);
});


db.close();
