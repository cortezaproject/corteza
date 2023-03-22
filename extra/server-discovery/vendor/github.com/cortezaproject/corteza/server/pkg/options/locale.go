package options

func (o *LocaleOpt) Defaults() {
	if Environment().IsDevelopment() || o.DevelopmentMode {
		o.Path = "../locale"
	}
}
