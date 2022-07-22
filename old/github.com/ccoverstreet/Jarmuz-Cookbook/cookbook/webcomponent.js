class extends HTMLElement {
	constructor() {
		super();

		this.getRecipeList = this.getRecipeList.bind(this);
		this.addRecipe = this.addRecipe.bind(this);
		this.getRecipe = this.getRecipe.bind(this);
		this.updateRecipe = this.updateRecipe.bind(this);

		this.attachShadow({mode: "open"});
	}

	init = (source, config) => {
		this.source = source;
		this.config = config;
		console.log(this.config);

		this.getRecipeList()

		this.shadowRoot.innerHTML = `
<link rel="stylesheet" href="/assets/standard.css"/>
<style>
svg > path {
	stroke: var(--clr-accent);
	stroke-width: 30px;
	stroke-linecap: round;
	fill: transparent;
}

.jmod-body {
	display: flex;
	flex-wrap: wrap;
}

#recipe-viewer {
	display: flex;
	flex-direction: column;
	width: 100%;
}
#recipe-viewer > textarea {
	height: 6em;
}

#new-recipe-viewer {
	display: flex;
	flex-direction: column;
	width: 100%;
}
#new-recipe-viewer > textarea {
	height: 6em;
}
</style>
<div class="jmod-wrapper">
	<div class="jmod-header" style="display:flex">
		<h1 id="title">${this.config.title}</h1>
		<svg viewBox="0 0 360 360">
			<path d="M60,60 L60,240  L180,300 L300,240 L300,60 L180,120 L60,60"/>
			<path d="M180,120 L180,300"/>
		</svg>
	</div>

	<hr>

	<div class="jmod-body">

		<div id="recipe-viewer">
			<div style="display: flex; width: 100%; height: 3em;">
				<select id="recipe-selector" style="font-size: 1.25em; flex: 1;"
					onchange="this.getRootNode().host.getRecipe()"></select>
				<button onclick="this.getRootNode().host.showNewRecipeViewer()"
					style="background-color: var(--clr-green)">Add</button>
			</div>

			<h2>Ingredients</h2>
			<textarea id="ingredients-viewer" ></textarea>
			<h2>Instructions</h2>
			<textarea id="instructions-viewer" ></textarea>

			<button onclick="this.getRootNode().host.updateRecipe()"
				style="background-color: var(--clr-green)">Update</button>
			<button onclick="this.getRootNode().host.removeRecipe()"
				style="background-color: var(--clr-red)">Remove Recipe</button>
		</div>

		<div id="new-recipe-viewer" style="display: none;">
			<div style="display: flex; width: 100%; height: 3em;">
				<h2 style="display: flex; align-items: center;">Name</h2>
				<input id="new-recipe-name" 
					style="background-color: var(--clr-background); color: var(--clr-font-high); margin-right: 1em;"></input>
			</div>

			<h2>Ingredients</h2>
			<textarea id="new-recipe-ingredients"></textarea>
			<h2>Instructions</h2>
			<textarea id="new-recipe-instructions"></textarea>
			<button onclick="this.getRootNode().host.addRecipe()"
				style="background-color: var(--clr-green)">Add Recipe</button>
			<button onclick="this.getRootNode().host.showRecipeViewer()"
				style="background-color: var(--clr-red)">Cancel</button>
		</div>
	</div>
</div>
		`
	}

	getRecipeList() {
		fetch(`/jmod/getRecipeList?JMOD-Source=${this.source}`)
			.then(async data => {
				const res = await data.json();
				const selectElem = this.shadowRoot.getElementById("recipe-selector");

				selectElem.innerHTML = "";

				// Create empty option for default
				var base = document.createElement("option");
				base.value = "";
				base.textContent = "";
				selectElem.appendChild(base);

				for (name of res) {
					var opt = document.createElement("option");
					opt.value = name;
					opt.textContent = name;
					selectElem.appendChild(opt);
				}
			})
			.catch(err => {
				console.error(err);
				console.log(err);
			})
	}

	addRecipe() {
		const name = this.shadowRoot.getElementById("new-recipe-name").value;
		if (name == "") {
			alert("Recipe name cannot be empty");
			return;
		}
		const ingredients = this.shadowRoot.getElementById("new-recipe-ingredients").value;
		const instructions = this.shadowRoot.getElementById("new-recipe-instructions").value;

		const body = {name, ingredients, instructions};
		console.log(body);

		fetch(`/jmod/addRecipe?JMOD-Source=${this.source}`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify(body)
		})
			.then(async data => {
				console.log(await data.text());
				this.getRecipeList();
				this.showRecipeViewer();
				this.shadowRoot.getElementById("new-recipe-name").value = "";
				this.shadowRoot.getElementById("new-recipe-ingredients").value = "";
				this.shadowRoot.getElementById("new-recipe-instructions").value = "";
			})
			.catch(err => {
				console.error(err);
				console.log(err);
			})
	}

	removeRecipe() {	
		const name = this.shadowRoot.getElementById("recipe-selector").value;
		if (name == "") {
			return;
		}
		
		if (!confirm(`Are you sure you want to delete "${name}"`)) {
			return
		}

		fetch(`/jmod/removeRecipe?JMOD-Source=${this.source}`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({name})
		})
			.then(async data => {
				console.log(await data.text());
				this.shadowRoot.getElementById("ingredients-viewer").value = "";
				this.shadowRoot.getElementById("instructions-viewer").value = "";
				this.getRecipeList();
			})
			.catch(err => {
				console.error(err);
				console.log(err);
			})
	}

	getRecipe() {
		const name = this.shadowRoot.getElementById("recipe-selector").value;
		if (name == "") {
			return;
		}

		fetch(`/jmod/getRecipe?JMOD-Source=${this.source}`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({name})
		})
			.then(async data => {
				var res = await data.json();
				console.log(res);

				this.shadowRoot.getElementById("ingredients-viewer").value = res.ingredients;
				this.shadowRoot.getElementById("instructions-viewer").value = res.instructions;
			})
			.catch(err => {
				console.error(err);
				console.log(err);
			})
	}

	updateRecipe() {
		const reqBody = {
			name: this.shadowRoot.getElementById("recipe-selector").value,
			ingredients: this.shadowRoot.getElementById("ingredients-viewer").value,
			instructions: this.shadowRoot.getElementById("instructions-viewer").value
		};

		fetch(`/jmod/updateRecipe?JMOD-Source=${this.source}`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify(reqBody)
		})
			.then(async data => {
				console.log(await data.text());
			})
			.catch(err => {
				console.error(err);
				console.log(err);
			})
	}

	showNewRecipeViewer() {
		this.shadowRoot.getElementById("recipe-viewer").style.display = "none";
		this.shadowRoot.getElementById("new-recipe-viewer").style.display = "flex";
	}

	showRecipeViewer() {
		this.shadowRoot.getElementById("recipe-viewer").style.display = "flex";
		this.shadowRoot.getElementById("new-recipe-viewer").style.display = "none";
	}
}
