package saasclient

import (
	"crypto/tls"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// NewSaaSClient creates a new client
func NewSaaSClient(client *resty.Client, baseURL string, apiKey string, clusterName string, log *zap.SugaredLogger) *API {
	client.SetHeader("Accept", "application/json")
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "c-n-d-e Agent",
	})
	return &API{client, log, baseURL, apiKey, clusterName}
}

// CreateClient create default client
func CreateClient() *resty.Client {

	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 10 * time.Second, errors.New("quota exceeded")
		})
	if _, err := os.Stat("/certs/client.pem"); !os.IsNotExist(err) {
		cert, err := tls.LoadX509KeyPair("/certs/client.pem", "/certs/client.key")
		if err != nil {
			log.Fatalf("ERROR client certificate: %s", err)
		}
		client.SetCertificates(cert)
	}
	return client
}

// GetDevEnvUsersForCluster retrieves DevEmv Links
func (api *API) GetDevEnvUsersForCluster() (*Cluster, error) {
	resp, err := api.client.R().
		SetResult(&Cluster{}).
		SetQueryParams(map[string]string{
			"apiKey": api.apiKey,
		}).SetPathParams(map[string]string{
		"clusterName": api.clusterName,
	}).Get(api.baseURL + "/clusters/{clusterName}")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("GetDevEnvUsersForCluster returned: " + strconv.Itoa(resp.StatusCode()))
	}
	return resp.Result().(*Cluster), err
}

// GetDevEnvUser returns a single dev env
func (api *API) GetDevEnvUser(url string) (*DevEnvUser, error) {
	resp, err := api.client.R().
		SetResult(&DevEnvUser{}).
		SetQueryParams(map[string]string{
			"apiKey": api.apiKey,
		}).Get(api.baseURL + url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("GetDevEnvUser returned: " + strconv.Itoa(resp.StatusCode()))
	}
	return resp.Result().(*DevEnvUser), err
}

// GetBuilder retrieves one builde instance
func (api *API) GetBuilder(url string) (*Builder, error) {
	resp, err := api.client.R().
		SetResult(&Builder{}).
		SetQueryParams(map[string]string{
			"apiKey": api.apiKey,
		}).Get(api.baseURL + url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("GetBuilder returned: " + strconv.Itoa(resp.StatusCode()))
	}
	return resp.Result().(*Builder), err
}

// GetBuildFile retrieves one Build File instance
func (api *API) GetBuildFile(url string) (*BuildFile, error) {
	resp, err := api.client.R().
		SetResult(&BuildFile{}).
		SetQueryParams(map[string]string{
			"apiKey": api.apiKey,
		}).Get(api.baseURL + url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("GetBuildFile returned: " + strconv.Itoa(resp.StatusCode()))
	}
	return resp.Result().(*BuildFile), err
}

// PutClusterStatus for a HTTP PUT as a cluster status
func (api *API) PutClusterStatus(stat *ClusterStatus) error {
	resp, err := api.client.R().
		SetBody(stat).
		SetQueryParams(map[string]string{
			"apiKey": api.apiKey,
		}).SetPathParams(map[string]string{
		"clusterName": api.clusterName,
	}).Put(api.baseURL + "/clusters/{clusterName}/devenvusers/metrics")
	if resp.StatusCode() != 204 {
		return errors.New("PutClusterStatus returned: " + strconv.Itoa(resp.StatusCode()))
	}
	if err != nil {
		return err
	}
	return err
}
