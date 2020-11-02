package controllers

import (
	"os"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"saas-controller.cloud-native-coding.dev/saasclient"
)

// Controller base struct
type Controller struct {
	api          *saasclient.API
	keycloakHost string
	k8sService   *K8sService
	log          *zap.SugaredLogger
}

var (
	keycloakHost = os.Getenv("CNDE_KEYCLOAK_HOST")
)

// NewController returns new controller instance
func NewController(api *saasclient.API, c *rest.Config, log *zap.SugaredLogger) *Controller {
	k8sService := NewK8sService(c, log)
	return &Controller{api, keycloakHost, k8sService, log}
}

// Reconcile is the reconciliation loop
func (c *Controller) Reconcile() error {

	//c.log.Infow("-- retrieving cluster from SaaS --")

	cluster, err := c.api.GetDevEnvUsersForCluster()
	if err != nil {
		c.log.Error(err, "could not read from API")
		return err
	}

	//c.log.Infow("retrieving dev env users from SaaS-API", "cluster", cluster)

	//c.log.Infow("-- retrieving K8s resources --")

	devEnvCRs, err := c.k8sService.GetDevEnvs()
	if statusError, isStatus := err.(*errors.StatusError); isStatus {
		c.log.Error(err, "Error getting DevEnvs %v\n", statusError.ErrStatus.Message)
		return err
	} else if err != nil {
		c.log.Error(err, "Unknown K8s Error")
		return err
	}
	c.log.Infow("Got DevEnv CRs", "num: ", len(devEnvCRs.Items))

	// builderCRs, err := c.k8sService.GetBuilders()
	// if statusError, isStatus := err.(*errors.StatusError); isStatus {
	// 	c.log.Error(err, "Error getting Builders %v\n", statusError.ErrStatus.Message)
	// 	return err
	// } else if err != nil {
	// 	c.log.Error(err, "Unknown K8s Error")
	// 	return err
	// }
	// c.log.Infow("Got Builder CRs", "num: ", len(builderCRs.Items))

	// buildFileCMs, err := c.k8sService.GetBuildFiles()
	// if statusError, isStatus := err.(*errors.StatusError); isStatus {
	// 	c.log.Error(err, "Error getting BuildFiles %v\n", statusError.ErrStatus.Message)
	// 	return err
	// } else if err != nil {
	// 	c.log.Error(err, "Unknown K8s Error")
	// 	return err
	// }
	// c.log.Infow("Got BuildFile CMs", "num: ", len(buildFileCMs.Items))

	userFromAPI := []*saasclient.DevEnvUser{}
	buildFilesFromAPI := []*saasclient.BuildFile{}
	builderFromAPI := []*saasclient.Builder{}
	for _, cde := range cluster.DevEnvUsers {

		//c.log.Infow("retrieving DevEnvUser", "User URL: ", cde.Links.Self)
		user, err := c.api.GetDevEnvUser(cde.Links.Self)
		if err != nil {
			c.log.Error(err, "Error retrieving single Dev Env User")
			return err
		}
		userFromAPI = append(userFromAPI, user)

		if cde.Links.Buildfile != "" {
			//c.log.Infow("retrieving Buildfile", "User URL: ", cde.Links.Buildfile)
			bf, err := c.api.GetBuildFile(cde.Links.Buildfile)
			if err != nil {
				c.log.Error(err, "Error retrieving single Buildfile")
				return err
			}
			buildFilesFromAPI = append(buildFilesFromAPI, bf)
			user.BuildFile = *bf

			//c.log.Infow("retrieving Builder", "User URL: ", bf.Links.Builder)
			bldr, err := c.api.GetBuilder(bf.Links.Builder)
			if err != nil {
				c.log.Error(err, "Error retrieving single Builder")
				return err
			}
			builderFromAPI = append(builderFromAPI, bldr)
			user.Builder = *bldr
		}
	}

	//c.log.Infow("-- Calculating Diffs --")

	//builderToCreate := []*saasclient.Builder{}

	devEnvsToCreate := c.devEnvsToCreate(userFromAPI, devEnvCRs)
	devEnvsToDelete := c.devEnvsToDelete(userFromAPI, devEnvCRs)

	// buildFilesToCreate := c.buildFileToCreate(buildFilesFromAPI, buildFileCMs)
	// buildFilesToDelete := c.buildFilesToDelete(buildFilesFromAPI, buildFileCMs)

	// buildersToCreate, buildFileToMap := c.buildersToCreate(builderFromAPI, builderCRs, buildFilesFromAPI)
	// buildersToDelete := c.buildersToDelete(builderFromAPI, builderCRs)

	/**
	* updating --------------------------------------------------------------
	**/

	//c.matchBuildFiles(buildFilesToCreate, buildFilesToDelete)
	//c.matchBuilders(buildersToCreate, buildersToDelete, buildFileToMap)
	c.matchDevEnvs(devEnvsToCreate, devEnvsToDelete, builderFromAPI, buildFilesFromAPI)

	//c.log.Infow("-- updating Dev Env K8s resources --")

	return nil
}
