#/bin/sh

docker run --rm -v /root/scripts:/scripts williamyeh/wrk -t1 -c1  -d1m  --timeout 30s http://lab.tiburon.su:8080/search -s /scripts/post.lua -- debug true

docker run --rm -v /root/scripts:/scripts williamyeh/wrk -t1 -c10 -d10m --timeout 30s http://lab.tiburon.su:8080/search -s /scripts/post.lua -- debug true

docker run --rm -v /root/scripts:/scripts williamyeh/wrk -t1 -c20 -d10m --timeout 30s http://lab.tiburon.su:8080/search -s /scripts/post.lua -- debug true

docker run --rm -v /root/scripts:/scripts williamyeh/wrk -t1 -c30 -d10m --timeout 30s http://lab.tiburon.su:8080/search -s /scripts/post.lua -- debug true

docker run --rm -v /root/scripts:/scripts williamyeh/wrk -t1 -c40 -d10m --timeout 30s http://lab.tiburon.su:8080/search -s /scripts/post.lua -- debug true

docker run --rm -v /root/scripts:/scripts williamyeh/wrk -t1 -c50 -d10m --timeout 30s http://lab.tiburon.su:8080/search -s /scripts/post.lua -- debug true
