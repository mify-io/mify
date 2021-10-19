package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chebykinn/mify/internal/mify"
)

func main() {
	createWorkspaceCmd := flag.NewFlagSet("create-workspace", flag.ExitOnError)

	createServiceCmd := flag.NewFlagSet("create-service", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Printf("usage: %s create-service|create-workspace <name>\n", os.Args[0])
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create-workspace":
		createWorkspaceCmd.Parse(os.Args[2:])
		if len(createWorkspaceCmd.Args()) != 1 {
			fmt.Printf("usage: %s create-workspace <workspace-name>\n", os.Args[0])
			os.Exit(1)
		}
		if err := mify.CreateWorkspace(createWorkspaceCmd.Args()[0]); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create workspace: %s\n", err)
			os.Exit(2)
		}

	case "create-service":
		createServiceCmd.Parse(os.Args[2:])
		if len(createServiceCmd.Args()) != 1 {
			fmt.Printf("usage: %s create-service <service-name>\n", os.Args[0])
			os.Exit(1)
		}
		if err := mify.CreateService(createServiceCmd.Args()[0]); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create service: %s\n", err)
			os.Exit(2)
		}
	default:
		flag.Usage()
		fmt.Printf("usage: %s create-service|create-workspace <name>\n", os.Args[0])
		os.Exit(1)
	}
}
