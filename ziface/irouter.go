package ziface

// IRouter 路由的抽象接口，路由里数据都是IRequest
type IRouter interface {
	// PreHandle 在处理conn业务之前的钩子方法hook
	PreHandle(request IRequest)
	// Handle 在处理conn业务的主方法hook
	Handle(request IRequest)
	// AfterHandle 在处理conn业务之后的钩子方法hook
	AfterHandle(request IRequest)
}
