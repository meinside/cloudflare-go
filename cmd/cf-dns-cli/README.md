# cf-dns-cli

Cloudflare DNS CLI is a command line tool that allows you to manage your DNS records in the Cloudflare DNS zones.

## install

```bash
$ go install github.com/meinside/cloudflare-go/cmd/cf-dns-cli@latest
```

## configuration

Create a config file `$XDG_CONFIG_HOME/cf-dns-cli/config.json` with content:

```json
{
  "email": "your-cloudflare-account@email.com",
  "api_key": "your-cloudflare-global-api-key"
}
```

## usage

See the (not so helpful) help messages with:

```bash
$ cf-dns-cli -h
# or
$ cf-dns-cli --help
```

## known issues

- [ ] (create/update) nested parameters are not supported yet

