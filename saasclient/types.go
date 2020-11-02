package saasclient

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// API base struct
type API struct {
	client      *resty.Client
	log         *zap.SugaredLogger
	baseURL     string
	apiKey      string
	clusterName string
}

// ClusterLinks is defined by
type ClusterLinks struct {
	Self      string
	Buildfile string
}

// ClusterDevEnv ist defined by
type ClusterDevEnv struct {
	ID          int
	BuildfileID int
	Links       ClusterLinks
}

// Cluster is the cluster endpoint respose
type Cluster struct {
	Name        string
	APIKey      string
	DevEnvUsers []ClusterDevEnv
}

// DevEnvUserLinks is defined by
type DevEnvUserLinks struct {
	Self      string
	Buildfile string
}

// DevEnvUser is the devenvuser endpoint response
type DevEnvUser struct {
	Name                string
	BuildfileID         int
	BuildFile           BuildFile // calculated property
	Builder             Builder   // calculated property
	DeleteVolume        bool
	ClusterRoleName     string
	RoleName            string
	DevEnvImage         string
	ContainerVolumeSize string
	HomeVolumeSize      string
	Email               string
	UserEnvDomain       string
	Links               DevEnvUserLinks
}

// BuildFileLinks is defined by
type BuildFileLinks struct {
	Self    string
	Builder string
}

// BuildFile Stores the build file for creating dev env container
type BuildFile struct {
	Name      string
	Value     string
	BuilderID int
	ID        int
	Links     BuildFileLinks
}

// Builder stores the attributes for building a BuildFile
type Builder struct {
	Name  string
	Value string
}

// DevEnvStatus is defined by
type DevEnvStatus struct {
	Devenvuser string `json:"devenvuser,omitempty"`
	Status     string `json:"status,omitempty"`
	CPU        string `json:"cpu,omitempty"`
	Memory     string `json:"memory,omitempty"`
}

// ClusterStatus is defined by
type ClusterStatus []DevEnvStatus
