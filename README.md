# osqueryi
 
This project mainly provides osquery remote management plan;
 
go build server.go
go build client.go
 
server :  server -h x.x.x.x -P 3993
client : client -h x.x.x.x -u frank -p 123456 -P 3993
 
127.0.0.1:3993> select * from users limit 3;
____________________________________________________________________________________________________________________________________________________
 
127.0.0.1:3993> .help

