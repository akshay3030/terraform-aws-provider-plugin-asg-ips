package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsasgips() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsasgipsCreate,
		Read:   resourceAwsasgipsRead,
		Update: resourceAwsasgipsUpdate,
		Delete: resourceAwsasgipsDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAwsasgipsCreate(d *schema.ResourceData, m interface{}) error {
	address := d.Get("address").(string)
	d.SetId(address)
	return resourceAwsasgipsRead(d, m)
}

func resourceAwsasgipsRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAwsasgipsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAwsasgipsRead(d, m)
}

func resourceAwsasgipsDelete(d *schema.ResourceData, m interface{}) error {
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")
	return nil
}
