function notepad(config) {
	// Properties
	this.id = config.id;
	this.title = config.title;

	this.saveNote = function() {
		var textElem = document.getElementById(`${this.id}_textarea`);
		console.log(textElem.value);

		fetch(`/jablkomods/${this.id}/saveNote`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({text: textElem.value})
		})
			.then(async data => {
				console.log(await data.text());
			})
			.catch(err => {
				console.error(err);
			})
	}.bind(this);

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
	<hr>

	<div style="display: flex; justify-content: center;">
		<textarea id="${this.id}_textarea" style="font-family: Arial; font-size: 14px; width: 90%; padding: 5px; height: 100px;"></textarea>
	</div>
</div>
`}.bind(this);
}
