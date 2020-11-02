package saasclient_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"net/http"
	"net/http/httptest"

	. "saas-controller.cloud-native-coding.dev/saasclient"
)

var _ = Describe("Saasclient", func() {

	var (
		log *zap.SugaredLogger
	)

	BeforeEach(func() {
		zapLog, err := zap.NewDevelopment()
		if err != nil {
			panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
		}
		log = zapLog.Sugar()
	})

	Describe("Cluster Endpoint", func() {
		Context("get all DevEnvUser", func() {
			It("should recieve 2 DevEnvs", func() {
				server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					Expect(req.URL.String()).To(Equal("/clusters/barCluster?apiKey=fooKey"))

					rw.Header().Set("Content-Type", "application/json")
					rw.WriteHeader(http.StatusOK)

					byt := []byte(`{"apiKey":"ffa3c633-1500-4943-aed0-9fe399163d62","name":"aa","devenvusers":[{"id":1,"buildfileId":1,"links":{"self":"/devenvusers/1","buildfile":"/buildfiles/1"}}]}`)
					rw.Write(byt)
				}))
				defer server.Close()

				saaSClient := NewSaaSClient(CreateClient(), server.URL, "fooKey", "barCluster", log)
				cluster, err := saaSClient.GetDevEnvUsersForCluster()

				Expect(err).NotTo(HaveOccurred())

				c := &Cluster{
					Name:   "aa",
					APIKey: "ffa3c633-1500-4943-aed0-9fe399163d62",
					DevEnvUsers: []ClusterDevEnv{
						{
							ID:          1,
							BuildfileID: 1,
							Links: ClusterLinks{
								Self:      "/devenvusers/1",
								Buildfile: "/buildfiles/1",
							},
						},
					},
				}
				Expect(cluster).To(Equal(c))
			})
		})
	})

	Describe("DevEnvUser Endpoint", func() {
		Context("get DevEnvUser", func() {
			It("should recieve 1 DevEnv", func() {
				server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					Expect(req.URL.String()).To(Equal("/a/b/c?apiKey=fooKey"))

					rw.Header().Set("Content-Type", "application/json")
					rw.WriteHeader(http.StatusOK)

					byt := []byte(`{"name":"barUser","buildfileID":1,"DeleteVolume":false,
					"ClusterRoleName":"ClusterRoleName", "RoleName":"RoleName","DevEnvImage":"DevEnvImage",
					"ContainerVolumeSize":"ContainerVolumeSize","HomeVolumeSize":"HomeVolumeSize", "Email":"Email"}`)
					rw.Write(byt)
				}))
				defer server.Close()

				saaSClient := NewSaaSClient(CreateClient(), server.URL, "fooKey", "barCluster", log)
				user, err := saaSClient.GetDevEnvUser("/a/b/c")

				Expect(err).NotTo(HaveOccurred())

				d := &DevEnvUser{
					Name:                "barUser",
					BuildfileID:         1,
					DeleteVolume:        false,
					ClusterRoleName:     "ClusterRoleName",
					RoleName:            "RoleName",
					DevEnvImage:         "DevEnvImage",
					ContainerVolumeSize: "ContainerVolumeSize",
					HomeVolumeSize:      "HomeVolumeSize",
					Email:               "Email",
				}
				Expect(user).To(Equal(d))
			})
		})
	})

	Describe("Metrics Endpoint", func() {
		Context("PUT Metrics", func() {
			It("should put right metrics JSON", func() {
				server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

					buf := new(bytes.Buffer)
					buf.ReadFrom(req.Body)

					Expect(req.URL.String()).To(Equal("/clusters/barCluster/devenvusers/metrics?apiKey=fooKey"))
					Expect(buf.String()).To(Equal("[{\"devenvuser\":\"whatever/1\",\"status\":\"running\",\"cpu\":\"1\",\"memory\":\"2\"}]"))

					rw.Header().Set("Content-Type", "application/json")
					rw.WriteHeader(http.StatusNoContent)

					byt := []byte(`ok`)
					rw.Write(byt)
				}))
				defer server.Close()

				saaSClient := NewSaaSClient(CreateClient(), server.URL, "fooKey", "barCluster", log)
				clusterStatus := ClusterStatus{
					DevEnvStatus{
						CPU:        "1",
						Devenvuser: "whatever/1",
						Memory:     "2",
						Status:     "running",
					},
				}
				err := saaSClient.PutClusterStatus(&clusterStatus)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

})
