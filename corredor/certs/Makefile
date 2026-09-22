# Helper tool to create CA, server and client certificate pairs

OPENSSL   ?= /usr/bin/openssl

HOSTNAME  ?= "localhost"

DAYS      ?= 3650
NUMBITS   ?= 4096

# change/set like this:
# $ make PASSPARSE=fff
PASSPARSE        ?= demo-pass-parse-must-change!

CA_SUBJ          ?= /O=Human/OU=Test/CN=localhost
CA_DIR           ?= ./ca
CA_KEY           ?= $(CA_DIR)/private.key
CA_CRT           ?= $(CA_DIR)/public.crt
CA_DAYS          ?= $(DAYS)
CA_PASSPARSE     ?= $(PASSPARSE)
CA_NUMBITS       ?= $(NUMBITS)

SERVER_SUBJ      ?= /O=Human/OU=Server/CN=$(HOSTNAME)
SERVER_DIR       ?= ./server
SERVER_KEY       ?= $(SERVER_DIR)/private.key
SERVER_CSR       ?= $(SERVER_DIR)/sign-req.csr
SERVER_CRT       ?= $(SERVER_DIR)/public.crt
SERVER_CA_CRT    ?= $(SERVER_DIR)/ca.crt
SERVER_DAYS      ?= $(DAYS)
SERVER_PASSPARSE ?= $(PASSPARSE)
SERVER_NUMBITS   ?= $(NUMBITS)

CLIENT_SUBJ      ?= /O=Human/OU=Client/CN=$(HOSTNAME)
CLIENT_DIR       ?= ./client
CLIENT_KEY       ?= $(CLIENT_DIR)/private.key
CLIENT_CSR       ?= $(CLIENT_DIR)/sign-req.csr
CLIENT_CRT       ?= $(CLIENT_DIR)/public.crt
CLIENT_CA_CRT    ?= $(CLIENT_DIR)/ca.crt
CLIENT_DAYS      ?= $(DAYS)
CLIENT_PASSPARSE ?= $(PASSPARSE)
CLIENT_NUMBITS   ?= $(NUMBITS)

all_certs: ca_certs server_certs client_certs

ca_certs: $(CA_CRT)

$(CA_DIR):
	@ mkdir -p $(CA_DIR)

$(CA_KEY): $(CA_DIR)
	@ $(OPENSSL) genrsa -passout "pass:$(CA_PASSPARSE)" -des3 -out "$(CA_KEY)" $(CA_NUMBITS)

$(CA_CRT): $(CA_DIR) $(CA_KEY)
	@ $(OPENSSL) req -passin "pass:$(CA_PASSPARSE)" -new -x509 -days $(CA_DAYS) -key "$(CA_KEY)" -out "$(CA_CRT)" -subj "$(CA_SUBJ)"


server_certs: $(SERVER_CRT)

$(SERVER_DIR):
	@ mkdir -p $(SERVER_DIR)

$(SERVER_KEY): $(SERVER_DIR)
	@ $(OPENSSL) genrsa -passout "pass:$(SERVER_PASSPARSE)" -des3 -out "$(SERVER_KEY)" $(SERVER_NUMBITS)

$(SERVER_CSR): $(SERVER_DIR) $(SERVER_KEY)
	@ $(OPENSSL) req -passin "pass:$(SERVER_PASSPARSE)" -new -key "$(SERVER_KEY)" -out "$(SERVER_CSR)" -subj "$(SERVER_SUBJ)"

$(SERVER_CRT): ca_certs $(SERVER_DIR) $(SERVER_CSR)
	@ $(OPENSSL) x509 -req -passin "pass:$(SERVER_PASSPARSE)" -days $(SERVER_DAYS) -in "$(SERVER_CSR)" -CA "$(CA_CRT)" -CAkey "$(CA_KEY)" -set_serial 01 -out $(SERVER_CRT)
	@ $(OPENSSL) rsa -passin "pass:$(SERVER_PASSPARSE)" -in "$(SERVER_KEY)" -out "$(SERVER_KEY)"
	@ cp $(CA_CRT) $(SERVER_CA_CRT)


client_certs: $(CLIENT_CRT)

$(CLIENT_DIR):
	@ mkdir -p $(CLIENT_DIR)

$(CLIENT_KEY): $(CLIENT_DIR)
	@ $(OPENSSL) genrsa -passout "pass:$(CLIENT_PASSPARSE)" -des3 -out "$(CLIENT_KEY)" $(CLIENT_NUMBITS)

$(CLIENT_CSR): $(CLIENT_DIR) $(CLIENT_KEY)
	@ $(OPENSSL) req -passin "pass:$(CLIENT_PASSPARSE)" -new -key "$(CLIENT_KEY)" -out "$(CLIENT_CSR)" -subj "$(CLIENT_SUBJ)"

$(CLIENT_CRT): ca_certs $(CLIENT_DIR) $(CLIENT_CSR)
	@ $(OPENSSL) x509 -passin "pass:$(CLIENT_PASSPARSE)" -req -days $(CLIENT_DAYS) -in "$(CLIENT_CSR)" -CA "$(CA_CRT)" -CAkey "$(CA_KEY)" -set_serial 01 -out "$(CLIENT_CRT)"
	@ $(OPENSSL) rsa -passin "pass:$(CLIENT_PASSPARSE)" -in "$(CLIENT_KEY)" -out "$(CLIENT_KEY)"
	@ cp $(CA_CRT) $(CLIENT_CA_CRT)


clean:
	rm -rf $(CA_DIR) $(SERVER_DIR) $(CLIENT_DIR)
