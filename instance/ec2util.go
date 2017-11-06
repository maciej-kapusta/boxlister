package instance

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
)

func DescribeInstances(profile *string) []*Instance {
	options := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if *profile != "" {
		options.Profile = *profile
	}
	sess := session.Must(session.NewSessionWithOptions(options))
	ec2Sess := ec2.New(sess)
	instances, error := ec2Sess.DescribeInstances(&ec2.DescribeInstancesInput{})
	if error != nil {
		panic(error)
	}

	ec2Machines := []*Instance{}
	for _, r := range instances.Reservations {
		for _, instance := range r.Instances {
			ec2Machine := ToE2Machine(instance)
			if ec2Machine.IsValid() {
				ec2Machines = append(ec2Machines, ec2Machine)
			}
		}
	}
	return ec2Machines
}
