import request from '@/utils/request'

export function getEnterpriseStats() {
  return request({
    url: '/enterpriseStats',
    method: 'post'
  })
}
// 企业端证据核验
export function enterpriseVerifyEvidence(data) {
  return request({
    url: '/enterprise/verify',
    method: 'post',
    data
  })
}
// 供应商关联
export function addSupplierLink(data) {
  return request({
    url: '/enterprise/supplier/add',
    method: 'post',
    data
  })
}
export function getSupplierLinks(traceabilityCode) {
  return request({
    url: `/enterprise/supplier/list?traceability_code=${traceabilityCode}`,
    method: 'get'
  })
}

// 合规事件
export function addComplianceEvent(data) {
  return request({
    url: '/enterprise/event/add',
    method: 'post',
    data
  })
}
export function getComplianceEvents(traceabilityCode) {
  return request({
    url: `/enterprise/event/list?traceability_code=${traceabilityCode}`,
    method: 'get'
  })
}
// 验证链上数据（根据交易哈希）
export function verifyChainData(txid) {
  return request({
    url: '/enterprise/verify/chaindata',
    method: 'post',
    data: { txid }
  })
}
