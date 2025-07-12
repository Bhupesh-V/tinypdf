package shared

type Config struct {
	OriginalFilePath string
	OriginalFileName string
	GSConfig         any
	IsDebug          bool
	// add more vendor-specific configs
}

type PipeFunc func(string, Config) string

func Pipe(funcs ...PipeFunc) PipeFunc {
	return func(input string, config Config) string {
		for _, fn := range funcs {
			input = fn(input, config)
		}
		return input
	}
}
