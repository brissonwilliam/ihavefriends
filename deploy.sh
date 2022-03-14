#!/bin/bash

# This script builds and deploys the app in PROD.
# It expects the local machine has the shh key to the server

if [ -z "$IHAVEFRIENDS_PROD_HOST"]; then
  echo -e "\033[31mIHAVEFRIENDS_PROD_HOST env var is empty\033[0m"
  exit 1
fi

echo -e "\033[32mBuilding Backend for linux\033[0m"
cd backend
export GOARCH=amd64
export GOOS=linux
go build -o ../backend-build
chmod +x ../backend-build
cd ..
echo -e "Done!\n"

echo -e "\033[32mBuilding Frontend\033[0m"
rm -rf frontend-build
cd frontend
rm -rf build
npm run build -- --production
chmod +xr build
mv build frontend-build
mv frontend-build ..
cd ..
echo -e "DONE!\n"

echo -e "\033[32mCreating build zip\033[0m"
rm build.tar.gz
tar -cvzf build.tar.gz frontend-build backend-build
rm backend-build
rm -rf frontend-build
echo -e "DONE!\n"

echo -e "\033[32mCopying file on remote\033[0m"
scp ./build.tar.gz $IHAVEFRIENDS_PROD_HOST:~
scp ./deploy_localprod.sh $IHAVEFRIENDS_PROD_HOST:~
echo -e "DONE!\n"

echo -e "\033[32mDeploying on remote\033[0m"
ssh $IHAVEFRIENDS_PROD_HOST "source ~/.bashrc;sh ~/deploy_localprod.sh"
echo -e "DONE!\n"

