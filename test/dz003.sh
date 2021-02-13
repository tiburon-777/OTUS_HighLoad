#/bin/sh

docker run --rm -v /root/scripts:/scripts williamyeh/wrk -t1 -c10 -d5m --timeout 30s http://localhost:8080/search -s /scripts/post.lua -- debug true
