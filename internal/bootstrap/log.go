package bootstrap

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func InitLog(debugLevel string) {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
	})
	if debugLevel, err := log.ParseLevel(debugLevel); err != nil {
		fmt.Printf("Fail to parse debug level \"%s\".\nFatal Error, Exiting...\n", debugLevel)
		os.Exit(1)
	} else {
		log.SetLevel(debugLevel)
	}
}
