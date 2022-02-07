package common

// 实现一致性哈希算法
import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

//声明新切片类型
type units []uint32

// Len 返回切片长度
func (x units) Len() int {
	return len(x)
}

// Less 比对两个数大小
func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

// Swap 切片中两个值的交换
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

//当hash环上没有数据时，提示错误
var errEmpty = errors.New("hash环没有数据")

// Consistent 创建结构体保存一致性hash信息
type Consistent struct {
	circle       map[uint32]string //hash环，key为哈希值，值存放节点的信息
	sortedHashes units             //已经排序的节点hash切片
	VirtualNode  int               //虚拟节点个数，用来增加hash的平衡性
	sync.RWMutex                   //map 读写锁
}

// NewConsistent 创建一致性hash算法结构体，设置默认节点数量
func NewConsistent() *Consistent {
	return &Consistent{
		circle:      make(map[uint32]string), //初始化变量
		VirtualNode: 20,                      //设置虚拟节点个数
	}
}

// Add 向hash环中添加节点
func (c *Consistent) Add(element string) {
	c.Lock()         //加锁
	defer c.Unlock() //解锁
	c.add(element)
}

// Remove 删除一个节点
func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

// Get 根据数据标示获取最近的服务器节点信息
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()               //添加锁
	defer c.RUnlock()       //解锁
	if len(c.circle) == 0 { //如果为零则返回错误
		return "", errEmpty
	}
	key := c.hashKey(name) //计算hash值
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}

// 自动生成key值
func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

// 获取hash位置
func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte                          //声明一个数组长度为64
		copy(scratch[:], key)                         //拷贝数据到数组中
		return crc32.ChecksumIEEE(scratch[:len(key)]) //使用IEEE 多项式返回数据的CRC-32校验和
	}
	return crc32.ChecksumIEEE([]byte(key))
}

// 更新排序，方便查找
func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) { //判断切片容量，是否过大，如果过大则重置
		hashes = nil
	}
	for k := range c.circle { //添加hashes
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)       //对所有节点hash值进行排序，方便之后进行二分查找
	c.sortedHashes = hashes //重新赋值
}

// 添加节点
func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ { //循环虚拟节点，设置副本
		c.circle[c.hashKey(c.generateKey(element, i))] = element //根据生成的节点添加到hash环中
	}
	c.updateSortedHashes() //更新排序
}

//删除节点
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

//顺时针查找最近的节点
func (c *Consistent) search(key uint32) int {
	//查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	i := sort.Search(len(c.sortedHashes), f) //使用"二分查找"算法来搜索指定切片满足条件的最小值
	if i >= len(c.sortedHashes) {            //如果超出范围则设置i=0
		i = 0
	}
	return i
}
