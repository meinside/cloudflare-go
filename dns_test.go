package cfgo

import (
	"log"
	"os"
	"testing"
)

func TestCRUDRequests(t *testing.T) {
	email := os.Getenv("EMAIL")
	apiKey := os.Getenv("API_KEY")
	zoneID := os.Getenv("ZONE_ID")
	verbose := false

	if os.Getenv("VERBOSE") == "true" {
		verbose = true
	}

	if email != "" && apiKey != "" && zoneID != "" {
		client := NewCloudflareClient(email, apiKey)
		client.Verbose = verbose

		// list zones
		if zones, err := client.ListZones(); err == nil {
			if len(zones.Result) <= 0 {
				t.Errorf("no zones found")
			}
		} else {
			t.Errorf("failed to list zones: %s", err)
		}

		// create a record
		cname := NewDNSRecordCNAME("testing", "test.somewhere.com")
		cname.SetComment("CNAME record for testing the library.")
		cname.SetTTL(3600)

		if created, err := client.CreateDNSRecord(zoneID, cname); err == nil {
			var createdCNAME DNSRecordCNAME
			if err := created.Result.Into(&createdCNAME); err != nil {
				t.Errorf("failed to parse created result: %s", err)
			} else {
				if verbose {
					log.Printf("created dns record = %+v", created)
				}

				// list records
				if retrieved, err := client.ListDNSRecords(zoneID, map[string]any{
					"comment.contains": "testing",
					"type":             createdCNAME.Type,
				}); err == nil {
					if verbose {
						log.Printf("retrieved dns records = %+v", retrieved)
					}

					exists := false
					var record DNSRecordCNAME
					for _, rawRecord := range retrieved.Result {
						if rawRecord.GetType() == createdCNAME.Type {
							if err := rawRecord.Into(&record); err == nil {
								// check if the created record exists
								if createdCNAME.ID == record.ID {
									exists = true

									if verbose {
										log.Printf("matched dns record = %+v", record)
									}

									break
								}
							} else {
								t.Errorf("failed to parse raw record: %s", err)
							}
						}
					}

					if !exists {
						t.Errorf("there was no newly-created dns record in the retrieved dns records")
					}
				} else {
					t.Errorf("failed to list records: %s", err)
				}

				// update a record
				cname.SetComment("Updated CNAME record for testing the library.")
				if updated, err := client.UpdateDNSRecord(zoneID, createdCNAME.ID, cname); err == nil {
					var updatedCNAME DNSRecordCNAME
					if err := updated.Result.Into(&updatedCNAME); err != nil {
						t.Errorf("failed to parse updated result: %s", err)
					} else {
						if verbose {
							log.Printf("updated dns record = %+v", updatedCNAME)
						}
					}
				} else {
					t.Errorf("failed to update dns record: %s", err)
				}

				// delete the record
				if deleted, err := client.DeleteDNSRecord(zoneID, createdCNAME.ID); err != nil {
					t.Errorf("failed to delete dns record: %s", err)
				} else {
					if verbose {
						log.Printf("deleted dns record = %+v", deleted)
					}
				}
			}
		} else {
			t.Errorf("failed to create dns record: %s", err)
		}
	} else {
		t.Errorf("Export `EMAIL`, API_KEY`, and `ZONE_ID` before running tests.")
	}

}
