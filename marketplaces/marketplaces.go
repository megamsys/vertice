package marketplaces

import (
  "github.com/megamsys/vertice/marketplaces/provision"
)
//Global provisioners set by the subd daemons.
var ProvisionerMap map[string]provision.Provisioner = make(map[string]provision.Provisioner)

/*
type RawImageAccess interface{
}
type MarketPlaceAccess interface{
}

type MarketplaceInterface interface {
   access(cat string) (interface, error)
}

type Global struct {
  AccountId string
  Category  string
  Action    string
  Access    interface{}
}

func (g *Global) access(cat string) (interface, error) {
  switch expression {
  case condition:

  }
}

category: rawimage             //Get rawimages build market
action: rawimage.iso.create

category: marketplaces        // Get Marketplaces build market
action: rawimage.iso.customize

category: marketplaces       // get marketplaces build market
action:   marketplaces.add.customized

func (p *ReqOperator) getAccess() (MarketplaceInterface, error) {
  switch p.Category {
  case condition:

  }
  return getAccess(p.Category,p.CartonsId,p.AccountId)
}

func getAccess(cat ,id ,email string)(*MarketplaceInterface, error)
  if err != nil {
    return nil, err
  }
  c, err := b.MkCloudMark()
  if err != nil {
    return nil, err
  }
  return c, nil

case marketplaces:
  d, err := GetMarkplace(p.CartonsId,p.AccountId)
  if err != nil {
    return nil, err
  }
  c, err := d.MkCloudMark()
  if err != nil {
    return nil, err
  }
  return c, nil

}
}

type RawImages struct {
  Id        string
  AccountId string
  OrgId     string
  Api       api.ApiArgs
  Inputs    pairs.JsonPairs
  Outputs   pairs.JsonPairs
  Repos     Repos
}

type Repos struct {
 Source string
 PubliceUrl string
 Properties pairs.JsonPairs
}

// struct for marketplaces and rawimages
type Marketplaces struct {
  Id        string
  AccountId string
  Api       api.ApiArgs
  Inputs    pairs.JsonPairs
  Outputs   pairs.JsonPairs
  Envs      pairs.JsonPairs
  Options   pairs.JsonPairs
  Name      string
  Flavor    string
  Image     string
  CatOrder  string
  Plans     map[string]string
  Status    string
  JsonClaz  string
}


{
Id: “RAW10001”,
name: “myfirst_iso”
org_id : “ORG123”,
account_id: “info@megam.io”,
Inputs: [region]
Json_claz: “”
repos: {
   source: “ISO”,
   public_url:  “”,
   properties: [key,value list],
}
outputs [image_id: 20 , ]
Status: “inprogress, active”
}


   CREATE TABLE marketplaces ( settings_name text,  cattype text, flavor text,  image text, catorder text, url text, json_claz text, envs list<text>,  options list<text>, plans map <text, text>,  PRIMARY KEY ((settings_name), flavor));

   ALTER TABLE marketplaces ADD  id text, account_id  text, status text;
   ALTER TABLE marketplaces ADD  inputs list<text>, outputs list<text>, acl_policies list<text>;

   data
   {
   id:”MKP0001 “,
   account_id: “ “,
   status: “inprogress”,
   Inputs: [raw_image_id: “RAW10001 “,
             image_id: 20, region: sydney],
   }
*/
