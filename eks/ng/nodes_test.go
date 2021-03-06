package ng

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"
)

func TestTemplateASG(t *testing.T) {
	tpl := template.Must(template.New("TemplateASG").Parse(TemplateASG))
	buf := bytes.NewBuffer(nil)
	if err := tpl.Execute(buf, templateASG{}); err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf.String())

	buf.Reset()
	if err := tpl.Execute(buf, templateASG{
		ImageID:             "abc",
		ImageIDSSMParameter: "",
		Metadata:            metadataAL2InstallSSM,
		UserData:            userDataAL2InstallSSM,
		ASGDesiredCapacity:  1,
	}); err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf.String())
	if strings.Contains(buf.String(), "AWS::SSM::Parameter") {
		t.Fatal("unexpected AWS::SSM::Parameter in CFN template")
	}
}
