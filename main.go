package main

import (
	"boxlister/cli"
	"boxlister/instance"
	"fmt"
	"os"
	"bytes"
	"boxlister/files"
	"os/user"
)

const template = `
Host %s
HostName %s
User %s
`

func main() {

	cliFlags := cli.ParseFlags()
	if cliFlags == nil {
		os.Exit(1)
	}

	fmt.Println(cliFlags)
	var instances []*instance.Instance

	for _, profile := range cliFlags.Profiles {
		profileInstances := instance.DescribeInstances(profile, cliFlags.Region)
		instances = append(instances, profileInstances...)
	}

	var outBuf bytes.Buffer
	for _, namePart := range cliFlags.InstanceNameParts {
		for _, inst := range instances {
			if inst.NameMatches(namePart, cliFlags.InstancePrefix) {
				serverConfigString := fmt.Sprintf(template, *inst.Name, *inst.DnsName, cliFlags.SshUser)
				outBuf.WriteString(serverConfigString)
			}
		}
	}
	if cliFlags.GenerateFile {
		current, e := user.Current()
		handleError(e)
		configPath := current.HomeDir + "/.ssh/config"

		files.FillGenerated(&configPath, outBuf)
	} else {
		fmt.Fprint(os.Stdout, outBuf.String())
	}
}

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}
