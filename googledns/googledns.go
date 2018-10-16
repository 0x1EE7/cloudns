package googledns

// Partially based on https://github.com/xenolf/lego/blob/master/providers/dns/gcloud/googlecloud.go

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
