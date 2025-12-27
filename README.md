# X Post Preview Generator

GoでX(旧Twitter)の埋め込みプレビュー風カードを画像/HTML/SVGとして生成するCLIです。

![](./samples/jack.svg)

## 使い方

ビルド:

```bash
make build
```

PNG出力:

```bash
./xpostgen \
  -text "今日はCLIを作りました" \
  -name "Example User" \
  -id "example" \
  -date "10:55 AM · Dec 6, 2017" \
  -location "Tokyo" \
  -output "out.png"
```

HTML出力:

```bash
./xpostgen \
  -text "HTMLでもプレビュー" \
  -name "Example User" \
  -id "example" \
  -output "out.html"
```

SVG出力:

```bash
./xpostgen \
  -text "SVGで書き出し" \
  -name "Example User" \
  -id "example" \
  -output "out.svg"
```

CTA非表示:

```bash
./xpostgen \
  -text "CTAを非表示" \
  -name "Example User" \
  -id "example" \
  -no-cta \
  -output "out.png"
```

## Makefile

- `make build`: CLIビルド
- `make test`: テスト実行
- `make fmt`: gofmt
- `make vet`: go vet
- `make lint`: lint (go vet)
- `make tidy`: go mod tidy
- `make sample`: サンプル生成 (jack / just setting up my twttr)
- `make`: fmt/lint/test/build

## オプション

- `-text` (必須): ツイート本文
- `-name` (必須): 表示名
- `-id` (必須): ユーザーID (@なし可)
- `-icon`: アイコン画像パスまたはURL
- `-date`: 日付 (任意)
- `-location`: 現在地 (任意)
- `-cta`: CTAボタン文言 (空で非表示、既定は英語文言)
- `-no-cta`: CTAを非表示
- `-verified`: 認証バッジを表示
- `-output`: 出力ファイルパス (拡張子から形式を推定)
- `-format`: 出力形式 `png|jpg|jpeg|gif|svg|html`
- `-width`: 出力幅(px)
- `-padding`: 余白(px)
- `-theme`: `light` または `dark`
- `-font`: 本文フォントパス(.ttf/.otf)
- `-font-bold`: 太字フォントパス(.ttf/.otf)
- `-font-family`: HTML/SVG用のCSS font-family

## フォントについて

HTML/SVGではシステムフォント優先のスタックを使用します。
PNG/JPG/GIFはGo側で描画するため、必要に応じて `-font` にCJK対応フォントを指定してください。
(例: Noto Sans JP など)

## ライセンス

MIT
