/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package minikube

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"k8s.io/minikube/pkg/minikube/constants"
	"k8s.io/minikube/pkg/minikube/notify"
)

const (
	downloadURL = "https://storage.googleapis.com/minikube/releases/%s/minikube-%s-amd64%s"
)

func getDownloadURL(version, platform string) string {
	switch platform {
	case "windows":
		return fmt.Sprintf(downloadURL, version, platform, ".exe")
	default:
		return fmt.Sprintf(downloadURL, version, platform, "")
	}
}

func getShaFromURL(url string) (string, error) {
	fmt.Println("Downloading: ", url)
	r, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	b := sha256.Sum256(body)
	return hex.EncodeToString(b[:]), nil
}

func TestReleasesJson(t *testing.T) {
	releases, err := notify.GetAllVersionsFromURL(constants.GithubMinikubeReleasesURL)
	if err != nil {
		t.Fatalf("Error getting releases.json: %s", err)
	}

	for _, r := range releases {
		fmt.Printf("Checking release: %s\n", r.Name)
		for platform, sha := range r.Checksums {
			fmt.Printf("Checking SHA for %s.\n", platform)
			actualSha, err := getShaFromURL(getDownloadURL(r.Name, platform))
			if err != nil {
				t.Errorf("Error calcuating SHA for %s-%s. Error: %s", r.Name, platform, err)
			}
			if actualSha != sha {
				t.Errorf("ERROR: SHA does not match for version %s, platform %s. Expected %s, got %s.", r.Name, platform, sha, actualSha)
			}
		}
	}
}
