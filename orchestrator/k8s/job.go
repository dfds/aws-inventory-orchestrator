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
	AccountId          string
	JobName            string
	JobNamespace       string
	ServiceAccountName string
	CredsVolName       string
	CredsVolPath       string
	OutputVolName      string
	OutputVolPath      string
	AssumeName         string
	AssumeImage        string
	AssumeCmd          []string
	AssumeArgs         []string
	InventoryName      string
	InventoryImage     string
	InventoryCmd       []string
	InventoryArgs      []string
	UploadName         string
	UploadImage        string
	UploadCmd          []string
	UploadArgs         []string
}

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
							Name: assumeJobSpec.CredsVolName,
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: assumeJobSpec.OutputVolName,
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
					},
					InitContainers: []v1.Container{
						{
							Name:            assumeJobSpec.AssumeName,
							Image:           assumeJobSpec.AssumeImage,
							ImagePullPolicy: PullAlways,
							Command:         assumeJobSpec.AssumeCmd,
							Args:            assumeJobSpec.AssumeArgs,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      assumeJobSpec.CredsVolName,
									MountPath: assumeJobSpec.CredsVolPath,
								},
							},
						},
						{
							Name:            assumeJobSpec.InventoryName,
							Image:           assumeJobSpec.InventoryImage,
							ImagePullPolicy: PullAlways,
							Command:         assumeJobSpec.InventoryCmd,
							Args:            assumeJobSpec.InventoryArgs,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      assumeJobSpec.CredsVolName,
									MountPath: assumeJobSpec.CredsVolPath,
								},
								{
									Name:      assumeJobSpec.OutputVolName,
									MountPath: assumeJobSpec.OutputVolPath,
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "AWS_SHARED_CREDENTIALS_FILE",
									Value: assumeJobSpec.CredsVolPath + "/creds",
								},
								{
									// Clear the web indentity token
									// so the mounted AWS profile is used
									// instead of IRSA
									Name:  "AWS_WEB_IDENTITY_TOKEN_FILE",
									Value: "",
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:            assumeJobSpec.UploadName,
							Image:           assumeJobSpec.UploadImage,
							ImagePullPolicy: PullAlways,
							Command:         assumeJobSpec.UploadCmd,
							Args:            assumeJobSpec.UploadArgs,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      assumeJobSpec.OutputVolName,
									MountPath: assumeJobSpec.OutputVolPath,
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
	fmt.Printf("Job \"%s\" created successfully\n", job.ObjectMeta.Name)
}
