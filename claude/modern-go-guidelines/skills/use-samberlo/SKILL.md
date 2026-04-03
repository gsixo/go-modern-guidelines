---
name: use-samberlo
description: Apply samber/lo library
---

**When writing Go code**, use ALL features from this document:

This block used for functionality in golang

imports "github.com/samber/lo"

required Go 1.18+

use this tips to refactor or write code

if u see something similar to before use function in after to simplify code

---

## Filter

**Before:**
```go
numbers := []int{1, 2, 3, 4}
var even []int
for _, x := range numbers {
    if x%2 == 0 {
        even = append(even, x)
    }
}
// []int{2, 4}
```

**After:**
```go
even := lo.Filter([]int{1, 2, 3, 4}, func(x int, index int) bool {
    return x%2 == 0
})
// []int{2, 4}
```

---

## Map

**Before:**
```go
ids := []int64{1, 2, 3, 4}
strs := make([]string, len(ids))
for i, x := range ids {
    strs[i] = strconv.FormatInt(x, 10)
}
// []string{"1", "2", "3", "4"}
```

**After:**
```go
lo.Map([]int64{1, 2, 3, 4}, func(x int64, index int) string {
    return strconv.FormatInt(x, 10)
})
// []string{"1", "2", "3", "4"}
```

---

## FilterMap

**Before:**
```go
items := []string{"cpu", "gpu", "mouse", "keyboard"}
var matching []string
for _, x := range items {
    if strings.HasSuffix(x, "pu") {
        matching = append(matching, "xpu")
    }
}
// []string{"xpu", "xpu"}
```

**After:**
```go
matching := lo.FilterMap([]string{"cpu", "gpu", "mouse", "keyboard"}, func(x string, _ int) (string, bool) {
    if strings.HasSuffix(x, "pu") {
        return "xpu", true
    }
    return "", false
})
// []string{"xpu", "xpu"}
```

---

## FlatMap

**Before:**
```go
nums := []int64{0, 1, 2}
var result []string
for _, x := range nums {
    s := strconv.FormatInt(x, 10)
    result = append(result, s, s)
}
// []string{"0", "0", "1", "1", "2", "2"}
```

**After:**
```go
lo.FlatMap([]int64{0, 1, 2}, func(x int64, _ int) []string {
    return []string{
        strconv.FormatInt(x, 10),
        strconv.FormatInt(x, 10),
    }
})
// []string{"0", "0", "1", "1", "2", "2"}
```

---

## Reduce

**Before:**
```go
nums := []int{1, 2, 3, 4}
sum := 0
for _, item := range nums {
    sum += item
}
// 10
```

**After:**
```go
sum := lo.Reduce([]int{1, 2, 3, 4}, func(agg int, item int, _ int) int {
    return agg + item
}, 0)
// 10
```

---

## ReduceRight

**Before:**
```go
slices := [][]int{{0, 1}, {2, 3}, {4, 5}}
var result []int
for i := len(slices) - 1; i >= 0; i-- {
    result = append(result, slices[i]...)
}
// []int{4, 5, 2, 3, 0, 1}
```

**After:**
```go
result := lo.ReduceRight([][]int{{0, 1}, {2, 3}, {4, 5}}, func(agg []int, item []int, _ int) []int {
    return append(agg, item...)
}, []int{})
// []int{4, 5, 2, 3, 0, 1}
```

---

## ForEach

**Before:**
```go
words := []string{"hello", "world"}
for _, x := range words {
    println(x)
}
// prints "hello\nworld\n"
```

**After:**
```go
lo.ForEach([]string{"hello", "world"}, func(x string, _ int) {
    println(x)
})
// prints "hello\nworld\n"
```

---

## ForEachWhile

**Before:**
```go
list := []int64{1, 2, -42, 4}
for _, x := range list {
    if x < 0 {
        break
    }
    fmt.Println(x)
}
// 1
// 2
```

**After:**
```go
list := []int64{1, 2, -42, 4}

lo.ForEachWhile(list, func(x int64, _ int) bool {
    if x < 0 {
        return false
    }
    fmt.Println(x)
    return true
})
// 1
// 2
```

---

## Times

**Before:**
```go
n := 3
result := make([]string, n)
for i := 0; i < n; i++ {
    result[i] = strconv.FormatInt(int64(i), 10)
}
// []string{"0", "1", "2"}
```

**After:**
```go
lo.Times(3, func(i int) string {
    return strconv.FormatInt(int64(i), 10)
})
// []string{"0", "1", "2"}
```

---

## Uniq

**Before:**
```go
input := []int{1, 2, 2, 1}
seen := make(map[int]struct{})
var result []int
for _, v := range input {
    if _, ok := seen[v]; !ok {
        seen[v] = struct{}{}
        result = append(result, v)
    }
}
// []int{1, 2}
```

**After:**
```go
uniqValues := lo.Uniq([]int{1, 2, 2, 1})
// []int{1, 2}
```

---

## UniqBy

**Before:**
```go
input := []int{0, 1, 2, 3, 4, 5}
seen := make(map[int]struct{})
var result []int
for _, i := range input {
    key := i % 3
    if _, ok := seen[key]; !ok {
        seen[key] = struct{}{}
        result = append(result, i)
    }
}
// []int{0, 1, 2}
```

**After:**
```go
uniqValues := lo.UniqBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// []int{0, 1, 2}
```

---

## GroupBy

**Before:**
```go
input := []int{0, 1, 2, 3, 4, 5}
groups := make(map[int][]int)
for _, i := range input {
    key := i % 3
    groups[key] = append(groups[key], i)
}
// map[int][]int{0: []int{0, 3}, 1: []int{1, 4}, 2: []int{2, 5}}
```

**After:**
```go
groups := lo.GroupBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// map[int][]int{0: []int{0, 3}, 1: []int{1, 4}, 2: []int{2, 5}}
```

---

## Chunk

**Before:**
```go
input := []int{0, 1, 2, 3, 4, 5}
size := 2
var result [][]int
for i := 0; i < len(input); i += size {
    end := i + size
    if end > len(input) {
        end = len(input)
    }
    result = append(result, input[i:end])
}
// [][]int{{0, 1}, {2, 3}, {4, 5}}
```

**After:**
```go
lo.Chunk([]int{0, 1, 2, 3, 4, 5}, 2)
// [][]int{{0, 1}, {2, 3}, {4, 5}}

lo.Chunk([]int{0, 1, 2, 3, 4, 5, 6}, 2)
// [][]int{{0, 1}, {2, 3}, {4, 5}, {6}}
```

---

## PartitionBy

**Before:**
```go
input := []int{-2, -1, 0, 1, 2, 3, 4, 5}
groupOrder := []string{}
groups := map[string][]int{}
for _, x := range input {
    var key string
    if x < 0 {
        key = "negative"
    } else if x%2 == 0 {
        key = "even"
    } else {
        key = "odd"
    }
    if _, exists := groups[key]; !exists {
        groupOrder = append(groupOrder, key)
    }
    groups[key] = append(groups[key], x)
}
var partitions [][]int
for _, k := range groupOrder {
    partitions = append(partitions, groups[k])
}
// [][]int{{-2, -1}, {0, 2, 4}, {1, 3, 5}}
```

**After:**
```go
partitions := lo.PartitionBy([]int{-2, -1, 0, 1, 2, 3, 4, 5}, func(x int) string {
    if x < 0 {
        return "negative"
    } else if x%2 == 0 {
        return "even"
    }
    return "odd"
})
// [][]int{{-2, -1}, {0, 2, 4}, {1, 3, 5}}
```

---

## Flatten

**Before:**
```go
nested := [][]int{{0, 1}, {2, 3, 4, 5}}
var flat []int
for _, inner := range nested {
    flat = append(flat, inner...)
}
// []int{0, 1, 2, 3, 4, 5}
```

**After:**
```go
flat := lo.Flatten([][]int{{0, 1}, {2, 3, 4, 5}})
// []int{0, 1, 2, 3, 4, 5}
```

---

## Concat

**Before:**
```go
a := []int{1, 2}
b := []int{3, 4}
result := make([]int, 0, len(a)+len(b))
result = append(result, a...)
result = append(result, b...)
// []int{1, 2, 3, 4}
```

**After:**
```go
slice := lo.Concat([]int{1, 2}, []int{3, 4})
// []int{1, 2, 3, 4}

slice := lo.Concat(nil, []int{1, 2}, nil, []int{3, 4}, nil)
// []int{1, 2, 3, 4}
```

---

## Interleave

**Before:**
```go
a := []int{1, 4, 7}
b := []int{2, 5, 8}
c := []int{3, 6, 9}
result := make([]int, 0, len(a)+len(b)+len(c))
for i := 0; i < len(a); i++ {
    result = append(result, a[i], b[i], c[i])
}
// []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
```

**After:**
```go
interleaved := lo.Interleave([]int{1, 4, 7}, []int{2, 5, 8}, []int{3, 6, 9})
// []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

interleaved := lo.Interleave([]int{1}, []int{2, 5, 8}, []int{3, 6}, []int{4, 7, 9, 10})
// []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
```

---

## Shuffle

**Before:**
```go
import "math/rand"

list := []int{0, 1, 2, 3, 4, 5}
rand.Shuffle(len(list), func(i, j int) {
    list[i], list[j] = list[j], list[i]
})
// list is now shuffled in place
```

**After:**
```go
import lom "github.com/samber/lo/mutable"

list := []int{0, 1, 2, 3, 4, 5}
lom.Shuffle(list)

list
// []int{1, 4, 0, 3, 5, 2}
```

---

## Reverse

**Before:**
```go
list := []int{0, 1, 2, 3, 4, 5}
for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
    list[i], list[j] = list[j], list[i]
}
// []int{5, 4, 3, 2, 1, 0}
```

**After:**
```go
import lom "github.com/samber/lo/mutable"

list := []int{0, 1, 2, 3, 4, 5}
lom.Reverse(list)

list
// []int{5, 4, 3, 2, 1, 0}
```

---

## Fill

**Before:**
```go
type foo struct{ bar string }

slice := []foo{{"a"}, {"a"}}
fill := foo{"b"}
for i := range slice {
    slice[i] = fill
}
// []foo{foo{"b"}, foo{"b"}}
```

**After:**
```go
type foo struct {
  bar string
}

func (f foo) Clone() foo {
  return foo{f.bar}
}

initializedSlice := lo.Fill([]foo{foo{"a"}, foo{"a"}}, foo{"b"})
// []foo{foo{"b"}, foo{"b"}}
```

---

## Repeat

**Before:**
```go
type foo struct{ bar string }

n := 2
slice := make([]foo, n)
for i := range slice {
    slice[i] = foo{"a"}
}
// []foo{foo{"a"}, foo{"a"}}
```

**After:**
```go
type foo struct {
  bar string
}

func (f foo) Clone() foo {
  return foo{f.bar}
}

slice := lo.Repeat(2, foo{"a"})
// []foo{foo{"a"}, foo{"a"}}
```

---

## RepeatBy

**Before:**
```go
n := 5
slice := make([]string, n)
for i := 0; i < n; i++ {
    slice[i] = strconv.FormatInt(int64(math.Pow(float64(i), 2)), 10)
}
// []string{"0", "1", "4", "9", "16"}
```

**After:**
```go
slice := lo.RepeatBy(0, func (i int) string {
    return strconv.FormatInt(int64(math.Pow(float64(i), 2)), 10)
})
// []string{}

slice := lo.RepeatBy(5, func(i int) string {
    return strconv.FormatInt(int64(math.Pow(float64(i), 2)), 10)
})
// []string{"0", "1", "4", "9", "16"}
```

---

## KeyBy

**Before:**
```go
input := []string{"a", "aa", "aaa"}
m := make(map[int]string, len(input))
for _, str := range input {
    m[len(str)] = str
}
// map[int]string{1: "a", 2: "aa", 3: "aaa"}
```

**After:**
```go
m := lo.KeyBy([]string{"a", "aa", "aaa"}, func(str string) int {
    return len(str)
})
// map[int]string{1: "a", 2: "aa", 3: "aaa"}
```

---

## SliceToMap

**Before:**
```go
type foo struct{ baz string; bar int }
in := []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}}

aMap := make(map[string]int, len(in))
for _, f := range in {
    aMap[f.baz] = f.bar
}
// map[string]int{"apple": 1, "banana": 2}
```

**After:**
```go
in := []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}}

aMap := lo.SliceToMap(in, func (f *foo) (string, int) {
    return f.baz, f.bar
})
// map[string][int]{ "apple":1, "banana":2 }
```

---

## Take

**Before:**
```go
input := []int{0, 1, 2, 3, 4, 5}
n := 3
if n > len(input) {
    n = len(input)
}
result := input[:n]
// []int{0, 1, 2}
```

**After:**
```go
l := lo.Take([]int{0, 1, 2, 3, 4, 5}, 3)
// []int{0, 1, 2}

l := lo.Take([]int{0, 1, 2}, 5)
// []int{0, 1, 2}
```

---

## TakeWhile

**Before:**
```go
input := []int{0, 1, 2, 3, 4, 5}
var result []int
for _, val := range input {
    if !(val < 3) {
        break
    }
    result = append(result, val)
}
// []int{0, 1, 2}
```

**After:**
```go
l := lo.TakeWhile([]int{0, 1, 2, 3, 4, 5}, func(val int) bool {
    return val < 3
})
// []int{0, 1, 2}

l := lo.TakeWhile([]string{"a", "aa", "aaa", "aa"}, func(val string) bool {
    return len(val) <= 2
})
// []string{"a", "aa"}
```

---

## Drop

**Before:**
```go
input := []int{0, 1, 2, 3, 4, 5}
n := 2
if n > len(input) {
    n = len(input)
}
result := input[n:]
// []int{2, 3, 4, 5}
```

**After:**
```go
l := lo.Drop([]int{0, 1, 2, 3, 4, 5}, 2)
// []int{2, 3, 4, 5}
```

---

## DropRight

**Before:**
```go
input := []int{0, 1, 2, 3, 4, 5}
n := 2
end := len(input) - n
if end < 0 {
    end = 0
}
result := input[:end]
// []int{0, 1, 2, 3}
```

**After:**
```go
l := lo.DropRight([]int{0, 1, 2, 3, 4, 5}, 2)
// []int{0, 1, 2, 3}
```

---

## DropWhile

**Before:**
```go
input := []string{"a", "aa", "aaa", "aa", "aa"}
i := 0
for i < len(input) && len(input[i]) <= 2 {
    i++
}
result := input[i:]
// []string{"aaa", "aa", "aa"}
```

**After:**
```go
l := lo.DropWhile([]string{"a", "aa", "aaa", "aa", "aa"}, func(val string) bool {
    return len(val) <= 2
})
// []string{"aaa", "aa", "aa"}
```

---

## Reject

**Before:**
```go
numbers := []int{1, 2, 3, 4}
var odd []int
for _, x := range numbers {
    if x%2 != 0 {
        odd = append(odd, x)
    }
}
// []int{1, 3}
```

**After:**
```go
odd := lo.Reject([]int{1, 2, 3, 4}, func(x int, _ int) bool {
    return x%2 == 0
})
// []int{1, 3}
```

---

## Count

**Before:**
```go
input := []int{1, 5, 1}
target := 1
count := 0
for _, v := range input {
    if v == target {
        count++
    }
}
// 2
```

**After:**
```go
count := lo.Count([]int{1, 5, 1}, 1)
// 2
```

---

## CountBy

**Before:**
```go
input := []int{1, 5, 1}
count := 0
for _, i := range input {
    if i < 4 {
        count++
    }
}
// 2
```

**After:**
```go
count := lo.CountBy([]int{1, 5, 1}, func(i int) bool {
    return i < 4
})
// 2
```

---

## CountValues

**Before:**
```go
input := []string{"foo", "bar", "bar"}
counts := make(map[string]int)
for _, v := range input {
    counts[v]++
}
// map[string]int{"foo": 1, "bar": 2}
```

**After:**
```go
lo.CountValues([]int{1, 2, 2})
// map[int]int{1: 1, 2: 2}

lo.CountValues([]string{"foo", "bar", "bar"})
// map[string]int{"foo": 1, "bar": 2}
```

---

## Subset

**Before:**
```go
in := []int{0, 1, 2, 3, 4}
start := 2
length := 3
end := start + length
if end > len(in) {
    end = len(in)
}
sub := in[start:end]
// []int{2, 3, 4}
```

**After:**
```go
in := []int{0, 1, 2, 3, 4}

sub := lo.Subset(in, 2, 3)
// []int{2, 3, 4}

sub := lo.Subset(in, -4, 3)
// []int{1, 2, 3}

sub := lo.Subset(in, -2, math.MaxUint)
// []int{3, 4}
```

---

## Slice

**Before:**
```go
in := []int{0, 1, 2, 3, 4}
start, end := 2, 3
if start > len(in) {
    start = len(in)
}
if end > len(in) {
    end = len(in)
}
if start > end {
    start = end
}
result := in[start:end]
// []int{2}
```

**After:**
```go
in := []int{0, 1, 2, 3, 4}

slice := lo.Slice(in, 0, 5)
// []int{0, 1, 2, 3, 4}

slice := lo.Slice(in, 2, 3)
// []int{2}

slice := lo.Slice(in, 2, 6)
// []int{2, 3, 4}

slice := lo.Slice(in, 4, 3)
// []int{}
```

---

## Replace

**Before:**
```go
in := []int{0, 1, 0, 1, 2, 3, 0}
old, newVal, n := 0, 42, 1
result := make([]int, len(in))
copy(result, in)
replaced := 0
for i, v := range result {
    if v == old && (n < 0 || replaced < n) {
        result[i] = newVal
        replaced++
    }
}
// []int{42, 1, 0, 1, 2, 3, 0}
```

**After:**
```go
in := []int{0, 1, 0, 1, 2, 3, 0}

slice := lo.Replace(in, 0, 42, 1)
// []int{42, 1, 0, 1, 2, 3, 0}

slice := lo.Replace(in, 0, 42, 2)
// []int{42, 1, 42, 1, 2, 3, 0}

slice := lo.Replace(in, 0, 42, -1)
// []int{42, 1, 42, 1, 2, 3, 42}
```

---

## ReplaceAll

**Before:**
```go
in := []int{0, 1, 0, 1, 2, 3, 0}
result := make([]int, len(in))
for i, v := range in {
    if v == 0 {
        result[i] = 42
    } else {
        result[i] = v
    }
}
// []int{42, 1, 42, 1, 2, 3, 42}
```

**After:**
```go
in := []int{0, 1, 0, 1, 2, 3, 0}

slice := lo.ReplaceAll(in, 0, 42)
// []int{42, 1, 42, 1, 2, 3, 42}

slice := lo.ReplaceAll(in, -1, 42)
// []int{0, 1, 0, 1, 2, 3, 0}
```

---

## Compact

**Before:**
```go
in := []string{"", "foo", "", "bar", ""}
var result []string
for _, v := range in {
    if v != "" {
        result = append(result, v)
    }
}
// []string{"foo", "bar"}
```

**After:**
```go
in := []string{"", "foo", "", "bar", ""}

slice := lo.Compact(in)
// []string{"foo", "bar"}
```

---

## IsSorted

**Before:**
```go
slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
sorted := true
for i := 1; i < len(slice); i++ {
    if slice[i] < slice[i-1] {
        sorted = false
        break
    }
}
// true
```

**After:**
```go
slice := lo.IsSorted([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
// true
```

---

## IsSortedBy

**Before:**
```go
strs := []string{"a", "bb", "ccc"}
sorted := true
for i := 1; i < len(strs); i++ {
    if len(strs[i]) < len(strs[i-1]) {
        sorted = false
        break
    }
}
// true
```

**After:**
```go
slice := lo.IsSortedBy([]string{"a", "bb", "ccc"}, func(s string) int {
    return len(s)
})
// true
```

---

## Keys

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2}
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
// []string{"foo", "bar"} (order not guaranteed)
```

**After:**
```go
keys := lo.Keys(map[string]int{"foo": 1, "bar": 2})
// []string{"foo", "bar"}

keys := lo.Keys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
// []string{"foo", "bar", "baz"}
```

---

## Values

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2}
values := make([]int, 0, len(m))
for _, v := range m {
    values = append(values, v)
}
// []int{1, 2} (order not guaranteed)
```

**After:**
```go
values := lo.Values(map[string]int{"foo": 1, "bar": 2})
// []int{1, 2}

values := lo.Values(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
// []int{1, 2, 3}
```

---

## PickBy

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2, "baz": 3}
result := make(map[string]int)
for k, v := range m {
    if v%2 == 1 {
        result[k] = v
    }
}
// map[string]int{"foo": 1, "baz": 3}
```

**After:**
```go
m := lo.PickBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
    return value%2 == 1
})
// map[string]int{"foo": 1, "baz": 3}
```

---

## PickByKeys

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2, "baz": 3}
keys := []string{"foo", "baz"}
result := make(map[string]int)
allowed := make(map[string]struct{}, len(keys))
for _, k := range keys {
    allowed[k] = struct{}{}
}
for k, v := range m {
    if _, ok := allowed[k]; ok {
        result[k] = v
    }
}
// map[string]int{"foo": 1, "baz": 3}
```

**After:**
```go
m := lo.PickByKeys(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})
// map[string]int{"foo": 1, "baz": 3}
```

---

## OmitBy

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2, "baz": 3}
result := make(map[string]int)
for k, v := range m {
    if v%2 != 1 {
        result[k] = v
    }
}
// map[string]int{"bar": 2}
```

**After:**
```go
m := lo.OmitBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
    return value%2 == 1
})
// map[string]int{"bar": 2}
```

---

## Entries

**Before:**
```go
type Entry struct {
    Key   string
    Value int
}
m := map[string]int{"foo": 1, "bar": 2}
entries := make([]Entry, 0, len(m))
for k, v := range m {
    entries = append(entries, Entry{Key: k, Value: v})
}
// []Entry{{Key:"foo", Value:1}, {Key:"bar", Value:2}}
```

**After:**
```go
entries := lo.Entries(map[string]int{"foo": 1, "bar": 2})
// []lo.Entry[string, int]{
//     {Key: "foo", Value: 1},
//     {Key: "bar", Value: 2},
// }
```

---

## FromEntries

**Before:**
```go
type Entry struct {
    Key   string
    Value int
}
pairs := []Entry{{Key: "foo", Value: 1}, {Key: "bar", Value: 2}}
m := make(map[string]int, len(pairs))
for _, e := range pairs {
    m[e.Key] = e.Value
}
// map[string]int{"foo": 1, "bar": 2}
```

**After:**
```go
m := lo.FromEntries([]lo.Entry[string, int]{
    {
        Key: "foo",
        Value: 1,
    },
    {
        Key: "bar",
        Value: 2,
    },
})
// map[string]int{"foo": 1, "bar": 2}
```

---

## Invert

**Before:**
```go
m := map[string]int{"a": 1, "b": 2}
inverted := make(map[int]string, len(m))
for k, v := range m {
    inverted[v] = k
}
// map[int]string{1: "a", 2: "b"}
```

**After:**
```go
m1 := lo.Invert(map[string]int{"a": 1, "b": 2})
// map[int]string{1: "a", 2: "b"}

m2 := lo.Invert(map[string]int{"a": 1, "b": 2, "c": 1})
// map[int]string{1: "c", 2: "b"}
```

---

## Assign

**Before:**
```go
maps := []map[string]int{
    {"a": 1, "b": 2},
    {"b": 3, "c": 4},
}
result := make(map[string]int)
for _, m := range maps {
    for k, v := range m {
        result[k] = v
    }
}
// map[string]int{"a": 1, "b": 3, "c": 4}
```

**After:**
```go
mergedMaps := lo.Assign(
    map[string]int{"a": 1, "b": 2},
    map[string]int{"b": 3, "c": 4},
)
// map[string]int{"a": 1, "b": 3, "c": 4}
```

---

## MapKeys

**Before:**
```go
m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
result := make(map[string]int, len(m))
for _, v := range m {
    result[strconv.FormatInt(int64(v), 10)] = v
}
// map[string]int{"1": 1, "2": 2, "3": 3, "4": 4}
```

**After:**
```go
m2 := lo.MapKeys(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(_ int, v int) string {
    return strconv.FormatInt(int64(v), 10)
})
// map[string]int{"1": 1, "2": 2, "3": 3, "4": 4}
```

---

## MapValues

**Before:**
```go
m1 := map[int]int64{1: 1, 2: 2, 3: 3}
m2 := make(map[int]string, len(m1))
for k, x := range m1 {
    m2[k] = strconv.FormatInt(x, 10)
}
// map[int]string{1: "1", 2: "2", 3: "3"}
```

**After:**
```go
m1 := map[int]int64{1: 1, 2: 2, 3: 3}

m2 := lo.MapValues(m1, func(x int64, _ int) string {
    return strconv.FormatInt(x, 10)
})
// map[int]string{1: "1", 2: "2", 3: "3"}
```

---

## MapToSlice

**Before:**
```go
m := map[int]int64{1: 4, 2: 5, 3: 6}
s := make([]string, 0, len(m))
for k, v := range m {
    s = append(s, fmt.Sprintf("%d_%d", k, v))
}
// []string{"1_4", "2_5", "3_6"} (order not guaranteed)
```

**After:**
```go
m := map[int]int64{1: 4, 2: 5, 3: 6}

s := lo.MapToSlice(m, func(k int, v int64) string {
    return fmt.Sprintf("%d_%d", k, v)
})
// []string{"1_4", "2_5", "3_6"}
```

---

## Range / RangeFrom / RangeWithSteps

**Before:**
```go
// Range(4)
result := make([]int, 4)
for i := range result {
    result[i] = i
}
// [0, 1, 2, 3]

// RangeWithSteps(0, 20, 5)
var stepped []int
for i := 0; i < 20; i += 5 {
    stepped = append(stepped, i)
}
// [0, 5, 10, 15]
```

**After:**
```go
result := lo.Range(4)
// [0, 1, 2, 3]

result := lo.Range(-4)
// [0, -1, -2, -3]

result := lo.RangeFrom(1, 5)
// [1, 2, 3, 4, 5]

result := lo.RangeFrom[float64](1.0, 5)
// [1.0, 2.0, 3.0, 4.0, 5.0]

result := lo.RangeWithSteps(0, 20, 5)
// [0, 5, 10, 15]

result := lo.RangeWithSteps[float32](-1.0, -4.0, -1.0)
// [-1.0, -2.0, -3.0]
```

---

## Clamp

**Before:**
```go
func clamp(val, min, max int) int {
    if val < min {
        return min
    }
    if val > max {
        return max
    }
    return val
}

r1 := clamp(0, -10, 10)   // 0
r2 := clamp(-42, -10, 10) // -10
r3 := clamp(42, -10, 10)  // 10
```

**After:**
```go
r1 := lo.Clamp(0, -10, 10)
// 0

r2 := lo.Clamp(-42, -10, 10)
// -10

r3 := lo.Clamp(42, -10, 10)
// 10
```

---

## Sum

**Before:**
```go
list := []int{1, 2, 3, 4, 5}
sum := 0
for _, v := range list {
    sum += v
}
// 15
```

**After:**
```go
list := []int{1, 2, 3, 4, 5}
sum := lo.Sum(list)
// 15
```

---

## SumBy

**Before:**
```go
strs := []string{"foo", "bar"}
sum := 0
for _, item := range strs {
    sum += len(item)
}
// 6
```

**After:**
```go
strings := []string{"foo", "bar"}
sum := lo.SumBy(strings, func(item string) int {
    return len(item)
})
// 6
```

---

## Mean

**Before:**
```go
nums := []int{2, 3, 4, 5}
if len(nums) == 0 {
    // return 0
}
sum := 0
for _, v := range nums {
    sum += v
}
mean := sum / len(nums)
// 3
```

**After:**
```go
mean := lo.Mean([]int{2, 3, 4, 5})
// 3

mean := lo.Mean([]float64{2, 3, 4, 5})
// 3.5

mean := lo.Mean([]float64{})
// 0
```

---

## MeanBy

**Before:**
```go
list := []string{"aa", "bbb", "cccc", "ddddd"}
if len(list) == 0 {
    // return 0
}
sum := 0.0
for _, item := range list {
    sum += float64(len(item))
}
mean := sum / float64(len(list))
// 3.5
```

**After:**
```go
list := []string{"aa", "bbb", "cccc", "ddddd"}
mapper := func(item string) float64 {
    return float64(len(item))
}

mean := lo.MeanBy(list, mapper)
// 3.5
```

---

## Contains

**Before:**
```go
collection := []int{0, 1, 2, 3, 4, 5}
target := 5
present := false
for _, v := range collection {
    if v == target {
        present = true
        break
    }
}
// true
```

**After:**
```go
present := lo.Contains([]int{0, 1, 2, 3, 4, 5}, 5)
// true
```

---

## ContainsBy

**Before:**
```go
collection := []int{0, 1, 2, 3, 4, 5}
present := false
for _, x := range collection {
    if x == 3 {
        present = true
        break
    }
}
// true
```

**After:**
```go
present := lo.ContainsBy([]int{0, 1, 2, 3, 4, 5}, func(x int) bool {
    return x == 3
})
// true
```

---

## Every

**Before:**
```go
collection := []int{0, 1, 2, 3, 4, 5}
subset := []int{0, 2}
lookup := make(map[int]struct{}, len(collection))
for _, v := range collection {
    lookup[v] = struct{}{}
}
ok := true
for _, v := range subset {
    if _, found := lookup[v]; !found {
        ok = false
        break
    }
}
// true
```

**After:**
```go
ok := lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// true

ok := lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// false
```

---

## EveryBy

**Before:**
```go
collection := []int{1, 2, 3, 4}
ok := true
for _, x := range collection {
    if !(x < 5) {
        ok = false
        break
    }
}
// true
```

**After:**
```go
b := lo.EveryBy([]int{1, 2, 3, 4}, func(x int) bool {
    return x < 5
})
// true
```

---

## Some

**Before:**
```go
collection := []int{0, 1, 2, 3, 4, 5}
subset := []int{0, 6}
lookup := make(map[int]struct{}, len(collection))
for _, v := range collection {
    lookup[v] = struct{}{}
}
ok := false
for _, v := range subset {
    if _, found := lookup[v]; found {
        ok = true
        break
    }
}
// true
```

**After:**
```go
ok := lo.Some([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// true
```

---

## SomeBy

**Before:**
```go
collection := []int{1, 2, 3, 4}
found := false
for _, x := range collection {
    if x < 3 {
        found = true
        break
    }
}
// true
```

**After:**
```go
b := lo.SomeBy([]int{1, 2, 3, 4}, func(x int) bool {
    return x < 3
})
// true
```

---

## Intersect

**Before:**
```go
a := []int{0, 1, 2, 3, 4, 5}
b := []int{0, 2}
lookup := make(map[int]struct{}, len(a))
for _, v := range a {
    lookup[v] = struct{}{}
}
var result []int
for _, v := range b {
    if _, ok := lookup[v]; ok {
        result = append(result, v)
    }
}
// []int{0, 2}
```

**After:**
```go
result1 := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// []int{0, 2}

result2 := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// []int{0}

result3 := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// []int{}
```

---

## Difference

**Before:**
```go
list1 := []int{0, 1, 2, 3, 4, 5}
list2 := []int{0, 2, 6}
set1 := make(map[int]struct{}, len(list1))
for _, v := range list1 { set1[v] = struct{}{} }
set2 := make(map[int]struct{}, len(list2))
for _, v := range list2 { set2[v] = struct{}{} }

var left, right []int
for _, v := range list1 {
    if _, ok := set2[v]; !ok {
        left = append(left, v)
    }
}
for _, v := range list2 {
    if _, ok := set1[v]; !ok {
        right = append(right, v)
    }
}
// left: []int{1, 3, 4, 5}, right: []int{6}
```

**After:**
```go
left, right := lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 6})
// []int{1, 3, 4, 5}, []int{6}

left, right := lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5})
// []int{}, []int{}
```

---

## Union

**Before:**
```go
collections := [][]int{{0, 1, 2, 3, 4, 5}, {0, 2}, {0, 10}}
seen := make(map[int]struct{})
var result []int
for _, col := range collections {
    for _, v := range col {
        if _, ok := seen[v]; !ok {
            seen[v] = struct{}{}
            result = append(result, v)
        }
    }
}
// []int{0, 1, 2, 3, 4, 5, 10}
```

**After:**
```go
union := lo.Union([]int{0, 1, 2, 3, 4, 5}, []int{0, 2}, []int{0, 10})
// []int{0, 1, 2, 3, 4, 5, 10}
```

---

## Without

**Before:**
```go
collection := []int{0, 2, 10}
exclude := map[int]struct{}{2: {}}
var result []int
for _, v := range collection {
    if _, ok := exclude[v]; !ok {
        result = append(result, v)
    }
}
// []int{0, 10}
```

**After:**
```go
subset := lo.Without([]int{0, 2, 10}, 2)
// []int{0, 10}

subset := lo.Without([]int{0, 2, 10}, 0, 1, 2, 3, 4, 5)
// []int{10}
```

---

## Find

**Before:**
```go
collection := []string{"a", "b", "c", "d"}
var found string
ok := false
for _, i := range collection {
    if i == "b" {
        found = i
        ok = true
        break
    }
}
// "b", true
```

**After:**
```go
str, ok := lo.Find([]string{"a", "b", "c", "d"}, func(i string) bool {
    return i == "b"
})
// "b", true

str, ok := lo.Find([]string{"foobar"}, func(i string) bool {
    return i == "b"
})
// "", false
```

---

## FindIndexOf

**Before:**
```go
collection := []string{"a", "b", "a", "b"}
var found string
index := -1
ok := false
for i, v := range collection {
    if v == "b" {
        found = v
        index = i
        ok = true
        break
    }
}
// "b", 1, true
```

**After:**
```go
str, index, ok := lo.FindIndexOf([]string{"a", "b", "a", "b"}, func(i string) bool {
    return i == "b"
})
// "b", 1, true

str, index, ok := lo.FindIndexOf([]string{"foobar"}, func(i string) bool {
    return i == "b"
})
// "", -1, false
```

---

## FindOrElse

**Before:**
```go
collection := []string{"a", "b", "c", "d"}
result := "x"
for _, i := range collection {
    if i == "b" {
        result = i
        break
    }
}
// "b"
```

**After:**
```go
str := lo.FindOrElse([]string{"a", "b", "c", "d"}, "x", func(i string) bool {
    return i == "b"
})
// "b"

str := lo.FindOrElse([]string{"foobar"}, "x", func(i string) bool {
    return i == "b"
})
// "x"
```

---

## Min

**Before:**
```go
nums := []int{1, 2, 3}
if len(nums) == 0 {
    // return zero value
}
min := nums[0]
for _, v := range nums[1:] {
    if v < min {
        min = v
    }
}
// 1
```

**After:**
```go
min := lo.Min([]int{1, 2, 3})
// 1

min := lo.Min([]int{})
// 0

min := lo.Min([]time.Duration{time.Second, time.Hour})
// 1s
```

---

## MinBy

**Before:**
```go
strs := []string{"s1", "string2", "s3"}
if len(strs) == 0 {
    // return ""
}
min := strs[0]
for _, item := range strs[1:] {
    if len(item) < len(min) {
        min = item
    }
}
// "s1"
```

**After:**
```go
min := lo.MinBy([]string{"s1", "string2", "s3"}, func(item string, min string) bool {
    return len(item) < len(min)
})
// "s1"

min := lo.MinBy([]string{}, func(item string, min string) bool {
    return len(item) < len(min)
})
// ""
```

---

## Max

**Before:**
```go
nums := []int{1, 2, 3}
if len(nums) == 0 {
    // return zero value
}
max := nums[0]
for _, v := range nums[1:] {
    if v > max {
        max = v
    }
}
// 3
```

**After:**
```go
max := lo.Max([]int{1, 2, 3})
// 3

max := lo.Max([]int{})
// 0

max := lo.Max([]time.Duration{time.Second, time.Hour})
// 1h
```

---

## MaxBy

**Before:**
```go
strs := []string{"string1", "s2", "string3"}
if len(strs) == 0 {
    // return ""
}
max := strs[0]
for _, item := range strs[1:] {
    if len(item) > len(max) {
        max = item
    }
}
// "string1"
```

**After:**
```go
max := lo.MaxBy([]string{"string1", "s2", "string3"}, func(item string, max string) bool {
    return len(item) > len(max)
})
// "string1"

max := lo.MaxBy([]string{}, func(item string, max string) bool {
    return len(item) > len(max)
})
// ""
```

---

## First

**Before:**
```go
collection := []int{1, 2, 3}
if len(collection) == 0 {
    // 0, false
}
first := collection[0]
// 1, true
```

**After:**
```go
first, ok := lo.First([]int{1, 2, 3})
// 1, true

first, ok := lo.First([]int{})
// 0, false
```

---

## Last

**Before:**
```go
collection := []int{1, 2, 3}
if len(collection) == 0 {
    // 0, false
}
last := collection[len(collection)-1]
// 3, true
```

**After:**
```go
last, ok := lo.Last([]int{1, 2, 3})
// 3
// true

last, ok := lo.Last([]int{})
// 0
// false
```

---

## Nth

**Before:**
```go
collection := []int{0, 1, 2, 3}
n := 2
if n < 0 {
    n = len(collection) + n
}
if n < 0 || n >= len(collection) {
    // return 0, error
}
nth := collection[n]
// 2
```

**After:**
```go
nth, err := lo.Nth([]int{0, 1, 2, 3}, 2)
// 2

nth, err := lo.Nth([]int{0, 1, 2, 3}, -2)
// 2
```

---

## Sample

**Before:**
```go
import "math/rand"

collection := []string{"a", "b", "c"}
if len(collection) == 0 {
    // return ""
}
item := collection[rand.Intn(len(collection))]
// a random string from {"a", "b", "c"}
```

**After:**
```go
lo.Sample([]string{"a", "b", "c"})
// a random string from []string{"a", "b", "c"}

lo.Sample([]string{})
// ""
```

---

## Samples

**Before:**
```go
import "math/rand"

collection := []string{"a", "b", "c"}
n := 3
shuffled := make([]string, len(collection))
copy(shuffled, collection)
rand.Shuffle(len(shuffled), func(i, j int) {
    shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
})
result := shuffled[:n]
// []string{"a", "b", "c"} in random order
```

**After:**
```go
lo.Samples([]string{"a", "b", "c"}, 3)
// []string{"a", "b", "c"} in random order
```

---

## Ternary

**Before:**
```go
var result string
if true {
    result = "a"
} else {
    result = "b"
}
// "a"
```

**After:**
```go
result := lo.Ternary(true, "a", "b")
// "a"

result := lo.Ternary(false, "a", "b")
// "b"
```

---

## If / ElseIf / Else

**Before:**
```go
var result int
if true {
    result = 1
} else if false {
    result = 2
} else {
    result = 3
}
// 1
```

**After:**
```go
result := lo.If(true, 1).
    ElseIf(false, 2).
    Else(3)
// 1

result := lo.If(false, 1).
    ElseIf(true, 2).
    Else(3)
// 2

result := lo.If(false, 1).
    ElseIf(false, 2).
    Else(3)
// 3
```

---

## Switch / Case / Default

**Before:**
```go
var result string
switch 1 {
case 1:
    result = "1"
case 2:
    result = "2"
default:
    result = "3"
}
// "1"
```

**After:**
```go
result := lo.Switch(1).
    Case(1, "1").
    Case(2, "2").
    Default("3")
// "1"

result := lo.Switch(2).
    Case(1, "1").
    Case(2, "2").
    Default("3")
// "2"

result := lo.Switch(42).
    Case(1, "1").
    Case(2, "2").
    Default("3")
// "3"
```

---

## ToPtr

**Before:**
```go
s := "hello world"
ptr := &s
// *string{"hello world"}
```

**After:**
```go
ptr := lo.ToPtr("hello world")
// *string{"hello world"}
```

---

## FromPtr

**Before:**
```go
str := "hello world"
ptr := &str

var value string
if ptr != nil {
    value = *ptr
}
// "hello world"

var nilPtr *string
var nilValue string
if nilPtr != nil {
    nilValue = *nilPtr
}
// ""
```

**After:**
```go
str := "hello world"
value := lo.FromPtr(&str)
// "hello world"

value := lo.FromPtr(nil)
// ""
```

---

## FromPtrOr

**Before:**
```go
str := "hello world"
ptr := &str
fallback := "empty"

var value string
if ptr != nil {
    value = *ptr
} else {
    value = fallback
}
// "hello world"
```

**After:**
```go
str := "hello world"
value := lo.FromPtrOr(&str, "empty")
// "hello world"

value := lo.FromPtrOr(nil, "empty")
// "empty"
```

---

## Empty

**Before:**
```go
var zeroInt int     // 0
var zeroStr string  // ""
var zeroBool bool   // false
```

**After:**
```go
lo.Empty[int]()
// 0
lo.Empty[string]()
// ""
lo.Empty[bool]()
// false
```

---

## IsEmpty

**Before:**
```go
val := 0
isEmpty := val == 0
// true

strVal := "foobar"
isStrEmpty := strVal == ""
// false
```

**After:**
```go
lo.IsEmpty(0)
// true
lo.IsEmpty(42)
// false

lo.IsEmpty("")
// true
lo.IsEmpty("foobar")
// false
```

---

## Coalesce

**Before:**
```go
args := []int{0, 1, 2, 3}
var result int
ok := false
for _, v := range args {
    if v != 0 {
        result = v
        ok = true
        break
    }
}
// 1, true
```

**After:**
```go
result, ok := lo.Coalesce(0, 1, 2, 3)
// 1 true

result, ok := lo.Coalesce("")
// "" false

var nilStr *string
str := "foobar"
result, ok := lo.Coalesce(nil, nilStr, &str)
// &"foobar" true
```

---

## Attempt

**Before:**
```go
maxAttempts := 42
var err error
var iter int
for iter = 1; iter <= maxAttempts; iter++ {
    if iter == 5 {
        err = nil
        break
    }
    err = errors.New("failed")
}
// iter=6, err=nil
```

**After:**
```go
iter, err := lo.Attempt(42, func(i int) error {
    if i == 5 {
        return nil
    }

    return errors.New("failed")
})
// 6
// nil

iter, err := lo.Attempt(2, func(i int) error {
    if i == 5 {
        return nil
    }

    return errors.New("failed")
})
// 2
// error "failed"
```

---

## Debounce

**Before:**
```go
import (
    "sync"
    "time"
)

type Debouncer struct {
    mu    sync.Mutex
    timer *time.Timer
    fn    func()
    delay time.Duration
}

func (d *Debouncer) Call() {
    d.mu.Lock()
    defer d.mu.Unlock()
    if d.timer != nil {
        d.timer.Stop()
    }
    d.timer = time.AfterFunc(d.delay, d.fn)
}

func (d *Debouncer) Cancel() {
    d.mu.Lock()
    defer d.mu.Unlock()
    if d.timer != nil {
        d.timer.Stop()
    }
}

f := func() { println("Called once after 100ms when debounce stopped invoking!") }
d := &Debouncer{fn: f, delay: 100 * time.Millisecond}
for j := 0; j < 10; j++ {
    d.Call()
}
time.Sleep(1 * time.Second)
d.Cancel()
```

**After:**
```go
f := func() {
    println("Called once after 100ms when debounce stopped invoking!")
}

debounce, cancel := lo.NewDebounce(100 * time.Millisecond, f)
for j := 0; j < 10; j++ {
    debounce()
}

time.Sleep(1 * time.Second)
cancel()
```

---

## Async

**Before:**
```go
ch := make(chan error, 1)
go func() {
    time.Sleep(10 * time.Second)
    ch <- nil
}()
// chan error
```

**After:**
```go
ch := lo.Async(func() error { time.Sleep(10 * time.Second); return nil })
// chan error (nil)
```

---

## Transaction

**Before:**
```go
state := -5
var finalErr error

// step 1
state += 10
fmt.Println("step 1")

// step 2
state += 15
fmt.Println("step 2")

// step 3
fmt.Println("step 3")
if true {
    finalErr = errors.New("error")
    // rollback step 2
    state -= 15
    fmt.Println("rollback 2")
    // rollback step 1
    state -= 10
    fmt.Println("rollback 1")
}
```

**After:**
```go
transaction := lo.NewTransaction().
    Then(
        func(state int) (int, error) {
            fmt.Println("step 1")
            return state + 10, nil
        },
        func(state int) int {
            fmt.Println("rollback 1")
            return state - 10
        },
    ).
    Then(
        func(state int) (int, error) {
            fmt.Println("step 2")
            return state + 15, nil
        },
        func(state int) int {
            fmt.Println("rollback 2")
            return state - 15
        },
    ).
    Then(
        func(state int) (int, error) {
            fmt.Println("step 3")

            if true {
                return state, errors.New("error")
            }

            return state + 42, nil
        },
        func(state int) int {
            fmt.Println("rollback 3")
            return state - 42
        },
    )

_, _ = transaction.Process(-5)

// Output:
// step 1
// step 2
// step 3
// rollback 2
// rollback 1
```

---

## Must

**Before:**
```go
val, err := time.Parse("2006-01-02", "2022-01-15")
if err != nil {
    panic(err)
}
// val = 2022-01-15
```

**After:**
```go
val := lo.Must(time.Parse("2006-01-02", "2022-01-15"))
// 2022-01-15

val := lo.Must(time.Parse("2006-01-02", "bad-value"))
// panics
```

---

## Try

**Before:**
```go
func safeCall(fn func() error) (ok bool) {
    defer func() {
        if r := recover(); r != nil {
            ok = false
        }
    }()
    err := fn()
    return err == nil
}

ok := safeCall(func() error {
    panic("error")
    return nil
})
// false

ok = safeCall(func() error {
    return nil
})
// true
```

**After:**
```go
ok := lo.Try(func() error {
    panic("error")
    return nil
})
// false

ok := lo.Try(func() error {
    return nil
})
// true

ok := lo.Try(func() error {
    return errors.New("error")
})
// false
```

---

## TryCatch

**Before:**
```go
caught := false

func() {
    defer func() {
        if r := recover(); r != nil {
            caught = true
        }
    }()
    panic("error")
}()
// caught == true
```

**After:**
```go
caught := false

ok := lo.TryCatch(func() error {
    panic("error")
    return nil
}, func() {
    caught = true
})
// false
// caught == true
```

---

## HasKey

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2}
_, exists := m["foo"]
// true

_, exists = m["baz"]
// false
```

**After:**
```go
exists := lo.HasKey(map[string]int{"foo": 1, "bar": 2}, "foo")
// true

exists := lo.HasKey(map[string]int{"foo": 1, "bar": 2}, "baz")
// false
```

---

## ValueOr

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2}
value, ok := m["foo"]
if !ok {
    value = 42
}
// 1

value, ok = m["baz"]
if !ok {
    value = 42
}
// 42
```

**After:**
```go
value := lo.ValueOr(map[string]int{"foo": 1, "bar": 2}, "foo", 42)
// 1

value := lo.ValueOr(map[string]int{"foo": 1, "bar": 2}, "baz", 42)
// 42
```

---

## IndexOf

**Before:**
```go
collection := []int{0, 1, 2, 1, 2, 3}
target := 2
found := -1
for i, v := range collection {
    if v == target {
        found = i
        break
    }
}
// 2
```

**After:**
```go
found := lo.IndexOf([]int{0, 1, 2, 1, 2, 3}, 2)
// 2

notFound := lo.IndexOf([]int{0, 1, 2, 1, 2, 3}, 6)
// -1
```

---

## LastIndexOf

**Before:**
```go
collection := []int{0, 1, 2, 1, 2, 3}
target := 2
found := -1
for i, v := range collection {
    if v == target {
        found = i
    }
}
// 4
```

**After:**
```go
found := lo.LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 2)
// 4

notFound := lo.LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 6)
// -1
```

---

## Validate

**Before:**
```go
slice := []string{"a"}
var err error
if len(slice) != 0 {
    err = fmt.Errorf("Slice should be empty but contains %v", slice)
}
// error("Slice should be empty but contains [a]")
```

**After:**
```go
slice := []string{"a"}
val := lo.Validate(len(slice) == 0, "Slice should be empty but contains %v", slice)
// error("Slice should be empty but contains [a]")

slice := []string{}
val := lo.Validate(len(slice) == 0, "Slice should be empty but contains %v", slice)
// nil
```

---

## ToSlicePtr

**Before:**
```go
strs := []string{"hello", "world"}
ptrs := make([]*string, len(strs))
for i := range strs {
    s := strs[i]
    ptrs[i] = &s
}
// []*string{"hello", "world"}
```

**After:**
```go
ptr := lo.ToSlicePtr([]string{"hello", "world"})
// []*string{"hello", "world"}
```

---

## FromSlicePtr

**Before:**
```go
str1 := "hello"
str2 := "world"
ptrs := []*string{&str1, &str2, nil}
result := make([]string, len(ptrs))
for i, p := range ptrs {
    if p != nil {
        result[i] = *p
    }
}
// []string{"hello", "world", ""}
```

**After:**
```go
str1 := "hello"
str2 := "world"

ptr := lo.FromSlicePtr[string]([]*string{&str1, &str2, nil})
// []string{"hello", "world", ""}

ptr := lo.Compact(
    lo.FromSlicePtr[string]([]*string{&str1, &str2, nil}),
)
// []string{"hello", "world"}
```

---

## ToAnySlice

**Before:**
```go
ints := []int{1, 5, 1}
result := make([]any, len(ints))
for i, v := range ints {
    result[i] = v
}
// []any{1, 5, 1}
```

**After:**
```go
elements := lo.ToAnySlice([]int{1, 5, 1})
// []any{1, 5, 1}
```

---

## FromAnySlice

**Before:**
```go
input := []any{"foobar", "42"}
result := make([]string, 0, len(input))
ok := true
for _, v := range input {
    s, isStr := v.(string)
    if !isStr {
        ok = false
        break
    }
    result = append(result, s)
}
// []string{"foobar", "42"}, true
```

**After:**
```go
elements, ok := lo.FromAnySlice([]any{"foobar", 42})
// []string{}, false

elements, ok := lo.FromAnySlice([]any{"foobar", "42"})
// []string{"foobar", "42"}, true
```

---

## IsNil

**Before:**
```go
import "reflect"

func isNilSafe(v any) bool {
    if v == nil {
        return true
    }
    rv := reflect.ValueOf(v)
    switch rv.Kind() {
    case reflect.Ptr, reflect.Interface, reflect.Chan,
         reflect.Func, reflect.Map, reflect.Slice:
        return rv.IsNil()
    }
    return false
}

var i *int
isNilSafe(i) // true

var ifaceWithNilValue any = (*string)(nil)
isNilSafe(ifaceWithNilValue) // true
ifaceWithNilValue == nil     // false
```

**After:**
```go
var x int
lo.IsNil(x)
// false

var i *int
lo.IsNil(i)
// true

var ifaceWithNilValue any = (*string)(nil)
lo.IsNil(ifaceWithNilValue)
// true
ifaceWithNilValue == nil
// false
```

---

## Partial

**Before:**
```go
add := func(x, y int) int { return x + y }

// manually create a closure with the first arg bound
addFive := func(y int) int {
    return add(5, y)
}

addFive(10) // 15
addFive(42) // 47
```

**After:**
```go
add := func(x, y int) int { return x + y }
f := lo.Partial(add, 5)

f(10)
// 15

f(42)
// 47
```

---

## ErrorsAs

**Before:**
```go
err := doSomething()

var rateLimitErr *RateLimitError
if ok := errors.As(err, &rateLimitErr); ok {
    // retry later
}
```

**After:**
```go
err := doSomething()

if rateLimitErr, ok := lo.ErrorsAs[*RateLimitError](err); ok {
    // retry later
}
```

---

## TryOr

**Before:**
```go
func safeGetString(fn func() (string, error), fallback string) (string, bool) {
    defer func() { recover() }()
    result, err := fn()
    if err != nil {
        return fallback, false
    }
    return result, true
}

str, ok := safeGetString(func() (string, error) {
    return "hello", nil
}, "world")
// "hello", true
```

**After:**
```go
str, ok := lo.TryOr(func() (string, error) {
    panic("error")
    return "hello", nil
}, "world")
// world
// false

str, ok := lo.TryOr(func() (string, error) {
    return "hello", nil
}, "world")
// hello
// true
```

---

## Assert / Assertf

**Before:**
```go
age := getUserAge()
if age < 15 {
    panic("assertion failed: age must be >= 15")
}

// with formatting:
if age < 15 {
    panic(fmt.Sprintf("user age must be >= 15, got %d", age))
}
```

**After:**
```go
age := getUserAge()

lo.Assert(age >= 15)

lo.Assert(age >= 15, "user age must be >= 15")

lo.Assertf(age >= 15, "user age must be >= 15, got %d", age)
```

---

## CountValuesBy

**Before:**
```go
isEven := func(v int) bool { return v%2 == 0 }

input := []int{1, 2, 2}
counts := make(map[bool]int)
for _, v := range input {
    counts[isEven(v)]++
}
// map[bool]int{false: 1, true: 2}
```

**After:**
```go
isEven := func(v int) bool {
    return v%2==0
}

lo.CountValuesBy([]int{1, 2}, isEven)
// map[bool]int{false: 1, true: 1}

lo.CountValuesBy([]int{1, 2, 2}, isEven)
// map[bool]int{false: 1, true: 2}

length := func(v string) int {
    return len(v)
}

lo.CountValuesBy([]string{"foo", "bar", ""}, length)
// map[int]int{0: 1, 3: 2}
```

---

## OmitByKeys

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2, "baz": 3}
exclude := map[string]struct{}{"foo": {}, "baz": {}}
result := make(map[string]int)
for k, v := range m {
    if _, ok := exclude[k]; !ok {
        result[k] = v
    }
}
// map[string]int{"bar": 2}
```

**After:**
```go
m := lo.OmitByKeys(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})
// map[string]int{"bar": 2}
```

---

## FindKey

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2, "baz": 3}
target := 2
var found string
ok := false
for k, v := range m {
    if v == target {
        found = k
        ok = true
        break
    }
}
// "bar", true
```

**After:**
```go
result1, ok1 := lo.FindKey(map[string]int{"foo": 1, "bar": 2, "baz": 3}, 2)
// "bar", true

result2, ok2 := lo.FindKey(map[string]int{"foo": 1, "bar": 2, "baz": 3}, 42)
// "", false
```

---

## Splice

**Before:**
```go
slice := []string{"a", "b"}
insert := []string{"1", "2"}
i := 1
result := make([]string, 0, len(slice)+len(insert))
result = append(result, slice[:i]...)
result = append(result, insert...)
result = append(result, slice[i:]...)
// []string{"a", "1", "2", "b"}
```

**After:**
```go
result := lo.Splice([]string{"a", "b"}, 1, "1", "2")
// []string{"a", "1", "2", "b"}

// negative
result = lo.Splice([]string{"a", "b"}, -1, "1", "2")
// []string{"a", "1", "2", "b"}

// overflow
result = lo.Splice([]string{"a", "b"}, 42, "1", "2")
// []string{"a", "b", "1", "2"}
```

---

## MapEntries

**Before:**
```go
in := map[string]int{"foo": 1, "bar": 2}
out := make(map[int]string, len(in))
for k, v := range in {
    out[v] = k
}
// map[int]string{1: "foo", 2: "bar"}
```

**After:**
```go
in := map[string]int{"foo": 1, "bar": 2}

out := lo.MapEntries(in, func(k string, v int) (int, string) {
    return v,k
})
// map[int]string{1: "foo", 2: "bar"}
```

---

## WithoutEmpty

**Before:**
```go
input := []int{0, 2, 10}
var result []int
for _, v := range input {
    if v != 0 {
        result = append(result, v)
    }
}
// []int{2, 10}
```

**After:**
```go
subset := lo.WithoutEmpty([]int{0, 2, 10})
// []int{2, 10}
```

---

## FindUniques

**Before:**
```go
input := []int{1, 2, 2, 1, 2, 3}
counts := make(map[int]int)
for _, v := range input {
    counts[v]++
}
var result []int
for _, v := range input {
    if counts[v] == 1 {
        result = append(result, v)
        break
    }
}
// []int{3}
```

**After:**
```go
uniqueValues := lo.FindUniques([]int{1, 2, 2, 1, 2, 3})
// []int{3}
```

---

## FindDuplicates

**Before:**
```go
input := []int{1, 2, 2, 1, 2, 3}
counts := make(map[int]int)
for _, v := range input {
    counts[v]++
}
seen := make(map[int]struct{})
var result []int
for _, v := range input {
    if counts[v] > 1 {
        if _, already := seen[v]; !already {
            seen[v] = struct{}{}
            result = append(result, v)
        }
    }
}
// []int{1, 2}
```

**After:**
```go
duplicatedValues := lo.FindDuplicates([]int{1, 2, 2, 1, 2, 3})
// []int{1, 2}
```

---

## Keyify

**Before:**
```go
input := []int{1, 1, 2, 3, 4}
set := make(map[int]struct{}, len(input))
for _, v := range input {
    set[v] = struct{}{}
}
// map[int]struct{}{1:{}, 2:{}, 3:{}, 4:{}}
```

**After:**
```go
set := lo.Keyify([]int{1, 1, 2, 3, 4})
// map[int]struct{}{1:{}, 2:{}, 3:{}, 4:{}}
```

---

## Product

**Before:**
```go
list := []int{1, 2, 3, 4, 5}
product := 1
for _, v := range list {
    product *= v
}
// 120
```

**After:**
```go
list := []int{1, 2, 3, 4, 5}
product := lo.Product(list)
// 120
```

---

## Ellipsis

**Before:**
```go
import "unicode/utf8"

func ellipsis(s string, maxRunes int) string {
    s = strings.TrimSpace(s)
    runes := []rune(s)
    if len(runes) <= maxRunes {
        return s
    }
    if maxRunes <= 3 {
        return "..."
    }
    return string(runes[:maxRunes-3]) + "..."
}

str := ellipsis("  Lorem Ipsum  ", 5)
// "Lo..."
```

**After:**
```go
str := lo.Ellipsis("  Lorem Ipsum  ", 5)
// Lo...

str := lo.Ellipsis("Lorem Ipsum", 100)
// Lorem Ipsum

str := lo.Ellipsis("Lorem Ipsum", 3)
// ...
```

---

## PascalCase / CamelCase / KebabCase / SnakeCase

**Before:**
```go
// manual string case conversion requires regex and custom logic
import (
    "strings"
    "unicode"
)

// PascalCase: "hello_world" -> "HelloWorld"
// ... complex manual implementation with multiple edge cases

// CamelCase: "hello_world" -> "helloWorld"
// ... split by separator, capitalize each word except first

// KebabCase: "helloWorld" -> "hello-world"
// ... scan runes, insert hyphens before uppercase letters

// SnakeCase: "HelloWorld" -> "hello_world"
// ... scan runes, insert underscores before uppercase letters
```

**After:**
```go
str := lo.PascalCase("hello_world")
// HelloWorld

str := lo.CamelCase("hello_world")
// helloWorld

str := lo.KebabCase("helloWorld")
// hello-world

str := lo.SnakeCase("HelloWorld")
// hello_world
```

---

## TernaryF

**Before:**
```go
var s *string

var someStr string
if s == nil {
    someStr = uuid.New().String()
} else {
    someStr = *s
}
```

**After:**
```go
var s *string

someStr := lo.TernaryF(s == nil, func() string { return uuid.New().String() }, func() string { return *s })
// ef782193-c30c-4e2e-a7ae-f8ab5e125e02
```

---

## AttemptWithDelay

**Before:**
```go
maxAttempts := 5
delay := 2 * time.Second
var err error
var iter int
for iter = 1; iter <= maxAttempts; iter++ {
    if iter == 3 {
        err = nil
        break
    }
    err = errors.New("failed")
    if iter < maxAttempts {
        time.Sleep(delay)
    }
}
// iter=3, ~4s, nil
```

**After:**
```go
iter, duration, err := lo.AttemptWithDelay(5, 2*time.Second, func(i int, duration time.Duration) error {
    if i == 2 {
        return nil
    }

    return errors.New("failed")
})
// 3
// ~ 4 seconds
// nil
```

---

## Synchronize

**Before:**
```go
import "sync"

mu := sync.Mutex{}

for i := 0; i < 10; i++ {
    go func() {
        mu.Lock()
        defer mu.Unlock()
        println("will be called sequentially")
    }()
}
```

**After:**
```go
s := lo.Synchronize()

for i := 0; i < 10; i++ {
    go s.Do(func () {
        println("will be called sequentially")
    })
}
```

---

## WaitFor

**Before:**
```go
start := time.Now()
timeout := 10 * time.Millisecond
tick := 2 * time.Millisecond
condition := func(i int) bool { return true }

var iterations int
ok := false
for time.Since(start) < timeout {
    iterations++
    if condition(iterations) {
        ok = true
        break
    }
    time.Sleep(tick)
}
duration := time.Since(start)
```

**After:**
```go
alwaysTrue := func(i int) bool { return true }
alwaysFalse := func(i int) bool { return false }

iterations, duration, ok := lo.WaitFor(alwaysTrue, 10*time.Millisecond, 2 * time.Millisecond)
// 1
// 1ms
// true

iterations, duration, ok := lo.WaitFor(alwaysFalse, 10*time.Millisecond, time.Millisecond)
// 10
// 10ms
// false
```

---

## TryCatchWithErrorValue

**Before:**
```go
caught := false

func() {
    defer func() {
        if r := recover(); r != nil {
            if r == "error" {
                caught = true
            }
        }
    }()
    panic("error")
}()
// caught == true
```

**After:**
```go
caught := false

ok := lo.TryCatchWithErrorValue(func() error {
    panic("error")
    return nil
}, func(val any) {
    caught = val == "error"
})
// false
// caught == true
```

---

## UniqMap

**Before:**
```go
type User struct {
    Name string
    Age  int
}
users := []User{{Name: "Alex", Age: 10}, {Name: "Alex", Age: 12}, {Name: "Bob", Age: 11}, {Name: "Alice", Age: 20}}

seen := map[string]struct{}{}
var names []string
for _, u := range users {
    if _, ok := seen[u.Name]; !ok {
        seen[u.Name] = struct{}{}
        names = append(names, u.Name)
    }
}
// []string{"Alex", "Bob", "Alice"}
```

**After:**
```go
type User struct {
    Name string
    Age  int
}
users := []User{{Name: "Alex", Age: 10}, {Name: "Alex", Age: 12}, {Name: "Bob", Age: 11}, {Name: "Alice", Age: 20}}

names := lo.UniqMap(users, func(u User, index int) string {
    return u.Name
})
// []string{"Alex", "Bob", "Alice"}
```

---

## GroupByMap

**Before:**
```go
result := map[int][]int{}
for _, i := range []int{0, 1, 2, 3, 4, 5} {
    key := i % 3
    result[key] = append(result[key], i*2)
}
// map[int][]int{0: []int{0, 6}, 1: []int{2, 8}, 2: []int{4, 10}}
```

**After:**
```go
groups := lo.GroupByMap([]int{0, 1, 2, 3, 4, 5}, func(i int) (int, int) {
    return i%3, i*2
})
// map[int][]int{0: []int{0, 6}, 1: []int{2, 8}, 2: []int{4, 10}}
```

---

## Window

**Before:**
```go
s := []int{1, 2, 3, 4, 5}
size := 3
var windows [][]int
for i := 0; i+size <= len(s); i++ {
    window := make([]int, size)
    copy(window, s[i:i+size])
    windows = append(windows, window)
}
// [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}
```

**After:**
```go
lo.Window([]int{1, 2, 3, 4, 5}, 3)
// [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}

lo.Window([]float64{20, 22, 21, 23, 24}, 3)
// [][]float64{{20, 22, 21}, {22, 21, 23}, {21, 23, 24}}
```

---

## Sliding

**Before:**
```go
s := []int{1, 2, 3, 4, 5, 6}
size, step := 3, 1
var windows [][]int
for i := 0; i+size <= len(s); i += step {
    window := make([]int, size)
    copy(window, s[i:i+size])
    windows = append(windows, window)
}
// [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}}
```

**After:**
```go
// Windows with shared elements (step < size)
lo.Sliding([]int{1, 2, 3, 4, 5, 6}, 3, 1)
// [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}}

// Windows with no shared elements (step == size, like Chunk)
lo.Sliding([]int{1, 2, 3, 4, 5, 6}, 3, 3)
// [][]int{{1, 2, 3}, {4, 5, 6}}

// Step > size (skipping elements)
lo.Sliding([]int{1, 2, 3, 4, 5, 6, 7, 8}, 2, 3)
// [][]int{{1, 2}, {4, 5}, {7, 8}}
```

---

## FilterSliceToMap

**Before:**
```go
list := []string{"a", "aa", "aaa"}
result := map[string]int{}
for _, str := range list {
    if len(str) > 1 {
        result[str] = len(str)
    }
}
// map[string][int]{"aa":2 "aaa":3}
```

**After:**
```go
list := []string{"a", "aa", "aaa"}

result := lo.FilterSliceToMap(list, func(str string) (string, int, bool) {
    return str, len(str), len(str) > 1
})
// map[string][int]{"aa":2 "aaa":3}
```

---

## TakeFilter

**Before:**
```go
s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
n := 3
var result []int
for _, v := range s {
    if v%2 == 0 {
        result = append(result, v)
        if len(result) == n {
            break
        }
    }
}
// []int{2, 4, 6}
```

**After:**
```go
l := lo.TakeFilter([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3, func(val int, index int) bool {
    return val%2 == 0
})
// []int{2, 4, 6}

l := lo.TakeFilter([]string{"a", "aa", "aaa", "aaaa"}, 2, func(val string, index int) bool {
    return len(val) > 1
})
// []string{"aa", "aaa"}
```

---

## DropRightWhile

**Before:**
```go
s := []string{"a", "aa", "aaa", "aa", "aa"}
end := len(s)
for end > 0 && len(s[end-1]) <= 2 {
    end--
}
result := s[:end]
// []string{"a", "aa", "aaa"}
```

**After:**
```go
l := lo.DropRightWhile([]string{"a", "aa", "aaa", "aa", "aa"}, func(val string) bool {
    return len(val) <= 2
})
// []string{"a", "aa", "aaa"}
```

---

## DropByIndex

**Before:**
```go
s := []int{0, 1, 2, 3, 4, 5}
indices := map[int]struct{}{2: {}, 4: {}, 5: {}} // 5 == -1 (last)
var result []int
for i, v := range s {
    if _, skip := indices[i]; !skip {
        result = append(result, v)
    }
}
// []int{0, 1, 3}
```

**After:**
```go
l := lo.DropByIndex([]int{0, 1, 2, 3, 4, 5}, 2, 4, -1)
// []int{0, 1, 3}
```

---

## RejectMap

**Before:**
```go
var result []int
for _, x := range []int{1, 2, 3, 4} {
    if x%2 != 0 {
        result = append(result, x*10)
    }
}
// []int{10, 30}
```

**After:**
```go
items := lo.RejectMap([]int{1, 2, 3, 4}, func(x int, _ int) (int, bool) {
    return x*10, x%2 == 0
})
// []int{10, 30}
```

---

## FilterReject

**Before:**
```go
var kept, rejected []int
for _, x := range []int{1, 2, 3, 4} {
    if x%2 == 0 {
        kept = append(kept, x)
    } else {
        rejected = append(rejected, x)
    }
}
// kept: []int{2, 4}
// rejected: []int{1, 3}
```

**After:**
```go
kept, rejected := lo.FilterReject([]int{1, 2, 3, 4}, func(x int, _ int) bool {
    return x%2 == 0
})
// []int{2, 4}
// []int{1, 3}
```

---

## Clone

**Before:**
```go
in := []int{1, 2, 3, 4, 5}
cloned := make([]int, len(in))
copy(cloned, in)
in[0] = 99
// cloned is []int{1, 2, 3, 4, 5}
```

**After:**
```go
in := []int{1, 2, 3, 4, 5}
cloned := lo.Clone(in)
// Verify it's a different slice by checking that modifying one doesn't affect the other
in[0] = 99
// cloned is []int{1, 2, 3, 4, 5}
```

---

## Cut

**Before:**
```go
s := []string{"a", "b", "c", "d", "e", "f", "g"}
sep := []string{"b", "c", "d"}
// find separator
left, right := s[:0], s[0:]
found := false
for i := 0; i <= len(s)-len(sep); i++ {
    match := true
    for j, v := range sep {
        if s[i+j] != v {
            match = false
            break
        }
    }
    if match {
        left = s[:i]
        right = s[i+len(sep):]
        found = true
        break
    }
}
// left: []string{"a"}, right: []string{"e", "f", "g"}, found: true
```

**After:**
```go
actualLeft, actualRight, result = lo.Cut([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"b", "c", "d"})
// actualLeft: []string{"a"}
// actualRight: []string{"e", "f", "g"}
// result: true

result = lo.Cut([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"z"})
// actualLeft: []string{"a", "b", "c", "d", "e", "f", "g"}
// actualRight: []string{}
// result: false
```

---

## CutPrefix

**Before:**
```go
s := []string{"a", "b", "c", "d", "e", "f", "g"}
prefix := []string{"a", "b", "c"}
found := false
result := s
if len(s) >= len(prefix) {
    match := true
    for i, v := range prefix {
        if s[i] != v {
            match = false
            break
        }
    }
    if match {
        result = s[len(prefix):]
        found = true
    }
}
// result: []string{"d", "e", "f", "g"}, found: true
```

**After:**
```go
actualRight, result = lo.CutPrefix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"a", "b", "c"})
// actualRight: []string{"d", "e", "f", "g"}
// result: true

result = lo.CutPrefix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"b"})
// actualRight: []string{"a", "b", "c", "d", "e", "f", "g"}
// result: false
```

---

## CutSuffix

**Before:**
```go
s := []string{"a", "b", "c", "d", "e", "f", "g"}
suffix := []string{"f", "g"}
result := s
found := false
if len(s) >= len(suffix) {
    start := len(s) - len(suffix)
    match := true
    for i, v := range suffix {
        if s[start+i] != v {
            match = false
            break
        }
    }
    if match {
        result = s[:start]
        found = true
    }
}
// result: []string{"a", "b", "c", "d", "e"}, found: true
```

**After:**
```go
actualLeft, result = lo.CutSuffix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"f", "g"})
// actualLeft: []string{"a", "b", "c", "d", "e"}
// result: true

actualLeft, result = lo.CutSuffix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"b"})
// actualLeft: []string{"a", "b", "c", "d", "e", "f", "g"}
// result: false
```

---

## Trim

**Before:**
```go
s := []int{0, 1, 2, 0, 3, 0}
cutset := map[int]struct{}{1: {}, 0: {}}
start := 0
for start < len(s) {
    if _, ok := cutset[s[start]]; ok {
        start++
    } else {
        break
    }
}
end := len(s)
for end > start {
    if _, ok := cutset[s[end-1]]; ok {
        end--
    } else {
        break
    }
}
result := s[start:end]
// []int{2, 0, 3}
```

**After:**
```go
result := lo.Trim([]int{0, 1, 2, 0, 3, 0}, []int{1, 0})
// []int{2, 0, 3}

result := lo.Trim([]string{"hello", "world", " "}, []string{" ", ""})
// []string{"hello", "world"}
```

---

## TrimLeft

**Before:**
```go
s := []int{0, 1, 2, 0, 3, 0}
cutset := map[int]struct{}{1: {}, 0: {}}
start := 0
for start < len(s) {
    if _, ok := cutset[s[start]]; ok {
        start++
    } else {
        break
    }
}
result := s[start:]
// []int{2, 0, 3, 0}
```

**After:**
```go
result := lo.TrimLeft([]int{0, 1, 2, 0, 3, 0}, []int{1, 0})
// []int{2, 0, 3, 0}

result := lo.TrimLeft([]string{"hello", "world", " "}, []string{" ", ""})
// []string{"hello", "world", " "}
```

---

## TrimPrefix

**Before:**
```go
s := []int{1, 2, 1, 2, 3, 1, 2, 4}
prefix := []int{1, 2}
result := s
for len(result) >= len(prefix) {
    match := true
    for i, v := range prefix {
        if result[i] != v {
            match = false
            break
        }
    }
    if match {
        result = result[len(prefix):]
    } else {
        break
    }
}
// []int{3, 1, 2, 4}
```

**After:**
```go
result := lo.TrimPrefix([]int{1, 2, 1, 2, 3, 1, 2, 4}, []int{1, 2})
// []int{3, 1, 2, 4}

result := lo.TrimPrefix([]string{"hello", "world", "hello", "test"}, []string{"hello"})
// []string{"world", "hello", "test"}
```

---

## TrimRight

**Before:**
```go
s := []int{0, 1, 2, 0, 3, 0}
cutset := map[int]struct{}{0: {}, 3: {}}
end := len(s)
for end > 0 {
    if _, ok := cutset[s[end-1]]; ok {
        end--
    } else {
        break
    }
}
result := s[:end]
// []int{0, 1, 2}
```

**After:**
```go
result := lo.TrimRight([]int{0, 1, 2, 0, 3, 0}, []int{0, 3})
// []int{0, 1, 2}

result := lo.TrimRight([]string{"hello", "world", "  "}, []string{" ", ""})
// []string{"hello", "world", ""}
```

---

## TrimSuffix

**Before:**
```go
s := []int{1, 2, 3, 1, 2, 4, 2, 4, 2, 4}
suffix := []int{2, 4}
result := s
for len(result) >= len(suffix) {
    start := len(result) - len(suffix)
    match := true
    for i, v := range suffix {
        if result[start+i] != v {
            match = false
            break
        }
    }
    if match {
        result = result[:start]
    } else {
        break
    }
}
// []int{1, 2, 3, 1}
```

**After:**
```go
result := lo.TrimSuffix([]int{1, 2, 3, 1, 2, 4, 2, 4, 2, 4}, []int{2, 4})
// []int{1, 2, 3, 1}

result := lo.TrimSuffix([]string{"hello", "world", "hello", "test"}, []string{"test"})
// []string{"hello", "world", "hello"}
```

---

## UniqKeys

**Before:**
```go
m1 := map[string]int{"foo": 1, "bar": 2}
m2 := map[string]int{"bar": 3}
seen := map[string]struct{}{}
var keys []string
for k := range m1 {
    if _, ok := seen[k]; !ok {
        seen[k] = struct{}{}
        keys = append(keys, k)
    }
}
for k := range m2 {
    if _, ok := seen[k]; !ok {
        seen[k] = struct{}{}
        keys = append(keys, k)
    }
}
// []string{"foo", "bar"}
```

**After:**
```go
keys := lo.UniqKeys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
// []string{"foo", "bar", "baz"}

keys := lo.UniqKeys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"bar": 3})
// []string{"foo", "bar"}
```

---

## UniqValues

**Before:**
```go
m1 := map[string]int{"foo": 1, "bar": 2}
m2 := map[string]int{"bar": 2}
seen := map[int]struct{}{}
var values []int
for _, v := range m1 {
    if _, ok := seen[v]; !ok {
        seen[v] = struct{}{}
        values = append(values, v)
    }
}
for _, v := range m2 {
    if _, ok := seen[v]; !ok {
        seen[v] = struct{}{}
        values = append(values, v)
    }
}
// []int{1, 2}
```

**After:**
```go
values := lo.UniqValues(map[string]int{"foo": 1, "bar": 2})
// []int{1, 2}

values := lo.UniqValues(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
// []int{1, 2, 3}

values := lo.UniqValues(map[string]int{"foo": 1, "bar": 2}, map[string]int{"bar": 2})
// []int{1, 2}
```

---

## PickByValues

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2, "baz": 3}
allowed := map[int]struct{}{1: {}, 3: {}}
result := map[string]int{}
for k, v := range m {
    if _, ok := allowed[v]; ok {
        result[k] = v
    }
}
// map[string]int{"foo": 1, "baz": 3}
```

**After:**
```go
m := lo.PickByValues(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})
// map[string]int{"foo": 1, "baz": 3}
```

---

## OmitByValues

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2, "baz": 3}
excluded := map[int]struct{}{1: {}, 3: {}}
result := map[string]int{}
for k, v := range m {
    if _, ok := excluded[v]; !ok {
        result[k] = v
    }
}
// map[string]int{"bar": 2}
```

**After:**
```go
m := lo.OmitByValues(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})
// map[string]int{"bar": 2}
```

---

## ChunkEntries

**Before:**
```go
m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
size := 3
entries := make([]map[string]int, 0)
chunk := map[string]int{}
i := 0
for k, v := range m {
    chunk[k] = v
    i++
    if i == size {
        entries = append(entries, chunk)
        chunk = map[string]int{}
        i = 0
    }
}
if len(chunk) > 0 {
    entries = append(entries, chunk)
}
// []map[string]int{{"a": 1, "b": 2, "c": 3}, {"d": 4, "e": 5}}
```

**After:**
```go
maps := lo.ChunkEntries(
    map[string]int{
        "a": 1,
        "b": 2,
        "c": 3,
        "d": 4,
        "e": 5,
    },
    3,
)
// []map[string]int{
//    {"a": 1, "b": 2, "c": 3},
//    {"d": 4, "e": 5},
// }
```

---

## FilterMapToSlice

**Before:**
```go
kv := map[int]int64{1: 1, 2: 2, 3: 3, 4: 4}
var result []string
for k, v := range kv {
    if k%2 == 0 {
        result = append(result, fmt.Sprintf("%d_%d", k, v))
    }
}
// []{"2_2", "4_4"}
```

**After:**
```go
kv := map[int]int64{1: 1, 2: 2, 3: 3, 4: 4}

result := lo.FilterMapToSlice(kv, func(k int, v int64) (string, bool) {
    return fmt.Sprintf("%d_%d", k, v), k%2 == 0
})
// []{"2_2", "4_4"}
```

---

## FilterKeys

**Before:**
```go
kv := map[int]string{1: "foo", 2: "bar", 3: "baz"}
var result []int
for k, v := range kv {
    if v == "foo" {
        result = append(result, k)
    }
}
// [1]
```

**After:**
```go
kv := map[int]string{1: "foo", 2: "bar", 3: "baz"}

result := FilterKeys(kv, func(k int, v string) bool {
    return v == "foo"
})
// [1]
```

---

## FilterValues

**Before:**
```go
kv := map[int]string{1: "foo", 2: "bar", 3: "baz"}
var result []string
for k, v := range kv {
    if v == "foo" {
        _ = k
        result = append(result, v)
    }
}
// ["foo"]
```

**After:**
```go
kv := map[int]string{1: "foo", 2: "bar", 3: "baz"}

result := FilterValues(kv, func(k int, v string) bool {
    return v == "foo"
})
// ["foo"]
```

---

## ProductBy

**Before:**
```go
strings := []string{"foo", "bar"}
product := 1
for _, item := range strings {
    product *= len(item)
}
// 9
```

**After:**
```go
strings := []string{"foo", "bar"}
product := lo.ProductBy(strings, func(item string) int {
    return len(item)
})
// 9
```

---

## Mode

**Before:**
```go
nums := []int{2, 2, 3, 4}
freq := map[int]int{}
maxFreq := 0
for _, v := range nums {
    freq[v]++
    if freq[v] > maxFreq {
        maxFreq = freq[v]
    }
}
var modes []int
for v, f := range freq {
    if f == maxFreq {
        modes = append(modes, v)
    }
}
// [2]
```

**After:**
```go
mode := lo.Mode([]int{2, 2, 3, 4})
// [2]

mode := lo.Mode([]float64{2, 2, 3, 3})
// [2, 3]

mode := lo.Mode([]float64{})
// []

mode := lo.Mode([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
// [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

---

## RandomString

**Before:**
```go
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
b := make([]byte, 5)
for i := range b {
    b[i] = charset[rand.Intn(len(charset))]
}
str := string(b)
// example: "eIGbt"
```

**After:**
```go
str := lo.RandomString(5, lo.LettersCharset)
// example: "eIGbt"
```

---

## Substring

**Before:**
```go
s := "hello"
runes := []rune(s)
start := 2
length := 3
if start < 0 {
    start = len(runes) + start
}
sub := string(runes[start : start+length])
// "llo"
```

**After:**
```go
sub := lo.Substring("hello", 2, 3)
// "llo"

sub := lo.Substring("hello", -4, 3)
// "ell"

sub := lo.Substring("hello", -2, math.MaxUint)
// "lo"
```

---

## ChunkString

**Before:**
```go
s := "123456"
size := 2
var chunks []string
for len(s) > 0 {
    if len(s) < size {
        chunks = append(chunks, s)
        break
    }
    chunks = append(chunks, s[:size])
    s = s[size:]
}
// []string{"12", "34", "56"}
```

**After:**
```go
lo.ChunkString("123456", 2)
// []string{"12", "34", "56"}

lo.ChunkString("1234567", 2)
// []string{"12", "34", "56", "7"}

lo.ChunkString("", 2)
// []string{""}

lo.ChunkString("1", 2)
// []string{"1"}
```

---

## RuneLength

**Before:**
```go
s := "hellô"
runeCount := 0
for range s {
    runeCount++
}
// 5 (vs len("hellô") == 6)
```

**After:**
```go
sub := lo.RuneLength("hellô")
// 5

sub := len("hellô")
// 6
```

---

## Words

**Before:**
```go
// splitting camelCase/PascalCase manually requires regex or manual parsing
import "regexp"
re := regexp.MustCompile(`[A-Z][a-z]*|[a-z]+`)
words := re.FindAllString("helloWorld", -1)
for i, w := range words {
    words[i] = strings.ToLower(w)
}
// []string{"hello", "world"}
```

**After:**
```go
str := lo.Words("helloWorld")
// []string{"hello", "world"}
```

---

## Capitalize

**Before:**
```go
s := "heLLO"
if len(s) == 0 {
    // ""
}
result := strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
// "Hello"
```

**After:**
```go
str := lo.Capitalize("heLLO")
// Hello
```

---

## T2 -> T9

**Before:**
```go
// No built-in tuple type in Go; typically use a struct
type Pair struct {
    A string
    B int
}
tuple1 := Pair{A: "x", B: 1}
```

**After:**
```go
tuple1 := lo.T2("x", 1)
// Tuple2[string, int]{A: "x", B: 1}

func example() (string, int) { return "y", 2 }
tuple2 := lo.T2(example())
// Tuple2[string, int]{A: "y", B: 2}
```

---

## Unpack2 -> Unpack9

**Before:**
```go
// Manually destructure a struct
type Pair struct{ A string; B int }
p := Pair{A: "a", B: 1}
r1, r2 := p.A, p.B
// "a", 1
```

**After:**
```go
r1, r2 := lo.Unpack2(lo.Tuple2[string, int]{"a", 1})
// "a", 1
```

---

## Zip2 -> Zip9

**Before:**
```go
as := []string{"a", "b"}
bs := []int{1, 2}
type Pair struct{ A string; B int }
tuples := make([]Pair, len(as))
for i := range as {
    tuples[i] = Pair{A: as[i], B: bs[i]}
}
```

**After:**
```go
tuples := lo.Zip2([]string{"a", "b"}, []int{1, 2})
// []Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}}
```

---

## ZipBy2 -> ZipBy9

**Before:**
```go
as := []string{"a", "b"}
bs := []int{1, 2}
items := make([]string, len(as))
for i := range as {
    items[i] = fmt.Sprintf("%s-%d", as[i], bs[i])
}
// []string{"a-1", "b-2"}
```

**After:**
```go
items := lo.ZipBy2([]string{"a", "b"}, []int{1, 2}, func(a string, b int) string {
    return fmt.Sprintf("%s-%d", a, b)
})
// []string{"a-1", "b-2"}
```

---

## Unzip2 -> Unzip9

**Before:**
```go
type Pair struct{ A string; B int }
pairs := []Pair{{A: "a", B: 1}, {A: "b", B: 2}}
as := make([]string, len(pairs))
bs := make([]int, len(pairs))
for i, p := range pairs {
    as[i] = p.A
    bs[i] = p.B
}
// []string{"a", "b"}, []int{1, 2}
```

**After:**
```go
a, b := lo.Unzip2([]Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}})
// []string{"a", "b"}
// []int{1, 2}
```

---

## UnzipBy2 -> UnzipBy9

**Before:**
```go
strs := []string{"hello", "john", "doe"}
names := make([]string, len(strs))
lengths := make([]int, len(strs))
for i, s := range strs {
    names[i] = s
    lengths[i] = len(s)
}
// []string{"hello", "john", "doe"}, []int{5, 4, 3}
```

**After:**
```go
a, b := lo.UnzipBy2([]string{"hello", "john", "doe"}, func(str string) (string, int) {
    return str, len(str)
})
// []string{"hello", "john", "doe"}
// []int{5, 4, 3}
```

---

## CrossJoin2 -> CrossJoin9

**Before:**
```go
type Pair struct{ A string; B int }
as := []string{"hello", "john", "doe"}
bs := []int{1, 2}
var result []Pair
for _, a := range as {
    for _, b := range bs {
        result = append(result, Pair{A: a, B: b})
    }
}
```

**After:**
```go
result := lo.CrossJoin2([]string{"hello", "john", "doe"}, []int{1, 2})
// lo.Tuple2{"hello", 1}
// lo.Tuple2{"hello", 2}
// lo.Tuple2{"john", 1}
// lo.Tuple2{"john", 2}
// lo.Tuple2{"doe", 1}
// lo.Tuple2{"doe", 2}
```

---

## CrossJoinBy2 -> CrossJoinBy9

**Before:**
```go
as := []string{"hello", "john", "doe"}
bs := []int{1, 2}
var result []string
for _, a := range as {
    for _, b := range bs {
        result = append(result, fmt.Sprintf("%s - %d", a, b))
    }
}
```

**After:**
```go
result := lo.CrossJoinBy2([]string{"hello", "john", "doe"}, []int{1, 2}, func(a A, b B) string {
    return fmt.Sprintf("%s - %d", a, b)
})
// "hello - 1"
// "hello - 2"
// "john - 1"
// "john - 2"
// "doe - 1"
// "doe - 2"
```

---

## Duration

**Before:**
```go
start := time.Now()
// very long job
duration := time.Since(start)
// 3s
```

**After:**
```go
duration := lo.Duration(func() {
    // very long job
})
// 3s
```

---

## Duration0 -> Duration10

**Before:**
```go
start := time.Now()
// very long job
err := errors.New("an error")
duration := time.Since(start)
// an error, 3s
```

**After:**
```go
duration := lo.Duration0(func() {
    // very long job
})
// 3s

err, duration := lo.Duration1(func() error {
    // very long job
    return errors.New("an error")
})
// an error
// 3s

str, nbr, err, duration := lo.Duration3(func() (string, int, error) {
    // very long job
    return "hello", 42, nil
})
// hello
// 42
// nil
// 3s
```

---

## ChannelDispatcher

**Before:**
```go
ch := make(chan int, 42)
for i := 0; i <= 10; i++ {
    ch <- i
}
// Manually create N child channels and distribute messages with a goroutine
children := make([]chan int, 5)
for i := range children {
    children[i] = make(chan int, 10)
}
go func() {
    i := 0
    for msg := range ch {
        children[i%len(children)] <- msg
        i++
    }
    for _, c := range children {
        close(c)
    }
}()
```

**After:**
```go
ch := make(chan int, 42)
for i := 0; i <= 10; i++ {
    ch <- i
}

children := lo.ChannelDispatcher(ch, 5, 10, DispatchingStrategyRoundRobin[int])
// []<-chan int{...}

consumer := func(c <-chan int) {
    for {
        msg, ok := <-c
        if !ok {
            println("closed")
            break
        }
        println(msg)
    }
}

for i := range children {
    go consumer(children[i])
}
```

---

## SliceToChannel

**Before:**
```go
list := []int{1, 2, 3, 4, 5}
ch := make(chan int, 2)
go func() {
    for _, v := range list {
        ch <- v
    }
    close(ch)
}()
for v := range ch {
    println(v)
}
```

**After:**
```go
list := []int{1, 2, 3, 4, 5}

for v := range lo.SliceToChannel(2, list) {
    println(v)
}
// prints 1, then 2, then 3, then 4, then 5
```

---

## ChannelToSlice

**Before:**
```go
list := []int{1, 2, 3, 4, 5}
ch := make(chan int)
go func() {
    for _, v := range list {
        ch <- v
    }
    close(ch)
}()
var items []int
for v := range ch {
    items = append(items, v)
}
// []int{1, 2, 3, 4, 5}
```

**After:**
```go
list := []int{1, 2, 3, 4, 5}
ch := lo.SliceToChannel(2, list)

items := ChannelToSlice(ch)
// []int{1, 2, 3, 4, 5}
```

---

## Generator

**Before:**
```go
ch := make(chan int, 2)
go func() {
    ch <- 1
    ch <- 2
    ch <- 3
    close(ch)
}()
for v := range ch {
    println(v)
}
// prints 1, then 2, then 3
```

**After:**
```go
generator := func(yield func(int)) {
    yield(1)
    yield(2)
    yield(3)
}

for v := range lo.Generator(2, generator) {
    println(v)
}
// prints 1, then 2, then 3
```

---

## Buffer

**Before:**
```go
ch := make(chan int)
// ... fill channel ...
items := make([]int, 0, 3)
for len(items) < 3 {
    v, ok := <-ch
    if !ok {
        break
    }
    items = append(items, v)
}
// read up to 3 items, track open/closed manually
```

**After:**
```go
ch := lo.SliceToChannel(2, []int{1, 2, 3, 4, 5})

items1, length1, duration1, ok1 := lo.Buffer(ch, 3)
// []int{1, 2, 3}, 3, 0s, true
items2, length2, duration2, ok2 := lo.Buffer(ch, 3)
// []int{4, 5}, 2, 0s, false
```

---

## BufferWithContext

**Before:**
```go
ctx, cancel := context.WithCancel(context.TODO())
// manually select on ctx.Done() and channel reads to batch with context cancellation
items := make([]int, 0, 3)
for len(items) < 3 {
    select {
    case <-ctx.Done():
        goto done
    case v, ok := <-ch:
        if !ok {
            goto done
        }
        items = append(items, v)
    }
}
done:
```

**After:**
```go
ctx, cancel := context.WithCancel(context.TODO())
go func() {
    ch <- 0
    time.Sleep(10*time.Millisecond)
    ch <- 1
    time.Sleep(10*time.Millisecond)
    ch <- 2
    time.Sleep(10*time.Millisecond)
    ch <- 3
    time.Sleep(10*time.Millisecond)
    ch <- 4
    time.Sleep(10*time.Millisecond)
    cancel()
}()

items1, length1, duration1, ok1 := lo.BufferWithContext(ctx, ch, 3)
// []int{0, 1, 2}, 3, 20ms, true
items2, length2, duration2, ok2 := lo.BufferWithContext(ctx, ch, 3)
// []int{3, 4}, 2, 30ms, false
```

---

## BufferWithTimeout

**Before:**
```go
// Manually implement timeout batching using time.After and channel reads
items := make([]int, 0, 3)
timeout := time.After(100 * time.Millisecond)
loop:
for len(items) < 3 {
    select {
    case <-timeout:
        break loop
    case v, ok := <-ch:
        if !ok {
            break loop
        }
        items = append(items, v)
    }
}
```

**After:**
```go
generator := func(yield func(int)) {
    for i := 0; i < 5; i++ {
        yield(i)
        time.Sleep(35*time.Millisecond)
    }
}

ch := lo.Generator(0, generator)

items1, length1, duration1, ok1 := lo.BufferWithTimeout(ch, 3, 100*time.Millisecond)
// []int{1, 2}, 2, 100ms, true
items2, length2, duration2, ok2 := lo.BufferWithTimeout(ch, 3, 100*time.Millisecond)
// []int{3, 4, 5}, 3, 75ms, true
items3, length3, duration2, ok3 := lo.BufferWithTimeout(ch, 3, 100*time.Millisecond)
// []int{}, 0, 10ms, false
```

---

## FanIn

**Before:**
```go
stream1 := make(chan int, 42)
stream2 := make(chan int, 42)
stream3 := make(chan int, 42)
all := make(chan int, 100)
var wg sync.WaitGroup
for _, s := range []<-chan int{stream1, stream2, stream3} {
    wg.Add(1)
    go func(c <-chan int) {
        defer wg.Done()
        for v := range c {
            all <- v
        }
    }(s)
}
go func() { wg.Wait(); close(all) }()
```

**After:**
```go
stream1 := make(chan int, 42)
stream2 := make(chan int, 42)
stream3 := make(chan int, 42)

all := lo.FanIn(100, stream1, stream2, stream3)
// <-chan int
```

---

## FanOut

**Before:**
```go
stream := make(chan int, 42)
children := make([]chan int, 5)
for i := range children {
    children[i] = make(chan int, 100)
}
go func() {
    for v := range stream {
        for _, c := range children {
            c <- v
        }
    }
    for _, c := range children {
        close(c)
    }
}()
```

**After:**
```go
stream := make(chan int, 42)

all := lo.FanOut(5, 100, stream)
// [5]<-chan int
```

---

## None

**Before:**
```go
collection := []int{0, 1, 2, 3, 4, 5}
subset := []int{0, 2}
subsetMap := map[int]struct{}{}
for _, v := range subset {
    subsetMap[v] = struct{}{}
}
none := true
for _, v := range collection {
    if _, ok := subsetMap[v]; ok {
        none = false
        break
    }
}
// false (0 and 2 are present)
```

**After:**
```go
b := None([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// false
b := None([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// true
```

---

## NoneBy

**Before:**
```go
collection := []int{1, 2, 3, 4}
noneMatch := true
for _, x := range collection {
    if x < 0 {
        noneMatch = false
        break
    }
}
// true
```

**After:**
```go
b := NoneBy([]int{1, 2, 3, 4}, func(x int) bool {
    return x < 0
})
// true
```

---

## IntersectBy

**Before:**
```go
list1 := []int{0, 1, 2, 3, 4, 5}
list2 := []int{0, 2}
seen := map[string]struct{}{}
for _, v := range list2 {
    seen[strconv.Itoa(v)] = struct{}{}
}
var result []int
for _, v := range list1 {
    if _, ok := seen[strconv.Itoa(v)]; ok {
        result = append(result, v)
    }
}
// []int{0, 2}
```

**After:**
```go
transform := func(v int) string {
    return strconv.Itoa(v)
}

result1 := lo.IntersectBy(transform, []int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// []int{0, 2}

result2 := lo.IntersectBy(transform, []int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// []int{0}

result3 := lo.IntersectBy(transform, []int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// []int{}
```

---

## WithoutBy

**Before:**
```go
type User struct {
    ID   int
    Name string
}
users := []User{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
    {ID: 3, Name: "Charlie"},
}
excludedIDs := map[int]struct{}{2: {}, 3: {}}
var filteredUsers []User
for _, u := range users {
    if _, skip := excludedIDs[u.ID]; !skip {
        filteredUsers = append(filteredUsers, u)
    }
}
// []User[{ID: 1, Name: "Alice"}]
```

**After:**
```go
type User struct {
    ID int
    Name string
}

users := []User{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
    {ID: 3, Name: "Charlie"},
}

getID := func(user User) int {
    return user.ID
}

excludedIDs := []int{2, 3}

filteredUsers := lo.WithoutBy(users, getID, excludedIDs...)
// []User[{ID: 1, Name: "Alice"}]
```

---

## WithoutNth

**Before:**
```go
s := []int{-2, -1, 0, 1, 2}
excludeIdx := map[int]struct{}{3: {}, 1: {}} // index 3 and index 1 (also -42 is out of bounds, ignored)
var result []int
for i, v := range s {
    if _, skip := excludeIdx[i]; !skip {
        result = append(result, v)
    }
}
// []int{-2, 0, 2}
```

**After:**
```go
subset := lo.WithoutNth([]int{-2, -1, 0, 1, 2}, 3, -42, 1)
// []int{-2, 0, 2}
```

---

## ElementsMatch

**Before:**
```go
list1 := []int{1, 1, 2}
list2 := []int{2, 1, 1}
freq1 := map[int]int{}
freq2 := map[int]int{}
for _, v := range list1 {
    freq1[v]++
}
for _, v := range list2 {
    freq2[v]++
}
match := reflect.DeepEqual(freq1, freq2)
// true
```

**After:**
```go
b := lo.ElementsMatch([]int{1, 1, 2}, []int{2, 1, 1})
// true
```

---

## ElementsMatchBy

**Before:**
```go
freq1 := map[string]int{}
freq2 := map[string]int{}
for _, item := range list1 {
    freq1[item.ID()]++
}
for _, item := range list2 {
    freq2[item.ID()]++
}
match := reflect.DeepEqual(freq1, freq2)
// true
```

**After:**
```go
b := lo.ElementsMatchBy(
    []someType{a, b},
    []someType{b, a},
    func(item someType) string { return item.ID() },
)
// true
```

---

## HasPrefix

**Before:**
```go
s := []int{1, 2, 3, 4}
prefix := []int{1, 2}
hasPrefix := len(s) >= len(prefix)
if hasPrefix {
    for i, v := range prefix {
        if s[i] != v {
            hasPrefix = false
            break
        }
    }
}
// true
```

**After:**
```go
ok := lo.HasPrefix([]int{1, 2, 3, 4}, []int{42})
// false

ok := lo.HasPrefix([]int{1, 2, 3, 4}, []int{1, 2})
// true
```

---

## HasSuffix

**Before:**
```go
s := []int{1, 2, 3, 4}
suffix := []int{3, 4}
hasSuffix := len(s) >= len(suffix)
if hasSuffix {
    offset := len(s) - len(suffix)
    for i, v := range suffix {
        if s[offset+i] != v {
            hasSuffix = false
            break
        }
    }
}
// true
```

**After:**
```go
ok := lo.HasSuffix([]int{1, 2, 3, 4}, []int{42})
// false

ok := lo.HasSuffix([]int{1, 2, 3, 4}, []int{3, 4})
// true
```

---

## FindLastIndexOf

**Before:**
```go
s := []string{"a", "b", "a", "b"}
lastIndex := -1
var lastVal string
for i, v := range s {
    if v == "b" {
        lastIndex = i
        lastVal = v
    }
}
// "b", 3, true
```

**After:**
```go
str, index, ok := lo.FindLastIndexOf([]string{"a", "b", "a", "b"}, func(i string) bool {
    return i == "b"
})
// "b", 4, true

str, index, ok := lo.FindLastIndexOf([]string{"foobar"}, func(i string) bool {
    return i == "b"
})
// "", -1, false
```

---

## FindKeyBy

**Before:**
```go
m := map[string]int{"foo": 1, "bar": 2, "baz": 3}
var foundKey string
var found bool
for k := range m {
    if k == "foo" {
        foundKey = k
        found = true
        break
    }
}
// "foo", true
```

**After:**
```go
result1, ok1 := lo.FindKeyBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(k string, v int) bool {
    return k == "foo"
})
// "foo", true

result2, ok2 := lo.FindKeyBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(k string, v int) bool {
    return false
})
// "", false
```

---

## FindUniquesBy

**Before:**
```go
s := []int{3, 4, 5, 6, 7}
freq := map[int]int{}
for _, v := range s {
    freq[v%3]++
}
var result []int
for _, v := range s {
    if freq[v%3] == 1 {
        result = append(result, v)
    }
}
// []int{5}
```

**After:**
```go
uniqueValues := lo.FindUniquesBy([]int{3, 4, 5, 6, 7}, func(i int) int {
    return i%3
})
// []int{5}
```

---

## FindDuplicatesBy

**Before:**
```go
s := []int{3, 4, 5, 6, 7}
freq := map[int]int{}
for _, v := range s {
    freq[v%3]++
}
seen := map[int]struct{}{}
var result []int
for _, v := range s {
    key := v % 3
    if freq[key] > 1 {
        if _, ok := seen[key]; !ok {
            seen[key] = struct{}{}
            result = append(result, v)
        }
    }
}
// []int{3, 4}
```

**After:**
```go
duplicatedValues := lo.FindDuplicatesBy([]int{3, 4, 5, 6, 7}, func(i int) int {
    return i%3
})
// []int{3, 4}
```

---

## MinIndex

**Before:**
```go
s := []int{1, 2, 3}
if len(s) == 0 {
    // 0, -1
}
minVal, minIdx := s[0], 0
for i, v := range s[1:] {
    if v < minVal {
        minVal = v
        minIdx = i + 1
    }
}
// 1, 0
```

**After:**
```go
min, index := lo.MinIndex([]int{1, 2, 3})
// 1, 0

min, index := lo.MinIndex([]int{})
// 0, -1

min, index := lo.MinIndex([]time.Duration{time.Second, time.Hour})
// 1s, 0
```

---

## MinIndexBy

**Before:**
```go
s := []string{"s1", "string2", "s3"}
if len(s) == 0 {
    // "", -1
}
minVal, minIdx := s[0], 0
for i, v := range s[1:] {
    if len(v) < len(minVal) {
        minVal = v
        minIdx = i + 1
    }
}
// "s1", 0
```

**After:**
```go
min, index := lo.MinIndexBy([]string{"s1", "string2", "s3"}, func(item string, min string) bool {
    return len(item) < len(min)
})
// "s1", 0

min, index := lo.MinIndexBy([]string{}, func(item string, min string) bool {
    return len(item) < len(min)
})
// "", -1
```

---

## Earliest

**Before:**
```go
times := []time.Time{time.Now(), time.Time{}}
earliest := times[0]
for _, t := range times[1:] {
    if t.Before(earliest) {
        earliest = t
    }
}
// 0001-01-01 00:00:00 +0000 UTC
```

**After:**
```go
earliest := lo.Earliest(time.Now(), time.Time{})
// 0001-01-01 00:00:00 +0000 UTC
```

---

## EarliestBy

**Before:**
```go
type foo struct{ bar time.Time }
items := []foo{{time.Now()}, {}}
earliest := items[0]
for _, item := range items[1:] {
    if item.bar.Before(earliest.bar) {
        earliest = item
    }
}
// {bar:{0001-01-01 00:00:00 +0000 UTC}}
```

**After:**
```go
type foo struct {
    bar time.Time
}

earliest := lo.EarliestBy([]foo{{time.Now()}, {}}, func(i foo) time.Time {
    return i.bar
})
// {bar:{2023-04-01 01:02:03 +0000 UTC}}
```

---

## MaxIndex

**Before:**
```go
s := []int{1, 2, 3}
if len(s) == 0 {
    // 0, -1
}
maxVal, maxIdx := s[0], 0
for i, v := range s[1:] {
    if v > maxVal {
        maxVal = v
        maxIdx = i + 1
    }
}
// 3, 2
```

**After:**
```go
max, index := lo.MaxIndex([]int{1, 2, 3})
// 3, 2

max, index := lo.MaxIndex([]int{})
// 0, -1

max, index := lo.MaxIndex([]time.Duration{time.Second, time.Hour})
// 1h, 1
```

---

## MaxIndexBy

**Before:**
```go
s := []string{"string1", "s2", "string3"}
if len(s) == 0 {
    // "", -1
}
maxVal, maxIdx := s[0], 0
for i, v := range s[1:] {
    if len(v) > len(maxVal) {
        maxVal = v
        maxIdx = i + 1
    }
}
// "string1" or "string3", 0 or 2
```

**After:**
```go
max, index := lo.MaxIndexBy([]string{"string1", "s2", "string3"}, func(item string, max string) bool {
    return len(item) > len(max)
})
// "string1", 0

max, index := lo.MaxIndexBy([]string{}, func(item string, max string) bool {
    return len(item) > len(max)
})
// "", -1
```

---

## Latest

**Before:**
```go
times := []time.Time{time.Now(), time.Time{}}
latest := times[0]
for _, t := range times[1:] {
    if t.After(latest) {
        latest = t
    }
}
// 2023-04-01 01:02:03 +0000 UTC
```

**After:**
```go
latest := lo.Latest(time.Now(), time.Time{})
// 2023-04-01 01:02:03 +0000 UTC
```

---

## LatestBy

**Before:**
```go
type foo struct{ bar time.Time }
items := []foo{{time.Now()}, {}}
latest := items[0]
for _, item := range items[1:] {
    if item.bar.After(latest.bar) {
        latest = item
    }
}
// {bar:{2023-04-01 01:02:03 +0000 UTC}}
```

**After:**
```go
type foo struct {
    bar time.Time
}

latest := lo.LatestBy([]foo{{time.Now()}, {}}, func(i foo) time.Time {
    return i.bar
})
// {bar:{2023-04-01 01:02:03 +0000 UTC}}
```

---

## FirstOrEmpty

**Before:**
```go
collection := []int{1, 2, 3}
var first int
if len(collection) > 0 {
    first = collection[0]
}
// 1
```

**After:**
```go
first := lo.FirstOrEmpty([]int{1, 2, 3})
// 1

first := lo.FirstOrEmpty([]int{})
// 0
```

---

## FirstOr

**Before:**
```go
collection := []int{}
fallback := 31
var first int
if len(collection) > 0 {
    first = collection[0]
} else {
    first = fallback
}
// 31
```

**After:**
```go
first := lo.FirstOr([]int{1, 2, 3}, 245)
// 1

first := lo.FirstOr([]int{}, 31)
// 31
```

---

## LastOrEmpty

**Before:**
```go
collection := []int{1, 2, 3}
var last int
if len(collection) > 0 {
    last = collection[len(collection)-1]
}
// 3
```

**After:**
```go
last := lo.LastOrEmpty([]int{1, 2, 3})
// 3

last := lo.LastOrEmpty([]int{})
// 0
```

---

## LastOr

**Before:**
```go
collection := []int{}
fallback := 31
var last int
if len(collection) > 0 {
    last = collection[len(collection)-1]
} else {
    last = fallback
}
// 31
```

**After:**
```go
last := lo.LastOr([]int{1, 2, 3}, 245)
// 3

last := lo.LastOr([]int{}, 31)
// 31
```

---

## NthOr

**Before:**
```go
s := []int{10, 20, 30, 40, 50}
n := 5
fallback := -1
var result int
if n < 0 {
    n = len(s) + n
}
if n >= 0 && n < len(s) {
    result = s[n]
} else {
    result = fallback
}
// -1
```

**After:**
```go
nth := lo.NthOr([]int{10, 20, 30, 40, 50}, 2, -1)
// 30

nth := lo.NthOr([]int{10, 20, 30, 40, 50}, -1, -1)
// 50

nth := lo.NthOr([]int{10, 20, 30, 40, 50}, 5, -1)
// -1 (fallback value)
```

---

## NthOrEmpty

**Before:**
```go
s := []int{10, 20, 30, 40, 50}
n := 5
var result int
if n < 0 {
    n = len(s) + n
}
if n >= 0 && n < len(s) {
    result = s[n]
}
// 0 (zero value for int)
```

**After:**
```go
nth := lo.NthOrEmpty([]int{10, 20, 30, 40, 50}, 2)
// 30

nth := lo.NthOrEmpty([]int{10, 20, 30, 40, 50}, -1)
// 50

nth := lo.NthOrEmpty([]int{10, 20, 30, 40, 50}, 5)
// 0 (zero value for int)

nth := lo.NthOrEmpty([]string{"apple", "banana", "cherry"}, 2)
// "cherry"

nth := lo.NthOrEmpty([]string{"apple", "banana", "cherry"}, 5)
// "" (zero value for string)
```

---

## SampleBy

**Before:**
```go
r := rand.New(rand.NewSource(42))
s := []string{"a", "b", "c"}
if len(s) == 0 {
    // ""
}
result := s[r.Intn(len(s))]
// a random string from []string{"a", "b", "c"}
```

**After:**
```go
import "math/rand"

r := rand.New(rand.NewSource(42))
lo.SampleBy([]string{"a", "b", "c"}, r.Intn)
// a random string from []string{"a", "b", "c"}, using a seeded random generator

lo.SampleBy([]string{}, r.Intn)
// ""
```

---

## SamplesBy

**Before:**
```go
r := rand.New(rand.NewSource(42))
s := []string{"a", "b", "c"}
n := 3
perm := r.Perm(len(s))
result := make([]string, n)
for i := 0; i < n; i++ {
    result[i] = s[perm[i]]
}
// []string{"a", "b", "c"} in random order
```

**After:**
```go
r := rand.New(rand.NewSource(42))
lo.SamplesBy([]string{"a", "b", "c"}, 3, r.Intn)
// []string{"a", "b", "c"} in random order, using a seeded random generator
```

---

## IsNotNil

**Before:**
```go
var i *int
result := i != nil
// false

var x int
result := true // non-pointer non-nil is always not nil
// true
```

**After:**
```go
var x int
lo.IsNotNil(x)
// true

var k struct{}
lo.IsNotNil(k)
// true

var i *int
lo.IsNotNil(i)
// false

var ifaceWithNilValue any = (*string)(nil)
lo.IsNotNil(ifaceWithNilValue)
// false
ifaceWithNilValue == nil
// true
```

---

## Nil

**Before:**
```go
var ptr *float64 = nil
// nil typed pointer
```

**After:**
```go
ptr := lo.Nil[float64]()
// nil
```

---

## EmptyableToPtr

**Before:**
```go
s := ""
var ptr *string
if s != "" {
    ptr = &s
}
// nil (because s is empty)

s2 := "hello world"
ptr2 := &s2
// *string{"hello world"}
```

**After:**
```go
ptr := lo.EmptyableToPtr(nil)
// nil

ptr := lo.EmptyableToPtr("")
// nil

ptr := lo.EmptyableToPtr([]int{})
// *[]int{}

ptr := lo.EmptyableToPtr("hello world")
// *string{"hello world"}
```

---

## FromSlicePtrOr

**Before:**
```go
str1 := "hello"
str2 := "world"
ptrs := []*string{&str1, nil, &str2}
fallback := "fallback value"
result := make([]string, len(ptrs))
for i, p := range ptrs {
    if p != nil {
        result[i] = *p
    } else {
        result[i] = fallback
    }
}
// []string{"hello", "fallback value", "world"}
```

**After:**
```go
str1 := "hello"
str2 := "world"

ptr := lo.FromSlicePtrOr([]*string{&str1, nil, &str2}, "fallback value")
// []string{"hello", "fallback value", "world"}
```

---

## IsNotEmpty

**Before:**
```go
x := 42
result := x != 0
// true

s := ""
result := s != ""
// false
```

**After:**
```go
lo.IsNotEmpty(0)
// false
lo.IsNotEmpty(42)
// true

lo.IsNotEmpty("")
// false
lo.IsNotEmpty("foobar")
// true

type test struct {
    foobar string
}

lo.IsNotEmpty(test{foobar: ""})
// false
lo.IsNotEmpty(test{foobar: "foobar"})
// true
```

---

## CoalesceOrEmpty

**Before:**
```go
vals := []int{0, 1, 2, 3}
var result int
for _, v := range vals {
    if v != 0 {
        result = v
        break
    }
}
// 1
```

**After:**
```go
result := lo.CoalesceOrEmpty(0, 1, 2, 3)
// 1

result := lo.CoalesceOrEmpty("")
// ""

var nilStr *string
str := "foobar"
result := lo.CoalesceOrEmpty(nil, nilStr, &str)
// &"foobar"
```

---

## CoalesceSlice

**Before:**
```go
slices := [][]int{nil, {1, 2, 3}, {4, 5, 6}}
var result []int
found := false
for _, s := range slices {
    if s != nil {
        result = s
        found = true
        break
    }
}
// [1, 2, 3], true
```

**After:**
```go
result, ok := lo.CoalesceSlice([]int{1, 2, 3}, []int{4, 5, 6})
// [1, 2, 3]
// true

result, ok := lo.CoalesceSlice(nil, []int{})
// []
// true

result, ok := lo.CoalesceSlice([]int(nil))
// []
// false
```

---

## CoalesceSliceOrEmpty

**Before:**
```go
slices := [][]int{nil, {1, 2, 3}}
var result []int
for _, s := range slices {
    if s != nil {
        result = s
        break
    }
}
// [1, 2, 3]
```

**After:**
```go
result := lo.CoalesceSliceOrEmpty([]int{1, 2, 3}, []int{4, 5, 6})
// [1, 2, 3]

result := lo.CoalesceSliceOrEmpty(nil, []int{})
// []
```

---

## CoalesceMap

**Before:**
```go
maps := []map[string]int{nil, {"1": 1, "2": 2}}
var result map[string]int
found := false
for _, m := range maps {
    if m != nil {
        result = m
        found = true
        break
    }
}
// {"1": 1, "2": 2}, true
```

**After:**
```go
result, ok := lo.CoalesceMap(map[string]int{"1": 1, "2": 2, "3": 3}, map[string]int{"4": 4, "5": 5, "6": 6})
// {"1": 1, "2": 2, "3": 3}
// true

result, ok := lo.CoalesceMap(nil, map[string]int{})
// {}
// true

result, ok := lo.CoalesceMap(map[string]int(nil))
// {}
// false
```

---

## CoalesceMapOrEmpty

**Before:**
```go
maps := []map[string]int{nil, {"1": 1, "2": 2}}
var result map[string]int
for _, m := range maps {
    if m != nil {
        result = m
        break
    }
}
// {"1": 1, "2": 2}
```

**After:**
```go
result := lo.CoalesceMapOrEmpty(map[string]int{"1": 1, "2": 2, "3": 3}, map[string]int{"4": 4, "5": 5, "6": 6})
// {"1": 1, "2": 2, "3": 3}

result := lo.CoalesceMapOrEmpty(nil, map[string]int{})
// {}
```

---

## AttemptWhile

**Before:**
```go
maxAttempts := 5
var count int
var err error
for i := 0; i < maxAttempts; i++ {
    err = doMockedHTTPRequest(i)
    count = i + 1
    if err != nil {
        if errors.Is(err, ErrBadRequest) {
            break // critical error, stop retrying
        }
        continue
    }
    break
}
```

**After:**
```go
count1, err1 := lo.AttemptWhile(5, func(i int) (error, bool) {
    err := doMockedHTTPRequest(i)
    if err != nil {
        if errors.Is(err, ErrBadRequest) { // let's assume ErrBadRequest is a critical error that needs to terminate the invoke
            return err, false // flag the second return value as false to terminate the invoke
        }

        return err, true
    }

    return nil, false
})
```

---

## AttemptWhileWithDelay

**Before:**
```go
maxAttempts := 5
delay := time.Millisecond
var count int
var elapsed time.Duration
var err error
start := time.Now()
for i := 0; i < maxAttempts; i++ {
    err = doMockedHTTPRequest(i)
    count = i + 1
    if err != nil {
        if errors.Is(err, ErrBadRequest) {
            break
        }
        time.Sleep(delay)
        continue
    }
    break
}
elapsed = time.Since(start)
```

**After:**
```go
count1, time1, err1 := lo.AttemptWhileWithDelay(5, time.Millisecond, func(i int, d time.Duration) (error, bool) {
    err := doMockedHTTPRequest(i)
    if err != nil {
        if errors.Is(err, ErrBadRequest) { // let's assume ErrBadRequest is a critical error that needs to terminate the invoke
            return err, false // flag the second return value as false to terminate the invoke
        }

        return err, true
    }

    return nil, false
})
```

---

## DebounceBy

**Before:**
```go
// Manually managing per-key debounce timers with a map and mutex
mu := sync.Mutex{}
timers := map[string]*time.Timer{}
f := func(key string, count int) {
    println(key + ": Called once after 100ms when debounce stopped invoking!")
}
debounce := func(key string) {
    mu.Lock()
    defer mu.Unlock()
    if t, ok := timers[key]; ok {
        t.Stop()
    }
    timers[key] = time.AfterFunc(100*time.Millisecond, func() { f(key, 0) })
}
```

**After:**
```go
f := func(key string, count int) {
    println(key + ": Called once after 100ms when debounce stopped invoking!")
}

debounce, cancel := lo.NewDebounceBy(100 * time.Millisecond, f)
for j := 0; j < 10; j++ {
    debounce("first key")
    debounce("second key")
}

time.Sleep(1 * time.Second)
cancel("first key")
cancel("second key")
```

---

## Throttle

**Before:**
```go
mu := sync.Mutex{}
lastCall := time.Time{}
f := func() {
    println("Called once in every 100ms")
}
throttle := func() {
    mu.Lock()
    defer mu.Unlock()
    if time.Since(lastCall) >= 100*time.Millisecond {
        lastCall = time.Now()
        f()
    }
}
```

**After:**
```go
f := func() {
	println("Called once in every 100ms")
}

throttle, reset := lo.NewThrottle(100 * time.Millisecond, f)

for j := 0; j < 10; j++ {
	throttle()
	time.Sleep(30 * time.Millisecond)
}

reset()
throttle()
```

---

## ThrottleWithCount

**Before:**
```go
mu := sync.Mutex{}
lastReset := time.Now()
callCount := 0
maxCount := 3
f := func() {
    println("Called three times in every 100ms")
}
throttle := func() {
    mu.Lock()
    defer mu.Unlock()
    if time.Since(lastReset) >= 100*time.Millisecond {
        lastReset = time.Now()
        callCount = 0
    }
    if callCount < maxCount {
        callCount++
        f()
    }
}
```

**After:**
```go
f := func() {
	println("Called three times in every 100ms")
}

throttle, reset := lo.NewThrottleWithCount(100 * time.Millisecond, f)

for j := 0; j < 10; j++ {
	throttle()
	time.Sleep(30 * time.Millisecond)
}

reset()
throttle()
```

---

## ThrottleBy

**Before:**
```go
// Manually maintain per-key throttle state with a map and mutex
mu := sync.Mutex{}
lastCalls := map[string]time.Time{}
f := func(key string) {
    println(key, "Called once in every 100ms")
}
throttle := func(key string) {
    mu.Lock()
    defer mu.Unlock()
    if time.Since(lastCalls[key]) >= 100*time.Millisecond {
        lastCalls[key] = time.Now()
        f(key)
    }
}
```

**After:**
```go
f := func(key string) {
	println(key, "Called three times in every 100ms")
}

throttle, reset := lo.NewThrottleByWithCount(100 * time.Millisecond, f)

for j := 0; j < 10; j++ {
	throttle("foo")
	time.Sleep(30 * time.Millisecond)
}

reset()
throttle()
```

---

## ThrottleByWithCount

**Before:**
```go
// Manually maintain per-key throttle state with count limit
mu := sync.Mutex{}
type state struct {
    lastReset time.Time
    count     int
}
states := map[string]*state{}
maxCount := 3
f := func(key string) {
    println(key, "Called three times in every 100ms")
}
throttle := func(key string) {
    mu.Lock()
    defer mu.Unlock()
    s, ok := states[key]
    if !ok {
        s = &state{}
        states[key] = s
    }
    if time.Since(s.lastReset) >= 100*time.Millisecond {
        s.lastReset = time.Now()
        s.count = 0
    }
    if s.count < maxCount {
        s.count++
        f(key)
    }
}
```

**After:**
```go
f := func(key string) {
	println(key, "Called three times in every 100ms")
}

throttle, reset := lo.NewThrottleByWithCount(100 * time.Millisecond, f)

for j := 0; j < 10; j++ {
	throttle("foo")
	time.Sleep(30 * time.Millisecond)
}

reset()
throttle()
```

---

## WaitForWithContext

**Before:**
```go
ctx := context.Background()
condition := func(i int) bool { return i >= 5 }
ticker := time.NewTicker(time.Millisecond)
defer ticker.Stop()
timeout := time.NewTimer(10 * time.Millisecond)
defer timeout.Stop()
i := 0
for {
    select {
    case <-ctx.Done():
        // context cancelled
        return
    case <-timeout.C:
        // timed out
        return
    case <-ticker.C:
        i++
        if condition(i) {
            return
        }
    }
}
```

**After:**
```go
ctx := context.Background()

alwaysTrue := func(_ context.Context, i int) bool { return true }
alwaysFalse := func(_ context.Context, i int) bool { return false }
laterTrue := func(_ context.Context, i int) bool {
    return i >= 5
}

iterations, duration, ok := lo.WaitForWithContext(ctx, alwaysTrue, 10*time.Millisecond, 2 * time.Millisecond)
// 1
// 1ms
// true

iterations, duration, ok := lo.WaitForWithContext(ctx, alwaysFalse, 10*time.Millisecond, time.Millisecond)
// 10
// 10ms
// false

iterations, duration, ok := lo.WaitForWithContext(ctx, laterTrue, 10*time.Millisecond, time.Millisecond)
// 5
// 5ms
// true
```

---

## TryWithErrorValue

**Before:**
```go
var panicValue any
func() {
    defer func() {
        panicValue = recover()
    }()
    panic("error")
}()
ok := panicValue == nil
// "error", false
```

**After:**
```go
err, ok := lo.TryWithErrorValue(func() error {
    panic("error")
    return nil
})
// "error", false
```

---

## Partial2 -> Partial5

**Before:**
```go
add := func(x, y, z int) int { return x + y + z }
// manually curry by wrapping
addWith42 := func(y, z int) int { return add(42, y, z) }

addWith42(10, 5)
// 57
```

**After:**
```go
add := func(x, y, z int) int { return x + y + z }
f := lo.Partial2(add, 42)

f(10, 5)
// 57

f(42, -4)
// 80
```