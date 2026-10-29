# 将教程发布到 GitHub Pages

## 站点与源码的关系

Markdown 是维护源，MkDocs Material 把它构建为 site/ 中的静态文件，GitHub Pages 托管构建产物。Go 示例用于下载学习与本地运行，Pages 不运行后台 Go 进程。

仓库：[01reportplus/golang-tutorial](https://github.com/01reportplus/golang-tutorial)。站点配置地址：[Go 全栈学习笔记](https://01reportplus.github.io/golang-tutorial/)。首次发布是否完成以 Actions 部署结果和实际网页为准。

## 本地预览

需要 Python 3.10+ 与 Go 1.26+。Windows PowerShell：

```powershell
python -m venv .venv
.\.venv\Scripts\python.exe -m pip install -r requirements.txt
.\.venv\Scripts\python.exe -m mkdocs serve
```

macOS/Linux：

```sh
python3 -m venv .venv
.venv/bin/python -m pip install -r requirements.txt
.venv/bin/python -m mkdocs serve
```

浏览器打开终端显示的本地地址，默认是 `http://127.0.0.1:8000/`。修改正文通常自动刷新。无需激活虚拟环境，直接使用虚拟环境中的 Python 可以减少 shell 权限与 PATH 干扰。

## 第一次启用 Pages

1. 打开仓库 Settings → Pages。
2. 在 Build and deployment 的 Source 选择 **GitHub Actions**。
3. 检查 Actions 功能已开启，组织策略允许所用官方 actions。
4. 将工程提交推送至仓库默认分支。
5. 打开 Actions → Deploy GitHub Pages，等待 build 和 deploy 成功。
6. 从 deployment 输出的 page_url 打开站点，抽查文章、中文搜索、代码与手机布局。

工作流使用 contents: read 构建，部署 job 单独获得 pages: write 与 id-token: write，并通过 github-pages environment 发布。默认使用 GITHUB_TOKEN，无需在代码里设置个人 token。

## 默认分支与环境保护

Pages 工作流对 push 触发，但 build 条件只允许当前仓库 default_branch。这样既支持 main，也支持未来调整默认分支；非默认分支仍由 CI 检查。手动 workflow_dispatch 也必须选默认分支，否则 build 会跳过。

github-pages 环境可能有允许分支或审批规则。工作流判断正确但 deploy 被阻止时，应检查 Settings → Environments → github-pages；不能通过修改条件绕过组织保护。如果默认分支被改名，环境规则可能也需同步。

## 项目路径与链接

普通仓库站点在 `https://账号.github.io/仓库名/`，用户站点则在根路径。mkdocs.yml 中 site_url 已设置为本仓库项目地址，内部导航使用相对链接，避免部署在 `/golang-tutorial/` 后资源指向域名根目录。

迁移仓库、账号或自定义域名时，更新 site_url、repo_url、repo_name 和 edit_uri。自定义域名还需 DNS 与 Pages 设置，不能只改文件。部署产物由 artifact 工作流管理，不需要另外维护 gh-pages 分支。

## 发布前检查

```sh
go vet ./...
go test ./...
python scripts/check_docs.py
python scripts/check_snippets.py
python -m mkdocs build --strict
python scripts/check_site.py
```

Windows 可以将 python 替换为 `.\.venv\Scripts\python.exe`。race 和模糊测试在 CI 中执行；本地有相应工具链也可运行。站点脚本检查生成 HTML 内部链接与锚点，并验证中文搜索索引存在，不依赖公网链接永远可用。

## 失败排查

build 失败先看最早错误：缺依赖、Markdown 链接、代码块编译或严格构建警告。deploy 权限失败先确认 Source 与环境规则。网页 404 检查仓库可见性、Pages 地址、构建 artifact 与部署完成情况。内容陈旧时核对已部署 commit 和默认分支，再排除缓存。

GitHub Free 的公共仓库可使用 Pages，私有仓库可用性取决于计划；组织限制也会影响功能。当前规则以[GitHub 官方说明](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages)为准。
