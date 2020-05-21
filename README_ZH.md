這是一個非常簡陋的記帳程式，目的是處理一般的債務來往以及當中的貨幣轉換。
帳本文件 `book.bk` 記載帳本的基本設定。文件以命令形式書寫：

    currencies {string} ...
此命令跟隨多個字串，指示帳本可用之貨幣符號。示例：
    currencies USD JPY CNY

    account {string} {string}
此命令跟隨二個字串，分別是帳戶名及記載帳戶條目之文件名。示例：
    account MyName MyName.ent

帳戶條目記載語法：
    {uuid} (credit/debit) {float64} {string} [(@/=) {float64} {string}]
     編號      方    向       金 額    貨幣符號  [轉換方式  金 額    貨幣符號 ] 
每一筆帳均有獨立編號，用 uuid 作唯一編碼。`credit` 及 `debit` 乃關鍵字，指示資產往來方向（下面有詳細討論）。貨幣符號必須依足帳本設定。貨幣轉換數據可選，`@` 指示跟隨值為匯率，`=` 指示跟隨值為實際轉換金額。示例：
    6ba7b810-9dad-11d1-80b4-00c04fd430c8 debit 20 JPY
    6259a3b7-50e4-411c-83bd-d35abf0f39cd credit 100 USD @ 7 CNY
    9fe0bfd2-fa82-4a22-b48b-10a35ce1ddc6 debit 100 JPY = 1 USD

討論：
1. 語法。關於為何不使用數學正負號表示，目的有二：1) 利於人讀，且保證書寫人清楚知道往來方向；2) 為免使貨幣轉換記錄出現不必要符號。
2. 明確往來方向。`credit` 是對於該 account 之持有人的資產流出；`debit` 則為流入。
3. 本記賬方式依賴借貸發生人的信用，並實時提交作公共監督。當債務發生時，資產流入方記為負數 liability，資產流出方記為正數 liability。理論上不設「還款」之概念，而任何負數 liability 一方有義務以 liability 清零為目的而產生資金往來。產生「還款」仍以一般的 `credit` 記錄（收款方為 `debit`），並允許不必須向原來的「借出方」償還。
4. 由於每個 account 的 liability 按貨幣不同區分記錄，因此資產來往產生的金額應按實際填寫，但允許可選添加貨幣轉換以實現填補某幣種 liability 以實現清零。