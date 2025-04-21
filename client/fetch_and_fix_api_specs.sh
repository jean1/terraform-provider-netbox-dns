#!/bin/sh

# install netbox and netbox dns plugin in a vm called "nb" and start netbox
# curl -H "Accept: application/json" http://nb/api/schema/ >openapi.json
[ -f openapi_broken.json ] || \
  curl -H "Accept: application/json" http://nb/api/schema/ >openapi_broken.json

# sed \
# 	-e '/^    BriefDNSSECKeyTemplate:/,/^    [^ ]/{/- null/d}' \
# 	-e '/^    BriefDNSSECKeyTemplateRequest:/,/^    [^ ]/{/- null/d}' \
# 	-e '/^    DNSSECKeyTemplate:/,/^    [^ ]/{/- null/d}' \
# 	-e '/^    DNSSECKeyTemplateRequest:/,/^    [^ ]/{/- null/d}' \
# 	-e '/^    PatchedDNSSECKeyTemplateRequest:/,/^    [^ ]/{/- null/d}' \
# 	'NetBox REST API (4.2) unfixed.yaml' > 'NetBox REST API (4.2).yaml'
# for i in BriefDNSSECKeyTemplate BriefDNSSECKeyTemplateRequest DNSSECKeyTemplate DNSSECKeyTemplateRequest PatchedDNSSECKeyTemplateRequest
# do
#   echo "/^            \"$i\": \{/,/^            \}/s/,\n\s+- null//"
# done > /tmp/mysedscript
# sed -zE -f /tmp/mysedscript openapi_broken.json > openapi.json
# rm -f /tmp/mysedscript
# 
# With jq:
#cat openapi_broken.json | jq \
#	'del(.components.schemas.BriefDNSSECKeyTemplate.properties.key_size.enum[] | select(. == null))' \
#	>/tmp/out
# jq .components.schemas.BriefDNSSECKeyTemplate.properties.key_size.enum < /tmp/out

# Each pass applies delete to each problematic API object
cp openapi_broken.json openapi_src.json
for object in BriefDNSSECKeyTemplate BriefDNSSECKeyTemplateRequest DNSSECKeyTemplate DNSSECKeyTemplateRequest PatchedDNSSECKeyTemplateRequest
do
  echo "DEBUG: pass $object" >&2
  jq < openapi_src.json \
	'del(.components.schemas.'$object'.properties.key_size.enum[] | select(. == null))' \
	> openapi_dst.json
  mv openapi_dst.json openapi_src.json
done
mv openapi_src.json openapi.json
for object in BriefDNSSECKeyTemplate BriefDNSSECKeyTemplateRequest DNSSECKeyTemplate DNSSECKeyTemplateRequest PatchedDNSSECKeyTemplateRequest
do
   echo "CHECK: $object" >&2
  jq .components.schemas.$object.properties.key_size.enum < openapi.json
done
echo "FINAL FILE: openapi.json" >&2
