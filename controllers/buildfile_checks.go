package controllers

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	"saas-controller.cloud-native-coding.dev/saasclient"
)

func (c *Controller) buildFileToCreate(buildFileFromAPI []*saasclient.BuildFile, buildFileCRs *corev1.ConfigMapList) []*saasclient.BuildFile {
	buildFilesToCreate := []*saasclient.BuildFile{}
	for _, buildFile := range buildFileFromAPI {
		found := false
		for _, de := range buildFileCRs.Items {
			if strings.ToLower(buildFile.Name) == de.Name {
				found = true
				break
			}
		}
		if !found {
			buildFilesToCreate = append(buildFilesToCreate, buildFile)
		}
	}
	c.log.Infow("BuildFile CMs to create:", "num", len(buildFilesToCreate))
	return buildFilesToCreate
}

func (c *Controller) buildFilesToDelete(buildFileFromAPI []*saasclient.BuildFile, buildFileCRs *corev1.ConfigMapList) []*corev1.ConfigMap {
	buildFilesToDelete := []*corev1.ConfigMap{}
	for _, bf := range buildFileCRs.Items {
		found := false
		for _, user := range buildFileFromAPI {
			if strings.ToLower(user.Name) == bf.Name {
				found = true
				break
			}
		}
		if !found {
			buildFilesToDelete = append(buildFilesToDelete, &bf)
		}
	}
	c.log.Infow("BuildFile CMs to delete:", "num", len(buildFilesToDelete))
	return buildFilesToDelete
}

func (c *Controller) matchBuildFiles(buildFilesToCreate []*saasclient.BuildFile, buildFilesToDelete []*corev1.ConfigMap) error {
	for _, bf := range buildFilesToCreate {
		if _, err := c.k8sService.CreateBuildFile(bf); err != nil {
			c.log.Error(err, "can not create BuildFile instance")
		} else {
			c.log.Infow("created BuildFile CM", "name", bf.Name)
		}
	}

	for _, bf := range buildFilesToDelete {
		if err := c.k8sService.DeleteBuildFile(bf); err != nil {
			c.log.Error(err, "can not delete BuildFile instance")
		} else {
			c.log.Infow("deleted BuildFile CM", "name", bf.Name)
		}
	}
	return nil
}
