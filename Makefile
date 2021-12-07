test:
	# TODO: dockerコマンドが使えなければ異常終了する処理を追加
	go test ./... -v

test-nocache:
	# TODO: dockerコマンドが使えなければ異常終了する処理を追加
	go test ./... -v -count=1