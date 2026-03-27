package inference

import (
	"context"
	apisixReq "gin-vue-admin/model/apisix/request"
)

// inferenceApisixService Apisix 服务接口
type inferenceApisixService interface {
	CreateRoute(ctx context.Context, req *apisixReq.CreateRouteReq) error
	DeleteRoute(ctx context.Context, req *apisixReq.DeleteRouteReq) error
}
