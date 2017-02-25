# MusicBox

#### 搜索接口API
GET `/api/search?keyword=<search keyword>&limit=<number per page>&page=<page number start from 1>&source=<data source: netease qq xiami>`

Params
- keyword: default Nome
- limit: default 20, number per page
- page: default 1, current page
- source: default 'netease', can use 'qq' or 'xiami'
- country: default 'china', can use 'japan' or 'us'
Return:
