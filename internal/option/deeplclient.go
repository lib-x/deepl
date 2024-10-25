package option

type DeepLClientOption struct {
	HttpProxy             string
	Socket5Proxy          string
	Socket5ProxyUser      string
	Socket5proxyPassword  string
	DlSession             string
	TagHandling           string
	UseDeepLPro           bool
	IgnoreSSLVerification bool
}
