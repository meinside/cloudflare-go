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

See the following (not so helpful) message with `cf-dns-cli -h` or `cf-dns-cli --help`.

```
Usage v0.0.1:

<Flags>

  -h / --help: Show this help message.

  -v / --verbose: Show verbose messages for debugging purpose.


<Commands and parameters>

List all zones for this account.

  $ cf-dns-cli zones

List all DNS records for given zone identifier.

  $ cf-dns-cli records [ZONE_ID]

Create a DNS record with given parameters.

  $ cf-dns-cli create [ZONE_ID] [RECORD_TYPE] [key1=value1 key2=value2 ...]

  e.g.: $ cf-dns-cli create abcd123456 CNAME name=sub.from.com content=dest.com comment="New record."

Update a DNS record with given parameters.

  $ cf-dns-cli update [ZONE_ID] [RECORD_ID] [key1=value1 key2=value2 ...]

  e.g.: $ cf-dns-cli update abcd123456 wxyz098765 type=CNAME name=sub.from.com content=updated-dest.com comment="Updated record."

Batch upsert all DNS records in the given JSON file.

  $ cf-dns-cli batch [RECORDS_FILEPATH]

  If a record has 'id' in it, it will be updated. Otherwise, it will be newly created instead.

Delete a DNS record with given zone & record identifier.

  $ cf-dns-cli delete [ZONE_ID] [RECORD_ID]

Generate a sample DNS records file in JSON format. (file used with 'batch' command)

  $ cf-dns-cli generate
```


## known issues

- [ ] (create/update) nested parameters are not supported yet

