// Webcomponent for Jablko Demo mod
// Cale Overstreet
// May 15, 2021

class extends HTMLElement {
	constructor() {
		super();

		this.attachShadow({mode: "open"});
		this.webSocketResHandler = this.webSocketResHandler.bind(this);
		this.getUDPState = this.getUDPState.bind(this);
		this.testSave = this.testSave.bind(this);

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
		<h1>Test Module</h1>
		<svg viewBox="0 0 360 360">
			<circle cx="180" 
				cy="180" 
				r="120" 
				stroke="var(--clr-accent)" 
				stroke-width="30"
				fill="transparent"/>		
		</svg>
	</div>
	<hr>
	<div class="jmod-body">
		<p>
		Demo for differenent web technologies in Jablko
		</p>
		<div id="websocket">
			<h2 style="width: 100%;">Web Socket</h2>
			<input onkeypress="this.getRootNode().host.inputEventHandler(this, event)" style="flex-grow: 1;"></input>		
			<div id="websocket-output"></div>
		</div>

		<div>
			<button onclick="this.getRootNode().host.getUDPState()">Get UDP State</button>
			<h3>UDP State:</h3>
			<div id="udpstate-output"></div>
		</div>

		<button onclick="this.getRootNode().host.talk()" style="border-color: var(--clr-red);">Talk</button>
		<button onclick="this.getRootNode().host.testSave()">Test Save</button>
	</div>
</div>
		`
	}

	init(source, instName, config) {
		// Setup WebSocket
		try {
			this.webSocket = new WebSocket(`ws://${document.location.host}/jmod/socket?JMOD-Source=${source}`);
			this.webSocket.onmessage = this.webSocketResHandler;
		} catch(err) {
			console.error(err);
			return;
		}

		this.source = source;
		this.instName = instName;
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

	getUDPState() {
		fetch(`/jmod/getUDPState?JMOD-Source=${this.source}`)
			.then((async data => {
				var res = await data.json();	
				var elem = this.shadowRoot.getElementById("udpstate-output");
				elem.textContent = res.state;
			}).bind(this))
			.catch((err => {
				console.log(this.shadowRoot);
				console.log(this.shadowRoot.getElementById("udpstate-output"));
				var elem = this.shadowRoot.getElementById("udpstate-output");
				elem.textContent = err;		
			}).bind(this));
	}

	talk() {
		alert("Hello");
	}

	testSave() {
		console.log("HelloWorld");
		fetch(`/jmod/testConfigSave?JMOD-Source=${this.source}`)
	}
}
