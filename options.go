package deepl

type deepLClientOption struct {
	httpProxy            string
	socket5Proxy         string
	socket5ProxyUser     string
	socket5proxyPassword string
}

type Option func(option *deepLClientOption)

// WithHttpProxy  set http proxy.if both httpProxy and sock5 proxy are set,
// http proxy will be over-wrote by sock5 proxy .example http://127.0.0.1:1080
func WithHttpProxy(proxy string) Option {
	return func(option *deepLClientOption) {
		option.httpProxy = proxy
	}
}

// WithSocket5Proxy  set socket5Proxy.if both httpProxy and sock5 proxy are set,
// http proxy will be over-wrote by sock5 proxy example 127.0.0.1:1080
func WithSocket5Proxy(socket5Proxy string, userName string, password string) Option {
	return func(option *deepLClientOption) {
		option.socket5Proxy = socket5Proxy
		option.socket5ProxyUser = userName
		option.socket5proxyPassword = password
	}
}
