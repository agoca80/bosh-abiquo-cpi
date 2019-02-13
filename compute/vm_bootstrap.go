package compute

import (
	"bytes"
	"encoding/base64"
	"text/template"

	"github.com/abiquo/ojal/abiquo"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

const bootstrap = `#!/bin/sh -x

# Make some traps 
usermod abiquo -s /bin/bash -G "adm,admin,bosh_users"

# Reconfigure SSH
sed \
  -e '/^ *AllowGroups/d'            \
  -e '/^ *DenyUsers/d'              \
  -e '/^ *PermitRootLogin/d'        \
  -e '/^ *PasswordAuthentication/d' \
  -i etc/ssh/sshd_config
systemctl restart ssh

# Instance data
mkdir /var/vcap/instance
echo "{{.UUID}}" >> /var/vcap/instance/id
echo "{{.Name}}" >> /var/vcap/instance/name

decode () { echo "$2" | base64 -d | jq . > $1; }
decode /var/vcap/bosh/Settings "{{.Settings}}"
decode /var/vcap/bosh/MetaData "{{.MetaData}}"
# decode /var/vcap/bosh/UserData "{{.UserData}}"

exit 0
`

const metadata = `{
	"instance-id": "{{.UUID}}",
	"public-keys": {
		"0": {
			"openssh-key": "{{.PublicSSHKey}}"
		}
	}
}`

const userdata = `{
	"server": {
		"name": "instance-id"
	},
	"registry": {
		"endpoint": "http://antxon:pasahitza@helbidea:39517"
	}
}`

type templateOptions struct {
	Name         string
	Settings     string
	MetaData     string
	UserData     string
	UUID         string
	PublicSSHKey string
}

var (
	metadataTemplate  = template.Must(template.New("metadata").Parse(metadata))
	userdataTemplate  = template.Must(template.New("userdata").Parse(userdata))
	bootstrapTemplate = template.Must(template.New("bootstrap").Parse(bootstrap))
)

func parse(template *template.Template, options templateOptions) ([]byte, error) {
	buffer := &bytes.Buffer{}
	err := template.Execute(buffer, options)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (v *vm) bootstrap(agentEnv apiv1.AgentEnv) (err error) {
	agentBytes, err := agentEnv.AsBytes()
	if err != nil {
		return
	}

	metadata, err := parse(metadataTemplate, templateOptions{UUID: v.Label})
	if err != nil {
		return
	}

	userdata, err := parse(userdataTemplate, templateOptions{})
	if err != nil {
		return
	}

	bootstrap, err := parse(bootstrapTemplate, templateOptions{
		MetaData: base64.StdEncoding.EncodeToString(metadata),
		Name:     v.Name,
		Settings: base64.StdEncoding.EncodeToString(agentBytes),
		UserData: base64.StdEncoding.EncodeToString(userdata),
		UUID:     v.Label,
	})
	if err != nil {
		return
	}

	return v.SetVMMetadata(&abiquo.VirtualMachineMetadata{
		Metadata: abiquo.VirtualMachineMetadataFields{
			StartupScript: string(bootstrap),
		},
	})
}
