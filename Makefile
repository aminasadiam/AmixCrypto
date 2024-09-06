build:
	go build . -o ./dist/Amix.exe

run:
	go build . -o ./dist/Amix.exe
	./dist/Amix.exe

test:
	go test