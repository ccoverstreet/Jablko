Jablko: main.go core/*/*.go
	go build -o Jablko main.go
	
run: Jablko
	./Jablko

clean:
	rm -r github.com log tmp
