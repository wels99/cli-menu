//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/wels99/climenu"
)

func main() {

	mitems := [][]string{}

	for i := 0; i < 1000; i++ {
		mitems = append(mitems, []string{
			fmt.Sprintf("item%03d", i),
			fmt.Sprintf("note%03d", i),
		})
	}

	m := climenu.New()

	m.SetIndex(true)
	m.SetSelectIcon(" \u27A4 ") // '➤'
	m.SetMessage("select one:")
	m.SetPagesize(5)
	// m.SetSelectedStyle(climenu.Style_Red, climenu.Style_White_bg)
	// m.SetSelectedStyle(climenu.Style_Reverse)
	// m.SetDelimiter("|")
	// m.Seti18n("当前", "页", "方向键移动，回车确认")

	for _, v := range mitems {
		name := v[0]
		note := fmt.Sprintf("%s_%s", v[1], strings.Repeat("🍕", rand.Intn(10)))
		note2 := fmt.Sprintf("%s_%s", v[1], strings.Repeat("🍀", rand.Intn(20)))

		m.Add(name, note, v, func(e *climenu.Item) error {
			fmt.Println("selected: ", e)
			return nil
		})
		m.AddItem(climenu.Item{
			Name: name,
			Note: note2,
			Tags: v,
			Act: func(e *climenu.Item) error {
				fmt.Println("new item selected: ", e)
				return nil
			},
		})
	}

	m.Sort(func(i, j *climenu.Item) bool {
		return i.Name > j.Name
	})

	ret, _ := m.Run()
	if ret != nil {
		fmt.Println("return:", ret)
		ret.Act(ret)
	}
}
