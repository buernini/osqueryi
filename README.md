# osqueryi
 
This project mainly provides osquery remote management plan;
 
go build server.go
go build client.go
 
server :  server -h 0.0.0.0 -p 3993
client : client -h  127.0.0.1 -p 3993
 
127.0.0.1:3993> select * from users limit 3;
 description        | directory       | gid | gid_signed | shell          | uid | uid_signed | username     | uuid                                 |   
____________________________________________________________________________________________________________________________________________________
                                                                                                                                                                                                            
 AMaViS Daemon      | /var/virusmails | 83  | 83         | /usr/bin/false | 83  | 83         | _amavisd     | FFFFEEEE-DDDD-CCCC-BBBB-AAAA00000053 |
 AppleEvents Daemon | /var/empty      | 55  | 55         | /usr/bin/false | 55  | 55         | _appleevents | FFFFEEEE-DDDD-CCCC-BBBB-AAAA00000037 |
 Application Owner  | /var/empty      | 87  | 87         | /usr/bin/false | 87  | 87         | _appowner    | FFFFEEEE-DDDD-CCCC-BBBB-AAAA00000057 |
____________________________________________________________________________________________________________________________________________________
 

