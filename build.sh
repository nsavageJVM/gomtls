#!/bin/bash
echo "cleaning up"
rm client/client
rm server/server
echo "building"
cd server
go build .
cd ..
cd client
go build .
