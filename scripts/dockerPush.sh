if [ ! -f build/info.json ]; then echo "ERROR Expected file: 'build/info.json'"; exit 1; fi

dockerRegistry=$(cat build/info.json | jq '.docker.registry' | xargs)
dockerImage=$(cat build/info.json | jq '.docker.image' | xargs)
dockerUsr=$(cat build/info.json | jq '.docker.usr' | xargs)
gitTagAtRev=$(cat build/info.json | jq '.git.tagAtRev' | xargs)
gitRevAtLatestTag=$(cat build/info.json | jq '.git.revAtLatestTag' | xargs)

if [ -z $dockerImage ]; then echo "ERROR Expected: 'docker.image'"; exit 1; fi
if [ -z $gitTagAtRev ]; then echo "ERROR Expected: 'git.tagAtRev'"; exit 1; fi

if [ ! -z $DOCKER_USER ]; then dockerUsr=$DOCKER_USER; fi
if [ ! -z $DOCKER_PASS ]; then dockerPwd=$DOCKER_PASS; fi

if [ -z $dockerUsr ]; then echo "ERROR Expected: 'docker.usr'"; exit 1; fi
if [ -z $dockerPwd ]; then echo "ERROR Expected: 'docker.pwd'"; exit 1; fi

docker login $dockerRegistry -u $dockerUsr -p $dockerPwd
if [ $? -ne 0 ]; then exit 1; fi

docker tag "$dockerImage:latest" "$dockerImage:$gitTagAtRev"
if [ $? -ne 0 ]; then exit 1; fi

docker push "$dockerImage:$gitTagAtRev"
if [ $? -ne 0 ]; then exit 1; fi

if [ "$gitTagAtRev" == "$gitRevAtLatestTag" ]; then 
  docker push "$dockerImage:latest"
  if [ $? -ne 0 ]; then exit 1; fi
fi