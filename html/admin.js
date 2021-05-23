class JMODEntry extends HTMLElement {
	constructor() {
		super();

		this.start = this.start.bind(this);
		this.stop = this.stop.bind(this);
		this.toggleEditor = this.toggleEditor.bind(this);

		this.attachShadow({mode: "open"});

		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="assets/standard.css"/>
<style>
.entry {
	display: flex;
	flex-wrap: wrap;
	justify-content:center;
	background-color: var(--clr-surface-2);
	border-radius: 5px;
	padding: 0.3em;
}
.entry > h3 {
	width: 100%;
}

#config-editor-panel {
	display: flex;
	width: 100%;
}
#config-editor-panel > textarea {
	flex: 1;
	height: 7em;
	background-color: var(--clr-surface-1);
	color: var(--clr-font-high);
}
</style>
<div class="entry">
	<h3 id="jmod-name" style="word-break: break-all;"></h3>
	<button onclick="this.getRootNode().host.start()">Start</button>
	<button onclick="this.getRootNode().host.stop()" style="border-color: var(--clr-red)">Stop</button>
	<button onclick="this.getRootNode().host.toggleEditor()" style="border-color: var(--clr-font-high)">Edit Config</button>

	<div id="config-editor-panel" style="display:none;">
		<textarea id="config-editor"></textarea>
	</div>
</div>
		`
	}

	init(name, config) {
		this.config = config;	
		this.jmodName = name;

		this.shadowRoot.getElementById("jmod-name").textContent = this.jmodName;
		this.shadowRoot.getElementById("config-editor").value = JSON.stringify(this.config, null, "  ");
	}

	start() {
		fetch("/admin/startJMOD", {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({"jmodName": this.jmodName})
		})
			.then(async data => {
				console.log(await data.text());
			})
			.catch(err => {
				console.error(err);
			});
	}

	stop() {
		fetch("/admin/stopJMOD", {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({jmodName: this.jmodName})
		})
			.then(async data => {
				console.log(await data.text());
			})
			.catch(err => {
				console.error(err);
			});
	}

	toggleEditor() {
		const editorPanel = this.shadowRoot.getElementById("config-editor-panel")
		if (editorPanel.style.display === "none") {
			editorPanel.style.display = "flex";
		} else {
			editorPanel.style.display = "none";
		}
	}
}

customElements.define("jmod-entry", JMODEntry);


function getJMODData() {
	fetch("/admin/getJMODData", {

	})	
		.then(async data => {
			const mods = await data.json();

			const holder = document.getElementById("jmod-entry-holder");
			holder.innerHTML = "";

			for (mod in mods) {
				console.log(mod);
				const newEntry = document.createElement("jmod-entry");
				newEntry.init(mod, mods[mod]);
				holder.appendChild(newEntry);
			}
		})
		.catch(err => {
			console.error(err);
			console.log(err);
		})
}

document.addEventListener("DOMContentLoaded", function() {
	getJMODData();
});

function getUserList() {
	fetch("/admin/getUserList", {
		method: "POST"
	})
		.then(async data => {
			const res = await data.json();
			const holder = document.getElementById("user-list");

			holder.innerHTML = "";

			for (user of res) {
				const userDisplay = document.createElement("p");
				userDisplay.textContent = "- " + user;

				holder.appendChild(userDisplay);
			}
		})
		.catch(err => {
			console.error(err);
		});
}

function createUser(event) {
	if (event.key != "Enter") {
		return;
	}

	const username = document.getElementById("create-user-username").value;
	const password1 = document.getElementById("create-user-password1").value;
	const password2 = document.getElementById("create-user-password2").value;

	if (password1 !== password2) {
		console.error(new Error("Passwords do not match"));
		alert("Passwords do not match");
		return
	}

	event.preventDefault();

	fetch("/admin/createUser", {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({username, password: password1})
	})
		.then(async data => {
			console.log(await data.text());
			getUserList()
		})
		.catch(err => {
			console.log(err);
			console.error(err);
		})
}

function deleteUser(event) {
	if (event.key != "Enter") {
		return;
	}

	event.preventDefault();

	const username = document.getElementById("delete-user-username").value;

	fetch("/admin/deleteUser", {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({username})
	})
		.then(async data => {
			console.log(await data.text());
			getUserList()
		})
		.catch(err => {
			console.log(err);
			console.error(err);
		})
}

document.addEventListener("DOMContentLoaded", function() {
	getUserList();
})
