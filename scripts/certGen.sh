#!/bin/bash
# call this script with an email address (valid or not).
# like:
# ./makecert.sh joe@random.com
mkdir certs
rm certs/*
echo "make server cert"
openssl req -new -nodes -x509 -out certs/server.pem -keyout certs/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=$1"
echo "make client cert"
openssl req -new -nodes -x509 -out certs/client.pem -keyout certs/client.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=$1"
# echo "make CA"
# openssl genrsa -out rootCA.key 4096
# openssl req -x509 -new -key certs/rootCA.key -days 3650 -out certs/rootCA.crt
# echo "make cert for CA"
# openssl genrsa -out backdoor.aws.araalinetworks.com.key 2048
# openssl req -new -key backdoor.aws.araalinetworks.com.key -out backdoor.aws.araalinetworks.com.csr
# openssl x509 -req -in backdoor.aws.araalinetworks.com.csr -CA certs/rootCA.crt -CAkey certs/rootCA.key -CAcreateserial -days 365 -out backdoor.aws.araalinetworks.com.crt