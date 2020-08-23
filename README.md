<div align="center">
  <br>
  <h1>cfaccessproxy</h1>
</div>

**cfaccessproxy** is a [Cloudflare Access](https://developers.cloudflare.com/access/) companion proxy.

## Installation

1. Install the latest version of [Go](https://golang.org) if you haven't yet.

2. Then run:

        $ go get github.com/astrophena/cfaccessproxy

   **cfaccessproxy** should be installed in your `$GOBIN`.

3. Set some environment variables:

| Name | Description |
| ---- | ----------- |
| CFACCESSPROXY_LISTEN_ADDR | Network address to listen on. |
| CFACCESSPROXY_CANONICAL_URL | Canonical URL to redirect. |
| CFACCESSPROXY_UPSTREAM | URL to proxy requests after JWT check. |
| CFACCESSPROXY_AUTH_DOMAIN | Cloudflare Access domain (e. g. \*.cloudflareaccess.com). |
| CFACCESSPROXY_POLICY_AUD | Application AUD from Cloudflare Access. |
| CFACCESSPROXY_BYPASS_URL_PREFIXES | Comma-separated list of prefixes that should bypass JWT check. |

4. Start **cfaccessproxy**:

        $ cfaccessproxy

## License

[MIT](LICENSE.md) © [Ilya Mateyko](https://github.com/astrophena)
