# shorturl
[原始連結](https://boards.greenhouse.io/dcard/jobs/3874841)

## 參考文獻
1. [完全GO语言实现的短网址微服务，可自用，可部署，美呆了](https://zhuanlan.zhihu.com/p/111573621)
  * 這篇文章算是最簡易的短網址產生程式
  * 使用套件為 redis, gin
  * 短網址產生方式為按照順序產生數字後轉換為 62 進位的字串
2. [How to Make a Custom URL Shortener Using Golang and Redis](https://intersog.com/blog/how-to-write-a-custom-url-shortener-using-golang-and-redis/)
  * 這篇文章有完整的分析及框架
  * 使用套件為 redis, fasthttp
  * 短網址產生方式為產生隨機亂數後判斷有無使用過，再轉為 62 進位的字串

## 假設與限制
* 轉成短網址時輸入使用 POST 
  * 輸入格式為 json, 需要輸入一個 url 及一個 RF3339 格式的有效時間
  * 執行後若無錯誤則回傳 json 資料，包含短網址字串以及產出的連結
  * 執行時若發生錯誤（無效的輸入資料），則回傳 400 狀態碼，及 Invalid request 的訊息
* 短網址導向原始網站時輸入是使用 GET，最後面為短網址
  * 直接導向原來的網站連結
  * 執行時若發生錯誤（無效或過期的輸入資料），則回傳 404 狀態碼，及 Page not found 的訊息
## 選用套件
* redis
  * 非關聯性資料庫符合短網址的需求
  * 使用簡單且普遍
* gin
  * 速度極快的 golang 網路框架
  * 目前在網路上使用非常廣泛

## 思路
* 短網址長度：大約 3~10 個字元，參數可以依需要而調整
  * 不會太短造成容易猜測
  * 不會太長符合以短網址的需求
  * 使用 3~10 個字元，可能性高達 62^3+62^4+...+62^10，一般情況綽綽有餘
* 短網址產生方式 : 直接產生隨機字串後檢查有無重複，可以減省 10 進位數字轉換 62 進位字元的步驟
* 隨機字串產生方式
  * 第一版：先隨機決定短址網有幾個字元 (3~10)，再依個數隨機產生短網址
    * 問題：產生長度為 3 的機率 = 長度為 10 的機率，造成長度為 3 時比較容易重複，需要更多次迴圈才能產生唯一的網址，降低效率
  * 第二版：先產生長度為 3 的隨機字串，之後最多跑 7 次迴圈，每次有一定的機率 (暫訂為 1/31，用參數調整) 跳出迴圈，若無跳出則每次增加一個隨機字元，最後長度一定介於 3~10 個字元
    * 成功解決第一版的問題，產生短網址的速度更快
  
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
## 測試方式
```bash
go test -v # 程式測試
go test -bench=. -run=none # 效能測試
```
* 測試內容包含
  * 測試加短網址，透過 http 方式呼叫
  * 測試加短網址，透過 handle 直接操作
  * 測試加短網址，輸入無效的資料，回傳 400
  * 測試加短網址，然後拜訪短網址
  * 測試加短網址，過時之後再拜訪短網址，回傳 404
  * 測試不存在的短網址，回傳 404
* 效能測試
  * 測試加短網址 (Handle) 的效能   
![image](https://user-images.githubusercontent.com/91168102/161761979-04a98845-3ca1-4ce1-b380-6f3375d8a847.png)
一秒內跑了 18747 筆，平均每次的回應時間為 63371 奈秒
