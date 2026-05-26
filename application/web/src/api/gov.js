import request from '@/utils/request'

export function govGetHistory({ certId, page = 1, limit = 10, startTime, endTime }) {
  const params = { page, limit }
  if (startTime) params.startTime = startTime
  if (endTime) params.endTime = endTime
  return request({
    url: `/gov/cert/${certId}/history`,
    method: 'get',
    params
  })
}

export function govExportEvidence(certId) {
  return `/gov/cert/${certId}/evidence`
}

export function govExportReport(certId) {
  return `/gov/reports/audit?certId=${certId}`
}

export function govVerifyEvidence(data) {
  return request({
    url: '/gov/verify/evidence',
    method: 'post',
    data
  })
}

export function govGetStats() {
  return request({
    url: '/gov/stats',
    method: 'get'
  })
}
