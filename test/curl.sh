#!/bin/bash

curl -i -XPOST http://0.0.0.0:8080/login \ 
--header 'application/x-www-form-urlencoded'
--data 'idtoken=GOOGLE_CLIENT_ID_TOKEN'

curl -i -XGET http://0.0.0.0:8080/v1/user \
--cookie 'whiteboardsession=WHITEBOARD_COOKIE'

curl -i -XPOST http://0.0.0.0:8080/v1/user \ 
--header 'Content-Type: Application/json' \ 
--data '{"name":"John","email":"foobar@example.com"}' \
--cookie 'whiteboardsession=WHITEBOARD_COOKIE'

curl -i -XDELETE http://0.0.0.0:8080/v1/user/eddie
--cookie 'whiteboardsession=WHITEBOARD_COOKIE'

curl -i http://0.0.0.0:8080/v1/user 
--cookie 'whiteboardsession=WHITEBOARD_COOKIE'


