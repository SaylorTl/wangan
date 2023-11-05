package Constant

// 任务进行到哪一步了 默认0 没开始或者是资产扫描未完成 1:正在同步数据 2:正在打标签 3:正在漏洞扫描 4:扫描完成
const STEP_SCANNING = 0
const STEP_SYNC_DATA = 1
const STEP_TAGGING = 2
const STEP_LOOPHOLE = 3
const STEP_FINISHED = 4

// 扫描IP类型
const IP_TYPE_V4 int = 0

const IP_TYPE_V6 int = 1

const IP_TYPE_DOMAIN int = 2
