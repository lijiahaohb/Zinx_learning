package znet

import "zinx/ziface"

type BaseRouter struct {
}


// 这里之所以BaseRouter的方法都为空是因为有的Router不希望有 PreHandler 或者 PostHandler这两个业务
// 之后所有的Router全部继承BaseRouter的好处就是，不需要全部实现 PreHandler 或者 PostHandler
func (b *BaseRouter) PreHandler(request ziface.IRequest) {

}
func (b *BaseRouter) Handler(request ziface.IRequest) {

}
func (b *BaseRouter) PostHandler(request ziface.IRequest) {

}
