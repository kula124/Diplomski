OUTPUT_FOLDER = ./bin
MAIN_FILE = ./src/main.go ./src/encryption.go
ARGS = ""
ifeq ($(OS),Windows_NT)
		RM = cmd.exe /c rmdir /s /q
		MKDIR = cmd.exe /c mkdir
		CLEAN_FOLER = bin
    OUTPUT_BIN = $(OUTPUT_FOLDER)/main.exe
else
    #Linux stuff
		RM = rm -rf
		MKDIR = mkdir
    OUTPUT_BIN = $(OUTPUT_FOLDER)/main
		CLEAN_FOLER = $(OUTPUT_FOLDER)
endif

test:
	go test -v  .\src\cli
	go test -v .\src\crypto
	go test -v .\src\utils
	go test -v .\src\config
	go test -v .\src\

build: test
	go build -o $(OUTPUT_BIN) $(MAIN_FILE)

run: build
	$(OUTPUT_BIN) $(ARGS)
exec:
	${OUTPUT_BIN} ${ARGS}

clean:
	$(RM) $(CLEAN_FOLER)
	$(MKDIR) $(CLEAN_FOLER)
#clearc:
 # go clean -cache -modcache -i -r
