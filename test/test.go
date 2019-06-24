package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"net/http"
	"time"
)

func main() {
	//x := 7
	//if x > 6 {
	//	fmt.Println("More than 6")
	//}

	//var inf interface{} // interface declaration

	//var ipList []string
	//var instanceList []string

	//var instanceInfo map[string][]string //this was resulting in "panic: assignment to entry in nil map"
	//created map of list
	instanceInfo := make(map[string][]string)

	svc := newAwsAsgService("us-west-2", "default")
	svcec2 := newEc2Service("us-west-2", "default")

	asg, err := describeScalingGroup("asg-green-dev-media20190314181204859800000009", svc)
	//fmt.Println(err)
	//fmt.Println(asg)
	//fmt.Println(asg.AutoScalingGroups)

	if err != nil {
		fmt.Println(err)
	}

	for _, autoScalingGroup := range asg.AutoScalingGroups {
		//fmt.Println(autoScalingGroup.Instances, autoScalingGroup.LaunchConfigurationName)
		//fmt.Println(autoScalingGroup.Instances)
		for _, Instance := range autoScalingGroup.Instances {

			// print the instance id
			//fmt.Println("Instance Id is --> ",*Instance.InstanceId)
			//instanceList = append(instanceList,*Instance.InstanceId)
			instanceInfo["id"] = append(instanceInfo["id"], *Instance.InstanceId)

			ins, err1 := describeEc2(aws.StringSlice([]string{*Instance.InstanceId}), svcec2)
			if err1 != nil {
				fmt.Println(err1)
			}

			//print the private ip
			//fmt.Println("Private IP  is --> ", *ins.Reservations[0].Instances[0].PrivateIpAddress)
			//var ip string = *ins.Reservations[0].Instances[0].PrivateIpAddress
			//inf = ins.Reservations[0].Instances[0].PublicIpAddress
			private_ip := *ins.Reservations[0].Instances[0].PrivateIpAddress

			//fmt.Println("ip address is",ip)
			if ins.Reservations[0].Instances[0].PublicIpAddress != nil {
				instanceInfo["public_ip"] = append(instanceInfo["public_ip"], *ins.Reservations[0].Instances[0].PublicIpAddress)

			}

			instanceInfo["private_ip"] = append(instanceInfo["private_ip"], private_ip)
			//fmt.Println(iplist)

		}

	}

	//instanceIdss := []string{"i-0d828ca28b40a2b27", "i-0fc4733bf128e2244"}
	//instanceIdss := []string{"i-0d828ca28b40a2b27"}

	//ins, err1 := describeEc2(aws.StringSlice(instanceIdss), svcec2)
	//fmt.Println(ins)
	//fmt.Println(*ins.Reservations[0].Instances[0].PrivateIpAddress)
	//fmt.Println(err1)

	//return []*autoscaling.Instance{}, nil, errors.New("asg not found")

	//for idx, res := range asg {
	//	fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
	//	for _, inst := range resp.Reservations[idx].Instances {
	//		fmt.Println("    - Instance ID: ", *inst.InstanceId)
	//	}

	fmt.Println(instanceInfo)

}

// newAwsAsgService returns a session object for the AWS autoscaling service.
func newAwsAsgService(region string, profile string) (Session *autoscaling.AutoScaling) {
	sess := session.Must(session.NewSession())
	//svc := autoscaling.New(sess, config(region, os.Getenv("ASG_ID"), os.Getenv("ASG_SECRET")))
	svc := autoscaling.New(sess, config(region, profile))

	return svc
}

// newAwsAsgService returns a session object for the AWS autoscaling service.
func newEc2Service(region string, profile string) (Session *ec2.EC2) {
	sessec2 := session.Must(session.NewSession())
	svcec2 := ec2.New(sessec2, config(region, profile))
	return svcec2
}

// Config produces a generic set of AWS configs
func config(region, profile string) *aws.Config {
	return aws.NewConfig().
		//WithCredentials(credentials.NewStaticCredentials(id, secret, "")).
		WithCredentials(credentials.NewSharedCredentials("", profile)).
		WithRegion(region).
		WithHTTPClient(http.DefaultClient).
		WithMaxRetries(aws.UseServiceDefaultRetries).
		WithLogger(aws.NewDefaultLogger()).
		WithLogLevel(aws.LogOff).
		WithSleepDelay(time.Sleep).
		WithEndpointResolver(endpoints.DefaultResolver())
}

// describeScalingGroup
func describeScalingGroup(asgName string,
	svc *autoscaling.AutoScaling) (
	asg *autoscaling.DescribeAutoScalingGroupsOutput, err error) {

	params := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(asgName),
		},
	}
	resp, err := svc.DescribeAutoScalingGroups(params)
	if err != nil {
		fmt.Println(err)
	}

	// If we failed to get exactly one ASG, raise an error.
	//if len(resp.AutoScalingGroups) != 1 {
	//	err = fmt.Errorf("the attempt to retrieve the current worker pool "+
	//		"autoscaling group configuration expected exaclty one result got %v",
	//		len(resp.AutoScalingGroups))
	//}

	return resp, err
}

// DescribeInstances returns a list of Instances, given a list of instance IDs.
//func (e EC2) DescribeInstances(instanceIds []string) ([]types.Instance, error) {
//	params := &ec2.DescribeInstancesInput{
//		InstanceIds: []*string{},
//		MaxResults:  aws.Int64(int64(len(instanceIds))),
//	}
//	for _, id := range instanceIds {
//		params.InstanceIds = append(params.InstanceIds, aws.String(id))
//	}
//	return e.describeInstances(params)
//}

// describeEc2
func describeEc2(instanceIds []*string, svcec2 *ec2.EC2) (ec2out *ec2.DescribeInstancesOutput, err error) {

	params := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
		//MaxResults:  aws.Int64(int64(len(instanceIds))),
	}

	resp, err := svcec2.DescribeInstances(params)

	// If we failed to get exactly one ASG, raise an error.
	//if len(resp.AutoScalingGroups) != 1 {
	//	err = fmt.Errorf("the attempt to retrieve the current worker pool "+
	//		"autoscaling group configuration expected exaclty one result got %v",
	//		len(resp.AutoScalingGroups))
	//}

	return resp, err
}

//AWS Example Code from https://docs.aws.amazon.com/sdk-for-go/api/service/autoscaling/#AutoScaling.DescribeAutoScalingGroups

//svc := autoscaling.New(session.New())
//input := &autoscaling.DescribeAutoScalingGroupsInput{
//AutoScalingGroupNames: []*string{
//aws.String("my-auto-scaling-group"),
//},
//}
//
//result, err := svc.DescribeAutoScalingGroups(input)
//if err != nil {
//if aerr, ok := err.(awserr.Error); ok {
//switch aerr.Code() {
//case autoscaling.ErrCodeInvalidNextToken:
//fmt.Println(autoscaling.ErrCodeInvalidNextToken, aerr.Error())
//case autoscaling.ErrCodeResourceContentionFault:
//fmt.Println(autoscaling.ErrCodeResourceContentionFault, aerr.Error())
//default:
//fmt.Println(aerr.Error())
//}
//} else {
//// Print the error, cast err to awserr.Error to get the Code and
//// Message from an error.
//fmt.Println(err.Error())
//}
//return
//}
//
//fmt.Println(result)

//type AutoScalingGroups struct{
//AutoScalingGroupName: input.AutoScalingGroupName,
//AvailabilityZones:    input.AvailabilityZones,
//CreatedTime:          &createdTime,
//DefaultCooldown:      input.DefaultCooldown,
//DesiredCapacity:      input.DesiredCapacity,
//// EnabledMetrics:          input.EnabledMetrics,
//HealthCheckGracePeriod:  input.HealthCheckGracePeriod,
//HealthCheckType:         input.HealthCheckType,
//Instances:               []*autoscaling.Instance{},
//LaunchConfigurationName: input.LaunchConfigurationName,
//LoadBalancerNames:       input.LoadBalancerNames,
//MaxSize:                 input.MaxSize,
//MinSize:                 input.MinSize,
//NewInstancesProtectedFromScaleIn: input.NewInstancesProtectedFromScaleIn,
//PlacementGroup:                   input.PlacementGroup,
//// Status:                           input.Status,
//// SuspendedProcesses:               input.SuspendedProcesses,
//// Tags:                input.Tags,
//TargetGroupARNs:     input.TargetGroupARNs,
//TerminationPolicies: input.TerminationPolicies,
//VPCZoneIdentifier:   input.VPCZoneIdentifier,
//}

//func listASGInstaces(ASAPI autoscalingiface.AutoScalingAPI, asgName string) ([]*autoscaling.Instance, *string, error) {
//	output, err := ASAPI.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
//		AutoScalingGroupNames: []*string{aws.String(asgName)},
//	})
//	if err != nil {
//		return []*autoscaling.Instance{}, nil, err
//	}
//
//	for _, autoScalingGroup := range output.AutoScalingGroups {
//		return autoScalingGroup.Instances, autoScalingGroup.LaunchConfigurationName, nil
//	}
//
//	return []*autoscaling.Instance{}, nil, errors.New("asg not found")
//}
