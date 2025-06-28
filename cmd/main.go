package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

type KernelVersions struct {
	LatestStable struct {
		Version string `json:"version"`
	} `json:"latest_stable"`
	Releases []struct {
		Iseol    bool        `json:"iseol"`
		Version  string      `json:"version"`
		Moniker  string      `json:"moniker"`
		Source   string      `json:"source"`
		Pgp      interface{} `json:"pgp"`
		Released struct {
			Timestamp int    `json:"timestamp"`
			Isodate   string `json:"isodate"`
		} `json:"released"`
		Gitweb    string      `json:"gitweb"`
		Changelog interface{} `json:"changelog"`
		Diffview  string      `json:"diffview"`
		Patch     struct {
			Full        string      `json:"full"`
			Incremental interface{} `json:"incremental"`
		} `json:"patch"`
	} `json:"releases"`
}

func main() {
	major, _ := strconv.Atoi(os.Getenv("KERNEL_MAJOR"))
	if major == 0 {
		major = 5
	}
	minor, _ := strconv.Atoi(os.Getenv("KERNEL_MINOR"))
	if minor == 0 {
		minor = 10
	}
	supportlevel := os.Getenv("KERNEL_SUPPORT")
	if supportlevel == "" {
		supportlevel = "longterm"
	}

	kv, ok := getNextKernelVersion(major, minor, supportlevel)
	if !ok {
		fmt.Println("no new kernel found")
		return
	}
	fmt.Printf("new kernel:%s\n", kv)

	err := createTag("5.10.48-65")
	if err != nil {
		panic(err)
	}
}

func getKernelVersions() (*KernelVersions, error) {
	url := "https://www.kernel.org/releases.json"
	c := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	kv := KernelVersions{}
	err = json.Unmarshal(body, &kv)
	if err != nil {
		return nil, err
	}
	return &kv, nil
}

func getNextKernelVersion(major, minor int, supportLevel string) (string, bool) {
	kv, err := getKernelVersions()
	if err != nil {
		return "", false
	}
	for _, r := range kv.Releases {
		if r.Moniker != supportLevel {
			continue
		}
		kernelVersion := r.Version
		version, err := semver.NewVersion(kernelVersion)
		if err != nil {
			continue
		}
		isOdd := version.Patch()%2 == 0
		if !isOdd {
			continue
		}
		if version.Major() == uint64(major) && version.Minor() == uint64(minor) {
			// This is the next desired kernel version we would like to release
			// TODO: create a tag with this version and modify the action to build the kernel based on the tag name
			return version.String(), true
		}
	}
	return "", false
}

func createTag(name string) error {
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(nil)
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		client = github.NewClient(tc)
	}

	tag, resp, err := client.Git.GetTag(ctx, "metal-stack", "kernel", name)
	if err != nil && resp.StatusCode != http.StatusNotFound {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("not found")
	} else {
	}
	fmt.Printf("tag exists: %v\n", tag)

	// tag := &github.Tag{
	// 	Tag: &name,
	// }
	// _, _, err := client.Git.CreateTag(ctx, "metal-stack", "kernel", tag)
	// if err != nil {
	// 	return err
	// }
	return nil
}
