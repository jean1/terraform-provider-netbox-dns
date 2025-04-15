#!/bin/bash
go get github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config config.yaml openapi.json

# fix types for SOA fields
## SoaSerial: int64→int32
sed -i -E '/^\s+SoaSerial\s+\*int64/s/int64/int32/' client.gen.go

## other Soa fields: int→int32
for field in DefaultTtl SoaTtl SoaRefresh SoaRetry SoaExpire SoaMinimum ; do
  sed -i -E '/^\s+'$field'\s+\*int/s/int/int32/' client.gen.go
done

# fix type for nested view and nested mname (primary namerserver) in zone request struct
sed -i -E -e '/^type WritableZoneRequest struct \{/,/^\}/s/View\s+\*WritableZoneRequest_View/View *int/' \
	client.gen.go

sed -i -E -e '/^type WritableZoneRequest struct \{/,/^\}/s/SoaMname\s+\*WritableZoneRequest_SoaMname/SoaMname *int/' \
	client.gen.go
