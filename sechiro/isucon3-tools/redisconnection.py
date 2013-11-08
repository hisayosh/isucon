# sudo pip install redis 

# connection pool
pool = redis.ConnectionPool(host='localhost', port=6379, db=0)
r = redis.StrictRedis(connection_pool=pool)

# unix socket
r = redis.StrictRedis(unix_socket_path='/tmp/redis.sock')

# pipeline
p = r.pipeline()

# sample
from time import sleep
def sample():
    p.set('hoge','hogehoge', px=10)
    p.append('hoge','hogehoge2')
    p.set('fuzz','buzz')
    p.get('hoge')
    p.get('fuzz')
    p.get('novalue')
    res = p.execute()
    print res 
    sleep(1)
    print r.get('hoge')
        # hash
    r.hset('hash','key1','value1')
    print r.hgetall('hash')
    r.expire('hash', 1)
    res = r.hset('hash','key2','value2')
    res = r.hgetall('hash')
    print r.ttl('hash')
    print type(res)
    print res
    sleep(2)
    res = r.hgetall('hash')
    print res['key2']