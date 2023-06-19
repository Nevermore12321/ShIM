package config

// customize config method and struct
type (
	options struct {
		env bool
	}

	Option func(opt *options)
)

// UseEnv use environment variables in config files
// return configuration func
func UseEnv() Option {
	return func(opt *options) {
		opt.env = true
	}
}
