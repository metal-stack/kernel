package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/mmcdole/gofeed"
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
		if version.Major() == uint64(major) && version.Minor() == uint64(minor) && support == supportlevel && isOdd {
			// This is the next desired kernel version we would like to release
			// TODO: create a tag with this version and modify the action to build the kernel based on the tag name
			fmt.Println(version)
			return
		}
	}
}
