package deepl

import "github.com/lib-x/deepl/internal/option"

type TagHandlingType int

const (
	TagHandlingHtml TagHandlingType = iota
	TagHandlingXml  TagHandlingType = iota
)

type Option func(option *option.DeepLClientOption)

// WithTagHandling TagHandling: type of tags to parse before translation, options are "html" and "xml"
// Todo: support xml options，see https://www.nuget.org/packages/DeepL.net/
func WithTagHandling(handingType TagHandlingType) Option {
	tagHandlingType := ""
	switch handingType {
	case TagHandlingHtml:
		tagHandlingType = "html"
	case TagHandlingXml:
		tagHandlingType = "xml"
	}
	return func(option *option.DeepLClientOption) {
		if tagHandlingType != "" {
			option.TagHandling = tagHandlingType
		}
	}
}

// WithDeeplProSession set deepl pro session.if session is set,
func WithDeeplProSession(dlSession string) Option {
	return func(option *option.DeepLClientOption) {
		if len(dlSession) == 36 {
			option.UseDeepLPro = true
			option.DlSession = dlSession
		}
	}
}

// WithHttpProxy  set http proxy.if both httpProxy and sock5 proxy are set,
// http proxy will be over-wrote by sock5 proxy .example http://127.0.0.1:1080
func WithHttpProxy(proxy string) Option {
	return WithHttpProxyEx(proxy, false)
}

// WithHttpProxyEx  set http proxy.if both httpProxy and sock5 proxy are set,
// http proxy will be over-wrote by sock5 proxy .example http://http://127.0.0.1:1080
// ignoreSSLVerification: ignore SSL verification
func WithHttpProxyEx(proxy string, ignoreSSLVerification bool) Option {
	return func(option *option.DeepLClientOption) {
		option.HttpProxy = proxy
		option.IgnoreSSLVerification = ignoreSSLVerification
	}
}

// WithSocket5Proxy  set socket5Proxy.if both httpProxy and sock5 proxy are set,
// http proxy will be over-wrote by sock5 proxy example 127.0.0.1:1080
func WithSocket5Proxy(socket5Proxy string, userName string, password string) Option {
	return WithSocket5ProxyEx(socket5Proxy, userName, password, false)
}

// WithSocket5ProxyEx  set socket5Proxy.if both httpProxy and sock5 proxy are set,
// http proxy will be over-wrote by sock5 proxy example 127.0.0.1:1080
// ignoreSSLVerification: ignore SSL verification
func WithSocket5ProxyEx(socket5Proxy string, userName string, password string, ignoreSSLVerification bool) Option {
	return func(option *option.DeepLClientOption) {
		option.Socket5Proxy = socket5Proxy
		option.Socket5ProxyUser = userName
		option.Socket5proxyPassword = password
		option.IgnoreSSLVerification = ignoreSSLVerification
	}
}
