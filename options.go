package deepl

type deepLClientOption struct {
	httpProxy string
}

type Option func(option *deepLClientOption)

func WithHttpProxy(proxy string) Option {
	return func(option *deepLClientOption) {
		option.httpProxy = proxy
	}
}
