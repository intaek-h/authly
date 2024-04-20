package templates

import "context"

// TODO: context 키 상수로 관리해야 함.

func GetPageTitle(ctx context.Context) string {
	title, ok := ctx.Value("head-title").(string)
	if !ok {
		// TODO: 사이트 기본 이름 넣기
		title = ""
	}

	return title
}

func IsProduction(ctx context.Context) bool {
	isPrd, ok := ctx.Value("isProduction").(bool)
	if !ok {
		isPrd = false
	}

	return isPrd
}

func IsHxRequest(ctx context.Context) bool {
	hxRequest, ok := ctx.Value("isHxRequest").(bool)
	if !ok {
		return false
	}
	return hxRequest
}
