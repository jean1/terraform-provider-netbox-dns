#!/bin/sh

# install netbox and netbox dns plugin in a vm called "nb" and start netbox
wget -O "NetBox REST API (4.2) unfixed.yaml" http://nb/api/schema/

# fix schema
sed \
	-e '/^    BriefDNSSECKeyTemplate:/,/^    [^ ]/{/- null/d}' \
	-e '/^    BriefDNSSECKeyTemplateRequest:/,/^    [^ ]/{/- null/d}' \
	-e '/^    DNSSECKeyTemplate:/,/^    [^ ]/{/- null/d}' \
	-e '/^    DNSSECKeyTemplateRequest:/,/^    [^ ]/{/- null/d}' \
	-e '/^    PatchedDNSSECKeyTemplateRequest:/,/^    [^ ]/{/- null/d}' \
	'NetBox REST API (4.2) unfixed.yaml' > 'NetBox REST API (4.2).yaml'
