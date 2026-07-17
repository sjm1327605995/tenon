package ui

import (
	"container/list"
	"testing"
)

type fakeBitmap struct{ w, h int }

func (f *fakeBitmap) Size() (int, int) { return f.w, f.h }

// 缓存必须有上限。曾经是只增不减的 map —— 翻遍一个大图库就把内存吃光。
func TestImageCacheEvictsOverBudget(t *testing.T) {
	resetImageCache(t)
	// 每张 100x100 -> 40000 字节；预算给 5 张的量
	const one = 100 * 100 * 4
	ImageCacheBudget(one * 5)

	for i := 0; i < 20; i++ {
		storeImage(string(rune('a'+i)), &fakeBitmap{100, 100})
	}
	n, bytes, evicts := ImageCacheStats()
	t.Logf("放入 20 张、预算 5 张：缓存 %d 张 / %d 字节 / 淘汰 %d 次", n, bytes, evicts)
	if bytes > one*5 {
		t.Errorf("占用 %d 超出预算 %d —— 没淘汰", bytes, one*5)
	}
	if n > 5 {
		t.Errorf("缓存 %d 张，超过预算能放的 5 张", n)
	}
	if evicts != 15 {
		t.Errorf("淘汰 %d 次，want 15（放 20 留 5）", evicts)
	}
}

// 淘汰必须是「最近最少使用」，不是随便挑一个：反复用到的图应当留下。
func TestImageCacheEvictsLeastRecentlyUsed(t *testing.T) {
	resetImageCache(t)
	const one = 100 * 100 * 4
	ImageCacheBudget(one * 3)

	storeImage("keep", &fakeBitmap{100, 100})
	storeImage("b", &fakeBitmap{100, 100})
	storeImage("c", &fakeBitmap{100, 100})

	// 反复命中 keep，让它成为最近使用；此时 b 变成最久未用
	for i := 0; i < 3; i++ {
		if _, ok := lookupImage("keep"); !ok {
			t.Fatal("keep 不在缓存里，前提不成立")
		}
	}
	storeImage("d", &fakeBitmap{100, 100}) // 超预算，应淘汰最久未用的 b

	if _, ok := lookupImage("keep"); !ok {
		t.Error("反复用到的 keep 被淘汰了 —— 不是 LRU")
	}
	if _, ok := lookupImage("b"); ok {
		t.Error("最久未用的 b 还在 —— 淘汰选错了对象")
	}
	for _, s := range []string{"c", "d"} {
		if _, ok := lookupImage(s); !ok {
			t.Errorf("%s 不该被淘汰", s)
		}
	}
}

// 预算按字节而非条数：一张 200x200 顶 16 张 50x50。
func TestImageCacheBudgetIsBytesNotCount(t *testing.T) {
	resetImageCache(t)
	const small = 50 * 50 * 4 // 10000
	ImageCacheBudget(small * 20)

	for i := 0; i < 8; i++ { // 8 张小图，占 8/20
		storeImage(string(rune('a'+i)), &fakeBitmap{50, 50})
	}
	storeImage("big", &fakeBitmap{200, 200}) // 160000 字节 = 16 张小图，会挤掉大部分
	n, b, ev := ImageCacheStats()
	t.Logf("8 张小图 + 1 张 200x200（预算 %d 字节）：缓存 %d 张 / %d 字节 / 淘汰 %d 次",
		small*20, n, b, ev)

	if b > small*20 {
		t.Errorf("占用 %d 超预算 %d", b, small*20)
	}
	if _, ok := lookupImage("big"); !ok {
		t.Error("大图没进缓存 —— 它并未超过总预算")
	}
	if ev < 4 {
		t.Errorf("只淘汰了 %d 次 —— 一张大图应当顶掉多张小图，预算没按字节算", ev)
	}
}

// 单张超过整个预算的图不进缓存，也不得牵连已缓存的其它图。
// 否则为给它腾地方会把缓存清空、它自己又放不下，净结果一张不剩（曾如此）。
func TestOversizedImageDoesNotWipeCache(t *testing.T) {
	resetImageCache(t)
	const small = 50 * 50 * 4
	ImageCacheBudget(small * 10)

	for i := 0; i < 8; i++ {
		storeImage(string(rune('a'+i)), &fakeBitmap{50, 50})
	}
	storeImage("huge", &fakeBitmap{500, 500}) // 1000000 字节，远超整个预算
	n, b, _ := ImageCacheStats()
	t.Logf("放入一张超预算巨图后：缓存 %d 张 / %d 字节", n, b)

	if _, ok := lookupImage("huge"); ok {
		t.Error("超预算的巨图进了缓存")
	}
	if n != 8 {
		t.Errorf("缓存还剩 %d 张，want 8 —— 一张巨图把别人冲掉了", n)
	}
}

// resetImageCache 清空缓存与统计，并在测试结束后还原默认预算。
func resetImageCache(t *testing.T) {
	t.Helper()
	old := imgCacheBudget
	imgCache = map[string]*list.Element{}
	imgLRU = list.New()
	imgBytes, imgEvicts = 0, 0
	t.Cleanup(func() {
		imgCache = map[string]*list.Element{}
		imgLRU = list.New()
		imgBytes, imgEvicts = 0, 0
		ImageCacheBudget(old)
	})
}
