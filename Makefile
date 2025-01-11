build:
	go build -o main.exe ./cmd/main.go

run: build
	./main.exe

clean:
	del /f main.exe