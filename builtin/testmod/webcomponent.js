class extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({mode: "open"});

		this.shadowRoot.innerHTML = `
		<link rel="stylesheet" href="/assets/standard.css"/>
		<div class="mod-card">
			<div class="mod-header"></div>

			<div class="mod-body">
				<form onsubmit="this.getRootNode().host.sendTestFunc(event)">
					<label>
						Test func:
						<input id="test-input"></input>
					</label>
				</form>
				<p id="test-output" style="white-space: pre-wrap">Test output</p>
			</div>
		</div>
		`
	}	

	init = (modName) => {
		this.modName = modName;
		console.log(modName);
	}

	sendTestFunc = (event) => {
		event.preventDefault();
		const funcName = this.shadowRoot.querySelector("#test-input").value

		fetch(`/mod/${funcName}?modName=${this.modName}`)
			.then(async res => {
				const data = await res.json();
				const print = JSON.stringify(data, 0, "    ");
				console.log(print);
				this.shadowRoot.querySelector("#test-output").textContent = print;
			})
	}
}
