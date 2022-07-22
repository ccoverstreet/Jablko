// Jarmuz RGB Light Webcomponent
// Cale Overstreet
// Jun. 27, 2021

/* UI component for Jablko. Interacts through WebSocket
 */

class extends HTMLElement {
	constructor() {
		super();

		this.setRGBA = this.setRGBA.bind(this);
		this.sliderUpdate = this.sliderUpdate.bind(this);
		this.changeTarget = this.changeTarget.bind(this);

		this.r = 120;
		this.g = 120;
		this.b = 120;
		this.a = 120;

		this.attachShadow({mode: "open"});
		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="/assets/standard.css"/>
<style>
#slider_holder {
	display: flex;
	flex-direction: column;
	height: 12em;
}
#slider_holder > input {
	-webkit-appearance: none;
	width: 80%;
	height: 0.2em;
	background-color: var(--clr-red);
}
</style>
<div class="jmod-wrapper">
	<div class="jmod-header" style="display:flex">
		<h1>RGB Light</h1>
		<svg viewBox="0 0 360 360">
			<path d="M150,300 A30,60,0,0,1,210,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
			<path d="M105,300 A60,90,0,0,1,255,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
			<path d="M60,300 A60,80,0,0,1,300,300" stroke="var(--clr-accent)" stroke-width="30" stroke-linecap="round" fill="transparent"/>
		</svg>
	</div>

	<hr>

	<div class="jmod-body" style="display: flex; justify-content: center; flex-direction: column">
		<select id="light-target" onchange="this.getRootNode().host.changeTarget(this.value)" style="margin: auto 4em; font-size: 1.25em;">
		</select>

		<div id="slider_holder">
			<input type="range" id="slider_r" hmin="0" max="255" value="120"
				oninput="this.getRootNode().host.sliderUpdate('r', parseInt(this.value, 10))"
				onchange="this.getRootNode().host.sliderFinal('r', parseInt(this.value, 10))"
				style="">
			</input>

			<input type="range" id="slider_g" hmin="0" max="255" value="120"
				oninput="this.getRootNode().host.sliderUpdate('g', parseInt(this.value, 10))"
				onchange="this.getRootNode().host.sliderFinal('g', parseInt(this.value, 10))"
				style="background-color: var(--clr-green);">
			</input>

			<input type="range" id="slider_b" hmin="0" max="255" value="120"
				oninput="this.getRootNode().host.sliderUpdate('b', parseInt(this.value, 10))"
				onchange="this.getRootNode().host.sliderFinal('b', parseInt(this.value, 10))"
				style="background-color: var(--clr-blue);">
			</input>

			<input type="range" id="slider_a" hmin="0" max="255" value="120"
				oninput="this.getRootNode().host.sliderUpdate('a', parseInt(this.value, 10))"
				onchange="this.getRootNode().host.sliderFinal('a', parseInt(this.value, 10))"
				style="background-color: var(--clr-font-med);">
			</input>
		</div>
	</div>
</div>
		`
	}

	init(source, config) {
		this.source = source;
		this.config = config;
		this.lastMessageTime = performance.now();

		try {
			this.currentTarget = this.config.lightIPs[0]
		} catch (err) {
			console.error(err);
			console.log(err);
		}

		this.socket = new WebSocket(`ws://${document.location.host}/jmod/socket?JMOD-Source=${this.source}`);

		const select = this.shadowRoot.getElementById("light-target");

		// Add options to select element for light target
		for (var x of this.config.lightIPs) {
			var opt = document.createElement("option");
			opt.value = x;
			opt.innerText = x;
			select.appendChild(opt);
		}
	}

	changeTarget(target) {
		this.currentTarget = target;
	}

	sliderUpdate(color, value) {
		this[color]	= value

		// Rate limiting to 10 Hz
		var curTime = performance.now();
		if (curTime - this.lastMessageTime < 100) {
			return 
		}

		this.lastMessageTime = curTime

		this.setRGBA();
	}

	sliderFinal(color, value) {
		this[color] = value;
		this.setRGBA();
	}

	setRGBA() {
		this.socket.send(`${this.currentTarget},${this.r},${this.g},${this.b},${this.a}`);
		console.log(this.r, this.g, this.b, this.a);
	}
}
