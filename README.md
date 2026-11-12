# kos-conformance

KOS 平台项目：契约测试套件 — 能力、Provider、扩展与应用的一致性。

kos-conformance 按发行版用例集验证能力、Provider、扩展与应用，输出逐条结果与用例集摘要。它不提供产品页面，不写能力状态：运行由运行者直接执行，结论留在报告里。

- 设计：`tsic/kos-design` 仓的 `kos-conformance/`（范围、契约、验收、决策记录）

## 模块

| 组 | 模块 | 职责 | 拥有的状态 |
|---|---|---|---|
| 能力 | `capability/` | 能力用例集：基线的卷、快照、HTTP 入口，以及能力作用域、信任与需求规则 | 无 |
| Provider | `provider/` | 注册、状态与退场用例 | 无 |
| 扩展 | `extension/` | 包、声明、对象类与控制器用例 | 无 |
| 应用 | `app/` | 清单与会话用例 | 无 |
| 验收 | `acceptance/` | 运行器与结果摘要 | 无 |

用例集模型的公共类型（组、用例、用例集、摘要）在模块根的 `conformance` 包；运行入口在 `cmd/kos-conformance`。

## 用例集

发行版用例集按四组组织。每条用例带所属组、它验证的契约条目与判据版本；没有契约归属的用例在构造时就被拒绝，因此报告里不存在未归属的行。

| 组 | 用例集 | 覆盖契约 |
|---|---|---|
| 能力 | `kos/conformance/storage-volume-v1` | KOSD-C-09 |
| 能力 | `kos/conformance/storage-snapshot-v1` | KOSD-C-09 |
| 能力 | `kos/conformance/network-http-endpoint-v1` | KOSD-C-09 |
| 能力 | `kos/conformance/capability-scope-v1` | KOSD-C-09 |
| 能力 | `kos/conformance/capability-trust-v1` | KOSD-C-10 |
| 能力 | `kos/conformance/capability-demand-v1` | KOSD-C-11 |
| Provider | `kos/conformance/provider-registration-v1` | KOSD-C-47 |
| Provider | `kos/conformance/provider-deprovision-v1` | KOSD-C-48 |
| 扩展 | `kos/conformance/extension-package-v1` | KOSD-C-17 |
| 扩展 | `kos/conformance/extension-declaration-v1` | KOSD-C-18 |
| 扩展 | `kos/conformance/extension-object-class-v1` | KOSD-C-19 |
| 扩展 | `kos/conformance/extension-controller-v1` | KOSD-C-20 |
| 应用 | `kos/conformance/app-manifest-v1` | KOSD-C-73 |
| 应用 | `kos/conformance/app-session-v1` | KOSD-C-75 |

能力组的卷、快照与 HTTP 入口三组用例集接手自 kosd 的发行版用例集，用例标识沿用。

## 用例集摘要

摘要覆盖用例集合与每条用例的判据版本，形如 `sha256:<hex>`，随报告输出，在构造用例集时算出。

- 同一用例集重复运行给出相同摘要。
- 集合变化（增删用例或用例集）或任一判据版本变化，摘要都不同。
- 描述文字不计入摘要：措辞不是判据；判据变化改判据版本。
- 用例集按标识编码，声明顺序不影响摘要；用例按执行顺序编码。

发行版把 `kos-conformance digest` 的输出作为记录，运行报告把它与实际执行的摘要并列。

## 报告与结论

一次运行输出一份报告：每条用例给出 `passed`、`failed` 或 `not-run` 之一，未通过项与未执行项分列。

- 单条用例失败不中止运行，其余用例照常给出结果。
- 用例因前置不满足而未运行（`acceptance.ErrPrerequisite`）记为未执行，不计入通过项。
- 没有注册实现的 Provider 全部记为未执行，绝不计为通过。
- 整体结论：运行摘要与发行版记录一致且全部用例通过才记 `passed`；一致但存在未通过或未执行记 `failed`；摘要不一致记 `incomparable`，报告同时给出运行摘要与发行版摘要。

## 运行入口

```sh
make build
bin/kos-conformance cases                     # 四组用例集与摘要
bin/kos-conformance digest                    # 发行版记录的用例集摘要
bin/kos-conformance run -provider <name> -distribution-digest <sha256:...> [-json]
```

`run` 只执行二进制内注册的实现，并要求给出发行版记录：没有记录的运行无法判定可比。本仓只发行套件本身，不链接任何 Provider 实现，因此这里 `run` 的逐条结果是未执行；发行版在自己的入口注册实现（`acceptance.Registry`）后，同一运行器给出真实结果。结论不是 `passed` 时退出码非零。运行入口以运行者身份执行，不写能力状态，也不提供产品页面。

## 自测

```sh
make lint        # gofmt 检查与 go vet
make test        # go test ./...
make acceptance  # tests/acceptance/CONF-TC-1..4，逐条复现设计场景
```

验收用例目录以设计编号命名，一条用例一个目录。
