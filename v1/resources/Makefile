EXE = $(shell go env GOEXE)
SOURCES = $(wildcard *.go)
TESTER = ./resources.test$(EXE)
ZIP = ./embed.zip
DIR = ./embed

test : $(TESTER)
	$(TESTER) -test.v

clean :
	go clean -x
	rm -f $(ZIP) $(EXE)
	rm -rf $(DIR)

$(TESTER) : $(ZIP) $(SOURCES)
	go test -c -o $@
	cat $(ZIP) >> $@

$(ZIP) : $(SOURCES)
	mkdir -p $(DIR)
	cp *.go $(DIR)
	zip -r $(ZIP) $(DIR)
	rm -rf $(DIR)
