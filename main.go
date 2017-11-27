package main

import (
	"boxlister/cli"
	"boxlister/files"
	"boxlister/instance"
	"bytes"
	"fmt"
	"os"
	"os/user"
)

const (
	fileTemplate    = "\n# %s\nHost%s\nHostName %s\nUser %s\n"
	consoleTemplate = "%s %s %s\n"
)

func main() {

	cliFlags := cli.ParseFlags()
	if cliFlags == nil {
		os.Exit(1)
	}

	fmt.Println(cliFlags)
	var instances []*instance.Instance

	for _, profile := range cliFlags.Profiles {
		profileInstances := instance.Fetch(profile, cliFlags.Region)
		instances = append(instances, profileInstances...)
	}

	var outBuf bytes.Buffer
	for _, namePart := range cliFlags.InstanceNameParts {
		for _, inst := range instances {
			if inst.NameMatches(namePart, cliFlags.InstancePrefix) {
				serverConfigString := fillTemplate(inst, cliFlags)
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

func fillTemplate(inst *instance.Instance, cliFlags *cli.CliFlags) string {
	if cliFlags.GenerateFile {
		return fmt.Sprintf(fileTemplate, *inst.Id, *inst.Name, *inst.DnsName, cliFlags.SshUser)
	}
	return fmt.Sprintf(consoleTemplate, *inst.Id, *inst.Name, *inst.DnsName)

}

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}
