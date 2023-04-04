(cd ..; docker build -t ferlease:local .)
docker run --rm \
           -v $(pwd)/known_host:/opt/known_host \
           -v ~/.ssh/id_rsa:/opt/id_rsa \
           -v $(pwd)/config-docker.yml:/opt/config-docker.yml \
           --network host \
           ferlease:local teardown --config=/opt/config-docker.yml