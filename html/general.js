const alertSound = () => {
	const snd = new Audio("data:audio/mpeg;base64,//NAxAANgIKGP0VIArY2QDjidvw1f1I1TnzknOgcAAUEwQDDE0YWOEJ/8Mcpzn/8o7/92woc4f///wQDMMF3h/lNTvv9/+AABfnbG9Jrf6NRdgvQWereUxTFRZ5URogmophl2ZgQoFn/80LEJR3qIspfiJgAgzhTuFpKTJlw84ypdOEPWakXMc0zMdZIHzBhzCqfpu3orY+iqbUUWMd0GTRezumbKdXsi6kLftZBkG9XVnCj9pF8MFEXlIcEDjkTISHn//+v/0JPNsMgT1vf9tX/80DECRaxRqHt25gAmzrG5eu6v6rb79WzUbIrcrMYh2AbfkNGRIEibmpo7OiapOyCkHWupJQsp1KZ0UE3Qremm9FFl9Ew5W9JY6DAhhMOuY5/c9ZbahDzosGTjkvSRld3iKBJZLtQBP/zQsQJFVGq6v5LRtczEzo9IgE5pvg+jgyUhVzY6gvQtvpInOpJ00HektJTGATlZ5PV16Fjtytz4QnvnShTci3IsIFiCGBXEAw//9T+U8YO6Gpk3VV7vt+mKlh5mMHddttuBF8xuJ8gbP/zQMQPFpGK9v5ixLY7Mz54iCEd0gjc38oAVBoxd/LZVWPX8TaSksmoD4XQ+5ksDLCGqCOooNXa35tladDiHnodcFRRvb63ZATlzI8WEwejK9ejutIx7mPqaIh3IWu2u2+BhLl9YWJD//NCxA8W4t7y/kjFZgTiwqqM1aj/5+BUAREea//TdNWjO768YZThsRI2SLFo3arijVxgVW8FFCDD/tHanXqhNEZZqSn////35ipWqUIv6P+2nr+DtAabX/7qW00AZxq2QC/m64SQNQtb//NAxA8UINqyPjYQiCCFd0H5dhLxhazEMjCJLp6i8QRC5v9fBAo5RIwes8et2bqoEoVEj0AAaHzq/Qx95lbpJ4V/uWtU7iWGgo/2f8WeqrrrEOGffbfAA4zl66eSNcbAuEwBhStpe/3/80LEGRRYxup+RlKqptKHlxI9A8ZjMM1b3g26ftX3OMkwoA4WAwqLnhEoBvgogUFcQgkbNDH+486zwsZALgxSF5Z6L0V1OtIAjlkkuABtuNRUGotePqndg1K3VmfrQiA1DFI1Mfy7+7r/80DEIxQZdqpeHgaQ/k9/21GLMw1y4KaTf7bySui/nw5rc47EQoCgAHYCYFSzP2YjchU0u5a0PYy9KnKmgJJLJLRAv+S6NM7eTafMAQbp57TdBEyROWZbDWXjicq9XSciurkdga15r//zQsQtFEpCxl4zRK4JXrUTMy0RUnooNFLy6UvRlS/o/t97eiOUd2JGYaewL1J///VVe8AAdu1vuwEa3lfzVh1xSbVIDw5AEgVmMb8rbJFFhzrWUDiKO1D0c50NY85Lo7PT7IpqoyO52//zQMQ3FDGurl9POAC8x71VtnOa6WXnEQwFlBFtuz0NbN7aP/9GNQIoE5XJ3S/7ivhGY9/mDgHVTqBPP/gBQCMQObgwMDc4ZLhiQgA0VEWItxekqbzEmRzTLzFM0IogRYgRiXS78gZF//NCxEEhEnpgqZqYAMniKFQxP1oqLxeR/tTPnzAuGhmXS6kkl1/83c8YonGm506YIIoqSMloo+l//myVF3PCA+FkBUFREerCQNf6AiH/6wVVTEFNRTMuMTAwVVVVVVVVVVVVVVVVVVVV//NAxBgAAANIAcAAAFVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVU=");
	snd.play();
}

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
	width: calc(100% - 2.5em);
	display: grid;
	grid-template-columns: calc(100% - 3em) 3em;
	background-color: var(--clr-gray);
	padding: 1em;
	margin: 0.25em;
	border-radius: 0.5em;
}

</style>

<div id="container">
	<p id="message"></p>
	<button style="background-color: var(--clr-red); font-weight: bold; height: 2em; width: 2em" onclick="this.getRootNode().host.close();">X</button>
</div>
		`
	}

	init = (message, lifetime) => {
		const elem = this.shadowRoot.getElementById("message");
		elem.textContent = message;

		if (!lifetime) {
			return
		}

		setTimeout(this.close, lifetime);
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

	alert: (message, lifetime) => {
		alertSound();
		// Check if alert holder exists
		var alertBox = document.getElementById("jablko-alert-box");
		if (alertBox === null) {
			alertBox = document.createElement("div");
			alertBox.id = "jablko-alert-box";
			alertBox.style.position = "fixed";
			alertBox.style.top = "3em";
			alertBox.style.right = "0";
			alertBox.style.width = "20em";
			alertBox.style.display = "flex";
			alertBox.style.flexDirection = "column";
			alertBox.style.zIndex = 1000;

			document.querySelector("body").appendChild(alertBox);
		}

		let newElem = document.createElement("j-alert");
		newElem.init(message, lifetime);
		setTimeout(() => alertBox.appendChild(newElem), 100);
	}
}
