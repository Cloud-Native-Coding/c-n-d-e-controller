package metrics

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/inf.v0"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"saas-controller.cloud-native-coding.dev/controllers"
	"saas-controller.cloud-native-coding.dev/saasclient"
)

// Metrics base struct
type Metrics struct {
	api        *saasclient.API
	k8sService *controllers.K8sService
	log        *zap.SugaredLogger
}

var (
	keycloakHost = os.Getenv("CNDE_KEYCLOAK_HOST")
)

// NewMetrics returns new Metrics instance
func NewMetrics(api *saasclient.API, c *rest.Config, log *zap.SugaredLogger) *Metrics {

	k8sService := controllers.NewK8sService(c, log)
	return &Metrics{api, k8sService, log}
}

// Calculate is the reconciliation loop
func (c *Metrics) Calculate() error {

	cluster, err := c.api.GetDevEnvUsersForCluster()
	if err != nil {
		c.log.Error(err, "could not read from API")
		return err
	}

	userFromAPI := []*saasclient.DevEnvUser{}
	for _, cde := range cluster.DevEnvUsers {
		user, err := c.api.GetDevEnvUser(cde.Links.Self)
		if err != nil {
			c.log.Error(err, "Error retrieving single Dev Env User")
			return err
		}
		userFromAPI = append(userFromAPI, user)
	}

	devEnvCRs, err := c.k8sService.GetDevEnvs()
	if statusError, isStatus := err.(*k8sErrors.StatusError); isStatus {
		c.log.Error(err, "Error getting DevEnvs %v\n", statusError.ErrStatus.Message)
		return err
	} else if err != nil {
		c.log.Error(err, "Unknown K8s Error")
		return err
	}
	c.log.Infow("Calculating Metrics for ", "num: ", len(devEnvCRs.Items))

	clusterStatus := saasclient.ClusterStatus{}

	for _, de := range devEnvCRs.Items {
		allPodMetrics, err := c.k8sService.GetPodMetrics(de.Status.Realm)
		if err != nil {
			c.log.Error(err, "Unknown K8s Error")
		}
		if de.Status.Realm == "" || len(allPodMetrics.Items) == 0 {
			if deu, err := c.searchDevEnv(de.Name, userFromAPI); err == nil {
				clusterStatus = append(clusterStatus, saasclient.DevEnvStatus{
					Devenvuser: deu,
					Status:     string(de.Status.Build),
					CPU:        "",
					Memory:     "",
				})
			}
		} else {
			for _, podMetric := range allPodMetrics.Items {
				cCPU := new(inf.Dec)
				cMem := int64(0)
				for _, container := range podMetric.Containers {
					cCPU = cCPU.Add(cCPU, container.Usage.Cpu().AsDec())
					if m, ok := container.Usage.Memory().AsInt64(); ok {
						cMem += m
					}
				}
				if deu, err := c.searchDevEnv(de.Name, userFromAPI); err == nil {
					clusterStatus = append(clusterStatus, saasclient.DevEnvStatus{
						Devenvuser: deu,
						Status:     string(de.Status.Build),
						CPU:        cCPU.String(),
						Memory:     strconv.FormatInt(cMem, 10),
					})
				}
			}
		}
	}
	//c.log.Infow("retrieved Metrics", "clusterStatus", clusterStatus)
	c.api.PutClusterStatus(&clusterStatus)

	return nil
}

func (c *Metrics) searchDevEnv(devEnvName string, userFromAPI []*saasclient.DevEnvUser) (string, error) {
	for _, user := range userFromAPI {
		if strings.ToLower(user.Name) == devEnvName {
			return user.Links.Self, nil
		}
	}
	return "", errors.New("DevEnv not fould in API response")
}
