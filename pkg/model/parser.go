package model

import (
	"errors"
	"fmt"
	"io"
	ioutil "io/ioutil"

	"gopkg.in/yaml.v2"
)

const cfYaml = `
	version: '1.0'
	steps:
	main_clone:
		title: Cloning main repository...
		type: git-clone
		repo: 'codefresh-contrib/golang-sample-app'
		revision: master
		git: github
	MyAppDockerImage:
		title: Building Docker Image
		type: build
		image_name: my-golang-image
		working_directory: ./
		tag: full
		dockerfile: Dockerfile`

//CF ...
type StepMetadata struct {
	Name  string `yaml:"name"`
	Type  string
	Title string
	Desc  string
}
type Step interface {
	GetType() StepMetadata
	Print(w io.Writer)
	Run() error
}
type FreeStyle struct {
	StepMetadata
	Name     string `yaml:"name"`
	Type     string
	Commands []string
	WorkDir  string
	Image    string
}
type Build struct {
	StepMetadata
	Dockerfile string
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

func (s *FreeStyle) GetType() StepMetadata {
	return StepMetadata{s.Name, s.Type, s.Title, s.Desc}
}

//Print ...
func (s *FreeStyle) Print(w io.Writer) {
	str := fmt.Sprintf("\n[name: %v\n type: %v\n title: %v\n description : %v\n", s.Name, s.Type, s.Title, s.Desc)
	w.Write([]byte(str))

	w.Write([]byte(fmt.Sprintf("Image: %v\n WorkdDir %v\n Commands =TBD %v\n]",
		s.Image,
		s.WorkDir,
		s.Commands)))
}
func defaultRun(s Step) error {
	fmt.Printf("\nexecuting step %v\n", s.GetType().Name)
	return nil
}
func (s *Build) Run() error {
	return defaultRun(s)

}
func (s *GitClone) Run() error {
	return defaultRun(s)
}
func (s *FreeStyle) Run() error {
	return defaultRun(s)
}

//GetType ...
func (s *Build) GetType() StepMetadata {
	return StepMetadata{s.Name, s.Type, s.Title, s.Desc}
}

//GetType GitClone ...
func (s *GitClone) GetType() StepMetadata {
	return StepMetadata{s.Name, s.Type, s.Title, s.Desc}
}

//Print ...
func (s *Build) Print(w io.Writer) {
	str := fmt.Sprintf("\n[name: %v\n type: %v\n title: %v\n description : %v\n", s.Name, s.Type, s.Title, s.Desc)
	w.Write([]byte(str))

	w.Write([]byte(fmt.Sprintf("Dockerfile: %v\n Image %v\n Tag %v\n]",
		s.Dockerfile,
		s.Image,
		s.Tag)))

}

//Print ...
func (s *GitClone) Print(w io.Writer) {
	str := fmt.Sprintf("\n[name: %v\n type: %v\n title: %v\n description : %v\n", s.Name, s.Type, s.Title, s.Desc)
	w.Write([]byte(str))

	w.Write([]byte(fmt.Sprintf(" Repo: %v\n Revision %v\n Git %v\n]",
		s.Repo,
		s.Revision,
		s.Git)))

}

func HandleBuildStep(s Build, stepData map[interface{}]interface{}) Step {

	for k, v := range stepData {
		fmt.Printf("k=%v, v=%v\n", k, v)

		fd, ok := v.(string)
		if !ok {
			continue
		}
		fmt.Println("simple case...")

		switch k {
		case "title":
			s.Title = fd
		case "image_name":
			s.Image = fd
		case "working_directory":
			s.WorkDir = fd
		case "tag":
			s.Tag = fd
		case "dockerfile":
			s.Dockerfile = fd
		}

	}
	fmt.Println("!!%v", s)
	return &s
}
func HandleGitClone(s GitClone, stepData map[interface{}]interface{}) Step {

	for k, v := range stepData {
		fmt.Printf("k=%v, v=%v\n", k, v)

		fd, ok := v.(string)
		if !ok {
			continue
		}
		fmt.Println("simple case...")
		switch k {
		case "repo":
			s.Repo = fd
		case "git":
			s.Git = fd
		case "revision":
			s.Revision = fd
		case "title":
			s.Title = fd
		}
		/*else {
			for k1 , v1 := range fieldData {
			//image := v1 .(string)
			fmt.Printf("k1=%v, v1=%v\n", k1 ,v1);
			step.Commands = stepData["commands"].([]string)
		   }
		}*/

	}
	fmt.Println("!!%v", s)
	return &s
}

//HandelStep ....
func HandleFreestyle(s FreeStyle, stepData map[interface{}]interface{}) Step {

	for k, v := range stepData {
		fmt.Printf("k=%v, v=%v\n", k, v)

		fd, ok := v.(string)
		if !ok {
			continue
		}
		fmt.Println("simple case...")
		switch k {
		case "image":
			s.Image = fd
		case "working_directory":
			s.WorkDir = fd
		case "description":
			s.Desc = fd
		case "title":
			s.Title = fd
		}
		/*else {
			for k1 , v1 := range fieldData {
			//image := v1 .(string)
			fmt.Printf("k1=%v, v1=%v\n", k1 ,v1);
			step.Commands = stepData["commands"].([]string)
		   }
		}*/

	}
	fmt.Println("!!%v", s)
	return &s
}

func HandleStep(n string, t string, s map[interface{}]interface{}) Step {

	var result Step
	switch t {

	case "build":
		step := Build{}
		step.Name = n
		step.Type = t
		result = HandleBuildStep(step, s)

	case "git-clone":
		step := GitClone{}
		step.Name = n
		step.Type = t
		result = HandleGitClone(step, s)

	default:
		step := FreeStyle{}
		step.Name = n
		step.Type = t
		result = HandleFreestyle(step, s)

	}
	fmt.Printf("step:%v\n", result)
	return result
}

func (cf *CF) UnmarshalYAML(unmarshal func(interface{}) error) error {
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
		stepName, ok := k.(string)
		if !ok {
			panic(k)
		}
		step := HandleStep(stepName, t, s)
		result = append(result, step)
		if step != nil {
			fmt.Printf("step - [%v]\n", step.GetType().Name)
		}
	}

	cf.Steps = result
	return err
}

//ParseYaml ...
func ParseYaml(path string) CF {
	t := CF{}
	yamlFile, err := ioutil.ReadFile("./codefresh.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		panic(err)
	}
	return t
}
