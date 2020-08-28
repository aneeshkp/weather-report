#! /bin/sh
cd ../../../
#make docker-build docker-push IMG=quay.io/aneeshkp/weather-report:latest
#make install
make deploy IMG=quay.io/aneeshkp/weather-report:latest