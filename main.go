package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	app := cli.NewApp()
	app.Name = "pbl"
	app.Usage = "List your Pinboard's bookmarks in terminal"
	app.Action = func(c *cli.Context) {
		posts, err := GetPosts(c.Args())
		if err != nil {
			color.Printf("@r%s\n", err)
			os.Exit(1)
		}
		for _, post := range posts.Posts {
			fmt.Printf("%s | %s | %s\n",
				color.Sprintf("@b%s", post.Time.Format("2006-01-02")),
				truncate(post.Desc, 70),
				color.Sprintf("@g%s", post.Url))
		}
	}
	app.Run(os.Args)
}

func truncate(str string, truncatePoint int) (truncated string) {
	length := 0
	var buf bytes.Buffer
	for _, r := range str {
		if utf8.RuneLen(r) == 1 {
			length += 1
		} else {
			length += 2
		}
		if length > truncatePoint {
			break
		}
		buf.WriteRune(r)
	}

	replacement := "..."
	if length > truncatePoint {
		buf.WriteString(replacement)
		return buf.String()
	}

	padSize := (truncatePoint - length + 1 + len(replacement))
	return strings.Repeat(" ", padSize) + buf.String()
}
