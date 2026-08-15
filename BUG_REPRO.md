# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

审计文件被改成两条记录挤在同一行之后，verify 判定审计链无效，但读取审计记录的路径却照常返回两条记录、报告一切正常，两个命令对同一个文件给出互相矛盾的结论。请修复审计读取，使它和校验采用同样的一行一条记录框定，遇到不合法的行返回错误，同时不改变正常审计文件的读取结果与 limit 行为，并保证全量测试通过。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/gogo-52
- 仓库地址：https://github.com/zhanglei10281852-gif/gogo-52.git
- parent SHA：a59a140d9872d0335abca4025cf241cfc120008b

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/gogo-52.git bug-repro
cd bug-repro
git checkout --detach a59a140d9872d0335abca4025cf241cfc120008b
go test ./internal/rail -run "^TestReadAuditRejectsMultipleRecordsOnOneLine$" -count=1 -v
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/rail -run "^TestReadAuditRejectsMultipleRecordsOnOneLine$" -count=1 -v
=== RUN   TestReadAuditRejectsMultipleRecordsOnOneLine
    audit_read_regression_test.go:39: ReadAudit accepted an audit file whose line framing verification rejects
--- FAIL: TestReadAuditRejectsMultipleRecordsOnOneLine (0.01s)
FAIL
FAIL	releaserail/internal/rail	0.016s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/rail -run "^TestReadAuditRejectsMultipleRecordsOnOneLine$" -count=1 -v
=== RUN   TestReadAuditRejectsMultipleRecordsOnOneLine
    audit_read_regression_test.go:39: ReadAudit accepted an audit file whose line framing verification rejects
--- FAIL: TestReadAuditRejectsMultipleRecordsOnOneLine (0.06s)
FAIL
FAIL	releaserail/internal/rail	0.195s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

同一行出现第二条记录或空行时读取返回错误；正常审计文件的记录内容、顺序与 limit 截断行为不回归；双架构定向、全量、build/vet 通过。
