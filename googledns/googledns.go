package googledns

// Partially based on https://github.com/xenolf/lego/blob/master/providers/dns/gcloud/googlecloud.go

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	dns "google.golang.org/api/dns/v1"
)

// Config is used to configure the creation of the DNSProvider
type Config struct {
	Project            string
	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPClient         *http.Client `json:"-"`
}

// NewDefaultConfig returns a default configuration for the DNSProvider
func NewDefaultConfig() *Config {
	return &Config{
		TTL:                120,
		PropagationTimeout: 180 * time.Second,
		PollingInterval:    5 * time.Second,
	}
}

// DNSProvider is an implementation of the DNSProvider interface.
type DNSProvider struct {
	Config *Config
	Client *dns.Service
}

// NewDNSProviderServiceAccount uses the supplied service account JSON file
// to return a DNSProvider instance configured for Google Cloud DNS.
func NewDNSProviderServiceAccount(saFile string) (*DNSProvider, error) {
	if saFile == "" {
		return nil, fmt.Errorf("googlecloud: Service Account file missing")
	}

	dat, err := ioutil.ReadFile(saFile)
	if err != nil {
		return nil, fmt.Errorf("googlecloud: unable to read Service Account file: %v", err)
	}

	// read project id from service account file
	var datJSON struct {
		ProjectID string `json:"project_id"`
	}
	err = json.Unmarshal(dat, &datJSON)
	if err != nil || datJSON.ProjectID == "" {
		return nil, fmt.Errorf("googlecloud: project ID not found in Google Cloud Service Account file")
	}
	project := datJSON.ProjectID

	conf, err := google.JWTConfigFromJSON(dat, dns.NdevClouddnsReadwriteScope)
	if err != nil {
		return nil, fmt.Errorf("googlecloud: unable to acquire config: %v", err)
	}
	client := conf.Client(context.Background())

	config := NewDefaultConfig()
	config.Project = project
	config.HTTPClient = client

	return NewDNSProviderConfig(config)
}

// NewDNSProviderConfig return a DNSProvider instance configured for Google Cloud DNS.
func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("googlecloud: the configuration of the DNS provider is nil")
	}

	svc, err := dns.New(config.HTTPClient)
	if err != nil {
		return nil, fmt.Errorf("googlecloud: unable to create Google Cloud DNS service: %v", err)
	}

	return &DNSProvider{Config: config, Client: svc}, nil
}

func NewDNSProvider() (*DNSProvider, error) {
	saFile := viper.GetString("SA_FILE")
	return NewDNSProviderServiceAccount(saFile)
}

type DNSRecord struct {
	Ips    *[]net.IP
	Domain *string
}

func MakeResourceRecordSet(domain string, ips []string, ttl int) []*dns.ResourceRecordSet {
	rrec := &dns.ResourceRecordSet{
		Name:    domain + ".",
		Rrdatas: ips,
		Ttl:     int64(ttl),
		Type:    "A",
	}

	return []*dns.ResourceRecordSet{rrec}
}

func (d *DNSProvider) MakeChange(rec DNSRecord, adding bool) error {
	zone := viper.GetString("DNS_ZONE")
	newIPs := make([]string, len(*rec.Ips))
	for i, v := range *rec.Ips {
		newIPs[i] = v.String()
	}
	project := d.Config.Project

	domain := *rec.Domain
	oldIPs, err := d.GetResourceRecordSets(domain)
	if err != nil {
		return fmt.Errorf("Could not get resource sets for %v", domain)
	}

	change := &dns.Change{}
	if adding {
		// Adding records
		newIPs = UniqueMerge(newIPs, oldIPs)
		change.Additions = MakeResourceRecordSet(domain, newIPs, d.Config.TTL)
		if len(oldIPs) > 0 {
			change.Deletions = MakeResourceRecordSet(domain, oldIPs, d.Config.TTL)
		}
		fmt.Printf("Up to date records after changes: %v\n", newIPs)
	} else {
		// Deleteing records
		change.Deletions = MakeResourceRecordSet(domain, oldIPs, d.Config.TTL)
		remainingIPs := Diff(oldIPs, newIPs)
		if len(remainingIPs) > 0 {
			change.Additions = MakeResourceRecordSet(domain, remainingIPs, d.Config.TTL)
		}
		fmt.Printf("Up to date records after changes: %v\n", remainingIPs)
	}

	cli := d.Client
	resp, err := cli.Changes.Create(project, zone, change).Do()
	if err != nil {
		return fmt.Errorf("googlecloud: %v", err)
	}

	for resp.Status == "pending" {
		time.Sleep(time.Second)

		resp, err = cli.Changes.Get(project, zone, resp.Id).Do()
		if err != nil {
			return fmt.Errorf("googlecloud: %v", err)
		}
	}

	return nil

}

func (d *DNSProvider) GetResourceRecordSets(domain string) ([]string, error) {
	project := d.Config.Project
	zone := viper.GetString("DNS_ZONE")

	resp, err := d.Client.ResourceRecordSets.List(project, zone).Name(domain + ".").Do()
	if err != nil {
		return []string{}, fmt.Errorf("googlecloud: %v", err)
	}

	if len(resp.Rrsets) == 1 {
		return resp.Rrsets[0].Rrdatas, nil
	} else {
		return []string{}, nil
	}
}

func UniqueMerge(s1 []string, s2 []string) []string {
	merged := append([]string{}, s1...)
	for _, s := range s2 {
		if !Contains(merged, s) {
			merged = append(merged, s)
		}
	}
	return merged

}

func Diff(s1 []string, s2 []string) []string {
	onlyInS1 := []string{}
	for _, s := range s1 {
		if !Contains(s2, s) {
			onlyInS1 = append(onlyInS1, s)
		}
	}
	return onlyInS1
}

func Contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
