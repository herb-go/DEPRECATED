# Member 用户系统

member 用户系统是一套基于 github.com/herb-go/herb 的用户系统。
系统主要是定义了用户系统的常用操作，用户驱动的接口，以相应的缓存操作。

## 功能

- 用户的状态/帐号/密码/令牌/权限/档案的驱动支持。可以混合使用 sql/ldap/第三方方案。
- 与缓存系统的良好配合
- 登录/登出/跳转登录等中间件的配合

## 依赖

- github.com/herb-go/usersystem/user 基本用户模块
- github.com/herb-go/deprecated/httpuser HTTP 用户模块
- github.com/herb-go/user/role 用户权限
- github.com/herb-go/user/role/roleservice 权限服务模块
- github.com/herb-go/herb/cache 缓存
- github.com/herb-go/session 会话模块
