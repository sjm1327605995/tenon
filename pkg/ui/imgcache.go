package ui

import "container/list"

// ---- 图片缓存（按字节预算的 LRU）----
//
// 解码后的位图按 src 缓存，避免同一张图重复解码。缓存必须有上限：曾经是一个只增不减的
// map，翻遍一个大图库就会把内存吃光（一张 177x254 的卡图解码后约 180KB，上万张即数 GB）。
//
// 预算按「解码后的字节数」估算（宽 x 高 x 4），而不是按条数 —— 条数对不上实际占用：
// 一张 4K 壁纸抵得上几百张缩略图。
//
// 只在渲染线程访问（reconcile 与 drainPosts），故无需加锁；后台 goroutine 只做 IO/解码。

// imgCacheBudget 是缓存的字节上限。默认 256MB：足够放下上千张卡图，又不至于让长时间
// 运行的应用把内存吃光。可用 ImageCacheBudget 调整。
var imgCacheBudget int64 = 256 << 20

// ImageCacheBudget 设置图片缓存的字节上限（默认 256MB）。超出时按最近最少使用淘汰。
//
// 调大：图多且反复出现（图库、卡牌游戏），愿意用内存换解码开销。
// 调小：内存紧张，或图片很少复用。设为 0 则每次都重新解码。
func ImageCacheBudget(bytes int64) {
	if bytes < 0 {
		bytes = 0
	}
	imgCacheBudget = bytes
	trimImageCache()
}

type imgEntry struct {
	src  string
	img  bitmap
	size int64
}

var (
	imgCache  = map[string]*list.Element{} // src -> lru 中的元素
	imgLRU    = list.New()                 // 队首=最近用过，队尾=最久未用
	imgBytes  int64                        // 当前缓存占用的估算字节数
	imgEvicts int                          // 累计淘汰次数（诊断用）
)

// imgSizeOf 估算一张位图解码后的占用：宽 x 高 x 4（RGBA）。
func imgSizeOf(b bitmap) int64 {
	w, h := b.Size()
	if w <= 0 || h <= 0 {
		return 0
	}
	return int64(w) * int64(h) * 4
}

// lookupImage 取缓存并把它标记为最近使用。
func lookupImage(src string) (bitmap, bool) {
	e, ok := imgCache[src]
	if !ok {
		return nil, false
	}
	imgLRU.MoveToFront(e)
	return e.Value.(*imgEntry).img, true
}

// storeImage 放入缓存并按预算淘汰最久未用的条目。
func storeImage(src string, img bitmap) {
	if e, ok := imgCache[src]; ok { // 已在缓存：替换并前移
		old := e.Value.(*imgEntry)
		imgBytes -= old.size
		old.img, old.size = img, imgSizeOf(img)
		imgBytes += old.size
		imgLRU.MoveToFront(e)
		return
	}
	size := imgSizeOf(img)
	// 单张就超过整个预算的图不进缓存：否则为给它腾地方会把其余全淘汰，
	// 然后它自己也放不下、同样被淘汰 —— 净结果是缓存被冲空，一张不剩。
	// 不缓存它只是每次重新解码，至少不牵连别人。
	if size > imgCacheBudget {
		return
	}
	ent := &imgEntry{src: src, img: img, size: size}
	imgCache[src] = imgLRU.PushFront(ent)
	imgBytes += ent.size
	trimImageCache()
}

// trimImageCache 从队尾淘汰，直到占用回到预算内。
//
// 正在显示的图也可能被淘汰 —— 那只是丢掉 Go 侧的引用，节点下一帧会重新加载；
// 预算小于一屏所需时会来回抖动，但那是预算设得太小，不是这里的问题。
func trimImageCache() {
	for imgBytes > imgCacheBudget {
		e := imgLRU.Back()
		if e == nil {
			return
		}
		ent := e.Value.(*imgEntry)
		imgLRU.Remove(e)
		delete(imgCache, ent.src)
		imgBytes -= ent.size
		imgEvicts++
	}
}

// ImageCacheStats 返回图片缓存的当前状态：条数、估算字节数、累计淘汰次数。
func ImageCacheStats() (entries int, bytes int64, evictions int) {
	return imgLRU.Len(), imgBytes, imgEvicts
}
