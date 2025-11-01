# Filo a File Discovery API Server
This project is a Go server that uses Fiber, Huma, and GORM (with SQLite) to provide a CRUD API for tracking file metadata.

## Running the Server

### Run the application:

```bash
make run_server
```

## Using the API

### API Documentation

To access the API documentation, open your browser and navigate to: http://localhost:8000/docs

You can explore all endpoints and even test them directly from this web page.

### Example curl Commands

Here are some curl commands you can use to test from your terminal:

1. Register a File (POST /files)
This will create a new file record.

```bash
curl -X POST 'http://localhost:8000/files' \
-H 'Content-Type: application/json' \
-d '{
    "directory_path": "/var/log/",
    "file_type": "filo.log",
    "size": 123456,
    "checksum": "a1b2c3d4e5f6..."
}'
```


2. Register the Same File (POST /files)
If you run this command, it will find the existing record and update its updatedAt timestamp (and size/checksum if they changed).

```bash
curl -X POST 'http://localhost:8000/files' \
-H 'Content-Type: application/json' \
-d '{
    "directory_path": "/var/log/",
    "file_type": "filo.log",
    "size": 789000,
    "checksum": "f6e5d4c3b2a1..."
}'
```

3. List All Files (GET /files)

```bash
curl 'http://localhost:8000/files'
```

4. List Files from a specific host (GET /files?directory_path=...

```bash
curl 'http://localhost:8000/files?directory_path=agent-001'
```

5. Get File by ID (GET /files/{id})
(Assuming the first file you created has ID: 1)

```bash
curl 'http://localhost:8000/files/1'
```

6. Delete a File (DELETE /files)
(This uses query parameters, which is safer for agents)

```bash
curl -X DELETE 'http://localhost:8000/files?directory_path=%2Fvar%2Flog%2F&filename=filo.log'
```

## Load testing
To load test the application, you can use [k6](github.com/grafana/k6)


```bash
k6 run create.js

         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/

     execution: local
        script: create.js
        output: -

     scenarios: (100.00%) 1 scenario, 100 max VUs, 1m15s max duration (incl. graceful stop):
              * default: Up to 100 looping VUs for 45s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


     ✗ created successfully (status 201)
      ↳  0% — ✓ 0 / ✗ 91550

     checks.........................: 0.00% 0 out of 91550
     data_received..................: 78 MB 1.7 MB/s
     data_sent......................: 24 MB 530 kB/s
     http_req_blocked...............: avg=1.15µs  min=0s       med=1µs    max=1.45ms p(90)=1µs     p(95)=2µs
     http_req_connecting............: avg=202ns   min=0s       med=0s     max=715µs  p(90)=0s      p(95)=0s
     http_req_duration..............: avg=41.28ms min=348µs    med=1.52ms max=5.13s  p(90)=36.89ms p(95)=142.58ms
       { expected_response:true }...: avg=41.22ms min=348µs    med=1.52ms max=4.97s  p(90)=36.89ms p(95)=142.58ms
     http_req_failed................: 0.00% 1 out of 91550
     http_req_receiving.............: avg=12.15µs min=7µs      med=10µs   max=2.1ms  p(90)=16µs    p(95)=20µs
     http_req_sending...............: avg=3.22µs  min=1µs      med=2µs    max=1.64ms p(90)=5µs     p(95)=6µs
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s     max=0s     p(90)=0s      p(95)=0s
     http_req_waiting...............: avg=41.26ms min=338µs    med=1.5ms  max=5.13s  p(90)=36.88ms p(95)=142.56ms
     http_reqs......................: 91550 2030.195786/s
     iteration_duration.............: avg=41.27ms min=370.83µs med=1.54ms max=5.12s  p(90)=36.92ms p(95)=142.54ms
     iterations.....................: 91550 2030.195786/s
     vus............................: 4     min=4          max=100
     vus_max........................: 100   min=100        max=100


running (0m45.1s), 000/100 VUs, 91550 complete and 0 interrupted iterations
default ✓ [======================================] 000/100 VUs  45s
```


```bash
k6 run get.js

         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/

     execution: local
        script: get.js
        output: -

     scenarios: (100.00%) 1 scenario, 50 max VUs, 1m15s max duration (incl. graceful stop):
              * default: Up to 50 looping VUs for 45s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


     ✓ retrieved successfully (status 200)

     checks.........................: 100.00% 1323853 out of 1323853
     data_received..................: 1.1 GB  25 MB/s
     data_sent......................: 115 MB  2.6 MB/s
     http_req_blocked...............: avg=2.05µs  min=0s     med=1µs      max=28.43ms p(90)=2µs    p(95)=3µs
     http_req_connecting............: avg=7ns     min=0s     med=0s       max=513µs   p(90)=0s     p(95)=0s
     http_req_duration..............: avg=1.35ms  min=51µs   med=658µs    max=49.71ms p(90)=3.12ms p(95)=5.15ms
       { expected_response:true }...: avg=1.35ms  min=51µs   med=658µs    max=49.71ms p(90)=3.12ms p(95)=5.15ms
     http_req_failed................: 0.00%   0 out of 1323853
     http_req_receiving.............: avg=30.23µs min=5µs    med=10µs     max=28.67ms p(90)=31µs   p(95)=42µs
     http_req_sending...............: avg=7.44µs  min=1µs    med=2µs      max=27.52ms p(90)=7µs    p(95)=10µs
     http_req_tls_handshaking.......: avg=0s      min=0s     med=0s       max=0s      p(90)=0s     p(95)=0s
     http_req_waiting...............: avg=1.31ms  min=42µs   med=633µs    max=49.7ms  p(90)=3.06ms p(95)=5.05ms
     http_reqs......................: 1323853 29418.945749/s
     iteration_duration.............: avg=1.41ms  min=65.2µs med=700.79µs max=49.73ms p(90)=3.22ms p(95)=5.28ms
     iterations.....................: 1323853 29418.945749/s
     vus............................: 1       min=1                  max=50
     vus_max........................: 50      min=50                 max=50


running (0m45.0s), 00/50 VUs, 1323853 complete and 0 interrupted iterations
default ✓ [======================================] 00/50 VUs  45s
```
