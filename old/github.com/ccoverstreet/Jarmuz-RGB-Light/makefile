out = jarmuzrgblight
build_command = go build -o $(out) .

run: $(out)
	./$(out)

build: $(out)

$(out): *.go *.js
	$(build_command)
