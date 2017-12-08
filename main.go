package main

import (
	"boxlister/cli"
	"boxlister/files"
	"boxlister/instance"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"time"
)

func main() {

	cliFlags := cli.ParseFlags()
	if cliFlags == nil {
		os.Exit(1)
	}

	fmt.Println(cliFlags)
	instances := fetchInstances(cliFlags)

	if cliFlags.GenerateFile {
		generateSshFile(cliFlags, instances)
	} else if cliFlags.HealthCheck != "" {
		printHealthCheck(instances, cliFlags)
	} else {
		for _, inst := range instances {
			fmt.Println(inst.PrintOut())
		}
	}
}
func printHealthCheck(instances []*instance.Instance, cliFlags *cli.CliFlags) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	result := make(chan string)
	for _, inst := range instances {
		go func() {
			resp, err := client.Get("http://" + inst.DnsName + cliFlags.HealthCheck)
			if err != nil {
				result <- err.Error()
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				result <- inst.PrintOut() + " ok"
				return
			} else {
				result <- inst.PrintOut() + " not ok"
			}

		}()
	}

	for i := 0; i < len(instances); i++ {
		status := <-result
		fmt.Println(status)
	}
}

func generateSshFile(cliFlags *cli.CliFlags, instances []*instance.Instance) {
	var outBuf bytes.Buffer
	for _, inst := range instances {
		serverConfigString := inst.PrintOutSshFormat(cliFlags.SshUser)
		outBuf.WriteString(serverConfigString)
	}
	current, e := user.Current()
	handleError(e)
	configPath := current.HomeDir + "/.ssh/config"
	files.FillGenerated(&configPath, outBuf)
}

func fetchInstances(cliFlags *cli.CliFlags) []*instance.Instance {
	var instances []*instance.Instance
	if len(cliFlags.Profiles) == 0 {
		instances = appendInstances("", cliFlags, instances)
	}

	for _, profile := range cliFlags.Profiles {
		instances = appendInstances(profile, cliFlags, instances)
	}
	return instances
}

func appendInstances(profile string, cliFlags *cli.CliFlags, instances []*instance.Instance) []*instance.Instance {
	profileInstances := instance.Fetch(profile, cliFlags.Region, cliFlags.InstancePrefix, cliFlags.InstanceNameParts)

	instances = append(instances, profileInstances...)
	return instances
}

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}
