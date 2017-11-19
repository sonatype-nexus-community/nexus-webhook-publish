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
	"fmt"
	"errors"
	"log"
	"net/http"
	"io"
	"os"
	"os/exec"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sonatype-nexus-community/nexus-webhook-publish/webhook"	
)

const NEXUS_REPO_BASE_URL = "http://localhost:8081/repository/"
const TEMP_DIR = ".tmp"

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
	downloadFile(&component)
}

func downloadFile(c *webhook.Component) {
	fileName, fullUrl, err := getFileNameAndDownloadUrl(c)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(TEMP_DIR); os.IsNotExist(err) {
		os.Mkdir(TEMP_DIR, 0755)
	}

	out, err := os.Create(TEMP_DIR + "/" + fileName)
	defer out.Close()

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(fullUrl)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if n, err := io.Copy(out, resp.Body); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Downloaded", fileName, "with", n, "bytes")
		publishNpmPackage(TEMP_DIR + "/" + fileName)
	}
}

func publishNpmPackage(name string) {
	cmd := exec.Command("npm","publish", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

}

func getFileNameAndDownloadUrl(c *webhook.Component) (fileName string, downloadUrl string, err error) {
	switch format := c.Component.Format; format {
	case "npm":
		fileName := c.Component.Name + "-" + c.Component.Version + ".tgz"
		fullUrl := NEXUS_REPO_BASE_URL + c.RepositoryName + "/" + c.Component.Name + "/-/" + fileName 
		return fileName, fullUrl, nil
	default:
		return "", "", errors.New("Unsupported format")
	}
}
