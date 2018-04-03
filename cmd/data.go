package cmd

type CertEntry struct {
	CertCn        string `json:"cert_cn"`
	DaysRemaining int    `json:"days_remaining"`
	Expiry        string `json:"expiry"`
	Health        string `json:"health"`
	Path          string `json:"path"`
	Serial        float64  `json:"serial"`
	SerialHex     string `json:"serial_hex"`
}

type Server struct {
	Etcd []*CertEntry `json:"etcd"`
	Kubeconfigs []*CertEntry `json:"kubeconfigs"`
	Meta struct {
		CheckedAtTime  string `json:"checked_at_time"`
		ShowAll        string `json:"show_all"`
		WarnBeforeDate string `json:"warn_before_date"`
		WarningDays    int    `json:"warning_days"`
	} `json:"meta"`
	OcpCerts []*CertEntry `json:"ocp_certs"`
	Registry []*CertEntry `json:"registry"`
	Router   []*CertEntry `json:"router"`
}

type CertExpiryReport struct {
	Data map[string]*Server `json:"data"`
	Summary struct {
		Expired int `json:"expired"`
		Ok      int `json:"ok"`
		Total   int `json:"total"`
		Warning int `json:"warning"`
	} `json:"summary"`
}
