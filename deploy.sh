docker build . -t registry.gitlab.com/kamackay/filer:$1 && \
    docker push registry.gitlab.com/kamackay/filer:$1 && \
    kubectl --context do-nyc3-keithmackay-cluster -n file-server \
      set image statefulset/file-server server=registry.gitlab.com/kamackay/filer:$1 && \
    kubectl --context do-nyc3-keithmackay-cluster -n file-server rollout restart statefulset file-server
