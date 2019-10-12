package model

import (
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

/*
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
		}

	}

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
		}

	}

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
}*/

//ParseYaml ...
func ParseYaml(path string) CF {
	t := CF{}
	if path == "" {
		path = "./codefresh.yaml"
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		panic(err)
	}
	return t
}
