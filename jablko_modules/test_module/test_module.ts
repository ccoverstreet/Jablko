// Jablko Modules: Test Module
// Cale Overstreet
// May 23, 2020
// This module is just a test module for testing formatting and new features

import { Context } from "https://deno.land/x/oak/mod.ts";

export function generate_card() {
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
<div class="jablko_module_card">
	<div class="card_title">Test Module</div>
	<hr>
	<div class="label_value_pair">
		<div class="label">Hello</div>
		<div class="value">World</div>
	</div>
	<br>
	<br>
	<br>
	<br>
	<br>
	<br>
	<br>
	<br>
	<br>
	<br>
	<br>
	<br>
	<button onclick="test_module.fire()">Push to Talk</button>
</div>
`
}

export function test_response(context: Context) {
	context.response.type = "json";
	context.response.body = {
		hello: "world"
	};
}
