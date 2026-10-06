package gear

import "testing"

func TestScoreWithoutPercent(t *testing.T) {
	if got := base(Row{AP: 100}); got != 100 {
		t.Fatalf("无百分比加成时应是 100，得到 %d", got)
	}
}

func TestBestPicksHigherScore(t *testing.T) {
	rows := []Row{
		{Slot: "weapon", Item: "axe", AP: 100},
		{Slot: "weapon", Item: "bow", AP: 50},
	}
	picks := Solve(rows, &Counter{})
	if len(picks) != 1 || picks[0].Item != "axe" {
		t.Fatalf("应挑出 axe，得到 %+v", picks)
	}
}

func TestBestKeepsSlotsApart(t *testing.T) {
	rows := []Row{
		{Slot: "head", Item: "helm", AP: 10},
		{Slot: "feet", Item: "boots", AP: 20},
	}
	if picks := Solve(rows, &Counter{}); len(picks) != 2 {
		t.Fatalf("两个槽位应各出一条，得到 %+v", picks)
	}
}

func TestBestOnEmptyInput(t *testing.T) {
	if picks := Solve(nil, &Counter{}); len(picks) != 0 {
		t.Fatalf("空表应无推荐，得到 %+v", picks)
	}
}

func TestUniqueOnDistinctRows(t *testing.T) {
	rows := []Row{{Slot: "head", Item: "a"}, {Slot: "head", Item: "b"}}
	if got := Unique(rows); got != 2 {
		t.Fatalf("两件不同装备应算 2，得到 %d", got)
	}
}
