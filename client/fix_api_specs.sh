#!/bin/sh

sed \
	-e '/^    BriefDNSSECKeyTemplate:/,/^    [^ ]/{/- null/d}' \
	-e '/^    BriefDNSSECKeyTemplateRequest:/,/^    [^ ]/{/- null/d}' \
	-e '/^    DNSSECKeyTemplate:/,/^    [^ ]/{/- null/d}' \
	-e '/^    DNSSECKeyTemplateRequest:/,/^    [^ ]/{/- null/d}' \
	-e '/^    PatchedDNSSECKeyTemplateRequest:/,/^    [^ ]/{/- null/d}' \
	'NetBox REST API (4.2).yaml' >new.yaml
