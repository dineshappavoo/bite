#!/bin/bash
User=""
touch pid.txt


#start the mysql server
cd /usr/local/mysql
sudo -i && ./bin/mysqld_safe &
#write pid to a file
echo $! > pid.txt


#start redis server
redis-server &


#Export the GOPATH
export GOPATH=/Users/Dany/Documents/FALL-2013-COURSES/Imp_Data_structures/workspace/projectcube/


#Start the elastic search server
cd /Users/Dany/Documents/SOFTWARES/elasticsearch-1.4.4/bin/
./elasticsearch &
#write pid to a file
echo $! > pid.txt


#start the application server
cd /Users/Dany/Documents/FALL-2013-COURSES/Imp_Data_structures/workspace/projectcube/src/app
go run server.go &
#write pid to a file
echo $! > pid.txt
