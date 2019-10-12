package model

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"gopkg.in/yaml.v2"
)

//load yaml
//build step order
//

//CF ...
type Metadata struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Title string `yaml:"title"`
	Desc  string `yaml:"description"`
}
type Scale struct {
	Scale  []string
	Matrix []string
}
type FlowControl struct {
	when struct {
		steps []struct {
			name string
			on   []string
		}

		//args map[interface{}]interface{}
	}
}

type StepMetadata struct {
	Metadata `yaml:",inline"`
	Scale
	FlowControl
}

func NewStepMetadata(n string, t string, title string, desc string) *StepMetadata {
	return &StepMetadata{
		Metadata{n, t, title, desc},
		Scale{[]string{}, []string{}},
		FlowControl{}}
}

type Step interface {
	GetType() StepMetadata
	Print(w io.Writer)
	Run() error
	//Before(ctx context.Context)
	//After(ctx context.Context)
}
type FreeStyle struct {
	StepMetadata
	Name     string `yaml:"name"`
	Type     string
	Commands []string `yaml:"commands"`
	WorkDir  string
	Image    string `yaml:"image"`
}
type Build struct {
	StepMetadata
	Dockerfile string `yaml:"dockerfile"`
	Image      string
	WorkDir    string
	Tag        string
}
type GitClone struct {
	StepMetadata
	Repo     string
	Revision string
	Git      string
}
type CF struct {
	Version string `yaml:"version"`
	Steps   []Step `yaml:"steps"`
}
type Extra map[string]interface{}
type StepAll struct {
	StepMetadata `yaml:",inline"`

	Extra map[string]interface{} `yaml:",inline"`
}

func (s *StepAll) GetType() StepMetadata {
	return StepMetadata{}
	//*NewStepMetadata(s.Name, s.Type, s.Title, s.Desc)
}

//Print ...
func (s *StepAll) Print(w io.Writer) {
	str := fmt.Sprintf("\n[name: %v\n type: %v\n title: %v\n description : %v\n", s.Name, s.Type, s.Title, s.Desc)
	w.Write([]byte(str))

}
func defaultRun(s Step) error {
	fmt.Printf("\nexecuting step %v\n", s.GetType().Name)
	return nil
}
func (s *StepAll) Run() error {
	return defaultRun(s)

}

//TestStep ...
func TestStep(data map[string]interface{}, steps []StepMetadata) {

}
func (cf *CF) UnmarshalYAML(unmarshal func(interface{}) error) error {
	fmt.Println("=====UnmarshalYAML=====")
	var c map[string]interface{}
	err := unmarshal(&c)
	if err != nil {
		panic(err)
	}
	m := make(map[string]StepAll)
	v := c["steps"]
	stepBytes := fmt.Sprintf("%v", v)

	fmt.Printf("%v\n", stepBytes)
	d, err := yaml.Marshal(&v)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(d, &m)
	if err != nil {
		panic(err)
	}

	for i, v := range m {
		v.Name = i

		if v.Type == "" {
			v.Type = "freestyle"
		}

		switch v.Type {
		case "freestyle":
			cmds, ok := v.Extra["commands"].([]interface{})
			t := reflect.TypeOf(v.Extra["commands"])
			if !ok {
				fmt.Printf("can't convert , real type is %v\n ", t)
			}
			fmt.Printf("[%v]\n - type  %v\n %v\n %v\n %v\n",
				v.Name, v.Type,
				v.Extra["image"],
				v.Extra["workdir"],
				cmds[1].(string))
		case "git-clone":
			fmt.Printf("[%v]\n - type %v\n %v\n %v\n %v\n",
				v.Name, v.Type,
				v.Extra["git"],
				v.Extra["repo"],
				v.Extra["revision"])
		case "build":
			fmt.Printf("[%v]\n - type %v\n %v\n %v\n %v\n %v\n",
				v.Name, v.Type, v.Extra["dockerfile"],
				v.Extra["working_directory"],
				v.Extra["image"],
				v.Extra["tag"])
		}

	}

	return err
}

//UnmarshalYAML ...
func (cf *CF) UnmarshalYAML1(unmarshal func(interface{}) error) error {
	fmt.Println("=====UnmarshalYAML=====")
	var c map[string]interface{}
	err := unmarshal(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
	fmt.Println("==========")
	cf.Version = c["version"].(string)
	steps, ok := c["steps"].(map[interface{}]interface{})

	if !ok {
		return errors.New("can't unmarshall steps")
	}
	fmt.Printf("[%v]\n", steps)
	var result []Step
	for k, v := range steps {
		fmt.Printf("[%v]\n", k.(string))
		fmt.Printf("[%v]\n", v)
		s, ok := v.(map[interface{}]interface{})
		if !ok {
			fmt.Printf("something wrong in conversioin")
			continue
		}

		t, ok := s["type"].(string)

		if !ok || t == "" {
			t = "freestyle"
		}
		_, ok = k.(string)
		if !ok {
			panic(k)
		}
		//step := HandleStep(stepName, t, s)
		var step Step
		result = append(result, step)
		if step != nil {
			fmt.Printf("step - [%v]\n", step.GetType().Name)
		}
	}

	cf.Steps = result
	return err
}
