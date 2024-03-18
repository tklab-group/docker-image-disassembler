# 開発について

## テスト

### テストの実行

dockerコマンドが使える環境が必要

```sh
# シンプルなテストの実行
make test

# キャッシュなしでのテスト実行
make test-nocache
```

### テストの更新

テストにはタグ記述を元に実際にイメージを取得・ビルドした上で処理を行うものが含まれているため、実行時期依存で失敗するものがある。

主な例:
- `cli/checkpkg/check_pkg_test.go`
- `disassembler/image/image_archive_test.go`
- `disassembler/pkginfo/apt_pkg_info_test.go`

テストの失敗内容やエラーメッセージを見て適宜テストの方を更新する。

テスト更新の例: [実行時期依存で通らなくなっていたテストを通るようにした by oribe1115 · Pull Request #14 · tklab-group/docker-image-disassembler](https://github.com/tklab-group/docker-image-disassembler/pull/14)

#### 対象イメージが取得できない

Docker Hubリポジトリ上で該当のタグが付与されたイメージがなくなったなどの理由で発生する。

同じリポジトリ上の、あまり付与対象が変更されなさそうなタグを選んで置き換える。
[ubuntu](https://hub.docker.com/_/ubuntu)であれば`{バージョン名}-{ビルド日時}`の形式のタグ(e.g. `mantic-20231011`)は基本的に付与対象が変更されることはないはず。

#### 対象パッケージが取得できない

ベースイメージの変更や、パッケージレジストリ上での配信停止などが理由で発生する。

対象パッケージの現時点で取得できるバージョンに手動で置き換える。

#### 取得されるパッケージバージョンが違う

ベースイメージの変更や、パッケージレジストリ上での更新などが理由で発生する。

多くのテストは[sebdah/goldie](https://github.com/sebdah/goldie)によるGoledn Testを採用しているため`make test-update`で一括更新ができる。
更新結果が形式的に問題がなければそのまま採用。

その他のテストについては、対象パッケージの現時点で取得できるバージョンに手動で置き換える。