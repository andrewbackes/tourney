#!/bin/bash

cat <<EOF > /tmp/data.json
{
	"testSeats" : 1,
	"carousel" : true,
	"rounds": 10,
	"contestants": [
		{
			"name": "one"
		},
		{
			"name": "two"
		}
	]
}
EOF

curl -X POST -d @/tmp/data.json http://localhost:8080/api/v2/tournaments
