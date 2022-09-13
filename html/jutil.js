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
		confirm.onclick = async function() {
			try {
				await callback(this.getRootNode().host);
			} catch(e) {
				alert(e);
			}
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

customElements.define("jablko-prompt", JablkoPrompt);
