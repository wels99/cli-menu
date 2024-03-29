# climenu

[README](README.md) | [中文文档](README_zh.md)

[![GoDoc](https://img.shields.io/badge/go-documentation-blue)](https://pkg.go.dev/github.com/wels99/climenu)

The selection menu for the command line program in terminal.

Supports various terminals under Windows, Linux, and MacOS.

Use the up and down arrow keys to select, and use the left and right arrow keys to turn pages. `<Esc>` and `<Ctrl>+C` exit, `<Enter>` confirm.

Key `<PageUp>` and `<PageDown>` can also be turn pages.

Input regular characters to filter menu item names. Case insensitive.

![pic](./img/pic001.jpg)

## Install

```bash
go get -u github.com/wels99/climenu
```

## Example

- Example 1

```go
...
m := climenu.New()
m.Add("name1", "note1", nil, func(e *climenu.Item) error {
    fmt.Println("selected: ", e)
    return nil
})
m.Add("name2", "note2", nil, func(e *climenu.Item) error {
    fmt.Println("selected: ", e)
    return nil
})
m.run()
...
```

- Example 2

```go
package main

import (
    "fmt"
    "math/rand"
    "strings"

    "github.com/wels99/climenu"
)

func main() {

    mitems := [][]string{
        {"item001", "note001"},
        {"item002", "note002"},
        {"item003", "note003"},
        {"item004", "note004"},
        {"item005", "note005"},
        {"item006", "note006"},
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

    e, _ := m.Run()
    if e != nil {
        fmt.Println("return:", e)
        ret.Act(e)
    }
}
```
