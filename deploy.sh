IMAGE=registry.gitlab.com/kamackay/filer:$1

time docker build . -t "$IMAGE" && \
    docker push "$IMAGE" && \
    kubectl --context do-nyc3-keithmackay-cluster -n file-server \
      set image statefulset/file-server server=$IMAGE && \
    kubectl --context do-nyc3-keithmackay-cluster -n file-server rollout restart statefulset file-server

sleep 1
ATTEMPTS=0
ROLLOUT_STATUS_CMD="kubectl --context do-nyc3-keithmackay-cluster rollout status statefulset/file-server -n file-server"
until $ROLLOUT_STATUS_CMD || [ $ATTEMPTS -eq 60 ]; do
  $ROLLOUT_STATUS_CMD
  ATTEMPTS=$((ATTEMPTS + 1))
  sleep 1
done

ECHO "Successfully deployed" $1