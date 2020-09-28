#!/usr/bin/env bash

curl http://localhost:8081/api/get_rsa


curl -X POST -d '{"data_list":[112,333,444,555],"pri_key":{"lambda":"2774103120","mu":"882170834"},"pub_key":{"n":"2774208617","g":"2774208618"}	}' http://localhost:8081/api/compute


curl -X POST -d '{"data_list":[111,333,444],"pri_key":{"lambda":"2761845840","mu":"1561715058"},"pub_key":{"n":"2761950953","g":"2761950954"}}' http://localhost:8081/api/compute?pri_key=xxxx
