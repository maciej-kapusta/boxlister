package instance

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"fmt"
	"strings"
)

type Instance struct {
	Name    *string
	DnsName *string
}

func ToE2Machine(intance *ec2.Instance) *Instance {
	return &Instance{
		DnsName: publicDns(intance),
		Name:    nameTag(intance),
	}
}

func (i *Instance) String() string {
	return fmt.Sprintf("%v:%v", *i.Name, *i.DnsName)
}

func (i *Instance) IsValid() bool {
	return i.DnsName != nil && i.Name != nil
}

func (i *Instance) NameContains(namePart string) bool {
	return i.Name != nil && strings.Contains(*i.Name, namePart)
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
