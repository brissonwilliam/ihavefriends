#!bin/bash

mkdir newapp
tar -xzf build.tar.gz -C newapp
mv newapp/frontend-build newapp/frontend
mv newapp/backend-build newapp/backend
chmod 700 -r newapp/*
mv app app.old
mv newapp app
pkill -f backend

nohup ~/app/backend < /dev/null > backend_logs.log 2>&1 &
rm -rf app.old