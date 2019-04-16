package main

import "net"

// ConfigFile used to construct adGuard Home config file
type ConfigFile struct {
	Letsencrypt struct {
		Enabled          bool              `yaml:"enabled"`
		Production       bool              `yaml:"production"`
		Timeout          int               `yaml:"timeout"`
		Email            string            `yaml:"email"`
		Provider         string            `yaml:"provider"`
		ProviderSettings map[string]string `yaml:"provider_settings"`
	}
	DynamicDNS struct {
		Enabled  bool   `yaml:"enabled"`
		APIToken string `yaml:"api_token"`
	}
	BindHost string `yaml:"bind_host"`
	BindPort int    `yaml:"bind_port"`
	UserName string `yaml:"auth_name"`
	Password string `yaml:"auth_pass"`
	Language string `yaml:"language"`
	DNS      struct {
		BindHost            net.IP   `yaml:"bind_host"`
		Port                int      `yaml:"port"`
		ProtectionEnabled   bool     `yaml:"protection_enabled"`
		FilteringEnabled    bool     `yaml:"filtering_enabled"`
		BlockedResponseTTL  int      `yaml:"blocked_response_ttl"`
		QuerylogEnabled     bool     `yaml:"querylog_enabled"`
		Ratelimit           int      `yaml:"ratelimit"`
		RatelimitWhitelist  []string `yaml:"ratelimit_whitelist"`
		RefuseAny           bool     `yaml:"refuse_any"`
		BootstrapDNS        string   `yaml:"bootstrap_dns"`
		ParentalSensitivity int      `yaml:"parental_sensitivity"`
		ParentalEnabled     bool     `yaml:"parental_enabled"`
		SafesearchEnabled   bool     `yaml:"safesearch_enabled"`
		SafebrowsingEnabled bool     `yaml:"safebrowsing_enabled"`
		UpstreamDNS         []string `yaml:"upstream_dns"`
	}
	TLS struct {
		Enabled          bool   `yaml:"enabled"`
		ServerName       string `yaml:"server_name"`
		ForceHTTPS       bool   `yaml:"force_https"`
		PortHTTPS        int    `yaml:"port_https"`
		PortDNSOverTLS   int    `yaml:"port_dns_over_tls"`
		CertificateChain string `yaml:"certificate_chain"`
		PrivateKey       string `yaml:"private_key"`
	}
	Filters []struct {
		Enabled bool   `yaml:"enabled"`
		URL     string `yaml:"url"`
		Name    string `yaml:"name"`
		ID      int    `yaml:"id"`
	}
	UserRules []string `yaml:"user_rules"`
	DHCP      struct {
		Enabled       bool   `yaml:"enabled"`
		InterfaceName string `yaml:"interface_name"`
		GatewayIP     string `yaml:"gateway_ip"`
		SubnetMask    string `yaml:"subnet_mask"`
		RangeStart    string `yaml:"range_start"`
		RangeEnd      string `yaml:"range_end"`
		LeaseDuration int    `yaml:"lease_duration"`
	}
	LogFile       string `yaml:"log_file"`
	Verbose       bool   `yaml:"verbose"`
	SchemaVersion int    `yaml:"schema_version"`
}
