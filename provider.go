package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsbase "github.com/hashicorp/aws-sdk-go-base"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sspinc/terraform-provider-credstash/credstash"
)

var _ terraform.ResourceProvider = provider()

const defaultAWSProfile = "default"

func provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"credstash_secret": dataSourceSecret(),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				Description: "The region where AWS operations will take place. Examples\n" +
					"are us-east-1, us-west-2, etc.",
			},
			"table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DynamoDB table where the secrets are stored.",
				Default:     "credential-store",
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     defaultAWSProfile,
				Description: "The profile that should be used to connect to AWS",
			},
			"assume_role": assumeRoleSchema(),
		},
		ConfigureFunc: providerConfig,
	}
}

func providerConfig(d *schema.ResourceData) (interface{}, error) {
	region := d.Get("region").(string)
	table := d.Get("table").(string)
	profile := d.Get("profile").(string)

	var sess *session.Session
	var err error
	if profile != defaultAWSProfile {
		log.Printf("[DEBUG] creating a session for profile: %s", profile)
		sess, err = session.NewSessionWithOptions(session.Options{
			Config:            aws.Config{Region: aws.String(region)},
			Profile:           profile,
			SharedConfigState: session.SharedConfigEnable,
		})
	} else if l, ok := d.Get("assume_role").([]interface{}); ok && len(l) > 0 && l[0] != nil {
		m := l[0].(map[string]interface{})
		log.Println("[DEBUG] creating a session with assume role")
		cfg := &awsbase.Config{Region: region}
		if v, ok := m["duration_seconds"].(int); ok && v != 0 {
			cfg.AssumeRoleDurationSeconds = v
		}
		if v, ok := m["role_arn"].(string); ok && v != "" {
			cfg.AssumeRoleARN = v
		}
		sess, err = awsbase.GetSession(cfg)
	} else {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(region)})
	}
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] configured credstash for table %s", table)
	return credstash.New(table, sess), nil
}

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"duration_seconds": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Seconds to restrict the assume role session duration.",
				},
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Amazon Resource Name of an IAM Role to assume prior to making API calls.",
				},
			},
		},
	}
}
