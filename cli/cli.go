package cli

import (
	"flag"
	"fmt"
	"os"
	user2 "os/user"
)

type CliFlags struct {
	Profile          *string
	SshUser          *string
	InstanceNamePart *string
	OutputFile       *string
	Region           *string
}

func ParseFlags() *CliFlags {
	current, e := user2.Current()
	if e != nil {
		panic(e)
	}
	profile := flag.String("profile", "", "Optional AWS Profile. If missing the profile from your AWS CLI used.")
	user := flag.String("user", current.Username, "ssh user. If empty your login will be used")
	instanceNamePart := flag.String("instance", "", "part of instance name to be matched(required)")
	outputFile := flag.String("out", "", "Optional output file to save instances and use in ssh. If missing std out will be used")
	region := flag.String("region", "us-east-1", "Optional AWS region. If absent the AWS CLI config is used.")
	flag.Parse()

	if *instanceNamePart == "" {
		fmt.Fprintln(os.Stderr, "Missing parameters. Usage: boxlister [-user=joe.doe] [-profile=acme-prod] -instance=db [-out=my_ssh_file]")
		flag.PrintDefaults()
		return nil
	}

	return &CliFlags{
		Profile:          profile,
		SshUser:          user,
		InstanceNamePart: instanceNamePart,
		OutputFile:       outputFile,
		Region:           region,
	}
}
