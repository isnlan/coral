package k8s_yaml

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/util/yaml"
)

//yaml to 字节数组
func Yamls2Bytes(rootPath string, files []string) ([][]byte, error) {
	yamls := make([][]byte, len(files))
	for i, name := range files {
		yamlBytes, err := ioutil.ReadFile(filepath.Join(rootPath, name))
		if err != nil {
			return nil, errors.Wrap(err, "error read file: "+name)
		}
		yamls[i] = yamlBytes

	}
	return yamls, nil
}

//yaml to json
func Yamls2Jsons(data [][]byte) [][]byte {
	jsons := make([][]byte, 0)
	for _, yamlBytes := range data {
		yamls := bytes.Split(yamlBytes, []byte("---"))
		for _, v := range yamls {
			if len(v) == 0 {
				continue
			}
			obj, err := yaml.ToJSON(v)
			if err != nil {
				log.Println(err.Error())
			}
			jsons = append(jsons, obj)
		}

	}
	return jsons
}
