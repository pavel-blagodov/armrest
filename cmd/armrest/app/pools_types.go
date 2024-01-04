package app

type Pools struct {
	IsAdminCreds          bool              `json:"isAdminCreds"`
	IsROAdminCreds        bool              `json:"isROAdminCreds"`
	IsEnterprise          bool              `json:"isEnterprise"`
	ConfigProfile         string            `json:"configProfile"`
	AllowedServices       []string          `json:"allowedServices"`
	IsDeveloperPreview    bool              `json:"isDeveloperPreview"`
	PackageVariant        string            `json:"packageVariant"`
	Pools                 []PoolsPools      `json:"pools"`
	Settings              Settings          `json:"settings"`
	UUID                  string            `json:"uuid"`
	ImplementationVersion string            `json:"implementationVersion"`
	ComponentsVersion     ComponentsVersion `json:"componentsVersion"`
}
type PoolsPools struct {
	Name         string `json:"name"`
	URI          string `json:"uri"`
	StreamingURI string `json:"streamingUri"`
}
type Settings struct {
	MaxParallelIndexers string `json:"maxParallelIndexers"`
	ViewUpdateDaemon    string `json:"viewUpdateDaemon"`
}
type ComponentsVersion struct {
	Esaml     string `json:"esaml"`
	Xmerl     string `json:"xmerl"`
	Inets     string `json:"inets"`
	OsMon     string `json:"os_mon"`
	Sasl      string `json:"sasl"`
	Ale       string `json:"ale"`
	Stdlib    string `json:"stdlib"`
	Asn1      string `json:"asn1"`
	Crypto    string `json:"crypto"`
	PublicKey string `json:"public_key"`
	Chronicle string `json:"chronicle"`
	Kernel    string `json:"kernel"`
	NsServer  string `json:"ns_server"`
	Ssl       string `json:"ssl"`
	Lhttpc    string `json:"lhttpc"`
}
