# ZENV

切り替える系の環境変数が多いときに使うツール

## Usage

管理する環境変数のキーと値の一覧を表示
```sh
zenv list
```

管理する環境変数のキーを追加
```sh
zenv add [KEY]
```

環境変数[KEY]に対して選べる値の一覧を取得
今選ばれているものは*がつく
```sh
zenv set [KEY]
```

環境変数[KEY]の値に[VALUE]をセットする
値一覧にも追加する
```sh
zenv set [KEY] [VALUE]
```

環境変数[KEY]の値を空にする
```sh
zenv unset[KEY]
```

環境変数[KEY]の管理をやめる
```sh
zenv rm AWS_REGION
```
