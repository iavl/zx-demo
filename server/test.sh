#!/usr/bin/env bash

curl http://localhost:8081/api/get_rsa


curl -X POST -d '{"data_list":[112,333,444,555],"pri_key":{"lambda":"2774103120","mu":"882170834"},"pub_key":{"n":"2774208617","g":"2774208618"}	}' http://localhost:8081/api/compute


curl -X POST -d '{"data_list":[111,222,333],"pri_key":{"lambda":"3070826472","mu":"1975704042"},"pub_key":{"n":"3070937591","g":"3070937592"}}' http://localhost:8081/api/compute?pri_key=xxxx