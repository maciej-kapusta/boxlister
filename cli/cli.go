package cli

import (
	"flag"
	"fmt"
	"os"
)

type CliFlags struct {
	Profile          *string
	SshUser          *string
	InstanceNamePart *string
	OutputFile       *string
	Region           *string
}

func ParseFlags() *CliFlags {

	profile := flag.String("profile", "", "Optional AWS Profile. If missing the profile from your AWS CLI used.")
	user := flag.String("user", "", "ssh user (required)")
	instanceNamePart := flag.String("instance", "", "part of instance name to be matched(required)")
	outputFile := flag.String("out", "", "Optional output file to save instances and use in ssh. If missing std out will be used")
	region := flag.String("region", "", "Optional AWS region. If absent the AWS CLI config is used.")
	flag.Parse()

	if *user == "" || *instanceNamePart == "" {
		fmt.Fprintln(os.Stderr, "Missing parameters. Usage: boxlister -user=joe.doe [-profile=acme-prod] -instance=db [-out=my_ssh_file]")
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
