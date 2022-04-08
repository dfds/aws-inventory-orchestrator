package k8s

import (
	"context"
	"log"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func CreateJob(jobName *string, jobNamespace *string, image *string, cmd *string) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// get jobs collection in the inventory namespace
	jobs := clientset.BatchV1().Jobs("inventory")
	var backOffLimit int32 = 0

	// create new job spec
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: *jobName,
			Namespace:    *jobNamespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    *jobName,
							Image:   *image,
							Command: strings.Split(*cmd, " "),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err = jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create K8s job. Error: %v\n", err)
	}

	//print job details
	log.Println("Created K8s job successfully")
}
