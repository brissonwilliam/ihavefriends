export DOCKER_BUILDKIT=1
export IHAVEFRIENDS_CONTAINER=ihavefriends_backend

echo -------------------------------------
echo BUILDING DOCKER IMAGE
echo -------------------------------------
docker build . -t ihavefriends_backend:latest

echo --------------------------------------
echo STOP AND REMOVE EXISTING CONTAINER
echo -------------------------------------
docker stop $IHAVEFRIENDS_CONTAINER
docker rm $IHAVEFRIENDS_CONTAINER

echo -------------------------------------
echo LAUNCHING CONTAINER
echo -------------------------------------
docker run --rm -d --env-file=$IHAVEFRIENDS_SECRETS --network host --name $IHAVEFRIENDS_CONTAINER ihavefriends_backend:latest

echo -------------------------------------
echo DONE
echo -------------------------------------