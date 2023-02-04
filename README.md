# Vot!ng
## URL
https://voting0195.an.r.appspot.com/<br/>
## 概要
Vot!ngは**投票Webアプリケーション**です。<br/>
ユーザーは大きく分けて**投稿**と**投票**の二つのアクションが可能です。<br/>

ホーム画面
<img width="1440" alt="Screen Shot 2023-02-04 at 20 41 58" src="https://user-images.githubusercontent.com/98641436/216765400-939faa42-97e7-4f4f-97fb-208bfe3f0b95.png">
投票画面
<img width="1440" alt="Screen Shot 2023-02-04 at 20 42 14" src="https://user-images.githubusercontent.com/98641436/216765425-a60e96eb-b256-49b1-b684-b693dd636373.png">

## 機能紹介
- アカウント
  - 登録、ログイン機能
  - アカウント情報変更機能
  - 退会機能
- 投稿者機能
  - 投稿の情報、投票対象者の設定
    - タイトル、選択肢、対象年齢、対象性別など
  - 投稿の削除
  - 投票の結果確認
- 投票者機能
  - 選択肢から投票
  - 結果だけ確認
- お問合せ機能

## ポイント
- セキュリティ面を考慮し、パスワードをハッシュ化
- 投稿をクリックした際の条件分岐
  - 投票したことがない場合
    - 指定の条件に一致する場合は投票、または結果だけの確認ができる
  - 投票したことがある、結果だけ確認したことがある、または指定の条件に当てはまっていない場合
    - 結果が表示される
  - ログインしていない
    - 選択肢が表示され、選択するとログインページにリダイレクト
  
## 環境・ツール・言語
- Mac OS
- Visual Studio Code
- Microsoft Azure
- TablePlus
- Go 1.19
  - Gin
  - GORM
  - sessions
  - multitemplate
- HTML5
- CSS
- JavaScript
