GOFILES := $(wildcard *.go src/*.go)

.PHONY: build run

build: $(GOFILES)
	go build

run: build
	./xtool

test1: build
	./xtool --pack -v -i D:\web4xboson\webservice \
		-o test/out.zip \
		--exclude-from test/exclude.txt

test2: build
	./xtool --diff -i D:\web4xboson\webservice \
		-b test/out.zip \
		-del test/dellist.txt \
		-o test/update.zip \
		--exclude-from test/exclude.txt \
