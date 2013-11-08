#!/bin/bash
cd /tmp
wget http://openresty.org/download/ngx_openresty-1.4.3.1.tar.gz
tar xf ngx_openresty-1.4.3.1.tar.gz
cd ngx_openresty-1.4.3.1/
sudo yum -y install pcre-devel

./configure --with-luajit --with-http_gzip_static_module
make && sudo make install

wget -O nginx https://fzrxefe.googlecode.com/files/openresty.init.d.script
chmod +x nginx
sudo cp nginx /etc/init.d
