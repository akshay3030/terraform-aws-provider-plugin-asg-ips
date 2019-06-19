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
	//client := m.(*MyClient)

	// Attempt to read from an upstream API
	//obj, ok := client.Get(d.Id())

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	//if !ok {
	//	d.SetId("")
	//	return nil
	//}

	//d.Set("address", obj.Address)
	return nil
}

func resourceAwsasgipsUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	if d.HasChange("address") {
		if err := updateAddress(d, m); err != nil {
			return err
		}
		d.SetPartial("address")
	}

	d.Partial(false)
	return resourceAwsasgipsRead(d, m)
}

func resourceAwsasgipsDelete(d *schema.ResourceData, m interface{}) error {
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")
	return nil
}

func updateAddress(d *schema.ResourceData, m interface{}) error {
	return nil
}
