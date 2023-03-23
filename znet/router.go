package znet

import "v1/ziface"

// BaseRouter 实现router时候，先嵌入这个BaseRouter，然后根据需要堆个基类的方法进行重写
// 为什么不直接实现那个接口呢？
/*
	这个地方理解是这样的，如果后面的类实现这个接口会有一个比较大的问题，
	就是它必须要实现全部的接口，有些Router不希望有PreHandle、AfterHandle，这样是没有必要的
	因此采用嵌套的方式来重写方法就可以很好得避免这个问题
*/
type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) AfterHandle(request ziface.IRequest) {}
