package ziface

type IRouter interface {
	// 处理业务之前的钩子方法
	PreHandler(request IRequest)
	// 处理业务的主钩子方法
	Handler(request	 IRequest)
	// 处理业务之后的钩子方法
	PostHandler(request IRequest)
}
