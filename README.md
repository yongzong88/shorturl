# shorturl
[原始連結](https://boards.greenhouse.io/dcard/jobs/3874841)

## 專案執行方式
1. 安裝redis
```bash
docker pull redis
docker run --name redis-lab -p 6379:6379 -d redis
```
2. git我的專案
```bash
git clone https://github.com/yongzong88/shorturl.git
cd shorturl/
go install
go run main.go
```

## 假設與限制
* 轉成短網址時輸入是使用POST 
* 輸入格式為json,需要輸入一個url及一個3339格式的時間
* 短網址導向原始網站輸入是使用GET

## 參考文獻
1. [完全GO语言实现的短网址微服务，可自用，可部署，美呆了](https://zhuanlan.zhihu.com/p/111573621)
2. [How to Make a Custom URL Shortener Using Golang and Redis](https://intersog.com/blog/how-to-write-a-custom-url-shortener-using-golang-and-redis/)

## 思路
在參考第二篇文章的分析後決定使用1-10,A-Z,a-z作為短網址的字元,並產生長度為3-10個字元的短網址
在參考過網路上其他人的做法之後,發現都是使用整數轉換為62進位的字串
其中一個是按照1,2,3,...的方式產生,這樣會讓短網址變成可預測的
另一個則是髓機產生亂數,再判斷是否使用過
我的方式是直接隨機產生一個字串,再判斷是否使用過
一開始我是隨機產生一個3-10的數字N
然後用for迴圈跑N次來產生一個3-10個字的隨機字串
但由於是先產生3-10的隨機數,所以選到3跟10的機率是相等的
3個字的不重複字串只有62^3次方種
10個字的有62^10次方種
這樣會造成N=3時容易重複,使的效率降低
所以改成先產生一個長度3的隨機字串
接者取亂數0-61+2=N
然後跑7次for迴圈
如果N>=61則跳出
否則增加一個隨機的字元
這樣可以讓3-10個字串中長度越長的機率越高,減少重複的情況發生

### 測試
```bash
go test -v # 語法測試
go test -bench=. -run=none # 測試效能
```

### 選用套件
* redis
  由於製作短網址的每一筆資料並沒有互相關聯,所以使用非關聯性資料庫
  其中redis相當普遍也較為簡單,所以選用redis作為資料庫
* gin
  gin是非常快速的golang網路框架
  也非常普遍且使用容易,所以選用gin作為網路框架
