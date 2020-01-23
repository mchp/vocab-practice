echo "building server binary..."
go build -ldflags "-X 'main.version=`date`'" server.go
echo "finished building server binary"

cd ui/add && npm run publish && cd ../..
cd ui/quiz && npm run publish && cd ../..

echo "clear previous version..."
aws2 s3 rm  s3://elasticbeanstalk-us-west-1-724990643256/vocab-practice-binary --recursive

echo "uploading server binary..."
aws2 s3 cp server s3://elasticbeanstalk-us-west-1-724990643256/vocab-practice-binary/server --acl public-read
echo "uploading public files..."
aws2 s3 cp public s3://elasticbeanstalk-us-west-1-724990643256/vocab-practice-binary/public --recursive --acl public-read

rm server
