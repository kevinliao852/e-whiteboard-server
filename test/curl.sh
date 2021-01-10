#!/bin/bash
curl -i -XPOST 0.0.0.0:8080/v1/user --header 'Content-Type: Application/json' --data '{"name":"eddie","password":"123","email":"123@gmail.com"}'

curl -i -XDELETE 0.0.0.0:8080/v1/user/eddie

curl -i 0.0.0.0:8080/v1/user 


echo "done"

