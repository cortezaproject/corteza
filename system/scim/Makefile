.PHONY: clean all

include ../../Makefile.inc

all: static.go

static.go: $(STATIK)
	$(STATIK) -p assets -m -Z -f -src=$(@D)/assets

clean:
	rm -f static.go
