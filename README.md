<div align="center">
  <h1>cfaccessproxy</h1>
</div>

`cfaccessproxy` is a simple reverse proxy that authenticates
[Cloudflare Access] requests.

## Installation

1. Install the latest version of [Go] if you haven't yet.

2. Install with `go get`:

        $ pushd $(mktemp -d); go mod init tmp; go get go.astrophena.name/cfaccessproxy; popd

   `go get` puts binaries by default to `$GOPATH/bin` (e.g.
   `~/go/bin`).

   Use `GOBIN` environment variable to change this behavior.

## Configuration

`cfaccessproxy` is configured by environment variables.

| Name | Description |
| ---- | ----------- |
| `CFACCESSPROXY_ADDR` | Address to listen on (**optional**, `:3000` by default). |
| `CFACCESSPROXY_BASE_URL` | Base URL (used for canonical redirection, *required*). |
| `CFACCESSPROXY_UPSTREAM` | URL to proxy requests (*required*). |
| `CFACCESSPROXY_AUTH_DOMAIN` | Cloudflare Access domain (e. g. \*.cloudflareaccess.com) (*required*). |
| `CFACCESSPROXY_POLICY_AUD` | Application AUD from Cloudflare Access (*required*). |
| `CFACCESSPROXY_BYPASS_PREFIXES` | Comma-separated list of URL prefixes that should bypass JWT check (*optional*). |

## License

[MIT] © [Ilya Mateyko]

[Cloudflare Access]: https://www.cloudflare.com/teams/access/
[Go]: https://golang.org
[MIT]: LICENSE.md
[Ilya Mateyko]: https://astrophena.name
