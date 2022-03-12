.PHONY: *.cue

DEF_ALL    = *.cue
DEF_SERVER = server.*.cue
DEF_DOCS   = docs.*.cue

all:    $(DEF_ALL)
server: $(DEF_SERVER)
docs:   $(DEF_DOCS)

include ../Makefile.inc

$(DEF_ALL): $(CUE) $(JSONTPLEXEC)
	@ echo "$(COLOUR_GREEN)Generating from $@$(COLOUR_END)"
	@ $(CUE) eval $@ --out json | $(JSONTPLEXEC) -v -p ../codegen/assets/templates -b ../

gen:
	@ rm $(JSONTPLEXEC)
