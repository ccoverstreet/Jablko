export async function send_message(message: string) {
	const interface_exports = await import("../jablko_interface.ts");
	const smtp_client = interface_exports.smtp_client;
	console.log(smtp_client);
	smtp_client.send_message("coverstreet", message);
}
