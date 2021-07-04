# arithmetic

本仓库用于作业和笔记提交

## 写在最前

- 学习从来没有速成之法.别人可以告诉你别人的思路,但那不意味着这个思路就属于你.
- 如果你从未接触过算法,那么你无法在短时间内学会它.除非你是天才.当然,如果你是天才,看到这里已经可以关闭浏览器的该标签页了
- 做题只是手段,而非目的.做题只是一种"检测自己有什么是不会的"的一种手段,因此应该以开放和拥抱变化的心态去看待做题这件事.而非看重自己刷了多少题,会做多少题,还有多少题不会做之类的
- 每一道你没有做出来的题,都证明了一件事:你至少对一个数据结构或算法的定义或性质并没有达到能够足以应用它们解决实际问题的地步
- 我不需要表演型人格
- 不要用战术上的勤奋,来掩盖战略上的懒惰
- 先实现一个数据结构,你才有掌握这个数据结构所具备的性质的可能性.然后你去做题才有意义.否则不过是井中月,水中花罢了.

## 目录说明

### 一级目录

```
./
├── README.md
└── class1
```

`README.md`:仓库说明

`class1 /`:本目录用于存放每课时的代码和笔记.命名风格为:`class + 课时编号`.课时编号从1开始累加计算.

### 二级目录

```
./class1
├── code
└── note

2 directories, 0 files
```

`code/`:本目录用于存放代码

`node/`:本目录用于存放笔记

### 三级目录

#### class1/code目录

```
./
├── 26-removeDuplicates-InPlace-advance.go
├── 26-removeDuplicates-InPlace.go
├── 26-removeDuplicates-useAnotherArray.go
├── 283-moveZeroes-InPlace.go
├── 283-moveZeroes-UseAnotherArray.go
├── 88-merge-In-Place.go
├── 88-merge-useAnotherArray.go
└── dataStructure

1 directory, 7 files
```

`88-merge-In-Place.go`:命名风格为:`题号-函数名-简要描述做法`

`dataStructure/`:本目录用于存放自己实现的数据结构

#### class1/note目录

```
./
├── image
└── 第1课 数组、链表、栈、队列.md

1 directory, 1 file
```

`第1课 数组、链表、栈、队列.md`:命名风格:`第N课 课时内容`.其中N为课时编号,从1开始累加计算.

`image/`:本目录用于存放笔记中的图片.

### 四级目录

#### class1/code/dataStructure目录

```
./
├── array
│   └── array.go
└── linkedList
    ├── linkedList.go
    └── node.go

2 directories, 3 files
```

`array/`: 代码以包形式组织.包名为数据结构名称.

`array/array.go`:包中的文件以类为单位组织.



