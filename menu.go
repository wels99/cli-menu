package climenu

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"

	"golang.org/x/text/width"
)

type Item struct {
	idx    int               // the sequence number
	idxstr string            // Display string for menu sequence number
	Name   string            // Menu item name, first item displayed
	Note   string            // Menu item description, second item displayed
	Tags   any               // Attached data
	Act    func(*Item) error // The execution function of menu items. The function parameter is a menu item
}

type LangTxt struct {
	Cur  string // current index
	Page string // page
	Help string // help text
} // Footer information

type Menu struct {
	items            []Item       // Menu items
	itemsdisplay     []*Item      // Filtered items, used to display
	message          string       // Menu prompt words
	pageSize         int          // Quantity per page
	curPage          int          // Current page
	curIndex         int          // The currently selected menu item
	cursor           string       // The indicator cursor for the current selection
	cursorspacer     string       // Indicator cursor blank placeholder string
	showIndex        bool         // Display menu sequence number
	showNamelen      int          // Display menu name width
	showNotelen      int          // Display menu description width
	cache            bytes.Buffer // Output cache
	rowcount         int          // Display Area Rows
	pagecount        int          // Total number of pages
	lastpagestart    int          // Start index of the last page
	delimiter        string       // Field delimiter
	spacer           string       // Blank placeholder characters
	itemcount        int          // Menu item count
	itemdisplaycount int          // displayed menu items count
	selectedColor    string       // Selected item color
	lang             LangTxt      // Display language
	filter           string       // Filter strings
} // Menu structure

// Create a new menu object
func New() *Menu {
	return &Menu{
		items:            []Item{},
		itemsdisplay:     []*Item{},
		message:          "",
		pageSize:         10,
		curPage:          0,
		curIndex:         0,
		cursor:           ">>> ",
		cursorspacer:     " ",
		showIndex:        false,
		showNamelen:      0,
		showNotelen:      0,
		cache:            bytes.Buffer{},
		rowcount:         0,
		pagecount:        0,
		lastpagestart:    0,
		delimiter:        "   ",
		spacer:           " ",
		itemcount:        0,
		itemdisplaycount: 0,
		selectedColor:    "\033[34;1m",
		lang:             LangTxt{Cur: "Cur", Page: "Page", Help: "Arrow key move, <Esc> or <Ctrl>+C exit, <Enter> confirm"},
		filter:           "",
	}
}

// Add menu item
func (m *Menu) AddItem(i ...Item) {
	m.items = append(m.items, i...)
}

// Add menu item
//
// tag - Attached parameter;
// act - Execute function, function parameter is menu item
func (m *Menu) Add(name, note string, tag any, act func(*Item) error) {
	m.items = append(m.items, Item{
		idxstr: "",
		Name:   name,
		Note:   note,
		Act:    act,
		Tags:   tag,
	})
}

// Is the sequence number of menu items displayed
func (m *Menu) SetIndex(si bool) {
	m.showIndex = si
}

// display width of a string
func getWidth(s string) (w int) {
	e := regexp.MustCompile(`\033\[[\d;]+m`)
	s = e.ReplaceAllString(s, "")
	for _, r := range s {
		switch width.LookupRune(r).Kind() {
		case width.EastAsianFullwidth, width.EastAsianWide:
			w += 2
		case width.EastAsianHalfwidth, width.EastAsianNarrow, width.Neutral, width.EastAsianAmbiguous:
			w += 1
		}
	}
	return w
}

// Set the indicator for the current selection to be at the front of the menu
func (m *Menu) SetSelectIcon(p string) {
	m.cursor = p
}

// Set menu prompt messages
func (m *Menu) SetMessage(msg string) {
	m.message = msg
}

// Set the number of menu items displayed per page
func (m *Menu) SetPagesize(n int) {
	m.pageSize = n
}

// Set menu item field separator
func (m *Menu) SetDelimiter(s string) {
	m.delimiter = s
}

// Sort menu items. If the sort function is not executed,
// it will be displayed in the order of addition.
func (m *Menu) Sort(s func(i, j *Item) bool) {
	sort.SliceStable(m.items, func(i, j int) bool {
		return s(&(m.items[i]), &(m.items[j]))
	})
}

// Run menu, loop endlessly.
// Exit after selecting a menu item and executing it, or interrupt with CtrlC.
// Returns the currently selected menu item.
func (m *Menu) Run() (*Item, error) {
	if len(m.items) == 0 {
		return nil, errors.New("no menu item")
	}

	cursorHide(os.Stdout)
	m.configure()
	isloop := true
	for isloop {
		m.update()
		m.renderer()
		keytype, c := m.getInput()
		switch keytype {
		case KEY_escape, KEY_ctrlC:
			isloop = false
		case KEY_up:
			if m.curIndex > 0 {
				m.curIndex--
			}
		case KEY_down:
			if m.curIndex < m.itemdisplaycount-1 {
				m.curIndex++
			}
		case KEY_left, KEY_pageup:
			if m.curIndex >= m.pageSize {
				m.curIndex -= m.pageSize
			}
		case KEY_right, KEY_pagedown:
			if m.curIndex < m.lastpagestart {
				newidx := m.curIndex + m.pageSize
				if newidx < m.itemdisplaycount {
					m.curIndex = newidx
				} else {
					m.curIndex = m.itemdisplaycount - 1
				}
			}
		case KEY_enter:
			if m.itemdisplaycount > 0 {
				act := m.itemsdisplay[m.curIndex].Act
				if act != nil {
					m.rendererclear()
					cursorShow(os.Stdout)
					return m.itemsdisplay[m.curIndex], act(m.itemsdisplay[m.curIndex])
				}
			}
		case KEY_backspace:
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.reconfigure()
			}
		case KEY_filterstring:
			m.filter += c
			m.reconfigure()
		}
	}
	m.rendererclear()
	cursorShow(os.Stdout)
	return nil, errors.New("exit without selection")
}

// Update display cache
func (m *Menu) update() {
	m.cache.Reset()
	w := io.Writer(&m.cache)
	// head
	fmt.Fprintf(w, "%s   %s\033[K\n", m.message, m.filter)
	//
	curpage := m.curIndex / m.pageSize
	istart := curpage * m.pageSize
	for i := 0; i < m.pageSize; i++ {
		idx := istart + i
		if idx < m.itemdisplaycount {
			v := m.itemsdisplay[idx]
			if idx == m.curIndex {
				fmt.Fprint(w, m.cursor, m.selectedColor)
			} else {
				fmt.Fprint(w, m.cursorspacer, "\033[0m")
			}
			if m.showIndex {
				fmt.Fprintf(w, "%s", v.idxstr)
			}
			fmt.Fprintf(w, "%s%s%s%s\033[0m\033[K\n", m.delimiter, v.Name, m.delimiter, v.Note)
		} else {
			fmt.Fprint(w, "\033[K\n")
		}
	}
	// foot
	pageidx := m.curIndex/m.pageSize + 1
	fmt.Fprintf(w, "%s: (%d/%d) %s: (%d/%d) %s\033[K\n", m.lang.Cur, m.items[m.curIndex].idx, m.itemdisplaycount, m.lang.Page, pageidx, m.pagecount, m.lang.Help)
}

// Output cache to screen
func (m *Menu) renderer() {
	fmt.Print(m.cache.String())
	cursorMoveup(os.Stdout, m.rowcount)
}

// Clear menu display area
func (m *Menu) rendererclear() {
	fmt.Print(strings.Repeat("\033[K\n", m.rowcount))
	cursorMoveup(os.Stdout, m.rowcount)
}

// Necessary configuration before run
func (m *Menu) configure() {
	m.cursorspacer = strings.Repeat(m.spacer, getWidth(m.cursor))
	m.itemcount = len(m.items)
	m.pageSize = min(m.pageSize, m.itemcount)
	m.itemsdisplay = make([]*Item, 0, m.itemcount)

	m.rowcount = m.pageSize + 2
	s := fmt.Sprintf("%d", m.itemcount)
	showIndexlen := getWidth(s)

	anamelen := make([]int, m.itemcount)
	anotelen := make([]int, m.itemcount)
	for i := 0; i < m.itemcount; i++ {
		anamelen[i] = getWidth(m.items[i].Name)
		anotelen[i] = getWidth(m.items[i].Note)
	}

	m.showNamelen = slices.Max(anamelen)
	m.showNotelen = slices.Max(anotelen)

	idxfmt := fmt.Sprintf("[%%%dd]:", showIndexlen)
	for i := 0; i < m.itemcount; i++ {
		m.items[i].idx = i + 1
		m.items[i].idxstr = fmt.Sprintf(idxfmt, i+1)
		m.items[i].Name += strings.Repeat(m.spacer, m.showNamelen-anamelen[i])
		m.items[i].Note += strings.Repeat(m.spacer, m.showNotelen-anotelen[i])
	}

	m.reconfigure()
}

// 根据当前显示内容重新计算参数
func (m *Menu) reconfigure() {
	m.itemsdisplay = m.itemsdisplay[0:0]
	if m.filter == "" {
		for i := range m.items {
			m.itemsdisplay = append(m.itemsdisplay, &m.items[i])
		}
	} else {
		filterlower := strings.ToLower(m.filter)
		for i := range m.items {
			if strings.Contains(
				strings.ToLower(m.items[i].Name),
				filterlower,
			) {
				m.itemsdisplay = append(m.itemsdisplay, &m.items[i])
			}
		}
	}

	m.itemdisplaycount = len(m.itemsdisplay)
	m.pagecount = (m.itemdisplaycount-1)/m.pageSize + 1
	m.lastpagestart = (m.pagecount - 1) * m.pageSize
	m.curIndex = 0
}

// Set selected item color.
//
// Deprecated: use 'SetSelectedStyle'.
func (m *Menu) SetSelectedColor(c int) {
	if c > Color_min && c < Color_max {
		m.selectedColor = fmt.Sprintf("\033[%d;1m", c)
	} else {
		m.selectedColor = fmt.Sprintf("\033[%d;1m", Color_Blue)
	}
}

// Set selected item style. Use regular expression `^\033\[[\d;]+m$` to check.
func (m *Menu) SetSelectedStyle(styles ...string) {
	re := regexp.MustCompile(`^\033\[[\d;]+m$`)
	m.selectedColor = ""
	for _, s := range styles {
		if re.MatchString(s) {
			m.selectedColor += s
		}
	}
	if m.selectedColor == "" {
		m.selectedColor = Style_Blue
	}
}

// Set language text for foot
func (m *Menu) Seti18n(cur, page, help string) {
	m.lang.Cur = cur
	m.lang.Page = page
	m.lang.Help = help
}
