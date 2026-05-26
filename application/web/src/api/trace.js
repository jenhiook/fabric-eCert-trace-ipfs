import request from '@/utils/request'

export function uplink(data) {
  return request({
    url: '/uplink',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function getFruitInfo(data) {
  return request({
    url: '/getFruitInfo',
    method: 'post',
    data
  })
}

export function getAllFruitInfo(data) {
  return request({
    url: '/getAllFruitInfo',
    method: 'post',
    data
  })
}

export function getStats() {
  return request({
    url: '/getStats',
    method: 'post'
  })
}

export function getFruitList() {
  return request({
    url: '/getFruitList',
    method: 'post'
  })
}

// ========== 新增：获取证照历史 ==========
export function getHistory(certId) {
  return request({
    url: `/cert/${certId}/history`,
    method: 'get'
  })
}

// 获取当前用户的证照统计（个人用户专用）
export function getUserStats() {
  return request({
    url: '/userStats',
    method: 'post'
  })
}

export function govtAudit(data) {
  return request({
    url: '/govtAudit',
    method: 'post',
    headers: { 'Content-Type': 'multipart/form-data' },
    data
  })
}
