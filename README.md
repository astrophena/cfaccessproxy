<div align="center">
  <h1>cfaccessproxy</h1>
</div>

**cfaccessproxy** is a [Cloudflare Access](https://developers.cloudflare.com/access/) companion proxy.

## Installation

1. Install the latest version of [Go](https://golang.org) if you haven't yet.

2. Then run:

        $ go get github.com/astrophena/cfaccessproxy

   **cfaccessproxy** should be installed in your `$GOBIN` (e. g.
   `~/go/bin`).

3. Set some environment variables:

| Name | Description |
| ---- | ----------- |
| `CFACCESSPROXY_LISTEN_ADDR` | Network address to listen on (**optional**, `:3000` by default). |
| `CFACCESSPROXY_CANONICAL_URL` | Canonical URL to redirect (**required**). |
| `CFACCESSPROXY_UPSTREAM` | URL to proxy requests after JWT check (**required**). |
| `CFACCESSPROXY_AUTH_DOMAIN` | Cloudflare Access domain (e. g. \*.cloudflareaccess.com) (**required**). |
| `CFACCESSPROXY_POLICY_AUD` | Application AUD from Cloudflare Access (**required**). |
| `CFACCESSPROXY_BYPASS_URL_PREFIXES` | Comma-separated list of URL prefixes that should bypass JWT check (**required**). |

4. Start **cfaccessproxy**:

        $ cfaccessproxy

   You will probably want to setup a [systemd](https://systemd.io)
   service or init script to autostart **cfaccessproxy**.

## License

[MIT](LICENSE.md) © [Ilya Mateyko](https://github.com/astrophena)
