# gcontainer -- Go 语言的全泛型数据结构(容器)实现

[English](README.md) | 简体中文

包含 Go 语言实现的各类数据结构(容器)，全面支持泛型，目标是复刻 Java 类库中的 java.util 和 java.util.concurrent 包里的数据结构实现。

[![License MIT](https://img.shields.io/badge/License-MIT-red.svg)](LICENSE)
[![Golang](https://img.shields.io/badge/Language-go1.20%2B-blue)](https://go.dev/)

```go
import "github.com/wesleywu/gcontainer"
```

## 前言

在 Golang 的生态中，有不少优秀的库都实现了各类数据结构，例如

- [stl4go](https://github.com/chen3feng/stl4go) 实现了大部分泛型容器和算法，类似 C++ 中的 STL
- [lo](https://github.com/samber/lo) Lodash风格的 Go 实现
- [goframe container](https://github.com/gogf/gf/tree/master/container) 实现了绝大多数数据结构，但不支持泛型

为什么还要再造一个轮子呢？
* 作为一个二十多年的Java程序员，习惯了Java类库中对各类容器接口、方法的精准命名和优雅实现，希望能将之前的肌肉记忆应用到 Go 语言当中
* 在使用 Go 开发的过程中，频繁的需要用到各种数据结构，比较看来 goframe 提供的库是最全的，完全满足需求，且都实现了线程安全。但可惜 goframe 会在相对较长的时间里，保持对 go 1.15 版本的兼容，而不会支持泛型。
 
本库将 goframe 的数据结构实现从 goframe 项目中剥离出来，改为全泛型支持。
参照 Java 类库做了 interface 的定义，并对实现 struct 和 method(function) 进行适当的更名。

**所有容器的实现，在创建时都可选择是否线程安全**

## 容器接口

已实现（实现中）的容器接口，与 Java 8 类库的对应关系如下：

| gcontainer      | [java.util](https://docs.oracle.com/javase/8/docs/api/java/util/package-summary.html)       | 用途                   |
|-----------------|---------------------------------------------------------------------------------------------|----------------------|
| Collection[T]   | [Collection<T>](https://docs.oracle.com/javase/8/docs/api/java/util/Collection.html)        | 集合和列表容器的基础接口      | 
| Set[T]          | [Set<T>](https://docs.oracle.com/javase/8/docs/api/java/util/Set.html)                      | 集合容器的接口              |
| SortedSet[T]    | [NavigableSet<T>](https://docs.oracle.com/javase/8/docs/api/java/util/NavigableSet.html)    | 有序集合容器的接口           |
| List[T]         | [List<T>](https://docs.oracle.com/javase/8/docs/api/java/util/List.html)                    | 列表容器的接口              |
| Map[K, V]       | [Map<K, V>](https://docs.oracle.com/javase/8/docs/api/java/util/Map.html)                   | key-value 关联容器的接口    |
| SortedMap[K, V] | [NavigableMap<K, V>](https://docs.oracle.com/javase/8/docs/api/java/util/NavigableMap.html) | 有序 key-value 关联容器的接口 |
| MapEntry[K, V]  | [Map.Entry<K, V>](https://docs.oracle.com/javase/8/docs/api/java/util/Map.Entry.html)       | 有序 key-value 关联容器的接口 |
| 实现中           | [Queue<T>](https://docs.oracle.com/javase/8/docs/api/java/util/Queue.html)                  | 先进先出的队列的接口 |
| 实现中           | [Deque<T>](https://docs.oracle.com/javase/8/docs/api/java/util/Deque.html)                  | 双端队列的接口 |

## 容器结构体

已实现的容器结构体，与 Java 8 类库的对应关系如下：

| gcontainer    | [java.util](https://docs.oracle.com/javase/8/docs/api/java/util/package-summary.html)      | [java.util.concurrent](https://docs.oracle.com/javase/8/docs/api/java/util/concurrent/package-summary.html)              | 用途                            |
|---------------|--------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|-------------------------------|
| ArrayList[T]  | [ArrayList<T>](https://docs.oracle.com/javase/8/docs/api/java/util/ArrayList.html)         | 无对应                                                                                                                    | 基于数组的列表                       |
| LinkedList[T] | [LinkedList<T>](https://docs.oracle.com/javase/8/docs/api/java/util/LinkedList.html)       | 无对应                                                                                                                    | 基于双向链表的列表                     | 
| HashSet[T]    | [HashSet<T>](https://docs.oracle.com/javase/8/docs/api/java/util/HashSet.html)             | 无对应                                                                                                                    | 基于map的去重集合                    |    
| TreeSet[T]    | [TreeSet<T>](https://docs.oracle.com/javase/8/docs/api/java/util/TreeSet.html)             | 无对应                                                                                                                    | 基于TreeMap红黑树的去重排序集合           | 
| HashMap[K, V] | [HashMap<K, V>](https://docs.oracle.com/javase/8/docs/api/java/util/HashMap.html)          | [ConcurrentHashMap<K, V>](https://docs.oracle.com/javase/8/docs/api/java/util/concurrent/ConcurrentHashMap.html)         | 基于map的key-value关联表            | 
| TreeMap[K, V] | [TreeMap<K, V>](https://docs.oracle.com/javase/8/docs/api/java/util/TreeMap.html)          | [ConcurrentSkipListMap<K, V>](https://docs.oracle.com/javase/8/docs/api/java/util/concurrent/ConcurrentSkipListMap.html) | 基于红黑树的排序key-value关联表          |  
| AVLTree[K, V] | 无对应                                                                                      |                                                                                                                          | 基于avl-tree的排序key-value关联表     |    
| BTree[K, V]   | 无对应                                                                                      |                                                                                                                          | 基于b-tree的排序key-value关联表       |    
| Queue[T]      | [PriorityQueue<T>](https://docs.oracle.com/javase/8/docs/api/java/util/PriorityQueue.html) | [LinkedBlockingQueue<T>](https://docs.oracle.com/javase/8/docs/api/java/util/concurrent/LinkedBlockingQueue.html)        | 基于channel的阻塞队列                |   
| Ring[T]       | 无对应                                                                                      |                                                                                                                          | 环容器                           |   
