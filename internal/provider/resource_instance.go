package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Sotoon Compute Instance.",
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // The name cannot be changed after creation
				Description: "The unique name of the instance.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // Changing the instance type requires a new instance
				Description: "The instance type, e.g., 'c1d30-1'.",
			},
			"image": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // Changing the OS image requires a new instance
				Description: "The OS image for the instance, e.g., 'ubuntu-22.04'.",
			},
			"network_interface": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the network interface, e.g., 'ens3'.",
						},
						"link": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the Link resource to attach.",
						},
					},
				},
			},
			"powered_on": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the instance should be powered on. This can be updated.",
			},
			"iam_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether IAM is enabled for the instance. This can be updated.",
			},
			// Computed values
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier (UID) of the instance.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The workspace ID.",
			},
		},
	}
}

// buildInstanceSpec is a helper to construct the spec from schema data
func buildInstanceSpec(d *schema.ResourceData) (client.InstanceSpec, error) {
	interfacesSpec := d.Get("network_interface").([]interface{})
	var apiInterfaces []struct {
		Link string `json:"link"`
		Name string `json:"name"`
	}
	for _, i := range interfacesSpec {
		ifaceMap, ok := i.(map[string]interface{})
		if !ok {
			return client.InstanceSpec{}, fmt.Errorf("could not parse network_interface block")
		}
		apiInterfaces = append(apiInterfaces, struct {
			Link string `json:"link"`
			Name string `json:"name"`
		}{
			Link: ifaceMap["link"].(string),
			Name: ifaceMap["name"].(string),
		})
	}

	spec := client.InstanceSpec{
		PoweredOn:  d.Get("powered_on").(bool),
		IAMEnabled: d.Get("iam_enabled").(bool),
		Type:       d.Get("type").(string),
		ImageSource: struct {
			Image string `json:"image"`
		}{
			Image: d.Get("image").(string),
		},
		Interfaces: apiInterfaces,
	}
	return spec, nil
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	name := d.Get("name").(string)

	spec, err := buildInstanceSpec(d)
	if err != nil {
		return diag.FromErr(err)
	}

	instance, err := c.CreateInstance(name, spec)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instance.Metadata.Name)

	return resourceInstanceRead(ctx, d, m)
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	instanceName := d.Id()

	instance, err := c.GetInstance(instanceName)
	if err != nil {
		return diag.FromErr(err)
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", instance.Metadata.Name)
	d.Set("uid", instance.Metadata.UID)
	d.Set("workspace_id", instance.Metadata.Workspace)
	d.Set("type", instance.Spec.Type)
	d.Set("image", instance.Spec.ImageSource.Image)
	d.Set("powered_on", instance.Spec.PoweredOn)
	d.Set("iam_enabled", instance.Spec.IAMEnabled)

	if len(instance.Spec.Interfaces) > 0 {
		interfaces := make([]map[string]interface{}, len(instance.Spec.Interfaces))
		for i, ni := range instance.Spec.Interfaces {
			interfaces[i] = map[string]interface{}{
				"name": ni.Name,
				"link": ni.Link,
			}
		}
		if err := d.Set("network_interface", interfaces); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	instanceName := d.Id()

	if d.HasChange("powered_on") || d.HasChange("iam_enabled") {
		spec, err := buildInstanceSpec(d)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = c.UpdateInstance(instanceName, spec)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceInstanceRead(ctx, d, m)
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	instanceName := d.Id()

	err := c.DeleteInstance(instanceName)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
