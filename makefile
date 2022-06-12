Jablko: rewrite.go core/*/*.go
	go build -o Jablko rewrite.go
	
run: Jablko
	./Jablko

clean:
	rm -r github.com log tmp
