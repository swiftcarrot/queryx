# Results

- orm-benchmark
```
golang-mysql-Reports:
InsertAll
queryx:         24      68688983 ns/op  279390 B/op     3296 allocs/op

Create
queryx:         50      24818608 ns/op  16685 B/op      253 allocs/op

Update
queryx:         1215    968346 ns/op    7980 B/op       114 allocs/op

Read
queryx:         1201    976828 ns/op    9196 B/op       155 allocs/op

ReadSlice
queryx:         762     2458304 ns/op   386941 B/op     12405 allocs/op


golang-sqlite-Reports:
ReadSlice
queryx:         358     3447322 ns/op   398725 B/op     16991 allocs/op

InsertAll
queryx:         48      31670682 ns/op  246893 B/op     3898 allocs/op

Create
queryx:         72      18938976 ns/op  11278 B/op      238 allocs/op

Update
queryx:         499     3277440 ns/op   8052 B/op       127 allocs/op

Read
queryx:         903     1488458 ns/op   9252 B/op       208 allocs/op


golang-postgresql-Reports:
InsertAll
queryx:         134     8033703 ns/op   289751 B/op     3505 allocs/op

Create
queryx:         194     5316417 ns/op   10501 B/op      184 allocs/op

Update
queryx:         208     5298013 ns/op   7944 B/op       116 allocs/op

Read
queryx:         1186    1187729 ns/op   9278 B/op       157 allocs/op

ReadSlice
queryx:         1008    2237629 ns/op   348255 B/op     10661 allocs/op

typescript-mysql-Reports:
create: x 790 ops/sec ±191.34% (36 runs sampled)
insertAll: x 591 ops/sec ±60.50% (20 runs sampled)
find x 4,631 ops/sec ±24.66% (34 runs sampled)
update: x 760 ops/sec ±87.89% (10 runs sampled)

typescript-postgresql-Reports:
create: x 1,071 ops/sec ±86.40% (9 runs sampled)
insertAll: x 8,521 ops/sec ±21.30% (35 runs sampled)
find: x 4,269 ops/sec ±16.80% (52 runs sampled)
update: x 1,860 ops/sec ±46.43% (38 runs sampled)

typescript-sqlite-Reports:
create:x 25,930 ops/sec +18.77%(44 runs sampled)
insertAll:x14,361 ops/sec +66.85%(29 runs sampled)
find x 1.087 ops/sec +78.85%(13 runs sampled)
update:x 2,225 ops/sec +18.76%(36 runs sampled)

```
