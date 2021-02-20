function interfaceStatus(config) {
	this.id = config.id;
	this.title = config.title;
	this.updateInterval = config.updateInterval;

	this.getStatus = function() {
		fetch(`/jablkomods/${this.id}/getStatus`, {
			method: "POST",
			header: {
				"Content-Type": "application/json"
			}
		})
			.then(async (data) => {
				res = await data.json();
				console.log(res);

				hours = Math.floor((res.uptime) / 3600);
				minutes = Math.floor((res.uptime - hours * 3600) / 60)
				seconds = res.uptime - hours * 3600 - minutes * 60;

				uptimeElem = document.getElementById(`${this.id}_uptime`);
				uptimeElem.textContent = `${hours} : ${minutes} : ${seconds} s`;
				uptimeElem.style.color = "var(--color-good)";

				messageElem = document.getElementById(`${this.id}_message`);
				messageElem.textContent = res.message;
				messageElem.style.color = "var(--color-font-med)";

				messageElem = document.getElementById(`${this.id}_curAlloc`);
				messageElem.textContent = (res.curAlloc / 1000000).toFixed(1) + " MB";

				messageElem.style.color = "var(--color-good)";

				messageElem = document.getElementById(`${this.id}_sysAlloc`);
				messageElem.textContent = (res.sysAlloc / 1000000).toFixed(1) + " MB";
				messageElem.style.color = "var(--color-good)";
			})
			.catch(err => {
				uptimeElem = document.getElementById(`${this.id}_uptime`);
				uptimeElem.textContent = "N/A"
				uptimeElem.style.color = "var(--color-bad)";

				curAllocElem = document.getElementById(`${this.id}_curAlloc`);
				curAllocElem.textContent = "N/A"
				curAllocElem.style.color = "var(--color-bad)";

				sysAllocElem = document.getElementById(`${this.id}_sysAlloc`);
				sysAllocElem.textContent = "N/A"
				sysAllocElem.style.color = "var(--color-bad)";

				messageElem = document.getElementById(`${this.id}_message`);
				messageElem.textContent = "Unable to communicate with interface.";
				messageElem.style.color = "var(--color-bad)";

				console.log("ERROR: Unable to use interfacestatus/getStatus.");
				console.log(err);
			});
	}.bind(this);

	this.speedTest = function() {
		startTime = new Date().getTime();
		fetch(`/jablkomods/${this.id}/speedTest`, {
			method: "POST",
			header: {
				"Content-Type": "application/json"
			}
		})
			.then(async (data) => {
				responseTime = new Date().getTime() - startTime
				speedElem = document.getElementById(`${this.id}_speed`)
				speedElem.textContent = responseTime + " ms"
				speedElem.style.color = "var(--color-good)"
				console.log(await data.json())
			})
			.catch((err) => {
				speedElem = document.getElementById(`${this.id}_speed`)
				speedElem.textContent = "N/A"
				speedElem.style.color = "var(--color-bad)"
			})
	}.bind(this);

	// Setting repeating tasks
	document.addEventListener("DOMContentLoaded", function() {
		this.getStatus();
		this.speedTest();
	}.bind(this));

	setInterval(this.getStatus, this.updateInterval * 1000);
	setInterval(this.speedTest, this.updateInterval * 1000);

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
	<div class="label_value_pair">
		<div>Uptime:</div>
		<div id="${this.id}_uptime" style="color: var(--color-warn)">N/A</div>
	</div>

	<div class="label_value_pair">
		<div>Memory in Use:</div>
		<div id="${this.id}_curAlloc" style="color: var(--color-warn)">N/A</div>
	</div>

	<div class="label_value_pair">
		<div>Virtual Space:</div>
		<div id="${this.id}_sysAlloc" style="color: var(--color-warn)">N/A</div>
	</div>

	<div class="label_value_pair">
		<div>Response Time:</div>
		<div id="${this.id}_speed" style="color: var(--color-warn)">N/A</div>
	</div>

	<div class="label_value_pair">
		<div>Message:</div>
		<div id="${this.id}_message" style="color: var(--color-warn)">N/A</div>
	</div>
</div>
`}.bind(this);
}
