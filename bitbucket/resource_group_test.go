package bitbucket

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

func TestAccBitbucketResourceGroup_basic(t *testing.T) {
	config := `
		resource "bitbucketserver_group" "test" {
			name = "test-group"
		}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucketserver_group.test", "name", "test-group"),
				),
			},
		},
	})
}

func TestAccBitbucketResourceGroup_DisallowImport(t *testing.T) {
	resourceName := "duplicate_group"
	groupName := "duplicate-group"
	config := fmt.Sprintf(`
		resource "bitbucketserver_group" "%s" {
			name = "%s"
		}
	`, resourceName, groupName)

	fmt.Printf("config: %s", config)

	createGroup(groupName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("bitbucketserver_group.%s", resourceName), "name", groupName),
				),
				ExpectError: regexp.MustCompile("API Error: 409"),
			},
		},
	})
}

func TestAccBitbucketResourceGroup_AllowImport(t *testing.T) {
	resourceName := "duplicate_group"
	groupName := "duplicate-group"
	config := fmt.Sprintf(`
		resource "bitbucketserver_group" "%s" {
			name = "%s"
			import_if_exists = true
		}
	`, resourceName, groupName)

	fmt.Printf("config: %s", config)

	createGroup(groupName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("bitbucketserver_group.%s", resourceName), "name", groupName),
				),
			},
		},
	})
}

func createGroup(groupName string) {
	client := newBitbucketClient()
	client.Post(fmt.Sprintf("/rest/api/1.0/admin/groups?name=%s", url.QueryEscape(groupName)), nil)
}

func newBitbucketClient() *BitbucketClient {
	//serverSanitized := d.Get("server").(string)
	serverSanitized := "http://localhost:7990"
	//if strings.HasSuffix(serverSanitized, "/") {
	//	serverSanitized = serverSanitized[0 : len(serverSanitized)-1]
	//}

	return &BitbucketClient{
		Server: serverSanitized,
		//Username:   d.Get("username").(string),
		Username: "admin",
		//Password:   d.Get("password").(string),
		Password:   "admin",
		HTTPClient: &http.Client{},
	}
}
