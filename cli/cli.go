package cli

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"
)

type CliFlags struct {
	Profiles          []string
	SshUser           string
	InstanceNameParts []string
	InstancePrefix    string
	GenerateFile      bool
	Region            string
	HealthCheck       string
}

func ParseFlags() *CliFlags {
	current, e := user.Current()
	if e != nil {
		panic(e)
	}
	rawProfiles := flag.String("profile", "", "Optional AWS Profiles. If missing the profile from your AWS CLI used.")
	userName := flag.String("user", current.Username, "ssh user. If empty your login will be used")
	instanceNamePartsRaw := flag.String("instance", "", "comma separated list of parts of instance name to be matched(required)")
	instancePrefix := flag.String("prefix", "", "instance name common prefix")
	generateFile := flag.Bool("generate_file", false, "generate ~/.ssh/config")
	region := flag.String("region", "us-east-1", "Optional AWS region. If absent the AWS CLI config is used.")
	healthcheck := flag.String("healthcheck", "", "healthcheck path after host name. :8080/healthcheck")
	flag.Parse()

	instanceNamePartsParsed := splitCommas(instanceNamePartsRaw)

	profiles := splitCommas(rawProfiles)

	if *instanceNamePartsRaw == "" {
		fmt.Fprintln(os.Stderr, "Missing parameters. Usage: boxlister [-user=joe.doe] [-profile=acme-prod] -prefix=zeropark -instance=db,bidder [-generate_file] [-healthcheck=:8080/healtcheck]")
		flag.PrintDefaults()
		return nil
	}

	if *healthcheck != "" && *generateFile {
		fmt.Println("Healthcheck and generate file together. Healthcheck overrides file generation.")
	}

	return &CliFlags{
		Profiles:          profiles,
		SshUser:           *userName,
		InstanceNameParts: instanceNamePartsParsed,
		InstancePrefix:    *instancePrefix,
		GenerateFile:      *generateFile,
		Region:            *region,
		HealthCheck:       *healthcheck,
	}
}

func splitCommas(partsRaw *string) []string {
	split := strings.Split(*partsRaw, ",")
	var parsed []string
	for _, part := range split {
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			continue
		}
		parsed = append(parsed, part)
	}
	return parsed
}
