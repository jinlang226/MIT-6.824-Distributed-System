# Big Storage

* performance -> sharding
* faucts -> toleranc
* tolerance -> replication
* repel -> in consistency
* consistency -> low performance
* strong consistancy

## Strong Consistency
* c1 -> w x 1
* c2 -> w x 2 (write at the same time)
* c3 -> Rx
* c4 -> Rx

## Bad Repel Design
later reading

# GFS
* Big, fast
* global
* sharding
* automatic recovery
* single data center
* internal use
* big sequential random access

Not gaurantee correct data, but performance

## Master Data
* (Non-volatile memory (NVM) or non-volatile storage is a type of computer memory that can retain stored information even after power is removed. In contrast, volatile memory needs constant power in order to retain data.)
* file name -> array of chunk (64 megaByte?bits) handles (on disk, nv, Non-Volatile CACHE)
* handle -> 
    * list of chunk servers(not writen in the disk, v), 
    * version# (nv), 
    * primary (v, forgets, 60s), 
    * lease expiration (v) 
* log (append is efficient), checkpoint -> Disk

## Read
1. file name and offset (read data from) -> Master
2. Master sends chunk handle and a list of servers (client cached the result)
3. Client -> Chunk servers (store each chunk in seperate linux file hard drive)
    * retuen the data

## Writes
1. No primary - on master(need to find the up to date replicas-- the version that master remembers)
2. Pick primary, secondary
3. Increments the version #, and write it to the disk
4. Then tells P, S, V# -- lease(60s)
5. M write its v# to the disk after telling the disk
6. Primary picks offset. All replicas told to write at offeset
7. if all of then "yes", then primary reply success to the client; else no to client
