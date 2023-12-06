## 项目启动

确保你安装了 `pnpm`：

```bash
npm install -g pnpm
```

依赖安装：

```bash
pnpm install
```

开发环境启动：

```bash
pnpm run dev
```

生产环境构建：

```bash
pnpm run build
```

可以通过 `pnpm run preview` 来预览生产环境构建结果。

打包产物位于 `dist` 目录下。

## 配置更改

开发环境域名配置位于 `vite.config.ts` 中的 `proxy`。

## 部署

项目使用客户端路由，因此需要后端配置 `index.html` 的路由规则。