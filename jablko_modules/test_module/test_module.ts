// test_module.ts: Test Module
// Cale Overstreet
// June 26, 2020
// This is a testbed for new Jablko Modules features


const smtp_client = (await import("../../jablko_interface.ts")).smtp_client;

export const info = {
	permissions: "All"
};

export function generate_card() {
	return `
	<script>
	const test_module = {
		send_message: async function() {
			const message_text = document.getElementById("test_module_input").value;

			const response = await fetch("/jablko_modules/test_module/send_message", {
				method: "POST",
				headers: {
					"Accept": "application/json",
					"Content-Type": "application/json"
				},
				body: JSON.stringify({message: message_text})
			})
			.catch(error => {
				console.log(error);
			})

		}
	}
	</script>
	<div id="test_module" class="jablko_module_card">
	<div class="card_title">Test Module</div>
	<div style="display: grid; grid-template-columns: 50% 50%">
	<div class="label">Send Message:</div>
	<input id="test_module_input">
	</div>
	<button onclick="test_module.send_message()">Send</button>
	</div>
	`;
}

export async function send_message(context: any) {
	smtp_client.send_message(context.user_data, context.json_data.message);

	context.response.type = "json";
	context.response.body = {
		status: "good"
	};
}
