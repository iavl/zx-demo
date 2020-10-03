#!/usr/bin/env bash

curl http://localhost:8081/api/get_rsa

curl -X POST -d '{"data_list":[111,222,333],"pri_key":{"lambda":"2774103120","mu":"882170834"},"pub_key":{"n":"2774208617","g":"2774208618"}}' http://localhost:8081/api/compute?pri_key=xxxx


curl -X POST -d '{"data_list":[7585, 6116, 4376, 2269, 3049, 105, 5543, 8237, 7597, 9907, 1324, 5472, 9081, 670, 6209, 2775, 7405, 6923, 2574, 2759, 3899, 841, 1817, 2169, 8294, 9673, 1037, 6285, 1665, 4770],"pri_key":{"lambda":"3133326168","mu":"1377571274"},"pub_key":{"n":"3133438163","g":"3133438164"}}' http://localhost:8081/api/compute?pri_key=xxxx