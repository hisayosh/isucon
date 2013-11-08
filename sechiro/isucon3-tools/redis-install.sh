wget http://download.redis.io/releases/redis-2.6.16.tar.gz
tar xf redis-2.6.16.tar.gz
cd redis-2.6.16
make && make install
sudo mkdir /var/lib/redis
echo "You must set /etc/redis.conf"