dockerImage="tmp/cleaner-ttl.build"

mkdir -p ./build 

echo "===== building docker to use for compilation"
docker build -t $dockerImage -f Dockerfile.build .

echo "===== compiling using built docker image"
docker run -v $(pwd)/build:/app/build $dockerImage