package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	cfgo "github.com/meinside/cloudflare-go"
	"github.com/meinside/version-go"
)

var _stdout = log.New(os.Stdout, "", 0)
var _stderr = log.New(os.Stderr, "", 0)

const (
	applicationName = "cf-dns-cli"
	configFilename  = "config.json"

	cmdZones    = "zones"
	cmdRecords  = "records"
	cmdCreate   = "create"
	cmdUpdate   = "update"
	cmdBatch    = "batch"
	cmdDelete   = "delete"
	cmdGenerate = "generate"

	regexKeyValue = `(.*?)=['"]?(.*?)['"]?$`
	regexFloat    = `[+-]?([0-9]*[.])?[0-9]+`
	regexInt      = `[+-]?[0-9]+`
)

// config struct for configuration
type config struct {
	Email  string `json:"email"`
	APIKey string `json:"api_key"`
}

// read config file
func readConfig() (conf config, err error) {
	configFilepath := strings.Join([]string{getConfigDir(), configFilename}, string(filepath.Separator))

	var bytes []byte
	if bytes, err = os.ReadFile(configFilepath); err == nil {
		if err = json.Unmarshal(bytes, &conf); err == nil {
			if conf.APIKey == "" || conf.Email == "" {
				return config{}, fmt.Errorf("missing properties `email` or `api_key` in: %s", configFilepath)
			}

			return conf, nil
		}
	}

	return config{}, err
}

// get "$XDG_CONFIG_HOME/cf-dns-updater"
func getConfigDir() string {
	// https://xdgbasedirectoryspecification.com
	configDir := os.Getenv("XDG_CONFIG_HOME")

	// If the value of the environment variable is unset, empty, or not an absolute path, use the default
	if configDir == "" || configDir[0:1] != "/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			_stderr.Fatalf("failed to get home directory (%s)\n", err)
		} else {
			configDir = filepath.Join(homeDir, ".config", applicationName)
		}
	} else {
		configDir = filepath.Join(configDir, applicationName)
	}

	return configDir
}

// check if short/long flag arguments exist in the args
func flagExists(args []string, short, long string) bool {
	for _, arg := range args {
		if arg == short || arg == long {
			return true
		}
	}

	return false
}

// encode json string for debugging
func jsonString(v any) string {
	if bytes, err := json.Marshal(v); err == nil {
		return string(bytes)
	} else {
		return fmt.Sprintf("<%s>", err)
	}
}

// show help message
func showHelp(applicationName string, err error) {
	if err != nil {
		_stdout.Printf("Error: %s\n\n", err)
	}

	_stdout.Printf(`Usage %[2]s:

<Flags>

  -h / --help: Show this help message.

  -v / --verbose: Show verbose messages for debugging purpose.


<Commands and parameters>

List all zones for this account.

  $ %[1]s %[3]s

List all DNS records for given zone identifier.

  $ %[1]s %[4]s [ZONE_ID]

Create a DNS record with given parameters.

  $ %[1]s %[5]s [ZONE_ID] [RECORD_TYPE] [key1=value1 key2=value2 ...]

  e.g.: $ %[1]s %[5]s abcd123456 CNAME name=sub.from.com content=dest.com comment="New record."

Update a DNS record with given parameters.

  $ %[1]s %[6]s [ZONE_ID] [RECORD_ID] [key1=value1 key2=value2 ...]

  e.g.: $ %[1]s %[6]s abcd123456 wxyz098765 type=CNAME name=sub.from.com content=updated-dest.com comment="Updated record."

Batch upsert all DNS records in the given JSON file.

  $ %[1]s %[7]s [RECORDS_FILEPATH]

  If a record has 'id' in it, it will be updated. Otherwise, it will be newly created instead.

Delete a DNS record with given zone & record identifier.

  $ %[1]s %[8]s [ZONE_ID] [RECORD_ID]

Generate a sample DNS records file in JSON format. (file used with '%[7]s' command)

  $ %[1]s %[9]s
`, applicationName, version.Minimum(),
		cmdZones, cmdRecords, cmdCreate, cmdUpdate, cmdBatch, cmdDelete, cmdGenerate)

	if err == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// print sample records (JSON) to stdout
func showSampleRecords() {
	const sampleZoneID = "--your-zone-id-here--"
	const sampleRecordID = "--id-of-an-already-existing-record-here--"

	records := []any{
		cfgo.NewDNSRecordA("sample1.com", "1.2.3.4").
			SetComment("A record for sampling").
			SetZoneID(sampleZoneID).
			SetID(sampleRecordID),
		cfgo.NewDNSRecordAAAA("sample2.com", "beef:beef:beef:beef:beef:beef:beef:beef").
			SetComment("AAAA record for sampling").
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordCAA("sample3.com", 0, "issue", "letsencrypt.org").
			SetComment("CAA record for sampling").
			SetID(sampleRecordID).
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordCERT("sample4.com", 8, "--certificate--", 1, 9).
			SetComment("CERT record for sampling").
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordCNAME("sample.sample5.com", "somewhere.com").
			SetComment("CNAME record for sampling").
			SetID(sampleRecordID).
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordDNSKEY("sample6.com", 5, 1, 3, "--public-key--").
			SetComment("DNSKEY record for sampling").
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordDS("sample7.com", 3, "--digest--", 1, 1).
			SetComment("DS record for sampling").
			SetID(sampleRecordID).
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordHTTPS("sample8.com", 1, ".", `alpn="h3,h2" ipv4hint="127.0.0.1" ipv6hint="::1"`).
			SetComment("HTTPS record for sampling").
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordLOC("sample9.com", 0, 37, cfgo.North, 46, 46, 122, cfgo.West, 23, 35, 0, 0, 100).
			SetComment("LOC record for sampling").
			SetID(sampleRecordID).
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordMX("sample10.com", "mx.sample10.com", 10).
			SetComment("MX record for sampling").
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordNAPTR("sample11.com", "flags", 100, 10, "regex", "replacement", "service").
			SetComment("NAPTR record for sampling").
			SetID(sampleRecordID).
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordNS("sample12.com", "ns1.sample12").
			SetComment("NS record for sampling").
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordPTR("sample13.com", "ptr.sample13.com").
			SetComment("PTR record for sampling").
			SetID(sampleRecordID).
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordSMIMEA("sample14.com", "--cretificate--", 0, 0, 3).
			SetComment("SMIMEA record for sampling").
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordSRV("sample15.com", 8806, 10, "_tcp", "_sip", "sample15.com", 5).
			SetComment("SRV record for sampling").
			SetID(sampleRecordID).
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordSSHFP("sample16.com", 2, "--fingerprint--", 1).
			SetComment("SSHFP record for sampling").
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordSVCB("sample17.com", 1, ".", `alpn="h3,h2" ipv4hint="127.0.0.1" ipv6hint="::1"`).
			SetComment("SVCB record for sampling").
			SetID(sampleRecordID).
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordTLSA("sample18.com", "--certificate--", 1, 0, 0).
			SetComment("TLSA record for sampling").
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordTXT("sample19.com", "sample text content").
			SetComment("TXT record for sampling").
			SetID(sampleRecordID).
			SetZoneID(sampleZoneID),
		cfgo.NewDNSRecordURI("sample20.com", "http://sample20.com/sample.html", 20).
			SetComment("URI record for sampling").
			SetZoneID(sampleZoneID),
	}

	if bytes, err := json.MarshalIndent(records, "", "  "); err == nil {
		_stdout.Printf("%s\n", string(bytes))

		os.Exit(0)
	} else {
		_stderr.Printf("failed to print sample records: %s\n", err)

		os.Exit(1)
	}
}

// list all zones
func listZones(client *cfgo.CloudflareClient) {
	if zones, err := client.ListZones(); err == nil {
		for _, zone := range zones.Result {
			_stdout.Printf("%s %s\n", zone.ID, zone.Name)
		}

		os.Exit(0)
	} else {
		_stderr.Printf("failed to list zones: %s\n", err)

		os.Exit(1)
	}
}

// list all DNS records for given zone identifier
func listDNSRecords(client *cfgo.CloudflareClient, zoneID string) {
	if records, err := client.ListDNSRecords(zoneID, nil); err == nil {
		for _, record := range records.Result {
			if name, err := record.StringFor("name"); err == nil {
				if id, err := record.StringFor("id"); err == nil {
					if typ3, err := record.StringFor("type"); err == nil {
						lines := []string{name}
						if content, err := record.StringFor("content"); err == nil {
							lines = append(lines, content)
						}
						if comment, err := record.StringFor("comment"); err == nil {
							lines = append(lines, comment)
						}

						_stdout.Printf("%s [%s] %s", id, typ3, strings.Join(lines, " | "))
					}
				}
			}
		}

		os.Exit(0)
	} else {
		_stderr.Printf("failed to list DNS records for zone %s: %s", zoneID, err)

		os.Exit(1)
	}
}

// create a DNS record with given parameters
func createDNSRecord(client *cfgo.CloudflareClient, zoneID, typ3 string, params map[string]any) {
	record := map[string]any{
		"type": typ3,
	}
	for k, v := range params {
		record[k] = v
	}

	// create
	if _, err := client.CreateDNSRecord(zoneID, record); err == nil {
		_stdout.Printf("created [%s] record with params %s\n", typ3, jsonString(record))
		os.Exit(0)
	} else {
		_stderr.Printf("failed to create [%s] record with params %s: %s\n", typ3, jsonString(record), err)
		os.Exit(1)
	}
}

// update a DNS record with given parameters
func updateDNSRecord(client *cfgo.CloudflareClient, zoneID, recordID string, params map[string]any) {
	record := map[string]any{
		"id": recordID,
	}
	for k, v := range params {
		record[k] = v
	}

	// update
	if updated, err := client.UpdateDNSRecord(zoneID, recordID, record); err == nil {
		if typ3, err := updated.Result.StringFor("type"); err == nil {
			_stdout.Printf("updated [%s] record with params %s\n", typ3, jsonString(record))
		}
		os.Exit(0)
	} else {
		_stderr.Printf("failed to update record with params %s: %s\n", jsonString(record), err)
		os.Exit(1)
	}
}

// upsert all DNS records with given JSON file
func upsertDNSRecords(client *cfgo.CloudflareClient, fpath string) {
	processed := 0
	failed := 0

	var records []cfgo.DNSRecordRaw
	if bytes, err := os.ReadFile(fpath); err == nil {
		if err := json.Unmarshal(bytes, &records); err == nil {
			for _, record := range records {
				var err error
				var zoneID, recordID string
				if zoneID, err = record.StringFor("zone_id"); err != nil {
					failed += 1

					_stderr.Printf("zone id not found in record: %s", err)
				} else {
					recordID, _ = record.StringFor("id") // record id can be null (when creating a new one)

					// upsert
					if recordID != "" {
						// update
						if _, err = client.UpdateDNSRecord(zoneID, recordID, record); err == nil {
							processed += 1

							if recordName, err := record.StringFor("name"); err == nil {
								_stdout.Printf("updated [%s] record '%s'\n", record.GetType(), recordName)
							}
						} else {
							failed += 1

							if recordName, e := record.StringFor("name"); e == nil {
								_stderr.Printf("failed to update [%s] record '%s': %s\n", record.GetType(), recordName, err)
							}
						}
					} else {
						// create
						if _, err = client.CreateDNSRecord(zoneID, record); err == nil {
							processed += 1

							if recordName, err := record.StringFor("name"); err == nil {
								_stdout.Printf("created [%s] record '%s'\n", record.GetType(), recordName)
							}
						} else {
							failed += 1

							if recordName, e := record.StringFor("name"); e == nil {
								_stderr.Printf("failed to create [%s] record '%s': %s\n", record.GetType(), recordName, err)
							}
						}
					}
				}
			}

			_stderr.Printf("processed %d DNS records (%d errors)\n", processed, failed)

			if failed == 0 {
				os.Exit(0)
			}
		} else {
			_stderr.Printf("failed to parse JSON file: %s\n", err)
		}
	} else {
		_stderr.Printf("failed to read file: %s\n", err)
	}

	os.Exit(1)
}

// delete a DNS record with given record identifier
func deleteDNSRecord(client *cfgo.CloudflareClient, zoneID, recordID string) {
	if deleted, err := client.DeleteDNSRecord(zoneID, recordID); err == nil {
		_stdout.Printf("successfully deleted DNS record %s", deleted.Result.ID)

		os.Exit(0)
	} else {
		_stderr.Printf("failed to delete DNS record %s for zone %s: %s", recordID, zoneID, err)

		os.Exit(1)
	}
}

// filter parameters only (drop flags)
func filterParams(args []string) (filtered []string) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		filtered = append(filtered, arg)
	}

	return filtered
}

// convert an array of strings like "key1=value1" into a map with desired value types
func convertKeyValueParams(params []string) (result map[string]any) {
	result = map[string]any{}

	regexKV := regexp.MustCompile(regexKeyValue)
	regexFloat := regexp.MustCompile(regexFloat)
	regexInt := regexp.MustCompile(regexInt)

	for _, param := range params {
		matches := regexKV.FindStringSubmatch(param)

		if len(matches) == 3 {
			k := matches[1]
			v := matches[2]

			// FIXME: not all types are handled properly
			if regexFloat.MatchString(v) { // float
				result[k], _ = strconv.ParseFloat(v, 32)
				continue
			} else if regexInt.MatchString(v) { // int
				result[k], _ = strconv.ParseInt(v, 10, 32)
				continue
			} else if arr := strings.Split(v, ","); len(arr) > 1 { // []string
				result[k] = arr
				continue
			} else if v == "true" { // bool (true)
				result[k] = true
				continue
			} else if v == "false" { // bool (false)
				result[k] = false
				continue
			}

			result[k] = v
		}
	}

	return result
}

// run with arguments
func run(application string, args []string) {
	// handle flags
	if flagExists(args, "-h", "--help") {
		showHelp(application, nil)
	}
	verbose := flagExists(args, "-v", "--verbose")

	if conf, err := readConfig(); err == nil {
		client := cfgo.NewCloudflareClient(conf.Email, conf.APIKey)
		client.Verbose = verbose

		// handle commands
		argsWithoutFlags := filterParams(args)
		if len(argsWithoutFlags) > 0 {
			cmd := argsWithoutFlags[0]
			params := argsWithoutFlags[1:]
			switch cmd {
			case cmdZones: // list zones
				listZones(client)
			case cmdRecords:
				if len(params) >= 1 {
					listDNSRecords(client, params[0])
				} else {
					showHelp(application, fmt.Errorf("zone identifier was not given"))
				}
			case cmdCreate:
				if len(params) >= 3 {
					kvs := convertKeyValueParams(params[2:])
					if len(kvs) > 0 {
						createDNSRecord(client, params[0], params[1], kvs)
					} else {
						showHelp(application, fmt.Errorf("parameters for a new DNS record were not given"))
					}
				} else {
					showHelp(application, fmt.Errorf("essential parameters were not given"))
				}
			case cmdUpdate:
				if len(params) >= 3 {
					kvs := convertKeyValueParams(params[2:])
					if len(kvs) > 0 {
						updateDNSRecord(client, params[0], params[1], kvs)
					} else {
						showHelp(application, fmt.Errorf("parameters for an updated DNS record were not given"))
					}
				} else {
					showHelp(application, fmt.Errorf("essential parameters were not given"))
				}
			case cmdBatch:
				if len(params) >= 1 {
					upsertDNSRecords(client, params[0])
				} else {
					showHelp(application, fmt.Errorf("JSON filepath was not given"))
				}
			case cmdDelete:
				if len(params) >= 2 {
					deleteDNSRecord(client, params[0], params[1])
				} else {
					showHelp(application, fmt.Errorf("zone identifier or DNS record identifier was not given"))
				}
			case cmdGenerate:
				showSampleRecords()
			}

			showHelp(application, fmt.Errorf("'%s' is not a supported command.", cmd))
		} else {
			showHelp(application, nil)
		}
	} else {
		_stderr.Fatalf("failed to read config: %s\n", err)
	}
}
