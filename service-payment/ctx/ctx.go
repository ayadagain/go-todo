package ctx

type ServiceCtx interface{}

type defaultServiceCtx struct {
	config *Config
}

func (ctx *defaultServiceCtx) Config() *Config {
	return ctx.config
}
