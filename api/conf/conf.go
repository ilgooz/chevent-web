package conf

import (
	"flag"
	"log"

	"gopkg.in/mgo.v2"
)

var (
	addr      = flag.String("addr", ":3000", "Server Address")
	mongoAddr = flag.String("mongo", "mongodb://127.0.0.1:27017/chevent", "Mongodb Address")
	hash      = flag.String("hash", "", "latest sha")
)

var (
	Addr string
	Hash string
	M    *mgo.Session
)

func Load() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	Addr = *addr
	Hash = *hash
	M = dialMongo()
}
