package main

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["access_key"],
			},

			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["secret_key"],
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["profile"],
			},

			"assume_role": assumeRoleSchema(),

			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["shared_credentials_file"],
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["token"],
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				Description:  descriptions["region"],
				InputDefault: "us-east-1",
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     25,
				Description: descriptions["max_retries"],
			},

			"allowed_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"forbidden_account_ids"},
				Set:           schema.HashString,
			},
			"forbidden_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"allowed_account_ids"},
				Set:           schema.HashString,
			},

			"dynamodb_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["dynamodb_endpoint"],
				Removed:     "Use `dynamodb` inside `endpoints` block instead",
			},

			"kinesis_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["kinesis_endpoint"],
				Removed:     "Use `kinesis` inside `endpoints` block instead",
			},

			"endpoints": endpointsSchema(),

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},

			"skip_credentials_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_credentials_validation"],
			},

			"skip_get_ec2_platforms": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_get_ec2_platforms"],
			},

			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_region_validation"],
			},

			"skip_requesting_account_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_requesting_account_id"],
			},

			"skip_metadata_api_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_metadata_api_check"],
			},
			"s3_force_path_style": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["s3_force_path_style"],
			},
		},

		//ResourcesMap: map[string]*schema.Resource{
		//	"awsasgips_provider": resourceAwsasgips(),
		//},
		DataSourcesMap: map[string]*schema.Resource{
			"awsasgips": dataSourceAwsasgips(),
		},
		ConfigureFunc: provideConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{}
}

func provideConfigure(d *schema.ResourceData) (interface{}, error) {

	config := Config{
		AccessKey:               d.Get("access_key").(string),
		SecretKey:               d.Get("secret_key").(string),
		Profile:                 d.Get("profile").(string),
		CredsFilename:           d.Get("shared_credentials_file").(string),
		Token:                   d.Get("token").(string),
		Region:                  d.Get("region").(string),
		MaxRetries:              d.Get("max_retries").(int),
		Insecure:                d.Get("insecure").(bool),
		SkipCredsValidation:     d.Get("skip_credentials_validation").(bool),
		SkipGetEC2Platforms:     d.Get("skip_get_ec2_platforms").(bool),
		SkipRegionValidation:    d.Get("skip_region_validation").(bool),
		SkipRequestingAccountId: d.Get("skip_requesting_account_id").(bool),
		SkipMetadataApiCheck:    d.Get("skip_metadata_api_check").(bool),
		S3ForcePathStyle:        d.Get("s3_force_path_style").(bool),
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		config.AssumeRoleARN = assumeRole["role_arn"].(string)
		config.AssumeRoleSessionName = assumeRole["session_name"].(string)
		config.AssumeRoleExternalID = assumeRole["external_id"].(string)

		if v := assumeRole["policy"].(string); v != "" {
			config.AssumeRolePolicy = v
		}

		log.Printf("[INFO] assume_role configuration set: (ARN: %q, SessionID: %q, ExternalID: %q, Policy: %q)",
			config.AssumeRoleARN, config.AssumeRoleSessionName, config.AssumeRoleExternalID, config.AssumeRolePolicy)
	} else {
		log.Printf("[INFO] No assume_role block read from configuration")
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		config.CloudFormationEndpoint = endpoints["cloudformation"].(string)
		config.CloudWatchEndpoint = endpoints["cloudwatch"].(string)
		config.CloudWatchEventsEndpoint = endpoints["cloudwatchevents"].(string)
		config.CloudWatchLogsEndpoint = endpoints["cloudwatchlogs"].(string)
		config.DeviceFarmEndpoint = endpoints["devicefarm"].(string)
		config.DynamoDBEndpoint = endpoints["dynamodb"].(string)
		config.Ec2Endpoint = endpoints["ec2"].(string)
		config.ElbEndpoint = endpoints["elb"].(string)
		config.IamEndpoint = endpoints["iam"].(string)
		config.KinesisEndpoint = endpoints["kinesis"].(string)
		config.KmsEndpoint = endpoints["kms"].(string)
		config.RdsEndpoint = endpoints["rds"].(string)
		config.S3Endpoint = endpoints["s3"].(string)
		config.SnsEndpoint = endpoints["sns"].(string)
		config.SqsEndpoint = endpoints["sqs"].(string)
	}

	if v, ok := d.GetOk("allowed_account_ids"); ok {
		config.AllowedAccountIds = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("forbidden_account_ids"); ok {
		config.ForbiddenAccountIds = v.(*schema.Set).List()
	}

	return config.Client()
}

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_role_arn"],
				},

				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
				},

				"external_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_external_id"],
				},

				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_policy"],
				},
			},
		},
		Set: assumeRoleToHash,
	}
}

func assumeRoleToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["role_arn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["session_name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["external_id"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["policy"].(string)))
	return hashcode.String(buf.String())
}

func endpointsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cloudwatch": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudwatch_endpoint"],
				},
				"cloudwatchevents": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudwatchevents_endpoint"],
				},
				"cloudwatchlogs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudwatchlogs_endpoint"],
				},
				"cloudformation": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudformation_endpoint"],
				},
				"devicefarm": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["devicefarm_endpoint"],
				},
				"dynamodb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dynamodb_endpoint"],
				},
				"iam": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["iam_endpoint"],
				},

				"ec2": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ec2_endpoint"],
				},

				"elb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["elb_endpoint"],
				},
				"kinesis": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["kinesis_endpoint"],
				},
				"kms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["rds_endpoint"],
				},
				"s3": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["s3_endpoint"],
				},
				"sns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sns_endpoint"],
				},
				"sqs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sqs_endpoint"],
				},
			},
		},
		Set: endpointsToHash,
	}
}

func endpointsToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["cloudwatch"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cloudwatchevents"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cloudwatchlogs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cloudformation"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["devicefarm"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dynamodb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["iam"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ec2"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["elb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kinesis"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["rds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["s3"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sqs"].(string)))

	return hashcode.String(buf.String())
}
