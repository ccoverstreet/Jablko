// Jablko Modules: Server Status Module
// Cale Overstreet
// May 3, 2020
// This module serves as an example for creating custom modules and provides information on the smart home dashboard related to the run condition, uptime, temperature, and memory usage

module.exports = {
	generate_card: generate_card
}

function generate_card() {
	return `
<script> 
</script>
<div class="jablko_module_card light_card">
	<div class="card_title">Server Status</div>
	<hr>
	<div class="label_value_pair">
		<div class="label">SMS Server Status:</div>
		<div class="value">N/A</div>
	</div>
	<div class="label_value_pair">
		<div class="label">SMS Server Message:</div>
		<div class="value">N/A</div>
	</div>
</div>
	`
}
