// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Command cfaccessproxy is a reverse proxy that authenticates
// Cloudflare Access requests.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
)

// Version is a version of cfaccessproxy.
var Version = "HEAD"

func main() {
	log.SetFlags(0)

	showVersion := flag.Bool("version", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Fprintf(os.Stderr, "%s\n", Version)
		os.Exit(0)
	}

	p, err := newProxy()
	if err != nil {
		log.Fatal(err)
	}

	h := handlers.CanonicalHost(p.BaseURL, http.StatusMovedPermanently)
	srv := &http.Server{
		Addr:         p.Addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      h(p),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("shutting down")
	srv.Shutdown(ctx)
}

func newProxy() (*proxy, error) {
	var p proxy

	if err := envconfig.Process("cfaccessproxy", &p); err != nil {
		return nil, err
	}

	url, err := url.Parse(p.Upstream)
	if err != nil {
		return nil, err
	}

	p.proxy = httputil.NewSingleHostReverseProxy(url)

	p.verifier = oidc.NewVerifier(
		p.AuthDomain,
		oidc.NewRemoteKeySet(context.Background(), p.certsURL()),
		&oidc.Config{ClientID: p.PolicyAUD},
	)

	return &p, nil
}

type proxy struct {
	Addr           string   `default:":3000" split_words:"true"`
	BaseURL        string   `required:"true" split_words:"true"`
	Upstream       string   `required:"true"`
	AuthDomain     string   `required:"true" split_words:"true"`
	PolicyAUD      string   `required:"true" split_words:"true"`
	BypassPrefixes []string `split_words:"true"`
	LogoutRedirect bool     `split_words:"true"`

	proxy    *httputil.ReverseProxy
	verifier *oidc.IDTokenVerifier
}

func (p *proxy) certsURL() string {
	return fmt.Sprintf("%s/cdn-cgi/access/certs", p.AuthDomain)
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(p.BypassPrefixes) > 0 {
		for _, pref := range p.BypassPrefixes {
			if strings.HasPrefix(r.URL.Path, pref) {
				p.proxy.ServeHTTP(w, r)
				return
			}
		}
	}

	jwt := r.Header.Get("Cf-Access-Jwt-Assertion")
	if jwt == "" {
		http.Error(w, "No token on the request.", http.StatusUnauthorized)
		return
	}

	_, err := p.verifier.Verify(r.Context(), jwt)
	if err != nil {
		http.Error(w, "Invalid token.", http.StatusUnauthorized)
		return
	}

	if p.LogoutRedirect && r.URL.Path == "/logout" {
		http.Redirect(w, r, p.AuthDomain+"/cdn-cgi/access/logout", http.StatusFound)
		return
	}

	p.proxy.ServeHTTP(w, r)
}
