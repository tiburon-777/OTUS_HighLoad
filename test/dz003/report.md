Встроен запуск prometeus, grafana и т.д.
Графана доступна на http://localhost:3001/

$ sudo docker run --rm -v /root/scripts:/scripts williamyeh/wrk -t1 -c10 -d5m --timeout 30s http://192.168.1.66:8080/search -s /scripts/post.lua -- debug true
Running 5m test @ http://192.168.1.66:8080/search
1 threads and 10 connections
Thread Stats   Avg      Stdev     Max   +/- Stdev
Latency     5.83s     1.74s   13.29s    77.58%
Req/Sec     3.57      4.74    30.00     82.50%
511 requests in 5.00m, 18.69MB read
Requests/sec:      1.70
Transfer/sec:     63.79KB

$ sudo make app-reload
...
...
...
$ sudo docker run --rm -v /root/scripts:/scripts williamyeh/wrk -t1 -c10 -d5m --timeout 30s http://192.168.1.66:8080/search -s /scripts/post.lua -- debug true
Running 5m test @ http://192.168.1.66:8080/search
1 threads and 10 connections
Thread Stats   Avg      Stdev     Max   +/- Stdev
Latency    32.12ms   31.59ms 561.62ms   90.65%
Req/Sec   373.34    147.37   610.00     68.45%
110834 requests in 5.00m, 433.90MB read
Non-2xx or 3xx responses: 110834
Requests/sec:    369.39
Transfer/sec:      1.45MB