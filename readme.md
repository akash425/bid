to run this project we need to run following commands:
<!-- Assuming go is already installed in the system -->

<!-- initialize project -->
go mod init bid

<!-- install dependencies -->
go get github.com/go-redis/redis/v8
go get github.com/gorilla/mux

<!-- run this project -->
go run main.go

<!-- Assuming redis is installed in the system -->
<!-- insert data into redis -->
hset id_123 impression "http://127.0.0.1/impression/123" click "http://127.0.0.1/click/123" video_url "http://127.0.0.1/video/123" video_start "http://127.0.0.1/start/123" video_end "http://127.0.0.1/end/123"

set id_123_js "<a href='{click}'><img src='{impression}'></a>"

set id_123_xml "<Response><impression><![CDATA[{impression}]]</impression><click><![CDATA[{click}]]</click><video_url><![CDATA[{video_url}]]</video_url><video_start><![CDATA[{video_start}]]</video_start><video_end><![CDATA[{video_end}]]</video_end></Response>"



<!-- curl to see the results -->
<!-- type can be 1 or 2 -->
curl -X POST "http://localhost:5004/bid" \
  -H "Accept: */*" \
  -H "User-Agent: Thunder Client (https://www.thunderclient.com)" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test1234" \
  -d '{
  "id": "id_123",
  "width": 600,
  "height": 328,
  "banner": {
    "type": 1
  }
}'

