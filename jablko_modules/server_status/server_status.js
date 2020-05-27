// Jablko Modules: Server Status Module
// Cale Overstreet
// May 3, 2020
// This module serves as an example for creating custom modules and provides information on the smart home dashboard related to the run condition, uptime, temperature, and memory usage

module.exports = {
	generate_card: generate_card,
	check_status: check_status
}

function generate_card() {
	return `
<script> 
	server_status = {
		check_status: function() {
			$.post("/jablko_modules/server_status/check_status", function(data) {
				console.log(data);
				const interface_status_element = document.getElementById("interface_status")
				const interface_uptime_element = document.getElementById("interface_uptime")

				interface_status_element.textContent = data.interface_status;
				interface_uptime_element.textContent = data.interface_uptime;

				switch (data.interface_status) {
					case "good":
						interface_status_element.style.color = "#44bd32";
						interface_uptime_element.style.color = "#f5f6fa";
						break;
					case "fail":
						interface_status_element.style.color = "#c23616";
						break;
					default:
						interface_status_element.style.color = "#e1b12c";
				}
			})
				.fail(function() {
					const interface_status_element = document.getElementById("interface_status")
					const interface_uptime_element = document.getElementById("interface_uptime")

					interface_status_element.textContent = "fail"
					interface_status_element.style.color = "#c23616"

					interface_uptime_element.textContent = "Server Down"
					interface_uptime_element.style.color = "#c23616"
				});
		}
	}	
	
	server_status.check_status();
	setInterval(server_status.check_status, 30000);
</script>
<div class="jablko_module_card">
	<div class="card_title">Server Status</div>
	<hr>
	<div class="label_value_pair">
		<div class="label">Interface Status:</div>
		<div id="interface_status" class="value">N/A</div>
	</div>
	<div class="label_value_pair">
		<div class="label">Interface Uptime:</div>
		<div id="interface_uptime" class="value">N/A</div>
	</div>
	<br>
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

function check_status(req, res) {
	res.json({interface_status: "good", interface_uptime: process.uptime()});
}
