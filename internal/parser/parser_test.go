package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSOP(t *testing.T) {
	// Create a temporary markdown file for testing
	testContent := `# Deploy Nginx

This SOP deploys nginx to the server.

## Step 1: Check if nginx is running
Check the status of nginx service.

` + "```bash" + `
curl -I localhost
` + "```" + `

## Step 2: Restart nginx if needed
Restart the nginx service if it's not running.

` + "```bash" + `
sudo systemctl restart nginx
` + "```" + `

Done with the SOP.
`

	// Create temp file
	tmpfile, err := os.CreateTemp("", "test-sop-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	if _, err := tmpfile.Write([]byte(testContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test the parser
	sop, err := ParseSOP(tmpfile.Name())
	if assert.NoError(t, err) {
		assert.Equal(t, "Deploy Nginx", sop.Title)
		assert.Equal(t, tmpfile.Name(), sop.Path)
		assert.Len(t, sop.Steps, 2) // Should find 2 command steps

		if len(sop.Steps) >= 2 {
			assert.Equal(t, "curl", sop.Steps[0].Title)
			assert.Equal(t, "bash", sop.Steps[0].CommandType)
			assert.Equal(t, "curl -I localhost", sop.Steps[0].Command)
			
			assert.Equal(t, "sudo", sop.Steps[1].Title)
			assert.Equal(t, "bash", sop.Steps[1].CommandType)
			assert.Equal(t, "sudo systemctl restart nginx", sop.Steps[1].Command)
		}
	}
}