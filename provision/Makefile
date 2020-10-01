.PHONY: clean statik

PARTS   ?= compose messaging system
OUTPUTS ?= $(addsuffix /static.go, $(PARTS))
STATIK   = $(GOPATH)/bin/statik

all: $(OUTPUTS)

%/static.go: $(STATIK)
	$(STATIK) -p $(@D) -m -Z -f -src=$(@D)/src

clean:
	rm -f $(OUTPUTS)

$(STATIK):
	$(MAKE) --directory=../ $(STATIK)
