class extends HTMLElement {
	constructor() {
		super();

		this.attachShadow({mode: "open"});
		this.webSocketResHandler = this.webSocketResHandler.bind(this);

		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="/assets/standard.css"></link>
<style>
#websocket {
	display: flex;
	flex-wrap: wrap;
}
#websocket > div {
	padding: 0.5em;
	flex: 1 1 10em;
}
</style>
<div class="jmod-wrapper">
	<div class="jmod-header" style="display: flex; ">
		<h1>My Module</h1>
		<div style="flex-grow: 1;"></div>
		<svg viewBox="0 0 360 360">
			<circle cx="180" 
				cy="180" 
				r="90" 
				stroke="var(--clr-accent)" 
				stroke-width="30"
				fill="transparent"/>		
		</svg>
	</div>
	<hr>
	<div class="jmod-body">
		<div id="websocket">
			<h2 style="width: 100%;">Web Socket</h2>
			<input onkeypress="this.getRootNode().host.inputEventHandler(this, event)" style="flex-grow: 1;"></input>		
			<div id="websocket-output"></div>
		</div>
		<button onclick="this.getRootNode().host.talk()" style="border-color: var(--clr-red);">Talk</button>
		<button onclick="this.getRootNode().host.talk()" style="border-color: var(--clr-green);">Talk</button>
		<button onclick="this.getRootNode().host.talk()" style="border-color: var(--clr-yellow);">Talk</button>
		<button onclick="this.getRootNode().host.talk()" style="border-color: var(--clr-accent);">Talk</button>
	</div>
</div>
		`
	}

	init(source, config) {
		// Setup WebSocket
		try {
			this.webSocket = new WebSocket(`ws://${document.location.host}/jmod/socket?JMOD-Source=${source}`);
			this.webSocket.onmessage = this.webSocketResHandler;
		} catch(err) {
			console.error(err);
			return;
		}
	}

	webSocketResHandler(event) {
		var elem = this.shadowRoot.getElementById("websocket-output");
		elem.textContent = event.data;
		console.log(event.data);
	}

	inputEventHandler(elem, event) {
		if (event.key == "Enter") {
			this.webSocket.send(elem.value);
			elem.value = "";
		}
	}

	talk() {
		alert("Hello");
	}
}
