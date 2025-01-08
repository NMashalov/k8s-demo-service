package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// getKubernetesClient creates a new Kubernetes clientset from the current context.
func getKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// createService creates a new Kubernetes service with the given name.
func createService(clientset *kubernetes.Clientset, serviceName string) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: corev1.ServiceSpec{
			Type:     "ClusterIP",
			Selector: map[string]string{"app": serviceName}, // Assume pods are labeled with 'app' key
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.IntOrString{IntVal: 8080},
				},
			},
		},
	}

	_, err := clientset.CoreV1().Services(metav1.NamespaceDefault).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Service %s created\n", serviceName)
	return nil
}

func main() {
	router := gin.Default()

	router.GET("/service/:serviceName", func(c *gin.Context) {
		serviceName := c.Param("serviceName")

		clientset, err := getKubernetesClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Kubernetes cluster"})
			return
		}

		err = createService(clientset, serviceName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create service %s: %v", serviceName, err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Service %s created successfully", serviceName),
		})
	})

	log.Fatal(router.Run(":8080"))
}
