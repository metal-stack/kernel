package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v37/github"
	"github.com/mmcdole/gofeed"
	"golang.org/x/oauth2"
)

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
	}
	fmt.Printf("new kernel:%s\n", kv)

	err := createTag("5.10.48-65")
	if err != nil {
		panic(err)
	}
}

func getNextKernelVersion(major, minor int, supportLevel string) (string, bool) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.kernel.org/feeds/kdist.xml")
	for _, item := range feed.Items {
		parts := strings.Split(item.Title, ": ")
		if len(parts) != 2 {
			continue
		}
		kernelVersion := parts[0]
		support := parts[1]
		version, err := semver.NewVersion(kernelVersion)
		if err != nil {
			continue
		}
		isOdd := version.Minor()%2 == 0
		if version.Major() == uint64(major) && version.Minor() == uint64(minor) && support == supportLevel && isOdd {
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
