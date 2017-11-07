package main

import (
	"boxlister/cli"
	"boxlister/instance"
	"fmt"
	"os"
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

	instances := instance.DescribeInstances(cliFlags.Profile, cliFlags.Region)

	out := os.Stdout
	if *cliFlags.OutputFile != "" {

		file, e := os.Create(*cliFlags.OutputFile)
		if e != nil {
			panic(e)
		}
		defer file.Close()
		out = file
	}
	for _, inst := range instances {
		if inst.NameContains(*cliFlags.InstanceNamePart) {
			serverConfigString := fmt.Sprintf(template, *inst.Name, *inst.DnsName, *cliFlags.SshUser)
			out.WriteString(serverConfigString)
		}
	}
}
