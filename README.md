# kube-job-cleaner

kube-job-cleaner - простое приложение, которое при запуске удаляет все завершенные Job из текущего kubernetes кластера.

## Локальная разработка

Для запуска через `go run main.go` или из IDE нужно выполнить:

```bash
$ mkdir -p $GOPATH/src/stash.alternativaplatform.com
$ ln -s /path/to/alternet/hosting/system/image/system/kube-job-cleaner $GOPATH/src/stash.alternativaplatform.com/
$ cd $GOPATH/src/stash.alternativaplatform.com/kube-job-cleaner
$ glide install --strip-vendor
```
