
package test

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
)

func TestTerraformHelloWorldExample(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/simple",
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	//output := terraform.Output(t, terraformOptions, "local_path")
	//assert.Equal(t, "kubectl-v1.17.0", output)
}
