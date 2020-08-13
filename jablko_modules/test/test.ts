export function permission_level() {
	return 0;
}

export async function generate_card() {
	return `
	<script>
	const test = {
		toggle_light: async function() {
			const raw_response = await fetch("/jablko_modules/test/toggle_light", {method: "POST"})
				.catch((error) => {
					console.log(error);
					return;	
				});
			const response = await raw_response.json();
			console.log(response);
		}
	}
	</script>
	<div class="jablko_module_card">
	<h1>Test</h1>
	<hr>
	<button onclick="test.toggle_light()">Toggle Light</button>
	<h2>Toggle Light</h2>
	</div>
	`
}

export async function toggle_light(context: any) {
	const raw_response: any = await fetch("http://10.0.0.87/toggle_light", {method: "POST"})
		.catch((error) => {
			console.log(error);
			context.response.type = "json";
			context.response.body = {status: "fail", message: "Unable to toggle"};
			return;
		});

	if (raw_response != undefined) {
		const response = await raw_response.json();	

		if (response.status == "good") {
			context.response.type = "json";
			context.response.body = {status: "good", message: "Toggled"};
		} else {
			context.response.type = "json";
			context.response.body = {status: "fail", message: "Unable to toggle"};
		}
	}
}
