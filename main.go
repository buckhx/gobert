package main

import (
	"fmt"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {
	path := "./python/output"
	tags := []string{"bert-uncased"}
	bm, err := tf.LoadSavedModel(path, tags, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", bm)
	fmt.Println("Session")
	fmt.Println("\tDevice")
	devs, err := bm.Session.ListDevices()
	if err != nil {
		fmt.Println(err)
	}
	for _, dev := range devs {
		fmt.Printf("\t\t%+v\n", dev)
	}
	fmt.Println("Graph")
	fmt.Println("\tOperation")
	for _, op := range bm.Graph.Operations() {
		fmt.Printf("\t\t%+v\n", op.Name())
	}
}
