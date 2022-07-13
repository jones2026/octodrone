/*
Copyright © 2022 John Jones

*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type DronePipeline struct {
	Kind     string `yaml:"kind"`
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Platform struct {
		Os   string `yaml:"os"`
		Arch string `yaml:"arch"`
	} `yaml:"platform"`
	Steps []struct {
		Name     string                 `yaml:"name"`
		Image    string                 `yaml:"image"`
		Commands []string               `yaml:"commands,omitempty"`
		Settings map[string]interface{} `yaml:"settings,omitempty"`
	} `yaml:"steps"`
	Trigger struct {
		Branch []string `yaml:"branch"`
	} `yaml:"trigger"`
}

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert .drone.yml file into github action workflows",
	Long: `This is a simple CLI tool to convert .drone.yml files to github action workflows.

By default each pipeline in your .drone.yml will be converted to a separate
github action workflow file.`,
	Run: func(cmd *cobra.Command, args []string) {
		parseDrone()
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func parseDrone() {
	yfile, err := ioutil.ReadFile(".drone.yml")
	if err != nil {
		log.Fatal(err)
	}

	dronePipelines := []DronePipeline{}
	if err := UnmarshalAllPipelines(yfile, &dronePipelines); err != nil {
		fmt.Println(err)
		return
	}

	for index, pipeline := range dronePipelines {
		fmt.Println("PIPELINE:", pipeline.Name)
		for _, step := range dronePipelines[index].Steps {
			fmt.Println("name:", step.Name)
			fmt.Println("image:", step.Image)
			fmt.Println(" -", step.Commands)
			fmt.Println("Settings:")
			for key, value := range step.Settings {
				if reflect.ValueOf(value).Kind() == reflect.Map {
					fmt.Println("SECRET", key, " :", value)
				} else {
					fmt.Println("  ", key, ":", value)
				}
			}
		}
		fmt.Println("END OF PIPELINE")
	}

}

func UnmarshalAllPipelines(in []byte, out *[]DronePipeline) error {
	r := bytes.NewReader(in)
	decoder := yaml.NewDecoder(r)
	for {
		var pipeline DronePipeline
		if err := decoder.Decode(&pipeline); err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		*out = append(*out, pipeline)
	}
	return nil
}
