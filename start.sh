#! /bin/bash
git pull
cd api
go build
cd ..
nohup ./api/api &