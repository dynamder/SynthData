# SynthData

基于大语言模型（LLM）的合成数据集生成工具。

[English Version](./README_en.md)

## 概述

SynthData 通过描述文件使用 LLM 生成合成数据集。支持多种数据格式（JSON、CSV），支持大规模批量生成，并提供交互式 CLI 向导。

## 功能特性

- **LLM 驱动生成**：使用 OpenAI 兼容 API 生成合成数据
- **多格式输出**：支持 JSON 和 CSV 输出
- **大规模生成**：支持批量处理和并发控制
- **交互模式**：CLI 向导引导完成数据生成
- **Schema 验证**：生成前验证描述文件

## 安装

```bash
git clone https://github.com/dynamder/synthdata.git
cd synthdata
go build -o synthdata ./cmd/synthdata
```

## 配置

创建配置文件（例如 `configs/default.toml`）：

```toml
[llm]
api_key = "your-api-key"
base_url = "https://api.openai.com/v1"
model = "gpt-4o-mini"
max_retries = 3
```

支持的 LLM 提供商：OpenAI、SiliconFlow、Azure OpenAI 及任何兼容 OpenAI 的 API。

## 快速开始

```bash
synthdata generate -d description.md -o output.json -s 100
```

## 使用方法

### 命令选项

| 参数 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| `--description` | `-d` | 描述文件路径 | （必填） |
| `--output` | `-o` | 输出文件路径 | （必填） |
| `--format` | `-f` | 输出格式 (json, csv) | json |
| `--scale` | `-s` | 生成记录数 | 10 |
| `--config` | `-c` | 配置文件路径 | configs/default.toml |
| `--batch-size` | | 每批记录数 | 10 |
| `--concurrency` | | 最大并发 LLM 调用数 | 5 |
| `--max-retries` | | 最大重试次数 | 3 |
| `--force` | | 覆盖已存在的输出文件 | false |
| `--verbose` | `-v` | 启用详细日志 | false |
| `--interactive` | `-i` | 启用交互式向导 | false |

### 交互模式

启动交互式向导：

```bash
synthdata generate --interactive
```

### 大规模生成

使用批量处理生成大规模数据集：

```bash
synthdata generate -d description.md -o output.json -s 10000 --batch-size 100 --concurrency 10
```

## 描述文件格式

描述文件定义数据结构：

```json
{
  "name": "数据集名称",
  "description_file": "description.md",
  "format": "json",
  "count": 100,
  "schema": {
    "name": "table_name",
    "type": "nested",
    "children": [
      { "name": "id", "type": "integer" },
      { "name": "username", "type": "string" },
      { "name": "email", "type": "string" }
    ]
  }
}
```

更多示例请参考 `examples/bilibili_chat_description/bilibili_chat.json`。

## 示例

生成 JSON 输出：
```bash
synthdata generate -d examples/bilibili_chat_description/bilibili_chat.json -o output.json -s 50
```

生成 CSV 输出：
```bash
synthdata generate -d description.md -o output.csv -f csv -s 100
```

使用自定义配置：
```bash
synthdata generate -d description.md -o output.json -c my_config.toml
```
