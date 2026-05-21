import { request } from "./request"

export interface ReportData {
  content_type: string
  content_id: number
  reason: string
}

const REPORT_REASONS = [
  "垃圾广告",
  "色情低俗",
  "人身攻击",
  "违法信息",
  "虚假信息",
  "刷屏",
  "其他",
]

export function getReportReasons(): string[] {
  return REPORT_REASONS
}

export async function createReport(data: ReportData): Promise<void> {
  await request({
    method: "POST",
    url: "/reports",
    data,
  })
}
