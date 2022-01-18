# measure

復元度合いを調べるためのコード

## matchrate

`find . -type f -ls | awk '{ print $11, $7 }'`の出力結果を比較する

結果が書かれたファイルのパス2つを引数に取る