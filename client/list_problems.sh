#!/bin/bash
for i in BriefDNSSECKeyTemplate BriefDNSSECKeyTemplateRequest DNSSECKeyTemplate DNSSECKeyTemplateRequest PatchedDNSSECKeyTemplateRequest
do
  jq .components.schemas.$i.properties.key_size.enum < openapi_broken.json
done
