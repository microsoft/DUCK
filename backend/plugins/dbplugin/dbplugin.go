// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package dbplugin

type MyDatabase struct {
}

func (mdb *MyDatabase) Init(url string) error {
	//fmt.Println("Testextension initialized")
	return nil
}

func init() {
	//db := &MyDatabase{}
	//pluginregistry.RegisterDatabase(db)
}
