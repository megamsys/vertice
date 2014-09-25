package cmd

import (
  "fmt"
  "log"
  "gopkg.in/yaml.v1"
  "io/ioutil"
"path/filepath"
)


const defaultTOSCAPath = "conf/tosca_schema.yaml"

func NewTOSCA() string {
        p, _ := filepath.Abs(defaultTOSCAPath)
        log.Println(fmt.Errorf("Conf: %s", p))

        data, err  := ioutil.ReadFile(p)

        if err != nil {
            log.Fatalf("error: %v", err)
        }

        m := make(map[interface{}]interface{})

        err = yaml.Unmarshal([]byte(data), &m)
        if err != nil {
                log.Fatalf("error: %v", err)
        }
        for key, value := range m {
        log.Printf("\n%v\n :=>\n", key)
          switch value.(type) {
            case string:
              log.Printf("===[%v]\n",value)
              case map[interface{}]interface{}:
              log.Printf("???[%v]\n",value)
              break
              default:
              log.Printf(">>>[%v] is unknown!!\n", value)
          }
       }

        d, err := yaml.Marshal(&m)
        if err != nil {
                log.Fatalf("error: %v", err)
        }
      //  log.Printf("--- m dump:\n%s\n\n", string(d))
        return string(d)
}
