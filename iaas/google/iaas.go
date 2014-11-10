package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/megamd/global"
	"github.com/megamsys/megamd/iaas"
	"github.com/megamsys/megamd/provisioner"
	"github.com/tsuru/config"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func Init() {
	iaas.RegisterIaasProvider("google", &GoogleIaaS{})
}

type GoogleIaaS struct{}

type GoogleCredential struct {
	AuthorizationURI   string `json:"authorization_uri"`
	TokenCredentailURI string `json:"token_credential_uri"`
	Scope              string `json:"scope"`
	RedirectURI        string `json:"redirect_uri"`
	ClientID           string `json:"client_id"`
	ClientSecret       string `json:"client_secret"`
	AccessToken        string `json:"access_token"`
	ExpiresIN          string `json:"expires_in"`
	RefreshToken       string `json:"refresh_token"`
	Project            string `json:"project"`
}

func (i *GoogleIaaS) DeleteMachine(*global.PredefClouds, *provisioner.AssemblyResult) (string, error) {

	return "", nil
}

func (i *GoogleIaaS) CreateMachine(pdc *global.PredefClouds, assembly *provisioner.AssemblyResult) (string, error) {
	cre, derr := downloadCredentials(pdc)
	if derr != nil {
		return "", derr
	}
	str, err := buildCommand(iaas.GetPlugins("google"), pdc, "create")
	if err != nil {
		return "", err
	}
	str = str + " -N " + assembly.Name + "." + assembly.Components[0].Inputs.Domain
	riak, err_riak := config.GetString("api:server")
	if err_riak != nil {
		return "", err_riak
	}
	
	recipe, err_recipe := config.GetString("knife:recipe")
	if err_recipe != nil {
		return "", err_recipe
	}
	
	str = str + " --run-list \"" + "recipe[" + recipe + "]" + "\""
	attributes := &iaas.Attributes{RiakHost: riak, AccountID: pdc.Accounts_id, AssemblyID: assembly.Id}
    b, aerr := json.Marshal(attributes)
    if aerr != nil {
        fmt.Println(aerr)
        return "", aerr
    }
	str = str + " --json-attributes " + string(b)
	str = strings.Replace(str, "create -f", "create "+assembly.Name+"."+assembly.Components[0].Inputs.Domain+" -f "+cre, -1)
	knifePath, kerr := config.GetString("knife:path")
	if kerr != nil {
		return "", kerr
	}
	str = strings.Replace(str, "-c", "-c "+knifePath, -1)
	return str, nil
}

func buildCommand(plugin *iaas.Plugins, pdc *global.PredefClouds, command string) (string, error) {
	var buffer bytes.Buffer
	if len(plugin.Tool) > 0 {
		buffer.WriteString(plugin.Tool)
	} else {
		return "", fmt.Errorf("Plugin tool doesn't loaded")
	}

	if command == "create" {
		if len(plugin.Command.Create) > 0 {
			buffer.WriteString(" " + plugin.Command.Create)
		} else {
			return "", fmt.Errorf("Plugin commands doesn't loaded")
		}
	}

	if len(pdc.Spec.Groups) > 0 {
		buffer.WriteString(" -n " + pdc.Spec.Groups)
	} else {
		return "", fmt.Errorf("Groups doesn't loaded")
	}

	if len(pdc.Spec.Image) > 0 {
		buffer.WriteString(" -I " + pdc.Spec.Image)
	} else {
		return "", fmt.Errorf("Image doesn't loaded")
	}

	if len(pdc.Spec.Flavor) > 0 {
		buffer.WriteString(" -m " + pdc.Spec.Flavor)
	} else {
		return "", fmt.Errorf("Flavor doesn't loaded")
	}

	if len(pdc.Access.Sshuser) > 0 {
		buffer.WriteString(" -x " + pdc.Access.Sshuser)
	} else {
		return "", fmt.Errorf("Ssh user value doesn't loaded")
	}

	if len(pdc.Access.IdentityFile) > 0 {
		buffer.WriteString(" --identity-file " + pdc.Access.IdentityFile)
	} else {
		return "", fmt.Errorf("Identity file doesn't loaded")
	}

	if len(pdc.Access.Zone) > 0 {
		buffer.WriteString(" -Z " + pdc.Access.Zone)
	} else {
		return "", fmt.Errorf("Zone doesn't loaded")
	}

	return buffer.String(), nil
}

func downloadCredentials(pdc *global.PredefClouds) (string, error) {
	s := make([]string, 2)
	s = strings.Split(pdc.Access.VaultLocation, "_")
	email, name := s[0], s[1]
	cloudkeysBucket, ckberr := config.GetString("buckets:CLOUDKEYS")
	if ckberr != nil {
		return "", ckberr
	}

	megam_home, ckberr := config.GetString("MEGAM_HOME")
	if ckberr != nil {
		return "", ckberr
	}

	basePath := megam_home + cloudkeysBucket
	cloudaccesskeysBucket, cakberr := config.GetString("buckets:CLOUDACCESSKEYS")
	if cakberr != nil {
		return "", cakberr
	}
	conn, err := db.Conn(cloudaccesskeysBucket)
	keys := &GoogleCredential{}
	if err != nil {
		return "", err
	}

	ferr := conn.FetchStruct(pdc.Access.VaultLocation, keys)
	if ferr != nil {
		return "", ferr
	}

	b, errk := json.Marshal(keys)
	if errk != nil {
		return "", errk
	}

	dir := path.Join(basePath, email)
	filePath := path.Join(dir, name+".json")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", dir)

		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return "", errm
		}
		// open output file
		_, err := os.Create(filePath)
		if err != nil {
			return "", err
		}
	}
	errf := ioutil.WriteFile(filePath, []byte(string(b)), 0644)
	if errf != nil {
		return "", errf
	}

	return cloudkeysBucket + "/" + email + "/" + name + ".json", nil
}
