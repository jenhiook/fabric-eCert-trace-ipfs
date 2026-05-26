import request from '@/utils/request'

export function getTechStats() {
  return request({
    url: '/techStats',
    method: 'post'
  })
}
