# cloudflare-go

A golang library for interacting with Cloudflare API.

## Usage

```go
import cfgo

const (
    email = "my-cloudflare-account@email.com"
    zoneID = "my-cloudflare-zone-id"
    apiKey = "my-cloudflare-api-key"
)

func main() {
    client := cfgo.NewCloudflareClient(email, apiKey)

    // do something with `client`
}
```

See more usages [here](https://github.com/meinside/cloudflare-go/tree/master/cmd/cf-dns-cli).

## Implementations

- [X] List/create/update/delete DNS records
- [ ] Other things that I need
- [ ] All others

## License

MIT

