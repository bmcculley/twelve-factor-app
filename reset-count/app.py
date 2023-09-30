from redis import Redis

redisDb = Redis(host="redis", db=0, socket_timeout=5)

redisDb.set('visitorCount', 0)