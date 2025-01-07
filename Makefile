build:
	go build -o main.exe

run: build
	./main.exe

clean:
	del /f main.exe