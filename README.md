zkdo
====

Apache ZooKeeper meets *NIX

---

`zkdo` currently does two simple things:

1. (Optionally) Wait on a lock in ZK before running a subcommand you provide.
2. (Optionally) Register a string you provide as an ephemeral sequence node so you can discover your subprocess from elsewhere.

It's mostly useful for bringing naive programs (especially ones you don't want to modify) into your ZK-managed deployment.

---

In the example below, imagine we have a simple service that should have one and only one instance running at a time. The problem is you want a new instance to take over if the original instance dies or loses network connectivity.

You might normally run your service like so:
```
my-service --port 7777
```

In order to run it "under" zkdo, you'd do this:
```
zkdo --lock=/locks/myservice \
     --register=/services/myservice \
     --data="{\"host\":\"$(hostname)\",\"port\":7777}" \
     -- \
     my-service --port 7777

2013/05/27 11:55:14 Connected to zk.
2013/05/27 11:55:14 Created lock: /locks/myservice
2013/05/27 11:55:14 Registered at: /services/myservice0000000002
2013/05/27 11:55:14 Registered data: {"host":"scumbag","port":7777}
2013/05/27 11:55:14 Running command.
[26595] 27 May 11:55:14.553 # Server started, MyService version 1.0.0
[26595] 27 May 11:55:14.553 * The server is now ready to accept connections on port 7777
```

Now, on another server:
```
zkdo --lock=/locks/myservice \
     --register=/services/myservice \
     --data="{\"host\":\"$(hostname)\",\"port\":7777}" \
     -- \
     my-service --port 7777

2013/05/27 11:58:32 Connected to zk.
2013/05/27 11:58:32 Couldn't obtain lock: /locks/myservice
2013/05/27 11:58:32 Waiting on lock watch.

```

Now kill the first instance of zkdo and check on the second instance, it will (once the first dies) take the lock and start its subprocess:
```
2013/05/27 11:59:18 Lock changed, retrying obtain.
2013/05/27 11:59:18 Created lock: /locks/myservice
2013/05/27 11:59:18 Registered at: /services/myservice0000000004
2013/05/27 11:59:18 Registered data: {"host":"scumbag","port":7777}
2013/05/27 11:59:18 Running command.
[27376] 27 May 11:59:18.218 # Server started, MyService version 1.0.0
[27376] 27 May 11:59:18.218 * The server is now ready to accept connections on port 7777
```

Ta-da!

Alternatively, if you don't care how many instances run at once, you can leave out the `--lock` argument and just use zkdo to register running services in ZK for discovery. You can also leave out `--register` (but use `--lock`) if you don't need discovery at all.
