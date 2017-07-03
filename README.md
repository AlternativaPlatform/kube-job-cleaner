# kube-job-cleaner

kube-job-cleaner - простое приложение, которое при запуске удаляет все завершенные Job из текущего Kubernetes кластера.

Актуально только для Kubernetes версии 1.5, в версии 1.6 у CronJob добавлены параметры `successfulJobsHistoryLimit` и `failedJobsHistoryLimit`, которые позволяют управлять историей выполненных задач.

## Локальная разработка

Для запуска через `go run main.go` или из IDE нужно выполнить:

```bash
$ mkdir -p $GOPATH/src/stash.alternativaplatform.com
$ ln -s /path/to/alternet/hosting/system/image/system/kube-job-cleaner $GOPATH/src/stash.alternativaplatform.com/
$ cd $GOPATH/src/stash.alternativaplatform.com/kube-job-cleaner
$ glide install --strip-vendor
```

## Ссылки

* [Kubernetes CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/)
