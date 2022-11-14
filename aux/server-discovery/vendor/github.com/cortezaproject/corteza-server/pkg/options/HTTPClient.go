package options

func (o *HTTPClientOpt) Defaults() {
	// just in case anyone used env var with the typo (before it was fixed)
	o.ClientTSLInsecure = EnvBool("HTTP_CLIENT_TSL_INSECURE", o.ClientTSLInsecure)

}
