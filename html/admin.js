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
			<button onclick="this.getRootNode().host.updateMod()">Update</button>
			<button onclick="this.getRootNode().host.removeMod()">Delete</button>
		</div>
		`
	}

	removeMod = () => {
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

	updateMod = () => {
		jablko.prompt({
			fields: [
				{
					type: "string",
					label: "Tag",
					id: "tag"
				}
			]
		},
			async function(prompt) {
				console.log(prompt);
				const data = prompt.collect();
				console.log(data);
				console.log(this);
				console.log(this.data);
				data.name = this.data.name;

				const x = await jablko.postJSON("/admin/updateMod", data)
					.then(async res => {
						res = await jablko.checkFetchError(res);
						jablko.getModList(document.getElementById("mod-edit-content"));
						prompt.close();
					})
					.catch(async res => {
						console.error(res, res.data);
						throw new Error(res.data.err);
					});
			}.bind(this));
	}
}
customElements.define("jablko-mod-card", JablkoModCard);

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
	
	postJSON: (url, data) => {
		return fetch(url, {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify(data)
		})
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
			async function(prompt) {
				console.log(prompt);
				const data = prompt.collect();

				const x = await fetch("/admin/addMod", {
					method: "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify(data)
				})
					.then(async res => {
						res = await jablko.checkFetchError(res);
						jablko.getModList(document.getElementById("mod-edit-content"));
					})
					.catch(async res => {
						console.error(res, res.data);
						throw new Error(res.data.err);
					});
			})


	},

	updateMod: (modName) => {
		
	}
}

document.addEventListener("DOMContentLoaded", () => {
	jablko.getModList(document.querySelector("#mod-edit-content"));
})
