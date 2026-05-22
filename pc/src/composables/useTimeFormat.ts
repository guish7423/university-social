import dayjs from "dayjs"
import relativeTime from "dayjs/plugin/relativeTime"
import "dayjs/locale/zh-cn"

dayjs.extend(relativeTime)
dayjs.locale("zh-cn")

export function useTimeFormat() {
  function formatTime(t: string): string {
    const d = dayjs(t)
    if (dayjs().diff(d, "day") < 1) return d.fromNow()
    return d.format("MM-DD HH:mm")
  }

  function formatDate(t: string): string {
    return dayjs(t).format("YYYY-MM-DD")
  }

  function formatDateTime(t: string): string {
    return dayjs(t).format("YYYY-MM-DD HH:mm")
  }

  return { formatTime, formatDate, formatDateTime }
}
