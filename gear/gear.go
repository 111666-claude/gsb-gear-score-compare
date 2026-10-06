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
	return r.AP*(100+r.APPct)/100 + r.Crit*2 + r.Haste*3
}

// dedup 去掉重复的（槽位, 装备）行，只保留第一次出现的那行。
func dedup(rows []Row) []Row {
	seen := map[string]bool{}
	out := make([]Row, 0, len(rows))
	for _, r := range rows {
		key := r.Slot + "\x00" + r.Item
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, r)
	}
	return out
}

// pickWinners 按 score 给每个槽位挑一行；同分取 item 名升序靠前者。
func pickWinners(rows []Row, c *Counter, score func(Row) int) map[string]Row {
	win := map[string]Row{}
	for _, r := range rows {
		c.Scanned++
		cur, ok := win[r.Slot]
		if !ok {
			win[r.Slot] = r
			continue
		}
		rScore, curScore := score(r), score(cur)
		if rScore > curScore || (rScore == curScore && r.Item < cur.Item) {
			win[r.Slot] = r
		}
	}
	return win
}

// Solve 给每个槽位挑一件装备。
func Solve(rows []Row, c *Counter) []Pick {
	uniq := dedup(rows)

	round1 := pickWinners(uniq, c, base)
	counts := map[string]int{}
	for _, r := range round1 {
		if r.Set != "" {
			counts[r.Set]++
		}
	}
	qualify := map[string]bool{}
	for set, n := range counts {
		if n >= SetSize {
			qualify[set] = true
		}
	}

	boosted := func(r Row) int {
		score := base(r)
		if qualify[r.Set] {
			score += SetBonus
		}
		return score
	}
	round2 := pickWinners(uniq, c, boosted)

	slots := make([]string, 0, len(round2))
	for slot := range round2 {
		slots = append(slots, slot)
	}
	sort.Strings(slots)

	picks := make([]Pick, 0, len(slots))
	for _, slot := range slots {
		r := round2[slot]
		picks = append(picks, Pick{Slot: slot, Item: r.Item, Set: r.Set, Score: boosted(r)})
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
