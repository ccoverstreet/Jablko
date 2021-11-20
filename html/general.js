class JPrompt extends HTMLElement {
	constructor() {
		super();

		this.attachShadow({mode: "open"});

		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="/assets/standard.css"/>
<style>
label {
	display: inline-block;
	font-weight: bold;
	width: calc(100% - 1em);
	text-align: end;
	padding: 0.25em 0.75em;
}


#grayout {
	width: 100vw;
	height: calc(100vh - 3em);
	background-color: rgba(0, 0, 0, 0.6);
	overflow: auto;
	display: flex;
	justify-content: center;
}

#box {
	display: flex;
	flex-direction: column;
	width: min(max((100%), 33vw), 30em);
	margin: 0px;
	padding: 0px;
}

#container {
	width: calc(100% - 4em);
	display: grid;
	grid-template-columns: 50% 50%;
	grid-row-gap: 0.25em;
	background-color: rgba(38, 38, 38, 1);
	padding: 1em;
	margin: 1em;
	border-radius: 0.25em;
}

input {
	width: calc(100% - 1em);
	font-size: 1em;
}

#controls {
	width: calc(100% - 2em);
	display: flex;
	justify-content: flex-end;
	margin: 0.25em 1em;
	gap: 0.5em;
}

#controls > button {
	margin: 0;
}
		</style>

<div id="grayout">
	<div id="box">
		<div id="container">
		</div>
		<div id="controls">
			<button style="background-color: var(--clr-green); font-weight: bold;" onclick="this.getRootNode().host.submit();">Submit</button>
			<button style="background-color: var(--clr-red); font-weight: bold;" onclick="this.getRootNode().host.remove();">X</button>
		</div>
	</div>
</div>
		`

		this.shadowRoot.addEventListener("keydown", (event) => {
			if (event.key == "Enter") {
				this.submit();
			}
		});
	}

	connectedCallback() {
		console.log("ADDING");
	}


	init = (config, callback) => {
		this.config = config;
		this.callback = callback;
		let container = this.shadowRoot.getElementById("container");
		for (const field of config) {
			let label = document.createElement("label");
			label.textContent = field.label + ":";

			if (field.type === "input") {
				let data = document.createElement("input");
				data.id = field.id;
				data.value = field.value ? field.value : "";
				container.appendChild(label);
				container.appendChild(data);
			} else if (field.type === "checkbox") {
				let data = document.createElement("input");
				data.type = "checkbox";
				data.style.textAlign = "left";
				data.style.marginLeft = "1em";
				data.style.width = "auto";
				data.id = field.id;
				data.checked = field.value ? true : false;
				container.appendChild(label);
				container.appendChild(data);
			} else if (field.type === "textarea") {
				let data = document.createElement("textarea");
				data.id = field.id;
				data.value = field.value ? data.value : "";
				label.style.textAlign = "left";
				label.style.gridColumn = "1 / span 2";
				data.style.gridColumn = "1 / span 2";
				container.appendChild(label);
				container.appendChild(data);
			}
		}
	}

	submit = () => {
		console.log("Submit");
		let output = {};
		for (const field of this.config) {
			let val = null;
			if (field.type === "input") {
				val = this.shadowRoot.getElementById(field.id).value;
			} else if (field.type === "checkbox") {
				val = this.shadowRoot.getElementById(field.id).checked;
			} else if (field.type === "textarea") {
				val = this.shadowRoot.getElementById(field.id).value;
			}
			output[field.id] = val;
		}

		// Send data to callback function
		this.callback(output, this);
	}
}

class JConfirm extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({mode: "open"});

		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="/assets/standard.css"/>
<style>
#grayout {
	width: 100vw;
	height: calc(100vh - 3em);
	background-color: rgba(0, 0, 0, 0.6);
	overflow: auto;
	display: flex;
	justify-content: center;
}i

#box {
	display: flex;
	flex-direction: column;
	width: min(max((100%), 33vw), 40em);
	margin: 0px;
	padding: 0px;
}

#container {
	width: calc(100% - 4em);
	display: block;
	background-color: rgba(38, 38, 38, 1);
	padding: 1em;
	margin: 1em;
	border-radius: 0.25em;
}

#controls {
	width: calc(100% - 2em);
	display: flex;
	justify-content: flex-end;
	margin: 0.25em 1em;
	gap: 0.5em;
}
</style>
<div id="grayout">
	<div id="box">
		<div id="container">
		</div>
		<div id="controls">
			<button style="background-color: var(--clr-green); font-weight: bold;" onclick="this.getRootNode().host.confirm();">Confirm</button>
			<button style="background-color: var(--clr-red); font-weight: bold;" onclick="this.getRootNode().host.cancel();">X</button>
		</div>
	</div> 
</div>
		`
	}

	init = (prompt) => {
		let container = this.shadowRoot.getElementById("container");
		let p = document.createElement("p");
		p.textContent = prompt;
		p.width = "100%";
		p.style.fontSize = "1em";
		container.appendChild(p);
	}

	promiseFunc = (resolve, reject) => {
		this.resolve = resolve;
		this.reject = reject;
	}

	confirm = () => {
		this.resolve(true);
		this.remove();
	}

	cancel = () => {
		this.resolve(false);
		this.remove();
	}
}

class JAlert extends HTMLElement {
	constructor() {
		super();

		this.attachShadow({mode: "open"});
		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="/assets/standard.css"/>
<style>
#grayout {
	width: 100vw;
	height: calc(100vh - 3em);
	background-color: rgba(0, 0, 0, 0.6);
	overflow: auto;
	display: flex;
	justify-content: center;
}i

#box {
	display: flex;
	flex-direction: column;
	width: min(max((100%), 33vw), 40em);
	margin: 0px;
	padding: 0px;
}

#container {
	width: calc(100% - 4em);
	display: block;
	background-color: rgba(38, 38, 38, 1);
	padding: 1em;
	margin: 1em;
	border-radius: 0.25em;
}

#controls {
	width: calc(100% - 2em);
	display: flex;
	justify-content: flex-end;
	margin: 0.25em 1em;
	gap: 0.5em;
}

</style>
<div id="grayout">
	<div id="box">
		<div id="container">
			<p id="message"></p>
		</div>
		<div id="controls">
			<button style="background-color: var(--clr-red); font-weight: bold;" onclick="this.getRootNode().host.close();">X</button>
		<div>
	</div> 
</div>
		`
	}

	init = (message) => {
		const elem = this.shadowRoot.getElementById("message");
		elem.textContent = message;
	}

	close = () => {
		this.remove();
	}
}

// Register element
customElements.define("j-prompt", JPrompt);
customElements.define("j-confirm", JConfirm);
customElements.define("j-alert", JAlert);

const jablko = {
	// Customizable prompt
	prompt: (config, callback) => {
		let newElem = document.createElement("j-prompt");
		newElem.init(config, callback);
		newElem.style.position = "fixed";
		newElem.style.top = "3em";
		newElem.style.left = "0";
		document.querySelector("body").appendChild(newElem);
	},

	// Jablko specific confirm prompt
	confirm: (prompt) => {
		let newElem = document.createElement("j-confirm");
		newElem.init(prompt);
		newElem.style.position = "fixed";
		newElem.style.top = "3em";
		newElem.style.left = "0";
		document.querySelector("body").appendChild(newElem);
		return new Promise(newElem.promiseFunc);
	},

	alert: (message) => {
		let newElem = document.createElement("j-alert");
		newElem.init(message);
		newElem.style.position = "fixed";
		newElem.style.top = "3em";
		newElem.style.left = "0";
		document.querySelector("body").appendChild(newElem);
	}
}
