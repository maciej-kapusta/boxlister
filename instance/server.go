package instance

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

//Instance instance representation
type Instance struct {
	Name    *string
	DnsName *string
	Id      *string
}

func asInstance(intance *ec2.Instance) *Instance {
	return &Instance{
		DnsName: publicDns(intance),
		Name:    nameTag(intance),
		Id:      intance.InstanceId,
	}
}

func (i *Instance) String() string {
	return fmt.Sprintf("%v:%v", *i.Name, *i.DnsName)
}

func (i *Instance) IsValid() bool {
	return i.DnsName != nil && i.Name != nil
}

func (i *Instance) NameMatches(namePart string, namePrefix string) bool {
	return i.Name != nil && strings.Index(*i.Name, namePrefix) == 0 && strings.Contains(*i.Name, namePart)
}

func publicDns(instance *ec2.Instance) *string {
	for _, networkInterface := range instance.NetworkInterfaces {
		if networkInterface != nil && networkInterface.Association != nil {
			return networkInterface.Association.PublicDnsName
		}
	}
	return nil
}
func nameTag(instance *ec2.Instance) *string {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return tag.Value
		}
	}

	return nil
}
