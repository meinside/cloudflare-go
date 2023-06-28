package cfgo

import (
	"encoding/json"
	"fmt"
)

// ListZones returns all zones.
func (c *CloudflareClient) ListZones() (response ResponseZones, err error) {
	var bytes []byte
	bytes, err = c.get(fmt.Sprintf("zones"), nil)

	if err == nil {
		err = json.Unmarshal(bytes, &response)
	}

	return response, err
}

// ListDNSRecords returns DNS records for given zone identifier and queries.
//
// The type of each `DNSRecordRaw` value in `Result` can be determined with `GetType()` function,
// then it can be converted into the determined struct type with `Into()` function.
//
// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-list-dns-records
func (c *CloudflareClient) ListDNSRecords(zoneID string, queries map[string]any) (response ResponseDNSRecords, err error) {
	var bytes []byte
	bytes, err = c.get(fmt.Sprintf("zones/%s/dns_records", zoneID), queries)

	if err == nil {
		err = json.Unmarshal(bytes, &response)
	}

	return response, err
}

// CreateDNSRecord creates a DNS record with given parameters.
//
// Generate a new record with NewDNSRecord* functions.
//
// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-create-dns-record
func (c *CloudflareClient) CreateDNSRecord(zoneID string, newOne any) (response ResponseDNSRecordCreation, err error) {
	var bytes []byte
	bytes, err = c.post(fmt.Sprintf("zones/%s/dns_records", zoneID), newOne)

	if err == nil {
		err = json.Unmarshal(bytes, &response)
	}

	return response, err
}

// DeleteDNSRecord deletes a DNS record with given identifiers.
//
// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-delete-dns-record
func (c *CloudflareClient) DeleteDNSRecord(zoneID, recordID string) (response ResponseDNSRecordDeletion, err error) {
	var bytes []byte
	bytes, err = c.delete(fmt.Sprintf("zones/%s/dns_records/%s", zoneID, recordID), nil)

	if err == nil {
		err = json.Unmarshal(bytes, &response)
	}

	return response, err
}

// UpdateDNSRecord updates a DNS record with given parameters.
//
// Updated record can be generated with NewDNSRecord* functions.
//
// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-update-dns-record
func (c *CloudflareClient) UpdateDNSRecord(zoneID, recordID string, updatedOne any) (response ResponseDNSRecordUpdate, err error) {
	var bytes []byte
	bytes, err = c.put(fmt.Sprintf("zones/%s/dns_records/%s", zoneID, recordID), updatedOne)

	if err == nil {
		err = json.Unmarshal(bytes, &response)
	}

	return response, err
}
