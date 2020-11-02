package controllers

import (
	"strings"

	v1alpha1 "saas-controller.cloud-native-coding.dev/pkg/apis/cndecontroller/v1alpha1"
	"saas-controller.cloud-native-coding.dev/saasclient"
)

func (c *Controller) devEnvsToCreate(userFromAPI []*saasclient.DevEnvUser, devEnvCRs *v1alpha1.DevEnvList) []*saasclient.DevEnvUser {
	devEnvsToCreate := []*saasclient.DevEnvUser{}
	for _, user := range userFromAPI {
		found := false
		for _, de := range devEnvCRs.Items {
			if strings.ToLower(user.Name) == de.Name {
				found = true
				break
			}
		}
		if !found {
			devEnvsToCreate = append(devEnvsToCreate, user)
		}
	}
	c.log.Infow("DevEnv CRs to create:", "num", len(devEnvsToCreate))
	return devEnvsToCreate
}

func (c *Controller) devEnvsToDelete(userFromAPI []*saasclient.DevEnvUser, devEnvCRs *v1alpha1.DevEnvList) []v1alpha1.DevEnv {
	devEnvsToDelete := []v1alpha1.DevEnv{}
	for _, devEnvCR := range devEnvCRs.Items {
		found := false
		for _, user := range userFromAPI {
			//c.log.Infow("devEnvsToDelete is validating", "user.Name", user.Name, "devEnvCR.Name", devEnvCR.Name)
			if strings.ToLower(user.Name) == devEnvCR.Name {
				found = true
				break
			}
		}
		if !found {
			//c.log.Infow("devEnvsToDelete is adding", "devEnvCR.Name", devEnvCR.Name)
			devEnvsToDelete = append(devEnvsToDelete, devEnvCR)
		}
	}
	c.log.Infow("DevEnv CRs to delete:", "num", len(devEnvsToDelete))
	return devEnvsToDelete
}

// func (c *Controller) matchDevEnvs(devEnvsToCreate []*saasclient.DevEnvUser, devEnvsToDelete []*v1alpha1.DevEnv, builderFromAPI []*saasclient.Builder, buildFilesAll []*saasclient.BuildFile) error {
// 	for _, de := range devEnvsToCreate {

// 		builderName := ""
// 		for i := range buildFilesAll {
// 			if buildFilesAll[i].Name == de.BuildfileName {
// 				builderName = builderFromAPI[i].Name
// 			}
// 		}

// 		if _, err := c.k8sService.CreateDevEnv(de, c.keycloakHost, builderName); err != nil {
// 			c.log.Error(err, "can not create DevEnv instance")
// 		} else {
// 			c.log.Infow("created DevEnv CR", "name", de.Name)
// 		}
// 	}

// 	for _, de := range devEnvsToDelete {
// 		if err := c.k8sService.DeleteDevEnv(de); err != nil {
// 			c.log.Error(err, "can not delete DevEnv instance")
// 		} else {
// 			c.log.Infow("deleted DevEnv CR", "name", de.Name)
// 		}
// 	}
// 	return nil
// }

func (c *Controller) matchDevEnvs(devEnvsToCreate []*saasclient.DevEnvUser, devEnvsToDelete []v1alpha1.DevEnv, builderFromAPI []*saasclient.Builder, buildFilesAll []*saasclient.BuildFile) error {
	for _, de := range devEnvsToCreate {

		if de.BuildfileID != 0 {
			de.BuildFile.Name = de.Name
			de.Builder.Name = de.Name

			if _, err := c.k8sService.CreateBuildFile(&de.BuildFile); err != nil {
				c.log.Error(err, "can not create BuildFile instance")
			} else {
				c.log.Infow("created BuildFile CM", "name", de.BuildFile.Name)
			}

			if _, err := c.k8sService.CreateBuilder(&de.Builder, de.BuildFile.Name); err != nil {
				c.log.Error(err, "can not create Builder instance")
			} else {
				c.log.Infow("created Builder CR", "name", de.Builder.Name)
			}
		}

		if _, err := c.k8sService.CreateDevEnv(de, c.keycloakHost, de.Builder.Name); err != nil {
			c.log.Error(err, "can not create DevEnv instance")
		} else {
			c.log.Infow("created DevEnv CR", "name", de.Name)
		}
	}

	for _, de := range devEnvsToDelete {
		if err := c.k8sService.DeleteDevEnv(&de); err != nil {
			c.log.Error(err, "can not delete DevEnv instance")
		} else {
			c.log.Infow("deleted DevEnv CR", "name", de.Name)
		}

		if err := c.k8sService.DeleteBuilderByName(de.Name); err != nil {
			//if its not there it may be never created
			//c.log.Error(err, "can not delete Builder instance")
		} else {
			c.log.Infow("deleted Builder CR", "name", de.Name)
		}

		if err := c.k8sService.DeleteBuildFileByName(de.Name); err != nil {
			//if its not there it may be never created
			//c.log.Error(err, "can not delete BuildFile instance")
		} else {
			c.log.Infow("deleted BuildFile CM", "name", de.Name)
		}
	}
	return nil
}
