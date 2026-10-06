// Package gear 给装备打分并挑出每个槽位的推荐件。
package gear

import "sort"

// 套装加成门槛与数值。
const (
	SetSize  = 3
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

func scored(r Row, bonus map[string]bool) int {
	s := base(r)
	if bonus[r.Set] {
		s += SetBonus
	}
	return s
}

// bestInSlot 在一个槽位的候选里挑分最高的一件，同分取 item 名升序靠前的。
func bestInSlot(cands []Row, bonus map[string]bool) Row {
	best := cands[0]
	for _, r := range cands[1:] {
		rs, bs := scored(r, bonus), scored(best, bonus)
		if rs > bs || (rs == bs && r.Item < best.Item) {
			best = r
		}
	}
	return best
}

// Solve 给每个槽位挑一件装备。
func Solve(rows []Row, c *Counter) []Pick {
	// 同一（槽位, 装备）只认第一行，同时单遍按槽位归组。
	seen := map[string]bool{}
	bySlot := map[string][]Row{}
	var slotOrder []string
	for _, r := range rows {
		key := r.Slot + "\x00" + r.Item
		if seen[key] {
			continue
		}
		seen[key] = true
		c.Scanned++
		if _, ok := bySlot[r.Slot]; !ok {
			slotOrder = append(slotOrder, r.Slot)
		}
		bySlot[r.Slot] = append(bySlot[r.Slot], r)
	}
	sort.Strings(slotOrder)

	noBonus := map[string]bool{}

	// 第一轮：按基础分选，数选中件里各套装的件数。
	setCount := map[string]int{}
	for _, slot := range slotOrder {
		best := bestInSlot(bySlot[slot], noBonus)
		if best.Set != "" {
			setCount[best.Set]++
		}
	}

	// 凑满 SetSize 件的套装，第二轮给该套装的候选件加分。
	bonus := map[string]bool{}
	for set, n := range setCount {
		if n >= SetSize {
			bonus[set] = true
		}
	}

	// 第二轮：按加成后的分数最终定推荐。
	var picks []Pick
	for _, slot := range slotOrder {
		best := bestInSlot(bySlot[slot], bonus)
		picks = append(picks, Pick{Slot: slot, Item: best.Item, Set: best.Set, Score: scored(best, bonus)})
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
