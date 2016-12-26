#GRIGORI
##Overview
Grigori will be reading events, processing them and sending them somewhere else.

The original idea is to have a "shipper" kind of logstash shipper. So, read from logs (tail like), transform them if necessary and send them, in my case to redis, where finally will be read by logstash indexer and sent to ES.

##Goals
1) Learn some golang myself
2) Try to use a DDD like with golang (that's the reason for the strange code structure)
3) Learn how to test with golang

##Motivation
Well, while I really like using some of their tools, I really don't like logstash as a shipper (I like as an indexer though). I think it is quite resource needy. So I decided to try to do one myself and see if it is better.

I know they already have beats now, but I was kind of discussed if they would keep supporting sending to redis, and that is a functionality that I really want.

##Config file example
```
{
  "config": {
    "writer": {
      "type": "redis",
      "n_writers": 2,
      "redis_host": "172.17.0.1",
      "redis_port": "6379",
      "redis_key": "grig"
    },
    "resources": [
      {
        "reader": "tail",
        "processor": "plaintologstash",
        "n_processors": 20,
        "resource": "/tmp/testgrigori89833.log",
        "maxlines": 33,
        "tags": ["tag1test","tag2test"],
        "type": "testtype",
        "version": 1
      },
      {
        "reader": "tail",
        "processor": "plaintojson",
        "n_processors": 2,
        "resource": "/tmp/testgrigori77263.log"
        "tags": ["tag1test","tag2test"],
        "type": "testtype"
      }
    ],
    "position_keeper": {
      "type": "file",
      "path": "/tmp/test_grigori_position_keeper.pos"
    }
  }
}

```

##Execution
###Help
```
./grigori -h
```
###Run
```
./grigori
```
```
./grigori -c [path to config file]
```