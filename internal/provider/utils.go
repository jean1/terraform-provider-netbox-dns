package provider

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jean1/terraform-provider-netbox-dns/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func toIntPointer(from *int64) *int {
	if from == nil {
		return nil
	}
	val := int(*from)
	return &val
}

func int32toIntPointer(from *int32) *int {
	if from == nil {
		return nil
	}
	val := int(*from)
	return &val
}

func booltoBoolPointer(from *bool) *bool {
	if from == nil {
		return nil
	}
	val := bool(*from)
	return &val
}

func toIntListPointer(ctx context.Context, from types.List) ([]int, diag.Diagnostics) {
	var values []int64
	diags := from.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]int, 0, len(values))
	for _, val := range values {
		out = append(out, int(val))
	}
	return out, diags
}

func maybeStringValue(in *string) types.String {
	if in == nil {
		return types.StringNull()
	}
	if *in == "" {
		return types.StringNull()
	}
	return types.StringPointerValue(in)
}

func maybeInt64Value(in *int) types.Int64 {
	if in == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*in))
}

func fromInt64Value(in types.Int64) *int {
	if in.IsNull() {
		return nil
	}
	return toIntPointer(in.ValueInt64Pointer())
}

func maybeInt32Value(in *int32) types.Int32 {
	if in == nil {
		return types.Int32Null()
	}
	return types.Int32Value(int32(*in))
}

func fromInt32Value(in types.Int32) *int {
	if in.IsNull() {
		return nil
	}
	return int32toIntPointer(in.ValueInt32Pointer())
}


func maybeBoolValue(in *bool) types.Bool {
	if in == nil {
		return types.BoolNull()
	}
	if *in {
		return types.BoolValue(true)
	} else {
		return types.BoolValue(false)
	}
}

func fromBoolValue(in types.Bool) *bool {
	if in.IsNull() {
		return nil
	}
	return booltoBoolPointer(in.ValueBoolPointer())
}

func httpError(res *http.Response, body []byte) string {
	return fmt.Sprintf("Bad response: Status %d with content type \"%s\"\n%s", res.StatusCode, res.Header.Get("Content-Type"), string(body))
}

func importByInt64ID(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "ID to import must be a number")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func appendPointerSlice[T any](s *[]T, vals ...T) *[]T {
	if s == nil {
		val := make([]T, 0, len(vals))
		s = &val
	}
	newS := append(*s, vals...)
	return &newS
}

func doPlainReq(ctx context.Context, req *http.Request, c *client.Client) (*http.Response, error) {
	req = req.WithContext(ctx)
	for _, e := range c.RequestEditors {
		if err := e(ctx, req); err != nil {
			return nil, err
		}
	}

	return c.Client.Do(req)
}
