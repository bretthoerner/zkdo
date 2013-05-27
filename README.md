zkdo
====

Apache ZooKeeper meets *NIX

Example below (note that Redis is only used as an example, I see no reason why someone would want to do this):

```
zkdo --lock=/locks/redis \
     --register=/services/redis \
     --data="{\"host\":\"$(hostname)\",\"port\":7777}" \
     -- \
     redis-server --port 7777

2013/05/27 11:55:14 Connected to zk.
2013/05/27 11:55:14 Created lock: /locks/redis
2013/05/27 11:55:14 Registered at: /services/redis0000000002
2013/05/27 11:55:14 Registered data: {"host":scumbag,"port":7777}
2013/05/27 11:55:14 Running command.
                _._                                                  
           _.-``__ ''-._                                             
      _.-``    `.  `_.  ''-._           Redis 2.6.13 (00000000/0) 64 bit
  .-`` .-```.  ```\/    _.,_ ''-._                                   
 (    '      ,       .-`  | `,    )     Running in stand alone mode
 |`-._`-...-` __...-.``-._|'` _.-'|     Port: 7777
 |    `-._   `._    /     _.-'    |     PID: 26595
  `-._    `-._  `-./  _.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |           http://redis.io        
  `-._    `-._`-.__.-'_.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |                                  
  `-._    `-._`-.__.-'_.-'    _.-'                                   
      `-._    `-.__.-'    _.-'                                       
          `-._        _.-'                                           
              `-.__.-'                                               

[26595] 27 May 11:55:14.553 # Server started, Redis version 2.6.13
[26595] 27 May 11:55:14.553 * The server is now ready to accept connections on port 7777
```

Now, in another terminal:
```
zkdo --lock=/locks/redis \
     --register=/services/redis \
     --data="{\"host\":\"$(hostname)\",\"port\":7777}" \
     -- \
     redis-server --port 7777

2013/05/27 11:58:32 Connected to zk.
2013/05/27 11:58:32 Couldn't obtain lock: /locks/redis
2013/05/27 11:58:32 Waiting on lock watch.

```

Now kill the original redis instance and check the second terminal:
```
2013/05/27 11:59:18 Lock changed, retrying obtain.
2013/05/27 11:59:18 Created lock: /locks/redis
2013/05/27 11:59:18 Registered at: /services/redis0000000004
2013/05/27 11:59:18 Registered data: {"host":"scumbag","port":7777}
2013/05/27 11:59:18 Running command.
                _._                                                  
           _.-``__ ''-._                                             
      _.-``    `.  `_.  ''-._           Redis 2.6.13 (00000000/0) 64 bit
  .-`` .-```.  ```\/    _.,_ ''-._                                   
 (    '      ,       .-`  | `,    )     Running in stand alone mode
 |`-._`-...-` __...-.``-._|'` _.-'|     Port: 7777
 |    `-._   `._    /     _.-'    |     PID: 27376
  `-._    `-._  `-./  _.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |           http://redis.io        
  `-._    `-._`-.__.-'_.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |                                  
  `-._    `-._`-.__.-'_.-'    _.-'                                   
      `-._    `-.__.-'    _.-'                                       
          `-._        _.-'                                           
              `-.__.-'                                               

[27376] 27 May 11:59:18.218 # Server started, Redis version 2.6.13
[27376] 27 May 11:59:18.218 * The server is now ready to accept connections on port 7777
```

You can also leave out `--lock` if you just want to register multiple instances of your service for discovery/load-balancing/failover.
