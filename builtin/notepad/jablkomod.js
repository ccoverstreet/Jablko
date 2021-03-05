function notepad(config) {
	this.id = config.id;
	this.title = config.title;

	// Setting repeating tasks
	document.addEventListener("DOMContentLoaded", function() {
	}.bind(this));

	this.card = function() {
		return `
<div id="${this.id}" class="module_card">
	<div class="module_title">
		<div>${this.title}</div>
		<svg class="module_icon" viewBox="0 0 150 150">
			<path d="M 20 75 H 40 L 60 30 L 90 120 L 110 75 H 130" fill="transparent" stroke="var(--color-primary)" stroke-width="20px" stroke-linejoin="round" stroke-linecap="round"/>
		</svg>
	</div>
</div>
`}.bind(this);
}
