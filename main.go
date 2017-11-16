/* Copyright 2017 Sonatype

Licensed under the Apache License, Version 2.0 (the "License"); 
you may not use this file except in compliance with the License. 
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software 
distributed under the License is distributed on an "AS IS" BASIS, 
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. 
See the License for the specific language governing permissions and 
limitations under the License. */

package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sonatype-nexus-community/nexus-webhook-publish/webhook"	
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/publish", publishPackage)

	return r
}

func main() {
	r := setupRouter()
	r.Use(gin.Logger())
	r.Run(":8000")
}

func publishPackage(c *gin.Context) {
	body, err := webhook.Handler(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
	}
	var component webhook.Component
	err = json.Unmarshal(body, &component)
	if err != nil {
		c.AbortWithError(500, err)
	}
	// do something with component
}
