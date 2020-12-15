package cloudflare

import (
	"goddns/lib/providers"
	"log"
	"os"
)

type Cloudflare struct {
	client *CloudflareAPI
}

func init() {
	providers.RegisterProvider("cloudflare", NewProvider)
}

var ZONEID, HOST string

func NewProvider() (providers.Provider, error) {
	// Check for use of any deprecated variables first, point to how to update
	if os.Getenv("CLOUDFLARE_APIKEY") != "" {
		log.Fatal("Do not use CLOUDFLARE_APIKEY, see https://github.com/hugomd/cloudflare-ddns#deprecated-environment-variables")
	}
	if os.Getenv("CLOUDFLARE_EMAIL") != "" {
		log.Fatal("Do not use CLOUDFLARE_EMAIL, see https://github.com/hugomd/cloudflare-ddns#deprecated-environment-variables")
	}
	if os.Getenv("CLOUDFLARE_ZONE") != "" {
		log.Fatal("Do not use CLOUDFLARE_ZONE, see https://github.com/hugomd/cloudflare-ddns#deprecated-environment-variables")
	}

	APITOKEN := os.Getenv("CLOUDFLARE_APITOKEN")
	if APITOKEN == "" {
		log.Fatal("CLOUDFLARE_APITOKEN env. variable is required")
	}

	ZONEID = os.Getenv("CLOUDFLARE_ZONEID")
	if ZONEID == "" {
		log.Fatal("CLOUDFLARE_ZONEID env. variable is required")
	}

	HOST = os.Getenv("CLOUDFLARE_HOST")
	if HOST == "" {
		log.Fatal("CLOUDFLARE_HOST env. variable is required")
	}

	api, err := NewCloudflareClient(APITOKEN, ZONEID, HOST)

	if err != nil {
		return nil, err
	}

	provider := &Cloudflare{
		client: api,
	}

	return provider, nil
}

func (api *Cloudflare) AddRecord(name, ip string) error {
	record := Record{
		Type:    "A",
		Content: ip,
		Name:    name,
		TTL:     600,
		Proxied: false,
	}
	err := api.client.AddDNSRecord(record)
	if err != nil {
		return err
	}
	return nil
}

func (api *Cloudflare) UpdateRecord(ip string) error {
	records, err := api.client.ListDNSRecords()
	if err != nil {
		return err
	}

	var record Record
	for i := range records {
		if records[i].Name == HOST {
			record = records[i]
		}
	}

	if record == (Record{}) {
		err = api.AddRecord(HOST, ip)
		if err != nil {
			return err
		}
		log.Printf("Host not found. Created %s point to %s", HOST, ip)
		return nil
	}

	if ip != record.Content {
		record.Content = ip
		err = api.client.UpdateDNSRecord(record)
		if err != nil {
			return err
		}
		log.Printf("IP changed, updated to %s", ip)
		return nil
	}
	log.Print("No change in IP, not updating record")

	return nil
}
