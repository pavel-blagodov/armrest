package app

import "time"

type TasksResponse []struct {
	StatusID                 string           `json:"statusId"`
	Type                     string           `json:"type"`
	Subtype                  string           `json:"subtype"`
	RecommendedRefreshPeriod float64          `json:"recommendedRefreshPeriod"`
	Status                   string           `json:"status"`
	Progress                 float64          `json:"progress"`
	PerNode                  PerNode          `json:"perNode"`
	DetailedProgress         DetailedProgress `json:"detailedProgress"`
	StageInfo                StageInfo        `json:"stageInfo"`
	RebalanceID              string           `json:"rebalanceId"`
	NodesInfo                NodesInfo        `json:"nodesInfo"`
	MasterNode               string           `json:"masterNode"`
}

type PerNode map[string]PerNodeDetails

type PerNodeDetails struct {
	Progress float64 `json:"progress"`
}

type DetailedProgress struct {
	Bucket       string  `json:"bucket"`
	BucketNumber int     `json:"bucketNumber"`
	BucketsCount int     `json:"bucketsCount"`
	PerNode      PerNode `json:"perNode"`
}
type Backup struct {
	StartTime     bool `json:"startTime"`
	CompletedTime bool `json:"completedTime"`
	TimeTaken     bool `json:"timeTaken"`
}
type Analytics struct {
	StartTime     bool `json:"startTime"`
	CompletedTime bool `json:"completedTime"`
	TimeTaken     bool `json:"timeTaken"`
}
type Eventing struct {
	StartTime     bool `json:"startTime"`
	CompletedTime bool `json:"completedTime"`
	TimeTaken     bool `json:"timeTaken"`
}
type Search struct {
	StartTime     bool `json:"startTime"`
	CompletedTime bool `json:"completedTime"`
	TimeTaken     bool `json:"timeTaken"`
}
type Index struct {
	StartTime     bool `json:"startTime"`
	CompletedTime bool `json:"completedTime"`
	TimeTaken     bool `json:"timeTaken"`
}
type PerNodeProgress map[string]float64
type Move struct {
	AverageTime    float64 `json:"averageTime"`
	TotalCount     int     `json:"totalCount"`
	RemainingCount int     `json:"remainingCount"`
}
type Backfill struct {
	AverageTime float64 `json:"averageTime"`
}
type Takeover struct {
	AverageTime float64 `json:"averageTime"`
}
type Persistence struct {
	AverageTime float64 `json:"averageTime"`
}
type VbucketLevelInfo struct {
	Move        Move        `json:"move"`
	Backfill    Backfill    `json:"backfill"`
	Takeover    Takeover    `json:"takeover"`
	Persistence Persistence `json:"persistence"`
}
type PerNodeReplicationInfo struct {
	InDocsTotal  int `json:"inDocsTotal"`
	InDocsLeft   int `json:"inDocsLeft"`
	OutDocsTotal int `json:"outDocsTotal"`
	OutDocsLeft  int `json:"outDocsLeft"`
}

type ReplicationInfo map[string]PerNodeReplicationInfo
type TravelSample struct {
	VbucketLevelInfo VbucketLevelInfo `json:"vbucketLevelInfo"`
	ReplicationInfo  ReplicationInfo  `json:"replicationInfo"`
	StartTime        time.Time        `json:"startTime"`
	CompletedTime    bool             `json:"completedTime"`
	TimeTaken        int              `json:"timeTaken"`
}
type Details struct {
	TravelSample TravelSample `json:"travel-sample"`
}
type Data struct {
	TotalProgress   float64         `json:"totalProgress"`
	PerNodeProgress PerNodeProgress `json:"perNodeProgress"`
	StartTime       time.Time       `json:"startTime"`
	CompletedTime   bool            `json:"completedTime"`
	TimeTaken       int             `json:"timeTaken"`
	Details         Details         `json:"details"`
}
type Query struct {
	StartTime     bool `json:"startTime"`
	CompletedTime bool `json:"completedTime"`
	TimeTaken     bool `json:"timeTaken"`
}
type StageInfo struct {
	Backup    Backup    `json:"backup"`
	Analytics Analytics `json:"analytics"`
	Eventing  Eventing  `json:"eventing"`
	Search    Search    `json:"search"`
	Index     Index     `json:"index"`
	Data      Data      `json:"data"`
	Query     Query     `json:"query"`
}
type NodesInfo struct {
	ActiveNodes []string `json:"active_nodes"`
	KeepNodes   []string `json:"keep_nodes"`
	EjectNodes  []string `json:"eject_nodes"`
	DeltaNodes  []any    `json:"delta_nodes"`
	FailedNodes []any    `json:"failed_nodes"`
}
