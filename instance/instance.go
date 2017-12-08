package instance

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//Instance instance representation
type Instance struct {
	Name    string
	DnsName string
	Id      string
	checked bool
}

func Fetch(profile string, region string, prefix string, nameParts []string) []*Instance {
	sess := awsSession(profile, region)
	ec2Sess := ec2.New(sess)
	instances, err := ec2Sess.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		panic(err)
	}

	var ec2Machines []*Instance
	for _, r := range instances.Reservations {
		for _, instance := range r.Instances {
			ec2Machine := asInstance(instance)
			if ec2Machine.IsValid() {
				for _, namePart := range nameParts {
					if ec2Machine.NameMatches(namePart, prefix) {
						ec2Machines = append(ec2Machines, ec2Machine)
					}
				}
			}
		}
	}
	return ec2Machines
}

func asInstance(inst *ec2.Instance) *Instance {
	nameTag := nameTag(inst)
	return &Instance{
		DnsName: publicDns(inst),
		Name:    nameTag,
		Id:      *inst.InstanceId,
	}
}

func awsSession(profile string, region string) *session.Session {
	options := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if profile != "" {
		options.Profile = profile
	}
	if region != "" {
		options.Config = aws.Config{
			Region: &region,
		}
	}
	return session.Must(session.NewSessionWithOptions(options))
}

func (i *Instance) String() string {
	return fmt.Sprintf("%v:%v", i.Name, i.DnsName)
}

func (i *Instance) IsValid() bool {
	return i.DnsName != "" && i.Name != ""
}

func (i *Instance) NameMatches(namePart string, namePrefix string) bool {
	return i.Name != "" && strings.Index(i.Name, namePrefix) == 0 && strings.Contains(i.Name, namePart)
}

func (i *Instance) PrintOut() string {
	return fmt.Sprintf("%s %s %s", i.Id, i.Name, i.DnsName)
}
func (i *Instance) PrintOutSshFormat(sshUser string) string {
	return fmt.Sprintf("\n# %s\nHost %s\nHostName %s\nUser %s\n", i.Id, i.Name, i.DnsName, sshUser)
}

func publicDns(instance *ec2.Instance) string {
	for _, networkInterface := range instance.NetworkInterfaces {
		if networkInterface != nil && networkInterface.Association != nil && networkInterface.Association.PublicDnsName != nil {
			return *networkInterface.Association.PublicDnsName
		}
	}
	return ""
}
func nameTag(instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}

	return ""
}
