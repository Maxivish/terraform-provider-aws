package connect_test

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
	// TIP: ==== IMPORT ====
	// This is a common set of imports but not customized to your code
	// since your code hasn't been written yet. Make sure you, your IDE, or
	// goimports -w <file> fixes these imports.
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfconnect "github.com/hashicorp/terraform-provider-aws/internal/service/connect"
)

// TIP: File Structure. The basic outline for all test files should be as
// follows. Improve this resource's maintainability by following this
// outline.
//
// 1. Package declaration (add "_test" since this is a test file)
// 2. Imports
// 3. Unit tests
// 4. Basic test
// 5. Disappears test
// 6. All the other tests
// 7. Helper functions (exists, destroy, check, etc.)
// 8. Functions that return Terraform configurations

func TestAccConnectInstanceStorageConfig_basic(t *testing.T) {
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var instanceStorageConfig connect.InstanceStorageConfig
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_connect_instance_storage_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, connect.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckInstanceStorageConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceStorageConfigConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceStorageConfigExists(resourceName, &instanceStorageConfig),
					resource.TestCheckResourceAttrSet(resourceName, "association_id"),
					resource.TestCheckResourceAttr(resourceName, "storage_config.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apply_immediately", "user"},
			},
		},
	})
}

func TestAccConnectInstanceStorageConfig_disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var instancestorageconfig connect.InstanceStorageConfig
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_connect_instance_storage_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(connect.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, connect.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckInstanceStorageConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceStorageConfigConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceStorageConfigExists(resourceName, &instancestorageconfig),
					acctest.CheckResourceDisappears(acctest.Provider, tfconnect.ResourceInstanceStorageConfig(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckInstanceStorageConfigDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ConnectConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_connect_instance_storage_config" {
			continue
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		associationId := rs.Primary.ID
		resourceType := rs.Primary.Attributes["resource_type"]
		_, err := tfconnect.FindInstanceStorageConfigByIDAndType(context.Background(), conn, instanceId, associationId, resourceType)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, connect.ErrCodeResourceNotFoundException) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Expected Connect InstanceStorageConfig to be destroyed, %s found", rs.Primary.ID)
	}

	return nil
}

func testAccCheckInstanceStorageConfigExists(name string, instanceStorageConfig *connect.InstanceStorageConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Connect InstanceStorageConfig is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ConnectConn

		instanceId := rs.Primary.Attributes["instance_id"]
		associationId := rs.Primary.ID
		resourceType := rs.Primary.Attributes["resource_type"]
		resp, err := tfconnect.FindInstanceStorageConfigByIDAndType(context.Background(), conn, instanceId, associationId, resourceType)

		if err != nil {
			return fmt.Errorf("Error describing Connect InstanceStorageConfig: %s", err.Error())
		}

		*instanceStorageConfig = *resp

		return nil
	}
}

func testAccPreCheck(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ConnectConn

	input := &connect.ListInstanceStorageConfigsInput{}

	_, err := conn.ListInstanceStorageConfigs(input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

// func testAccCheckInstanceStorageConfigNotRecreated(before, after *connect.InstanceStorageConfig) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		if before, after := aws.StringValue(before.InstanceStorageConfigId), aws.StringValue(after.InstanceStorageConfigId); before != after {
// 			return fmt.Errorf("Connect InstanceStorageConfig (%s/%s) recreated", before, after)
// 		}

// 		return nil
// 	}
// }

func testAccInstanceStorageConfigConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_connect_instance_storage_config" "test" {
//   instance_id   = data.aws_connect_instance.test.instance_id
  instance_id   = "bfef786e-63e0-4dae-a88d-629f667f8538"
  resource_type = "CHAT_TRANSCRIPTS"

  storage_config {
    s3_config {
		bucket_name = aws_s3_bucket.test.bucket
		bucket_prefix = "tf-test-Chat-Transcripts"
	}
	storage_type = "S3"
  }
}

// data "aws_connect_instance" "test" {
//   instance_id = "bfef786e-63e0-4dae-a88d-629f667f8538"
// }

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}
`, rName)
}
