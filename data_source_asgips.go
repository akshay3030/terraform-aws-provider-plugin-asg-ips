package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rs/xid"
	"log"
	"net/http"
	"time"
)

func dataSourceAwsasgips() *schema.Resource {
	log.Println("[INFO] $$$$$$$$$$$$$$$$$$$$ we are here1********************************************")

	return &schema.Resource{
		Read: dataSourceAwsasgipsRead,

		Schema: map[string]*schema.Schema{
			"asgname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			//"output": &schema.Schema{
			//	Type:     schema.TypeMap,
			//	//Required: true,
			//	Computed: true,
			//	Elem:     schema.TypeList,
			//
			//},
			"private_ip": &schema.Schema{
				Type: schema.TypeList,
				//Required: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ip": &schema.Schema{
				Type: schema.TypeList,
				//Required: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_id": &schema.Schema{
				Type: schema.TypeList,
				//Required: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			//"ip": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
		},
	}
}

func dataSourceAwsasgipsRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[INFO] we are here2")

	instanceInfo := make(map[string][]string)

	asgname := d.Get("asgname").(string)
	log.Println("[INFO] asgname:-->", asgname)

	//region := d.Get("region").(string)

	//svc := newAwsAsgService(region, "default")
	//svcec2 := newEc2Service(region, "default")

	svc := m.(*AWSClient).autoscalingconn
	svcec2 := m.(*AWSClient).ec2conn

	asg, err := describeScalingGroup(asgname, svc)
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
			instanceInfo["instance_id"] = append(instanceInfo["instance_id"], *Instance.InstanceId)

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

	log.Println("[INFO] instanceInfo:-->", instanceInfo)

	//d.Set("output",instanceInfo)
	d.Set("instance_id", instanceInfo["instance_id"])
	d.Set("private_ip", instanceInfo["private_ip"])
	d.Set("public_ip", instanceInfo["public_ip"])
	//d.Set("ip",instanceInfo["private_ip"][0])

	log.Println("[INFO] d-full is :-->", d)
	//log.Println("[INFO] d-output is :-->",d.Get("output"))
	log.Println("[INFO] id-is :-->", d.Get("id"))
	log.Println("[INFO] private_ip-is :-->", d.Get("private_ip"))
	log.Println("[INFO] public_ip-is :-->", d.Get("public_ip"))
	//log.Println("[INFO] ip1-is :-->",d.Get("ip"))

	//create random uuid for the id, without this object won't get mapped in the data resources
	id := xid.New().String()
	d.SetId(id)
	return nil
}

// newAwsAsgService returns a session object for the AWS autoscaling service.
//func newAwsAsgService(region string, profile string) (Session *autoscaling.AutoScaling) {
//	sess := session.Must(session.NewSession())
//	//svc := autoscaling.New(sess, config(region, os.Getenv("ASG_ID"), os.Getenv("ASG_SECRET")))
//	svc := autoscaling.New(sess, config(region, profile))
//
//	return svc
//}

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
