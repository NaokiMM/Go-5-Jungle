# Go-5-Jungle

現在AWSに接続しているIAMユーザーまたはロールを表示するGo CLIです。AWS SDK for Go v2の標準認証情報チェーンとSTS `GetCallerIdentity`を使い、アカウントID・ARN・ユーザーIDを取得します。AWSリソースやIAM設定は変更しません。

## 設定

Go 1.26.4以降、有効なAWS認証情報、リージョン、AWSへのネットワーク接続が必要です。

SDKが環境変数、`~/.aws/config`・`~/.aws/credentials`、SSO、Web Identity、ECS・EC2ロールなどから認証情報を解決します。キーをコードに埋め込まないでください。認証情報ファイルをGitに追加しないでください。

SSOプロファイルを準備する例（AWS CLIは設定・ログイン時のみ必要）:

```sh
aws configure sso --profile development
aws sso login --profile development
```

既存の認証情報も使用できます。環境変数で一時認証情報を設定する場合は`AWS_ACCESS_KEY_ID`・`AWS_SECRET_ACCESS_KEY`に加え`AWS_SESSION_TOKEN`が必要です。`-profile`を明示するとSDKが指定プロファイルを選択します。`AWS_PROFILE`を使う場合、環境変数のアクセスキーが優先されることがあります。実際の接続主体は出力ARNで確認してください。

## 実行

```sh
go run . -profile development -region ap-northeast-1
```

認証情報とリージョンが設定済みの場合:

```sh
go run .
```

```sh
AWS_PROFILE=development AWS_REGION=ap-northeast-1 go run .
go run . -help
go run . -profile development -region ap-northeast-1 -timeout 60s
```

出力例（架空の値）:

```text
Account: 123456789012
ARN:     arn:aws:sts::123456789012:assumed-role/Developer/example-session
User ID: AROAEXAMPLE:example-session
```

`iam::…:user/…`はIAMユーザー、`sts::…:assumed-role/…`はロールとセッションを示します。SSOも通常はロールとして表示されます。

成功時は終了コード0、失敗時は標準エラー出力に説明を表示して終了コード1を返します。リージョン未設定なら`-region`または`AWS_REGION`を指定してください。認証失敗時はプロファイルや有効期限を確認し、SSOの場合は再ログインしてください。タイムアウト時はネットワークを確認してください（既定30秒）。

## ビルド・検証

```sh
go test ./...
go vet ./...
go build -o /tmp/go-5-jungle .
/tmp/go-5-jungle -help
```

テストはローカルの模擬STSを利用し、AWS認証情報は不要です。ビルド・テストの成功だけでは実AWSへの接続成功は保証されません。有効な認証情報でCLIを実行して確認してください。

参考: [AWS SDK for Go v2の設定](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-gosdk.html)
