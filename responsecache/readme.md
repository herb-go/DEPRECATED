# ResponseCache 页面响应缓存组件

提供将整个http响应进行缓存功能的组件

## 功能
* 缓存后的页面能够继续添加或者维护Header
* 通过缓存的Loader功能防止缓存雪崩
* 能状态码进行判断页面是否需要需要缓存

## 上下文

上下文指在一次请求中与缓存相关的内容。通过维护上下文来决定缓存组件的行为

上下文结构为

	http.ResponseWriter
    响应的写入数据，一般用于判断响应状态
	Request    *http.Request
    当前请求
	Identifier func(*http.Request) string
    请求识别器。由ContextBuilder设置。返回值为空字符串时，不缓存，否则以返回值为主键进行缓存
	TTL        time.Duration
    缓存声明周期。由ContextBuilder设置。
	Buffer     []byte
    缓存写入的数据
	StatusCode int
    当前响应的状态码。一般用于Validator判断
	Validator  func(*Context) bool
    响应验证器。由ContextBuilder设置。返回值为true时才进行缓存。为空时使用DefaultValidator
	Cache      cache.Cacheable
    使用的缓存组件，由ContextBuilder设置

## 上下文构建器 ContextBuilder

上下文构建器负责设置维护上下文

上下文中的 Identifier,TTL,Validator,Cache由构建器负责构建

### 简单上下文构建器 PlainContextBuilder

简单上下文构建器通过设置固定的Identifier,TTL,Validator,Cache来进行使用，是一种基础的上下文构建器

### 参数上下文构建器 ParamContextBuilder

#### 常用参数

## 使用方式
