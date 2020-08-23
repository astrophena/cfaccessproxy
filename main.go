/*
       __
  ___ / _| __ _  ___ ___ ___  ___ ___ _ __  _ __ _____  ___   _
 / __| |_ / _` |/ __/ __/ _ \/ __/ __| '_ \| '__/ _ \ \/ / | | |
| (__|  _| (_| | (_| (_|  __/\__ \__ \ |_) | | | (_) >  <| |_| |
 \___|_|  \__,_|\___\___\___||___/___/ .__/|_|  \___/_/\_\\__, |
                                     |_|                  |___/
*/

// cfaccessproxy is a Cloudflare Access companion proxy.
// See README.md for documentation.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
)

const (
	jwtHeader = "Cf-Access-Jwt-Assertion"

	noTokenErr      = "No token on the request."
	invalidTokenErr = "Invalid token."
)

// Config represents a cfaccessproxy configuration.
type Config struct {
	ListenAddr        string   `default:":3000" split_words:"true"`
	CanonicalURL      string   `required:"true" split_words:"true"`
	Upstream          string   `required:"true" split_words:"true"`
	AuthDomain        string   `required:"true" split_words:"true"`
	PolicyAUD         string   `required:"true" split_words:"true"`
	BypassURLPrefixes []string `split_words:"true"`
}

// CertsURL returns the URL of the endpoint which is used for JWT
// verification.
func (c *Config) CertsURL() string {
	return fmt.Sprintf("%s/cdn-cgi/access/certs", c.AuthDomain)
}

// ReverseProxy implements a reverse proxy.
type ReverseProxy struct {
	c *Config
}

// ServeHTTP makes ReverseProxy satisfy the http.Handler interface.
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(rp.c.Upstream)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	pr := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Host = url.Host

	pr.ServeHTTP(w, r)
}

// JWTVerifier is a HTTP middleware that verifies JWTs issued by
// Cloudflare Access.
func JWTVerifier(c *Config, next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		verifier := oidc.NewVerifier(
			c.AuthDomain,
			oidc.NewRemoteKeySet(context.TODO(), c.CertsURL()),
			&oidc.Config{ClientID: c.PolicyAUD},
		)

		if len(c.BypassURLPrefixes) > 0 {
			for _, p := range c.BypassURLPrefixes {
				if strings.HasPrefix(r.URL.Path, p) {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		jwt := r.Header.Get(jwtHeader)
		if jwt == "" {
			http.Error(w, noTokenErr, http.StatusUnauthorized)
			return
		}

		tok, err := verifier.Verify(r.Context(), jwt)
		if err != nil {
			http.Error(w, invalidTokenErr, http.StatusUnauthorized)
			return
		}

		if r.URL.Path == "/_/token_info" {
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			if err := enc.Encode(tok); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func main() {
	log.SetFlags(0)

	var c Config
	if err := envconfig.Process("cfaccessproxy", &c); err != nil {
		log.Fatal(err)
	}

	var (
		canonicalHost = handlers.CanonicalHost(c.CanonicalURL, http.StatusMovedPermanently)
		reverseProxy  = &ReverseProxy{c: &c}

		srv = &http.Server{
			Addr:         c.ListenAddr,
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 15,
			Handler:      canonicalHost(JWTVerifier(&c, reverseProxy)),
		}
	)

	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
