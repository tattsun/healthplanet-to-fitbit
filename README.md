# healthplanet-to-fitbit

[![dockeri.co](https://dockeri.co/image/tattsun/healthplanet-to-fitbit)](https://hub.docker.com/r/tattsun/healthplanet-to-fitbit)

[HealthPlanet](https://www.healthplanet.jp/)に登録された体重・体脂肪率情報を Fitbit へ転送する。

## 環境変数

`.env`ファイルまたは、環境変数に以下を定義する。

| 環境変数名                 | 内容                                                                          |
| -------------------------- | ----------------------------------------------------------------------------- |
| HEALTHPLANET_CLIENT_ID     | HealthPlanet 公式サイトで発行できるクライアント ID                            |
| HEALTHPLANET_CLIENT_SECRET | HealthPlanet 公式サイトで発行できるクライアントシークレット                   |
| HEALTHPLANET_ACCESS_TOKEN  | HealthPlanet のアクセストークン（`healthplanet-gettoken` を使用して取得する） |
| FITBIT_CLIENT_ID           | Fitbit 公式サイトで発行できるクライアント ID                                  |
| FITBIT_CLIENT_SECRET       | Fitbit 公式サイトで発行できるクライアントシークレット                         |
| FITBIT_ACCESS_TOKEN        | Fitbit のアクセストークン（`fitbit-gettoken` を使用して取得する）             |
| FITBIT_REFRESH_TOKEN       | Fitbit のリフレッシュトークン（`fitbit-gettoken` を使用して取得する）         |

## 事前準備

- HealthPlant, Fitbit の公式サイトから各種 API キーを取得し、環境変数に登録する。
- `fitbit-gettoken` と `healthplanet-gettoken` を使用して、各種トークンを取得し、環境変数に登録する。

## 使用方法

`healthplanet-to-fitbit` を実行する。

直近３か月の情報（体重・体脂肪率）が HeathPlanet から取得され、Fitbit へ登録される。繰り返し起動するとアクセス数の制限に引っかかる場合があるため、時間をおいて起動することを推奨する。
