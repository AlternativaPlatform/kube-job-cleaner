FROM scratch

COPY kube-job-cleaner /kube-job-cleaner

ENTRYPOINT ["/kube-job-cleaner"]
