// Package gear 给装备打分并挑出每个槽位的推荐件。
package gear

import "sort"

// 套装加成门槛与数值。
const (
	SetSize = 3
	SetBonus = 10
)

// Row 是装备表里的一行。
type Row struct {
	Slot  string
	Item  string
	Set   string
	AP    int
	APPct int
	Crit  int
	Haste int
}

// Counter 统计本轮比较过的行数。
type Counter struct {
	Scanned int
}

// Pick 是一个槽位的推荐结果。
type Pick struct {
	Slot  string `json:"slot"`
	Item  string `json:"item"`
	Set   string `json:"set"`
	Score int    `json:"score"`
}

func base(r Row) int {
	return (r.AP + r.Crit*2 + r.Haste*3) * (100 + r.APPct) / 100
}

func slots(rows []Row) []string {
	var out []string
	for _, r := range rows {
		dup := false
		for _, s := range out {
			if s == r.Slot {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, r.Slot)
		}
	}
	sort.Strings(out)
	return out
}

// Solve 给每个槽位挑一件装备。
func Solve(rows []Row, c *Counter) []Pick {
	var picks []Pick
	for _, slot := range slots(rows) {
		best := Row{}
		found := false
		for _, r := range rows {
			c.Scanned++
			if r.Slot != slot {
				continue
			}
			if !found || base(r) > base(best) {
				best, found = r, true
			}
		}
		if found {
			picks = append(picks, Pick{Slot: slot, Item: best.Item, Set: best.Set, Score: base(best)})
		}
	}
	return picks
}

// Unique 统计不同的（槽位, 装备）组合数。
func Unique(rows []Row) int {
	seen := map[string]bool{}
	n := 0
	for _, r := range rows {
		key := r.Slot + "\x00" + r.Item
		if !seen[key] {
			seen[key] = true
			n++
		}
	}
	return n
}
