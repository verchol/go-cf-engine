package main

import (
	"fmt"
	"os"

	model "github.com/verchol/go-cf-engine/pkg/model"
)

func main() {
	t := model.ParseYaml("")
	fmt.Println(t.Version)

	for _, v := range t.Steps {
		v.Print(os.Stdout)
		v.Run()
		//fmt.Printf("stepName [%v] -  type [%s]\n", v.GetType().Name, v.GetType().Type)
	}

}
