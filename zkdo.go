package main

import (
	"flag"
	// "fmt"
	zk "github.com/bretthoerner/gozk"
	"log"
	"time"
	"os/exec"
	"os"
	// "strings"
)

var (
	zkServers = flag.String(
		"zk",
		"localhost:2181",
		"the comma separated ZK connection string")

	zkTimeout = flag.Int(
		"zk-timeout",
		5,
		"the ZK connection timeout (in seconds)")

	lock = flag.String(
		"lock",
		"",
		"path of lock to hold before running cmd")

	noblock = flag.Bool(
		"noblock",
		false,
		"instead of waiting for the lock, exit immediately")

	register = flag.String(
		"register",
		"",
		"path of ephemeral node to register")

	data = flag.String(
		"data",
		"",
		"string data to register in ephemeral node")
)

func main() {
	flag.Parse()

	// if len(flag.Args()) < 1 {
	// 	log.Fatalln("Command required")
	// }

	if len(*lock) == 0 && len(*register) == 0 {
		log.Fatalln("At least one of --lock or --register is required.")
	}

	if len(*register) > 0 && len(*data) == 0 {
		log.Fatalln("Must provide --data to to set with --register.")
	}

	if len(*data) > 0 && len(*register) == 0 {
		log.Fatalln("Must provide --register path where --data will be set.")
	}

	timeoutDuration := time.Duration(*zkTimeout) * time.Second

	conn, session, err := zk.Dial(*zkServers, timeoutDuration)
	if err != nil {
		log.Fatalf("Can't connect: %v", err)
	}
	defer conn.Close()

	// wait for connection.
	event := <-session
	if event.State != zk.STATE_CONNECTED {
		log.Fatalf("Can't connect: %v", event)
	} else {
		log.Println("Connected to zk.")
	}

	if len(*lock) > 0 {
		for {
			// ensure we have the lock before proceeding
			_, err = conn.Create(*lock, "", zk.EPHEMERAL, zk.WorldACL(zk.PERM_ALL))
			if err != nil {
				if *noblock {
					log.Fatalf("Couldn't obtain lock: %v", *lock)					
				} else {
					log.Printf("Couldn't obtain lock: %v", *lock)	
				}

				if !zk.IsError(err, zk.ZNODEEXISTS) {
					log.Fatalf("Unknown error: %v", err)
				}

				_, _, watch, err := conn.GetW(*lock)
				if err != nil {
					log.Fatalf("Unknown error: %v", err)
				}

				log.Println("Waiting on lock watch.")
				for {
					select {
					case <-watch:
						log.Println("Lock changed, retrying obtain.")
						break
					default:
						time.Sleep(1 * 1e9)
						continue
					}
					break
				}

			} else {
				log.Printf("Created lock: %v", *lock)
				break
			}
		}
	}

	// register ourselves if necessary
	if len(*register) > 0 {
		path, err := conn.Create(*register, *data, zk.SEQUENCE|zk.EPHEMERAL, zk.WorldACL(zk.PERM_ALL))
		if err != nil {
			log.Fatalf("Error while registering: %v", err)
		} else {
			log.Printf("Registered at: %v", path)
			log.Printf("Registered data: %v", *data)
		}
	}

	// run subprocess
	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
	
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error running subprocess: %v", err)
	} else {
		log.Println("Running command.")
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Error from command: %v", err)
	}
}
