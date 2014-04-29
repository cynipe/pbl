package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "pbl"
	app.Usage = "List your Pinboard's bookmarks in terminal"
	app.Action = func(c *cli.Context) {
		posts, err := GetPosts(c.Args())
		if err != nil {
			println(err)
		}
		for _, post := range posts.Posts {
			fmt.Printf("%s: %s\n", post.Desc, post.Url)
		}
	}
	app.Run(os.Args)
}
