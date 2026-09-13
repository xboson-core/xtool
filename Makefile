GOFILES := $(wildcard *.go src/*.go)

.PHONY: build run

build: $(GOFILES)
	go build

run: build
	./xtool

test1: build
	./xtool --pack -v -i D:\web4xboson\webservice \
		-o test/out.ign.zip \
		--exclude-from test/exclude.txt

test2: build
	./xtool --diff -i D:\web4xboson\webservice \
		-b test/out.ign.zip \
		-del test/dellist.ign.txt \
		-o test/update.ign.zip \
		--exclude-from test/exclude.txt \

test3: build
	./xtool --sqld -b test/b1.sql -i test/b2.sql -o test/b3.sql

test4: build
	./xtool --sqld -b test/a1.sql -i test/a2.sql -o test/a3.sql