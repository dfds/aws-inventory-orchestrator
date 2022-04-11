package k8s

import (
	"context"
	"fmt"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	// PullAlways means that kubelet always attempts to pull the latest image. Container will fail If the pull fails.
	PullAlways v1.PullPolicy = "Always"
	// PullNever means that kubelet never pulls an image, but only uses a local image. Container will fail if the image isn't present
	PullNever v1.PullPolicy = "Never"
	// PullIfNotPresent means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
	PullIfNotPresent v1.PullPolicy = "IfNotPresent"
)

type AssumeJobSpec struct {
	JobName            string
	JobNamespace       string
	ServiceAccountName string
	InitName           string
	InitImage          string
	InitCmd            []string
	InitArgs           []string
	ContainerName      string
	ContainerImage     string
	ContainerCmd       []string
	ContainerArgs      []string
}

// func CreateJob(jobName *string, jobNamespace *string, image *string, cmd *string, roleArn *string) {
func CreateJob(assumeJobSpec *AssumeJobSpec) {

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
	jobs := clientset.BatchV1().Jobs(assumeJobSpec.JobNamespace)
	var backOffLimit int32 = 0

	// create new job spec
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", assumeJobSpec.JobName),
			Namespace:    assumeJobSpec.JobNamespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "aws-creds",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
					},
					InitContainers: []v1.Container{
						{
							Name:            assumeJobSpec.InitName,
							Image:           assumeJobSpec.InitImage,
							ImagePullPolicy: PullAlways,
							Command:         assumeJobSpec.InitCmd,
							Args:            assumeJobSpec.InitArgs,
							// Env: []v1.EnvVar{
							// 	{
							// 		Name:  "AWS_ROLE_SESSION_NAME",
							// 		Value: "inventory",
							// 	},
							// },
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "aws-creds",
									MountPath: "/aws",
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:            assumeJobSpec.ContainerName,
							Image:           assumeJobSpec.ContainerImage,
							ImagePullPolicy: PullAlways,
							Command:         assumeJobSpec.ContainerCmd,
							Args:            assumeJobSpec.ContainerArgs,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "aws-creds",
									MountPath: "/aws",
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "AWS_SHARED_CREDENTIALS_FILE",
									Value: "/aws/creds",
								},
								{
									Name:  "AWS_WEB_IDENTITY_TOKEN_FILE",
									Value: "",
								},
							},
						},
					},
					RestartPolicy:      v1.RestartPolicyNever,
					ServiceAccountName: assumeJobSpec.ServiceAccountName,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	job, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create K8s job. Error: %v\n", err)
	}

	//print job details
	log.Printf("Job \"%s\" created successfully\n", job.ObjectMeta.Name)
}
