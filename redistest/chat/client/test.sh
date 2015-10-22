#!/bin/bash -x

for i in {1..5}
do
sh -c "go run client.go" & 
sleep 2
done
