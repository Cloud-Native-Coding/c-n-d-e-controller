package controllers

import (
	"strings"

	v1alpha1 "saas-controller.cloud-native-coding.dev/pkg/apis/cndecontroller/v1alpha1"
	"saas-controller.cloud-native-coding.dev/saasclient"
)

func (c *Controller) buildersToCreate(builderFromAPI []*saasclient.Builder, builderCRs *v1alpha1.BuilderList, buildFilesAll []*saasclient.BuildFile) (buildersToCreate []*saasclient.Builder, buildFiles []*saasclient.BuildFile) {
	for i, builder := range builderFromAPI {
		found := false
		for _, de := range builderCRs.Items {
			if strings.ToLower(builder.Name) == de.Name {
				found = true
				break
			}
		}
		if !found {
			buildersToCreate = append(buildersToCreate, builder)
			buildFiles = append(buildFiles, buildFilesAll[i])
		}
	}
	c.log.Infow("Builder CMs to create:", "num", len(buildersToCreate))
	return buildersToCreate, buildFiles
}

func (c *Controller) buildersToDelete(builderFromAPI []*saasclient.Builder, builderCRs *v1alpha1.BuilderList) (buildersToDelete []*v1alpha1.Builder) {
	for _, bf := range builderCRs.Items {
		found := false
		for _, user := range builderFromAPI {
			if strings.ToLower(user.Name) == bf.Name {
				found = true
				break
			}
		}
		if !found {
			buildersToDelete = append(buildersToDelete, &bf)
		}
	}
	c.log.Infow("Builder CMs to delete:", "num", len(buildersToDelete))
	return buildersToDelete
}

func (c *Controller) matchBuilders(buildersToCreate []*saasclient.Builder, buildersToDelete []*v1alpha1.Builder, buildFileToMap []*saasclient.BuildFile) error {
	for i, b := range buildersToCreate {
		if _, err := c.k8sService.CreateBuilder(b, buildFileToMap[i].Name); err != nil {
			c.log.Error(err, "can not create Builder instance")
		} else {
			c.log.Infow("created Builder CR", "name", b.Name)
		}
	}

	for _, bf := range buildersToDelete {
		if err := c.k8sService.DeleteBuilder(bf); err != nil {
			c.log.Error(err, "can not delete Builder instance")
		} else {
			c.log.Infow("deleted Builder CR", "name", bf.Name)
		}
	}
	return nil
}
