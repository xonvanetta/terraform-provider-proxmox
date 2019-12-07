/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkResourceVirtualEnvironmentAccessGroupComment = "comment"
	mkResourceVirtualEnvironmentAccessGroupID      = "group_id"
	mkResourceVirtualEnvironmentAccessGroupMembers = "members"
)

func resourceVirtualEnvironmentAccessGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentAccessGroupComment: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The group comment",
				Optional:    true,
				Default:     "",
			},
			mkResourceVirtualEnvironmentAccessGroupID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The group id",
				Required:    true,
			},
			mkResourceVirtualEnvironmentAccessGroupMembers: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The group members",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Create: resourceVirtualEnvironmentAccessGroupCreate,
		Read:   resourceVirtualEnvironmentAccessGroupRead,
		Update: resourceVirtualEnvironmentAccessGroupUpdate,
		Delete: resourceVirtualEnvironmentAccessGroupDelete,
	}
}

func resourceVirtualEnvironmentAccessGroupCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)

	if config.veClient == nil {
		return errors.New("You must specify the virtual environment details in the provider configuration to use this data source")
	}

	groupID := d.Get(mkResourceVirtualEnvironmentAccessGroupID).(string)
	body := &proxmox.VirtualEnvironmentAccessGroupCreateRequestBody{
		Comment: d.Get(mkResourceVirtualEnvironmentAccessGroupComment).(string),
		ID:      groupID,
	}

	err := config.veClient.CreateAccessGroup(body)

	if err != nil {
		return err
	}

	d.SetId(groupID)

	return resourceVirtualEnvironmentAccessGroupRead(d, m)
}

func resourceVirtualEnvironmentAccessGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)

	if config.veClient == nil {
		return errors.New("You must specify the virtual environment details in the provider configuration to use this data source")
	}

	groupID := d.Get(mkResourceVirtualEnvironmentAccessGroupID).(string)
	accessGroup, err := config.veClient.GetAccessGroup(groupID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId(groupID)

	d.Set(mkResourceVirtualEnvironmentAccessGroupComment, accessGroup.Comment)
	d.Set(mkResourceVirtualEnvironmentAccessGroupMembers, accessGroup.Members)

	return nil
}

func resourceVirtualEnvironmentAccessGroupUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)

	if config.veClient == nil {
		return errors.New("You must specify the virtual environment details in the provider configuration to use this data source")
	}

	body := &proxmox.VirtualEnvironmentAccessGroupUpdateRequestBody{
		Comment: d.Get(mkResourceVirtualEnvironmentAccessGroupComment).(string),
	}

	groupID := d.Get(mkResourceVirtualEnvironmentAccessGroupID).(string)
	err := config.veClient.UpdateAccessGroup(groupID, body)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentAccessGroupRead(d, m)
}

func resourceVirtualEnvironmentAccessGroupDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)

	if config.veClient == nil {
		return errors.New("You must specify the virtual environment details in the provider configuration to use this data source")
	}

	groupID := d.Get(mkResourceVirtualEnvironmentAccessGroupID).(string)
	err := config.veClient.DeleteAccessGroup(groupID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId("")

	return nil
}