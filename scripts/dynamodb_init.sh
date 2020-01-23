aws2 dynamodb create-table --table-name Vocabulary --attribute-definitions AttributeName=vocab,AttributeType=S AttributeName=translation,AttributeType=S --key-schema AttributeName=vocab,KeyType=HASH AttributeName=translation,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=10,WriteCapacityUnits=5 --endpoint-url http://localhost:8000

aws2 dynamodb update-table \
--table-name Vocabulary \
--attribute-definitions AttributeName=globalKey,AttributeType=S AttributeName=lastTested,AttributeType=N \
--global-secondary-index-updates '[{"Create": {"IndexName": "globalKey-lastTested-index", "KeySchema": [{"AttributeName": "globalKey", "KeyType": "HASH"},{"AttributeName": "lastTested", "KeyType": "RANGE"}], "Projection": {"ProjectionType":"ALL"},"ProvisionedThroughput": {"ReadCapacityUnits": 1, "WriteCapacityUnits": 1}}}]' \
--endpoint-url http://localhost:8000
