// Jablko Modules: Test Module
// Cale Overstreet
// May 23, 2020
// This module is just a test module for testing formatting and new features

module.exports = {
	generate_card: generate_card,
	test_response: test_response
}

function generate_card() {
	return `
<script>
	var test_module = {
		fire: function() {
			console.log("Burn");
			$.post("/jablko_modules/test_module/test_response", function(data) {
				console.log(data);
			});
		}
	}	
</script>
<div class="jablko_module_card light_card">
	<div class="card_title">Test Module</div>
	<hr>
	<div class="label_value_pair">
		<div class="label">Hello</div>
		<div class="value">World</div>
	</div>
	<hr>
	<hr>
	<hr>
	<button onclick="test_module.fire()">Push to Talk</button>
</div>
`
}

function test_response(req, res) {
	res.send("This is the server speaking.");
}
