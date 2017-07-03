package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/apis/batch/v1"
	"k8s.io/client-go/1.5/pkg/labels"
	"k8s.io/client-go/1.5/rest"
	"k8s.io/client-go/1.5/tools/clientcmd"
)

var dryRun = flag.Bool("dry-run", false, "Не удалять Job'ы и Pod'ы, только залогировать, что будет удалено")

type JobCleaner struct {
	DryRun bool
	Client *kubernetes.Clientset
}

func init() {
	flag.Parse()

	if *dryRun {
		log.Println("[!] Kube-job-cleaner запущен в режиме dry-run")
	}
}

func main() {
	clientset, err := GetKubeClientset()
	if err != nil {
		log.Fatal(err)
	}

	cleaner := JobCleaner{
		DryRun: *dryRun,
		Client: clientset,
	}

	if err := cleaner.DeleteSucceededJobs(); err != nil {
		log.Fatal(err)
	}
}

// Запускае удаление всех Job'ов с status.succeeded == 1, а так же их Pod'ов, из всех
// нэймспэйсов kubernetes кластера.
//
// Кластер выбирается из serviceaccount или current-context файла ~/.kube/config
func (j *JobCleaner) DeleteSucceededJobs() error {
	namespaces, err := j.Client.Core().Namespaces().List(api.ListOptions{})
	if err != nil {
		return err
	}
	if len(namespaces.Items) == 0 {
		return errors.New("Не найдено ни одного Namespace в кластере")
	}
	for _, namespace := range namespaces.Items {
		nsJobsClient := j.Client.Batch().Jobs(namespace.Name)
		jobs, err := nsJobsClient.List(api.ListOptions{})
		if err != nil {
			return err
		}
		if len(jobs.Items) == 0 {
			log.Printf("Список Job'ов в Namespace %s пуст, нечего удалять", namespace.Name)
			continue
		}
		for _, job := range jobs.Items {
			if job.Status.Succeeded != 1 {
				log.Printf("Job %s не завершен (status.succeeded != 1). Удаление невозможно", job.Name)
				continue
			}
			if err := j.DeleteJob(job); err != nil {
				fmt.Println(err)
				continue
			}
			if err := j.DeleteJobPods(job); err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

// Удаление конкретного джоба
func (j *JobCleaner) DeleteJob(job v1.Job) error {
	log.Printf("Job %s успешно завершен. Удаляем его...", job.Name)

	nsJobsClient := j.Client.Batch().Jobs(job.Namespace)

	if !j.DryRun {
		if err := nsJobsClient.Delete(job.Name, &api.DeleteOptions{}); err != nil {
			return fmt.Errorf("Не удалось удалить Job %s: %v", job.Name, err)
		}
	}

	return nil
}

// Удаление всех Pod'ов, которые создал Job
func (j *JobCleaner) DeleteJobPods(job v1.Job) error {
	log.Printf("Получаем список Pod'ов Job'а %s...", job.Name)

	nsPodsClinet := j.Client.Pods(job.Namespace)

	labelMap := labels.Set(map[string]string{"job-name": job.Name})
	podsList, err := nsPodsClinet.List(api.ListOptions{LabelSelector: labelMap.AsSelector()})
	if err != nil {
		return fmt.Errorf("Не удалось получить список Pod'ов Job'а %s: %v", job.Name, err)
	}
	if len(podsList.Items) == 0 {
		log.Println("Список Pod'ов пуст, нечего удалять")
		return nil
	}

	for _, pod := range podsList.Items {
		log.Printf("Удаляем Pod %s...", pod.Name)
		if *dryRun {
			continue
		}
		if err := nsPodsClinet.Delete(pod.Name, &api.DeleteOptions{}); err != nil {
			return fmt.Errorf("Не удалось удалить Pod %s: %v", pod.Name, err)
		}
	}

	return nil
}

// Создать клиент для kube-apiserver
// Сначала пробуем из serviceaccount
// Если не получилось пробуем kubeconfig из файла ~/.kube/config
func GetKubeClientset() (*kubernetes.Clientset, error) {
	var config *rest.Config

	log.Println("Create kube-apiserver client...")
	log.Println("Create kube-apiserver InCluster client with ServiceAccount...")

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("InCluster kube-apiserver creation failed. Try to create client from kubeconfig... %s", err)
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
		if err != nil {
			return nil, err
		}
	}

	return kubernetes.NewForConfig(config)
}
