package simplivity

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSimplivityVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceSimplivityVirtualMachineCreate,
		Read:   resourceSimplivityVirtualMachineRead,
		Update: resourceSimplivityVirtualMachineUpdate,
		Delete: resourceSimplivityVirtualMachineDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"power_state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (warning []string, errors []error) {
					val := v.(string)
					if val != "on" && val != "off" {
						errors = append(errors, fmt.Errorf("%q must be 'on' or 'off'", k))
					}
					return
				},
			},
		},
	}
}

func resourceSimplivityVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	d.SetId(name)
	return resourceSimplivityVirtualMachineRead(d, meta)
}

func resourceSimplivityVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).Client
	name := d.Id()

	log.Printf("[DEBUG] Reading Virtual Machine: %s", name)
	vm, err := client.VirtualMachines.GetByName(name)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("name", vm.Name)

	return nil
}

func resourceSimplivityVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).Client
	name := d.Get("name").(string)

	vm, err := client.VirtualMachines.GetByName(name)
	if err != nil {
		return err
	}

	if val, ok := d.GetOk("power_state"); ok {
		vm.UpdatePowerState(val.(string))
	}

	return resourceSimplivityVirtualMachineRead(d, meta)
}

func resourceSimplivityVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] VM delete endpoint doesnt exist")
	return nil
}