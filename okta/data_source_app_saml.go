package okta

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAppSaml() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppSamlRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
			},
			"label": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
			},
			"label_prefix": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
			},
			"active_only": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppSamlRead(d *schema.ResourceData, m interface{}) error {
	id := d.Get("id").(string)
	label := d.Get("label").(string)
	labelPrefix := d.Get("label_prefix").(string)
	filters := &appFilters{ID: id, Label: label, LabelPrefix: labelPrefix}

	if d.Get("active_only").(bool) {
		filters.ApiFilter = `status eq "ACTIVE"`
	}

	if id == "" && label == "" && labelPrefix == "" {
		return errors.New("you must provide either an label_prefix, id, or label to search with")
	}
	appList, err := listApps(m.(*Config), filters)
	if err != nil {
		return err
	}
	if len(appList) < 1 {
		return fmt.Errorf(`No application found with provided filter. id: "%s", label: "%s", label_prefix: "%s"`, id, label, labelPrefix)
	} else if len(appList) > 1 {
		fmt.Println("Found multiple applications with the criteria supplied, using the first one, sorted by creation date.")
	}
	app := appList[0]
	d.SetId(app.ID)
	d.Set("label", app.Label)
	d.Set("description", app.Description)
	d.Set("name", app.Name)
	d.Set("status", app.Status)

	return nil
}