# 校园闲置物品交换系统

本项目是课程设计"校园闲置物品交换系统"的完整前后端分离演示版，选题避开商城、外卖、教学、课程、酒店、图书等禁用题目，围绕校园闲置物品循环利用、技能互助、信用评价与安全面交构建完整应用。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Go 1.21 + Gin + JWT + JSON 文件存储 |
| 前端 | Vue 3 + Vite + Pinia + Vue Router + Axios |
| 认证 | JWT Bearer Token |
| 通信 | RESTful API |
| 版本控制 | Git |

## 项目结构

```
├── server/                  # Go 后端 (端口 8080)
│   ├── cmd/main.go          # 入口文件，路由注册
│   ├── internal/
│   │   ├── config/          # 配置常量 (JWT密钥、端口)
│   │   ├── model/           # 数据模型 (8个结构体)
│   │   ├── service/         # 数据持久化 (JSON读写+种子数据)
│   │   ├── middleware/      # JWT认证中间件 + CORS中间件
│   │   └── handler/         # 业务处理 (auth/item/exchange/message/favorite/misc)
│   └── data/db.json         # 运行时自动生成的数据文件
│
├── client/                  # Vue 3 前端 (端口 3000)
│   ├── src/
│   │   ├── views/           # 10个页面组件
│   │   ├── components/      # 4个公共组件
│   │   ├── stores/          # Pinia 状态管理
│   │   ├── api/             # Axios API 层 (6个模块)
│   │   └── router/          # Vue Router 路由配置
│   └── vite.config.js       # Vite配置 (含API代理)
│
├── database/schema.sql      # SQL建表参考
└── docs/                    # 文档 (分工、测试用例)
```

## 运行方式

### 前置要求

- Go >= 1.21
- Node.js >= 18
- npm >= 9

### 启动后端

```bash
cd server
go run cmd/main.go
```

后端启动后监听 `http://localhost:8080`，首次运行自动生成种子数据。

### 启动前端

```bash
cd client
npm install
npm run dev
```

前端启动后访问 `http://localhost:3000`，Vite 自动将 `/api` 请求代理到后端 8080 端口。

### 生产构建

```bash
cd client
npm run build
```

构建产物在 `client/dist/` 目录。

## 演示账号

| 角色 | 用户名 | 密码 |
| --- | --- | --- |
| 学生 | 202401001 | 123456 |
| 学生 | 202401002 | 123456 |
| 学生 | 202401003 | 123456 |
| 管理员 | admin | admin123 |

## 功能模块（10个）

1. **用户认证**：登录/注册/JWT认证，学生与管理员角色区分。
2. **闲置物品发布**：发布物品或技能服务，填写分类、校区、成色、期望交换内容。
3. **闲置广场检索**：支持关键词、分类、校区筛选，瀑布流卡片布局。
4. **交换申请管理**：发起交换申请，物主可同意/拒绝，双方可完成交换。
5. **站内消息沟通**：围绕交换申请的实时聊天，气泡样式，保留完整记录。
6. **收藏与智能推荐**：根据收藏分类偏好进行个性化推荐。
7. **信用评价**：交换完成后双方评价，1-5分评分。
8. **举报与审核**：发现违规物品/用户可举报，管理员后台处理。
9. **后台运营看板**：用户数、物品数、交换数、热门分类统计。
10. **移动端适配**：响应式布局，适配手机浏览器。

## 数据表设计

`database/schema.sql` 共设计 11 张表：

- `users` 用户表
- `categories` 分类表
- `items` 物品表
- `item_tags` 标签表
- `exchanges` 交换申请表
- `messages` 消息表
- `favorites` 收藏表
- `reports` 举报表
- `ratings` 评价表
- `notifications` 通知表
- `audit_logs` 审计日志表

## 技术特色与创新点

- **前后端分离**：Go Gin 后端 + Vue 3 前端，RESTful API 通信
- **JWT 认证**：无状态 Token 认证，支持中间件鉴权
- **响应式设计**：PC/平板/手机多端适配
- **创新场景**：物品换物 + 技能互助、校园信用分、安全面交提示、收藏驱动推荐、低碳校园循环理念

## 五人分工（远程协作，每人独立模块）

| 成员 | 负责模块 | 具体工作 |
| --- | --- | --- |
| 成员A | 用户认证 + 物品发布 + 项目管理 | Go后端auth/item接口、Login.vue、Publish.vue、架构搭建、代码合并 |
| 成员B | 交换流程 + 站内消息 | Go后端exchange/message接口、Exchanges.vue、Messages.vue |
| 成员C | 首页广场 + 物品详情 + 收藏 | Go后端item查询/favorite接口、Home.vue、ItemDetail.vue、Favorites.vue |
| 成员D | 数据库设计 + 文档 + UI | schema.sql、课程报告、UML图、Profile.vue、CSS样式调整 |
| 成员E | 推荐+评价+举报+后台看板 | Go后端misc接口、Recommend.vue、Admin.vue、测试用例、答辩PPT |
