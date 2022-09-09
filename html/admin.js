class JablkoModCard extends HTMLElement {
	constructor() {
		super();

		this.attachShadow({mode: "open"});
	}

	init = (name, config) => {
		this.data = {};
		this.data.name = name;
		this.data.config = config;

		this.shadowRoot.innerHTML = `
		<link rel="stylesheet" href="/assets/standard.css"/>
		<div class="card" style="background-color: var(--clr-surface-2);
		display: flex; height: 3em;">
			<h3 style="line-height: 2em">${this.data.name}</h3>
			<div style="flex-grow: 1"></div>
			<button onclick="this.getRootNode().host.remove()">Delete</button>
		</div>
		`
	}

	remove = () => {
		const payload = {
			name: this.data.name
		};

		fetch("/admin/removeMod", {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify(payload)
		})
			.then(async res => {
				res = await jablko.checkFetchError(res);
				jablko.getModList(document.getElementById("mod-edit-content"));
			})
			.catch(res => {
				alert(res.data.err);
				console.log(res.data);
			});
	}
}

class JablkoPrompt extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({mode: "open"});
	}

	init = (conf, callback) => {
		this.config = conf;

		this.shadowRoot.innerHTML = `
		<link rel="stylesheet" href="/assets/standard.css"/>
		<style>
			#container {
				position: absolute;
				top: 0;
				left: 0;
				background-color: rgba(0, 0, 0, 0.2);
				z-index:100;
				width: 100vw;
				height: 100vh;
			}	

			#display {
				height: max(100%, 100vh);
			}

			.pair {
				 display: flex;
				 align-items: center;
			}
			.pair > p {
				margin-right: 1em;
			}

			input {
				height: 1.5em;
				margin: 0.25em;
			}
		</style>
		<div id="container">
			<div id="display">
			</div>
		</div>
		`

		const disp = this.shadowRoot.querySelector("#display");

		const promptBody = document.createElement("div");
		promptBody.style.display = "flex";
		promptBody.style.flexDirection = "column";
		promptBody.style.backgroundColor = "var(--clr-surface-1)";
		promptBody.style.margin = "0.5em 0.5em"
		promptBody.style.padding = "0.5em"


		for (const field of conf.fields) {
			switch (field.type) {
				case "number":  {
					let pair = document.createElement("div");
					pair.classList.add("pair");

					let input = document.createElement("input");
					if (!field.id)
						throw "Jablko Prompt: No id specified";
					input.id = field.id;

					let label = document.createElement("p");

					label.textContent = field.label + ":";

					pair.appendChild(label);
					pair.appendChild(input);
					promptBody.appendChild(pair)
					break;
				}

				case "string": {
					let pair = document.createElement("div");
					pair.classList.add("pair");

					let input = document.createElement("input");
					if (!field.id)
						throw "Jablko Prompt: No id specified";
					input.id = field.id;

					let label = document.createElement("p");

					label.textContent = field.label + ":";

					pair.appendChild(label);
					pair.appendChild(input);
					promptBody.appendChild(pair)
					break;
				}

				default: 
					console.log(2);
			}
		}

		const cancel = document.createElement("button");
		cancel.textContent = "Cancel";
		cancel.onclick = this.close;

		const confirm = document.createElement("button");
		confirm.textContent = "Confirm";
		confirm.onclick = function() {
			callback(this.getRootNode().host);
		}

		promptBody.appendChild(cancel);
		promptBody.appendChild(confirm);

		disp.appendChild(promptBody);
	}

	close = () => {
		console.log("Closing prompt");
		this.remove();
	}

	// Returns JSON object containing prompt fields
	collect = () => {
		const out = {}
		for (const field of this.config.fields)	{
			console.log(field);
			switch (field.type) {
				case "number": {
					const rawVal = this.shadowRoot.querySelector("#"+field.id).value;
					console.log(rawVal);
					out[field.id] = parseInt(rawVal);
					break;
				} 

				case "string": {
					out[field.id] = this.shadowRoot.querySelector("#"+field.id).value;
					break;
				}
			}
		}

		return out;
	}
}

customElements.define("jablko-mod-card", JablkoModCard);
customElements.define("jablko-prompt", JablkoPrompt);

const jablko = {
	showTab: (tabId) => {
		const pages = document.querySelectorAll("#content > div");
		for (const p of pages) {
			console.log(p.id, p.id === tabId);
			if (p.id === tabId) {
				p.style.display = "block";
			} else {
				p.style.display = "none";
			}
		}
	},

	prompt: (conf, callback) => {
		const p = document.createElement("jablko-prompt");
		p.init(conf, callback);
		document.querySelector("body").appendChild(p);
	},

	checkFetchError: async (res) => {
		res.data = await res.json();
		if (res.status >= 400 || res.status < 200) 
			throw res;

		return res
	},

	getModList: (modContentElem) => {
		fetch("/admin/getModList")
			.then(async res => {
				res = await jablko.checkFetchError(res);

				modContentElem.innerHTML = "";
				for (const key in res.data) {
					console.log(key);
					const newCard = document.createElement("jablko-mod-card");
					newCard.init(key, res.data[key]);
					modContentElem.appendChild(newCard);
					/*
					const hr = document.createElement("hr");
					hr.style.borderColor = "var(--clr-font-low)";
					modContentElem.appendChild(hr);
					*/
				}
			})
			.catch(async err => {
				console.error(err.url + " " + err.status, err.data, err);
			})
	},

	addMod: () => {
		jablko.prompt({
			fields: [
				{
					type: "string",
					label: "Type",
					id: "type"
				},
				{
					type: "string",
					label: "Name",
					id: "name"
				},
				{
					type: "string",
					label: "Tag",
					id: "tag"
				},
				{
					type: "number",
					label: "Port",
					id: "port"
				}
			]
		},
			function(prompt) {
				console.log(prompt);
				const data = prompt.collect();

				fetch("/admin/addMod", {
					method: "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify(data)
				})
					.then(async res => {
						res = await jablko.checkFetchError(res);
					})
					.catch(res => {
						console.error(res);
					});
			})
	}
}

document.addEventListener("DOMContentLoaded", () => {
	jablko.getModList(document.querySelector("#mod-edit-content"));
})
