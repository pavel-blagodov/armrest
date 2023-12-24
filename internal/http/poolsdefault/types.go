package types

type PoolsDefault struct {
	Etag                   string                 `json:"etag"`
	Alerts                 []any                  `json:"alerts"`
	AlertsSilenceURL       string                 `json:"alertsSilenceURL"`
	AutoCompactionSettings AutoCompactionSettings `json:"autoCompactionSettings"`
	Buckets                Buckets                `json:"buckets"`
	Controllers            Controllers            `json:"controllers"`
	Counters               Counters               `json:"counters"`
	FastWarmupSettings     FastWarmupSettings     `json:"fastWarmupSettings"`
	MaxBucketCount         int                    `json:"maxBucketCount"`
	Name                   string                 `json:"name"`
	NodeStatusesURI        string                 `json:"nodeStatusesUri"`
	Nodes                  []Nodes                `json:"nodes"`
	RebalanceProgressURI   string                 `json:"rebalanceProgressUri"`
	RebalanceStatus        string                 `json:"rebalanceStatus"`
	RemoteClusters         RemoteClusters         `json:"remoteClusters"`
	ServerGroupsURI        string                 `json:"serverGroupsUri"`
	StopRebalanceURI       string                 `json:"stopRebalanceUri"`
	StorageTotals          StorageTotals          `json:"storageTotals"`
	Tasks                  Tasks                  `json:"tasks"`
	VisualSettingsURI      string                 `json:"visualSettingsUri"`
}
type DatabaseFragmentationThreshold struct {
	Percentage int    `json:"percentage"`
	Size       string `json:"size"`
}
type ViewFragmentationThreshold struct {
	Percentage int    `json:"percentage"`
	Size       string `json:"size"`
}
type AutoCompactionSettings struct {
	DatabaseFragmentationThreshold DatabaseFragmentationThreshold `json:"databaseFragmentationThreshold"`
	ParallelDBAndViewCompaction    bool                           `json:"parallelDBAndViewCompaction"`
	ViewFragmentationThreshold     ViewFragmentationThreshold     `json:"viewFragmentationThreshold"`
}
type Buckets struct {
	TerseBucketsBase          string `json:"terseBucketsBase"`
	TerseStreamingBucketsBase string `json:"terseStreamingBucketsBase"`
	URI                       string `json:"uri"`
}
type AddNode struct {
	URI string `json:"uri"`
}
type ClusterLogsCollection struct {
	CancelURI string `json:"cancelURI"`
	StartURI  string `json:"startURI"`
}
type EjectNode struct {
	URI string `json:"uri"`
}
type FailOver struct {
	URI string `json:"uri"`
}
type ReAddNode struct {
	URI string `json:"uri"`
}
type ReFailOver struct {
	URI string `json:"uri"`
}
type Rebalance struct {
	URI string `json:"uri"`
}
type Replication struct {
	CreateURI   string `json:"createURI"`
	ValidateURI string `json:"validateURI"`
}
type SetAutoCompaction struct {
	URI         string `json:"uri"`
	ValidateURI string `json:"validateURI"`
}
type SetFastWarmup struct {
	URI         string `json:"uri"`
	ValidateURI string `json:"validateURI"`
}
type SetRecoveryType struct {
	URI string `json:"uri"`
}
type StartGracefulFailover struct {
	URI string `json:"uri"`
}
type Controllers struct {
	AddNode               AddNode               `json:"addNode"`
	ClusterLogsCollection ClusterLogsCollection `json:"clusterLogsCollection"`
	EjectNode             EjectNode             `json:"ejectNode"`
	FailOver              FailOver              `json:"failOver"`
	ReAddNode             ReAddNode             `json:"reAddNode"`
	ReFailOver            ReFailOver            `json:"reFailOver"`
	Rebalance             Rebalance             `json:"rebalance"`
	Replication           Replication           `json:"replication"`
	SetAutoCompaction     SetAutoCompaction     `json:"setAutoCompaction"`
	SetFastWarmup         SetFastWarmup         `json:"setFastWarmup"`
	SetRecoveryType       SetRecoveryType       `json:"setRecoveryType"`
	StartGracefulFailover StartGracefulFailover `json:"startGracefulFailover"`
}
type Counters struct {
}
type FastWarmupSettings struct {
	FastWarmupEnabled  bool `json:"fastWarmupEnabled"`
	MinItemsThreshold  int  `json:"minItemsThreshold"`
	MinMemoryThreshold int  `json:"minMemoryThreshold"`
}
type InterestingStats struct {
	CmdGet                   int `json:"cmd_get"`
	CouchDocsActualDiskSize  int `json:"couch_docs_actual_disk_size"`
	CouchDocsDataSize        int `json:"couch_docs_data_size"`
	CouchViewsActualDiskSize int `json:"couch_views_actual_disk_size"`
	CouchViewsDataSize       int `json:"couch_views_data_size"`
	CurrItems                int `json:"curr_items"`
	CurrItemsTot             int `json:"curr_items_tot"`
	EpBgFetched              int `json:"ep_bg_fetched"`
	GetHits                  int `json:"get_hits"`
	MemUsed                  int `json:"mem_used"`
	Ops                      int `json:"ops"`
	VbReplicaCurrItems       int `json:"vb_replica_curr_items"`
}
type Ports struct {
	Direct    int `json:"direct"`
	HTTPSCAPI int `json:"httpsCAPI"`
	HTTPSMgmt int `json:"httpsMgmt"`
	Proxy     int `json:"proxy"`
	SslProxy  int `json:"sslProxy"`
}
type SystemStats struct {
	CPUUtilizationRate float64 `json:"cpu_utilization_rate"`
	MemFree            int64   `json:"mem_free"`
	MemTotal           int64   `json:"mem_total"`
	SwapTotal          int64   `json:"swap_total"`
	SwapUsed           int     `json:"swap_used"`
}
type Nodes struct {
	ClusterCompatibility int              `json:"clusterCompatibility"`
	ClusterMembership    string           `json:"clusterMembership"`
	CouchAPIBase         string           `json:"couchApiBase"`
	Hostname             string           `json:"hostname"`
	InterestingStats     InterestingStats `json:"interestingStats"`
	McdMemoryAllocated   int              `json:"mcdMemoryAllocated"`
	McdMemoryReserved    int              `json:"mcdMemoryReserved"`
	MemoryFree           int64            `json:"memoryFree"`
	MemoryTotal          int64            `json:"memoryTotal"`
	Os                   string           `json:"os"`
	OtpCookie            string           `json:"otpCookie"`
	OtpNode              string           `json:"otpNode"`
	Ports                Ports            `json:"ports"`
	RecoveryType         string           `json:"recoveryType"`
	Status               string           `json:"status"`
	SystemStats          SystemStats      `json:"systemStats"`
	ThisNode             bool             `json:"thisNode"`
	Uptime               string           `json:"uptime"`
	Version              string           `json:"version"`
}
type RemoteClusters struct {
	URI         string `json:"uri"`
	ValidateURI string `json:"validateURI"`
}
type Hdd struct {
	Free       int64 `json:"free"`
	QuotaTotal int64 `json:"quotaTotal"`
	Total      int64 `json:"total"`
	Used       int64 `json:"used"`
	UsedByData int   `json:"usedByData"`
}
type RAM struct {
	QuotaTotal        int   `json:"quotaTotal"`
	QuotaTotalPerNode int   `json:"quotaTotalPerNode"`
	QuotaUsed         int   `json:"quotaUsed"`
	QuotaUsedPerNode  int   `json:"quotaUsedPerNode"`
	Total             int64 `json:"total"`
	Used              int64 `json:"used"`
	UsedByData        int   `json:"usedByData"`
}
type StorageTotals struct {
	Hdd Hdd `json:"hdd"`
	RAM RAM `json:"ram"`
}
type Tasks struct {
	URI string `json:"uri"`
}
