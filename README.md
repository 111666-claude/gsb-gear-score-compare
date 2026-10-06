# gear-score-compare

装备评分与换装推荐：读一张装备 CSV，给每个槽位挑出推荐件，打印对齐表格。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/gear --csv items.csv
go run ./cmd/gear --sample=set
go run ./cmd/gear --sample=pct
go run ./cmd/gear --sample=tie
go run ./cmd/gear --sample=dupe
go run ./cmd/gear --sample=work
```

CSV 列为 `slot`、`item`、`set`、`ap`、`ap_pct`、`crit`、`haste`，第一行是表头。

## 口径

- **基础分**：`ap` 先按 `ap_pct` 放大，再叠加 `crit×2` 与 `haste×3`。
- **套装加成**：按基础分先选一轮，数出选中集合里各套装的件数；件数达到 3 的套装，其选中件 +10 分；再按加成后的分数选第二轮，最终以第二轮为准。
- **推荐**：每个槽位取第二轮里分数最高的一件；同分取 `item` 名升序里靠前的那件。
- **去重**：同一（槽位, 装备）出现多行时只认第一行，后出现的同名行忽略。
- **槽位**：推荐表按 `slot` 名升序打印。

## 不变量

- 每个槽位最多一条推荐。
- 同一份输入反复求解结果一致。
- 同名装备不会因为重复行而被抬分。
- `scanned` 不随行数平方级放大：两千行的 `scanned` 不超过 9000。

## 输出契约

推荐表列顺序为 `slot`、`item`、`set`、`score`。
`--sample=set|pct|tie|dupe` 只打印一条 `{"slot": ..., "item": ..., "set": ..., "score": N}`；
`--sample=work` 只打印 `{"rows": N, "scanned": N}`。
