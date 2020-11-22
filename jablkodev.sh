#!/bin/bash

cd test_modules/interfacestatus 
ls
go build -buildmode=plugin -o jablkomod.so .
go test
if [ $? -gt 0 ]
then
	exit 1
fi
cd ../..

echo "Finished Compiling Test Modules"
go run jablko.go
