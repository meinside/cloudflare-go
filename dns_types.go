package cfgo

import (
	"encoding/json"
	"fmt"
)

// ResponseCommon struct for common response
type ResponseCommon struct {
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Messages []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"messages"`

	Success bool `json:"success"`
}

// Zone struct for all the zones
type Zone struct {
	Name    string `json:"name"`
	Account struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	CreatedOn       string `json:"created_on"`
	ActivatedOn     string `json:"activated_on"`
	ModifiedOn      string `json:"modified_on"`
	DevelopmentMode int    `json:"development_mode"`
	ID              string `json:"id"`
	Meta            struct {
		CustomCertificateQuoti  int  `json:"custom_certificate_quota"`
		MultipleRailgunsAllowed bool `json:"multiple_railguns_allowed"`
		PageRuleQuota           int  `json:"page_rule_quota"`
		PhishingDetected        bool `json:"phishing_detected"`
		Step                    int  `json:"step"`
	} `json:"meta"`
	NameServers         []string `json:"name_servers"`
	OriginalDNSHost     string   `json:"original_dnshost"`
	OriginalNameServers []string `json:"original_name_servers"`
	OriginalRegistrar   string   `json:"original_registrar"`
	Owner               struct {
		Email string `json:"email"`
		ID    string `json:"id"`
		Type  string `json:"type"`
	} `json:"owner"`
	Paused      bool     `json:"paused"`
	Permissions []string `json:"permissions"`
	Plan        struct {
		CanSubscribe      bool   `json:"can_subscribe"`
		Currency          string `json:"currency"`
		ExternallyManaged bool   `json:"externally_managed"`
		Frequency         string `json:"frequency"`
		ID                string `json:"id"`
		IsSubscribed      bool   `json:"is_subscribed"`
		LegacyDiscount    bool   `json:"legacy_discount"`
		LegacyID          string `json:"legacy_id"`
		Name              string `json:"name"`
		Price             int    `json:"price"`
	} `json:"plan"`
	Status string `json:"status"`
	Tenant struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tenant"`
	TenantUnit struct {
		ID string `json:"id"`
	} `json:"tenant_unit"`
	Type string `json:"type"`
}

// ResponseZones for zones response
type ResponseZones struct {
	ResponseCommon

	Result []Zone `json:"result"`
}

// DNSRecordType for the type of DNSRecords
type DNSRecordType string

const (
	A      DNSRecordType = "A"
	AAAA   DNSRecordType = "AAAA"
	CAA    DNSRecordType = "CAA"
	CERT   DNSRecordType = "CERT"
	CNAME  DNSRecordType = "CNAME"
	DNSKEY DNSRecordType = "DNSKEY"
	DS     DNSRecordType = "DS"
	HTTPS  DNSRecordType = "HTTPS"
	LOC    DNSRecordType = "LOC"
	MX     DNSRecordType = "MX"
	NAPTR  DNSRecordType = "NAPTR"
	NS     DNSRecordType = "NS"
	PTR    DNSRecordType = "PTR"
	SMIMEA DNSRecordType = "SMIMEA"
	SRV    DNSRecordType = "SRV"
	SSHFP  DNSRecordType = "SSHFP"
	SVCB   DNSRecordType = "SVCB"
	TLSA   DNSRecordType = "TLSA"
	TXT    DNSRecordType = "TXT"
	URI    DNSRecordType = "URI"

	Undefined DNSRecordType = ""
)

// DNSRecordRaw struct for DNS records in various types (A, CNAME, MX, ...)
type DNSRecordRaw map[string]any

// DNSRecordsRaw type for arrays of DNSRecordRaw structs
type DNSRecordsRaw []DNSRecordRaw

// generic function for returning a value for given key from the records
func valueFor[T any](r DNSRecordRaw, key string) (value T, err error) {
	var exists bool
	var v any
	if v, exists = r[key]; exists {
		var ok bool
		if value, ok = v.(T); ok {
			return value, nil
		} else {
			return value, fmt.Errorf("value for key: '%s' could not be converted to %T", value)
		}
	} else {
		return value, fmt.Errorf("no such value for key: '%s'", key)
	}
}

// StringFor returns the string value with given key string.
func (r DNSRecordRaw) StringFor(key string) (value string, err error) {
	return valueFor[string](r, key)
}

// StringsFor returns the string array with given key string.
func (r DNSRecordRaw) StringsFor(key string) (value []string, err error) {
	return valueFor[[]string](r, key)
}

// IntFor returns the int value with given key string.
func (r DNSRecordRaw) IntFor(key string) (value int, err error) {
	return valueFor[int](r, key)
}

// GetType returns the type of DNSRecord.
func (r DNSRecordRaw) GetType() DNSRecordType {
	if t, err := r.StringFor("type"); err == nil {
		return DNSRecordType(t)
	}

	return Undefined
}

// Into converts any DNSRecord into given type.
func (r DNSRecordRaw) Into(v any) (err error) {
	var encoded []byte
	if encoded, err = json.Marshal(r); err == nil {
		err = json.Unmarshal(encoded, v)
	}

	return err
}

// DNSRecordCommon struct for comman values of DNSRecord types
type DNSRecordCommon struct {
	Content   string        `json:"content,omitempty"`
	Name      string        `json:"name"`
	Type      DNSRecordType `json:"type"`
	Comment   string        `json:"comment,omitempty"`
	CreatedOn string        `json:"created_on,omitempty"`
	ID        string        `json:"id,omitempty"`
	Locked    bool          `json:"locked,omitempty"`
	Meta      struct {
		AutoAdded bool   `json:"auto_added,omitempty"`
		Source    string `json:"source,omitempty"`
	} `json:"meta,omitempty"`
	ModifiedOn string   `json:"modified_on,omitempty"`
	Proxiable  bool     `json:"proxiable,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	TTL        int      `json:"ttl,omitempty"`
	ZoneID     string   `json:"zone_id,omitempty"`
	ZoneName   string   `json:"zone_name,omitempty"`
}

// DNSRecordA struct for parsing A record
type DNSRecordA struct {
	DNSRecordCommon

	Proxied bool `json:"proxied,omitempty"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordA) SetID(recordID string) *DNSRecordA {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordA) SetComment(comment string) *DNSRecordA {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordA) SetTags(tags []string) *DNSRecordA {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordA) SetTTL(ttl int) *DNSRecordA {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordA) SetZoneID(zoneID string) *DNSRecordA {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordA creates a new A record.
func NewDNSRecordA(name, content string) *DNSRecordA {
	r := DNSRecordA{}
	r.Type = A
	r.Name = name
	r.Content = content

	return &r
}

// DNSRecordAAAA struct for parsing AAAA record
type DNSRecordAAAA DNSRecordA

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordAAAA) SetID(recordID string) *DNSRecordAAAA {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordAAAA) SetComment(comment string) *DNSRecordAAAA {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordAAAA) SetTags(tags []string) *DNSRecordAAAA {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordAAAA) SetTTL(ttl int) *DNSRecordAAAA {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordAAAA) SetZoneID(zoneID string) *DNSRecordAAAA {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordAAA creates a new AAAA record.
func NewDNSRecordAAAA(name, content string) *DNSRecordAAAA {
	r := DNSRecordAAAA{}
	r.Type = AAAA
	r.Name = name
	r.Content = content

	return &r
}

// DNSRecordCAA struct for parsing CAA record
type DNSRecordCAA struct {
	DNSRecordCommon

	Data struct {
		Flags int    `json:"flags"`
		Tag   string `json:"tag"`
		Value string `json:"value"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCAA) SetID(recordID string) *DNSRecordCAA {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCAA) SetComment(comment string) *DNSRecordCAA {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCAA) SetTags(tags []string) *DNSRecordCAA {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCAA) SetTTL(ttl int) *DNSRecordCAA {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCAA) SetZoneID(zoneID string) *DNSRecordCAA {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordCAA creates a new CAA record.
func NewDNSRecordCAA(name string, flags int, tag, value string) *DNSRecordCAA {
	r := DNSRecordCAA{}
	r.Type = CAA
	r.Name = name
	r.Data.Flags = flags
	r.Data.Tag = tag
	r.Data.Value = value

	return &r
}

// DNSRecordCERT struct for parsing CERT record
type DNSRecordCERT struct {
	DNSRecordCommon

	Data struct {
		Algorithm   int    `json:"algorithm"`
		Certificate string `json:"certificate"`
		KeyTag      int    `json:"key_tag"`
		Type        int    `json:"type"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCERT) SetID(recordID string) *DNSRecordCERT {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCERT) SetComment(comment string) *DNSRecordCERT {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCERT) SetTags(tags []string) *DNSRecordCERT {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCERT) SetTTL(ttl int) *DNSRecordCERT {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCERT) SetZoneID(zoneID string) *DNSRecordCERT {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordCERT creates a new CERT record.
func NewDNSRecordCERT(name string, algorithm int, certificate string, keyTag, typ3 int) *DNSRecordCERT {
	r := DNSRecordCERT{}
	r.Type = CERT
	r.Name = name
	r.Data.Algorithm = algorithm
	r.Data.Certificate = certificate
	r.Data.KeyTag = keyTag
	r.Data.Type = typ3

	return &r
}

// DNSRecordCNAME struct for parsing CNAME record
type DNSRecordCNAME struct {
	DNSRecordCommon

	Proxied bool `json:"proxied,omitempty"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCNAME) SetID(recordID string) *DNSRecordCNAME {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCNAME) SetComment(comment string) *DNSRecordCNAME {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCNAME) SetTags(tags []string) *DNSRecordCNAME {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCNAME) SetTTL(ttl int) *DNSRecordCNAME {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordCNAME) SetZoneID(zoneID string) *DNSRecordCNAME {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordCNAME creates a new CNAME record.
func NewDNSRecordCNAME(name, content string) *DNSRecordCNAME {
	r := DNSRecordCNAME{}
	r.Type = CNAME
	r.Name = name
	r.Content = content

	return &r
}

// DNSRecordDNSKEY struct for parsing DNSKEY record
type DNSRecordDNSKEY struct {
	DNSRecordCommon

	Data struct {
		Algorithm int    `json:"algorithm"`
		Flags     int    `json:"flags"`
		Protocol  int    `json:"protocol"`
		PublicKey string `json:"public_key"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDNSKEY) SetID(recordID string) *DNSRecordDNSKEY {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDNSKEY) SetComment(comment string) *DNSRecordDNSKEY {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDNSKEY) SetTags(tags []string) *DNSRecordDNSKEY {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDNSKEY) SetTTL(ttl int) *DNSRecordDNSKEY {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDNSKEY) SetZoneID(zoneID string) *DNSRecordDNSKEY {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordDNSKEY creates a new DNSKEY record.
func NewDNSRecordDNSKEY(name string, algorithm, flags, protocl int, publicKey string) *DNSRecordDNSKEY {
	r := DNSRecordDNSKEY{}
	r.Type = DNSKEY
	r.Name = name
	r.Data.Algorithm = algorithm
	r.Data.Flags = flags
	r.Data.Protocol = protocl
	r.Data.PublicKey = publicKey

	return &r
}

// DNSRecordDS struct for parsing DS record
type DNSRecordDS struct {
	DNSRecordCommon

	Data struct {
		Algorithm  int    `json:"algorithm"`
		Digest     string `json:"digest"`
		DigestType int    `json:"digest_type"`
		KeyTag     int    `json:"key_tag"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDS) SetID(recordID string) *DNSRecordDS {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDS) SetComment(comment string) *DNSRecordDS {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDS) SetTags(tags []string) *DNSRecordDS {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDS) SetTTL(ttl int) *DNSRecordDS {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordDS) SetZoneID(zoneID string) *DNSRecordDS {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordDS creates a new DS record.
func NewDNSRecordDS(name string, algorithm int, digest string, digestType, keyTag int) *DNSRecordDS {
	r := DNSRecordDS{}
	r.Type = DS
	r.Name = name
	r.Data.Algorithm = algorithm
	r.Data.Digest = digest
	r.Data.DigestType = digestType
	r.Data.KeyTag = keyTag

	return &r
}

// DNSRecordHTTPS struct for parsing HTTPS record
type DNSRecordHTTPS struct {
	DNSRecordCommon

	Data struct {
		Priority int    `json:"priority"`
		Target   string `json:"target"`
		Value    string `json:"value"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordHTTPS) SetID(recordID string) *DNSRecordHTTPS {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordHTTPS) SetComment(comment string) *DNSRecordHTTPS {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordHTTPS) SetTags(tags []string) *DNSRecordHTTPS {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordHTTPS) SetTTL(ttl int) *DNSRecordHTTPS {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordHTTPS) SetZoneID(zoneID string) *DNSRecordHTTPS {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordHTTPS creates a new HTTPS record.
func NewDNSRecordHTTPS(name string, priority int, target, value string) *DNSRecordHTTPS {
	r := DNSRecordHTTPS{}
	r.Type = HTTPS
	r.Name = name
	r.Data.Priority = priority
	r.Data.Target = target
	r.Data.Value = value

	return &r
}

type LatitudeDirection string

const (
	North LatitudeDirection = "N"
	South LatitudeDirection = "S"
)

type LongitudeDirection string

const (
	East LongitudeDirection = "E"
	West LongitudeDirection = "W"
)

// DNSRecordLOC struct for parsing LOC record
type DNSRecordLOC struct {
	DNSRecordCommon

	Data struct {
		Altitude            int                `json:"altitude"`
		LatDegrees          int                `json:"lat_degrees"`
		LatDirection        LatitudeDirection  `json:"lat_direction"`
		LatMinutes          int                `json:"lat_minutes"`
		LatSeconds          int                `json:"lat_seconds"`
		LongDegrees         int                `json:"long_degrees"`
		LongDirection       LongitudeDirection `json:"long_direction"`
		LongMinutes         int                `json:"long_minutes"`
		LongSeconds         int                `json:"long_seconds"`
		PrecisionHorizontal int                `json:"precision_horz"`
		PrecisionVertical   int                `json:"precision_vert"`
		Size                int                `json:"size"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordLOC) SetID(recordID string) *DNSRecordLOC {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordLOC) SetComment(comment string) *DNSRecordLOC {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordLOC) SetTags(tags []string) *DNSRecordLOC {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordLOC) SetTTL(ttl int) *DNSRecordLOC {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordLOC) SetZoneID(zoneID string) *DNSRecordLOC {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordLOC creates a new LOC record.
func NewDNSRecordLOC(name string, altitude, latDeg int, latDir LatitudeDirection, latMin, latSec, longDeg int, longDir LongitudeDirection, longMin, longSec, precHorizontal, precVertical, size int) *DNSRecordLOC {
	r := DNSRecordLOC{}
	r.Type = LOC
	r.Name = name
	r.Data.Altitude = altitude
	r.Data.LatDegrees = latDeg
	r.Data.LatDirection = latDir
	r.Data.LatMinutes = latMin
	r.Data.LatSeconds = latSec
	r.Data.LongDegrees = longDeg
	r.Data.LongDirection = longDir
	r.Data.LongMinutes = longMin
	r.Data.LongSeconds = longSec
	r.Data.PrecisionHorizontal = precHorizontal
	r.Data.PrecisionVertical = precVertical
	r.Data.Size = size

	return &r
}

// DNSRecordMX struct for parsing MX record
type DNSRecordMX struct {
	DNSRecordCommon

	Priority int `json:"priority"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordMX) SetID(recordID string) *DNSRecordMX {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordMX) SetComment(comment string) *DNSRecordMX {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordMX) SetTags(tags []string) *DNSRecordMX {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordMX) SetTTL(ttl int) *DNSRecordMX {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordMX) SetZoneID(zoneID string) *DNSRecordMX {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordMX creates a new MX record.
func NewDNSRecordMX(name, content string, priority int) *DNSRecordMX {
	r := DNSRecordMX{}
	r.Type = MX
	r.Name = name
	r.Content = content
	r.Priority = priority

	return &r
}

// DNSRecordNAPTR struct for parsing NAPTR record
type DNSRecordNAPTR struct {
	DNSRecordCommon

	Data struct {
		Flags       string `json:"flags"`
		Order       int    `json:"order"`
		Preference  int    `json:"preference"`
		Regex       string `json:"regex"`
		Replacement string `json:"replacement"`
		Service     string `json:"service"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNAPTR) SetID(recordID string) *DNSRecordNAPTR {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNAPTR) SetComment(comment string) *DNSRecordNAPTR {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNAPTR) SetTags(tags []string) *DNSRecordNAPTR {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNAPTR) SetTTL(ttl int) *DNSRecordNAPTR {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNAPTR) SetZoneID(zoneID string) *DNSRecordNAPTR {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordNAPTR creates a new NAPTR record.
func NewDNSRecordNAPTR(name, flags string, order, preference int, regex, replacement, service string) *DNSRecordNAPTR {
	r := DNSRecordNAPTR{}
	r.Type = NAPTR
	r.Name = name
	r.Data.Flags = flags
	r.Data.Order = order
	r.Data.Preference = preference
	r.Data.Regex = regex
	r.Data.Replacement = replacement
	r.Data.Service = service

	return &r
}

// DNSRecordNS struct for parsing NS record
type DNSRecordNS DNSRecordCommon

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNS) SetID(recordID string) *DNSRecordNS {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNS) SetComment(comment string) *DNSRecordNS {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNS) SetTags(tags []string) *DNSRecordNS {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNS) SetTTL(ttl int) *DNSRecordNS {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordNS) SetZoneID(zoneID string) *DNSRecordNS {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordNS creates a new NS record.
func NewDNSRecordNS(name, content string) *DNSRecordNS {
	r := DNSRecordNS{}
	r.Type = NS
	r.Name = name
	r.Content = content

	return &r
}

// DNSRecordPTR struct for parsing PTR record
type DNSRecordPTR DNSRecordCommon

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordPTR) SetID(recordID string) *DNSRecordPTR {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordPTR) SetComment(comment string) *DNSRecordPTR {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordPTR) SetTags(tags []string) *DNSRecordPTR {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordPTR) SetTTL(ttl int) *DNSRecordPTR {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordPTR) SetZoneID(zoneID string) *DNSRecordPTR {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordPTR creates a new PTR record.
func NewDNSRecordPTR(name, content string) *DNSRecordPTR {
	r := DNSRecordPTR{}
	r.Type = PTR
	r.Name = name
	r.Content = content

	return &r
}

// DNSRecordSMIMEA struct for parsing SMIMEA record
type DNSRecordSMIMEA struct {
	DNSRecordCommon

	Data struct {
		Certificate  string `json:"certificate"`
		MatchingType int    `json:"matching_type"`
		Selector     int    `json:"selector"`
		Usage        int    `json:"usage"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSMIMEA) SetID(recordID string) *DNSRecordSMIMEA {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSMIMEA) SetComment(comment string) *DNSRecordSMIMEA {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSMIMEA) SetTags(tags []string) *DNSRecordSMIMEA {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSMIMEA) SetTTL(ttl int) *DNSRecordSMIMEA {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSMIMEA) SetZoneID(zoneID string) *DNSRecordSMIMEA {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordSMIMEA creates a new SMIMEA record.
func NewDNSRecordSMIMEA(name, certificate string, matchingType, selector, usage int) *DNSRecordSMIMEA {
	r := DNSRecordSMIMEA{}
	r.Type = SMIMEA
	r.Name = name
	r.Data.Certificate = certificate
	r.Data.MatchingType = matchingType
	r.Data.Selector = selector
	r.Data.Usage = usage

	return &r
}

// DNSRecordSRV struct for parsing SRV record
type DNSRecordSRV struct {
	DNSRecordCommon

	Data struct {
		Name     string `json:"name"`
		Port     int    `json:"port"`
		Priority int    `json:"priority"`
		Proto    string `json:"proto"`
		Service  string `json:"service"`
		Target   string `json:"target"`
		Weight   int    `json:"weight"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSRV) SetID(recordID string) *DNSRecordSRV {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSRV) SetComment(comment string) *DNSRecordSRV {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSRV) SetTags(tags []string) *DNSRecordSRV {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSRV) SetTTL(ttl int) *DNSRecordSRV {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSRV) SetZoneID(zoneID string) *DNSRecordSRV {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordSRV creates a new SRV record.
func NewDNSRecordSRV(name string, port, priority int, proto, service, target string, weight int) *DNSRecordSRV {
	r := DNSRecordSRV{}
	r.Type = SRV
	r.Name = name // FIXME
	r.Data.Name = name
	r.Data.Port = port
	r.Data.Priority = priority
	r.Data.Proto = proto
	r.Data.Service = service
	r.Data.Target = target
	r.Data.Weight = weight

	return &r
}

// DNSRecordSSHFP struct for parsing SSHFP record
type DNSRecordSSHFP struct {
	DNSRecordCommon

	Data struct {
		Algorithm   int    `json:"algorithm"`
		Fingerprint string `json:"fingerprint"`
		Type        int    `json:"type"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSSHFP) SetID(recordID string) *DNSRecordSSHFP {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSSHFP) SetComment(comment string) *DNSRecordSSHFP {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSSHFP) SetTags(tags []string) *DNSRecordSSHFP {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSSHFP) SetTTL(ttl int) *DNSRecordSSHFP {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSSHFP) SetZoneID(zoneID string) *DNSRecordSSHFP {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordSSHFP creates a new SSHFP record.
func NewDNSRecordSSHFP(name string, algorithm int, fingerprint string, typ3 int) *DNSRecordSSHFP {
	r := DNSRecordSSHFP{}
	r.Type = SSHFP
	r.Name = name
	r.Data.Algorithm = algorithm
	r.Data.Fingerprint = fingerprint
	r.Data.Type = typ3

	return &r
}

// DNSRecordSVCB struct for parsing SVCB record
type DNSRecordSVCB struct {
	DNSRecordCommon

	Data struct {
		Priority int    `json:"priority"`
		Target   string `json:"target"`
		Value    string `json:"value"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSVCB) SetID(recordID string) *DNSRecordSVCB {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSVCB) SetComment(comment string) *DNSRecordSVCB {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSVCB) SetTags(tags []string) *DNSRecordSVCB {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSVCB) SetTTL(ttl int) *DNSRecordSVCB {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordSVCB) SetZoneID(zoneID string) *DNSRecordSVCB {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordSVCB creates a new SVCB record.
func NewDNSRecordSVCB(name string, priority int, target, value string) *DNSRecordSVCB {
	r := DNSRecordSVCB{}
	r.Type = SVCB
	r.Name = name
	r.Data.Priority = priority
	r.Data.Target = target
	r.Data.Value = value

	return &r
}

// DNSRecordTLSA struct for parsing TLSA record
type DNSRecordTLSA struct {
	DNSRecordCommon

	Data struct {
		Certificate  string `json:"certificate"`
		MatchingType int    `json:"matching_type"`
		Selector     int    `json:"selector"`
		Usage        int    `json:"usage"`
	} `json:"data"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTLSA) SetID(recordID string) *DNSRecordTLSA {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTLSA) SetComment(comment string) *DNSRecordTLSA {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTLSA) SetTags(tags []string) *DNSRecordTLSA {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTLSA) SetTTL(ttl int) *DNSRecordTLSA {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTLSA) SetZoneID(zoneID string) *DNSRecordTLSA {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordTLSA creates a new TLSA record.
func NewDNSRecordTLSA(name, certificate string, matchingType, selector, usage int) *DNSRecordTLSA {
	r := DNSRecordTLSA{}
	r.Type = TLSA
	r.Name = name
	r.Data.Certificate = certificate
	r.Data.MatchingType = matchingType
	r.Data.Selector = selector
	r.Data.Usage = usage

	return &r
}

// DNSRecordTXT struct for parsing TXT record
type DNSRecordTXT DNSRecordCommon

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTXT) SetID(recordID string) *DNSRecordTXT {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTXT) SetComment(comment string) *DNSRecordTXT {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTXT) SetTags(tags []string) *DNSRecordTXT {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTXT) SetTTL(ttl int) *DNSRecordTXT {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordTXT) SetZoneID(zoneID string) *DNSRecordTXT {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordTXT creates a new TXT record.
func NewDNSRecordTXT(name, content string) *DNSRecordTXT {
	r := DNSRecordTXT{}
	r.Type = TXT
	r.Name = name
	r.Content = content

	return &r
}

// DNSRecordURI struct for parsing URI record
type DNSRecordURI struct {
	DNSRecordCommon

	Data struct {
		Content string `json:"content"`
		Weight  int    `json:"weight"`
	} `json:"data"`
	Priority int `json:"priority"`
}

// SetID sets the `id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordURI) SetID(recordID string) *DNSRecordURI {
	r.ID = recordID
	return r
}

// SetComment sets the `comment` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordURI) SetComment(comment string) *DNSRecordURI {
	r.Comment = comment
	return r
}

// SetTags sets the `tags` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordURI) SetTags(tags []string) *DNSRecordURI {
	r.Tags = tags
	return r
}

// SetTTL sets the `ttl` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordURI) SetTTL(ttl int) *DNSRecordURI {
	r.TTL = ttl
	return r
}

// SetZoneID sets the `zone_id` value of DNS record.
//
// FIXME: repeated all over the record types
func (r *DNSRecordURI) SetZoneID(zoneID string) *DNSRecordURI {
	r.ZoneID = zoneID
	return r
}

// NewDNSRecordURI creates a new URI record.
func NewDNSRecordURI(name, content string, weight int) *DNSRecordURI {
	r := DNSRecordURI{}
	r.Type = URI
	r.Name = name
	r.Data.Content = content
	r.Data.Weight = weight

	return &r
}

// ResponseDNSRecords struct for the responses of `ListDNSRecords` function
type ResponseDNSRecords struct {
	ResponseCommon

	Result     []DNSRecordRaw `json:"result"`
	ResultInfo struct {
		Count      int `json:"count,omitempty"`
		Page       int `json:"page,omitempty"`
		PerPage    int `json:"per_page,omitempty"`
		TotalCount int `json:"total_count,omitempty"`
	} `json:"result_info,omitempty"`
}

// ResponseDNSRecordCreation struct for the responses of `CreateDNSRecord` function
type ResponseDNSRecordCreation struct {
	ResponseCommon

	Result DNSRecordRaw `json:"result"`
}

// ResponseDNSRecordDeletion struct for the responses of `DeleteDNSRecord` function
type ResponseDNSRecordDeletion struct {
	Result struct {
		ID string `json:"id"`
	} `json:"result"`
}

// ResponseDNSRecordUpdate struct for the responses of `UpdateDNSRecord` function
type ResponseDNSRecordUpdate struct {
	ResponseCommon

	Result DNSRecordRaw `json:"result"`
}
