/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	controller "saas-controller.cloud-native-coding.dev/controllers"
	"saas-controller.cloud-native-coding.dev/metrics"
	"saas-controller.cloud-native-coding.dev/saasclient"
)

var (
	APIKey        = os.Getenv("CNDE_API_KEY")
	ClusterName   = os.Getenv("CNDE_CLUSTER_NAME")
	CndeURL       = os.Getenv("CNDE_URL")
	KeycloakHost  = os.Getenv("CNDE_KEYCLOAK_HOST")
	CndeNamespace = os.Getenv("CNDE_NS")
)

func main() {

	zapLog, err := zap.NewDevelopment()
	defer zapLog.Sync() // flushes buffer, if any
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log := zapLog.Sugar() //zapr.NewLogger(zapLog)

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err, "is this running inside of a K8s Cluster?")
		panic(err.Error())
	}

	for {
		// // get pods in all the namespaces by omitting namespace
		// // Or specify namespace to get pods in particular namespace
		// pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		// if err != nil {
		// 	panic(err.Error())
		// }
		// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// // Examples for error handling:
		// // - Use helper functions e.g. errors.IsNotFound()
		// // - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		// _, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
		// if errors.IsNotFound(err) {
		// 	fmt.Printf("Pod example-xxxxx not found in default namespace\n")
		// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		// } else if err != nil {
		// 	panic(err.Error())
		// } else {
		// 	fmt.Printf("Found example-xxxxx pod in default namespace\n")
		// }

		time.Sleep(10 * time.Second)

		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		apiClient := saasclient.NewSaaSClient(saasclient.CreateClient(), CndeURL, APIKey, ClusterName, log)
		controller := controller.NewController(apiClient, config, log)
		metrics := metrics.NewMetrics(apiClient, config, log)

		go func() {
			sig := <-sigs
			fmt.Println("received signal: ", sig)
			done <- true
		}()

		ticker := time.NewTicker(time.Second * 10)
		go func() {
			for {
				<-ticker.C
				log.Infow("--- Running control loop ---")

				if err := controller.Reconcile(); err != nil {
					log.Infow("Failed Controller Reconsile Loop - doing nothing")
				}

				if err := metrics.Calculate(); err != nil {
					log.Infow("Failed Metrics Calculate Loop - doing nothing")
				}

				// progress:
				// calculate and generate K8s CRs
				// if CR is already present, check equality
				//   if equal then do nothing
				//   if not equal delete CR
				// if CR is not present then create

			}
		}()

		/////////////////
		// Exit "handler"

		log.Infow("awaiting signal")
		<-done
		log.Infow("cleaning up and exiting")
		ticker.Stop()
	}
}
