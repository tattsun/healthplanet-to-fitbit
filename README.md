# healthplanet-to-fitbit

[HealthPlanet](https://www.healthplanet.jp/)に登録された体重・体脂肪率情報を Fitbit へ転送します。

## 環境変数

`.env`ファイルまたは、環境変数に以下を定義してください。

| 環境変数名                 | 内容                                                                          |
| -------------------------- | ----------------------------------------------------------------------------- |
| HEALTHPLANET_CLIENT_ID     | HealthPlanet 公式サイトで発行できるクライアント ID                            |
| HEALTHPLANET_CLIENT_SECRET | HealthPlanet 公式サイトで発行できるクライアントシークレット                   |
| HEALTHPLANET_ACCESS_TOKEN  | HealthPlanet のアクセストークン（`healthplanet-gettoken` を使用して取得する） |
