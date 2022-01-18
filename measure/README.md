# measure

復元度合いを調べるためのコード

## lsexec

指定したimageによるコンテナ上で`find . -type f -ls | awk '{ print $11, $7 }'`を実行する

現在の実装では正確には`awk`はコンテナの外で行なっている

引数で対象のimageを指定

オプショナルなフラグでfindの対象ディレクトリを指定

## matchrate

`find . -type f -ls | awk '{ print $11, $7 }'`の出力結果を比較する

結果が書かれたファイルのパス2つを引数に取る