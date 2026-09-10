package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out, errOut io.Writer) error {
	flags := flag.NewFlagSet("go-5-jungle", flag.ContinueOnError)
	flags.SetOutput(errOut)
	profile := flags.String("profile", "", "AWSプロファイル（省略時はSDKの標準設定）")
	region := flags.String("region", "", "AWSリージョン（省略時は環境変数・共有設定）")
	timeout := flags.Duration("timeout", 30*time.Second, "認証とSTS呼び出しの制限時間")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("位置引数は使えません。-helpで使い方を確認してください")
	}
	if *timeout <= 0 {
		return errors.New("-timeoutは0より大きい時間を指定してください（例: 30s）")
	}
	ctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()
	var options []func(*config.LoadOptions) error
	if *profile != "" {
		options = append(options, config.WithSharedConfigProfile(*profile))
	}
	if *region != "" {
		options = append(options, config.WithRegion(*region))
	}
	cfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return fmt.Errorf("AWS設定を読み込めません。プロファイルと共有設定ファイルを確認してください: %w", err)
	}
	if cfg.Region == "" {
		return errors.New("AWSリージョンが未設定です。-region ap-northeast-1、AWS_REGION、または共有設定のregionを指定してください")
	}
	identity, err := sts.NewFromConfig(cfg).GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("AWS接続確認がタイムアウトしました。ネットワークを確認するか-timeoutを延長してください: %w", err)
		}
		if errors.Is(err, context.Canceled) {
			return errors.New("AWS接続確認を中止しました")
		}
		return fmt.Errorf("AWS接続主体を確認できません。認証情報の設定・有効期限とネットワークを確認してください。SSOの場合はaws sso login --profile <名前>を実行してください: %w", err)
	}
	_, err = fmt.Fprintf(out, "Account: %s\nARN:     %s\nUser ID: %s\n", aws.ToString(identity.Account), aws.ToString(identity.Arn), aws.ToString(identity.UserId))
	return err
}
