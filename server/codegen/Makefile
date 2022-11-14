.PHONY: *.cue

DEF_ALL    = *.cue
DEF_SERVER = server.*.cue
DEF_DOCS   = docs.*.cue

server: $(DEF_SERVER)
docs:   $(DEF_DOCS)
all:    $(DEF_SERVER) $(DEF_DOCS)

include ../Makefile.inc

$(DEF_SERVER): $(CUE) $(JSONTPLEXEC)
	@ echo "$(COLOUR_GREEN)Generating server files from $@$(COLOUR_END)"
	@ $(CUE) eval $@ --out json | $(JSONTPLEXEC) -v -p $(ROOT_DIR)/codegen/assets/templates -b $(ROOT_DIR)

ifeq ("${DOCS_DIR}", "")
$(DEF_DOCS):
	@ echo "$(COLOUR_YELLOW)Skipping docs generation: DOCS_DIR is not set.$(COLOUR_END)"
	@ echo "$Run task like this: make "$@" DOCS_DIR=/path/to/corteza-docs"
	@ echo ""
else
$(DEF_DOCS): $(CUE) $(JSONTPLEXEC)
	@ echo "$(COLOUR_GREEN)Generating doc files from $@$(COLOUR_END) $(COLOUR_BLUE)(dst: $(DOCS_DIR))$(COLOUR_END)"
	@ $(CUE) eval $@ --out json | $(JSONTPLEXEC) -v -p $(ROOT_DIR)/codegen/assets/templates -b $(DOCS_DIR)
endif
