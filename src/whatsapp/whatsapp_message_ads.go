package whatsapp

/*

"externalAdReply": {
	"title": "NÃO PERCA ESSA OPORTUNIDADE DE 10 REAIS",
	"body": "🎉 Não perca! GOHTUBE oferece 15 dias por R$6 ou 30 dias por apenas R$10. Mergulhe nas melhores histórias de filmes e séries. ➡ Clique agora!",
	"mediaType": "VIDEO",
	"thumbnailUrl": "https://scontent.xx.fbcdn.net/v/t15.5256-10/438055601_457003190341135_2845626957696104170_n.jpg?stp=dst-jpg_p180x540&_nc_cat=107&ccb=1-7&_nc_sid=63bd09&_nc_ohc=wLCyJwHE-MMQ7kNvgG7f8cm&_nc_ad=z-m&_nc_cid=0&_nc_ht=scontent.xx&_nc_gid=AtMgyL6ySP0hn-lgBXJMLPC&oh=00_AYDKLXXbIPXTNHXRxk8IxDncB45Mrty8IoztmAkychWXSA&oe=66FCD3B1",
	"mediaUrl": "https://www.facebook.com/expinformativa/videos/429637216365666/",
	"thumbnail": "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",
	"sourceType": "ad",
	"sourceId": "120212732498440753",
	"sourceUrl": "https://fb.me/4xR0JVZ3B",
	"containsAutoReply": false,
	"renderLargerThumbnail": true,
	"showAdAttribution": true,
	"ctwaClid": "ARACGmmFDy48IyrcMV2igjacnY39gcMjcUUlc_WBbfSBaoJU5Kxb7WjxI7j03Dx8Vtkf-o7D8sxvM1eZcVEEa6a7eN4rfUWtMXcU_TlGYA7jHtruvMaXGovdYqobnxLL_uN4kQ2NGA"
}

{
   "text":"Olá! Tenho interesse e queria mais informações, por favor.",
   "contextInfo":{
      "conversionSource":"FB_Post",
      "externalAdReply":{
         "title":"Garanta seu Desconto na Blue Friday",
         "body":"",
         "mediaType":1,
         "thumbnailURL":"https://scontent.xx.fbcdn.net/v/t45.1600-4/466893297_548105151336486_5338841918225127393_n.png?stp=c3.3.300.300a_dst-png_p306x306\\u0026_nc_cat=103\\u0026ccb=1-7\\u0026_nc_sid=e37a05\\u0026_nc_ohc=MNPayzx84dIQ7kNvgEJX-4K\\u0026_nc_ad=z-m\\u0026_nc_cid=0\\u0026_nc_zt=1\\u0026_nc_ht=scontent.xx\\u0026_nc_gid=AesiRDy7lFZhDMfXXTmGzHI\\u0026oh=00_AYCRuOfNo9T3TINWakizgWf-3SbK4I0hUM8ZnOetHWDsVQ\\u0026oe=675E4EAA",
         "thumbnail":"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",
         "sourceType":"post",
         "sourceID":"122144079350331080",
         "sourceURL":"https://fb.me/4zHlqO9J7",
         "containsAutoReply":false,
         "sourceApp":"facebook",
         "originalImageURL":"https://www.facebook.com/ads/image/?d=AQK6A0iwVIT_NoFYfM-vgIvgBo6mtkQ_RJRBhhu_z7FUEmdOPgDk6To0ZYWBnMp69zuvvrFzUAEEQ5tWFp6A9SHtbntYK7L441m368mBkg6WCj_EFXpHRgUd8riM8_lM0Z1V2ef4v3vYJfxp8yqMl9Xk"
      },
      "entryPointConversionSource":"post_cta",
      "entryPointConversionApp":"facebook",
      "entryPointConversionExternalSource":"FB_Post",
      "entryPointConversionExternalMedium":"unavailable"
   }
}


*/

type WhatsappMessageAds struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	SourceId  string `json:"sourceid,omitempty"`
	SourceUrl string `json:"sourceurl,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"` // base64 thumbnail
	App       string `json:"app,omitempty"`
	Type      string `json:"type,omitempty"`
}
