build_instruction = go build -o demomod .
file_dependencies = demomod.go jablkohelpers.go webcomponent.js

run: demomod
	./demomod

build: demomod

demomod: $(file_dependencies)
	$(build_instruction)
