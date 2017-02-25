# MusicBox

## 整合了网易,QQ音乐,虾米的音乐API服务，支持境外客户端使用。

FrontEnd:https://github.com/acgotaku/MusicPlayer

#### 搜索接口API

GET `/api/search?keyword=<search keyword>&limit=<number per page>&page=<page number start from 1>&source=<data source: netease qq xiami>`

Params
- keyword: default Nome
- limit: default 20, number per page
- page: default 1, current page
- source: default 'netease', can use 'qq' or 'xiami'

Return:

status: 200
body:
```json
{
  "code": 200,
  "data": [
    {
      "id": "593040",
      "name": "tsubasa",
      "artist": "梶浦由記",
      "album": "NHKアニメーション ｢ツバサ・クロニクル｣ オリジナルサウンドトラック Future Soundscape I"
    },
    {
      ...
    }
  ],
  "total": 629
}
```

#### 音乐接口API

GET `/api/track?id=<music track id>&source=<data source: netease qq xiami>&country=<your country>`

Params
- keyword: default Nome
- source: default 'netease', can use 'qq' or 'xiami'
- country: default 'china', can use 'japan' or 'us'


Return:

status: 200
body:
```json
{
  "code": 200,
  "mp3Url": "http://m10.music.126.net/7256/bcdd/9a8dd536a7cdf84abc1cc6a4a833cab4.mp3"
}
```
