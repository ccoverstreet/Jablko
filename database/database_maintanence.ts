// Jablko: Database Maintanence Tool
// Cale Overstreet
// May 30, 2020
// Used for creating the database used by Jablko. Can be used to add users and modify user data.

import { DB } from "https://deno.land/x/sqlite/mod.ts"; // SQLite3 module
import { readLines } from "https://deno.land/std/io/bufio.ts";
import * as banana from "https://deno.land/x/input/index.ts";

const input = new banana();
console.log(await input.question("ASDASDA: "));
