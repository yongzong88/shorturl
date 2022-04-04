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
* 輸入格式為json,需要輸入一個url及一個3339格式的過期時間
* 短網址導向原始網站輸入是使用GET

## 參考文獻
1. [完全GO语言实现的短网址微服务，可自用，可部署，美呆了](https://zhuanlan.zhihu.com/p/111573621)
2. [How to Make a Custom URL Shortener Using Golang and Redis](https://intersog.com/blog/how-to-write-a-custom-url-shortener-using-golang-and-redis/)



這兩篇在產生短網址的方式都是產生一個整數再轉成62進位的字串


A. 短網址產生方式
B. 我的短網址產生方式
C. 優劣比較

### 資料庫選用



### 測試
```bash
go test -v # 語法測試
go test -bench=. -run=none # 測試效能
```

### 選用套件

