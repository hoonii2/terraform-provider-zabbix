package provider

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kgeroczi/go-zabbix-api"
)

// resourceUser terraform resource handler
func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Username",
				Required:     true,
			},
			"password": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Password",
				Required:     true,
				Sensitive:    true,
			},
		},
	}
}

// dataUser terraform data handler
func dataUser() *schema.Resource {
	return &schema.Resource{
		Read: dataUserRead,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "Username",
				Required:     true,
			},
		},
	}
}

// terraform user create function
func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*zabbix.API)

	item := zabbix.User{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	items := []zabbix.User{item}

	err := api.UsersCreate(items)

	if err != nil {
		return err
	}

	log.Trace("created User: %+v", items[0])

	d.SetId(items[0].UserID)

	return resourceUserRead(d, m)
}

// userRead terraform user read function
func userRead(d *schema.ResourceData, m interface{}, params zabbix.Params) error {
	api := m.(*zabbix.API)

	Users, err := api.UsersGet(params)

	if err != nil {
		return err
	}

	if len(Users) < 1 {
		d.SetId("")
		return nil
	}
	if len(Users) > 1 {
		return errors.New("multiple Users found")
	}
	t := Users[0]

	log.Debug("Got User: %+v", t)

	d.SetId(t.UserID)
	d.Set("username", t.Username)

	return nil
}

// dataUserRead terraform data resource read handler
func dataUserRead(d *schema.ResourceData, m interface{}) error {
	return userRead(d, m, zabbix.Params{
		"filter": map[string]interface{}{
			"name": d.Get("name"),
		},
	})
}

// resourceUserRead terraform resource read handler
func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	log.Debug("Lookup of User with id %s", d.Id())

	return userRead(d, m, zabbix.Params{
		"userids": d.Id(),
	})
}

// resourceUserUpdate terraform resource update handler
func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*zabbix.API)

	item := zabbix.User{
		UserID:   d.Id(),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	items := []zabbix.User{item}

	err := api.UsersUpdate(items)

	if err != nil {
		return err
	}

	return resourceUserRead(d, m)
}

// resourceUserDelete terraform resource delete handler
func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*zabbix.API)
	return api.UsersDeleteByIds([]string{d.Id()})
}
