const worker = new Worker(new URL("jablko_interface.ts", import.meta.url).href, {type: "module", deno: true});
