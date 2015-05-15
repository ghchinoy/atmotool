
# Community Manager API calls

## Upload `custom.less`

```
curl 'http://ent.akana-dev.net:9900/resources/theme/default/less?unpack=false&wrapInHTML=true' 
	-H 'Pragma: no-cache' 
	-H 'Origin: http://ent.akana-dev.net:9900' 
	-H 'Accept-Encoding: gzip, deflate' 
	-H 'Accept-Language: en-US,en;q=0.8' 
	-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.37 Safari/537.36' 
	-H 'Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryhEzDilJAZUUQwpoM' 
	-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' 
	-H 'Cache-Control: no-cache' 
	-H 'Referer: http://ent.akana-dev.net:9900/atmosphere/' 
	-H 'Cookie: AtmoAuthToken_enterpriseapi=TokenID%3Dbc89fd65-ead9-11e4-ba51-a820a7104c65%2Cclaimed_id%3Durn%3Aatmosphere%3Auser%3Aenterpriseapi%3A049a10ba-64b5-4940-bb52-7826f50971c6%2CissueTime%3D1429918162016%2CexpirationTime%3D1429919961991%2CUserName%3DPlatform+Admin%2CUserFDN%3D049a10ba-64b5-4940-bb52-7826f50971c6%252Eenterpriseapi%2CAttributesIncluded%3Dfalse%2Csig%3DaxGSegvMTig7nP0I1LLrkc8zzM6MP81m3HOvdP7Q18D1Of8FN4Ww_raoVcileRA1P40R2dUxhbiRNk6ve511LHYZDyp1JFotfyHGGw1S-JreU42hxUOeaVxMK40uEdNwKkRf3SVA65p5sW9eq6GoKmySYNT1qV7buZZNw3pbqmWsypZGboyIAjStc_NBwMnHezht-JmvmPEetsD1_JpPnw0AoYZhWIdKirdqMA3Ktiq8kB-bIJ4Qtjmu325Ey7r3B5sc0-hNqsYXEy5taiuFvwXREbeyjZp75o5_OVC47-XIE0-1OKU7wD85c8jAz8IguMcc6Afw7fkY52NEa3MyjA' 
	-H 'Connection: keep-alive' 
	--data-binary $'------WebKitFormBoundaryhEzDilJAZUUQwpoM\r\nContent-Disposition: form-data; name="File"; filename="custom.less"\r\nContent-Type: application/octet-stream\r\n\r\n\r\n------WebKitFormBoundaryhEzDilJAZUUQwpoM--\r\n' 
	--compressed
```

## Regenerate Styles

```
curl 'http://ent.akana-dev.net:9900/resources/branding/generatestyles' -H 'Accept: application/json' -H 'Referer: http://ent.akana-dev.net:9900/atmosphere/' -H 'Origin: http://ent.akana-dev.net:9900' 
	-H 'X-Requested-With: XMLHttpRequest' 
	-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.37 Safari/537.36' 
	-H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' 
	--data 'theme=default' 
	--compressed
```

## Upload to `/resources/theme/default`

```
curl 'http://ent.akana-dev.net:9900/resources/theme/default?unpack=true&wrapInHTML=true' 
	-H 'Pragma: no-cache' 
	-H 'Origin: http://ent.akana-dev.net:9900' 
	-H 'Accept-Encoding: gzip, deflate' 
	-H 'Accept-Language: en-US,en;q=0.8' 
	-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.37 Safari/537.36' 
	-H 'Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryQQPruABYBUmeoA4P' 
	-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' 
	-H 'Cache-Control: no-cache' -H 'Referer: http://ent.akana-dev.net:9900/ui/apps/atmosphere/_VSPUvNSKEu9CFztvkjorQcg/resources/console/widgets/CKDocumentWidget/filemanager/index.html?dynamic=true&direct=true&docContext=/resources&docBase=/theme/default/&tocEnabled=false&tocTitle=Show%20in%20TOC&displayDefault=false&hasEditor=false&hasDownload=true&rootTOCDir=/resources' 
	-H 'Cookie: AtmoAuthToken_enterpriseapi=TokenID%3Dbc89fd65-ead9-11e4-ba51-a820a7104c65%2Cclaimed_id%3Durn%3Aatmosphere%3Auser%3Aenterpriseapi%3A049a10ba-64b5-4940-bb52-7826f50971c6%2CissueTime%3D1429918162016%2CexpirationTime%3D1429919961991%2CUserName%3DPlatform+Admin%2CUserFDN%3D049a10ba-64b5-4940-bb52-7826f50971c6%252Eenterpriseapi%2CAttributesIncluded%3Dfalse%2Csig%3DaxGSegvMTig7nP0I1LLrkc8zzM6MP81m3HOvdP7Q18D1Of8FN4Ww_raoVcileRA1P40R2dUxhbiRNk6ve511LHYZDyp1JFotfyHGGw1S-JreU42hxUOeaVxMK40uEdNwKkRf3SVA65p5sW9eq6GoKmySYNT1qV7buZZNw3pbqmWsypZGboyIAjStc_NBwMnHezht-JmvmPEetsD1_JpPnw0AoYZhWIdKirdqMA3Ktiq8kB-bIJ4Qtjmu325Ey7r3B5sc0-hNqsYXEy5taiuFvwXREbeyjZp75o5_OVC47-XIE0-1OKU7wD85c8jAz8IguMcc6Afw7fkY52NEa3MyjA' -H 'Connection: keep-alive' 
	--data-binary $'------WebKitFormBoundaryQQPruABYBUmeoA4P\r\nContent-Disposition: form-data; name="File"; filename="marriott_resourcesThemeDefault.zip"\r\nContent-Type: application/zip\r\n\r\n\r\n------WebKitFormBoundaryQQPruABYBUmeoA4P--\r\n' --compressed
```

## Upload to `/content/home/landing`

```
curl 'http://ent.akana-dev.net:9900/content/home/landing?unpack=true&wrapInHTML=true' -H 'Pragma: no-cache' -H 'Origin: http://ent.akana-dev.net:9900' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: en-US,en;q=0.8' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.37 Safari/537.36' -H 'Content-Type: multipart/form-data; boundary=----WebKitFormBoundarysUu5Y5fPg31JfVik' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' -H 'Cache-Control: no-cache' -H 'Referer: http://ent.akana-dev.net:9900/ui/apps/atmosphere/_VSPUvNSKEu9CFztvkjorQcg/resources/console/widgets/CKDocumentWidget/filemanager/index.html?dynamic=true&direct=true&docBase=/home/landing&docContext=/content&tocEnabled=false&tocTitle=Show%20in%20TOC&displayDefault=false&hasEditor=false&hasDownload=true' -H 'Cookie: AtmoAuthToken_enterpriseapi=TokenID%3Dbc89fd65-ead9-11e4-ba51-a820a7104c65%2Cclaimed_id%3Durn%3Aatmosphere%3Auser%3Aenterpriseapi%3A049a10ba-64b5-4940-bb52-7826f50971c6%2CissueTime%3D1429918162016%2CexpirationTime%3D1429919961991%2CUserName%3DPlatform+Admin%2CUserFDN%3D049a10ba-64b5-4940-bb52-7826f50971c6%252Eenterpriseapi%2CAttributesIncluded%3Dfalse%2Csig%3DaxGSegvMTig7nP0I1LLrkc8zzM6MP81m3HOvdP7Q18D1Of8FN4Ww_raoVcileRA1P40R2dUxhbiRNk6ve511LHYZDyp1JFotfyHGGw1S-JreU42hxUOeaVxMK40uEdNwKkRf3SVA65p5sW9eq6GoKmySYNT1qV7buZZNw3pbqmWsypZGboyIAjStc_NBwMnHezht-JmvmPEetsD1_JpPnw0AoYZhWIdKirdqMA3Ktiq8kB-bIJ4Qtjmu325Ey7r3B5sc0-hNqsYXEy5taiuFvwXREbeyjZp75o5_OVC47-XIE0-1OKU7wD85c8jAz8IguMcc6Afw7fkY52NEa3MyjA' -H 'Connection: keep-alive' --data-binary $'------WebKitFormBoundarysUu5Y5fPg31JfVik\r\nContent-Disposition: form-data; name="File"; filename="marriott_contentHomeLanding.zip"\r\nContent-Type: application/zip\r\n\r\n\r\n------WebKitFormBoundarysUu5Y5fPg31JfVik--\r\n' --compressed
```

# Test calls


## Less

```
 go run atmosphere.go upload less ~/dev/Akana/prospects/cm_custom_starter/custom.less --config local.conf

 atmotool upload less custom.less --config local.conf
 ```

## Single file
```
go run atmosphere.go upload file --path /content/home/landing --config local.conf ~/dev/Akana/prospects/cm_custom_starter/starter_contentHomeLanding.zip

upload file --path /content/home/landing --config local.conf starter_contentHomeLanding.zip
```
