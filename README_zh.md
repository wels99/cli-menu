# climenu

[README](README.md) | [中文文档](README_zh.md)

[![GoDoc](https://img.shields.io/badge/go-documentation-blue)](https://pkg.go.dev/github.com/wels99/climenu)

终端中命令行程序的选择菜单。

支持 Windows, Linux, and MacOS 下的多种终端。

方向键上下选择，左右翻页。`<Esc>`和`<Ctrl>+C`退出，`<Enter>`确认。

![pic](./img/pic001.jpg)

## 安装

```bash
go get -u github.com/wels99/climenu
```

## 代码示例

- 示例 1

```go
...
m := climenu.New()
m.Add("name1", "note1", nil, nil)
m.Add("name2", "note2", nil, func(e *climenu.Item) error {
    fmt.Println("selected: ", e)
    return nil
})
m.AddItem(climenu.Item{
    Name: "name3",
    Note: "note3",
    Tags: v,
    Act: func(e *climenu.Item) error {
        fmt.Println("selected: ", e)
        return nil
    },
}
m.run()
...
```

- 示例 2

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
    // m.SetSelectedColor(climenu.Color_Red)
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
    if ret > 0 {
        e, _ := m.GetItem(ret)
        fmt.Println("return:", e)
    }
}
```
