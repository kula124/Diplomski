OUTPUT_FOLDER = ./bin
MAIN_FILE = ./src/main.go
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

build:
	go build -o $(OUTPUT_BIN) $(MAIN_FILE)

run: build
	$(OUTPUT_BIN) $(ARGS)

clean:
	$(RM) $(CLEAN_FOLER)
	$(MKDIR) $(CLEAN_FOLER)
