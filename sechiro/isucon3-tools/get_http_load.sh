#!/bin/bash
cd /tmp
wget http://acme.com/software/http_load/http_load-12mar2006.tar.gz
tar xf http_load-12mar2006.tar.gz
cd http_load-12mar2006/
make
sudo cp http_load /usr/local/bin

