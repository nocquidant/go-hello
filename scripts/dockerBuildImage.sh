if [ ! -f build/build-info.json ]; then echo "ERROR Expected file: 'build/build-info.json'"; exit 1; fi

dockerImage=$(cat build/build-info.json | jq '.docker.image' | xargs)
if [ -z $dockerImage ]; then echo "ERROR Expected: 'docker.image'"; exit 1; fi

docker build -t $dockerImage .
