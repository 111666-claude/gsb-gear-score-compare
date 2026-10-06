// Command gear 读装备 CSV 并打印推荐表。
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"example.com/gear-score-compare/gear"
)

func atoi(s string) int {
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return v
}

func readCSV(path string) ([]gear.Row, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	records, err := csv.NewReader(fh).ReadAll()
	if err != nil {
		return nil, err
	}
	var rows []gear.Row
	for i, rec := range records {
		if i == 0 || len(rec) < 7 {
			continue
		}
		rows = append(rows, gear.Row{
			Slot: strings.TrimSpace(rec[0]), Item: strings.TrimSpace(rec[1]),
			Set: strings.TrimSpace(rec[2]), AP: atoi(rec[3]), APPct: atoi(rec[4]),
			Crit: atoi(rec[5]), Haste: atoi(rec[6]),
		})
	}
	return rows, nil
}

func render(picks []gear.Pick) string {
	rows := [][4]string{{"slot", "item", "set", "score"}}
	for _, p := range picks {
		rows = append(rows, [4]string{p.Slot, p.Item, p.Set, strconv.Itoa(p.Score)})
	}
	widths := [4]int{}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	var b strings.Builder
	for _, row := range rows {
		for i, cell := range row {
			b.WriteString(cell)
			if i < 3 {
				b.WriteString(strings.Repeat(" ", widths[i]-len(cell)+2))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func emit(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main() {
	csvPath := flag.String("csv", "", "装备 CSV 路径")
	sample := flag.String("sample", "", "set|pct|tie|dupe|work")
	flag.Parse()

	switch *sample {
	case "set":
		rows := []gear.Row{
			{Slot: "a", Item: "a1", Set: "S", AP: 100},
			{Slot: "b", Item: "b1", Set: "S", AP: 100},
			{Slot: "c", Item: "c1", Set: "S", AP: 100},
			{Slot: "d", Item: "d1", Set: "S", AP: 100},
			{Slot: "d", Item: "z1", Set: "T", AP: 103},
		}
		picks := gear.Solve(rows, &gear.Counter{})
		for _, p := range picks {
			if p.Slot == "d" {
				emit(p)
			}
		}
	case "pct":
		rows := []gear.Row{{Slot: "weapon", Item: "axe", AP: 100, APPct: 50, Crit: 10}}
		emit(gear.Solve(rows, &gear.Counter{})[0])
	case "tie":
		rows := []gear.Row{
			{Slot: "weapon", Item: "zeta", AP: 100},
			{Slot: "weapon", Item: "alpha", AP: 100},
		}
		emit(gear.Solve(rows, &gear.Counter{})[0])
	case "dupe":
		rows := []gear.Row{
			{Slot: "weapon", Item: "axe", AP: 10},
			{Slot: "weapon", Item: "axe", AP: 99},
		}
		emit(gear.Solve(rows, &gear.Counter{})[0])
	case "work":
		var rows []gear.Row
		for i := 0; i < 2000; i++ {
			rows = append(rows, gear.Row{Slot: fmt.Sprintf("s%02d", i%50),
				Item: fmt.Sprintf("i%04d", i), AP: i})
		}
		c := &gear.Counter{}
		gear.Solve(rows, c)
		emit(map[string]int{"rows": 2000, "scanned": c.Scanned})
	case "":
		if *csvPath == "" {
			flag.Usage()
			os.Exit(2)
		}
		rows, err := readCSV(*csvPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Print(render(gear.Solve(rows, &gear.Counter{})))
	default:
		fmt.Fprintln(os.Stderr, "未知场景："+*sample)
		os.Exit(2)
	}
}
