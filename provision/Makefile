.PHONY: clean all

include ../Makefile.inc

PARTS   ?= compose messaging system
OUTPUTS ?= $(addsuffix /static.go, $(PARTS))

all: $(OUTPUTS)

%/static.go: $(STATIK)
	$(STATIK) -p $(@D) -m -Z -f -src=$(@D)/src

clean:
	rm -f $(OUTPUTS)
