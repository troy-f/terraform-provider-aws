package ecr_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestAccAWSEcrAuthorizationTokenDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	dataSourceName := "data.aws_ecr_authorization_token.repo"

	resource.Test(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ecr.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsEcrAuthorizationTokenDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "authorization_token"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "expires_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "user_name"),
					resource.TestMatchResourceAttr(dataSourceName, "user_name", regexp.MustCompile(`AWS`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "password"),
				),
			},
			{
				Config: testAccCheckAwsEcrAuthorizationTokenDataSourceRepositoryConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "registry_id", "aws_ecr_repository.repo", "registry_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorization_token"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "expires_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "user_name"),
					resource.TestMatchResourceAttr(dataSourceName, "user_name", regexp.MustCompile(`AWS`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "password"),
				),
			},
		},
	})
}

var testAccCheckAwsEcrAuthorizationTokenDataSourceBasicConfig = `
data "aws_ecr_authorization_token" "repo" {}
`

func testAccCheckAwsEcrAuthorizationTokenDataSourceRepositoryConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_ecr_repository" "repo" {
  name = %q
}

data "aws_ecr_authorization_token" "repo" {
  registry_id = aws_ecr_repository.repo.registry_id
}
`, rName)
}
