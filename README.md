# Go-5-Jungle

## 概要

Goを使って、AWS環境のセキュリティ設定を確認するCLIツールです。

AWSの各サービスから設定情報を取得し、ターミナル上で確認できるようにします。

## 目的

AWS環境では、IAM・S3・Security Groupなど複数のサービスにセキュリティ設定が分散しています。

このツールは、それらの設定情報をGoから取得し、AWS環境のセキュリティ状態を確認しやすくすることを目的としています。

まずはIAMから実装し、対象サービスを順次追加していきます。

## 対象サービス

| サービス | 確認内容 | 状況 |
|---|---|---|
| IAM | 接続中のIAMユーザー・ロールを確認 | 対応済み |
| S3 | セキュリティ設定を確認 | 今後対応 |
| Security Group | 通信ルールを確認 | 今後対応 |
| CloudTrail | ログ設定を確認 | 今後対応 |
| KMS | 暗号化設定を確認 | 今後対応 |

## 実行方法

- VS Codeでターミナルを開きます。
- `Go-5-Jungle` ディレクトリに移動します。
- 実行コマンドを入力します。
- 実行結果はターミナルに表示されます。

### 実行コマンド

AWSの認証情報とリージョンが設定済みの場合：

```sh
go run .
```

AWSプロファイルとリージョンを指定する場合：

```sh
go run . -profile development -region ap-northeast-1
```

### 実行結果の例

```text
Account: 123456789012
ARN:     arn:aws:sts::123456789012:assumed-role/Developer/example-session
User ID: AROAEXAMPLE:example-session
```

- `Account`：接続しているAWSアカウントID
- `ARN`：接続しているIAMユーザーまたはロール
- `User ID`：AWS上のユーザーまたはセッションの識別子

## 動作確認

`main.go` の処理が想定どおり動作するか、`main_test.go` を使って確認します。

### テスト対象

- 対象ディレクトリ：`Go-5-Jungle`
- 実装コード：`main.go`
- テストコード：`main_test.go`

### テスト実行コマンド

`Go-5-Jungle` ディレクトリで以下を実行します。

```sh
go test .
```

テストに成功すると、ターミナルに `ok` と表示されます。

## 必要な環境

- Go
- AWSの有効な認証情報
- AWSリージョンの設定
- AWSへ接続できるネットワーク環境

SSOを利用する場合は、事前にAWSへログインしておきます。

```sh
aws sso login --profile development
```

## 補足

このツールはAWSの設定やリソースを変更しません。

AWSから設定情報を取得し、セキュリティ状態を確認することを目的として開発しています。