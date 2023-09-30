# Admin Reset process

If the count needs to be reset this can be run seperately.

First build it.

```bash
 docker run -it --rm --link <redis>:redis --net <network> reset-count
 ```

If the twelve-factor app was started using docker-compose it will be
running on a different network. In order to attached to the redis 
container the network name will also need to be supplied.

```bash
docker network ls
```