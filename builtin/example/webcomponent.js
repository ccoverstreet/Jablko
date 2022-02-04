class extends HTMLElement {
	constructor() {
		super();

		this.attachShadow({mode: open})
		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="/assets/common.css"/>
<div>
</div>
		`

	}

}
