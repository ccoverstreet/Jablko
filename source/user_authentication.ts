import { Context } from "https://deno.land/x/oak/mod.ts";

/***
 *	Checks if request is authenticated and handles accordingly 
 *	@param context: Oak Context
 *	@param next: Handle for next function in Oak middleware
 */
export async function check_authentication(context: any, next: any) {
	if (context.cookies.get("key_1") == null) {
		context.cookies.set("key_1", "asd");
		context.response.body = "fart";
	} else {
		console.log(context.cookies.get("key_1"));
		await next();
	}
}
