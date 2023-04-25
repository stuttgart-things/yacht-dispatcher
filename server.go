/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package main

import (
	"os"

	// "github.com/go-git/plumbing/color"

	// internal "codehub.sva.de/Lab/stuttgart-things/yacht/yacht-application-dispatcher/internal"
	"github.com/fatih/color"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	// redis "github.com/go-redis/redis/v7"
	// redisqueue "github.com/robinjoseph08/redisqueue/v2"
	goVersion "go.hein.dev/go-version"
)

var (
	redisAddress    = os.Getenv("REDIS_SERVER")
	redisPort       = os.Getenv("REDIS_PORT")
	redisPassword   = os.Getenv("REDIS_PASSWORD")
	redisQueue      = os.Getenv("REDIS_QUEUE")
	tektonNamespace = os.Getenv("TEKTON_NAMESPACE")
	shortened       = false
	version         = "unset"
	date            = "unknown"
	commit          = "unknown"
	output          = "yaml"
	// redisClient     = redis.NewClient(&redis.Options{
	// 	Addr:     redisAddress + ":" + redisPort,
	// 	Password: redisPassword,
	// 	DB:       0,
	// })
	log         = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	logfilePath = "YD.log"
)

const banner = `
___    ___ ________
|\  \  /  /|\   ___ \
\ \  \/  / \ \  \_|\ \
 \ \    / / \ \  \ \\ \
  \/  /  /   \ \  \_\\ \
__/  / /      \ \_______\
|\___/ /        \|_______|
\|___|/

`

func main() {

	// c, err := redisqueue.NewConsumerWithOptions(&redisqueue.ConsumerOptions{
	// 	VisibilityTimeout: 60 * time.Second,
	// 	BlockingTimeout:   15 * time.Second,
	// 	ReclaimInterval:   1 * time.Second,
	// 	BufferSize:        100,
	// 	Concurrency:       1,
	// 	RedisClient:       redisClient,
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// c.Register(redisQueue, internal.CreateApplicationWorkerJobs)

	// handle errors accordingly
	// go func() {
	// 	for err := range c.Errors {
	// 		fmt.Printf("err: %+v\n", err)
	// 	}
	// }()

	// Output banner + version output
	color.Cyan(banner)
	color.Cyan("YACHT DISPATCHER")
	resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)
	color.Cyan(resp + "\n")

	log.Info("YD server started - waiting for messages in ", redisQueue)

	// c.Run()

	log.Warn("YD server stopped")
}
