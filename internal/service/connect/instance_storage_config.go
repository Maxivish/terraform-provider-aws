package connect

// **PLEASE DELETE THIS AND ALL TIP COMMENTS BEFORE SUBMITTING A PR FOR REVIEW!**
//
// TIP: ==== INTRODUCTION ====
// Thank you for trying the skaff tool!
//
// You have opted to include these helpful comments. They all include "TIP:"
// to help you find and remove them when you're done with them.
//
// While some aspects of this file are customized to your input, the
// scaffold tool does *not* look at the AWS API and ensure it has correct
// function, structure, and variable names. It makes guesses based on
// commonalities. You will need to make significant adjustments.
//
// In other words, as generated, this is a rough outline of the work you will
// need to do. If something doesn't make sense for your situation, get rid of
// it.
//
// Remember to register this new resource in the provider
// (internal/provider/provider.go) once you finish. Otherwise, Terraform won't
// know about it.

import (
	// TIP: ==== IMPORTS ====
	// This is a common set of imports but not customized to your code since
	// your code hasn't been written yet. Make sure you, your IDE, or
	// goimports -w <file> fixes these imports.
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceInstanceStorageConfig() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceInstanceStorageConfigCreate,
		ReadWithoutTimeout:   resourceInstanceStorageConfigRead,
		UpdateWithoutTimeout: resourceInstanceStorageConfigUpdate,
		DeleteWithoutTimeout: resourceInstanceStorageConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// TIP: ==== CONFIGURABLE TIMEOUTS ====
		// Users can configure timeout lengths (if you use the times they
		// provide). These are the defaults if they don't configure timeouts.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		// TIP: ==== SCHEMA ====
		// In the schema, add each of the arguments and attributes in snake
		// case (e.g., delete_automated_backups).
		// * Alphabetize arguments to make them easier to find.
		// * Do not add a blank line between arguments/attributes.
		//
		// Users can configure argument values while attribute values cannot be
		// configured and are used as output. Arguments have either:
		// Required: true,
		// Optional: true,
		//
		// All attributes will be computed and some arguments. If users will
		// want to read updated information or detect drift for an argument,
		// it should be computed:
		// Computed: true,
		//
		// You will typically find arguments in the input struct
		// (e.g., CreateDBInstanceInput) for the create operation. Sometimes
		// they are only in the input struct (e.g., ModifyDBInstanceInput) for
		// the modify operation.
		//
		// For more about schema options, visit
		// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#Schema
		Schema: map[string]*schema.Schema{
			"association_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(connect.InstanceStorageResourceType_Values(), false),
			},
			"storage_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// TODO: Implement Kinesis fields
						"s3_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"bucket_prefix": {
										Type:     schema.TypeString,
										Required: true,
									},
									"encryption_config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"encryption_type": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(connect.EncryptionType_Values(), false),
												},
												"key_id": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: verify.ValidARN,
												},
											},
										},
									},
								},
							},
						},
						"storage_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(connect.StorageType_Values(), false),
						},
					},
				},
			},
		},
	}
}

func resourceInstanceStorageConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectConn

	in := &connect.AssociateInstanceStorageConfigInput{
		InstanceId:    aws.String(d.Get("instance_id").(string)),
		ResourceType:  aws.String(d.Get("resource_type").(string)),
		StorageConfig: expandStorageConfig(d.Get("storage_config")),
	}

	var out *connect.AssociateInstanceStorageConfigOutput
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		out, err = conn.AssociateInstanceStorageConfigWithContext(ctx, in)

		if tfawserr.ErrCodeEquals(err, connect.ErrCodeAccessDeniedException) {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if tfresource.TimedOut(err) {
		out, err = conn.AssociateInstanceStorageConfigWithContext(ctx, in)
	}

	if err != nil {
		return diag.Errorf("creating Amazon Connect InstanceStorageConfig (%s,%s): %s", d.Get("instance_id").(string), d.Get("resource_type").(string), err)
	}

	if out == nil || out.AssociationId == nil {
		return diag.Errorf("creating Amazon Connect InstanceStorageConfig (%s,%s): empty output", d.Get("instance_id").(string), d.Get("resource_type").(string))
	}

	d.SetId(aws.StringValue(out.AssociationId))

	return resourceInstanceStorageConfigRead(ctx, d, meta)
}

func resourceInstanceStorageConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectConn

	instanceId := d.Get("instance_id").(string)
	associationId := d.Id()
	resourceType := d.Get("resource_type").(string)
	out, err := FindInstanceStorageConfigByIDAndType(ctx, conn, instanceId, associationId, resourceType)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Connect InstanceStorageConfig (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading Connect InstanceStorageConfig (%s): %s", d.Id(), err)
	}
	d.Set("association_id", out.AssociationId)
	d.Set("resource_type", resourceType)
	if err := d.Set("storage_config", flattenStorageConfig(out)); err != nil {
		return diag.Errorf("setting storage_config: %s", err)
	}

	return nil
}

func resourceInstanceStorageConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectConn

	if !d.HasChanges("storage_config") {
		return nil
	}

	in := &connect.UpdateInstanceStorageConfigInput{
		AssociationId: aws.String(d.Id()),
		InstanceId:    aws.String(d.Get("instance_id").(string)),
		ResourceType:  aws.String(d.Get("resource_type").(string)),
		StorageConfig: expandStorageConfig(d.Get("storage_config")),
	}

	_, err := conn.UpdateInstanceStorageConfigWithContext(ctx, in)
	if err != nil {
		return diag.Errorf("updating Connect InstanceStorageConfig (%s): %s", d.Id(), err)
	}

	return resourceInstanceStorageConfigRead(ctx, d, meta)
}

func resourceInstanceStorageConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectConn

	in := &connect.DisassociateInstanceStorageConfigInput{
		AssociationId: aws.String(d.Id()),
		InstanceId:    aws.String(d.Get("instance_id").(string)),
		ResourceType:  aws.String(d.Get("resource_type").(string)),
	}

	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := conn.DisassociateInstanceStorageConfigWithContext(ctx, in)

		if tfawserr.ErrCodeEquals(err, connect.ErrCodeAccessDeniedException) {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.DisassociateInstanceStorageConfigWithContext(ctx, in)
	}

	if tfawserr.ErrCodeEquals(err, connect.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting Connect InstanceStorageConfig (%s): %s", d.Id(), err)
	}

	return nil
}

func FindInstanceStorageConfigByIDAndType(ctx context.Context, conn *connect.Connect, instanceId, associationId, resourceType string) (*connect.InstanceStorageConfig, error) {
	in := &connect.DescribeInstanceStorageConfigInput{
		AssociationId: aws.String(associationId),
		InstanceId:    aws.String(instanceId),
		ResourceType:  aws.String(resourceType),
	}

	out, err := conn.DescribeInstanceStorageConfigWithContext(ctx, in)
	if tfawserr.ErrCodeEquals(err, connect.ErrCodeResourceNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: in,
		}
	}

	if err != nil {
		return nil, err
	}

	if out == nil || out.StorageConfig == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out.StorageConfig, nil
}

// TIP: Even when you have a list with max length of 1, this plural function
// works brilliantly. However, if the AWS API takes a structure rather than a
// slice of structures, you will not need it.
func expandStorageConfig(v interface{}) *connect.InstanceStorageConfig {
	l := v.([]interface{})
	if len(l) != 1 {
		return nil
	}

	m := l[0].(map[string]interface{})

	storageType := m["storage_type"].(string)

	result := &connect.InstanceStorageConfig{
		StorageType: aws.String(storageType),
	}

	switch storageType {
	case connect.StorageTypeS3:
		result.S3Config = exapandS3Config(m["s3_config"])
	}

	return result
}

func flattenStorageConfig(apiObject *connect.InstanceStorageConfig) []interface{} {
	if apiObject == nil {
		return nil
	}

	storageType := aws.StringValue(apiObject.StorageType)

	m := map[string]interface{}{
		"storage_type": storageType,
	}

	switch storageType {
	case connect.StorageTypeS3:
		m["s3_config"] = flattenS3Config(apiObject.S3Config)
	}

	return []interface{}{m}
}

func exapandS3Config(v interface{}) *connect.S3Config {
	l := v.([]interface{})
	if len(l) != 1 {
		return nil
	}

	m := l[0].(map[string]interface{})

	result := &connect.S3Config{
		BucketName:   aws.String(m["bucket_name"].(string)),
		BucketPrefix: aws.String(m["bucket_prefix"].(string)),
	}

	if m["encryption_config"] != nil {
		result.EncryptionConfig = expandEncryptionConfig(m["encryption_config"])
	}

	return result
}

func flattenS3Config(apiObject *connect.S3Config) []interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{
		"bucket_name":   aws.StringValue(apiObject.BucketName),
		"bucket_prefix": aws.StringValue(apiObject.BucketPrefix),
	}

	if apiObject.EncryptionConfig != nil {
		m["encryption_config"] = flattenEncryptionConfig(apiObject.EncryptionConfig)
	}

	return []interface{}{m}
}

func expandEncryptionConfig(v interface{}) *connect.EncryptionConfig {
	l := v.([]interface{})
	if len(l) != 1 {
		return nil
	}

	m := l[0].(map[string]interface{})

	result := &connect.EncryptionConfig{
		EncryptionType: aws.String(m["encryption_type"].(string)),
		KeyId:          aws.String(m["key_id"].(string)),
	}

	return result
}

func flattenEncryptionConfig(apiObject *connect.EncryptionConfig) []interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{
		"encryption_type": apiObject.EncryptionType,
		"key_id":          apiObject.KeyId,
	}

	return []interface{}{m}
}
