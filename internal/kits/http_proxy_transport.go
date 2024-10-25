package kits

import (
	"context"
	"crypto/tls"
	"github.com/lib-x/deepl/internal/option"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
)

func BuildHttpTransportWith(opt *option.DeepLClientOption) *http.Transport {
	var transport *http.Transport
	if opt.HttpProxy != "" {
		httpProxy, _ := url.Parse(opt.HttpProxy)
		if httpProxy != nil {
			transport = &http.Transport{Proxy: http.ProxyURL(httpProxy)}
		}
	}
	if opt.Socket5Proxy != "" {
		var auth *proxy.Auth
		if opt.Socket5ProxyUser != "" || opt.Socket5proxyPassword != "" {
			auth = &proxy.Auth{User: opt.Socket5ProxyUser, Password: opt.Socket5proxyPassword}
		}
		dialer, err := proxy.SOCKS5("tcp", opt.Socket5Proxy, auth, proxy.Direct)
		if err == nil {
			dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
				return dialer.Dial(network, address)
			}
			transport = &http.Transport{DialContext: dialContext}
		}
	}
	if opt.IgnoreSSLVerification && transport != nil {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return nil
}
