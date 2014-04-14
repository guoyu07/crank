package main

import (
	//	"../../pkg/crank"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Command func(*rpc.Client) error
type CommandSetup func(*flag.FlagSet) Command

var run string
var commands map[string]CommandSetup

func init() {
	defaultFlags(flag.CommandLine)

	commands = make(map[string]CommandSetup)
	commands["echo"] = Echo
}

func defaultFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&run, "run", run, "path to control socket")
}

func fail(reason string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, reason, args...)
	flag.Usage()
	os.Exit(1)
}

func main() {
	var err error
	flag.Parse()

	command := flag.Arg(0)

	if command == "" {
		fail("command missing\n")
	}

	cmdSetup, ok := commands[command]
	if !ok {
		fail("unknown command %s\n", command)
	}

	flagSet := flag.NewFlagSet(os.Args[0]+" "+command, flag.ExitOnError)
	defaultFlags(flagSet)
	cmd := cmdSetup(flagSet)

	if err = flagSet.Parse(flag.Args()); err != nil {
		fail("oops: %s\n", err)
	}

	client, err := rpc.Dial("unix", run)
	if err != nil {
		fail("Couldn't connect: %s\n", err)
	}

	if err = cmd(client); err != nil {
		fail("command failed: %v\n", err)
	}
}

func Echo(flag *flag.FlagSet) Command {
	var msg string
	flag.StringVar(&msg, "msg", "foo", "message to send")

	return func(client *rpc.Client) (err error) {
		var reply string

		if err = client.Call("crank.Echo", &msg, &reply); err != nil {
			return
		}

		log.Println("echo reply: ", reply)
		return
	}
}
