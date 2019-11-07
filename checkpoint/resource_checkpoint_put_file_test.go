package checkpoint

import (
	"fmt"
	chkp "github.com/Checkpoint/api_go_sdk/APIFiles"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"os"
	"strconv"
	"testing"
)

// Resource put file acceptance test:
// 1. Create resource
// 2. Check if resource exists
// 3. Check resource attributes are the same as in configure
// 4. Check resource destroy
func TestAccChkpPutFile_basic(t *testing.T){
	var put_file_res map[string]interface{}
	resourceName := "chkp_put_file.test"
	objName := "/home/admin/terrafile.txt"
	objContent := "It's terrafile 114"
	objOverride := true

	context := os.Getenv("CHKP_CONTEXT")
	if context != "gaia_api" {
		t.Skip("Skipping Gaia test")
	} else if context == "" {
		t.Skip("Env CHKP_CONTEXT must be specified to run this acc test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			{
				Config: testAccPutFileConfig(objName, objContent, objOverride),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChkpPutFileExists(resourceName,&put_file_res),
					testAccCheckChkpPutFileAttributes(&put_file_res,objName, objContent, objOverride),
				),
			},
		},
	})
}

// verifies resource exists by ID and init res with response data
func testAccCheckChkpPutFileExists(resourceTfName string, res *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Println("Enter testAccCheckChkpPutFileExists...")
		rs, ok := s.RootModule().Resources[resourceTfName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceTfName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		// retrieve the client from test provider. client is after providerConfigure()
		client := testAccProvider.Meta().(*chkp.ApiClient)

		payload := make(map[string]interface{})

		payload["file-name"] = rs.Primary.Attributes["file_name"]
		payload["text-content"] = rs.Primary.Attributes["text_content"]
		payload["override"], _ = strconv.ParseBool(rs.Primary.Attributes["override"])

		response, _ := client.ApiCall("put-file",payload,client.GetSessionID(),true,false)
		if !response.Success {
			return fmt.Errorf(response.ErrorMsg)
		}
		// init res with response data for next step (CheckAttributes)
		*res = payload
		log.Println("Exit testAccCheckChkpPutFileExists...")
		return nil
	}
}

// verifies resource attributes are same as in configure
func testAccCheckChkpPutFileAttributes(piRes *map[string]interface{}, fname string, content string, override bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Println("Enter testAccCheckChkpPutFileAttributes")
		putFileMap := *piRes
		if putFileMap == nil {
			return fmt.Errorf("putFileMap is nil")
		}

		res_fname := putFileMap["file-name"].(string)
		if res_fname != fname {
			return fmt.Errorf("fname is %s, expected %s", res_fname, fname)
		}

		res_content := putFileMap["text-content"].(string)
		if res_content != content {
			return fmt.Errorf("content is %s, expected %s", res_content, content)
		}

		res_override := putFileMap["override"].(bool)
		if res_override != override {
			return fmt.Errorf("fname is %t, expected %t", res_override, override)
		}

		log.Println("Exit testAccCheckChkpPutFileAttributes")
		return nil
	}
}

// return a string of the resource like define in a .tf file
func testAccPutFileConfig(file_name string, content string, override bool) string {
	return fmt.Sprintf(`
resource "chkp_put_file" "test" {
      file_name = "%s"
      text_content = "%s"
      override = %t
}
`, file_name, content, override)
}
