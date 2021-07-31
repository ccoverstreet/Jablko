function InstallJMOD() {
	jmodPath = document.getElementById("install-jmod-input").value.trim();
	console.log(jmodPath);

	fetch("/admin/installJMOD", {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({jmodPath: jmodPath})
	})
		.then(async data => {
			console.log(await data.text());
			getJMODData();
		})
		.catch(err => {
			alert(err);
			console.error(err);
			console.log(err);
		})
}

class JMODEntry extends HTMLElement {
	constructor() {
		super();

		this.attachShadow({mode: "open"});
		this.start = this.start.bind(this);
		this.stop = this.stop.bind(this);
		this.build = this.build.bind(this);
		this.toggleEditor = this.toggleEditor.bind(this);
		this.getJMODLog = this.getJMODLog.bind(this);
		this.deleteJMOD = this.deleteJMOD.bind(this);

		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="assets/standard.css"></link>
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

#jmod-controls {
	display: flex;
	flex-wrap: wrap;
}
#jmod-controls > button {
	margin: 0.2em;
	flex: 1 1;
}

#jmod-error-output {
	width: 100%;
	color: var(--clr-red);
}

#config-editor-panel {
	display: flex;
	flex-wrap: wrap;
	width: 100%;
	justify-content: flex-end;
}
#config-editor {
	width:100%;
	height: 7em;
	background-color: var(--clr-surface-1);
	color: var(--clr-font-high);
}
</style>
<div class="entry">
	<h3 id="jmod-name" style="word-break: break-all;"></h3>
	<div id="jmod-controls">
		<button onclick="this.getRootNode().host.start()" style="background-color: var(--clr-green)">Start</button>
		<button onclick="this.getRootNode().host.stop()" style="background-color: var(--clr-red)">Stop</button>
		<button onclick="this.getRootNode().host.build()" style="background-color: var(--clr-blue)">Build</button>

		<button onclick="this.getRootNode().host.toggleEditor()" style="background-color: var(--clr-green)">Config</button>
		<button onclick="this.getRootNode().host.getJMODLog()" style="background-color: var(--clr-purple)">Log</button>
		<button onclick="this.getRootNode().host.deleteJMOD()" style="background-color: var(--clr-red)">Delete</button>
	</div>

	<div id="config-editor-panel" style="display:none;">
		<textarea id="config-editor"></textarea>
		<button onclick="this.getRootNode().host.applyConfig();">Apply</button>
		<button onclick="this.getRootNode().host.cancelConfigChange()" style="border-color: var(--clr-red)">Cancel</button>
	</div>

	<div id="jmod-error-output"></div>
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
				if (data.status < 200 || data.status >= 400) {
					throw new Error(await data.text());
				}

				const output = this.shadowRoot.getElementById("jmod-error-output");
				output.textContent = "";

				console.log(await data.text());
			})
			.catch(err => {
				const output = this.shadowRoot.getElementById("jmod-error-output");
				output.textContent = err.message;
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
				if (data.status < 200 || data.status >= 400) {
					throw new Error(await data.text());
				}

				const output = this.shadowRoot.getElementById("jmod-error-output");
				output.textContent = "";

				console.log(await data.text());
			})
			.catch(err => {
				const output = this.shadowRoot.getElementById("jmod-error-output");
				output.textContent = err.message;
				console.error(err);
			});
	}
	
	build() {
		fetch("/admin/buildJMOD", {
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
				console.log(err);
			})
	}

	toggleEditor() {
		const editorPanel = this.shadowRoot.getElementById("config-editor-panel")
		if (editorPanel.style.display === "none") {
			editorPanel.style.display = "flex";
		} else {
			editorPanel.style.display = "none";
		}
	}

	applyConfig() {
		const editor = this.shadowRoot.getElementById("config-editor");

		fetch("/admin/applyJMODConfig", {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				jmodName: this.jmodName,
				newConfig: editor.value
			})
		})
			.then(async data => {
				if (data.status < 200 || data.status >= 400) {
					throw new Error(await data.text());
				}

				const output = this.shadowRoot.getElementById("jmod-error-output");
				output.textContent = "";

				console.log(await data.text());
			})
			.catch(err => {
				const output = this.shadowRoot.getElementById("jmod-error-output");
				output.textContent = err.message;
				console.error(err);
				console.log(err);
			});
	}

	cancelConfigChange() {
		const editor = this.shadowRoot.getElementById("config-editor");
		editor.value = JSON.stringify(this.config, null, "  ");

		this.shadowRoot.getElementById("config-editor-panel").style.display = "none";
	}

	getJMODLog() {
		fetch("/admin/getJMODLog", {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({jmodName: this.jmodName})
		})
			.then(async data => {
				var tab = window.open("about:blank", "_blank");
				tab.document.write(`<p style="white-space: pre;">${await data.text()}</p>`)
				tab.document.close();
			})
			.catch(err => {
				console.error(err);
				console.log(err);
			})
	}

	deleteJMOD() {
		var yes = confirm(`Are you sure you want to delete ${this.jmodName}`);

		if (yes) {
			fetch("/admin/deleteJMOD", {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({jmodName: this.jmodName})
			})
				.then(async data => {
					console.log(await data.text());
					getJMODData();
				})
				.catch(err => {
					console.error(err);
					console.log(err);
				})
		}
	}
}

customElements.define("jmod-entry", JMODEntry);

function getJMODData() {
	fetch("/admin/getJMODData", {})	
		.then(async data => {
			const mods = await data.json();

			const holder = document.getElementById("jmod-entry-holder");
			holder.innerHTML = "";

			Object.entries(mods).forEach((entry) => {
				console.log(entry[0]);
				const newEntry = document.createElement("jmod-entry");
				newEntry.init(entry[0], entry[1]);
				holder.appendChild(newEntry);
			})
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

function createUser(event, formNode) {
	if (event.key != "Enter") {
		return;
	}

	event.preventDefault();

	//console.log(formNode.getElementById("create-user-username"))
	const username = formNode.querySelector("#create-user-username").value;
	const password1 = formNode.querySelector("#create-user-password1").value;
	const password2 = formNode.querySelector("#create-user-password2").value;

	if (password1 !== password2) {
		console.error(new Error("Passwords do not match"));
		alert("Passwords do not match");
		return
	}

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
