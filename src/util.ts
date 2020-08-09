export async function readFileStr(filepath: string) {
	const decoder = new TextDecoder("utf-8");
	const text = decoder.decode(await Deno.readFile(filepath))
	return text;
}
