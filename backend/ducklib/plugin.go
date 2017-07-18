// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package ducklib

type Plugin interface {
	Intialize()
	Shutdown()
}
