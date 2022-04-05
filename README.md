# shorturl
[原始連結](https://boards.greenhouse.io/dcard/jobs/3874841)

## 參考文獻
1. [完全GO语言实现的短网址微服务，可自用，可部署，美呆了](https://zhuanlan.zhihu.com/p/111573621)
  * 這篇文章算是最簡易的短網址產生程式
  * 使用套件為redis,gin
  * 短網址產生方式為按照順序產生數字後轉換為62進位的字串
2. [How to Make a Custom URL Shortener Using Golang and Redis](https://intersog.com/blog/how-to-write-a-custom-url-shortener-using-golang-and-redis/)
  * 這篇文章有完整的分析及框架
  * 使用套件為redis,fasthttp
  * 短網址產生方式為產生隨機亂數後判斷有無使用過,再轉為62進位的字串

## 假設與限制
* 轉成短網址時輸入是使用POST 
* 輸入格式為json,需要輸入一個url及一個3339格式的時間
* 執行後會產生一個json格式,裡面包含3-10個字的字串短網址及一個短網址連結
* 短網址導向原始網站時輸入是使用GET

## 選用套件
* redis
  * 非關聯性資料庫符合短網址的需求
  * 使用簡單且普遍
* gin
  * 速度極快的golang網路框架
  * 使用簡單且普遍

## 思路
* 產生短網址長度:3-10個字元較為剛好且數量足夠
* 短網址產生方式:直接產生隨機字串後檢查有無重複,可以減省10進位數字轉換62近位字元的步驟
* 隨機字串產生方式
  * 隨機跑3-10次的迴圈每次加一隨機字元,產生短網址
    * 問題:產生長度為3的機率=長度為10的機率,造成長度為3時容易重復降低效率
  * 先產生長度為3的隨機字串後跑最多7次的迴圈,每次有1/31的機率跳出迴圈,若無跳出則加一個隨機字元
    * 成功解決第一個方法的問題
  
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

### 測試
```bash
go test -v # 語法測試
go test -bench=. -run=none # 測試效能
```

